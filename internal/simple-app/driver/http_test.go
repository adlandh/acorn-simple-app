package driver

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/adlandh/acorn-simple-app/internal/simple-app/domain"

	"github.com/adlandh/acorn-simple-app/internal/simple-app/domain/mocks"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gavv/httpexpect/v2"
	"github.com/labstack/echo/v4"
	"github.com/phayes/freeport"
	"github.com/stretchr/testify/suite"
)

var fakeError = errors.New("fake error")

const (
	apiUser = "/api/user"
)

type HttpServerTestSuite struct {
	suite.Suite
	e      *echo.Echo
	tester *httpexpect.Expect
	app    *mocks.ApplicationInterface
	url    string
}

func (s *HttpServerTestSuite) SetupSuite() {
	s.app = new(mocks.ApplicationInterface)
	s.e = echo.New()
	RegisterHandlers(s.e, NewHTTPServer(s.app))
	port, err := freeport.GetFreePort()
	s.Require().NoError(err)
	go func() {
		err := s.e.Start(":" + strconv.Itoa(port))
		s.Require().True(err == nil || errors.Is(err, http.ErrServerClosed))
	}()
	time.Sleep(time.Second)
	s.url = "http://localhost:" + strconv.Itoa(port)
	s.tester = httpexpect.Default(s.T(), s.url)
}

func (s *HttpServerTestSuite) TearDownSuite() {
	err := s.e.Shutdown(context.Background())
	s.NoError(err)
}

func (s *HttpServerTestSuite) TearDownTest() {
	s.app.AssertExpectations(s.T())
}

func (s *HttpServerTestSuite) TestHealthCheck() {
	s.tester.GET("/").WithHeader(echo.HeaderContentType, echo.MIMETextPlain).
		Expect().
		Status(http.StatusOK).
		Text().Contains("Ok")
}

func (s *HttpServerTestSuite) TestCreateUser() {
	id, err := uuid.NewUUID()
	s.Require().NoError(err)
	name := gofakeit.Username()

	s.Run("happy case", func() {
		s.app.On("CreateUser", mock.Anything, name).Return(id, nil).Once()
		s.tester.POST(apiUser).
			WithJSON(UserRequest{Name: name}).
			Expect().
			Status(http.StatusOK).JSON().Object().HasValue("id", id).HasValue("name", name)
	})

	s.Run("error in app", func() {
		s.app.On("CreateUser", mock.Anything, name).Return(id, fakeError).Once()
		s.tester.POST(apiUser).
			WithJSON(UserRequest{Name: name}).
			Expect().
			Status(http.StatusInternalServerError).JSON().Object().HasValue("message", fakeError.Error())
	})
}

func (s *HttpServerTestSuite) TestGetUser() {
	id, err := uuid.NewUUID()
	s.Require().NoError(err)
	name := gofakeit.Username()

	s.Run("happy case", func() {
		s.app.On("GetUser", mock.Anything, id).Return(name, nil).Once()
		s.tester.GET(apiUser+"/"+id.String()).
			Expect().
			Status(http.StatusOK).JSON().Object().HasValue("id", id).HasValue("name", name)
	})

	s.Run("invalid id", func() {
		invalidId := gofakeit.Word()
		s.tester.GET(apiUser + "/" + invalidId).
			Expect().
			Status(http.StatusBadRequest).JSON().Object().ContainsKey("message")
	})

	s.Run("not found", func() {
		s.app.On("GetUser", mock.Anything, id).Return("", domain.ErrorNotFound).Once()
		s.tester.GET(apiUser + "/" + id.String()).
			Expect().
			Status(http.StatusNotFound).NoContent()
	})

	s.Run("error in app", func() {
		s.app.On("GetUser", mock.Anything, id).Return("", fakeError).Once()
		s.tester.GET(apiUser+"/"+id.String()).
			Expect().
			Status(http.StatusInternalServerError).JSON().Object().HasValue("message", fakeError.Error())
	})
}

func (s *HttpServerTestSuite) TestUpdateUser() {
	id, err := uuid.NewUUID()
	s.Require().NoError(err)
	name := gofakeit.Username()

	s.Run("happy case", func() {
		s.app.On("UpdateUser", mock.Anything, id, name).Return(nil).Once()
		s.tester.POST(apiUser+"/"+id.String()).
			WithJSON(UserRequest{Name: name}).
			Expect().
			Status(http.StatusOK).JSON().Object().HasValue("id", id).HasValue("name", name)
	})

	s.Run("not found", func() {
		s.app.On("UpdateUser", mock.Anything, id, name).Return(domain.ErrorNotFound).Once()
		s.tester.POST(apiUser + "/" + id.String()).
			WithJSON(UserRequest{Name: name}).
			Expect().
			Status(http.StatusNotFound).NoContent()
	})

	s.Run("error in app", func() {
		s.app.On("UpdateUser", mock.Anything, id, name).Return(fakeError).Once()
		s.tester.POST(apiUser+"/"+id.String()).
			WithJSON(UserRequest{Name: name}).
			Expect().
			Status(http.StatusInternalServerError).JSON().Object().HasValue("message", fakeError.Error())
	})
}

func (s *HttpServerTestSuite) TestDeleteUser() {
	id, err := uuid.NewUUID()
	s.Require().NoError(err)

	s.Run("happy case", func() {
		s.app.On("DeleteUser", mock.Anything, id).Return(nil).Once()
		s.tester.DELETE(apiUser + "/" + id.String()).
			Expect().
			Status(http.StatusOK).NoContent()
	})

	s.Run("not found", func() {
		s.app.On("DeleteUser", mock.Anything, id).Return(domain.ErrorNotFound).Once()
		s.tester.DELETE(apiUser + "/" + id.String()).
			Expect().
			Status(http.StatusNotFound).NoContent()
	})

	s.Run("error in app", func() {
		s.app.On("DeleteUser", mock.Anything, id).Return(fakeError).Once()
		s.tester.DELETE(apiUser+"/"+id.String()).
			Expect().
			Status(http.StatusInternalServerError).JSON().Object().HasValue("message", fakeError.Error())
	})
}

func TestHttpServer(t *testing.T) {
	suite.Run(t, new(HttpServerTestSuite))
}
