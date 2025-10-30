// Package testutil provides shared utilities for API testing, particularly
// OpenAPI specification validation.
package testutil

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/getkin/kin-openapi/routers"
	"github.com/getkin/kin-openapi/routers/gorillamux"
	"github.com/stretchr/testify/require"
)

// APIRequest wraps an HTTP request with its method, path, and body
// Can be used for both OpenAPI validation and actual HTTP requests
type APIRequest struct {
	Method string
	Path   string
	Body   string // JSON string
}

// LoadOpenAPISpec loads and validates the OpenAPI specification from the given file path.
// Returns the parsed OpenAPI document and a router for matching requests to operations.
func LoadOpenAPISpec(t *testing.T, ctx context.Context, specPath string) (*openapi3.T, routers.Router) {
	t.Helper()

	// Load OpenAPI spec from file
	loader := &openapi3.Loader{Context: ctx, IsExternalRefsAllowed: true}
	doc, err := loader.LoadFromFile(specPath)
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

// EnsureReqMatchesSpecs validates an APIRequest against the OpenAPI spec.
// This should be called for unauthenticated requests.
func EnsureReqMatchesSpecs(t *testing.T, ctx context.Context, router routers.Router, apiReq APIRequest) {
	t.Helper()

	// Create HTTP request
	httpReq, err := http.NewRequest(apiReq.Method, apiReq.Path, strings.NewReader(apiReq.Body))
	require.NoError(t, err, "Failed to create HTTP request")
	httpReq.Header.Set("Content-Type", "application/json")

	ValidateHTTPRequest(t, ctx, router, httpReq)
}

// AuthenticatedRequestRespectsSpecs validates an authenticated APIRequest against the OpenAPI spec.
// This should be called for requests that require authentication with a JWT token.
func AuthenticatedRequestRespectsSpecs(t *testing.T, ctx context.Context, router routers.Router, apiReq APIRequest, authToken string) {
	t.Helper()

	// Create HTTP request with Authorization header
	httpReq, err := http.NewRequest(apiReq.Method, apiReq.Path, strings.NewReader(apiReq.Body))
	require.NoError(t, err, "Failed to create HTTP request")
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Token "+authToken)

	ValidateHTTPRequest(t, ctx, router, httpReq)
}

// ValidateHTTPRequest validates an http.Request against the OpenAPI spec.
// This is the core validation function used by RequestRespectsSpecs and AuthenticatedRequestRespectsSpecs.
func ValidateHTTPRequest(t *testing.T, ctx context.Context, router routers.Router, httpReq *http.Request) {
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

// EnsureRespMatchesSpecs validates an HTTP response against the OpenAPI spec.
// This should be called after receiving a response to ensure it matches the spec.
func EnsureRespMatchesSpecs(t *testing.T, ctx context.Context, router routers.Router, apiReq APIRequest, resp *http.Response) {
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
		Status:                 resp.StatusCode,
		Header:                 resp.Header,
	}
	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	responseValidationInput.SetBodyBytes(respBody)
	err = openapi3filter.ValidateResponse(ctx, responseValidationInput)
	require.NoError(t, err, "invalid response")
}

func EnsureRequestAndResponseMatchSpecs(t *testing.T, ctx context.Context, router routers.Router, apiReq APIRequest, resp *http.Response) {
	EnsureReqMatchesSpecs(t, ctx, router, apiReq)
	EnsureRespMatchesSpecs(t, ctx, router, apiReq, resp)
}

// CloneResponse makes a deep-ish copy of an *http.Response for testing.
// It reads and buffers the entire body, then replaces both the original
// and clone bodies with new readers so they can be read independently.
func CloneResponse(t testing.TB, resp *http.Response) *http.Response {
	t.Helper()

	// Read and copy body
	bodyBytes, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	require.NoError(t, resp.Body.Close())

	// Restore original body
	resp.Body = io.NopCloser(bytes.NewReader(bodyBytes))

	// Create a shallow copy of resp (safe for most testing use cases)
	clone := *resp

	// Assign a separate reader to the clone
	clone.Body = io.NopCloser(bytes.NewReader(bodyBytes))

	return &clone
}
