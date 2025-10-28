package auth_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
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
	"github.com/stretchr/testify/assert"
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

func TestRegisterUser(t *testing.T) {
	kratosContainer := startKratosContainer(t.Context())
	defer func() {
		if err := kratosContainer.Terminate(t.Context()); err != nil {
			log.Fatal("failed to stop kratos container", err)
		}
	}()

	server := setupTestServer(t, getKratosURL(t.Context(), kratosContainer))
	defer server.Close()

	// Load OpenAPI spec for validation
	_, router := loadOpenAPISpec(t, t.Context())
	apiReq := APIRequest{
		Method: "POST",
		Path:   "/api/users",
	}

	t.Run("valid registration request", func(t *testing.T) {
		// given - create request using typed struct
		reqBody := auth.RegisterReqBody{
			User: auth.RegisterUser{
				Username: "testuser",
				Email:    "test@example.com",
				Password: uuid.NewString(),
			},
		}

		// Marshal to JSON for APIRequest
		bodyJSON, err := json.Marshal(reqBody)
		require.NoError(t, err)
		apiReq.Body = string(bodyJSON)

		requestRespectsSpecs(t, t.Context(), router, apiReq)

		// when
		response := sendAPIRequest(t, server.URL, apiReq)
		defer response.Body.Close()

		// then
		assert.Equal(t, http.StatusCreated, response.StatusCode)

		responseBody := readBody(t, response.Body)
		responseRespectSpecs(t, t.Context(), router, apiReq, response.StatusCode, response.Header, responseBody)

		// Parse response and assert on values
		var userResp auth.UserResponse
		require.NoError(t, json.Unmarshal(responseBody, &userResp))

		assert.Equal(t, reqBody.User.Username, userResp.User.Username)
		assert.Equal(t, reqBody.User.Email, userResp.User.Email)
		assert.NotEmpty(t, userResp.User.Token, "should return a JWT token")
		assert.Equal(t, "", userResp.User.Bio)
		assert.Equal(t, "", userResp.User.Image)
	})

	t.Run("invalid request - missing required fields", func(t *testing.T) {
		testCases := []struct {
			name        string
			requestBody string
		}{
			{
				name:        "missing username",
				requestBody: `{"user": {"email": "test@example.com", "password": "password123"}}`,
			},
			{
				name:        "missing email",
				requestBody: `{"user": {"username": "testuser", "password": "password123"}}`,
			},
			{
				name:        "missing password",
				requestBody: `{"user": {"username": "testuser", "email": "test@example.com"}}`,
			},
			{
				name:        "empty user object",
				requestBody: `{"user": {}}`,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// given
				apiReq.Body = tc.requestBody

				// when
				response := sendAPIRequest(t, server.URL, apiReq)
				defer response.Body.Close()

				// then
				assert.Equal(t, http.StatusUnprocessableEntity, response.StatusCode)
				responseRespectSpecs(t, t.Context(), router, apiReq, response.StatusCode, response.Header, readBody(t, response.Body))
			})
		}
	})

	t.Run("invalid request - malformed JSON", func(t *testing.T) {
		// given
		apiReq.Body = `{"user": {"username": "testuser", "email": "test@example.com"`

		// when
		response := sendAPIRequest(t, server.URL, apiReq)
		defer response.Body.Close()

		// then
		assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	})
}

