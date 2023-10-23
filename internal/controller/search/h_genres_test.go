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

func TestGetGenreList(t *testing.T) {

	// given
	clientGiven := &clientMock{}

	wGiven := httptest.NewRecorder()
	gGiven := getContextGenreList(t, wGiven)

	genreListGiven := &pb.GenreList{
		Genres: []string{"genre1", "genre2"},
	}

	c := getController(t, clientGiven, nil)

	// when
	clientGiven.On("GetGenreList").Return(genreListGiven, nil)

	c.GetGenreList(gGiven)

	// then
	clientGiven.AssertExpectations(t)
	assert.Equal(t, http.StatusOK, gGiven.Writer.Status())

	genreListActual := getGenreListActual(t, wGiven)
	assert.Equal(t, genreListGiven.Genres, genreListActual.Genres)
}

func TestGetGenreList_withInternalError(t *testing.T) {

	// given
	clientGiven := &clientMock{}

	wGiven := httptest.NewRecorder()
	gGiven := getContextGenreList(t, wGiven)

	c := getController(t, clientGiven, nil)

	// when
	errorGiven := fmt.Errorf("test genre list error")
	clientGiven.On("GetGenreList").Return(&pb.GenreList{}, errorGiven)

	c.GetGenreList(gGiven)

	// then
	clientGiven.AssertExpectations(t)
	assert.Equal(t,
		http.StatusInternalServerError, gGiven.Writer.Status())
}

func TestGetGenreList_withAuthenticateError(t *testing.T) {

	// given
	clientGiven := &clientMock{}
	connGiven := &connectionMock{}
	c := getController(t, clientGiven, connGiven)

	wGiven := httptest.NewRecorder()
	gGiven := getContextGenreList(t, wGiven)

	errGiven := fmt.Errorf("test authentication error")
	connGiven.
		On("AuthenticateContext", mock.Anything).
		Return(context.Background(), errGiven)

	// when
	c.GetGenreList(gGiven)

	// then
	clientGiven.AssertNotCalled(t, "GetGenreList")
	assert.Equal(t,
		http.StatusInternalServerError, gGiven.Writer.Status())
}
