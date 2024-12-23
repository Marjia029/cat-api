package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/beego/beego/v2/server/web"
)

type VotingController struct {
	web.Controller
}

type VotingCatImage struct {
	URL string `json:"url"`
}

var favorites []string

// Channel for fetching random cat images
var fetchImageChan = make(chan string)
var favoriteActionChan = make(chan string)

// Fetch a random cat image concurrently
func fetchRandomCatImage(apiKey string) {
	url := "https://api.thecatapi.com/v1/images/search"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fetchImageChan <- ""
		fmt.Println("Failed to create request:", err)
		return
	}
	req.Header.Set("x-api-key", apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fetchImageChan <- ""
		fmt.Println("Failed to fetch cat image:", err)
		return
	}
	defer resp.Body.Close()

	var images []VotingCatImage
	if err := json.NewDecoder(resp.Body).Decode(&images); err != nil {
		fetchImageChan <- ""
		fmt.Println("Failed to decode cat image response:", err)
		return
	}

	if len(images) > 0 {
		fetchImageChan <- images[0].URL
	} else {
		fetchImageChan <- ""
	}
}

// Get method to fetch a random cat image and return it as JSON
func (c *VotingController) Get() {
	apiKey := "live_GWXcPdnWze27MNMJSjinKshtfsnVsi4EdrXfKUNhOmXsLakl5N7MwJCShLvC5Rxo"

	// Launch goroutine to fetch random cat image
	go fetchRandomCatImage(apiKey)

	imageURL := <-fetchImageChan // Wait for the result from the channel

	if imageURL == "" {
		c.Data["json"] = map[string]interface{}{"error": "Failed to fetch cat image"}
	} else {
		c.Data["json"] = map[string]interface{}{
			"image_url": imageURL,
			"favorites": favorites,
		}
	}

	c.ServeJSON()
}

// Post method to handle like, dislike, and saving to favorites
func (c *VotingController) Post() {
	action := c.GetString("action")
	imageURL := c.GetString("image_url")

	// Handle favorite action
	if action == "favorite" {
		go func() {
			favorites = append(favorites, imageURL)
			favoriteActionChan <- "done"
		}()
	}

	// Wait for favorite action to complete (if applicable)
	if action == "favorite" {
		<-favoriteActionChan
	}

	// Fetch a new random cat image if like/dislike action
	if action == "like" || action == "dislike" || action == "favorite" {
		c.Get() // Fetch a new image and update `c.Data["ImageURL"]`
	}

	// Send the updated image URL and favorites list as a JSON response
	c.Data["json"] = map[string]interface{}{"image_url": c.Data["ImageURL"], "favorites": favorites}
	c.ServeJSON()
}
