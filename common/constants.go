package common

type Boolean uint8
type MIMEType string

func (b Boolean) True() bool {
	return b == T
}

const (
	UserInfoKey = "user_info"

	F Boolean = 0
	T Boolean = 1

	MimeRaw  MIMEType = "application/octet-stream"
	MimeCSV  MIMEType = "text/csv"
	MimeText MIMEType = "text/plain"
	MimePDF  MIMEType = "application/pdf"
	MimeXLS  MIMEType = "application/vnd.ms-excel"
	MimeXLSX MIMEType = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	MimeZip  MIMEType = "application/zip"
	MimeHTML MIMEType = "text/html"
)
