package filestore

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"io"
	"net/url"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/tangvis/erp/common"
	logutil "github.com/tangvis/erp/pkg/log"
)

// FileStore s3 store wrapper.
type FileStore struct {
	client        *s3.Client
	defaultBucket string
	publicRead    bool
	sess          *session.Session
}

// NewFileStore create new s3 filestore.
func NewFileStore(options *Options) (*FileStore, error) {
	// Load the AWS configuration with custom endpoint and credentials
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(options.AccessKey, options.SecretKey, "")),
		config.WithRegion(options.Region), // You may need to adjust the region as necessary
	)
	if err != nil {
		return nil, err
	}

	// Create an S3 client
	s3Client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = options.S3ForcePathStyle // Force path style
	})

	return &FileStore{
		client:        s3Client,
		defaultBucket: options.Bucket,
	}, nil
}

func (fs *FileStore) GetDefaultBucketName() string {
	return fs.defaultBucket
}

// DeleteFile delete file.
func (fs *FileStore) DeleteFile(ctx context.Context, bucket string, objectKeys []string) error {
	var objectIds []types.ObjectIdentifier
	for _, key := range objectKeys {
		objectIds = append(objectIds, types.ObjectIdentifier{Key: aws.String(key)})
	}
	output, err := fs.client.DeleteObjects(ctx, &s3.DeleteObjectsInput{
		Bucket: aws.String(bucket),
		Delete: &types.Delete{Objects: objectIds},
	})
	if err != nil {
		logutil.CtxErrorF(ctx, "Couldn't delete objects from bucket %v. Here's why: %v", bucket, err)
	} else {
		logutil.CtxInfoF(ctx, "Deleted %v objects, bucket %s, keys %+v", len(output.Deleted), bucket, objectKeys)
	}
	return err
}

// ParseBucketAndKeyFromURL parse bucket and file key from file url.
func ParseBucketAndKeyFromURL(rawURL string) (string, string, error) {
	uri, err := url.Parse(rawURL)
	if err != nil {
		return "", "", err
	}

	res := strings.SplitN(uri.Path, "/", 3)
	if len(res) != 3 {
		return "", "", errors.New("rawUrl format error")
	}
	return res[1], res[2], nil
}

func (fs *FileStore) UploadBytes(ctx context.Context, filename string, data []byte) (string, error) {
	return filename, fs.Upload(ctx, fs.GetDefaultBucketName(), filename, bytes.NewBuffer(data))
}

func (fs *FileStore) Upload(ctx context.Context, bucketName, filename string, body io.Reader) error {
	_, err := fs.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(filename),
		Body:   body,
	})
	if err != nil {
		logutil.CtxErrorF(ctx, "Couldn't upload file %v to %v:%v. Here's why: %v",
			filename, bucketName, filename, err)
	} else {
		logutil.CtxInfoF(ctx, "successfully uploaded file %v to %v:%v", filename, bucketName, filename)
	}

	return err
}

func (fs *FileStore) DownloadFileByURL(ctx context.Context, location string) ([]byte, error) {
	bucketName, objectKey, err := ParseBucketAndKeyFromURL(location)
	if err != nil {
		return nil, err
	}
	output, err := fs.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = output.Body.Close()
	}()

	// Read the object data into a byte slice
	data, err := io.ReadAll(output.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (fs *FileStore) DeleteFileByURL(ctx context.Context, url string) error {
	bucket, key, err := ParseBucketAndKeyFromURL(url)
	if err != nil {
		return err
	}
	return fs.DeleteFile(ctx, bucket, []string{key})
}

func (fs *FileStore) UploadExcelFileWithName(bucket, key string, filename string, body io.Reader) (string, error) {
	uploader := s3manager.NewUploader(fs.sess)

	uploadInput := &s3manager.UploadInput{
		Bucket:             aws.String(bucket),
		Key:                aws.String(key),
		Body:               body,
		ACL:                aws.String("public-read"),
		ContentType:        aws.String("application/vnd.ms-excel"),
		ContentDisposition: aws.String("attachment; filename=" + filename),
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	out, err := uploader.UploadWithContext(ctx, uploadInput)
	if err != nil {
		return "", err
	}

	return out.Location, nil
}

func setHostname(addr, hostname string) (string, error) {
	u, err := url.Parse(addr)
	if err != nil {
		return "", err
	}
	u.Host = hostname
	return u.String(), nil

}

func (fs *FileStore) UploadWithFileName(ctx context.Context, filename string, body io.Reader) (string, error) {
	uploader := s3manager.NewUploader(fs.sess)
	suffix := filepath.Ext(filename)
	contentType := SuffixMap[suffix]
	if contentType == "" {
		contentType = string(common.MimeRaw)
	}
	bucket := getBucketByFileName(filename)
	logutil.CtxInfoF(ctx, "upload filename %s  to bucket: %s", filename, bucket)
	uploadInput := &s3manager.UploadInput{
		Bucket:             aws.String(bucket),
		Key:                aws.String(filename),
		Body:               body,
		ACL:                aws.String("public-read"),
		ContentType:        aws.String(contentType),
		ContentDisposition: aws.String("attachment; filename=" + filename),
		Expires:            aws.Time(time.Now().Add(time.Hour * 24 * 365)),
	}
	out, err := uploader.UploadWithContext(ctx, uploadInput)
	if err != nil {
		//data := err.Error()
		return "", err
	}
	// todo 前端鉴权以后可以删除替换host逻辑
	UnAuthHost := LiveUSSHttpHost
	//if setting.IsNonLiveEnv() {
	//	UnAuthHost = NonLiveUSSHttpHost
	//}
	return setHostname(out.Location, UnAuthHost)
}

func getBucketByFileName(filename string) string {
	for _, reStr := range BucketRepList {
		re := regexp.MustCompile(reStr)
		bucketType := BucketMap[reStr]

		if re.MatchString(filename) {
			if bucketType == Permanent {
				return BucketNameMap[Permanent]
			} else if bucketType == PermanentDate {
				tempList := re.FindStringSubmatch(filename)
				// python 中如果存在错误就返回这个默认值
				if len(tempList) < 2 {
					return BucketNameMap[TempBucket]
				}
				return fmt.Sprintf("%s-%s%s", BucketNameMap[bucketType], tempList[1], tempList[2])
			}
		}
	}
	return BucketNameMap[TempBucket]
}
