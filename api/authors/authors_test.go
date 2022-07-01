package authors

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bquenin/microservice/cmd/microservice/config"
	"github.com/bquenin/microservice/internal/database"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

type apiError struct {
	Error string
}

type ServiceTestSuite struct {
	suite.Suite
	router  *gin.Engine
	queries *database.Queries
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ServiceTestSuite))
}

func (suite *ServiceTestSuite) SetupSuite() {
	cfg, err := config.Read()
	suite.Require().NoError(err)

	postgres, err := database.NewPostgres(cfg.Postgres.Host, cfg.Postgres.User, cfg.Postgres.Password)
	suite.Require().NoError(err)

	suite.queries = database.New(postgres.DB)
	service := NewService(suite.queries)

	suite.router = gin.Default()
	service.RegisterHandlers(suite.router)
}

func (suite *ServiceTestSuite) SetupTest() {
	suite.queries.TruncateAuthor(context.Background())
}

func (suite *ServiceTestSuite) TestCreate() {
	request := apiAuthor{
		Name: "test author",
		Bio:  "A test bio",
	}
	var buffer bytes.Buffer
	suite.Require().NoError(json.NewEncoder(&buffer).Encode(request))

	req, err := http.NewRequest("POST", "/authors", &buffer)
	suite.Require().NoError(err)

	rec := httptest.NewRecorder()
	suite.router.ServeHTTP(rec, req)

	suite.Require().Equal(http.StatusCreated, rec.Result().StatusCode)
	var author apiAuthor
	suite.Require().NoError(json.NewDecoder(rec.Result().Body).Decode(&author))
	suite.Require().Equal(request.Name, author.Name)
	suite.Require().Equal(request.Bio, author.Bio)
}

func (suite *ServiceTestSuite) TestCreateBadRequest() {
	request := apiAuthor{
		Name: "the name of this author should be too long",
		Bio:  "A test bio",
	}
	var buffer bytes.Buffer
	suite.Require().NoError(json.NewEncoder(&buffer).Encode(request))

	req, err := http.NewRequest("POST", "/authors", &buffer)
	suite.Require().NoError(err)

	rec := httptest.NewRecorder()
	suite.router.ServeHTTP(rec, req)

	suite.Require().Equal(http.StatusBadRequest, rec.Result().StatusCode)
	var apiErr apiError
	suite.Require().NoError(json.NewDecoder(rec.Result().Body).Decode(&apiErr))
	suite.Require().Contains(apiErr.Error, "max")
}

func (suite *ServiceTestSuite) TestGet() {
	author, err := suite.queries.CreateAuthor(context.Background(), database.CreateAuthorParams{
		Name: "test author",
		Bio:  "A test bio",
	})
	suite.Require().NoError(err)

	req, err := http.NewRequest("GET", fmt.Sprintf("/authors/%d", author.ID), nil)
	suite.Require().NoError(err)

	rec := httptest.NewRecorder()
	suite.router.ServeHTTP(rec, req)

	suite.Require().Equal(http.StatusOK, rec.Result().StatusCode)
	var got apiAuthor
	suite.Require().NoError(json.NewDecoder(rec.Result().Body).Decode(&got))
	suite.Require().Equal(author.ID, got.ID)
	suite.Require().Equal(author.Name, got.Name)
	suite.Require().Equal(author.Bio, got.Bio)
}

func (suite *ServiceTestSuite) TestGetNotFound() {
	req, err := http.NewRequest("GET", "/authors/123", nil)
	suite.Require().NoError(err)

	rec := httptest.NewRecorder()
	suite.router.ServeHTTP(rec, req)

	suite.Require().Equal(http.StatusNotFound, rec.Result().StatusCode)
}

func (suite *ServiceTestSuite) TestGetBadRequest() {
	req, err := http.NewRequest("GET", "/authors/bad-request", nil)
	suite.Require().NoError(err)

	rec := httptest.NewRecorder()
	suite.router.ServeHTTP(rec, req)

	suite.Require().Equal(http.StatusBadRequest, rec.Result().StatusCode)
}

