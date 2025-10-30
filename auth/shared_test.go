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
	"github.com/err0r500/go-realworld-clean/driving/api/testutil"
	"github.com/getkin/kin-openapi/routers"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

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

// sendAPIRequest sends an APIRequest to the server and returns the response
func sendAPIRequest(t *testing.T, baseURL string, apiReq testutil.APIRequest) *http.Response {
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

	// Restore the body so it can be read again (e.g., for validation)
	response.Body = io.NopCloser(strings.NewReader(string(responseBody)))

	return response, responseBody
}

func toJSON(t *testing.T, any any) string {
	bodyJSON, err := json.Marshal(any)
	if err != nil {
		t.Fatal(err)
	}
	return string(bodyJSON)
}

// authTestContext manages test server, container, and OpenAPI validation
type authTestContext struct {
	t          *testing.T
	server     *httptest.Server
	container  testcontainers.Container
	specRouter routers.Router
}

// newAuthTestContext creates a new test context with Kratos container and server
func newAuthTestContext(t *testing.T) *authTestContext {
	t.Helper()

	kratosContainer, err := startKratosContainer(t.Context())
	require.NoError(t, err, "Failed to start Kratos container")

	server := setupTestServer(t, kratosContainer)
	_, specRouter := testutil.LoadOpenAPISpec(t, t.Context(), "../_specs/openapi.yml")

	return &authTestContext{
		t:          t,
		server:     server,
		container:  kratosContainer,
		specRouter: specRouter,
	}
}

// cleanup closes server and terminates container
func (ctx *authTestContext) cleanup() {
	ctx.server.Close()
	if err := ctx.container.Terminate(ctx.t.Context()); err != nil {
		ctx.t.Errorf("failed to stop kratos container: %v", err)
	}
}

// givenUserExists registers a user and returns their token and password
func (ctx *authTestContext) givenUserExists(username, email string) (token, password string) {
	ctx.t.Helper()
	return registerTestUser(ctx.t, ctx.server.URL, username, email)
}

// whenGettingCurrentUser sends GET /api/user request
func (ctx *authTestContext) whenGettingCurrentUser(token string) *http.Response {
	ctx.t.Helper()

	req, err := http.NewRequest("GET", ctx.server.URL+"/api/user", nil)
	require.NoError(ctx.t, err)
	req.Header.Set("Authorization", "Token "+token)

	response, err := http.DefaultClient.Do(req)
	require.NoError(ctx.t, err)

	return response
}

// whenGettingCurrentUserWithoutAuth sends GET /api/user without auth header
func (ctx *authTestContext) whenGettingCurrentUserWithoutAuth() *http.Response {
	ctx.t.Helper()

	req, err := http.NewRequest("GET", ctx.server.URL+"/api/user", nil)
	require.NoError(ctx.t, err)

	response, err := http.DefaultClient.Do(req)
	require.NoError(ctx.t, err)

	return response
}

// whenGettingCurrentUserWithHeader sends GET /api/user with custom auth header
func (ctx *authTestContext) whenGettingCurrentUserWithHeader(authHeader string) *http.Response {
	ctx.t.Helper()

	req, err := http.NewRequest("GET", ctx.server.URL+"/api/user", nil)
	require.NoError(ctx.t, err)
	req.Header.Set("Authorization", authHeader)

	response, err := http.DefaultClient.Do(req)
	require.NoError(ctx.t, err)

	return response
}

// whenLoggingInWith sends POST /api/users/login request
func (ctx *authTestContext) whenLoggingInWith(email, password string) *http.Response {
	ctx.t.Helper()

	loginBody := auth.LoginReqBody{
		User: auth.LoginUser{
			Email:    email,
			Password: password,
		},
	}

	bodyJSON, err := json.Marshal(loginBody)
	require.NoError(ctx.t, err)

	response, err := http.Post(ctx.server.URL+"/api/users/login", "application/json", strings.NewReader(string(bodyJSON)))
	require.NoError(ctx.t, err)

	return response
}

