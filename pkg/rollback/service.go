package rollback

import (
	"sort"
	"time"

	"github.com/BaderEddineBenhirt/stable-galaxy/pkg/deployment"
	"github.com/BaderEddineBenhirt/stable-galaxy/pkg/errors"
	"github.com/BaderEddineBenhirt/stable-galaxy/pkg/logging"
)

type Service struct {
	config   RollbackConfig
	versions []string
	strategy deployment.Strategy
	logger   *logging.Logger
}

func NewService(config RollbackConfig, strategy deployment.Strategy, logger *logging.Logger) *Service {
	if logger == nil {
		logger = logging.NewLogger("info", false)
	}
	return &Service{
		config:   config,
		versions: make([]string, 0),
		strategy: strategy,
		logger:   logger,
	}
}

func (s *Service) RegisterVersion(version string) {
	s.logger.Debug().Str("version", version).Msg("Registering version")
	s.versions = append(s.versions, version)
	sort.Strings(s.versions)
}

func (s *Service) findPreviousStableVersion(currentVersion string) (string, error) {
	s.logger.Debug().Str("current_version", currentVersion).Msg("Finding previous stable version")
	for i := len(s.versions) - 1; i >= 0; i-- {
		if s.versions[i] < currentVersion {
			if s.config.ValidateVersion != nil && !s.config.ValidateVersion(s.versions[i]) {
				s.logger.Debug().Str("version", s.versions[i]).Msg("Version validation failed")
				continue
			}
			s.logger.Debug().Str("found_version", s.versions[i]).Msg("Found stable version")
			return s.versions[i], nil
		}
	}
	return "", errors.NewValidationError("no stable previous version found", nil)
}

func (s *Service) executeRollback(from, to string) error {
	s.logger.Debug().Str("from", from).Str("to", to).Msg("Executing rollback")
	meta := map[string]interface{}{
		"from_version": from,
		"to_version":   to,
		"strategy":     s.strategy.StrategyName(),
	}

	if err := s.strategy.Rollback(from, to); err != nil {
		return errors.NewDeploymentError("rollback execution failed", err, meta)
	}
	return nil
}

func (s *Service) Rollback(currentVersion string) error {
	s.logger.Info().Str("from_version", currentVersion).Msg("Starting rollback")

	targetVersion, err := s.findPreviousStableVersion(currentVersion)
	if err != nil {
		s.logger.Error().Err(err).Str("current_version", currentVersion).Msg("Failed to find stable version")
		return errors.NewValidationError("failed to find stable version", err)
	}

	if s.config.PreRollbackHook != nil {
		s.logger.Debug().Msg("Executing pre-rollback hook")
		if err := s.config.PreRollbackHook(); err != nil {
			s.logger.Error().Err(err).Msg("Pre-rollback hook failed")
			return errors.NewDeploymentError("pre-rollback hook failed", err, nil)
		}
	}

	for attempt := 1; attempt <= s.config.MaxAttempts; attempt++ {
		s.logger.Info().Int("attempt", attempt).Int("max_attempts", s.config.MaxAttempts).Msg("Attempting rollback")

		if err := s.executeRollback(currentVersion, targetVersion); err != nil {
			if attempt == s.config.MaxAttempts {
				s.logger.Error().Err(err).Int("attempts", attempt).Msg("Rollback failed after all attempts")
				if s.config.OnFailureHook != nil {
					s.config.OnFailureHook(err)
				}
				return errors.NewDeploymentError("rollback failed after all attempts", err, nil)
			}
			s.logger.Warn().Err(err).Int("attempt", attempt).Dur("backoff", s.config.BackoffDuration).Msg("Retrying after backoff")
			time.Sleep(s.config.BackoffDuration)
			continue
		}
		break
	}

	if s.config.PostRollbackHook != nil {
		s.logger.Debug().Msg("Executing post-rollback hook")
		if err := s.config.PostRollbackHook(); err != nil {
			s.logger.Error().Err(err).Msg("Post-rollback hook failed")
			return errors.NewDeploymentError("post-rollback hook failed", err, nil)
		}
	}

	s.logger.Info().Str("from", currentVersion).Str("to", targetVersion).Msg("Rollback completed successfully")
	return nil
}
