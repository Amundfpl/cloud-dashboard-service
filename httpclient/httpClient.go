package httpclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/amundfpl/Assignment-2/utils"
	"io"
	"net/http"
	"time"
)

// Client provides a wrapper around http.Client with convenience methods.
type Client struct {
	httpClient *http.Client
}

// NewClient initializes an HTTP client with a timeout to prevent hanging requests.
func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// Get performs a GET request and returns the response body as bytes.
// Returns an error if the request fails or returns a non-200 status.
func (c *Client) Get(url string) ([]byte, error) {
	resp, reqErr := c.httpClient.Get(url)
	if reqErr != nil {
		return nil, fmt.Errorf(utils.ErrHTTPGetFailed, reqErr)
	}
	defer utils.CloseBody(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(utils.ErrHTTPGetStatus, resp.StatusCode)
	}

	bodyBytes, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return nil, fmt.Errorf(utils.ErrHTTPReadBody, readErr)
	}

	return bodyBytes, nil
}

// GetStatusCode performs a GET request and returns only the status code.
func (c *Client) GetStatusCode(url string) (int, error) {
	resp, reqErr := c.httpClient.Get(url)
	if reqErr != nil {
		return 0, fmt.Errorf(utils.ErrHTTPGetFailed, reqErr)
	}
	defer utils.CloseBody(resp.Body)

	return resp.StatusCode, nil
}

// Post sends a POST request with a JSON payload and returns the response body.
// Errors if the request fails or returns a non-200 response.
func (c *Client) Post(url string, body map[string]string) ([]byte, error) {
	requestBody, marshalErr := json.Marshal(body)
	if marshalErr != nil {
		return nil, fmt.Errorf(utils.ErrHTTPPostMarshal, marshalErr)
	}

	resp, postErr := c.httpClient.Post(url, utils.ContentTypeJSON, bytes.NewBuffer(requestBody))
	if postErr != nil {
		return nil, fmt.Errorf(utils.ErrHTTPPostFailed, postErr)
	}
	defer utils.CloseBody(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(utils.ErrHTTPPostStatus, resp.StatusCode)
	}

	responseBytes, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return nil, fmt.Errorf(utils.ErrHTTPReadBody, readErr)
	}

	return responseBytes, nil
}
