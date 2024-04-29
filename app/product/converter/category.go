package converter

import (
	"github.com/tangvis/erp/app/product/define"
	"github.com/tangvis/erp/app/product/repository/meta"
)

func CategoriesConvert(categories []meta.CategoryTab) []*define.Category {
	ret := make([]*define.Category, 0)
	m := make(map[uint64]*define.Category)
	for _, category := range categories {
		node := &define.Category{
			CateBrief: define.CateBrief{
				ID:   category.ID,
				Name: category.Name,
				Desc: category.Desc,
				URL:  category.URL,
				PID:  category.PID,
			},
			Path:       category.Path,
			CateStatus: category.CateStatus,
			Ctime:      category.Ctime,
			Mtime:      category.Mtime,
			Children:   make([]*define.Category, 0),
		}
		m[category.ID] = node
		if category.PID == 0 { // layer 1
			ret = append(ret, node)
		}
	}
	for _, category := range categories {
		parent, ok := m[category.PID]
		if !ok {
			continue
		}
		parent.Children = append(parent.Children, m[category.ID])
	}

	return ret
}

func CategoryConvert(cate meta.CategoryTab) *define.Category {
	return &define.Category{
		CateBrief: define.CateBrief{
			ID:   cate.ID,
			Name: cate.Name,
			Desc: cate.Desc,
			URL:  cate.URL,
			PID:  cate.PID,
		},
		Path:       cate.Path,
		CateStatus: cate.CateStatus,
		Ctime:      cate.Ctime,
		Mtime:      cate.Mtime,
	}
}
