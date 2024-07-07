package term

import (
	"context"

	"github.com/hashicorp/go-hclog"

	"github.com/cligpt/shai/config"
)

type Term interface {
	Init(context.Context) error
	Deinit(context.Context) error
	Run(context.Context) error
}

type Config struct {
	Config config.Config
	Logger hclog.Logger
}

type term struct {
	cfg *Config
}

func New(_ context.Context, cfg *Config) Term {
	return &term{
		cfg: cfg,
	}
}

func DefaultConfig() *Config {
	return &Config{}
}

func (t *term) Init(_ context.Context) error {
	return nil
}

func (t *term) Deinit(_ context.Context) error {
	return nil
}

func (t *term) Run(_ context.Context) error {
	return nil
}
