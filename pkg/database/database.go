package database

import (
	"fmt"
	"log"

	"little-vote/pkg/model"
	"little-vote/pkg/viper"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

type RunOptions struct {
	User string
	Pass string
	Port string
	Name string
}

func Init() {
	NewRunOptions().Init()
}

func NewRunOptions() *RunOptions {
	return &RunOptions{
		User: viper.Config.GetString("db.user"),
		Pass: viper.Config.GetString("db.password"),
		Port: viper.Config.GetString("db.address"),
		Name: viper.Config.GetString("db.name"),
	}
}

func (options *RunOptions) Init() { // 初始化数据库
	dsn := fmt.Sprintf("%v:%v@tcp(%v)/%v?charset=utf8mb4&parseTime=True&loc=Local",
		options.User,
		options.Pass,
		options.Port,
		options.Name,
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Panicln("Database Error: ", err)
	} else {
		fmt.Printf("database start")
	}
	err = autoMigrate(db)
	if err != nil {
		log.Fatal("DatabaseMigrateFailed", err)
	}
	DB = db
}

func autoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&model.User{},
	)
}
