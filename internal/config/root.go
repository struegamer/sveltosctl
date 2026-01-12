package config

import (
	"context"

	"github.com/go-logr/logr"

	ictx "github.com/projectsveltos/sveltosctl/internal/ctx"
	"github.com/projectsveltos/sveltosctl/internal/logging"
)

type CtlConfig struct {
	CfgFile         string
	NoConsoleOutput bool
	Verbose         bool
	logger          *logging.CliLogger
}

func NewCtlConfig() *CtlConfig {
	return &CtlConfig{
		NoConsoleOutput: false,
		Verbose:         false,
		logger:          logging.NewLogger(),
	}
}

func NewCtlConfigWithContext(ctx context.Context) context.Context {
	ctlCfg := NewCtlConfig()
	ctx = context.WithValue(ctx, ictx.CtlConfigCtxKey, ctlCfg)
	return ctx
}

func CtlConfigFromContext(ctx context.Context) *CtlConfig {
	cfg, ok := ctx.Value(ictx.CtlConfigCtxKey).(*CtlConfig)
	if !ok || cfg == nil {
		panic("no config in context")
	}
	return cfg
}

func (cfg *CtlConfig) SetLogger(logger *logging.CliLogger) {
	cfg.logger = logger
}

func (cfg *CtlConfig) Logger() logr.Logger {
	return cfg.logger.Logger()
}
