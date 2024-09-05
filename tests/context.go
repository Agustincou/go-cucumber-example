package tests

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
)

type _APIContainerCtxKey struct{} //testcontainers.Container
type _APIEndpointCtxKey struct{}  //string
type _responseBodyCtxKey struct{} //[]byte
type _responseCodeCtxKey struct{} //int

func tearDownScenario(ctx context.Context) error {
	apiContainer, ok := ctx.Value(_APIContainerCtxKey{}).(testcontainers.Container)
	if !ok || apiContainer == nil {
		return fmt.Errorf("API container not found in context")
	}

	if err := apiContainer.Terminate(ctx); err != nil {
		return fmt.Errorf("failed to terminate API container: %w", err)
	}

	log.Printf("API container stopped and removed")

	return nil
}

func givenAPIRunning(ctx context.Context) (context.Context, error) {
	// Start container
	apiContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			ExposedPorts: []string{"8080/tcp"},
			FromDockerfile: testcontainers.FromDockerfile{
				Context:       "../",
				Dockerfile:    "Dockerfile",
				PrintBuildLog: true,
				KeepImage:     false,
			},
		},
		Started: true,
	})

	if err != nil {
		log.Fatalf("Failed to start container: %v", err)

		return ctx, err
	}

	//containerInfo, _ := apiContainer.Inspect(ctx)

	var host string
	var port nat.Port
	maxRetries := 10
	for i := 0; i < maxRetries; i++ {
		host, err = apiContainer.Host(ctx)
		if err == nil {
			port, err = apiContainer.MappedPort(ctx, "8080")
			if err == nil {
				break
			}
		}
		log.Printf("Retrying to get container info (attempt %d/%d)...", i+1, maxRetries)
		time.Sleep(2 * time.Second)
	}

	apiEndpoint := fmt.Sprintf("http://%s:%s", host, port.Port())

	if err := waitForPingEndpoint(apiEndpoint); err != nil {
		return ctx, err
	}

	ctx = context.WithValue(ctx, _APIContainerCtxKey{}, apiContainer)
	ctx = context.WithValue(ctx, _APIEndpointCtxKey{}, apiEndpoint)

	return ctx, nil
}

func waitForPingEndpoint(apiEndpoint string) error {
	timeout := 10 * time.Second
	start := time.Now()

	for time.Since(start) < timeout {
		resp, err := http.Get(apiEndpoint + "/ping")
		if err != nil {
			log.Printf("Error checking /ping endpoint: %v \n", err)
		} else {
			defer resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				return nil // Endpoint disponible y correcto
			}
			fmt.Println("api response", resp.StatusCode)
		}
		time.Sleep(2 * time.Second)
	}

	return fmt.Errorf("service not available or incorrect response after: %s", timeout)
}

func whenReceiveRequest(ctx context.Context, method string, path string, body string) (context.Context, error) {
	apiEndpoint, ok := ctx.Value(_APIEndpointCtxKey{}).(string)

	if !ok || apiEndpoint == "" {
		return ctx, fmt.Errorf("API endpoint not found in context")
	}

	// Crear la URL completa con el path recibido
	url := apiEndpoint + path

	// Crear la request dependiendo del método recibido
	req, err := http.NewRequestWithContext(ctx, method, url, strings.NewReader(body))
	if err != nil {
		return ctx, fmt.Errorf("failed to create request: %w", err)
	}

	// Hacer la solicitud
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return ctx, fmt.Errorf("failed to make request: %w", err)
	}

	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return ctx, err
	}

	ctx = context.WithValue(ctx, _responseBodyCtxKey{}, responseBody)
	ctx = context.WithValue(ctx, _responseCodeCtxKey{}, resp.StatusCode)

	return ctx, nil
}

func thenResponseHTTPStatusCodeShouldBe(ctx context.Context, httpStatusCode int) error {
	respHttpStatusCode, ok := ctx.Value(_responseCodeCtxKey{}).(int)
	if !ok || httpStatusCode != respHttpStatusCode {
		return fmt.Errorf("expected code: %d, but receive: %d", httpStatusCode, respHttpStatusCode)
	}

	return nil
}

func thenResponseBodyShouldBe(ctx context.Context, body string) error {
	respBody, ok := ctx.Value(_responseBodyCtxKey{}).([]byte)
	if !ok || body != strings.TrimSpace(string(respBody)) {
		return fmt.Errorf("expected body: %s, but receive: %s", body, respBody)
	}

	return nil
}
