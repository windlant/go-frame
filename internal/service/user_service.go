package service

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/windlant/go-frame/internal/model"
	"github.com/windlant/go-frame/internal/repository"
)

type IUserService interface {
	GetAll(ctx context.Context) ([]*model.User, error)
	CreateBatch(ctx context.Context, users []*model.User) (int64, error)
	GetBatchByID(ctx context.Context, ids []int) ([]*model.User, error)
	GetBatchByEmail(ctx context.Context, emails []string) ([]*model.User, error)
	UpdateBatch(ctx context.Context, users []*model.User) error
	DeleteBatch(ctx context.Context, ids []int) error
}

type UserService struct {
	repo repository.IUserRepository
}

func NewUserService() IUserService {
	return &UserService{
		repo: repository.NewUserRepository(),
	}
}

// 查询所有用户
func (s *UserService) GetAll(ctx context.Context) ([]*model.User, error) {
	return s.repo.GetAll(ctx)
}

// 批量创建：可在此加入业务校验（如邮箱格式、去重等）
func (s *UserService) CreateBatch(ctx context.Context, users []*model.User) (int64, error) {
	if len(users) == 0 {
		return 0, nil
	}

	// 简单校验：非空、邮箱格式
	for _, u := range users {
		if u.Name == "" || u.Email == "" {
			return 0, gerror.New("name and email are required")
		}
	}

	return s.repo.CreateBatch(ctx, users)
}

// 批量按 ID 查询
func (s *UserService) GetBatchByID(ctx context.Context, ids []int) ([]*model.User, error) {
	if len(ids) == 0 {
		return []*model.User{}, nil
	}
	return s.repo.GetBatchByID(ctx, ids)
}

// 批量按邮箱查询
func (s *UserService) GetBatchByEmail(ctx context.Context, emails []string) ([]*model.User, error) {
	if len(emails) == 0 {
		return []*model.User{}, nil
	}
	return s.repo.GetBatchByEmail(ctx, emails)
}

// 批量更新
func (s *UserService) UpdateBatch(ctx context.Context, users []*model.User) error {
	if len(users) == 0 {
		return nil
	}
	// 可在此校验 ID 是否有效
	for _, u := range users {
		if u.ID == 0 {
			return gerror.New("user ID is required for update")
		}
	}
	return s.repo.UpdateBatch(ctx, users)
}

// 批量删除
func (s *UserService) DeleteBatch(ctx context.Context, ids []int) error {
	if len(ids) == 0 {
		return nil
	}
	return s.repo.DeleteBatch(ctx, ids)
}
