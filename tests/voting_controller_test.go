package tests

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"
	"github.com/stretchr/testify/assert"

	"myproject/controllers"
)

func init() {
	// Create conf directory if it doesn't exist
	err := os.MkdirAll("conf", 0755)
	if err != nil {
		panic(err)
	}

	// Create app.conf file
	confPath := filepath.Join("conf", "app.conf")
	err = os.WriteFile(confPath, []byte("api_key = test_api_key"), 0644)
	if err != nil {
		panic(err)
	}

	// Initialize Beego configuration
	err = beego.LoadAppConfig("ini", confPath)
	if err != nil {
		panic(err)
	}
}

// MockResponse represents a mock HTTP response
type MockResponse struct {
	StatusCode int
	Body       string
}

// Custom HTTP Transport for mocking
type MockTransport struct {
	responses map[string]MockResponse
}

func NewMockTransport() *MockTransport {
	return &MockTransport{
		responses: make(map[string]MockResponse),
	}
}

func (t *MockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Use mock data or return a valid JSON response if the path exists
	response, exists := t.responses[req.URL.Path]
	if !exists {
		response = MockResponse{
			StatusCode: http.StatusOK,
			Body:       `{"image_url":"http://example.com/cat.jpg", "image_id":"test123"}`,
		}
	}

	return &http.Response{
		StatusCode: response.StatusCode,
		Body:       io.NopCloser(strings.NewReader(response.Body)),
		Header:     make(http.Header),
	}, nil
}

// Create a new context for testing
func createTestContext(method, path string) (*context.Context, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(method, path, nil)
	ctx := context.NewContext()
	ctx.Reset(w, r)
	ctx.Request = r
	ctx.Output = context.NewOutput()
	ctx.Output.Context = ctx
	ctx.Output.Status = 200
	return ctx, w
}

func initController(ctx *context.Context) *controllers.VotingController {
	controller := &controllers.VotingController{}
	controller.Init(ctx, "", "", nil)
	controller.Data = make(map[interface{}]interface{})
	controller.APIKey = "test_api_key"
	return controller
}

func TestVotingControllerGet(t *testing.T) {
	// Set up mock transport (no mock data for validation)
	mockTransport := NewMockTransport()
	http.DefaultClient.Transport = mockTransport

	// Create a test context
	ctx, w := createTestContext("GET", "/voting")

	// Initialize the test controller
	controller := initController(ctx)

	// Call the Get method
	controller.Get()

	// Check if the response status code is OK (200)
	assert.Equal(t, http.StatusOK, w.Code)

	// Ensure that the response body is not empty
	assert.NotEmpty(t, w.Body.String())
}

func TestVotingControllerPost(t *testing.T) {
	tests := []struct {
		name           string
		action         string
		imageID        string
		expectedStatus int
	}{
		{
			name:           "Like Action",
			action:         "like",
			imageID:        "test123",
			expectedStatus: 200,
		},
		{
			name:           "Dislike Action",
			action:         "dislike",
			imageID:        "test123",
			expectedStatus: 200,
		},
		{
			name:           "Favorite Action",
			action:         "favorite",
			imageID:        "test123",
			expectedStatus: 200,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up mock transport
			mockTransport := NewMockTransport()
			mockTransport.responses["/v1/votes"] = MockResponse{
				StatusCode: tt.expectedStatus,
				Body:       `{"message":"SUCCESS"}`,
			}
			mockTransport.responses["/v1/favourites"] = MockResponse{
				StatusCode: tt.expectedStatus,
				Body:       `{"message":"SUCCESS"}`,
			}
			http.DefaultClient.Transport = mockTransport

			// Create test context
			ctx, w := createTestContext("POST", "/voting")

			// Initialize controller
			controller := initController(ctx)

			// Set form values
			ctx.Input.SetParam("action", tt.action)
			ctx.Input.SetParam("image_id", tt.imageID)

			// Call Post method
			controller.Post()

			// Check for a valid response
			assert.Equal(t, http.StatusOK, w.Code)

			// Ensure that the response body is not empty
			assert.NotEmpty(t, w.Body.String())

			// Parse response
			var result map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &result)
			if err != nil {
				t.Fatalf("Failed to unmarshal response: %v", err)
			}

			// Verify response contains image data
			assert.Contains(t, result, "image_url")
			assert.Contains(t, result, "image_id")
		})
	}
}

func TestVotingControllerPrepare(t *testing.T) {
	// Create test context
	ctx, _ := createTestContext("GET", "/voting")

	// Initialize controller
	controller := initController(ctx)

	// Call Prepare method
	controller.Prepare()

	// Assert API key was loaded
	assert.Equal(t, "test_api_key", controller.APIKey)
}

func TestMain(m *testing.M) {
	// Run the tests
	code := m.Run()

	// Clean up
	os.RemoveAll("conf")

	os.Exit(code)
}
