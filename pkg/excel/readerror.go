package excel

import (
	"fmt"
	"strings"
)

type ReadRowError interface {
	error
	GetRowIdx() int
	GetErr() error
	GetColumns() []string
}

type readRowError struct {
	rowIdx  int
	columns []string
	err     error
}

type ReadRowErrors []ReadRowError

func NewReadRowError(rowIdx int, columns []string, err error) ReadRowError {
	return &readRowError{
		rowIdx:  rowIdx,
		columns: columns,
		err:     err,
	}
}

func (r *readRowError) GetRowIdx() int {
	return r.rowIdx
}

func (r *readRowError) GetErr() error {
	return r.err
}

func (r *readRowError) GetColumns() []string {
	return r.columns
}

func (r *readRowError) Error() string {
	return fmt.Sprintf("Row[%d]: %s", r.rowIdx, r.err)
}

func (r *ReadRowErrors) Error() string {
	res := make([]string, len(*r))
	for i, v := range *r {
		res[i] = v.Error()
	}
	return strings.Join(res, ",\n")
}
