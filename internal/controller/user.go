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

// GetUser 获取单个用户
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

// CreateUser 创建用户
func (u *UserController) CreateUser(r *ghttp.Request) {
	var input model.User
	if err := r.Parse(&input); err != nil {
		r.Response.WriteStatusExit(400, err.Error())
	}
	input.ID = nextID
	nextID++
	users = append(users, input)
	r.Response.WriteJsonExit(input)
}

// UpdateUser 更新用户
func (u *UserController) UpdateUser(r *ghttp.Request) {
	id := r.Get("id").Int()
	if id <= 0 {
		r.Response.WriteStatusExit(400, "Invalid user ID")
	}
	var input model.User
	if err := r.Parse(&input); err != nil {
		r.Response.WriteStatusExit(400, err.Error())
	}

	found := false
	for i, user := range users {
		if user.ID == id {
			input.ID = id
			users[i] = input
			r.Response.WriteJsonExit(input)
			found = true
			break
		}
	}
	if !found {
		r.Response.WriteStatusExit(404, "User not found")
	}
}

// DeleteUser 删除用户
func (u *UserController) DeleteUser(r *ghttp.Request) {
	id := r.Get("id").Int()
	if id <= 0 {
		r.Response.WriteStatusExit(400, "Invalid user ID")
	}

	found := false
	for i, user := range users {
		if user.ID == id {
			users = append(users[:i], users[i+1:]...)
			found = true
			break
		}
	}

	if !found {
		r.Response.WriteStatusExit(404, "User not found")
	}
	r.Response.WriteStatusExit(204)
}
