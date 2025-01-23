/*
Package logging provides structured logging functionality for the stable-galaxy library.

It wraps zerolog to provide consistent logging across all components with support
for both JSON and console formats.

Basic usage:

	logger := logging.NewLogger("info", false)
	logger.Info().Str("version", "v1.0.0").Msg("Starting rollback")
*/
package logging
