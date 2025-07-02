package routers

import (
	"employee-salary-gin/controllers"
	"employee-salary-gin/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	api := r.Group("/api")
	{
		api.POST("/login", controllers.Login)
		api.POST("/register", controllers.Register)
		api.POST("/clothingImageUpload", controllers.UploadClothingImage)
		api.POST("/clothingImageDelete", controllers.DeleteClothingImage)
		api.GET("/clothingDetail", controllers.DetailsClothing)
		api.GET("/factoryDetail", controllers.DetailsFactory)
		api.GET("/userDetail", controllers.DetailsUser)
		api.GET("/dailySalaryDetail", controllers.DetailsDailySalary)
		api.GET("/dailySalarySpecialDetail", controllers.DetailsDailySalarySpecial)
		api.GET("/userSalaryList", controllers.UserSalaryList)
		api.GET("/salaryClothingList", controllers.SalaryClothingList)
		api.GET("/userFactoryList", controllers.UserFactoryList)

	}

	// 需要认证的 API 单独分一个 group
	auth := api.Group("").Use(middleware.AuthMiddleware())
	{
		auth.GET("/userList", controllers.UserList)
		auth.GET("/userMonthSalaryList", controllers.UserMonthSalaryList)
		auth.GET("/clothingList", controllers.ClothingList)
		auth.GET("/factoryList", controllers.FactoryList)
		auth.POST("/clothingAdd", controllers.AddClothing)
		auth.POST("/factoryAdd", controllers.AddFactory)
		auth.POST("/clothingUpdate", controllers.UpdateClothing)
		auth.POST("/factoryUpdate", controllers.UpdateFactory)
		auth.POST("/userUpdate", controllers.UpdateUser)
		auth.POST("/clothingDelete", controllers.DeleteClothing)
		auth.POST("/userDelete", controllers.DeleteUser)
		auth.POST("/factoryDelete", controllers.DeleteFactory)
		auth.GET("/usersByFactory", controllers.UsersByFactory)
		auth.POST("/dailySalaryAdd", controllers.AddDailySalary)
		auth.POST("/dailySalarySpecialAdd", controllers.AddDailySalarySpecial)
		auth.POST("/dailySalaryDelete", controllers.DeleteDailySalary)
		auth.POST("/dailySalarySpecialDelete", controllers.DeleteDailySalarySpecial)
		auth.POST("/dailySalaryUpdate", controllers.UpdateDailySalary)
		auth.POST("/dailySalarySpecialUpdate", controllers.UpdateDailySalarySpecial)
	}

	r.Static("/images", "../images")

	return r
}