func (suite *ServiceTestSuite) TestFullUpdateBadRequest() {
	author, err := suite.queries.CreateAuthor(context.Background(), database.CreateAuthorParams{
		Name: "test author",
		Bio:  "A test bio",
	})
	suite.Require().NoError(err)

	request := apiAuthor{
		Name: "the name of this author should be too long",
		Bio:  "A test bio",
	}
	var buffer bytes.Buffer
	suite.Require().NoError(json.NewEncoder(&buffer).Encode(request))

	req, err := http.NewRequest("PUT", fmt.Sprintf("/authors/%d", author.ID), &buffer)
	suite.Require().NoError(err)

	rec := httptest.NewRecorder()
	suite.router.ServeHTTP(rec, req)

	suite.Require().Equal(http.StatusBadRequest, rec.Result().StatusCode)
}

func (suite *ServiceTestSuite) TestFullUpdate() {
	author, err := suite.queries.CreateAuthor(context.Background(), database.CreateAuthorParams{
		Name: "test author",
		Bio:  "A test bio",
	})
	suite.Require().NoError(err)

	update := apiAuthor{
		Name: "this is a better name",
		Bio:  "This is a way better updated bio",
	}
	var buffer bytes.Buffer
	suite.Require().NoError(json.NewEncoder(&buffer).Encode(update))

	req, err := http.NewRequest("PUT", fmt.Sprintf("/authors/%d", author.ID), &buffer)
	suite.Require().NoError(err)

	rec := httptest.NewRecorder()
	suite.router.ServeHTTP(rec, req)

	suite.Require().Equal(http.StatusOK, rec.Result().StatusCode)
	var got apiAuthor
	suite.Require().NoError(json.NewDecoder(rec.Result().Body).Decode(&got))
	suite.Require().Equal(author.ID, got.ID)
	suite.Require().Equal(update.Name, got.Name)
	suite.Require().Equal(update.Bio, got.Bio)
}

func (suite *ServiceTestSuite) TestPartialUpdateBadRequest() {
	author, err := suite.queries.CreateAuthor(context.Background(), database.CreateAuthorParams{
		Name: "test author",
		Bio:  "A test bio",
	})
	suite.Require().NoError(err)

	request := apiAuthor{
		Name: "the name of this author should be too long",
	}
	var buffer bytes.Buffer
	suite.Require().NoError(json.NewEncoder(&buffer).Encode(request))

	req, err := http.NewRequest("PATCH", fmt.Sprintf("/authors/%d", author.ID), &buffer)
	suite.Require().NoError(err)

	rec := httptest.NewRecorder()
	suite.router.ServeHTTP(rec, req)

	suite.Require().Equal(http.StatusBadRequest, rec.Result().StatusCode)
}

func (suite *ServiceTestSuite) TestPartialUpdate() {
	author, err := suite.queries.CreateAuthor(context.Background(), database.CreateAuthorParams{
		Name: "test author",
		Bio:  "A test bio",
	})
	suite.Require().NoError(err)

	update := apiAuthor{
		Bio: "This is a way better updated bio",
	}
	var buffer bytes.Buffer
	suite.Require().NoError(json.NewEncoder(&buffer).Encode(update))

	req, err := http.NewRequest("PATCH", fmt.Sprintf("/authors/%d", author.ID), &buffer)
	suite.Require().NoError(err)

	rec := httptest.NewRecorder()
	suite.router.ServeHTTP(rec, req)

	suite.Require().Equal(http.StatusOK, rec.Result().StatusCode)
	var got apiAuthor
	suite.Require().NoError(json.NewDecoder(rec.Result().Body).Decode(&got))
	suite.Require().Equal(author.ID, got.ID)
	suite.Require().Equal(author.Name, got.Name)
	suite.Require().Equal(update.Bio, got.Bio)
}

func (suite *ServiceTestSuite) TestDelete() {
	author, err := suite.queries.CreateAuthor(context.Background(), database.CreateAuthorParams{
		Name: "test author",
		Bio:  "A test bio",
	})
	suite.Require().NoError(err)

	req, err := http.NewRequest("DELETE", fmt.Sprintf("/authors/%d", author.ID), nil)
	suite.Require().NoError(err)

	rec := httptest.NewRecorder()
	suite.router.ServeHTTP(rec, req)

	suite.Require().Equal(http.StatusOK, rec.Result().StatusCode)
	_, err = suite.queries.GetAuthor(context.Background(), author.ID)
	suite.Require().Error(err)
}

func (suite *ServiceTestSuite) TestList() {
	req, err := http.NewRequest("GET", "/authors", nil)
	suite.Require().NoError(err)

	rec := httptest.NewRecorder()
	suite.router.ServeHTTP(rec, req)

	suite.Require().Equal(http.StatusNotFound, rec.Result().StatusCode)
}
