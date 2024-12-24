package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/beego/beego/v2/core/config"
	"github.com/beego/beego/v2/server/web"
)

type VotingController struct {
	web.Controller
	APIKey string
}

type VotingCatImage struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}

// Channels for fetching random cat images and handling favorite actions
var fetchImageChan = make(chan string)
var favoriteActionChan = make(chan string)

// Initialize the controller with the API key
func (c *VotingController) Prepare() {
	apiKey, err := config.String("api_key")
	if err != nil {
		c.Data["json"] = map[string]interface{}{"error": "Failed to load API key from configuration"}
		c.ServeJSON()
		return
	}
	c.APIKey = apiKey // Store the API key in the controller's field
}

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
		fetchImageChan <- images[0].URL + "|" + images[0].ID // Send URL and ID as a combined string
	} else {
		fetchImageChan <- ""
	}
}

// Get method to fetch a random cat image and return it as JSON
func (c *VotingController) Get() {

	go fetchRandomCatImage(c.APIKey)

	imageData := <-fetchImageChan // Wait for the result from the channel
	if imageData == "" {
		c.Data["json"] = map[string]interface{}{"error": "Failed to fetch cat image"}
	} else {
		dataParts := strings.Split(imageData, "|") // Split URL and ID
		imageURL, imageID := dataParts[0], dataParts[1]

		c.Data["json"] = map[string]interface{}{
			"image_url": imageURL,
			"image_id":  imageID,
		}
	}

	c.ServeJSON()
}

// Post method to handle like, dislike, and saving to favorites
func (c *VotingController) Post() {
	action := c.GetString("action")
	imageID := c.GetString("image_id")
	userID := "user-123" // Use unique user ID here (you can replace it)

	// Handle favorite action
	if action == "favorite" {
		go func() {
			favoriteBody := map[string]string{
				"image_id": imageID,
				"sub_id":   userID,
			}

			bodyBytes, err := json.Marshal(favoriteBody)
			if err != nil {
				fmt.Println("Failed to marshal favorite request body:", err)
				favoriteActionChan <- "error"
				return
			}

			url := "https://api.thecatapi.com/v1/favourites"
			req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyBytes))
			if err != nil {
				fmt.Println("Failed to create favorite request:", err)
				favoriteActionChan <- "error"
				return
			}
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("x-api-key", c.APIKey)

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				fmt.Println("Failed to make favorite request:", err)
				favoriteActionChan <- "error"
				return
			}
			defer resp.Body.Close()

			body, _ := io.ReadAll(resp.Body)
			fmt.Println("Response Status:", resp.Status)
			fmt.Println("Response Body:", string(body))

			if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
				favoriteActionChan <- "done"
			} else {
				favoriteActionChan <- "error"
			}
		}()
	}

	// Wait for favorite action to complete (if applicable)
	if action == "favorite" {
		status := <-favoriteActionChan
		if status == "error" {
			c.Data["json"] = map[string]interface{}{"error": "Failed to favorite the image"}
			c.ServeJSON()
			return
		}
	}

	// Handle voting (like or dislike)
	if action == "like" || action == "dislike" {
		voteValue := 0
		if action == "like" {
			voteValue = 1
		} else if action == "dislike" {
			voteValue = -1
		}

		// Send vote to The Cat API
		go func() {
			voteBody := map[string]interface{}{
				"image_id": imageID,
				"sub_id":   userID,
				"value":    voteValue,
			}

			bodyBytes, err := json.Marshal(voteBody)
			if err != nil {
				fmt.Println("Failed to marshal vote request body:", err)
				return
			}

			voteURL := "https://api.thecatapi.com/v1/votes"
			req, err := http.NewRequest("POST", voteURL, bytes.NewBuffer(bodyBytes))
			if err != nil {
				fmt.Println("Failed to create vote request:", err)
				return
			}
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("x-api-key", c.APIKey)

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				fmt.Println("Failed to send vote:", err)
				return
			}
			defer resp.Body.Close()

			body, _ := io.ReadAll(resp.Body)
			fmt.Println("Vote Response Status:", resp.Status)
			fmt.Println("Vote Response Body:", string(body))
		}()
	}

	// Fetch a new random cat image after action (if applicable)
	c.Get()

	c.Data["json"] = map[string]interface{}{
		"image_url": c.Data["ImageURL"],
		"image_id":  c.Data["ImageID"],
	}
	c.ServeJSON()
}
