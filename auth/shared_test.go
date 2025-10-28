// Package auth_test provides integration tests for authentication endpoints.
// Tests validate request/response against OpenAPI specifications and use
// testcontainers for real Kratos integration.
//
// OpenAPI Validation Strategy:
//   - All valid requests should be validated against the spec using requestRespectsSpecs
//   - Intentionally malformed requests (invalid JSON, wrong format) may skip validation
//     since they're designed to test error handling before reaching the spec layer
//   - All responses should be validated using responseRespectSpecs
//
// Container Management:
// Note: Each test file currently starts its own Kratos container for isolation.
// For faster test execution, consider using TestMain to share a container:
//
//	func TestMain(m *testing.M) {
//	    ctx := context.Background()
//	    container, err := startKratosContainer(ctx)
//	    if err != nil {
//	        log.Fatal(err)
//	    }
//	    defer container.Terminate(ctx)
//	    // Store URLs in package variables
//	    os.Exit(m.Run())
//	}
//
// Pros: 5-10x faster test execution
// Cons: Reduced test isolation (state may leak between tests)
package auth_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/err0r500/go-realworld-clean/auth"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/getkin/kin-openapi/routers"
	"github.com/getkin/kin-openapi/routers/gorillamux"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// APIRequest wraps an HTTP request with its method, path, and body
// Can be used for both OpenAPI validation and actual HTTP requests
type APIRequest struct {
	Method string
	Path   string
	Body   string // JSON string
}

// setupTestServer creates a test server with all auth endpoints registered
func setupTestServer(t *testing.T, kratosContainer testcontainers.Container) *httptest.Server {
	t.Helper()
	ctx := t.Context()

	kratosURL, err := getKratosURL(ctx, kratosContainer)
	require.NoError(t, err, "Failed to get Kratos URL")

	kratosAdminURL, err := getKratosAdminURL(ctx, kratosContainer)
	require.NoError(t, err, "Failed to get Kratos admin URL")

	mux := http.NewServeMux()
	handler := auth.NewAuthHandlerWithAdmin(kratosURL, kratosAdminURL)
	handler.RegisterUser(mux)
	handler.Login(mux)
	handler.GetCurrentUser(mux)
	handler.UpdateUser(mux)

	// Wrap with auth middleware to handle JWT tokens
	return httptest.NewServer(handler.AuthMiddleware(mux))
}

func startKratosContainer(ctx context.Context) (testcontainers.Container, error) {
	// Add timeout for container startup
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "oryd/kratos:v1.3.0",
			ExposedPorts: []string{"4433/tcp", "4434/tcp"},
			Cmd:          []string{"serve", "-c", "/etc/config/kratos/kratos.yml", "--dev"},
			Files: []testcontainers.ContainerFile{
				{
					HostFilePath:      "./_kratos_config/kratos.yml",
					ContainerFilePath: "/etc/config/kratos/kratos.yml",
					FileMode:          0644,
				},
				{
					HostFilePath:      "./_kratos_config/identity.schema.json",
					ContainerFilePath: "/etc/config/kratos/identity.schema.json",
					FileMode:          0644,
				},
			},
			WaitingFor: wait.ForListeningPort("4433/tcp").WithStartupTimeout(30 * time.Second),
		},
		Started: true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to start Kratos container: %w", err)
	}

	return container, nil
}

func getKratosURL(ctx context.Context, container testcontainers.Container) (string, error) {
	host, err := container.Host(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get container host: %w", err)
	}

	port, err := container.MappedPort(ctx, "4433")
	if err != nil {
		return "", fmt.Errorf("failed to get container port: %w", err)
	}

	// External endpoint for host to access
	return fmt.Sprintf("http://%s:%s", host, port.Port()), nil
}

func getKratosAdminURL(ctx context.Context, container testcontainers.Container) (string, error) {
	host, err := container.Host(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get container host: %w", err)
	}

	port, err := container.MappedPort(ctx, "4434")
	if err != nil {
		return "", fmt.Errorf("failed to get container admin port: %w", err)
	}

	// External endpoint for host to access admin API
	return fmt.Sprintf("http://%s:%s", host, port.Port()), nil
}

// loadOpenAPISpec loads and validates the OpenAPI specification
func loadOpenAPISpec(t *testing.T, ctx context.Context) (*openapi3.T, routers.Router) {
	t.Helper()

	// Load OpenAPI spec from file
	loader := &openapi3.Loader{Context: ctx, IsExternalRefsAllowed: true}
	doc, err := loader.LoadFromFile("../_specs/openapi.yml")
	require.NoError(t, err, "Failed to load OpenAPI spec")

	// Override servers to match our test implementation
	// Our endpoints use /api prefix, so set server base path to /api
	doc.Servers = openapi3.Servers{
		&openapi3.Server{URL: "/api"},
	}

	// Validate the document structure
	err = doc.Validate(ctx)
	require.NoError(t, err, "OpenAPI document validation failed")

	// Create router for matching requests to operations
	router, err := gorillamux.NewRouter(doc)
	require.NoError(t, err, "Failed to create OpenAPI router")

	return doc, router
}

// requestRespectsSpecs validates an APIRequest against the OpenAPI spec
func requestRespectsSpecs(t *testing.T, ctx context.Context, router routers.Router, apiReq APIRequest) {
	t.Helper()

	// Create HTTP request
	httpReq, err := http.NewRequest(apiReq.Method, apiReq.Path, strings.NewReader(apiReq.Body))
	require.NoError(t, err, "Failed to create HTTP request")
	httpReq.Header.Set("Content-Type", "application/json")

	validateHTTPRequest(t, ctx, router, httpReq)
}

