package controller

import (
	"net/http"
	"strconv"

	"github.com/chtavares592/server_go/model"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	validator "gopkg.in/go-playground/validator.v9"
)

type Worker struct {
	Db       *gorm.DB
	Validate *validator.Validate
}

var response struct {
	Posts []model.Post `json:"posts"`
}

func (w Worker) GetPosts(c echo.Context) error {
	var response struct {
		Posts []model.Post `json:"posts"`
	}

	posts := []model.Post{}
	if c.QueryParam("title") != "" {
		w.Db.Where("title ILIKE ?", "%"+c.QueryParam("title")+"%").Find(&posts)
		response.Posts = posts
		return c.JSON(http.StatusOK, response)
	}

	w.Db.Find(&posts)

	response.Posts = posts

	return c.JSON(http.StatusOK, response)
}

func (w Worker) GetPostId(c echo.Context) error {

	u64, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return err
	}

	id := uint(u64)

	posts := model.Post{}

	w.Db.First(&posts, id)
	w.Db.Where(model.Comment{PostID: id}).Find(&posts.Comments)

	return c.JSON(http.StatusOK, posts)
}

func (w Worker) ReceivePost(c echo.Context) error {
	recv := &model.Post{}
	err := c.Bind(recv)
	if err != nil {
		return c.String(http.StatusNotAcceptable, "ERROR1")
	}

	err = w.Validate.Struct(recv)
	if err != nil {
		return c.String(http.StatusNotAcceptable, "ERROR2")
	}

	if w.Db.Create(&recv).Error != nil {
		return c.String(http.StatusNotAcceptable, "ERROR3")
	}

	return c.JSON(http.StatusCreated, recv)
}

func (w Worker) ReceiveComment(c echo.Context) error {
	u64, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return err
	}
	id := uint(u64)

	recv := &model.Comment{}

	err = c.Bind(recv)
	if err != nil {
		return c.String(http.StatusNotAcceptable, "ERROR1")
	}

	err = w.Validate.Struct(recv)
	if err != nil {
		return c.String(http.StatusNotAcceptable, "ERROR2")
	}

	recv.PostID = id

	if w.Db.Create(&recv).Error != nil {
		return c.String(http.StatusNotAcceptable, "ERROR3")
	}

	return c.JSON(http.StatusCreated, recv)
}
