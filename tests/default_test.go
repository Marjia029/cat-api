package tests

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"
	"github.com/stretchr/testify/assert"
	"myproject/controllers"
)

// setupViewsPath creates the views directory and template file
func setupViewsPath() string {
	pwd, _ := os.Getwd()
	viewsPath := filepath.Join(filepath.Dir(pwd), "tests/views")
	err := os.MkdirAll(viewsPath, 0755)
	if err != nil {
		panic(err)
	}

	// Create a simple template file for testing
	templatePath := filepath.Join(viewsPath, "index.tpl")
	if _, err := os.Stat(templatePath); os.IsNotExist(err) {
		err = os.WriteFile(templatePath, []byte(`<!DOCTYPE html>
<html>
<head>
    <title>Test Template</title>
</head>
<body>
    <h1>Test Template</h1>
</body>
</html>`), 0644)
		if err != nil {
			panic(err)
		}
	}
	return viewsPath
}

func TestMainControllerGet(t *testing.T) {
	viewsPath := setupViewsPath()
	defer os.RemoveAll(viewsPath)

	// Set the view path globally for the application
	beego.BConfig.WebConfig.ViewsPath = viewsPath

	tests := []struct {
		name           string
		expectedCode   int
		expectedTpl    string
		setupTestCase  func()
		validateResult func(*httptest.ResponseRecorder)
	}{
		{
			name:         "Basic Get Request",
			expectedCode: http.StatusOK,
			expectedTpl:  "index.tpl",
			setupTestCase: func() {
			},
			validateResult: func(w *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, w.Code)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupTestCase != nil {
				tt.setupTestCase()
			}

			w := httptest.NewRecorder()
			r, _ := http.NewRequest("GET", "/", nil)

			// Create a new Beego context
			ctx := context.NewContext()
			ctx.Reset(w, r)

			// Initialize controller
			c := &controllers.MainController{}
			c.Init(ctx, "", "", nil)

			// Call the Get method
			c.Get()

			// Validate the template name
			assert.Equal(t, tt.expectedTpl, c.TplName)

			// Validate the result if a check function is provided
			if tt.validateResult != nil {
				tt.validateResult(w)
			}
		})
	}
}
