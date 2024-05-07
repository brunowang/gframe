package gfs3

import (
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/brunowang/gframe/gfid"
	"github.com/brunowang/gframe/gflog"
	"go.uber.org/zap"
	"io"
	"mime"
	"strings"
	"time"
)

type S3Mgr struct {
	conf   *aws.Config
	s3Conn *s3.S3
}

type AwsS3Config struct {
	AccessKey string `toml:"ak"`
	SecretKey string `toml:"sk"`
	Endpoint  string `toml:"endpoint"`
	Region    string `toml:"region"`
}

func NewS3Mgr(c AwsS3Config) *S3Mgr {
	return &S3Mgr{conf: &aws.Config{
		Credentials:      credentials.NewStaticCredentials(c.AccessKey, c.SecretKey, ""),
		Endpoint:         aws.String(c.Endpoint),
		Region:           aws.String(c.Region),
		DisableSSL:       aws.Bool(true),
		S3ForcePathStyle: aws.Bool(true),
	}}
}

func (m *S3Mgr) GetS3Conn(ctx context.Context) (*s3.S3, error) {
	if m.s3Conn == nil {
		newSession, err := session.NewSession(m.conf)
		if err != nil {
			gflog.Error(ctx, "new s3 session failed",
				zap.Any("config", m.conf), zap.Error(err))
			return nil, err
		}
		m.s3Conn = s3.New(newSession)
	}
	return m.s3Conn, nil
}

func (m *S3Mgr) Upload(ctx context.Context, bucket, key string, file io.Reader, opts ...UploadOption) (string, error) {
	options := uploadOptions{
		ACL: aws.String("public-read"),
	}
	for _, opt := range opts {
		opt(&options)
	}

	start := time.Now()
	defer func() {
		latency := time.Since(start)
		if latency > 500*time.Millisecond {
			gflog.Warn(ctx, "upload file slowly", zap.Duration("latency", latency),
				zap.String("bucket", bucket), zap.String("key", key))
		}
	}()
	newSession, err := session.NewSession(m.conf)
	if err != nil {
		gflog.Error(ctx, "upload new s3 session failed",
			zap.Any("config", m.conf), zap.Error(err))
		return "", err
	}
	uploader := s3manager.NewUploader(newSession)
	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket:             aws.String(bucket),
		Key:                aws.String(key),
		Body:               file,
		ACL:                options.ACL,
		ContentType:        aws.String(getContentTypeByFileExtension(getFileExtension(key))),
		ContentDisposition: options.ContentDisposition,
	}, func(u *s3manager.Uploader) {
		u.PartSize = 10 * 1024 * 1024 // 10MB part size
		u.LeavePartsOnError = true
		u.Concurrency = 3
	})
	if err != nil {
		gflog.Error(ctx, "upload file to s3 failed",
			zap.String("bucket", bucket), zap.String("key", key), zap.Error(err))
		return "", err
	}
	return result.Location, nil
}

func (m *S3Mgr) Download(ctx context.Context, bucket, key string, out io.WriterAt, opts ...DownloadOption) (int64, error) {
	options := downloadOptions{}
	for _, opt := range opts {
		opt(&options)
	}

	start := time.Now()
	defer func() {
		latency := time.Since(start)
		if latency > 500*time.Millisecond {
			gflog.Warn(ctx, "download file slowly", zap.Duration("latency", latency),
				zap.String("bucket", bucket), zap.String("key", key))
		}
	}()
	newSession, err := session.NewSession(m.conf)
	if err != nil {
		gflog.Error(ctx, "download new s3 session failed",
			zap.Any("config", m.conf), zap.Error(err))
		return 0, err
	}
	downloader := s3manager.NewDownloader(newSession)
	fileSize, err := downloader.Download(out, &s3.GetObjectInput{
		Bucket:              aws.String(bucket),
		Key:                 aws.String(key),
		ExpectedBucketOwner: options.ExpectedBucketOwner,
	}, func(d *s3manager.Downloader) {
		d.PartSize = 10 * 1024 * 1024 // 10MB part size
		d.Concurrency = 3
	})
	if err != nil {
		gflog.Error(ctx, "download file from s3 failed",
			zap.String("bucket", bucket), zap.String("key", key), zap.Error(err))
		return 0, err
	}
	return fileSize, nil
}

func (m *S3Mgr) CreateBucket(ctx context.Context, bucket string) error {
	s3Conn, err := m.GetS3Conn(ctx)
	if err != nil {
		return err
	}
	if _, err := s3Conn.HeadBucket(&s3.HeadBucketInput{
		Bucket: aws.String(bucket),
	}); err == nil {
		// Bucket已存在，直接return
		gflog.Warn(ctx, "create bucket but already exist",
			zap.String("bucket", bucket), zap.Error(err))
		return nil
	}
	if _, err := s3Conn.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(bucket),
		ACL:    aws.String("public-read"),
	}); err != nil {
		gflog.Error(ctx, "create bucket failed",
			zap.String("bucket", bucket), zap.Error(err))
		return err
	}
	return nil
}

