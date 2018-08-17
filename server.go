package main

import (
	"github.com/chtavares592/server_go/controller"
	"github.com/chtavares592/server_go/model"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/labstack/echo"
	"gopkg.in/go-playground/validator.v9"
)

func setupDB() (*gorm.DB, error) {
	db, err := gorm.Open("postgres", "dbname=blogdb")
	if err != nil {
		return nil, err
	}

	if db.AutoMigrate(&model.Post{}, &model.Comment{}).Error != nil {
		return nil, err
	}

	return db, nil
}

func main() {
	var err error
	worker := controller.Worker{}
	worker.Validate = validator.New()

	worker.Db, err = setupDB()
	if err != nil {
		panic(err)
	}

	e := echo.New()

	e.GET("/posts", worker.GetPosts)

	e.GET("/posts/:id", worker.GetPostId)

	e.POST("/posts", worker.ReceivePost)

	e.POST("/posts/:id/comments", worker.ReceiveComment)

	e.Logger.Fatal(e.Start(":1323"))
}
