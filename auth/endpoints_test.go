package auth_test

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/err0r500/go-realworld-clean/auth"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestHello(t *testing.T) {
	kratosContainer := startKratosContainer(t.Context())
	defer func() {
		if err := kratosContainer.Terminate(t.Context()); err != nil {
			log.Fatal("failed to stop kratos container", err)
		}
	}()

	mux := http.NewServeMux()
	auth.
		NewAuthHandler(getKratosURL(t.Context(), kratosContainer)).
		RegisterUser(mux)

	server := httptest.NewServer(mux)
	defer server.Close()

	requestBody := strings.NewReader(`{"user": {"username": "testuser", "email": "test@example.com", "password": "2349874389lsediosdifu"}}`)
	response, err := http.Post(fmt.Sprintf("%v/api/users", server.URL), "application/json", requestBody)
	if err != nil {
		t.Errorf("failed to [/help] %v", err)
	}

	assert.Equal(t, response.StatusCode, http.StatusCreated)
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
