package controller

import (
	"fmt"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/windlant/go-frame/internal/consts"
	"github.com/windlant/go-frame/internal/model"
	"github.com/windlant/go-frame/internal/service"
)

type UserController struct {
	userService  service.IUserService
	batchMaxSize int
}

func NewUserController() *UserController {
	// 从配置文件读取 batchMaxSize
	batchMaxSize := g.Cfg().MustGet(gctx.New(), "server.batchMaxSize", 1000).Int()
	if batchMaxSize <= 0 {
		batchMaxSize = 1000 // 默认值
	}

	return &UserController{
		userService:  service.NewUserService(),
		batchMaxSize: int(batchMaxSize),
	}
}

// 查询所有用户
func (c *UserController) ListUsers(r *ghttp.Request) {
	ctx := r.Context()
	users, err := c.userService.GetAll(ctx)
	if err != nil {
		writeError(r, gerror.NewCode(consts.InternalError, "failed to fetch users"))
		return
	}
	writeSuccess(r, users)
}

// 统一成功响应
func writeSuccess(r *ghttp.Request, data interface{}) {
	r.Response.WriteJson(g.Map{
		"code":    consts.CodeOK,
		"message": "success",
		"data":    data,
	})
}

// 统一错误响应
func writeError(r *ghttp.Request, err error) {
	// gerror.Code(err) 返回 int
	code := gerror.Code(err)
	if code.Code() == 0 {
		code = consts.InternalError
	}
	r.Response.WriteJsonExit(g.Map{
		"code":    code.Code(),
		"message": err.Error(),
		"data":    nil,
	})
}

// 批量创建
func (c *UserController) CreateUsers(r *ghttp.Request) {
	var users []*model.User
	if err := r.Parse(&users); err != nil {
		writeError(r, gerror.NewCode(consts.InvalidParams, "invalid json format"))
		return
	}
	if len(users) == 0 {
		writeError(r, gerror.NewCode(consts.InvalidParams, "user list is empty"))
		return
	}

	// 校验批量大小
	if len(users) > c.batchMaxSize {
		writeError(r, gerror.NewCode(consts.BatchTooLarge,
			fmt.Sprintf("batch size %d exceeds maximum limit %d", len(users), c.batchMaxSize)))
		return
	}

	ctx := r.Context()
	firstID, err := c.userService.CreateBatch(ctx, users)
	if err != nil {
		writeError(r, gerror.NewCode(consts.InternalError, err.Error()))
		return
	}

	writeSuccess(r, g.Map{
		"first_id": firstID,
		"count":    len(users),
	})
}

// 批量查询
func (c *UserController) GetUsers(r *ghttp.Request) {
	var req struct {
		IDs    []int    `json:"ids,omitempty"`
		Emails []string `json:"emails,omitempty"`
	}
	if err := r.Parse(&req); err != nil {
		writeError(r, gerror.NewCode(consts.InvalidParams, "invalid request body"))
		return
	}

	// 校验批量大小
	totalQuerySize := len(req.IDs) + len(req.Emails)
	if totalQuerySize == 0 {
		writeError(r, gerror.NewCode(consts.InvalidParams, "either 'ids' or 'emails' must be provided"))
		return
	}

	if totalQuerySize > c.batchMaxSize {
		writeError(r, gerror.NewCode(consts.BatchTooLarge,
			fmt.Sprintf("total query size %d exceeds maximum limit %d", totalQuerySize, c.batchMaxSize)))
		return
	}

	resp := g.Map{
		"by_id":    []*model.User{},
		"by_email": []*model.User{},
	}

	ctx := r.Context()

	if len(req.IDs) > 0 {
		usersByID, err := c.userService.GetBatchByID(ctx, req.IDs)
		if err != nil {
			writeError(r, gerror.NewCode(consts.InternalError, err.Error()))
			return
		}
		resp["by_id"] = usersByID
	}

	if len(req.Emails) > 0 {
		usersByEmail, err := c.userService.GetBatchByEmail(ctx, req.Emails)
		if err != nil {
			writeError(r, gerror.NewCode(consts.InternalError, err.Error()))
			return
		}
		resp["by_email"] = usersByEmail
	}

	writeSuccess(r, resp)
}

// 批量更新
func (c *UserController) UpdateUsers(r *ghttp.Request) {
	var users []*model.User
	if err := r.Parse(&users); err != nil {
		writeError(r, gerror.NewCode(consts.InvalidParams, "invalid json format"))
		return
	}
	if len(users) == 0 {
		writeError(r, gerror.NewCode(consts.InvalidParams, "user list is empty"))
		return
	}

	// 校验批量大小
	if len(users) > c.batchMaxSize {
		writeError(r, gerror.NewCode(consts.BatchTooLarge,
			fmt.Sprintf("batch size %d exceeds maximum limit %d", len(users), c.batchMaxSize)))
		return
	}

	ctx := r.Context()
	if err := c.userService.UpdateBatch(ctx, users); err != nil {
		writeError(r, gerror.NewCode(consts.InternalError, err.Error()))
		return
	}

	writeSuccess(r, g.Map{"updated": len(users)})
}

// 批量删除
func (c *UserController) DeleteUsers(r *ghttp.Request) {
	var req struct {
		IDs []int `json:"ids"`
	}
	if err := r.Parse(&req); err != nil {
		writeError(r, gerror.NewCode(consts.InvalidParams, "invalid request body"))
		return
	}
	if len(req.IDs) == 0 {
		writeError(r, gerror.NewCode(consts.InvalidParams, "missing 'ids' in request body"))
		return
	}

	// 校验批量大小
	if len(req.IDs) > c.batchMaxSize {
		writeError(r, gerror.NewCode(consts.BatchTooLarge,
			fmt.Sprintf("batch size %d exceeds maximum limit %d", len(req.IDs), c.batchMaxSize)))
		return
	}

	ctx := r.Context()
	if err := c.userService.DeleteBatch(ctx, req.IDs); err != nil {
		writeError(r, gerror.NewCode(consts.InternalError, err.Error()))
		return
	}

	writeSuccess(r, g.Map{"deleted": len(req.IDs)})
}