// whenLoggingInWithJSON sends POST /api/users/login with raw JSON
func (ctx *authTestContext) whenLoggingInWithJSON(jsonBody string) *http.Response {
	ctx.t.Helper()

	response, err := http.Post(ctx.server.URL+"/api/users/login", "application/json", strings.NewReader(jsonBody))
	require.NoError(ctx.t, err)

	return response
}

// whenRegisteringWith sends POST /api/users request
func (ctx *authTestContext) whenRegisteringWith(username, email, password string) *http.Response {
	ctx.t.Helper()

	reqBody := auth.RegisterReqBody{
		User: auth.RegisterUser{
			Username: username,
			Email:    email,
			Password: password,
		},
	}

	bodyJSON, err := json.Marshal(reqBody)
	require.NoError(ctx.t, err)

	response, err := http.Post(ctx.server.URL+"/api/users", "application/json", strings.NewReader(string(bodyJSON)))
	require.NoError(ctx.t, err)

	return response
}

// whenRegisteringWithJSON sends POST /api/users with raw JSON
func (ctx *authTestContext) whenRegisteringWithJSON(jsonBody string) *http.Response {
	ctx.t.Helper()

	response, err := http.Post(ctx.server.URL+"/api/users", "application/json", strings.NewReader(jsonBody))
	require.NoError(ctx.t, err)

	return response
}

// updateUserBuilder builds update user requests fluently
type updateUserBuilder struct {
	ctx    *authTestContext
	token  string
	fields auth.UpdateUserFields
}

// whenUpdatingUser starts building an update user request
func (ctx *authTestContext) whenUpdatingUser(token string) *updateUserBuilder {
	return &updateUserBuilder{
		ctx:   ctx,
		token: token,
	}
}

func (b *updateUserBuilder) withEmail(email string) *updateUserBuilder {
	b.fields.Email = &email
	return b
}

func (b *updateUserBuilder) withUsername(username string) *updateUserBuilder {
	b.fields.Username = &username
	return b
}

func (b *updateUserBuilder) withBio(bio string) *updateUserBuilder {
	b.fields.Bio = &bio
	return b
}

func (b *updateUserBuilder) withImage(image string) *updateUserBuilder {
	b.fields.Image = &image
	return b
}

func (b *updateUserBuilder) withPassword(password string) *updateUserBuilder {
	b.fields.Password = &password
	return b
}