// authenticatedRequestRespectsSpecs validates an authenticated APIRequest against the OpenAPI spec
func authenticatedRequestRespectsSpecs(t *testing.T, ctx context.Context, router routers.Router, apiReq APIRequest, authToken string) {
	t.Helper()

	// Create HTTP request with Authorization header
	httpReq, err := http.NewRequest(apiReq.Method, apiReq.Path, strings.NewReader(apiReq.Body))
	require.NoError(t, err, "Failed to create HTTP request")
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Token "+authToken)

	validateHTTPRequest(t, ctx, router, httpReq)
}

// validateHTTPRequest validates an http.Request against the OpenAPI spec
func validateHTTPRequest(t *testing.T, ctx context.Context, router routers.Router, httpReq *http.Request) {
	t.Helper()

	// Find matching route
	route, pathParams, err := router.FindRoute(httpReq)
	require.NoError(t, err, "Route not found in OpenAPI spec")

	// Configure validation options with authentication handler
	options := &openapi3filter.Options{
		AuthenticationFunc: func(ctx context.Context, input *openapi3filter.AuthenticationInput) error {
			// Validate Token security scheme (apiKey in Authorization header)
			if input.SecuritySchemeName == "Token" {
				authHeader := input.RequestValidationInput.Request.Header.Get("Authorization")
				// Check that header has format "Token <value>"
				if !strings.HasPrefix(authHeader, "Token ") {
					return fmt.Errorf("invalid Authorization header format, expected 'Token <jwt>'")
				}
				// Structure is valid - actual token validation happens in middleware
				return nil
			}
			return fmt.Errorf("unknown security scheme: %s", input.SecuritySchemeName)
		},
	}

	// Validate request
	requestValidationInput := &openapi3filter.RequestValidationInput{
		Request:    httpReq,
		PathParams: pathParams,
		Route:      route,
		Options:    options,
	}

	err = openapi3filter.ValidateRequest(ctx, requestValidationInput)
	require.NoError(t, err, "invalid request")
}

// responseRespectSpecs validates an HTTP response against the OpenAPI spec
func responseRespectSpecs(t *testing.T, ctx context.Context, router routers.Router, apiReq APIRequest, statusCode int, headers http.Header, responseBody []byte) {
	t.Helper()

	// Create HTTP request for context
	httpReq, err := http.NewRequest(apiReq.Method, apiReq.Path, strings.NewReader(apiReq.Body))
	require.NoError(t, err, "Failed to create HTTP request")
	httpReq.Header.Set("Content-Type", "application/json")

	// Find matching route
	route, pathParams, err := router.FindRoute(httpReq)
	require.NoError(t, err, "Route not found in OpenAPI spec")

	// Create request validation input
	requestValidationInput := &openapi3filter.RequestValidationInput{
		Request:    httpReq,
		PathParams: pathParams,
		Route:      route,
	}

	// Validate response
	responseValidationInput := &openapi3filter.ResponseValidationInput{
		RequestValidationInput: requestValidationInput,
		Status:                 statusCode,
		Header:                 headers,
	}
	responseValidationInput.SetBodyBytes(responseBody)

	err = openapi3filter.ValidateResponse(ctx, responseValidationInput)
	require.NoError(t, err, "invalid response")
}

// sendAPIRequest sends an APIRequest to the server and returns the response
func sendAPIRequest(t *testing.T, baseURL string, apiReq APIRequest) *http.Response {
	t.Helper()

	// Create and send request
	url := baseURL + apiReq.Path
	req, err := http.NewRequest(apiReq.Method, url, strings.NewReader(apiReq.Body))
	require.NoError(t, err, "Failed to create HTTP request")
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err, "Failed to send HTTP request")

	return resp
}

// readBody reads and returns the response body
func readBody(t *testing.T, body io.Reader) []byte {
	t.Helper()
	bodyBytes, err := io.ReadAll(body)
	require.NoError(t, err, "Failed to read response body")
	return bodyBytes
}

// registerTestUser registers a new user and returns their token and password
// This reduces duplication across test files that need authenticated users
func registerTestUser(t *testing.T, baseURL, username, email string) (token string, password string) {
	t.Helper()

	password = uuid.NewString()
	registerReq := auth.RegisterReqBody{
		User: auth.RegisterUser{
			Username: username,
			Email:    email,
			Password: password,
		},
	}

	regBody, err := json.Marshal(registerReq)
	require.NoError(t, err, "Failed to marshal registration request")

	regResp, err := http.Post(baseURL+"/api/users", "application/json", strings.NewReader(string(regBody)))
	require.NoError(t, err, "Failed to send registration request")
	defer regResp.Body.Close()

	require.Equal(t, http.StatusCreated, regResp.StatusCode, "Registration should succeed")

	var regUserResp auth.UserResponse
	err = json.NewDecoder(regResp.Body).Decode(&regUserResp)
	require.NoError(t, err, "Failed to decode registration response")

	return regUserResp.User.Token, password
}

// sendAuthenticatedRequest sends an authenticated HTTP request and returns response with body
// This centralizes the pattern of sending requests with JWT tokens
func sendAuthenticatedRequest(t *testing.T, method, url, body, token string) (*http.Response, []byte) {
	t.Helper()

	req, err := http.NewRequest(method, url, strings.NewReader(body))
	require.NoError(t, err, "Failed to create HTTP request")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Token "+token)

	response, err := http.DefaultClient.Do(req)
	require.NoError(t, err, "Failed to send HTTP request")

	responseBody := readBody(t, response.Body)
	response.Body.Close()

	return response, responseBody
}
