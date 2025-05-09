package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/labstack/echo"

	// Import your server package with the correct module path
	myserver "github.com/Lxkota95/website/server/main"
)

// serverAdapter converts Lambda events to Echo requests and responses
func serverAdapter(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	// Create a new Echo instance
	e := echo.New()

	// Create a context for the request
	req, _ := http.NewRequest(request.HTTPMethod, request.Path, nil)

	// Add headers
	for key, value := range request.Headers {
		req.Header.Add(key, value)
	}

	// Add query parameters
	q := req.URL.Query()
	for key, value := range request.QueryStringParameters {
		q.Add(key, value)
	}
	req.URL.RawQuery = q.Encode()

	// Set the request body if present
	if request.Body != "" {
		req.Body = ioutil.NopCloser(strings.NewReader(request.Body))
	}

	// Create Echo context
	resp := echo.NewResponse(httptest.NewRecorder(), e)
	c := e.NewContext(req, resp)

	// Ensure we're in the right directory for content & template access
	cwd, _ := os.Getwd()
	fmt.Printf("Current working directory: %s\n", cwd)
	// fmt.Printf("Directory contents: %v\n", filepath.Glob("*"))

	// Call your server's handler
	if err := myserver.HandleRequest(c); err != nil {
		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Internal Server Error: " + err.Error(),
		}, nil
	}

	// Get the recorder from the response
	recorder := resp.Writer.(*httptest.ResponseRecorder)

	// Convert headers to the format expected by AWS
	headers := make(map[string]string)
	for key, values := range recorder.Header() {
		if len(values) > 0 {
			headers[key] = values[0]
		}
	}

	// Return the Lambda response
	return &events.APIGatewayProxyResponse{
		StatusCode: recorder.Code,
		Headers:    headers,
		Body:       recorder.Body.String(),
	}, nil
}

func main() {
	lambda.Start(serverAdapter)
}
