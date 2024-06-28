package filestore

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseBucketAndKeyFromURL(t *testing.T) {
	assert := require.New(t)
	_, _, err := ParseBucketAndKeyFromURL("abc")
	assert.NotNil(err)

	_, _, err = ParseBucketAndKeyFromURL("https://www.baidu.com/xixi")
	assert.NotNil(err)

	const URL = "http://s3.i.test.sz.shopee.io/shopee_srm_test/file_task/8e701390-71e2-49f3-aa03-296e9729a3a4"
	bucket, key, err := ParseBucketAndKeyFromURL(URL)
	assert.Nil(err)
	assert.Equal("shopee_srm_test", bucket)
	assert.Equal("file_task/8e701390-71e2-49f3-aa03-296e9729a3a4", key)
}

func TestFilestore(t *testing.T) {
	store, _ := NewFileStore(&Options{
		AccessKey: "", // check the apollo to get the config
		SecretKey: "",
		Bucket:    "",
		Endpoint:  "",
	})
	ctx := context.TODO()
	const testString = "test upload and download"
	location, err := store.Upload(ctx, "xxx.txt", bytes.NewBufferString(testString))
	if err != nil {
		t.Fatal(err)
	}
	res, err := store.DownloadFileByURL(ctx, location)
	if err != nil {
		t.Fatal(err)
	}
	if string(res) != testString {
		t.Fatalf("should be equal, but got %s", res)
	}

	if err = store.DeleteFileByURL(ctx, location); err != nil {
		t.Fatal(err)
	}

	_, err = store.DownloadFileByURL(ctx, location)
	t.Logf("err is %+v(%T)", err, err)
	if err == nil {
		t.Fatalf("should not be nil")
	}

}
