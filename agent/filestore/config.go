package filestore

import "github.com/tangvis/erp/common"

type Options struct {
	AccessKey string `yaml:"AccessKey" json:"AccessKey"`
	SecretKey string `yaml:"SecretKey" json:"SecretKey"`
	//Endpoint         string `yaml:"Endpoint" json:"Endpoint"`
	Bucket           string `yaml:"Bucket" json:"Bucket"`
	Region           string `yaml:"Region" json:"Region"`
	S3ForcePathStyle bool   `yaml:"S3ForcePathStyle" json:"S3ForcePathStyle"`
	PublicRead       bool   `yaml:"PublicRead" json:"PublicRead"`
}

const (
	S3Operation        = "S3Operation"
	NonLiveUSSHttpHost = ""
	LiveUSSHttpHost    = ""
)

var (
	SuffixMap = map[string]string{
		".csv":  string(common.MimeCSV),
		".txt":  string(common.MimeText),
		".pdf":  string(common.MimePDF),
		".xls":  string(common.MimeXLS),
		".xlsx": string(common.MimeXLSX),
		".zip":  string(common.MimeZip),
		".html": string(common.MimeHTML),
	}
)

type BucketType int8

const (
	TempBucket    BucketType = 0
	PermanentDate BucketType = 1
	Permanent     BucketType = 2
)

// 控制map key的顺序
var BucketRepList = []string{
	"downloads/(\\d{4})/(\\d{2})/(\\d{1,2})/.*jpeg",
	"downloads/(\\d{4})/(\\d{2})/(\\d{1,2})/.*jpg",
	"downloads/(\\d{4})/(\\d{1,2})/(\\d{1,2})/.*",
}

var BucketMap = map[string]BucketType{
	"downloads/(\\d{4})/(\\d{2})/(\\d{1,2})/.*jpeg": PermanentDate,
	"downloads/(\\d{4})/(\\d{2})/(\\d{1,2})/.*jpg":  PermanentDate,
	"downloads/(\\d{4})/(\\d{1,2})/(\\d{1,2})/.*":   Permanent,
}

var BucketNameMap = map[BucketType]string{
	TempBucket:    "linus-temp-data",
	PermanentDate: "linus-perm",
	Permanent:     "linus-perm-data",
}
