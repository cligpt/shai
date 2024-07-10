package drive

import (
	"context"

	"github.com/hashicorp/go-hclog"
)

type Drive interface {
	Init(context.Context) error
	Deinit(context.Context) error
	Run(context.Context) error
}

type Config struct {
	Logger hclog.Logger
}

type drive struct {
	cfg *Config
}

func New(_ context.Context, cfg *Config) Drive {
	return &drive{
		cfg: cfg,
	}
}

func DefaultConfig() *Config {
	return &Config{}
}

func (d *drive) Init(_ context.Context) error {
	return nil
}

func (d *drive) Deinit(_ context.Context) error {
	return nil
}

func (d *drive) Run(_ context.Context) error {
	return nil
}