func (m *S3Mgr) ListObjectKeysPages(ctx context.Context, bucket string, dir string, pageSize int64) ([][]string, error) {
	s3Conn, err := m.GetS3Conn(ctx)
	if err != nil {
		return nil, err
	}
	if _, err := s3Conn.HeadBucket(&s3.HeadBucketInput{
		Bucket: aws.String(bucket),
	}); err != nil {
		return nil, err
	}

	objPages := make([][]string, 0, 1)
	if err := s3Conn.ListObjectsPages(&s3.ListObjectsInput{
		Bucket:  aws.String(bucket),
		Prefix:  aws.String(dir),
		MaxKeys: aws.Int64(pageSize),
	}, func(page *s3.ListObjectsOutput, lastPage bool) bool {
		objs := make([]string, 0, pageSize)
		for _, obj := range page.Contents {
			if obj.Key != nil {
				objs = append(objs, *obj.Key)
			}
		}
		if len(objs) > 0 {
			objPages = append(objPages, objs)
		}
		return !lastPage
	}); err != nil {
		gflog.Error(ctx, "list s3 dir objs failed",
			zap.String("bucket", bucket), zap.Error(err))
		return nil, err
	}

	return objPages, nil
}

func (m *S3Mgr) TruncateDir(ctx context.Context, bucket string, dir string) error {
	s3Conn, err := m.GetS3Conn(ctx)
	if err != nil {
		return err
	}
	if _, err := s3Conn.HeadBucket(&s3.HeadBucketInput{
		Bucket: aws.String(bucket),
	}); err != nil {
		return err
	}

	const pageSize = 200
	objPages, err := m.ListObjectKeysPages(ctx, bucket, dir, pageSize)
	if err != nil {
		gflog.Error(ctx, "truncate s3 dir list objs failed",
			zap.String("bucket", bucket), zap.Error(err))
		return err
	}

	for _, objs := range objPages {
		if err := m.DeleteObjects(ctx, bucket, objs); err != nil {
			gflog.Error(ctx, "truncate s3 dir del objs failed",
				zap.String("bucket", bucket), zap.Error(err))
			return err
		}
	}

	return nil
}

func (m *S3Mgr) DeleteObjects(ctx context.Context, bucket string, objKeys []string) error {
	if len(objKeys) == 0 {
		return nil
	}
	s3Conn, err := m.GetS3Conn(ctx)
	if err != nil {
		return err
	}
	if _, err := s3Conn.HeadBucket(&s3.HeadBucketInput{
		Bucket: aws.String(bucket),
	}); err != nil {
		return err
	}
	objs := make([]*s3.ObjectIdentifier, 0, len(objKeys))
	for _, objKey := range objKeys {
		objs = append(objs, &s3.ObjectIdentifier{
			Key: aws.String(objKey),
		})
	}
	if _, err := s3Conn.DeleteObjects(&s3.DeleteObjectsInput{
		Bucket: aws.String(bucket),
		Delete: &s3.Delete{
			Objects: objs,
		},
	}); err != nil {
		gflog.Error(ctx, "delete objects failed",
			zap.String("bucket", bucket), zap.Any("objects", objKeys), zap.Error(err))
		return err
	}
	return nil
}

func (m *S3Mgr) DeleteBucket(ctx context.Context, bucket string) error {
	s3Conn, err := m.GetS3Conn(ctx)
	if err != nil {
		return err
	}
	if _, err := s3Conn.HeadBucket(&s3.HeadBucketInput{
		Bucket: aws.String(bucket),
	}); err != nil {
		return err
	}
	if _, err := s3Conn.DeleteBucket(&s3.DeleteBucketInput{
		Bucket: aws.String(bucket),
	}); err != nil {
		gflog.Error(ctx, "delete bucket failed",
			zap.String("bucket", bucket), zap.Error(err))
		return err
	}
	return nil
}

func (m *S3Mgr) ExpireBucket(ctx context.Context, bucket string, days int64) error {
	s3Conn, err := m.GetS3Conn(ctx)
	if err != nil {
		return err
	}
	if _, err := s3Conn.HeadBucket(&s3.HeadBucketInput{
		Bucket: aws.String(bucket),
	}); err != nil {
		return err
	}

	if _, err := s3Conn.PutBucketLifecycleConfiguration(&s3.PutBucketLifecycleConfigurationInput{
		Bucket: aws.String(bucket),
		LifecycleConfiguration: &s3.BucketLifecycleConfiguration{
			Rules: []*s3.LifecycleRule{
				{
					Expiration: &s3.LifecycleExpiration{
						Days: aws.Int64(days),
					},
					Filter: &s3.LifecycleRuleFilter{
						Prefix: aws.String("/"),
					},
					ID:     aws.String(gfid.GenID()),
					Status: aws.String("Enabled"),
				},
			},
		},
	}); err != nil {
		gflog.Error(ctx, "expire bucket failed",
			zap.String("bucket", bucket), zap.Error(err))
		return err
	}
	return nil
}

func getFileExtension(filePath string) string {
	if !strings.Contains(filePath, ".") {
		return ""
	}
	arr := strings.Split(filePath, ".")
	ext := "." + arr[len(arr)-1]
	return ext
}

func getContentTypeByFileExtension(fileExt string) string {
	contentType := mime.TypeByExtension(fileExt)
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	return contentType
}
