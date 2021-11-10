package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"log"
	"sebuntin/ginessential/common"
)

func main() {
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
	panic(r.Run(":8008"))
}
