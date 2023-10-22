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

func getContext(t *testing.T, wGiven *httptest.ResponseRecorder,
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

func getResultsActual(
	t *testing.T, wGiven *httptest.ResponseRecorder) *pb.Results {

	body, err := io.ReadAll(wGiven.Body)
	assert.Nil(t, err)

	var resultsActual *pb.Results
	err = json.Unmarshal(body, &resultsActual)
	assert.Nil(t, err)

	return resultsActual
}

func TestSearch(t *testing.T) {

	clientGiven := &clientMock{}

	queryGiven := "query-value"
	limitGiven := int32(10)
	genreListGiven := []string{"genre1", "genre2"}

	wGiven := httptest.NewRecorder()
	gGiven := getContext(t, wGiven, queryGiven, limitGiven, genreListGiven)

	trackIdGiven := "track-id"
	resultsGiven := &pb.Results{
		Tracks: []*pb.Track{
			{
				ID: trackIdGiven,
			},
		},
	}

	c := getController(t, clientGiven, nil)

	// when
	clientGiven.
		On("Search", queryGiven, limitGiven, genreListGiven).
		Return(resultsGiven, nil)

	c.Search(gGiven)

	// then
	clientGiven.AssertExpectations(t)
	assert.Equal(t, http.StatusOK, gGiven.Writer.Status())

	resultsActual := getResultsActual(t, wGiven)
	assert.Equal(t, len(resultsActual.Tracks), len(resultsGiven.Tracks))
	assert.Equal(t, resultsActual.Tracks[0].ID, resultsGiven.Tracks[0].ID)
}

func TestSearch_withInternalError(t *testing.T) {

	// given
	clientGiven := &clientMock{}

	queryGiven := "query-value"
	limitGiven := int32(10)
	genreListGiven := []string{"genre1", "genre2"}

	wGiven := httptest.NewRecorder()
	gGiven := getContext(t, wGiven, queryGiven, limitGiven, genreListGiven)

	c := getController(t, clientGiven, nil)

	// when
	errorGiven := fmt.Errorf("test search error")
	clientGiven.
		On("Search", queryGiven, limitGiven, genreListGiven).
		Return(&pb.Results{}, errorGiven)

	c.Search(gGiven)

	// then
	clientGiven.AssertExpectations(t)
	assert.Equal(t, http.StatusInternalServerError,
		gGiven.Writer.Status())
}

func TestSearch_withInvalidLimit(t *testing.T) {

	// given
	clientGiven := &clientMock{}
	c := getController(t, clientGiven, nil)

	queryGiven := "query-value"
	genreListGiven := []string{"genre1", "genre2"}

	wGiven := httptest.NewRecorder()

	limitListGiven := []interface{}{"salut", true, false}
	for _, limitGiven := range limitListGiven {
		gGiven := getContext(t, wGiven, queryGiven, limitGiven, genreListGiven)

		// when
		c.Search(gGiven)

		// then
		clientGiven.AssertNotCalled(t, "Search")
		assert.Equal(t, http.StatusBadRequest, gGiven.Writer.Status())
	}
}

func TestSearch_withAuthenticateError(t *testing.T) {

	// given
	connGiven := &connectionMock{}
	clientGiven := &clientMock{}
	c := getController(t, clientGiven, connGiven)

	wGiven := httptest.NewRecorder()
	gGiven := getContext(t, wGiven, "", "", nil)

	errMessageGiven := "test authentication error"
	connGiven.
		On("AuthenticateContext", mock.Anything, mock.Anything).
		Return(context.Background(), fmt.Errorf(errMessageGiven))

	// when
	c.Search(gGiven)

	// then
	clientGiven.AssertNotCalled(t, "Search")
	assert.Equal(t, http.StatusInternalServerError, gGiven.Writer.Status())
}
