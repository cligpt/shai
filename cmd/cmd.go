package cmd

import (
	"context"
	"github.com/alecthomas/kingpin/v2"
	"github.com/hashicorp/go-hclog"
	"github.com/pkg/errors"
	"os"

	"github.com/cligpt/shai/config"
	"github.com/cligpt/shai/gpt"
	"github.com/cligpt/shai/term"
)

const (
	logName    = "shai"
	routineNum = -1
)

var (
	app      = kingpin.New("shai", "shell with ai").Version(config.Version + "-build-" + config.Build)
	logLevel = app.Flag("log-level", "Log level (DEBUG|INFO|WARN|ERROR)").Default("WARN").String()
)

func Run(ctx context.Context) error {
	kingpin.MustParse(app.Parse(os.Args[1:]))

	logger, err := initLogger(ctx, *logLevel)
	if err != nil {
		return errors.Wrap(err, "failed to init logger")
	}

	c, err := initConfig(ctx, logger)
	if err != nil {
		return errors.Wrap(err, "failed to init config")
	}

	t, err := initTerm(ctx, logger, c)
	if err != nil {
		return errors.Wrap(err, "failed to init term")
	}

	g, err := initGpt(ctx, logger, c)
	if err != nil {
		return errors.Wrap(err, "failed to init gpt")
	}

	if err := runTerm(ctx, logger, t, g); err != nil {
		return errors.Wrap(err, "failed to run term")
	}

	return nil
}

func initLogger(_ context.Context, level string) (hclog.Logger, error) {
	return hclog.New(&hclog.LoggerOptions{
		Name:  logName,
		Level: hclog.LevelFromString(level),
	}), nil
}

func initConfig(_ context.Context, _ hclog.Logger) (*config.Config, error) {
	c := config.New()
	return c, nil
}

func initTerm(ctx context.Context, logger hclog.Logger, _ *config.Config) (term.Term, error) {
	c := term.DefaultConfig()
	if c == nil {
		return nil, errors.New("failed to config")
	}

	c.Logger = logger

	return term.New(ctx, c), nil
}

func initGpt(ctx context.Context, logger hclog.Logger, _ *config.Config) (gpt.Gpt, error) {
	c := gpt.DefaultConfig()
	if c == nil {
		return nil, errors.New("failed to config")
	}

	c.Logger = logger

	return gpt.New(ctx, c), nil
}

func runTerm(ctx context.Context, _ hclog.Logger, _term term.Term, _gpt gpt.Gpt) error {
	if err := _term.Init(ctx, _gpt); err != nil {
		return errors.New("failed to init")
	}

	defer func(_term term.Term, ctx context.Context) {
		_ = _term.Deinit(ctx)
	}(_term, ctx)

	if err := _term.Run(ctx); err != nil {
		return errors.Wrap(err, "failed to run")
	}

	return nil
}
