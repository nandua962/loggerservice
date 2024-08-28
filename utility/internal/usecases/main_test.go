package usecases

import (
	"os"
	"testing"
	"utility/internal/consts"

	"github.com/gin-gonic/gin"
	"gitlab.com/tuneverse/toolkit/core/logger"
)

// TestMain is the entry point for running tests in this test suite.
// It configures the testing environment, initializes the logger, and executes the test suite.
// The function sets Gin web framework to test mode, configures the logger with specified options,
// and then runs all the test functions. The exit code reflects the success or failure of the tests.
func TestMain(t *testing.M) {
	// Set Gin to test mode for optimized testing behavior.
	gin.SetMode(gin.TestMode)

	// Configure logger options, including service name, log level, and data inclusion in logs.
	clientOpt := &logger.ClientOptions{
		Service:             consts.AppName, // Service name.
		LogLevel:            "info",         // Log level.
		IncludeRequestDump:  false,          // Include request data in logs.
		IncludeResponseDump: false,          // Include response data in logs.
	}

	// Initialize the logger with the specified options.
	_ = logger.InitLogger(clientOpt)

	// Run all test functions and exit with an appropriate status code.
	os.Exit(t.Run())
}