func (b *updateUserBuilder) send() *http.Response {
	b.ctx.t.Helper()

	updateBody := auth.UpdateUserReqBody{User: b.fields}
	bodyJSON, err := json.Marshal(updateBody)
	require.NoError(b.ctx.t, err)

	req, err := http.NewRequest("PUT", b.ctx.server.URL+"/api/user", strings.NewReader(string(bodyJSON)))
	require.NoError(b.ctx.t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Token "+b.token)

	response, err := http.DefaultClient.Do(req)
	require.NoError(b.ctx.t, err)

	return response
}

// whenUpdatingUserWithoutAuth sends PUT /api/user without auth
func (ctx *authTestContext) whenUpdatingUserWithoutAuth(fields auth.UpdateUserFields) *http.Response {
	ctx.t.Helper()

	updateBody := auth.UpdateUserReqBody{User: fields}
	bodyJSON, err := json.Marshal(updateBody)
	require.NoError(ctx.t, err)

	req, err := http.NewRequest("PUT", ctx.server.URL+"/api/user", strings.NewReader(string(bodyJSON)))
	require.NoError(ctx.t, err)
	req.Header.Set("Content-Type", "application/json")

	response, err := http.DefaultClient.Do(req)
	require.NoError(ctx.t, err)

	return response
}

// whenUpdatingUserWithHeader sends PUT /api/user with custom auth header
func (ctx *authTestContext) whenUpdatingUserWithHeader(authHeader string, fields auth.UpdateUserFields) *http.Response {
	ctx.t.Helper()

	updateBody := auth.UpdateUserReqBody{User: fields}
	bodyJSON, err := json.Marshal(updateBody)
	require.NoError(ctx.t, err)

	req, err := http.NewRequest("PUT", ctx.server.URL+"/api/user", strings.NewReader(string(bodyJSON)))
	require.NoError(ctx.t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", authHeader)

	response, err := http.DefaultClient.Do(req)
	require.NoError(ctx.t, err)

	return response
}

// whenUpdatingUserWithJSON sends PUT /api/user with raw JSON
func (ctx *authTestContext) whenUpdatingUserWithJSON(token, jsonBody string) *http.Response {
	ctx.t.Helper()

	req, err := http.NewRequest("PUT", ctx.server.URL+"/api/user", strings.NewReader(jsonBody))
	require.NoError(ctx.t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Token "+token)

	response, err := http.DefaultClient.Do(req)
	require.NoError(ctx.t, err)

	return response
}

// Assertion helpers

func assertSuccessfulLogin(t *testing.T, resp *http.Response) auth.UserResponse {
	t.Helper()
	require.Equal(t, http.StatusOK, resp.StatusCode)

	bodyBytes := readBody(t, resp.Body)
	var userResp auth.UserResponse
	err := json.Unmarshal(bodyBytes, &userResp)
	require.NoError(t, err)
	require.NotEmpty(t, userResp.User.Token)

	// Restore body for later reads
	resp.Body = io.NopCloser(strings.NewReader(string(bodyBytes)))

	return userResp
}

func assertSuccessfulRegistration(t *testing.T, resp *http.Response) auth.UserResponse {
	t.Helper()
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	bodyBytes := readBody(t, resp.Body)
	var userResp auth.UserResponse
	err := json.Unmarshal(bodyBytes, &userResp)
	require.NoError(t, err)
	require.NotEmpty(t, userResp.User.Token)

	// Restore body for later reads
	resp.Body = io.NopCloser(strings.NewReader(string(bodyBytes)))

	return userResp
}

func assertSuccessfulUpdate(t *testing.T, resp *http.Response) auth.UserResponse {
	t.Helper()
	require.Equal(t, http.StatusOK, resp.StatusCode)

	bodyBytes := readBody(t, resp.Body)
	var userResp auth.UserResponse
	err := json.Unmarshal(bodyBytes, &userResp)
	require.NoError(t, err)
	require.NotEmpty(t, userResp.User.Token)

	// Restore body for later reads
	resp.Body = io.NopCloser(strings.NewReader(string(bodyBytes)))

	return userResp
}

func assertSuccessfulGetCurrentUser(t *testing.T, resp *http.Response) auth.UserResponse {
	t.Helper()
	require.Equal(t, http.StatusOK, resp.StatusCode)

	bodyBytes := readBody(t, resp.Body)
	var userResp auth.UserResponse
	err := json.Unmarshal(bodyBytes, &userResp)
	require.NoError(t, err)
	require.NotEmpty(t, userResp.User.Token)

	// Restore body for later reads
	resp.Body = io.NopCloser(strings.NewReader(string(bodyBytes)))

	return userResp
}

func assertUnauthorized(t *testing.T, resp *http.Response) {
	t.Helper()
	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func assertBadRequest(t *testing.T, resp *http.Response) {
	t.Helper()
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func assertUnprocessableEntity(t *testing.T, resp *http.Response) {
	t.Helper()
	require.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
}

// ensureResponseMatchesSpec validates response against OpenAPI spec
func (ctx *authTestContext) ensureResponseMatchesSpec(t *testing.T, apiReq testutil.APIRequest, resp *http.Response) {
	t.Helper()
	cloned := testutil.CloneResponse(t, resp)
	testutil.EnsureRespMatchesSpecs(t, t.Context(), ctx.specRouter, apiReq, cloned)
}

// ensureRequestMatchesSpec validates request against OpenAPI spec
func (ctx *authTestContext) ensureRequestMatchesSpec(t *testing.T, apiReq testutil.APIRequest) {
	t.Helper()
	testutil.EnsureReqMatchesSpecs(t, t.Context(), ctx.specRouter, apiReq)
}
