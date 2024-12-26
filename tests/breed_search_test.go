package tests

import (
	"encoding/json"
	"github.com/beego/beego/v2/server/web/context"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"myproject/controllers"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestBreedSearchController_Get(t *testing.T) {
	// Initialize httpmock
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Mock the external API response
	httpmock.RegisterResponder("GET", "https://api.thecatapi.com/v1/breeds",
		httpmock.NewStringResponder(200, `[
			{
				"id": "abys",
				"name": "Abyssinian",
				"description": "The Abyssinian is easy to care for, and a joy to have in your home. They're affectionate cats and love both people and other animals.",
				"origin": "Egypt",
				"wikipedia_url": "https://en.wikipedia.org/wiki/Abyssinian_(cat)"
			}
		]`))

	// Create a test controller with a test API key
	controller := &controllers.BreedSearchController{
		APIKey: "test-api-key",
	}

	// Create a new request and recorder
	r, _ := http.NewRequest("GET", "/breed-search", nil)
	w := httptest.NewRecorder()

	// Create a new context
	ctx := context.NewContext()
	ctx.Reset(w, r)

	// Initialize the controller with the context
	controller.Init(ctx, "", "", controller)

	// Call the Get method
	controller.Get()

	// Assert the status code
	assert.Equal(t, http.StatusOK, w.Code)

	// Decode and verify the response
	var breeds []controllers.CatBreed
	err := json.Unmarshal(w.Body.Bytes(), &breeds)
	assert.NoError(t, err)

	// Verify the response content
	assert.NotEmpty(t, breeds)
	assert.Equal(t, "abys", breeds[0].ID)
	assert.Equal(t, "Abyssinian", breeds[0].Name)
	assert.Equal(t, "The Abyssinian is easy to care for, and a joy to have in your home. They're affectionate cats and love both people and other animals.", breeds[0].Description)
	assert.Equal(t, "Egypt", breeds[0].Origin)
	assert.Equal(t, "https://en.wikipedia.org/wiki/Abyssinian_(cat)", breeds[0].WikipediaURL)
}

func TestBreedSearchController_Post(t *testing.T) {
	// Initialize httpmock
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Mock the breeds list endpoint
	httpmock.RegisterResponder("GET", "https://api.thecatapi.com/v1/breeds",
		httpmock.NewStringResponder(200, `[
			{
				"id": "abys",
				"name": "Abyssinian",
				"description": "The Abyssinian is easy to care for, and a joy to have in your home.",
				"origin": "Egypt",
				"wikipedia_url": "https://en.wikipedia.org/wiki/Abyssinian_(cat)"
			}
		]`))

	// Mock the images search endpoint
	httpmock.RegisterResponder("GET", `=~^https://api.thecatapi.com/v1/images/search\?breed_ids=abys`,
		httpmock.NewStringResponder(200, `[
			{"url": "https://example.com/cat1.jpg"},
			{"url": "https://example.com/cat2.jpg"}
		]`))

	// Create a test controller with a test API key
	controller := &controllers.BreedSearchController{
		APIKey: "test-api-key",
	}

	// Create a new POST request with breed_id parameter
	r, _ := http.NewRequest("POST", "/breed-search?breed_id=abys", nil)
	w := httptest.NewRecorder()

	// Create a new context
	ctx := context.NewContext()
	ctx.Reset(w, r)

	// Initialize the controller with the context
	controller.Init(ctx, "", "", controller)

	// Call the Post method
	controller.Post()

	// Assert the status code
	assert.Equal(t, http.StatusOK, w.Code)

	// Parse the response
	var response struct {
		Breed  controllers.CatBreed     `json:"breed"`
		Images []controllers.BreedImage `json:"images"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Verify breed details
	assert.Equal(t, "abys", response.Breed.ID)
	assert.Equal(t, "Abyssinian", response.Breed.Name)
	assert.Equal(t, "Egypt", response.Breed.Origin)
	assert.Equal(t, "https://en.wikipedia.org/wiki/Abyssinian_(cat)", response.Breed.WikipediaURL)
	assert.True(t, strings.Contains(response.Breed.Description, "The Abyssinian is easy to care for"))

	// Verify images
	assert.Len(t, response.Images, 2)
	assert.Equal(t, "https://example.com/cat1.jpg", response.Images[0].URL)
	assert.Equal(t, "https://example.com/cat2.jpg", response.Images[1].URL)
}
