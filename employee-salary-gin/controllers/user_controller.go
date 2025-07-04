package controllers

import (
	"employee-salary-gin/config"
	"employee-salary-gin/models"
	"employee-salary-gin/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"net/http"
	"strconv"
	"time"
)

// 登录
func Login(c *gin.Context) {
	var req struct {
		Username string `form:"username" json:"username" binding:"required"`
		Password string `form:"password" json:"password" binding:"required"`
	}

	if err := c.ShouldBind(&req); err != nil {
		// 判断是否是绑定错误
		if _, ok := err.(validator.ValidationErrors); ok {
			c.JSON(http.StatusOK, gin.H{"message": "用户名或密码不能为空", "code": 400})
			return
		}
		// 其他类型的错误
		c.JSON(http.StatusBadRequest, gin.H{"message": "请求数据有误", "code": 400})
		return
	}

	var user models.User

	if err := config.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"message": "该用户名不存在", "code": 401})
		return
	}

	if user.Password != req.Password {
		c.JSON(http.StatusOK, gin.H{"message": "密码错误", "code": 401})
		return
	}

	// 登录成功后生成 token（比如 UUID）
	token := uuid.New().String()

	// 存入 Redis，key=token，value=userID
	utils.RedisClient.Set(c, token, user.ID, 6*time.Hour)
	userData := gin.H{
		"id":       user.ID,
		"username": user.Username,
		"name":     user.Name,
		"type":     user.Type, // 职位类型
	}

	c.JSON(http.StatusOK, gin.H{"message": "登录成功", "token": token, "code": 200, "data": userData})
}

// 注册
func Register(c *gin.Context) {
	var req struct {
		Username string `form:"username" json:"username" binding:"required"`
		Name     string `form:"name" json:"name" binding:"required"`
		Password string `form:"password" json:"password" binding:"required"`
		FID      uint   `form:"f_id" json:"f_id" binding:"omitempty"`
		Type     int    `form:"type" json:"type" binding:"omitempty"`
	}

	if err := c.ShouldBind(&req); err != nil {
		// 判断是否是验证错误
		if _, ok := err.(validator.ValidationErrors); ok {
			c.JSON(http.StatusOK, gin.H{"message": "用户名、密码、真实姓名不能为空", "code": 400})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"message": "请求数据有误", "code": 400})
		return
	}

	var user models.User
	if err := config.DB.Where("username = ?", req.Username).First(&user).Error; err == nil {
		c.JSON(http.StatusOK, gin.H{"message": "该用户名已存在,换一个试试", "code": 400})
		return
	}

	// 设置默认值
	if req.FID == 0 {
		req.FID = 0 // 默认值
	}
	if req.Type == 0 {
		req.Type = 0 // 默认值
	}

	now := time.Now().Unix()

	newUser := models.User{
		Username:   req.Username,
		Password:   req.Password,
		Name:       req.Name,
		FID:        req.FID,
		Type:       req.Type,
		Status:     1,
		InsertTime: now,
		UpdateTime: now,
	}

	if err := config.DB.Create(&newUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "注册失败，请稍后再试", "code": 500})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "注册成功", "code": 200})
}

// UserList 用户列表
func UserList(c *gin.Context) {

	// 获取 user_id（由 AuthMiddleware 设置）
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "用户身份验证失败", "code": 401})
		return
	}

	userIDStr, ok := userIDInterface.(string)
	if !ok {
		c.JSON(http.StatusOK, gin.H{"message": "用户ID类型错误", "code": 500})
		return
	}
	userID, _ := strconv.Atoi(userIDStr)

	fmt.Println("user_id:", userID)

	//查询用户类型
	var user models.User
	if err := config.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"message": "用户不存在", "code": 500})
		return
	}

	var req struct {
		PageNum  int    `form:"pageNum" json:"pageNum" binding:"required,min=1"`
		PageSize int    `form:"pageSize" json:"pageSize" binding:"required,min=1,max=1000"`
		Name     string `form:"name" json:"name" binding:"omitempty"` // 可选参数
	}

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "请求参数错误：" + err.Error(),
		})
		return
	}

	db := config.DB.Model(&models.User{})

	// 如果有 name 参数，则添加模糊查询条件
	if req.Name != "" {
		db = db.Where("user.name LIKE ?", "%"+req.Name+"%")
	}

	if user.ID != 1 {
		if user.Type == 1 {
			db = db.Where("user.f_id = ?", user.FID)
		} else {
			db = db.Where("user.id = ?", user.ID)
		}

	}

	// 联查 employee 表获取 name
	type UserWithFactory struct {
		models.User
		FName string `json:"f_name"`
	}

	// 查询总记录数
	var total int64
	if err := db.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "查询失败",
		})
		return
	}

	// 分页查询
	var results []UserWithFactory
	offset := (req.PageNum - 1) * req.PageSize
	if err := db.Joins("LEFT JOIN factory ON user.f_id = factory.id").
		Select("user.*, factory.name as f_name").
		Offset(offset).Limit(req.PageSize).Find(&results).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "查询失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "查询用户列表成功",
		"data": gin.H{
			"UserData": results,
			"total":    total,
		},
	})
}

