/*
Package monitor provides health checking and monitoring capabilities.

It includes functionality for checking application health, collecting metrics,
and monitoring deployment status across different versions.

Basic usage:

	monitor := monitor.New(config)
	health := monitor.CheckHealth("v1.0.0")
*/
package monitor
