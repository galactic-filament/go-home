package Util

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
)

// ErrorResponse - error response body in json
type ErrorResponse struct {
	Error string `json:"error"`
}

// ValidateError - error messages from ValidateEnvironment
type ValidateError struct {
	Messages []string
}

// WriteJSONErrorResponse - writes error out in json
func WriteJSONErrorResponse(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	errResponse := ErrorResponse{Error: err.Error()}
	if err := json.NewEncoder(w).Encode(errResponse); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Could not encode error response body")
		return
	}
	return
}

// ValidateEnvironment - validates the environment
func ValidateEnvironment() (envVars map[string]string, errorMessages []string) {
	envVars = map[string]string{}

	// validating that env vars are available
	envVarNames := []string{
		"APP_PORT",
		"APP_LOG_DIR",
		"DATABASE_HOST",
	}
	for _, name := range envVarNames {
		envVars[name] = os.Getenv(name)
	}
	missingEnvVars := []string{}
	for key, value := range envVars {
		if len(value) == 0 {
			missingEnvVars = append(missingEnvVars, key)
		}
	}
	if len(missingEnvVars) > 0 {
		messages := []string{}
		for _, key := range missingEnvVars {
			messages = append(messages, fmt.Sprintf("%s could not be found", key))
		}

		return nil, messages
	}

	// validating that the database port is accessible
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:5432", envVars["DATABASE_HOST"]))
	if err != nil {
		return nil, []string{fmt.Sprintf("Could not connect to %s", envVars["DATABASE_HOST"])}
	}
	if err = conn.Close(); err != nil {
		return nil, []string{fmt.Sprintf("Could not close connection to %s", envVars["DATABASE_HOST"])}
	}

	// validating that the log dir exists
	_, err = os.Stat(envVars["APP_LOG_DIR"])
	if err != nil {
		msg := fmt.Sprintf("Could not stat log dir %s", err.Error())
		if os.IsNotExist(err) {
			msg = fmt.Sprintf("%s log dir does not exist", envVars["APP_LOG_DIR"])
		}

		return nil, []string{msg}
	}

	return envVars, nil
}
