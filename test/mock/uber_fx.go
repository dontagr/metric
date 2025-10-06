package mocks

import (
	"context"

	"go.uber.org/fx"
)

type MockLifecycle struct {
	hooks []fx.Hook
}

func NewMockLifecycle() *MockLifecycle {
	return &MockLifecycle{
		hooks: []fx.Hook{},
	}
}

func (m *MockLifecycle) Append(hook fx.Hook) {
	m.hooks = append(m.hooks, hook)
}

func (m *MockLifecycle) Start(ctx context.Context) error {
	for _, hook := range m.hooks {
		if hook.OnStart != nil {
			if err := hook.OnStart(ctx); err != nil {
				return err
			}
		}
	}
	return nil
}

func (m *MockLifecycle) Stop(ctx context.Context) error {
	for i := len(m.hooks) - 1; i >= 0; i-- {
		if m.hooks[i].OnStop != nil {
			if err := m.hooks[i].OnStop(ctx); err != nil {
				return err
			}
		}
	}
	return nil
}
