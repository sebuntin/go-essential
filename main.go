package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
	"log"
	"os"
	"sebuntin/ginessential/common"
)

func main() {
	InitConfig()
	db := common.InitDb()
	defer func(db *gorm.DB) {
		log.Println("closing database connection...")
		err := db.Close()
		if err != nil {
			log.Println("close db failed, err: ", err.Error())
		}
	}(db)

	r := gin.Default()
	r = CollectRoute(r)
	port := viper.GetString("server.port")
	panic(r.Run(":" + port))
}

// InitConfig initialize configuration
func InitConfig() {
	workDir, _ := os.Getwd()
	viper.SetConfigName("application")
	viper.AddConfigPath(workDir + "/config")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}
