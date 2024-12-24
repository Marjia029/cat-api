package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/beego/beego/v2/core/config"
	"github.com/beego/beego/v2/server/web"
)

type BreedImage struct {
	URL string `json:"url"`
}

type CatBreed struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	Origin       string `json:"origin"`
	WikipediaURL string `json:"wikipedia_url"`
}

type BreedSearchController struct {
	web.Controller
	APIKey string
}

// Initialize the controller with the API key
func (c *BreedSearchController) Prepare() {
	apiKey, err := config.String("api_key")
	if err != nil {
		c.CustomAbort(500, "Failed to load API key from configuration")
		return
	}
	c.APIKey = apiKey // Store the API key in the controller's field
}

// Fetch all breeds
func (c *BreedSearchController) Get() {
	url := "https://api.thecatapi.com/v1/breeds"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		c.CustomAbort(500, "Failed to create API request")
		return
	}
	req.Header.Set("x-api-key", c.APIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.CustomAbort(500, "Failed to fetch breed list")
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.CustomAbort(500, "Failed to read API response")
		return
	}

	var breeds []CatBreed
	if err := json.Unmarshal(body, &breeds); err != nil {
		c.CustomAbort(500, "Failed to parse breed list")
		return
	}

	c.Data["json"] = breeds
	c.ServeJSON()
}

// Fetch breed details and images concurrently
func (c *BreedSearchController) Post() {
	breedID := c.GetString("breed_id")
	if breedID == "" {
		c.CustomAbort(400, "Breed ID is required")
		return
	}

	// Channels for concurrent API calls
	breedDetailsChan := make(chan CatBreed)
	imagesChan := make(chan []BreedImage)
	errChan := make(chan error)

	// Fetch breed details concurrently
	go func() {
		url := "https://api.thecatapi.com/v1/breeds"
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("x-api-key", c.APIKey)

		client := &http.Client{}
		resp, _ := client.Do(req)
		defer resp.Body.Close()

		var breeds []CatBreed
		body, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal(body, &breeds)

		for _, breed := range breeds {
			if breed.ID == breedID {
				breedDetailsChan <- breed
				return
			}
		}
		errChan <- fmt.Errorf("Breed not found")
	}()

	// Fetch breed images concurrently
	go func() {
		url := fmt.Sprintf("https://api.thecatapi.com/v1/images/search?breed_ids=%s&limit=8", breedID)
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("x-api-key", c.APIKey)

		client := &http.Client{}
		resp, _ := client.Do(req)
		defer resp.Body.Close()

		var images []BreedImage
		body, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal(body, &images)

		imagesChan <- images
	}()

	// Wait for results or errors
	var selectedBreed CatBreed
	var breedImages []BreedImage

	select {
	case selectedBreed = <-breedDetailsChan:
	case err := <-errChan:
		c.CustomAbort(500, err.Error())
		return
	}

	select {
	case breedImages = <-imagesChan:
	case err := <-errChan:
		c.CustomAbort(500, err.Error())
		return
	}

	// Return the combined data as JSON
	c.Data["json"] = map[string]interface{}{
		"breed":  selectedBreed,
		"images": breedImages,
	}
	c.ServeJSON()
}
