// File: tests/voting_controller_test.go
package tests

import (
    "encoding/json"
    "fmt"
    "net/http"
    "net/http/httptest"
    "testing"

    beego "github.com/beego/beego/v2/server/web"
    "github.com/beego/beego/v2/server/web/context"
    "github.com/stretchr/testify/assert"
    
    "myproject/controllers"
)

func init() {
    // Initialize the router for testing
    beego.Router("/voting", &controllers.VotingController{})
}

// setupMockCatAPI creates and configures a mock Cat API server
func setupMockCatAPI(response string) *httptest.Server {
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(response))
    }))
    return server
}

// createTestRequest creates a new HTTP request for testing
func createTestRequest(method, path string) (*http.Request, *httptest.ResponseRecorder) {
    r, _ := http.NewRequest(method, path, nil)
    r.Header.Set("Content-Type", "application/json")
    return r, httptest.NewRecorder()
}

func TestVotingControllerGet(t *testing.T) {
    tests := []struct {
        name           string
        mockResponse   string
        expectedImage  string
        expectedError  bool
    }{
        {
            name:           "successful image fetch",
            mockResponse:   `[{"url": "http://example.com/cat.jpg"}]`,
            expectedImage:  "http://example.com/cat.jpg",
            expectedError:  false,
        },
        {
            name:           "empty response",
            mockResponse:   `[]`,
            expectedImage:  "",
            expectedError:  true,
        },
        {
            name:           "invalid JSON",
            mockResponse:   `invalid json`,
            expectedImage:  "",
            expectedError:  true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Setup mock API
            mockServer := setupMockCatAPI(tt.mockResponse)
            defer mockServer.Close()

            // Create test request
            r, w := createTestRequest("GET", "/voting")
            
            // Create context
            ctx := context.NewContext()
            ctx.Reset(w, r)
            
            // Create and initialize controller
            controller := &controllers.VotingController{}
            controller.Init(ctx, "", "", controller)

            // Execute Get method
            controller.Get()

            // Parse response
            var response map[string]interface{}
            err := json.Unmarshal(w.Body.Bytes(), &response)
            assert.NoError(t, err)

            if tt.expectedError {
                assert.Contains(t, response, "error")
            } else {
                assert.Equal(t, tt.expectedImage, response["image_url"])
                assert.Contains(t, response, "favorites")
            }
        })
    }
}

func TestVotingControllerPost(t *testing.T) {
    tests := []struct {
        name           string
        action         string
        imageURL       string
        expectedFavLen int
    }{
        {
            name:           "add to favorites",
            action:         "favorite",
            imageURL:       "http://example.com/cat1.jpg",
            expectedFavLen: 1,
        },
        {
            name:           "like action",
            action:         "like",
            imageURL:       "http://example.com/cat2.jpg",
            expectedFavLen: 0,
        },
        {
            name:           "dislike action",
            action:         "dislike",
            imageURL:       "http://example.com/cat3.jpg",
            expectedFavLen: 0,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Setup mock API
            mockServer := setupMockCatAPI(`[{"url": "http://example.com/new-cat.jpg"}]`)
            defer mockServer.Close()

            // Create test request
            r, w := createTestRequest("POST", "/voting")
            
            // Create context
            ctx := context.NewContext()
            ctx.Reset(w, r)
            
            // Set the parameters
            ctx.Input.SetParam("action", tt.action)
            ctx.Input.SetParam("image_url", tt.imageURL)

            // Create and initialize controller
            controller := &controllers.VotingController{}
            controller.Init(ctx, "", "", controller)

            // Execute Post method
            controller.Post()

            // Parse response
            var response map[string]interface{}
            err := json.Unmarshal(w.Body.Bytes(), &response)
            assert.NoError(t, err)

            // Verify response
            assert.Contains(t, response, "image_url")
            assert.Contains(t, response, "favorites")
            
            favoritesResponse, ok := response["favorites"].([]interface{})
            assert.True(t, ok)
            assert.Equal(t, tt.expectedFavLen, len(favoritesResponse))

            if tt.action == "favorite" {
                found := false
                for _, fav := range favoritesResponse {
                    if fav.(string) == tt.imageURL {
                        found = true
                        break
                    }
                }
                assert.True(t, found, "Expected image URL not found in favorites")
            }
        })
    }
}

func TestConcurrentFavorites(t *testing.T) {
    const numConcurrent = 10
    done := make(chan bool)

    for i := 0; i < numConcurrent; i++ {
        go func(index int) {
            // Create test request
            r, w := createTestRequest("POST", "/voting")
            
            // Create context
            ctx := context.NewContext()
            ctx.Reset(w, r)
            
            // Set parameters
            imageURL := fmt.Sprintf("http://example.com/cat%d.jpg", index)
            ctx.Input.SetParam("action", "favorite")
            ctx.Input.SetParam("image_url", imageURL)

            // Create and initialize controller
            controller := &controllers.VotingController{}
            controller.Init(ctx, "", "", controller)
            
            controller.Post()
            done <- true
        }(i)
    }

    // Wait for all goroutines to complete
    for i := 0; i < numConcurrent; i++ {
        <-done
    }

    // Create final request to check results
    r, w := createTestRequest("GET", "/voting")
    ctx := context.NewContext()
    ctx.Reset(w, r)
    
    controller := &controllers.VotingController{}
    controller.Init(ctx, "", "", controller)
    controller.Get()

    // Parse final response
    var response map[string]interface{}
    err := json.Unmarshal(w.Body.Bytes(), &response)
    assert.NoError(t, err)
    
    favoritesResponse := response["favorites"].([]interface{})
    assert.Equal(t, numConcurrent, len(favoritesResponse))
}