func TestLogin(t *testing.T) {
	kratosContainer := startKratosContainer(t.Context())
	defer func() {
		if err := kratosContainer.Terminate(t.Context()); err != nil {
			log.Fatal("failed to stop kratos container", err)
		}
	}()

	server := setupTestServer(t, getKratosURL(t.Context(), kratosContainer))
	defer server.Close()

	// Load OpenAPI spec for validation
	_, router := loadOpenAPISpec(t, t.Context())
	apiReq := APIRequest{
		Method: "POST",
		Path:   "/api/users/login",
	}

	password := uuid.NewString()
	registerReq := auth.RegisterReqBody{
		User: auth.RegisterUser{
			Username: "loginuser",
			Email:    "login@example.com",
			Password: password,
		},
	}
	regBody, _ := json.Marshal(registerReq)
	regResp, err := http.Post(server.URL+"/api/users", "application/json", strings.NewReader(string(regBody)))
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, regResp.StatusCode)
	regResp.Body.Close()

	t.Run("valid login request", func(t *testing.T) {
		// given - login with registered credentials
		loginBody := auth.LoginReqBody{
			User: auth.LoginUser{
				Email:    registerReq.User.Email,
				Password: registerReq.User.Password,
			},
		}

		bodyJSON, err := json.Marshal(loginBody)
		require.NoError(t, err)
		apiReq.Body = string(bodyJSON)

		requestRespectsSpecs(t, t.Context(), router, apiReq)

		// when
		response := sendAPIRequest(t, server.URL, apiReq)
		defer response.Body.Close()

		// then
		assert.Equal(t, http.StatusOK, response.StatusCode)

		responseBody := readBody(t, response.Body)
		responseRespectSpecs(t, t.Context(), router, apiReq, response.StatusCode, response.Header, responseBody)

		// Parse response and assert on values
		var userResp auth.UserResponse
		require.NoError(t, json.Unmarshal(responseBody, &userResp))

		assert.Equal(t, loginBody.User.Email, userResp.User.Email)
		assert.NotEmpty(t, userResp.User.Token, "should return a JWT token")
		assert.Equal(t, registerReq.User.Username, userResp.User.Username)
	})

	t.Run("invalid credentials", func(t *testing.T) {
		testCases := []struct {
			name     string
			email    string
			password string
		}{
			{
				name:     "non-existent user",
				email:    "nonexistent@example.com",
				password: registerReq.User.Password,
			},
			{
				name:     "wrong password",
				email:    registerReq.User.Email,
				password: "wrongpassword",
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// given
				loginBody := auth.LoginReqBody{
					User: auth.LoginUser{
						Email:    tc.email,
						Password: tc.password,
					},
				}
				bodyJSON, _ := json.Marshal(loginBody)
				apiReq.Body = string(bodyJSON)

				// when
				response := sendAPIRequest(t, server.URL, apiReq)
				defer response.Body.Close()

				// then
				assert.Equal(t, http.StatusUnauthorized, response.StatusCode)
			})
		}
	})

	t.Run("invalid request - missing required fields", func(t *testing.T) {
		testCases := []struct {
			name        string
			requestBody string
		}{
			{
				name:        "missing email",
				requestBody: `{"user": {"password": "password123"}}`,
			},
			{
				name:        "missing password",
				requestBody: `{"user": {"email": "test@example.com"}}`,
			},
			{
				name:        "empty user object",
				requestBody: `{"user": {}}`,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// given
				apiReq.Body = tc.requestBody

				// when
				response := sendAPIRequest(t, server.URL, apiReq)
				defer response.Body.Close()

				// then
				assert.Equal(t, http.StatusUnprocessableEntity, response.StatusCode)
				responseRespectSpecs(t, t.Context(), router, apiReq, response.StatusCode, response.Header, readBody(t, response.Body))
			})
		}
	})

	t.Run("invalid request - malformed JSON", func(t *testing.T) {
		// given
		apiReq.Body = `{"user": {"email": "test@example.com"`

		// when
		response := sendAPIRequest(t, server.URL, apiReq)
		defer response.Body.Close()

		// then
		assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	})

}

// setupTestServer creates a test server with all auth endpoints registered
func setupTestServer(t *testing.T, kratosURL string) *httptest.Server {
	t.Helper()
	mux := http.NewServeMux()
	handler := auth.NewAuthHandler(kratosURL)
	handler.RegisterUser(mux)
	handler.Login(mux)
	return httptest.NewServer(mux)
}

func startKratosContainer(ctx context.Context) testcontainers.Container {
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
		log.Fatal(err.Error())
	}

	// External endpoint for host to access
	return container
}

func getKratosURL(ctx context.Context, container testcontainers.Container) string {
	host, err := container.Host(ctx)
	if err != nil {
		log.Fatalf("Failed to get container host: %v", err)
	}

	port, err := container.MappedPort(ctx, "4433")
	if err != nil {
		log.Fatalf("Failed to get container port: %v", err)
	}

	// External endpoint for host to access
	return fmt.Sprintf("http://%s:%s", host, port.Port())
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

	// Find matching route
	route, pathParams, err := router.FindRoute(httpReq)
	require.NoError(t, err, "Route not found in OpenAPI spec")

	// Validate request
	requestValidationInput := &openapi3filter.RequestValidationInput{
		Request:    httpReq,
		PathParams: pathParams,
		Route:      route,
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
