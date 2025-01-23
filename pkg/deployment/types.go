package deployment

type Strategy interface {
	Rollback(from, to string) error
	Deploy(version string) error
	GetCurrentVersion() (string, error)
	StrategyName() string
}
