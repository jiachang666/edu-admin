package authz

type Enforcer interface {
	Can(subject string, object string, action string) bool
}

type NoopEnforcer struct{}

func NewNoopEnforcer() *NoopEnforcer {
	return &NoopEnforcer{}
}

func (e *NoopEnforcer) Can(subject string, object string, action string) bool {
	return true
}
