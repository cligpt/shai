package cmd

import (
	"context"

	"github.com/hashicorp/go-hclog"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cligpt/shai/config"
	"github.com/cligpt/shai/drive"
	"github.com/cligpt/shai/gpt"
	"github.com/cligpt/shai/term"
)

const (
	aiName = "shai"
)

var (
	configFile string
	logLevel   string
)

var rootCmd = &cobra.Command{
	Use:     aiName,
	Version: config.Version + "-build-" + config.Build,
	Short:   "shell with ai",
	Long:    "shell with ai",
	Run: func(cmd *cobra.Command, args []string) {
		cobra.CheckErr(loadConfig(context.Background()))
	},
}

// nolint: gochecknoinits
func init() {
	helper := func() {
		if configFile != "" {
			viper.SetConfigFile(configFile)
		} else {
			home, _ := homedir.Dir()
			viper.AddConfigPath(home)
			viper.AddConfigPath(".shai")
			viper.SetConfigName(aiName)
			viper.SetConfigType("yml")
		}
	}

	cobra.OnInitialize(helper)

	rootCmd.Flags().StringVarP(&configFile, "config-file", "f", "$HOME/.shai/shai.yml", "config file")
	rootCmd.Flags().StringVarP(&logLevel, "log-level", "l", "WRAN", "log level (DEBUG|INFO|WARN|ERROR)")
}

func Execute() error {
	return rootCmd.Execute()
}

func loadConfig(ctx context.Context) error {
	logger, err := initLogger(ctx, logLevel)
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
	c.Drive = _drive
	c.Gpt = _gpt

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
