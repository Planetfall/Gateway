package search_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	pb "github.com/planetfall/genproto/pkg/musicresearcher/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSearch(t *testing.T) {

	clientGiven := &clientMock{}

	queryGiven := "query-value"
	limitGiven := int32(10)
	genreListGiven := []string{"genre1", "genre2"}

	wGiven := httptest.NewRecorder()
	gGiven := getContextSearch(t, wGiven, queryGiven, limitGiven, genreListGiven)

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

	resultsActual := getSearchResultsActual(t, wGiven)
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
	gGiven := getContextSearch(t, wGiven, queryGiven, limitGiven, genreListGiven)

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
		gGiven := getContextSearch(t, wGiven, queryGiven, limitGiven, genreListGiven)

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
	gGiven := getContextSearch(t, wGiven, "", "", nil)

	errGiven := fmt.Errorf("test authentication error")
	connGiven.
		On("AuthenticateContext", mock.Anything).
		Return(context.Background(), errGiven)

	// when
	c.Search(gGiven)

	// then
	clientGiven.AssertNotCalled(t, "Search")
	assert.Equal(t,
		http.StatusInternalServerError, gGiven.Writer.Status())
}
