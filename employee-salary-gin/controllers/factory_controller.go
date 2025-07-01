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
			c.JSON(http.StatusOK, gin.H{"message": "fs编号，单价不能为空", "code": 400})
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

//func UpdateClothing(c *gin.Context) {
//	// 获取 user_id（由 AuthMiddleware 设置）
//	userIDInterface, exists := c.Get("user_id")
//	if !exists {
//		c.JSON(http.StatusUnauthorized, gin.H{"message": "用户身份验证失败", "code": 401})
//		return
//	}
//
//	userIDStr, ok := userIDInterface.(string)
//	if !ok {
//		c.JSON(http.StatusOK, gin.H{"message": "用户ID类型错误", "code": 500})
//		return
//	}
//	userID, _ := strconv.Atoi(userIDStr)
//
//	var req struct {
//		ID     uint   `form:"id" json:"id" binding:"required"`
//		Code   string `form:"code" json:"code" binding:"required"`
//		Price  string `form:"price" json:"price" binding:"required,gt=0"`
//		Source string `form:"source" json:"source,omitempty"`
//		Image  string `form:"image" json:"image" binding:"omitempty"`
//	}
//
//	if err := c.ShouldBind(&req); err != nil {
//		fmt.Println("req:", req)
//		c.JSON(http.StatusBadRequest, gin.H{"error": "请求数据有误", "detail": err.Error()})
//		return
//	}
//
//	//price, err := strconv.ParseFloat(req.Price, 64)
//	//if err != nil {
//	//	c.JSON(http.StatusBadRequest, gin.H{
//	//		"message": "单价格式不正确",
//	//		"code":    400,
//	//	})
//	//	return
//	//}
//
//	now := time.Now().Unix()
//
//	var clothing models.Clothing
//	if err := config.DB.Where("id = ?", req.ID).First(&clothing).Error; err != nil {
//		c.JSON(http.StatusOK, gin.H{"error": "服装不存在", "code": 1001})
//		return
//	}
//
//	clothing.Code = req.Code
//	clothing.Price = req.Price
//	clothing.Source = req.Source
//	clothing.UserID = uint(userID)
//	clothing.Image = req.Image
//	clothing.UpdateTime = now
//
//	if err := config.DB.Save(&clothing).Error; err != nil {
//		c.JSON(http.StatusOK, gin.H{"message": "更新服装信息失败", "code": 500})
//		return
//	}
//
//	c.JSON(http.StatusOK, gin.H{
//		"message": "服装信息更新成功",
//		"code":    200,
//	})
//}

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

//func DetailsClothing(c *gin.Context) {
//	// 获取要删除的服装 ID
//	var req struct {
//		ID uint `form:"id" json:"id" binding:"required"`
//	}
//
//	if err := c.ShouldBind(&req); err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{
//			"message": "无效的请求数据",
//			"code":    400,
//		})
//		return
//	}
//
//	id := req.ID
//
//	// 查询服装信息
//	var clothing models.Clothing
//	if err := config.DB.Where("id = ?", id).First(&clothing).Error; err != nil {
//		c.JSON(http.StatusOK, gin.H{
//			"message": "该服饰不存在",
//			"code":    404,
//		})
//		return
//	}
//
//	// 返回成功响应
//	c.JSON(http.StatusOK, gin.H{
//		"message": "查询服饰详情成功",
//		"code":    200,
//		"data":    clothing,
//	})
//}

//func SalaryClothingList(c *gin.Context) {
//	type ClothingSimple struct {
//		ID   uint   `json:"id"`
//		Code string `json:"code"`
//	}
//	var clothing []ClothingSimple
//
//	_ = config.DB.Model(&models.Clothing{}).Order("insert_time desc").Find(&clothing)
//
//	c.JSON(http.StatusOK, gin.H{
//		"code":    200,
//		"message": "查询服饰成功",
//		"data":    clothing,
//	})
//}

//func UploadClothingImage(c *gin.Context) {
//	// 获取上传的文件
//	file, err := c.FormFile("file")
//	if err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{
//			"code":    400,
//			"message": "文件上传失败",
//		})
//		return
//	}
//
//	// 获取文件扩展名
//	ext := ""
//	filenameWithExt := file.Filename
//	parts := strings.Split(filenameWithExt, ".")
//	if len(parts) > 1 {
//		ext = parts[len(parts)-1]
//	} else {
//		ext = "jpg" // 默认值或根据类型判断
//	}
//
//	// 生成新文件名：8位随机字符 + _ + 时间戳
//	rand.Seed(time.Now().UnixNano())
//	newFileName := fmt.Sprintf("%s_%d.%s", utils.RandStringBytes(16), time.Now().Unix(), ext)
//
//	// 定义保存路径（确保目录存在）
//	dst := fmt.Sprintf("../images/%s", newFileName)
//
//	// 保存文件
//	if err := c.SaveUploadedFile(file, dst); err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{
//			"code":    500,
//			"message": "无法保存文件",
//		})
//		return
//	}
//
//	// 返回文件访问路径
//	c.JSON(http.StatusOK, gin.H{
//		"code":    200,
//		"message": "上传成功",
//		"data": gin.H{
//			"filePath": "/images/" + newFileName,
//		},
//	})
//}

//func DeleteClothingImage(c *gin.Context) {
//	// 接收前端传来的 filePath 参数
//	var req struct {
//		FilePath string `json:"filePath" form:"filePath" binding:"required"`
//	}
//
//	if err := c.ShouldBind(&req); err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{
//			"code":    400,
//			"message": "请求参数错误：" + err.Error(),
//		})
//		return
//	}
//
//	// 拼接实际的文件物理路径
//	baseDir := ".."
//	filePath := baseDir + req.FilePath
//	fmt.Println("filePath:", filePath)
//
//	// 判断文件是否存在
//	if _, err := os.Stat(filePath); os.IsNotExist(err) {
//		c.JSON(http.StatusOK, gin.H{
//			"code":    404,
//			"message": "图片不存在",
//		})
//		return
//	}
//
//	// 删除文件
//	if err := os.Remove(filePath); err != nil {
//		c.JSON(http.StatusOK, gin.H{
//			"code":    500,
//			"message": "删除图片失败：" + err.Error(),
//		})
//		return
//	}
//
//	// 返回成功响应
//	c.JSON(http.StatusOK, gin.H{
//		"code":    200,
//		"message": "图片删除成功",
//	})
//}
