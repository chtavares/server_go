package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/chtavares592/server_go/controller"
	"github.com/chtavares592/server_go/model"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	validator "gopkg.in/go-playground/validator.v9"
)

type TestSuite struct {
	suite.Suite
	w controller.Worker
}

func setupDbTest() (*gorm.DB, error) {
	db, err := gorm.Open("postgres", "dbname=testeblogdb")
	if err != nil {
		return nil, err
	}

	if db.AutoMigrate(&model.Post{}, &model.Comment{}).Error != nil {
		return nil, err
	}

	return db, nil
}

func (suite *TestSuite) SetupTest() {
	suite.w.Validate = validator.New()
	var err error

	suite.w.Db, err = setupDbTest()
	if err != nil {
		panic(err)
	}
}

func (suite *TestSuite) TearDownTest() {
	suite.w.Db.Exec("DROP TABLE comments;")
	suite.w.Db.Exec("DROP TABLE posts;")

}

func connectGet() (echo.Context, *httptest.ResponseRecorder) {
	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	return c, rec
}

func connectPost(json string) (echo.Context, *httptest.ResponseRecorder) {
	e := echo.New()
	req := httptest.NewRequest(echo.POST, "/", strings.NewReader(json))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	return c, rec
}

func (suite *TestSuite) TestGetPosts() {

	posts := model.Post{Title: "Ola mundo", Body: "Bem vindo"}
	if suite.w.Db.Create(&posts).Error != nil {
		fmt.Printf("Error to create")
	}
	posts = model.Post{Title: "Tchau mundo", Body: "Boa viagem"}
	if suite.w.Db.Create(&posts).Error != nil {
		fmt.Printf("Error to create")
	}

	c, rec := connectGet()
	c.SetPath("/posts")

	expectPost := `{"posts":[{"id":1,"title":"Ola mundo","body":"Bem vindo","comments":null},{"id":2,"title":"Tchau mundo","body":"Boa viagem","comments":null}]}`

	if assert.NoError(suite.T(), suite.w.GetPosts(c)) {
		assert.Equal(suite.T(), http.StatusOK, rec.Code)
		assert.Equal(suite.T(), expectPost, rec.Body.String())
	}

}

func (suite *TestSuite) TestGetTitlePost() {
	post := &model.Post{ID: 7, Title: "Passeio de uma vida", Body: "Estava caminhando quando"}

	if suite.w.Db.Create(&post).Error != nil {
		fmt.Printf("ERROR to create post")
	}

	c, rec := connectGet()
	q := make(url.Values)
	q.Set("title", "vida")
	c.SetPath("/posts?" + q.Encode())

	expectPostsJSON := `{"posts":[{"id":7,"title":"Passeio de uma vida","body":"Estava caminhando quando","comments":null}]}`

	if assert.NoError(suite.T(), suite.w.GetPosts(c)) {
		assert.Equal(suite.T(), http.StatusOK, rec.Code)
		assert.Equal(suite.T(), string(expectPostsJSON), rec.Body.String())
	}

}

func (suite *TestSuite) TestCreatePost() {
	postJSON := `{"title":"Testando","body":"Realmente o dia está lindo"}`

	c, rec := connectPost(postJSON)
	c.SetPath("/posts")

	expectPostJSON := `{"id":1,"title":"Testando","body":"Realmente o dia está lindo","comments":null}`

	if assert.NoError(suite.T(), suite.w.ReceivePost(c)) {
		assert.Equal(suite.T(), http.StatusCreated, rec.Code)
		assert.Equal(suite.T(), expectPostJSON, rec.Body.String())
	}

}

func (suite *TestSuite) TestCreateComment() {
	posts := &model.Post{ID: 8, Title: "Belo dia", Body: "O sol está radiante e feliz"}
	if suite.w.Db.Create(&posts).Error != nil {
		fmt.Printf("ERROR to create post")
	}

	commentJSON := `{"name":"Rodrigo","body":"Realmente está um belo dia"}`

	c, rec := connectPost(commentJSON)
	c.SetPath("/posts/:id/comments")
	c.SetParamNames("id")
	c.SetParamValues("8")

	expectCommentJSON := `{"id":1,"name":"Rodrigo","body":"Realmente está um belo dia","postId":8}`

	if assert.NoError(suite.T(), suite.w.ReceiveComment(c)) {
		assert.Equal(suite.T(), http.StatusCreated, rec.Code)
		assert.Equal(suite.T(), string(expectCommentJSON), rec.Body.String())
	}

}

func (suite *TestSuite) TestGetPost() {
	post := &model.Post{ID: 5, Title: "GEBrasil", Body: "O melhor clube do brasil"}

	if suite.w.Db.Create(&post).Error != nil {
		fmt.Printf("ERROR to create post")
	}

	c, rec := connectGet()
	c.SetPath("/posts/:id")
	c.SetParamNames("id")
	c.SetParamValues("5")

	expectPostJSON := `{"id":5,"title":"GEBrasil","body":"O melhor clube do brasil","comments":[]}`

	if assert.NoError(suite.T(), suite.w.GetPostId(c)) {
		assert.Equal(suite.T(), http.StatusOK, rec.Code)
		assert.Equal(suite.T(), string(expectPostJSON), rec.Body.String())
	}

}

func TestServerGoSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
