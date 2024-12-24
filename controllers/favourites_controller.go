package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/beego/beego/v2/core/config"
	"github.com/beego/beego/v2/server/web"
)

type FavoritesController struct {
	web.Controller
	APIKey string
}

type FavoriteResponse struct {
	ID        int    `json:"id"`
	ImageID   string `json:"image_id"`
	SubID     string `json:"sub_id"`
	CreatedAt string `json:"created_at"`
	Image     struct {
		ID  string `json:"id"`
		URL string `json:"url"`
	} `json:"image"`
}

func (c *FavoritesController) Prepare() {
	apiKey, err := config.String("api_key")
	if err != nil {
		c.Data["json"] = map[string]interface{}{"error": "Failed to load API key from configuration"}
		c.ServeJSON()
		return
	}
	c.APIKey = apiKey // Store the API key in the controller's field
}

// Get retrieves all favorites
func (c *FavoritesController) Get() {
	// Create HTTP request
	req, err := http.NewRequest("GET", "https://api.thecatapi.com/v1/favourites", nil)
	if err != nil {
		c.Data["json"] = map[string]interface{}{"error": "Failed to create request"}
		c.ServeJSON()
		return
	}

	req.Header.Set("x-api-key", c.APIKey)

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.Data["json"] = map[string]interface{}{"error": "Failed to fetch favorites"}
		c.ServeJSON()
		return
	}
	defer resp.Body.Close()

	// Parse response
	var favorites []FavoriteResponse
	if err := json.NewDecoder(resp.Body).Decode(&favorites); err != nil {
		c.Data["json"] = map[string]interface{}{"error": "Failed to parse response"}
		c.ServeJSON()
		return
	}

	c.Data["json"] = favorites
	c.ServeJSON()
}
