package filestore

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/url"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/google/uuid"

	logutil "github.com/tangvis/erp/pkg/log"
)

// FileStore s3 store wrapper.
type FileStore struct {
	sess          *session.Session
	srv           *s3.S3
	defaultBucket string
	publicRead    bool
}

// NewFileStore create new s3 filestore.
func NewFileStore(options *Options) (Client, error) {
	c := &aws.Config{
		Credentials:      credentials.NewStaticCredentials(options.AccessKey, options.SecretKey, ""),
		Endpoint:         aws.String(options.Endpoint),
		Region:           aws.String("s3"),
		S3ForcePathStyle: aws.Bool(true), // 设为true「bucket-in-url的方式访问」
	}
	sess, err := session.NewSession(c)
	if err != nil {
		return nil, err
	}

	s := s3.New(sess)

	return &FileStore{srv: s, sess: sess, defaultBucket: options.Bucket}, nil
}

func (fs *FileStore) GetDefaultBucketName() string {
	return fs.defaultBucket
}

// DeleteFile delete file.
func (fs *FileStore) DeleteFile(bucket, key string) error {
	delOpt := &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}
	_, err := fs.srv.DeleteObject(delOpt)
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
	return fs.Upload(ctx, filename, bytes.NewBuffer(data))
}

func (fs *FileStore) Upload(ctx context.Context, filename string, body io.Reader) (string, error) {
	uploader := s3manager.NewUploader(fs.sess)
	suffix := filepath.Ext(filename)
	contentType := SuffixMap[suffix]
	if contentType == "" {
		contentType = string(common.MimeRaw)
	}
	path := fmt.Sprintf("downloads/spx_in_station/%s/%s", uuid.New().String(), filename)
	uploadInput := &s3manager.UploadInput{
		Bucket: aws.String(fs.GetDefaultBucketName()),
		// br的转发应该是匹配了downloads的path进行转发，必须要要有downloads/才能正常转发到存储
		Key:                aws.String(path),
		Body:               body,
		ContentType:        aws.String(contentType),
		ContentDisposition: aws.String("attachment; filename=" + filename),
	}
	if fs.publicRead {
		uploadInput.ACL = aws.String("public-read") // br使用gcp的存储，不允许public
	}
	out, err := uploader.UploadWithContext(ctx, uploadInput)
	if err != nil {
		data = err.Error()
		return "", err
	}
	u, err := url.Parse(out.Location)
	if err != nil {
		return "", err
	}
	return u.Path, nil
}

func (fs *FileStore) DownloadFileByURL(ctx context.Context, location string) ([]byte, error) {
	downloader := s3manager.NewDownloader(fs.sess)
	buf := aws.NewWriteAtBuffer(nil)
	bucket, key, err := ParseBucketAndKeyFromURL(location)
	if err != nil {
		return nil, err
	}
	_, err = downloader.DownloadWithContext(ctx, buf, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (fs *FileStore) DeleteFileByURL(ctx context.Context, location string) error {
	bucket, key, err := ParseBucketAndKeyFromURL(location)
	if err != nil {
		return err
	}
	return fs.DeleteFile(bucket, key)
}

func (fs *FileStore) GetObject(ctx context.Context, bucket, key string) (*s3.GetObjectOutput, error) {
	return fs.srv.GetObjectWithContext(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
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
		data = err.Error()
		return "", err
	}
	// todo 前端鉴权以后可以删除替换host逻辑
	UnAuthHost := LiveUSSHttpHost
	if setting.IsNonLiveEnv() {
		UnAuthHost = NonLiveUSSHttpHost
	}
	return setHostname(out.Location, UnAuthHost)
}

func getBucketByFileName(filename string) string {
	if setting.Env() == setting.ENV_UAT {
		return "shopee-spx-uat"
	}
	if setting.Env() == setting.ENV_STAGING {
		return "shopee-spx-staging"
	}
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
