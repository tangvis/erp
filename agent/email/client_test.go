package email

import (
	"context"
	"testing"
)

func TestSendHTML(t *testing.T) {
	err := NewDefaultClient().SendTxtMail(
		context.TODO(),
		[]string{"shaofeng.lin@google.com"},
		"test",
		"wwoawewqehqw",
	)
	if err != nil {
		t.Fatalf("must nil, but got %s", err)
	}
}

func TestSendMailWithFilepath(t *testing.T) {
	err := NewDefaultClient().SendMailWithFilepath(
		context.TODO(),
		[]string{"shaofeng.lin@google.com"},
		"test",
		"wwoawewqehqw",
		"",
		"./email_test.txt",
	)
	if err != nil {
		t.Fatalf("must nil, but got %s", err)
	}
}

type Dog struct {
	Name string `csv:"name"`
	Age  int    `csv:"age"`
}
