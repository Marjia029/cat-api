package controllers

import (
	"encoding/json"
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

// Get method to fetch a random cat image and return it as JSON
func (c *VotingController) Get() {
	apiKey := "live_GWXcPdnWze27MNMJSjinKshtfsnVsi4EdrXfKUNhOmXsLakl5N7MwJCShLvC5Rxo"
	url := "https://api.thecatapi.com/v1/images/search"

	// Fetch random cat image
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		c.Data["json"] = map[string]interface{}{"error": "Failed to fetch cat image: " + err.Error()}
		c.ServeJSON()
		return
	}
	req.Header.Set("x-api-key", apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.Data["json"] = map[string]interface{}{"error": "Failed to fetch cat image: " + err.Error()}
		c.ServeJSON()
		return
	}
	defer resp.Body.Close()

	var images []VotingCatImage
	if err := json.NewDecoder(resp.Body).Decode(&images); err != nil {
		c.Data["json"] = map[string]interface{}{"error": "Failed to decode cat image response: " + err.Error()}
		c.ServeJSON()
		return
	}

	// Send the first image in the response
	if len(images) > 0 {
		c.Data["json"] = map[string]interface{}{
			"image_url": images[0].URL,
			"favorites": favorites,
		}
	} else {
		c.Data["json"] = map[string]interface{}{"error": "No image found"}
	}

	c.ServeJSON()
}

// Post method to handle like, dislike, and saving to favorites
func (c *VotingController) Post() {
	action := c.GetString("action")
	imageURL := c.GetString("image_url")

	// Handle favorite action
	if action == "favorite" {
		// Add the image to favorites
		favorites = append(favorites, imageURL)
	}

	// Fetch a new random cat image if like/dislike action
	if action == "like" || action == "dislike" || action == "favorite" {
		c.Get() // This will fetch a new random image and update `c.Data["ImageURL"]`
	}

	// Send the updated image URL and favorites list as a JSON response
	c.Data["json"] = map[string]interface{}{"image_url": c.Data["ImageURL"], "favorites": favorites}
	c.ServeJSON()
}
