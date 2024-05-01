package impl

import (
	"context"
	"fmt"
	"github.com/tangvis/erp/app/system/actionlog/define"
	"github.com/tangvis/erp/app/system/actionlog/repository"
	"github.com/tangvis/erp/app/system/actionlog/service"
	"reflect"
	"strings"
)

type ServiceActionLog struct {
	repo repository.Repo
}

func (s ServiceActionLog) Create(ctx context.Context,
	operator string,
	moduleID, bizID uint64,
	action define.Action,
	before, after any,
) error {
	tab := repository.ActionLogTab{
		ModuleID:   moduleID,
		BizID:      bizID,
		ActionType: action,
		Operator:   operator,
	}
	switch action {
	case define.Add:
		tab.Content = "创建"
	case define.Delete:
		tab.Content = "删除"
	case define.Update:
		content, err := Compare(before, after)
		if err != nil {
			return err
		}
		tab.Content = updateDetail(content)
	default:

	}
	return s.repo.Save(ctx, tab)
}

func NewActionLogAPP(
	repo repository.Repo,
) service.APPActionLog {
	return &ServiceActionLog{
		repo: repo,
	}
}

// stringValue converts an interface{} to a string, checking if it implements the S interface.
func stringValue(v any) string {
	if stringer, ok := v.(fmt.Stringer); ok {
		return stringer.String()
	}
	return fmt.Sprint(v)
}

func Compare(before, after any) (map[string]string, error) {
	if before == nil || after == nil {
		return nil, fmt.Errorf("one of the input values is nil")
	}

	beforeType := reflect.TypeOf(before)
	afterType := reflect.TypeOf(after)
	if beforeType != afterType {
		return nil, fmt.Errorf("type mismatch between before and after")
	}

	beforeVal := reflect.ValueOf(before)
	afterVal := reflect.ValueOf(after)
	changes := make(map[string]string)

	for i := 0; i < beforeVal.NumField(); i++ {
		field := beforeType.Field(i)
		alTag, hasAlTag := field.Tag.Lookup("al")
		if !hasAlTag {
			continue // Skip fields without 'al' tag
		}

		beforeFieldVal := beforeVal.Field(i).Interface()
		afterFieldVal := afterVal.Field(i).Interface()

		beforeStr := stringValue(beforeFieldVal)
		afterStr := stringValue(afterFieldVal)

		if beforeStr != afterStr {
			changeDescription := fmt.Sprintf("changed from [%v] to [%v]", beforeStr, afterStr)
			changes[alTag] = changeDescription
		}
	}

	return changes, nil
}

func updateDetail(m map[string]string) string {
	l := make([]string, 0)
	for k, v := range m {
		l = append(l, k+":"+v)
	}
	return strings.Join(l, ";\n")
}
