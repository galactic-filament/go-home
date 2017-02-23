package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	// validating that env vars are available
	envVars := map[string]string{}
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
		for _, key := range missingEnvVars {
			fmt.Println(fmt.Sprintf("%s could not be found", key))
		}

		os.Exit(1)
	}

	// validating that the database port is accessible
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:5432", envVars["DATABASE_HOST"]))
	if err != nil {
		fmt.Println(fmt.Sprintf("Could not connect to %s", envVars["DATABASE_HOST"]))
		os.Exit(1)
	}
	if err = conn.Close(); err != nil {
		fmt.Println(fmt.Sprintf("Could not close connection to %s", envVars["DATABASE_HOST"]))
		os.Exit(1)
	}

	// validating that the log dir exists
	_, err = os.Stat(envVars["APP_LOG_DIR"])
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println(fmt.Sprintf("%s log dir does not exist", envVars["APP_LOG_DIR"]))
			os.Exit(1)
		}

		fmt.Println(fmt.Sprintf("Could not stat log dir %s", err.Error()))
		os.Exit(1)
	}

	os.Exit(0)
}
