package service

import (
	"context"
	"fmt"
	"go-demo-server/bootstrap"
	"go-demo-server/dao"
	"go-demo-server/dto"
	"go-demo-server/model"
	"go-demo-server/pkg/httpclient"
	"go-demo-server/utils"
	"time"
)

type userService interface {
	Login(username string, password string) (dto.LoginResponse, error)
	CreateUser(username string, password string) (dto.UserDTO, error)
	GetUserDetail(id uint) (dto.UserDTO, error)
	UpdateUser(info model.User, id uint) (dto.UserDTO, error)
}

func Login(username string, password string) (dto.LoginResponse, error) {
	// fmt.Println(username, password)
	// fmt.Printf("用户名=%s,密码=%s", username, password)

	// if username == "admin" && password == "123456" {
	// 	return LoginResponse{
	// 		User: model.User{
	// 			ID:   1,
	// 			Name: username,
	// 		},
	// 		Token: "token",
	// 	}
	// }

	// 1. 获取用户
	user, err := dao.GetSeatchInfo(username)
	if err != nil {
		// 如果 DAO 报错（比如“用户不存在”或“数据库挂了”），直接透传给 Controller
		// 不需要再写 if user == nil 了，因为如果有 err，user 肯定是 nil
		return dto.LoginResponse{}, err
	}

	// 2. 代码能走到这里，说明 user 肯定是有值的！
	// 校验密码
	if !utils.CheckPassword(password, user.Password) {
		return dto.LoginResponse{}, fmt.Errorf("密码错误")
	}

	token := utils.GenerateToken(uint(user.ID), user.Username)

	return dto.LoginResponse{
		User:  dto.ToUserDTO(user),
		Token: token,
	}, nil
}

func Create(username string, password string) (dto.UserDTO, error) {
	// 加密
	hashedPwd, err := utils.HashPassword(password)
	if err != nil {
		return dto.UserDTO{}, err
	}

	user, err := dao.CreateUser(username, hashedPwd)
	if err != nil {
		// 如果 DAO 报错（比如“用户不存在”或“数据库挂了”），直接透传给 Controller
		// 不需要再写 if user == nil 了，因为如果有 err，user 肯定是 nil
		return dto.UserDTO{}, err
	}
	return dto.ToUserDTO(user), nil
}

func GetUserDetail(id uint) (dto.UserDTO, error) {
	user, err := dao.GetUserByID(id)
	if err != nil {
		return dto.UserDTO{}, err
	}
	return dto.ToUserDTO(user), nil
}

func UpdateUser(info model.User, id uint) (dto.UserDTO, error) {
	user, err := dao.UpdateUser(info, id)
	if err != nil {
		return dto.UserDTO{}, err
	}
	return dto.ToUserDTO(user), nil
}

func DeleteUser(id uint) error {
	_, err_ := GetUserDetail(id)
	if err_ != nil {
		return fmt.Errorf("用户不存在")
	}

	err := dao.DeleteUser(id)
	if err != nil {
		return err
	}
	return nil
}

func GetUserList() (any, error) {

	var u interface{}
	err := bootstrap.GetJSON("user:1002", &u)

	// 没缓存去获取数据保存
	if err != nil {
		// 1. 创建实例
		client := httpclient.New("https://randomuser.me/api", 30*time.Second)
		// 2. 尝试调用 (这里如果不报错，说明代码没问题，是你调用处的环境问题)
		// 注意：这里只是测试编译，不一定会运行成功
		var result interface{}

		fmt.Printf("httpclient 进行请求！\n")
		// 如果这里 IDE 报错说没有 Get 方法，说明上面的第 1、2 点有问题
		err := client.Get(context.Background(), "", &result)
		if err != nil {
			return nil, err
		}
		bootstrap.SetJSON("user:1002", result, 10*time.Second)
		fmt.Printf("httpclient 请求成功！接收到的数据")
		// 方法存在，编译通过
		u = &result
	}
	return u, nil
}
