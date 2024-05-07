package gfs3

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
)

const (
	TEST_UPLOAD_FILE   = "./upload.data"
	TEST_DOWNLOAD_FILE = "./download.data"
	TEST_S3_DIR        = "test-dir/"
	TEST_BUCKET        = "test-bucket"
)

var (
	s3Conf = AwsS3Config{
		AccessKey: "minioadmin",
		SecretKey: "minioadmin",
		Endpoint:  "http://127.0.0.1:9000",
		Region:    "ap-beijing",
	}
	bgCtx = context.Background()
)

func calcS3Path() string {
	arr := strings.Split(TEST_UPLOAD_FILE, "/")
	fileName := arr[len(arr)-1]
	return TEST_S3_DIR + fileName
}

func TestS3Mgr_CreateBucket(t *testing.T) {
	ctx := bgCtx
	s3Mgr := NewS3Mgr(s3Conf)
	if err := s3Mgr.CreateBucket(ctx, TEST_BUCKET); err != nil {
		panic(err)
	}
	fmt.Println("test s3 create bucket ok")
}

func TestS3Mgr_Upload(t *testing.T) {
	ctx := bgCtx
	s3Mgr := NewS3Mgr(s3Conf)
	file, err := os.Open(TEST_UPLOAD_FILE)
	if err != nil {
		panic(err)
	}
	url, err := s3Mgr.Upload(ctx, TEST_BUCKET, calcS3Path(), file)
	if err != nil {
		panic(err)
	}
	fmt.Println("test s3 upload ok, url:", url)
}

func TestS3Mgr_Download(t *testing.T) {
	ctx := bgCtx
	s3Mgr := NewS3Mgr(s3Conf)
	out, err := os.Create(TEST_DOWNLOAD_FILE)
	if err != nil {
		panic(err)
	}
	url, err := s3Mgr.Download(ctx, TEST_BUCKET, calcS3Path(), out)
	if err != nil {
		panic(err)
	}
	fmt.Println("test s3 download ok, url:", url)
}

func TestS3Mgr_ListObjectKeysPages(t *testing.T) {
	ctx := bgCtx
	s3Mgr := NewS3Mgr(s3Conf)
	objPages, err := s3Mgr.ListObjectKeysPages(ctx, TEST_BUCKET, TEST_S3_DIR, 200)
	if err != nil {
		panic(err)
	}
	if len(objPages) == 0 {
		fmt.Println("test s3 list objs by page got empty pages")
		return
	}
	for i, objs := range objPages {
		fmt.Printf("page-%d objs: %+v\n", i, objs)
	}
	fmt.Println("test s3 list objs by page ok")
}

func TestS3Mgr_ExpireBucket(t *testing.T) {
	ctx := bgCtx
	s3Mgr := NewS3Mgr(s3Conf)
	if err := s3Mgr.ExpireBucket(ctx, TEST_BUCKET, 1); err != nil {
		panic(err)
	}
	fmt.Println("test s3 expire bucket ok")
}

func TestS3Mgr_TruncateDir(t *testing.T) {
	ctx := bgCtx
	s3Mgr := NewS3Mgr(s3Conf)
	if err := s3Mgr.TruncateDir(ctx, TEST_BUCKET, TEST_S3_DIR); err != nil {
		panic(err)
	}
	fmt.Println("test s3 truncate dir ok")
}

func TestS3Mgr_DeleteObjects(t *testing.T) {
	ctx := bgCtx
	s3Mgr := NewS3Mgr(s3Conf)
	if err := s3Mgr.DeleteObjects(ctx, TEST_BUCKET, []string{calcS3Path()}); err != nil {
		panic(err)
	}
	fmt.Println("test s3 delete objects ok")
}

func TestS3Mgr_DeleteBucket(t *testing.T) {
	ctx := bgCtx
	s3Mgr := NewS3Mgr(s3Conf)
	if err := s3Mgr.DeleteBucket(ctx, TEST_BUCKET); err != nil {
		panic(err)
	}
	fmt.Println("test s3 delete bucket ok")
}
