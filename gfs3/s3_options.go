package gfs3

import "github.com/aws/aws-sdk-go/aws"

type (
	UploadOption func(options *uploadOptions)

	uploadOptions struct {
		// The canned ACL to apply to the object. For more information, see Canned ACL
		// (https://docs.aws.amazon.com/AmazonS3/latest/dev/acl-overview.html#CannedACL).
		//
		// This action is not supported by Amazon S3 on Outposts.
		ACL *string `location:"header" locationName:"x-amz-acl" type:"string" enum:"ObjectCannedACL"`

		// Specifies presentational information for the object. For more information,
		// see http://www.w3.org/Protocols/rfc2616/rfc2616-sec19.html#sec19.5.1 (http://www.w3.org/Protocols/rfc2616/rfc2616-sec19.html#sec19.5.1).
		ContentDisposition *string `location:"header" locationName:"Content-Disposition" type:"string"`
	}

	DownloadOption func(options *downloadOptions)

	downloadOptions struct {
		// The account ID of the expected bucket owner. If the bucket is owned by a
		// different account, the request will fail with an HTTP 403 (Access Denied)
		// error.
		ExpectedBucketOwner *string `location:"header" locationName:"x-amz-expected-bucket-owner" type:"string"`
	}
)

func WithACL(acl string) UploadOption {
	return func(options *uploadOptions) {
		if acl != "" {
			options.ACL = aws.String(acl)
		}
	}
}

func WithContentDisposition(contentDisposition string) UploadOption {
	return func(options *uploadOptions) {
		if contentDisposition != "" {
			options.ContentDisposition = aws.String(contentDisposition)
		}
	}
}

func WithExpectedBucketOwner(expectedBucketOwner string) DownloadOption {
	return func(options *downloadOptions) {
		if expectedBucketOwner != "" {
			options.ExpectedBucketOwner = aws.String(expectedBucketOwner)
		}
	}
}
