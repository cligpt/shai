package cmd

import (
	"context"
	"os"

	"github.com/alecthomas/kingpin/v2"
	"github.com/hashicorp/go-hclog"
	"github.com/pkg/errors"

	"github.com/cligpt/shai/config"
	"github.com/cligpt/shai/drive"
	"github.com/cligpt/shai/gpt"
	"github.com/cligpt/shai/term"
)

const (
	aiName     = "shai"
	routineNum = -1
)

var (
	app      = kingpin.New(aiName, "shell with ai").Version(config.Version + "-build-" + config.Build)
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

	d, err := initDrive(ctx, logger, c)
	if err != nil {
		return errors.Wrap(err, "failed to init drive")
	}

	g, err := initGpt(ctx, logger, c)
	if err != nil {
		return errors.Wrap(err, "failed to init gpt")
	}

	t, err := initTerm(ctx, logger, c, d, g)
	if err != nil {
		return errors.Wrap(err, "failed to init term")
	}

	if err := runTerm(ctx, logger, t); err != nil {
		return errors.Wrap(err, "failed to run term")
	}

	return nil
}

func initLogger(_ context.Context, level string) (hclog.Logger, error) {
	return hclog.New(&hclog.LoggerOptions{
		Name:  aiName,
		Level: hclog.LevelFromString(level),
	}), nil
}

func initConfig(_ context.Context, _ hclog.Logger) (*config.Config, error) {
	c := config.New()
	return c, nil
}

func initDrive(ctx context.Context, logger hclog.Logger, cfg *config.Config) (drive.Drive, error) {
	c := drive.DefaultConfig()
	if c == nil {
		return nil, errors.New("failed to config")
	}

	c.Logger = logger
	c.Config = *cfg

	return drive.New(ctx, c), nil
}

func initGpt(ctx context.Context, logger hclog.Logger, cfg *config.Config) (gpt.Gpt, error) {
	c := gpt.DefaultConfig()
	if c == nil {
		return nil, errors.New("failed to config")
	}

	c.Logger = logger
	c.Config = *cfg

	return gpt.New(ctx, c), nil
}

func initTerm(ctx context.Context, logger hclog.Logger, cfg *config.Config, _drive drive.Drive, _gpt gpt.Gpt) (term.Term, error) {
	c := term.DefaultConfig()
	if c == nil {
		return nil, errors.New("failed to config")
	}

	c.Logger = logger
	c.Config = *cfg

	return term.New(ctx, c), nil
}

func runTerm(ctx context.Context, _ hclog.Logger, _term term.Term) error {
	if err := _term.Init(ctx); err != nil {
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
