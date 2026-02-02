package repository

import (
	"context"
	"encoding/json"

	"github.com/gogf/gf/v2/database/gredis"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/windlant/go-frame/internal/model"
)

const userCachePrefix = "user:id:"

type IUserRepository interface {
	GetAll(ctx context.Context) ([]*model.User, error)
	GetByID(ctx context.Context, id int) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)

	CreateBatch(ctx context.Context, users []*model.User) (int64, error)
	GetBatchByID(ctx context.Context, ids []int) ([]*model.User, error)
	GetBatchByEmail(ctx context.Context, emails []string) ([]*model.User, error)
	UpdateBatch(ctx context.Context, users []*model.User) error
	DeleteBatch(ctx context.Context, ids []int) error
}

type UserRepository struct{}

func NewUserRepository() IUserRepository {
	return &UserRepository{}
}

func (r *UserRepository) safeRedisClient() *gredis.Redis {
	ctx := gctx.New()

	enabled := g.Cfg().MustGet(ctx, "redis.default.enable").Bool()
	if !enabled {
		// g.Log().Debug(ctx, "Redis is disabled by config 'redis.default.enable=false'")
		return nil
	}

	addr := g.Cfg().MustGet(ctx, "redis.default.address").String()
	if addr == "" {
		g.Log().Warning(ctx, "Redis enabled but 'redis.default.address' is empty")
		return nil
	}

	return g.Redis()
}

func (r *UserRepository) getUserFromCache(ctx context.Context, redis *gredis.Redis, id int) (*model.User, bool) {
	if redis == nil {
		return nil, false
	}
	key := userCachePrefix + gconv.String(id)
	val, err := redis.Get(ctx, key)
	if err != nil || val.IsEmpty() {
		return nil, false
	}
	user := &model.User{}
	if err := json.Unmarshal([]byte(val.String()), user); err != nil {
		return nil, false
	}
	return user, true
}

func (r *UserRepository) setUsersToCache(ctx context.Context, redis *gredis.Redis, users []*model.User) {
	if redis == nil {
		return
	}
	for _, u := range users {
		key := userCachePrefix + gconv.String(u.ID)
		bytes, _ := json.Marshal(u)
		ttl := int64(300)
		_, _ = redis.Set(ctx, key, bytes, gredis.SetOption{
			TTLOption: gredis.TTLOption{
				EX: &ttl,
			},
		})
	}
}

// func (r *UserRepository) delUserCache(ctx context.Context, redis *gredis.Redis, id int) {
// 	if redis == nil {
// 		return
// 	}
// 	key := userCachePrefix + gconv.String(id)
// 	_, _ = redis.Del(ctx, key)
// }

func (r *UserRepository) delBatchUserCache(ctx context.Context, redis *gredis.Redis, ids []int) {
	if redis == nil || len(ids) == 0 {
		return
	}
	keys := make([]string, len(ids))
	for i, id := range ids {
		keys[i] = userCachePrefix + gconv.String(id)
	}
	_, _ = redis.Del(ctx, keys...)
}

func (r *UserRepository) GetAll(ctx context.Context) ([]*model.User, error) {
	var users []*model.User
	err := g.DB().Model("users").Scan(&users)
	return users, err
}

// 保留单个查询接口
func (r *UserRepository) GetByID(ctx context.Context, id int) (*model.User, error) {
	redis := r.safeRedisClient()

	// 尝试从缓存获取
	if user, ok := r.getUserFromCache(ctx, redis, id); ok {
		return user, nil
	}

	// 缓存未命中，查 DB
	var user model.User
	err := g.DB().Model("users").Where("id", id).Scan(&user)
	if err != nil {
		return nil, err
	}
	if user.ID == 0 {
		return nil, nil
	}

	// 回填缓存
	if redis != nil {
		r.setUsersToCache(ctx, redis, []*model.User{&user})
	}
	return &user, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	err := g.DB().Model("users").Where("email", email).Scan(&user)
	if err != nil {
		return nil, err
	}
	if user.ID == 0 {
		return nil, nil
	}
	return &user, nil
}

func (r *UserRepository) GetBatchByID(ctx context.Context, ids []int) ([]*model.User, error) {
	if len(ids) == 0 {
		return []*model.User{}, nil
	}

	redis := r.safeRedisClient()
	var missingIDs []int
	hitMap := make(map[int]*model.User)

	// 尝试从缓存批量获取
	if redis != nil {
		keys := make([]string, len(ids))
		for i, id := range ids {
			keys[i] = userCachePrefix + gconv.String(id)
		}
		cacheResults, err := redis.MGet(ctx, keys...)
		if err == nil {
			for idx, key := range keys {
				id := ids[idx]
				if val, exists := cacheResults[key]; exists && val != nil && !val.IsEmpty() {
					user := &model.User{}
					if err := json.Unmarshal([]byte(val.String()), user); err == nil {
						hitMap[id] = user
						continue
					}
				}
				missingIDs = append(missingIDs, id)
			}
		} else {
			// Redis 出错，全部走 DB
			missingIDs = ids
		}
	} else {
		missingIDs = ids
	}

	// 未命中的查 DB
	if len(missingIDs) > 0 {
		var dbUsers []*model.User
		err := g.DB().Model("users").WhereIn("id", missingIDs).Scan(&dbUsers)
		if err != nil {
			return nil, err
		}
		// 合并结果
		for _, u := range dbUsers {
			hitMap[u.ID] = u
		}
		// 回填缓存
		if redis != nil {
			r.setUsersToCache(ctx, redis, dbUsers)
		}
	}

	// 按原始顺序返回
	var users []*model.User
	for _, id := range ids {
		if user, ok := hitMap[id]; ok {
			users = append(users, user)
		}
	}
	return users, nil
}

func (r *UserRepository) GetBatchByEmail(ctx context.Context, emails []string) ([]*model.User, error) {
	if len(emails) == 0 {
		return []*model.User{}, nil
	}
	var users []*model.User
	err := g.DB().Model("users").WhereIn("email", emails).Scan(&users)
	return users, err
}

func (r *UserRepository) CreateBatch(ctx context.Context, users []*model.User) (int64, error) {
	if len(users) == 0 {
		return 0, nil
	}

	var data []interface{}
	for _, u := range users {
		data = append(data, g.Map{
			"name":  u.Name,
			"email": u.Email,
		})
	}

	result, err := g.DB().Model("users").Data(data).Insert()
	if err != nil {
		return 0, err
	}

	firstID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return firstID, nil
}

func (r *UserRepository) UpdateBatch(ctx context.Context, users []*model.User) error {
	if len(users) == 0 {
		return nil
	}

	db := g.DB()
	tx, err := db.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err := recover(); err != nil {
			_ = tx.Rollback()
			panic(err)
		}
	}()

	for _, user := range users {
		data := g.Map{
			"name":  user.Name,
			"email": user.Email,
		}
		_, err := tx.Model("users").Ctx(ctx).Data(data).Where("id", user.ID).Update()
		if err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	// 更新成功后，删除缓存
	redis := r.safeRedisClient()
	var ids []int
	for _, u := range users {
		ids = append(ids, u.ID)
	}
	r.delBatchUserCache(ctx, redis, ids)
	return nil
}

func (r *UserRepository) DeleteBatch(ctx context.Context, ids []int) error {
	if len(ids) == 0 {
		return nil
	}

	_, err := g.DB().Model("users").WhereIn("id", ids).Delete()
	if err != nil {
		return err
	}

	// 删除成功后，清理缓存
	redis := r.safeRedisClient()
	r.delBatchUserCache(ctx, redis, ids)
	return nil
}
