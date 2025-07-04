package controllers

import (
	"employee-salary-gin/config"
	"employee-salary-gin/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
	"strconv"
	"time"
)

func FactoryList(c *gin.Context) {
	var req struct {
		PageNum  int    `form:"pageNum" json:"pageNum" binding:"required,min=1"`
		PageSize int    `form:"pageSize" json:"pageSize" binding:"required,min=1,max=1000"`
		Name     string `form:"name" json:"name" binding:"omitempty"`
	}

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "请求参数错误：" + err.Error(),
		})
		return
	}

	db := config.DB.Model(&models.Factory{})

	// 如果有 Code 参数，则添加模糊查询条件
	if req.Name != "" {
		db = db.Where("name LIKE ?", "%"+req.Name+"%")
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
	offset := (req.PageNum - 1) * req.PageSize
	var results []models.Factory
	if err := db.Offset(offset).Limit(req.PageSize).Order("insert_time DESC").Find(&results).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "查询失败",
		})
		return
	}
	fmt.Println(results)

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "查询工厂成功",
		"data": gin.H{
			"FactoryData": results,
			"total":       total,
		},
	})
}

func AddFactory(c *gin.Context) {
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

	_, _ = strconv.Atoi(userIDStr)

	var req struct {
		Name string `form:"name" json:"name" binding:"required"`
	}

	if err := c.ShouldBind(&req); err != nil {

		// 判断是否是验证错误
		if _, ok := err.(validator.ValidationErrors); ok {
			c.JSON(http.StatusOK, gin.H{"message": "工厂名称不能为空", "code": 400})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"message": "请求数据有误", "code": 400})
		return
	}

	now := time.Now().Unix() // 获取当前 Unix 时间戳（秒）

	factory := models.Factory{
		Name:       req.Name,
		InsertTime: now,
		UpdateTime: now,
	}

	if err := config.DB.Create(&factory).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"message": "保存工厂信息失败", "code": 500})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "添加服装信息成功",
		"code":    200,
		"data":    factory,
	})
}

func UpdateFactory(c *gin.Context) {
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
		ID   uint   `form:"id" json:"id" binding:"required"`
		Name string `form:"code" json:"name" binding:"required"`
	}

	if err := c.ShouldBind(&req); err != nil {
		fmt.Println("req:", req)
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求数据有误", "detail": err.Error()})
		return
	}

	now := time.Now().Unix()

	var factory models.Factory
	if err := config.DB.Where("id = ?", req.ID).First(&factory).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"error": "工厂不存在", "code": 1001})
		return
	}

	factory.Name = req.Name
	factory.UpdateTime = now

	if err := config.DB.Save(&factory).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"message": "更新工厂信息失败", "code": 500})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "工厂信息更新成功",
		"code":    200,
	})
}

func DeleteFactory(c *gin.Context) {
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

	// 查询该服装是否存在
	var factory models.Factory
	if err := config.DB.Where("id = ?", id).First(&factory).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"message": "该工厂不存在", "code": 404})
		return
	}

	// 执行删除
	if err := config.DB.Delete(&factory).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"message": "删除失败，请重试", "code": 500})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"message": "删除工厂成功",
		"code":    200,
	})
}

func UserFactoryList(c *gin.Context) {
	type FactorySimple struct {
		ID   uint   `json:"id"`
		Name string `json:"name"`
	}
	var factory []FactorySimple

	_ = config.DB.Model(&models.Factory{}).Order("insert_time desc").Find(&factory)

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "查询工厂成功",
		"data":    factory,
	})
}

func DetailsFactory(c *gin.Context) {
	// 获取要查询的服装 ID
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

	// 查询工厂信息
	var factory models.Factory
	if err := config.DB.Where("id = ?", id).First(&factory).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "该工厂不存在",
			"code":    404,
		})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"message": "查询工厂详情成功",
		"code":    200,
		"data":    factory,
	})
}
