package search_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/planetfall/gateway/internal/controller"
	"github.com/planetfall/gateway/internal/controller/search"
	pb "github.com/planetfall/genproto/pkg/musicresearcher/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	grpc "google.golang.org/grpc"
)

type clientMock struct {
	mock.Mock
}

func (m *clientMock) Search(ctx context.Context, p *pb.Parameters,
	opts ...grpc.CallOption) (*pb.Results, error) {

	args := m.Called(p.Query, p.Limit, p.GenreFilters)
	return args.Get(0).(*pb.Results), args.Error(1)
}

func (m *clientMock) GetGenreList(context.Context, *pb.Empty,
	...grpc.CallOption) (*pb.GenreList, error) {

	args := m.Called()
	return args.Get(0).(*pb.GenreList), args.Error(1)
}

type connectionMock struct {
	mock.Mock
}

func (m *connectionMock) GrpcConn() *grpc.ClientConn {
	args := m.Called()
	return args.Get(0).(*grpc.ClientConn)
}

func (m *connectionMock) Close() error {
	args := m.Called()
	return args.Error(0)
}

func (m *connectionMock) AuthenticateContext(
	ctx context.Context) (context.Context, error) {

	args := m.Called(ctx)
	return args.Get(0).(context.Context), args.Error(1)
}

func getController(
	t *testing.T,
	clientGiven *clientMock,
	connGiven *connectionMock) *search.SearchController {

	loggerGiven := log.Default()
	reportError := func(err error) {
		loggerGiven.Println(err)
	}

	optGiven := search.SearchControllerOptions{
		ControllerOptions: controller.ControllerOptions{
			Name:        "name",
			Target:      "target",
			ReportError: reportError,
			Logger:      loggerGiven,
		},
		Insecure: true,
	}

	if connGiven != nil {
		optGiven.Conn = connGiven
	}

	if clientGiven != nil {
		optGiven.Client = clientGiven
	}

	c, err := search.NewSearchController(optGiven)
	assert.Nil(t, err)
	assert.NotNil(t, c)

	return c
}

func getContextSearch(t *testing.T, wGiven *httptest.ResponseRecorder,
	query interface{}, limit interface{}, genreList []string) *gin.Context {

	gGiven, _ := gin.CreateTestContext(wGiven)
	req, err := http.NewRequest(http.MethodGet, "", nil)
	assert.Nil(t, err)
	q := req.URL.Query()
	q.Set("q", fmt.Sprint(query))
	q.Set("limit", fmt.Sprint(limit))
	for _, genre := range genreList {
		q.Add("genre", genre)
	}

	req.URL.RawQuery = q.Encode()
	gGiven.Request = req

	return gGiven
}

func getContextGenreList(
	t *testing.T, wGiven *httptest.ResponseRecorder) *gin.Context {

	gGiven, _ := gin.CreateTestContext(wGiven)
	req, err := http.NewRequest(http.MethodGet, "", nil)
	assert.Nil(t, err)

	gGiven.Request = req
	return gGiven
}

func getSearchResultsActual(
	t *testing.T, wGiven *httptest.ResponseRecorder) *pb.Results {

	body, err := io.ReadAll(wGiven.Body)
	assert.Nil(t, err)

	var resultsActual *pb.Results
	err = json.Unmarshal(body, &resultsActual)
	assert.Nil(t, err)

	return resultsActual
}

func getGenreListActual(
	t *testing.T, wGiven *httptest.ResponseRecorder) *pb.GenreList {

	body, err := io.ReadAll(wGiven.Body)
	assert.Nil(t, err)

	var genreListActual *pb.GenreList
	err = json.Unmarshal(body, &genreListActual)
	assert.Nil(t, err)

	return genreListActual
}
