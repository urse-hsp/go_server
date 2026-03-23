package controller

import (
	"go-demo-server/dto"
	"go-demo-server/model"
	"go-demo-server/service"
	"go-demo-server/utils"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
)

func Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBind(&req); err != nil {
		// 这里会捕获到 JSON 解析错误 或 字段缺失错误
		utils.BadRequest(c, err.Error())
		return
	}

	user, err := service.Login(req.Username, req.Password)
	if err != nil {
		// service 返回错误
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, user)
}

func Login1(c *gin.Context) {
	type LoginRequest2 struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var req LoginRequest2
	// 1. 先判断 JSON 解析是否成功
	if err := c.ShouldBindJSON(&req); err != nil {
		// 如果 JSON 格式不对（比如传了空字符串、格式错误的 JSON）
		utils.BadRequest(c, "请求参数格式错误")
		return
	}
	// 2. 再判断字段内容是否为空（业务逻辑校验）
	if req.Username == "" || req.Password == "" {
		utils.BadRequest(c, "用户名或密码不能为空")
		return
	}

	// c.PostForm("username")
	// c.PostForm("password")

	user, err := service.Login(req.Username, req.Password)
	if err == nil {
		// service 返回错误
		utils.Unauthorized(c)
		return
	}

	utils.Success(c, user)
}

func Create(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBind(&req); err != nil {
		// 这里会捕获到 JSON 解析错误 或 字段缺失错误
		utils.BadRequest(c, err.Error())
		return
	}

	user, err := service.Create(req.Username, req.Password)
	if err != nil {
		// service 返回错误
		utils.BadRequest(c, err.Error())
		return
	}
	utils.Created(c, user)
}

func GetUserDetail(c *gin.Context) {
	userID := c.Query("id")

	// 1. 将字符串转换为整数
	id, err := strconv.Atoi(userID)
	if err != nil {
		// 2. 必须处理转换可能失败的错误
		// 例如，如果 userID 是 "abc"，转换就会失败
		log.Println("无效的 user ID:", err)
		return
	}

	user, _ := service.GetUserDetail(uint(id))
	utils.Success(c, user)
}

func GetUserInfo(c *gin.Context) {
	userId := utils.GetUserID(c)

	user, _ := service.GetUserDetail(userId)
	utils.Success(c, user)
}

func UpdateUser(c *gin.Context) {
	var req dto.UserUpdateRequest
	if err := c.ShouldBind(&req); err != nil {
		// 这里会捕获到 JSON 解析错误 或 字段缺失错误
		utils.BadRequest(c, err.Error())
		return
	}
	var reqModel model.User

	// 1. 自动复制所有同名字段 (Username, Avatar 等)
	if err := copier.Copy(&reqModel, &req); err != nil {
		utils.BadRequest(c, "数据转换错误")
		return
	}

	// id := uint(reqModel.ID)
	// 2. 【关键】强制覆盖 ID，防止前端传参篡改
	// reqModel.ID = uint(0)
	// reqModel.UpdatedAt = time.Now()

	userId := utils.GetUserID(c)

	// 3. 调用 Service
	user, err := service.UpdateUser(reqModel, userId)
	if err != nil {
		// service 返回错误
		utils.BadRequest(c, err.Error())
		return
	}
	utils.Success(c, user)
}

func DeleteUser(c *gin.Context) {
	// 1. 从 URL 路径中获取 id
	idStr := c.Param("id")

	// 2. 将字符串转换为 uint（或 int）
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的用户ID")
		return
	}

	// 3. 调用 Service/DAO 删除用户
	err = service.DeleteUser(uint(id))
	if err != nil {
		// c.JSON(500, gin.H{"error": "删除失败"})
		utils.BadRequest(c, err.Error())
		return
	}
	utils.Success(c)
}

func GetUserList(c *gin.Context) {
	user, err := service.GetUserList()
	if err != nil {
		// service 返回错误
		utils.BadRequest(c, err.Error())
		return
	}
	utils.Success(c, user)
}
