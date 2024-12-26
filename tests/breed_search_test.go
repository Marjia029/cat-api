package tests

import (
	"encoding/json"
	"github.com/beego/beego/v2/core/config"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"myproject/controllers"
	"net/http"
	"net/http/httptest"
	"testing"

	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	// Mock the API Key in configuration
	config.Set("api_key", "test-api-key")
}

func TestBreedSearchController_Get(t *testing.T) {
	// Initialize httpmock to mock external HTTP requests
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Set up the mock response for the external API
	httpmock.RegisterResponder("GET", "https://api.thecatapi.com/v1/breeds",
		httpmock.NewStringResponder(200, `[
			{
				"id": "abys",
				"name": "Abyssinian",
				"description": "The Abyssinian is easy to care for, and a joy to have in your home. They’re affectionate cats and love both people and other animals.",
				"origin": "Egypt",
				"wikipedia_url": "https://en.wikipedia.org/wiki/Abyssinian_(cat)"
			}
		]`))

	// Create a mock Beego controller and set up the route
	beego.Router("/breed-search", &controllers.BreedSearchController{})

	// Create a new GET request to simulate a request to "/breed-search"
	req, err := http.NewRequest("GET", "/breed-search", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a mock recorder to capture the response
	rr := httptest.NewRecorder()

	// Handle the request using Beego's HTTP handler
	beego.BeeApp.Handlers.ServeHTTP(rr, req)

	// Assert the status code is 200
	assert.Equal(t, http.StatusOK, rr.Code)

	// Print the raw response body to inspect it
	t.Log("Response Body:", rr.Body.String())

	// Decode the response body into a slice of CatBreed
	var breeds []controllers.CatBreed
	if err := json.Unmarshal(rr.Body.Bytes(), &breeds); err != nil {
		t.Fatal(err)
	}

	// Assert that the response body contains the expected structure (i.e., a list of breeds)
	assert.NotEmpty(t, breeds)
	assert.Equal(t, "abys", breeds[0].ID)
	assert.Equal(t, "Abyssinian", breeds[0].Name)
	assert.Equal(t, "The Abyssinian is easy to care for, and a joy to have in your home. They’re affectionate cats and love both people and other animals.", breeds[0].Description)
	assert.Equal(t, "Egypt", breeds[0].Origin)
	assert.Equal(t, "https://en.wikipedia.org/wiki/Abyssinian_(cat)", breeds[0].WikipediaURL)
}
