package tests

import (
	"encoding/json"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"myproject/controllers"
	"net/http"
	"net/http/httptest"
	"testing"

	beego "github.com/beego/beego/v2/server/web"
)

func TestFavoritesController_Get(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Set up the mock response for the external API
	httpmock.RegisterResponder("GET", "https://api.thecatapi.com/v1/favourites",
		httpmock.NewStringResponder(200, `[
			{
				"id": 232505403,
				"image_id": "MjAyMjUwMw",
				"sub_id": "user-123",
				"created_at": "2024-12-24T05:33:06.000Z",
				"image": {
					"id": "MjAyMjUwMw",
					"url": "https://cdn2.thecatapi.com/images/MjAyMjUwMw.jpg"
				}
			}
		]`))

	// Initialize the controller with the mock API key
	controller := &controllers.FavoritesController{}
	controller.APIKey = "test-api-key" // Directly inject API key

	// Register the route
	beego.Router("/favourites", controller)

	// Simulate a GET request
	req, err := http.NewRequest("GET", "/favourites", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(rr, req)

	// Assert the status code is 200
	assert.Equal(t, http.StatusOK, rr.Code, "Expected status 200, got %d", rr.Code)

	// Decode the response body
	var favorites []controllers.FavoriteResponse
	err = json.Unmarshal(rr.Body.Bytes(), &favorites)
	assert.NoError(t, err, "Error decoding JSON response")
	assert.NotEmpty(t, favorites)

	// Validate response content
	assert.Equal(t, 232505403, favorites[0].ID)
	assert.Equal(t, "MjAyMjUwMw", favorites[0].ImageID)
	assert.Equal(t, "user-123", favorites[0].SubID)
	assert.Equal(t, "https://cdn2.thecatapi.com/images/MjAyMjUwMw.jpg", favorites[0].Image.URL)
}
