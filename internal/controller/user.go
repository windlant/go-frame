package controller

import (
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/windlant/go-frame/internal/model"
)

var (
	users  = []model.User{{ID: 1, Name: "Alice", Email: "alice@example.com"}}
	nextID = 2
)

type UserController struct{}

// GetUsers 获取所有用户
func (u *UserController) GetUsers(r *ghttp.Request) {
	r.Response.WriteJson(users)
}

// GetUser 根据 ID 获取单个用户
func (u *UserController) GetUser(r *ghttp.Request) {
	id := r.Get("id").Int()
	if id <= 0 {
		r.Response.WriteStatusExit(400, "Invalid user ID")
	}
	for _, user := range users {
		if user.ID == id {
			r.Response.WriteJson(user)
			return
		}
	}
	r.Response.WriteStatusExit(404, "User not found")
}

// CreateUsers 批量创建用户
func (u *UserController) CreateUsers(r *ghttp.Request) {
	var inputs []model.User
	if err := r.Parse(&inputs); err != nil {
		r.Response.WriteStatusExit(400, err.Error())
	}

	var created []model.User
	var errors []string

	for _, input := range inputs {
		if input.Name == "" || input.Email == "" {
			errors = append(errors, "Name and Email are required")
			continue
		}
		input.ID = nextID
		nextID++
		users = append(users, input)
		created = append(created, input)
	}

	r.Response.WriteJsonExit(map[string]interface{}{
		"created": created,
		"errors":  errors,
	})
}

// UpdateUsers 批量更新用户
func (u *UserController) UpdateUsers(r *ghttp.Request) {
	var inputs []model.User
	if err := r.Parse(&inputs); err != nil {
		r.Response.WriteStatusExit(400, err.Error())
	}

	var updated []model.User
	var notFound []int

	for _, input := range inputs {
		if input.ID <= 0 {
			notFound = append(notFound, input.ID)
			continue
		}

		found := false
		for i, user := range users {
			if user.ID == input.ID {
				users[i].Name = input.Name
				users[i].Email = input.Email
				updated = append(updated, users[i])
				found = true
				break
			}
		}
		if !found {
			notFound = append(notFound, input.ID)
		}
	}

	r.Response.WriteJsonExit(map[string]interface{}{
		"updated":   updated,
		"not_found": notFound,
	})
}

// DeleteUsers 批量删除用户
func (u *UserController) DeleteUsers(r *ghttp.Request) {
	var req struct {
		IDs []int `json:"ids"`
	}
	if err := r.Parse(&req); err != nil {
		r.Response.WriteStatusExit(400, err.Error())
	}

	var deleted []int
	delSet := make(map[int]bool)
	for _, id := range req.IDs {
		delSet[id] = true
	}

	toKeep := make([]model.User, 0, len(users))
	for _, u := range users {
		if delSet[u.ID] {
			deleted = append(deleted, u.ID)
		} else {
			toKeep = append(toKeep, u)
		}
	}
	users = toKeep

	foundSet := make(map[int]bool)
	for _, id := range deleted {
		foundSet[id] = true
	}
	var notFound []int
	for _, id := range req.IDs {
		if !foundSet[id] {
			notFound = append(notFound, id)
		}
	}

	r.Response.WriteJsonExit(map[string]interface{}{
		"deleted":   deleted,
		"not_found": notFound,
	})
}
