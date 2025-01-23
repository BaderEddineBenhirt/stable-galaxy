package monitor

import (
	"fmt"
	"time"
)

type MonitorConfig struct {
	CPUThreshold     float64
	MemoryThreshold  float64
	ErrorThreshold   float64
	LatencyThreshold time.Duration
}

type Service struct {
	versions map[string]*Version
	config   MonitorConfig
}

func NewService(config MonitorConfig) *Service {
	return &Service{
		versions: make(map[string]*Version),
		config:   config,
	}
}

func (s *Service) AddVersion(number string) error {
	if _, exists := s.versions[number]; exists {
		return fmt.Errorf("version %s already exists", number)
	}

	s.versions[number] = &Version{
		Number:      number,
		Status:      StatusHealthy,
		Metrics:     &Metrics{},
		LastChecked: time.Now(),
		Errors:      make([]Error, 0),
	}
	return nil
}

func (s *Service) CheckHealth(version string) (Status, error) {
	v, exists := s.versions[version]
	if !exists {
		return "", fmt.Errorf("version %s not found", version)
	}

	metrics := v.Metrics
	if metrics.CPUUsage > s.config.CPUThreshold ||
		metrics.MemoryUsage > s.config.MemoryThreshold ||
		metrics.ErrorRate > s.config.ErrorThreshold ||
		metrics.Latency > s.config.LatencyThreshold {
		return StatusError, nil
	}

	return StatusHealthy, nil
}

func (s *Service) UpdateMetrics(version string, metrics *Metrics) error {
	v, exists := s.versions[version]
	if !exists {
		return fmt.Errorf("version %s not found", version)
	}

	v.Metrics = metrics
	v.LastChecked = time.Now()

	status, err := s.CheckHealth(version)
	if err != nil {
		return err
	}

	v.Status = status
	return nil
}
