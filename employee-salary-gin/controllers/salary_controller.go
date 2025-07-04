package controllers

import (
	"bytes"
	"employee-salary-gin/config"
	"employee-salary-gin/models"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"io"
	"net/http"
	"sort"
	"strconv"
	"time"
)

type UnifiedSalary struct {
	ID       uint   `json:"id"`
	Date     string `json:"date"`
	Price    string `json:"price,omitempty"`
	Type     int    `json:"type"` // 标记类型：1=普通薪资，2=特殊薪资
	Quantity int    `json:"quantity"`
	Total    string `json:"total"`

	// 普通薪资字段
	Code  string `json:"code,omitempty"` // 只有 type == 1 时存在
	Image string `json:"image"`

	// 特殊薪资字段
	Name string `json:"name,omitempty"` // 只有 type == 2 时存在
}

func UserSalaryList(c *gin.Context) {
	var req struct {
		UserID string `form:"user_id" json:"user_id" binding:"required"`
		Month  string `form:"month" json:"month" binding:"required"`
	}

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "请求参数错误：" + err.Error(),
		})
		return
	}
	// 获取上海时区
	shanghaiLoc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "无法加载时区信息",
		})
		return
	}

	// 解析传入的月份为上海时间
	startDateStr := req.Month + "-01"
	startDate, err := time.ParseInLocation("2006-01-02", startDateStr, shanghaiLoc)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "日期格式解析失败",
		})
		return
	}

	// 构造下个月第一天，使用上海时间
	endDate := time.Date(startDate.Year(), startDate.Month()+1, 1, 0, 0, 0, 0, shanghaiLoc)

	type SalaryWithClothing struct {
		models.DailySalary
		Code  string `json:"code"`  // 从 clothing 表中查询到的服饰名称
		Price string `json:"price"` // 从 clothing 表中查询到的单价
		Total string `json:"total"`
		Image string `json:"image"`
	}

	var result []SalaryWithClothing

	db := config.DB.Table("daily_salary AS ds").
		Select("ds.*, c.code, c.price,c.price * ds.quantity AS total,c.image").
		Joins("LEFT JOIN clothing AS c ON ds.c_id = c.id").
		Where("ds.u_id = ? AND ds.date >= ? AND ds.date < ?", req.UserID, startDate, endDate)

	if err := db.Find(&result).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "查询失败",
		})
		return
	}

	type SalaryWithSpecial struct {
		models.DailySalarySpecial
		Total string `json:"total"`
	}
	var result1 []SalaryWithSpecial
	db1 := config.DB.Table("daily_salary_special").
		Select("*, price * quantity AS total").
		Where("u_id = ? AND date >= ? AND date < ?", req.UserID, startDate, endDate)

	if err := db1.Find(&result1).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "查询失败",
		})
		return
	}

	unifiedList := make([]UnifiedSalary, 0)

	// 添加 type=1 的数据
	for _, item := range result {
		price, err := strconv.ParseFloat(item.Price, 64)
		if err != nil {
			// 如果 price 字段不是合法数字，可以设为 0 或者跳过这条记录
			price = 0
		}
		itemTotal := price * float64(item.Quantity)
		unifiedList = append(unifiedList, UnifiedSalary{
			ID:       item.ID,
			Date:     item.Date,
			Type:     1,
			Code:     item.Code,
			Price:    item.Price,
			Image:    item.Image,
			Quantity: item.Quantity,
			Total:    fmt.Sprintf("%.1f", itemTotal),
		})
	}

	// 添加 type=2 的数据
	for _, item := range result1 {
		var price float64
		price, _ = strconv.ParseFloat(item.Price, 64)
		itemTotal := price * float64(item.Quantity)

		unifiedList = append(unifiedList, UnifiedSalary{
			ID:       item.ID,
			Date:     item.Date,
			Type:     2,
			Name:     item.Name,
			Price:    item.Price,
			Quantity: item.Quantity,
			Total:    fmt.Sprintf("%.1f", itemTotal),
		})
	}

	// 合并完成后排序
	sort.Slice(unifiedList, func(i, j int) bool {
		return unifiedList[i].Date < unifiedList[j].Date
	})

	//总薪资计算 AllTotal
	var allTotal float64
	for _, item := range unifiedList {
		totalFloat, _ := strconv.ParseFloat(item.Total, 64)
		allTotal += totalFloat
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "查询成功",
		"data": gin.H{
			"UserSalaryList": unifiedList,
			"AllTotal":       fmt.Sprintf("%.1f", allTotal), // 总薪资也格式化成字符串返回给前端
		},
	})
}

