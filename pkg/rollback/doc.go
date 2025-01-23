/*
Package rollback provides automated application version rollback functionality.

It supports multiple deployment platforms through a strategy interface and includes
features like retry mechanisms, health checks, and logging.

Basic usage:

	logger := logging.NewLogger("info", false)
	strategy := deployment.NewDockerStrategy(config)
	service := rollback.NewService(rollback.DefaultConfig(), strategy, logger)

	service.RegisterVersion("v1.0.0")
	service.RegisterVersion("v1.1.0")

	err := service.Rollback("v1.1.0")
*/
package rollback
