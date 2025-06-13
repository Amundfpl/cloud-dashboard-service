package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

// WriteErrorResponse sends a structured JSON error response with the given HTTP status code.
// It ensures the Content-Type is set and logs any failure to write the response.
func WriteErrorResponse(w http.ResponseWriter, errorMessage string, statusCode int) {
	w.Header().Set(HeaderContentType, ContentTypeJSON)
	w.WriteHeader(statusCode)

	response := map[string]string{KeyError: errorMessage}
	prettyJSON, _ := json.MarshalIndent(response, "", "  ")

	if _, writeErr := w.Write(prettyJSON); writeErr != nil {
		log.Printf(LogWriteErrorResponseFailed, writeErr)
	}
}

// WriteSuccessResponse sends a successful response with pretty-printed JSON,
// using the provided HTTP status code.
func WriteSuccessResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set(HeaderContentType, ContentTypeJSON)
	w.WriteHeader(statusCode)

	prettyJSON, _ := json.MarshalIndent(data, "", "  ")

	if _, writeErr := w.Write(prettyJSON); writeErr != nil {
		log.Printf(LogWriteSuccessResponseFailed, writeErr)
	}
}

// WriteIDResponse returns a simple JSON payload with only an "id" field.
// Commonly used for resource creation acknowledgments.
func WriteIDResponse(w http.ResponseWriter, id string, statusCode int) {
	w.Header().Set(HeaderContentType, ContentTypeJSON)
	w.WriteHeader(statusCode)

	response := map[string]string{KeyID: id}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf(LogWriteIDResponseFailed, err)
	}
}

// EnforceMethod ensures that the incoming HTTP request uses the expected method.
// If not, it writes a "method not allowed" response and returns false.
func EnforceMethod(w http.ResponseWriter, r *http.Request, expectedMethod string) bool {
	if r.Method != expectedMethod {
		WriteErrorResponse(w, ErrMethodNotAllowed, http.StatusMethodNotAllowed)
		return false
	}
	return true
}

// ExtractIDFromPath retrieves the final segment of the URL path, assuming a minimum
// number of segments. Returns an error if the structure is too short.
func ExtractIDFromPath(r *http.Request, expectedParts int) (string, error) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < expectedParts {
		return "", fmt.Errorf(ErrMissingPathID)
	}
	return parts[len(parts)-1], nil
}

// CloseBody safely closes an HTTP response body and logs any errors that occur.
// Useful for cleanup after performing an HTTP request.
func CloseBody(body io.Closer) {
	if err := body.Close(); err != nil {
		log.Printf(LogCloseBodyFailed, err)
	}
}

// LogInfo is a generic utility for printing info-level logs to stdout.
// Can be replaced with a structured logging system later.
func LogInfo(msg string) {
	fmt.Println(msg)
}
