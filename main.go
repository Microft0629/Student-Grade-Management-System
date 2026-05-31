package main

import (
	"Student-Grade-Management-System/backend/api"
	"Student-Grade-Management-System/backend/config"
	"Student-Grade-Management-System/backend/service"
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	config.InitDatabase()

	err := service.LoadAllCSVData()
	if err != nil {
		println("CSV加载失败:", err.Error())
	}

	config.CreateDefaultAdmin()

	authAPI := api.NewAuthAPI()
	studentAPI := api.NewStudentAPI()
	courseAPI := api.NewCourseAPI()
	gradeAPI := api.NewGradeAPI()
	statisticsAPI := api.NewStatisticsAPI()
	gpaAPI := api.NewGpaAPI()
	backupAPI := api.NewBackupAPI()
	logAPI := api.NewLogAPI()
	excelAPI := api.NewExcelAPI()
	userAPI := api.NewUserAPI()

	err = wails.Run(&options.App{
		Title:  "Student grade Management System",
		Width:  1024,
		Height: 768,

		AssetServer: &assetserver.Options{
			Assets: assets,
		},

		Bind: []interface{}{
			authAPI,
			studentAPI,
			courseAPI,
			gradeAPI,
			statisticsAPI,
			gpaAPI,
			backupAPI,
			logAPI,
			excelAPI,
			userAPI,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