func UserMonthSalaryList(c *gin.Context) {
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
		Month string `form:"month" json:"month" binding:"required"`
		Name  string `form:"name" json:"name" binding:"omitempty"`
	}

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "请求参数错误：" + err.Error(),
		})
		return
	}

	db := config.DB.Model(&models.MonthlySalary{})

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

	// 联查 employee 表获取 用户相关信息
	type MonthSalaryWithUser struct {
		models.MonthlySalary
		Name string `json:"name"` // 从 clothing 表中查询到的服饰名称
		FID  uint   `json:"f_id"` // 从 clothing 表中查询到的单价
		Type int    `json:"type"`
	}

	var results []MonthSalaryWithUser

	if err := db.Joins("LEFT JOIN user ON monthly_salary.u_id = user.id").
		Select("monthly_salary.*, user.name, user.f_id, user.type").
		Where("monthly_salary.month = ?", req.Month).
		Find(&results).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "查询失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "查询每月薪资成功",
		"data": gin.H{
			"UserMonthSalaryList": results,
		},
	})
}

func AddDailySalary(c *gin.Context) {
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

	// userID, _ := strconv.Atoi(userIDStr)
	_, _ = strconv.Atoi(userIDStr)

	// 读取请求体
	//body, _ := io.ReadAll(c.Request.Body)
	//fmt.Println("Raw request body:", string(body))
	//c.Request.Body = io.NopCloser(bytes.NewBuffer(body)) // 重置 Body 以便后续绑定

	var req struct {
		Date       string `json:"date" binding:"required"` // 将 Date 变为 string 类型，方便后续处理
		UserID     uint   `json:"user_id" binding:"required"`
		ClothingID uint   `json:"clothing_id" binding:"required"`
		Quantity   int    `json:"quantity" binding:"required"`
	}

	if err := c.ShouldBind(&req); err != nil {
		fmt.Printf("Binding error: %+v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "参数绑定错误", "code": 400})
		return
	}

	// 获取当前 Unix 时间戳（秒）
	now := time.Now().Unix()

	dailySalary := models.DailySalary{
		Date:       req.Date,
		UserID:     req.UserID,
		ClothingID: req.ClothingID,
		Quantity:   req.Quantity,
		InsertTime: now,
		UpdateTime: now,
	}

	// 保存数据到数据库
	if err := config.DB.Create(&dailySalary).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "保存用户每日薪资失败", "code": 500})
		return
	}

	allTotal, err := GetAndSaveMonthlyTotal(req.UserID, req.Date)
	if err != nil {
		fmt.Println("计算月工资失败:", err)
	}
	fmt.Println("allTotal:", allTotal)

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"message": "添加用户每日薪资成功",
		"code":    200,
	})
}

func AddDailySalarySpecial(c *gin.Context) {
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

	// userID, _ := strconv.Atoi(userIDStr)
	_, _ = strconv.Atoi(userIDStr)

	// 读取请求体
	body, _ := io.ReadAll(c.Request.Body)
	fmt.Println("Raw request body:", string(body))
	c.Request.Body = io.NopCloser(bytes.NewBuffer(body)) // 重置 Body 以便后续绑定

	var req struct {
		Date     string `json:"date" binding:"required"` // 将 Date 变为 string 类型，方便后续处理
		UserID   uint   `json:"user_id" binding:"required"`
		Name     string `json:"name" binding:"required"`
		Price    string `json:"price" binding:"required"`
		Quantity int    `json:"quantity" binding:"required"`
	}

	if err := c.ShouldBind(&req); err != nil {
		fmt.Printf("Binding error: %+v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "参数绑定错误", "code": 400})
		return
	}

	// 获取当前 Unix 时间戳（秒）
	now := time.Now().Unix()

	dailySalarySpecial := models.DailySalarySpecial{
		Date:       req.Date,
		UserID:     req.UserID,
		Name:       req.Name,
		Price:      req.Price,
		Quantity:   req.Quantity,
		InsertTime: now,
		UpdateTime: now,
	}

	// 保存数据到数据库
	if err := config.DB.Create(&dailySalarySpecial).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "保存用户特殊每日薪资失败", "code": 500})
		return
	}

	allTotal, err := GetAndSaveMonthlyTotal(req.UserID, req.Date)
	if err != nil {
		fmt.Println("计算月工资失败:", err)
	}
	fmt.Println("allTotal:", allTotal)

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"message": "添加用户特殊每日薪资成功",
		"code":    200,
	})
}

