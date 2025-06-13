package utils

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

// CurrentTimestamp returns the current local time formatted using the layout "yyyyMMdd HH:mm".
// This format is compact yet readable for logs or Firestore fields.
func CurrentTimestamp() string {
	return time.Now().Format(TimestampLayout)
}

// DefaultCredentialsPath returns the location of the Firebase service account credentials used for testing.
// It supports two mechanisms:
//  1. Uses the GO_FIREBASE_CREDENTIALS environment variable if set.
//  2. Otherwise, falls back to a static path under "/credentials/test-serviceAccountKey.json"
//     relative to the root of the project.
func DefaultCredentialsPath() string {
	// Use environment override if available
	if customPath := os.Getenv(CredentialsEnvVar); customPath != "" {
		return customPath
	}

	// Resolve path based on source file location
	_, currentFilePath, _, _ := runtime.Caller(0)
	projectRoot := filepath.Dir(filepath.Dir(currentFilePath)) // Navigate up 2 directories
	defaultPath := filepath.Join(projectRoot, CredentialsDir, TestCredentialsFile)

	log.Printf(LogFallbackCredentialUsed, defaultPath)
	return defaultPath
}