func UsersByFactory(c *gin.Context) {
	// 从上下文中获取 user_id（由 AuthMiddleware 设置）
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "用户身份验证失败", "code": 401})
		return
	}

	userIDStr, ok := userIDInterface.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "用户ID类型错误", "code": 500})
		return
	}

	userID, _ := strconv.Atoi(userIDStr)
	//获取用户信息
	user := models.User{}
	err := config.DB.Where("id = ?", userID).First(&user).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"message": "用户不存在", "code": 500})
	}

	if user.Type != 0 && user.Type != 1 {
		c.JSON(http.StatusOK, gin.H{"message": "无任何员工", "code": 403})
		return
	}

	type UserSimple struct {
		ID   uint   `json:"id"`
		Name string `json:"name"`
	}
	var users []UserSimple

	if user.Type == 0 {
		//为管理员 获取所有
		_ = config.DB.Model(&models.User{}).Where("type = ?", 2).Find(&users)
	}

	if user.Type == 1 {
		//为厂长 获取该厂长的员工
		_ = config.DB.Model(&models.User{}).Where("f_id = ? and type = ?", user.FID, 2).Find(&users)
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "获取成功",
		"code":      200,
		"UsersData": users,
	})

}

func DetailsUser(c *gin.Context) {
	// 获取要获取的用户 ID
	var req struct {
		ID uint `form:"id" json:"id" binding:"required"`
	}

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "无效的请求数据",
			"code":    400,
		})
		return
	}

	id := req.ID

	// 查询服装信息
	var user models.User
	if err := config.DB.Where("id = ?", id).First(&user).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "该用户不存在",
			"code":    404,
		})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"message": "查询用户详情成功",
		"code":    200,
		"data":    user,
	})
}

func UpdateUser(c *gin.Context) {
	// 获取 user_id（由 AuthMiddleware 设置）
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "用户身份验证失败", "code": 401})
		return
	}

	userIDStr, ok := userIDInterface.(string)
	if !ok {
		c.JSON(http.StatusOK, gin.H{"message": "用户ID类型错误", "code": 500})
		return
	}
	_, _ = strconv.Atoi(userIDStr)

	var req struct {
		ID       uint   `form:"id" json:"id" binding:"required"`
		Name     string `form:"name" json:"name" binding:"required"`
		UserName string `form:"UserName" json:"UserName" binding:"required"`
		Password string `form:"password" json:"password" binding:"required"`
		FID      uint   `form:"f_id" json:"f_id" binding:"omitempty"`
		Type     int    `form:"type" json:"type" binding:"omitempty"`
		Status   int    `form:"status" json:"status" binding:"omitempty"`
	}

	if err := c.ShouldBind(&req); err != nil {
		fmt.Println("req:", req)
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求数据有误", "detail": err.Error()})
		return
	}

	now := time.Now().Unix()

	var user models.User
	if err := config.DB.Where("id = ?", req.ID).First(&user).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"error": "用户不存在", "code": 1001})
		return
	}
	if req.FID == 0 {
		req.Type = 0
	}

	user.Name = req.Name
	user.Username = req.UserName
	user.Password = req.Password
	user.FID = req.FID
	user.Type = req.Type
	user.Status = req.Status
	user.UpdateTime = now

	if err := config.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"message": "更新用户信息失败", "code": 500})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "用户信息更新成功",
		"code":    200,
	})
}

func DeleteUser(c *gin.Context) {
	// 获取 user_id（由 AuthMiddleware 设置）
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "用户身份验证失败", "code": 401})
		return
	}

	userIDStr, ok := userIDInterface.(string)
	if !ok {
		c.JSON(http.StatusOK, gin.H{"message": "用户ID类型错误", "code": 500})
		return
	}
	userID, _ := strconv.Atoi(userIDStr)
	//查询用户类型
	var user models.User
	if err := config.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"message": "用户不存在", "code": 500})
		return
	}

	if user.Type != 0 && user.Type != 1 {
		c.JSON(http.StatusOK, gin.H{"message": "无权限删除", "code": 403})
		return
	}

	// 获取要删除的服装 ID
	var req struct {
		ID uint `form:"id" json:"id" binding:"required"`
	}

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "无效的请求数据",
			"code":    400,
		})
		return
	}

	id := req.ID

	// 查询该用户是否存在
	var user1 models.User
	if err := config.DB.Where("id = ?", id).First(&user1).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"message": "该用户不存在", "code": 404})
		return
	}

	// 执行删除
	if err := config.DB.Delete(&user1).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"message": "删除失败，请重试", "code": 500})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"message": "删除用户成功",
		"code":    200,
	})
}
