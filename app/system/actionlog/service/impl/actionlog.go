package impl

import (
	"context"
	"fmt"
	"github.com/tangvis/erp/app/system/actionlog/define"
	"github.com/tangvis/erp/app/system/actionlog/repository"
	"github.com/tangvis/erp/app/system/actionlog/service"
	ctxUtil "github.com/tangvis/erp/pkg/context"
	logutil "github.com/tangvis/erp/pkg/log"
	"reflect"
	"strings"
)

type ServiceActionLog struct {
	repo repository.Repo
}

func (s ServiceActionLog) List(ctx context.Context, req *define.ListRequest) ([]define.ActionLog, error) {
	result, err := s.repo.List(ctx, repository.ListQuery{
		ModuleID: req.ModuleID,
		BizID:    req.BizID,
		Offset:   req.Offset,
		Limit:    req.Count,
	})
	if err != nil {
		return nil, err
	}
	ret := make([]define.ActionLog, len(result))
	for i, v := range result {
		ret[i] = define.ActionLog{
			ID:       v.ID,
			ModuleID: v.ModuleID,
			BizID:    v.BizID,
			Action:   v.ActionType.String(),
			Operator: v.Operator,
			Content:  v.Content,
			Ctime:    v.Ctime,
		}
	}

	return ret, nil
}

func (s ServiceActionLog) AsyncCreate(ctx context.Context,
	operator string,
	moduleID define.Module,
	bizID uint64,
	action define.Action,
	before, after any,
) {
	go func() {
		neCtx := ctxUtil.ForkContext(ctx)
		if err := s.Create(neCtx, operator, moduleID, bizID, action, before, after); err != nil {
			logutil.CtxError(neCtx, err.Error())
		}
	}()
}

func (s ServiceActionLog) Create(ctx context.Context,
	operator string,
	moduleID define.Module,
	bizID uint64,
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
	case define.UPDATE:
		content, err := Compare(before, after)
		if err != nil {
			return err
		}
		if len(content) == 0 {
			return nil
		}
		tab.Content = updateDetail(content)
	default:

	}
	return s.repo.Save(ctx, tab)
}

func NewActionLogAPP(
	repo repository.Repo,
) service.APP {
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

		beforeStr := stringValue(beforeVal.Field(i).Interface())
		afterStr := stringValue(afterVal.Field(i).Interface())

		if beforeStr != afterStr {
			changeDescription := fmt.Sprintf("[%v] has been changed to [%v]", beforeStr, afterStr)
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