func DeleteDailySalary(c *gin.Context) {
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
	var dailySalary models.DailySalary
	if err := config.DB.Where("id = ?", id).First(&dailySalary).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"message": "该每日薪资不存在", "code": 404})
		return
	}

	// 执行删除
	if err := config.DB.Delete(&dailySalary).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"message": "删除失败，请重试", "code": 500})
		return
	}

	allTotal, err := GetAndSaveMonthlyTotal(dailySalary.UserID, dailySalary.Date)
	if err != nil {
		fmt.Println("计算月工资失败:", err)
	}
	fmt.Println("allTotal:", allTotal)

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"message": "删除用户每日薪资成功",
		"code":    200,
	})
}

func DeleteDailySalarySpecial(c *gin.Context) {
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

	// 获取要删除的特殊 ID
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
	var dailySalarySpecial models.DailySalarySpecial
	if err := config.DB.Where("id = ?", id).First(&dailySalarySpecial).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"message": "该特殊每日薪资不存在", "code": 404})
		return
	}

	// 执行删除
	if err := config.DB.Delete(&dailySalarySpecial).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"message": "删除失败，请重试", "code": 500})
		return
	}

	allTotal, err := GetAndSaveMonthlyTotal(dailySalarySpecial.UserID, dailySalarySpecial.Date)
	if err != nil {
		fmt.Println("计算月工资失败:", err)
	}
	fmt.Println("allTotal:", allTotal)

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"message": "删除用户每日薪资成功",
		"code":    200,
	})
}

func DetailsDailySalary(c *gin.Context) {
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
	var dailySalary models.DailySalary
	if err := config.DB.Where("id = ?", id).First(&dailySalary).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "该每日薪资不存在",
			"code":    404,
		})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"message": "查询该日薪资详情成功",
		"code":    200,
		"data":    dailySalary,
	})
}

func DetailsDailySalarySpecial(c *gin.Context) {
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
	var dailySalarySpecial models.DailySalarySpecial
	if err := config.DB.Where("id = ?", id).First(&dailySalarySpecial).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "该每日薪资不存在",
			"code":    404,
		})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"message": "查询特殊薪资详情成功",
		"code":    200,
		"data":    dailySalarySpecial,
	})
}

func UpdateDailySalary(c *gin.Context) {
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
		ID         uint   `form:"id" json:"id" binding:"required"`
		Date       string `form:"date" json:"date" binding:"required"`
		UserID     uint   `form:"u_id" json:"u_id" binding:"required"`
		ClothingID uint   `form:"c_id" json:"c_id" binding:"required"`
		Quantity   int    `form:"quantity" json:"quantity" binding:"required"`
	}
	if err := c.ShouldBind(&req); err != nil {
		fmt.Println("req:", req)
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求数据有误", "detail": err.Error()})
		return
	}

	now := time.Now().Unix()

	var dailySalary models.DailySalary
	if err := config.DB.Where("id = ?", req.ID).First(&dailySalary).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"error": "每日薪资信息不存在", "code": 1001})
		return
	}

	dailySalary.Date = req.Date
	dailySalary.UserID = req.UserID
	dailySalary.ClothingID = req.ClothingID
	dailySalary.Quantity = req.Quantity

	dailySalary.UpdateTime = now

	if err := config.DB.Save(&dailySalary).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"message": "更新每日薪资信息失败", "code": 500})
		return
	}

	allTotal, err := GetAndSaveMonthlyTotal(req.UserID, req.Date)
	if err != nil {
		fmt.Println("计算月工资失败:", err)
	}
	fmt.Println("allTotal:", allTotal)

	c.JSON(http.StatusOK, gin.H{
		"message": "每日薪资更新成功",
		"code":    200,
	})
}

func UpdateDailySalarySpecial(c *gin.Context) {
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
		Date     string `form:"date" json:"date" binding:"required"`
		UserID   uint   `form:"u_id" json:"u_id" binding:"required"`
		Name     string `form:"name" json:"name" binding:"required"`
		Price    string `form:"price" json:"price" binding:"required"`
		Quantity int    `form:"quantity" json:"quantity" binding:"required"`
	}

	if err := c.ShouldBind(&req); err != nil {
		fmt.Println("req:", req)
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求数据有误", "detail": err.Error()})
		return
	}

	now := time.Now().Unix()

	var dailySalarySpecial models.DailySalarySpecial
	if err := config.DB.Where("id = ?", req.ID).First(&dailySalarySpecial).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"error": "每日特殊薪资信息不存在", "code": 1001})
		return
	}

	dailySalarySpecial.Date = req.Date
	dailySalarySpecial.UserID = req.UserID
	dailySalarySpecial.Name = req.Name
	dailySalarySpecial.Price = req.Price
	dailySalarySpecial.Quantity = req.Quantity

	dailySalarySpecial.UpdateTime = now

	if err := config.DB.Save(&dailySalarySpecial).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"message": "更新每日特殊薪资信息失败", "code": 500})
		return
	}

	allTotal, err := GetAndSaveMonthlyTotal(req.UserID, req.Date)
	if err != nil {
		fmt.Println("计算月工资失败:", err)
	}
	fmt.Println("allTotal:", allTotal)

	c.JSON(http.StatusOK, gin.H{
		"message": "每日特殊薪资更新成功",
		"code":    200,
	})
}

