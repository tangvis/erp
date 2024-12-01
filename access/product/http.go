package product

import (
	"net/http"

	"github.com/tangvis/erp/app/product/define"
	"github.com/tangvis/erp/app/product/service"
	"github.com/tangvis/erp/common"
	"github.com/tangvis/erp/middleware/engine"
)

type Controller struct {
	engine   engine.HTTPEngine
	cateAPP  service.Category
	brandAPP service.Brand
}

func NewController(
	engine engine.HTTPEngine,
	cateAPP service.Category,
	brandAPP service.Brand,
) *Controller {
	return &Controller{
		engine:   engine,
		cateAPP:  cateAPP,
		brandAPP: brandAPP,
	}
}

func (c *Controller) URLPatterns() []engine.Router {
	return []engine.Router{
		// category
		engine.NewRouter(http.MethodPost, "/category/add", c.engine.JSONAuth(c.CateAdd)),
		engine.NewRouter(http.MethodPost, "/category/update", c.engine.JSONAuth(c.CateUpdate)),
		engine.NewRouter(http.MethodGet, "/category/list", c.engine.JSONAuth(c.CateList)),
		engine.NewRouter(http.MethodPost, "/category/delete", c.engine.JSONAuth(c.CateRemove)),

		// brand
		engine.NewRouter(http.MethodPost, "/brand/add", c.engine.JSONAuth(c.BrandAdd)),
		engine.NewRouter(http.MethodPost, "/brand/update", c.engine.JSONAuth(c.BrandUpdate)),
		engine.NewRouter(http.MethodGet, "/brand/list", c.engine.JSONAuth(c.BrandList)),
		engine.NewRouter(http.MethodPost, "/brand/delete", c.engine.JSONAuth(c.BrandRemove)),
	}
}

func (c *Controller) CateAdd(ctx engine.Context, userInfo *common.UserInfo) (any, error) {
	var req define.AddCateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		return nil, err
	}
	return c.cateAPP.Add(ctx.GetCtx(), userInfo, &req)
}

func (c *Controller) CateList(ctx engine.Context, userInfo *common.UserInfo) (any, error) {
	return c.cateAPP.List(ctx.GetCtx(), userInfo)
}

func (c *Controller) CateUpdate(ctx engine.Context, userInfo *common.UserInfo) (any, error) {
	var req define.UpdateCateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		return nil, err
	}
	return c.cateAPP.Update(ctx.GetCtx(), userInfo, &req)
}

func (c *Controller) CateRemove(ctx engine.Context, userInfo *common.UserInfo) (any, error) {
	var req define.RemoveRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		return nil, err
	}
	return nil, c.cateAPP.Remove(ctx.GetCtx(), userInfo, req.IDs...)
}

func (c *Controller) BrandAdd(ctx engine.Context, userInfo *common.UserInfo) (any, error) {
	var req define.AddBrandRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		return nil, err
	}
	return c.brandAPP.Add(ctx.GetCtx(), userInfo, &req)
}

func (c *Controller) BrandList(ctx engine.Context, userInfo *common.UserInfo) (any, error) {
	var req define.ListBrandRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		return nil, err
	}
	return c.brandAPP.List(ctx.GetCtx(), &req, userInfo)
}

func (c *Controller) BrandUpdate(ctx engine.Context, userInfo *common.UserInfo) (any, error) {
	var req define.UpdateBrandRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		return nil, err
	}
	return c.brandAPP.Update(ctx.GetCtx(), userInfo, &req)
}

func (c *Controller) BrandRemove(ctx engine.Context, userInfo *common.UserInfo) (any, error) {
	var req define.RemoveRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		return nil, err
	}
	return nil, c.brandAPP.Remove(ctx.GetCtx(), userInfo, req.IDs...)
}