// CalculateMonthlyTotal 仅用于计算指定用户某月的总薪资
func CalculateMonthlyTotal(userID uint, day string) (float64, error) {
	// 解析日期字符串为 time.Time
	date, err := time.Parse("2006-01-02", day)
	if err != nil {
		return 0, fmt.Errorf("日期格式错误：%v", err)
	}

	// 获取该日期所在月份（格式：YYYY-MM）
	month := date.Format("2006-01")

	// 获取上海时区
	shanghaiLoc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return 0, fmt.Errorf("无法加载时区信息")
	}

	// 构造开始和结束时间
	startDateStr := month + "-01"
	startDate, err := time.ParseInLocation("2006-01-02", startDateStr, shanghaiLoc)
	if err != nil {
		return 0, fmt.Errorf("日期格式解析失败：%v", err)
	}
	endDate := time.Date(startDate.Year(), startDate.Month()+1, 1, 0, 0, 0, 0, shanghaiLoc)

	// 查询普通薪资总额
	var totalResult struct {
		Total string `json:"total"`
	}

	err = config.DB.Table("daily_salary AS ds").
		Select("COALESCE(SUM(c.price * ds.quantity), 0) AS total"). // 使用 COALESCE 避免 NULL
		Joins("LEFT JOIN clothing AS c ON ds.c_id = c.id").
		Where("ds.u_id = ? AND ds.date >= ? AND ds.date < ?", userID, startDate, endDate).
		Scan(&totalResult).Error

	if err != nil {
		return 0, fmt.Errorf("查询普通薪资失败: %v", err)
	}

	// 查询特殊薪资总额
	var specialTotal float64
	err = config.DB.Table("daily_salary_special").
		Select("COALESCE(SUM(price * quantity), 0) AS total"). // 使用 COALESCE 避免 NULL
		Where("u_id = ? AND date >= ? AND date < ?", userID, startDate, endDate).
		Row().
		Scan(&specialTotal)

	if err != nil && err.Error() != "sql: no rows in result set" {
		return 0, fmt.Errorf("查询特殊薪资失败: %v", err)
	}

	// 合并计算总金额
	var totalFloat float64

	// 处理普通薪资
	if totalResult.Total != "" {
		val, parseErr := strconv.ParseFloat(totalResult.Total, 64)
		if parseErr == nil {
			totalFloat += val
		} else {
			return 0, fmt.Errorf("解析普通薪资总额失败: %v", parseErr)
		}
	}

	// 特殊薪资已经是 float64，无需额外处理空值，COALESCE 已确保非空
	totalFloat += specialTotal

	return totalFloat, nil
}

// UpdateOrInsertMonthlySalary 存在则更新，不存在且 total > 0 则插入
func UpdateOrInsertMonthlySalary(userID uint, month string, total float64) error {
	var monthlySalary models.MonthlySalary
	err := config.DB.Where("month = ? AND u_id = ?", month, userID).First(&monthlySalary).Error

	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("查询月度薪资失败: %v", err)
		}

		// 如果没有记录，并且 total > 0，则插入新数据
		if total > 0 {
			newMonthlySalary := models.MonthlySalary{
				Month:  month,
				Total:  total,
				UserID: userID,
			}
			if err := config.DB.Create(&newMonthlySalary).Error; err != nil {
				return fmt.Errorf("插入月度薪资失败: %v", err)
			}
		}
	} else {
		// 如果已有记录，则更新 Total
		monthlySalary.Total = total
		if err := config.DB.Save(&monthlySalary).Error; err != nil {
			return fmt.Errorf("更新月度薪资失败: %v", err)
		}

	}

	return nil
}

// GetAndSaveMonthlyTotal 提供给外部调用的接口，封装计算和入库逻辑
func GetAndSaveMonthlyTotal(userID uint, day string) (string, error) {
	totalFloat, err := CalculateMonthlyTotal(userID, day)
	if err != nil {
		return "", err
	}

	// 解析 day 得到 month
	date, _ := time.Parse("2006-01-02", day)
	month := date.Format("2006-01")

	// 执行入库逻辑
	if err := UpdateOrInsertMonthlySalary(userID, month, totalFloat); err != nil {
		return "", err
	}

	return fmt.Sprintf("%.1f", totalFloat), nil
}
