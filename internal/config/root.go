package config

import (
	"context"

	"github.com/projectsveltos/sveltosctl/internal/k8s"

	ictx "github.com/projectsveltos/sveltosctl/internal/ctx"
	"github.com/projectsveltos/sveltosctl/internal/logging"
)

const ctlConfigCtxKey = ictx.ContextKey("ctlConfig")
const ProjectSveltos = "projectsveltos"

type CtlConfig struct {
	CfgFile         string
	NoConsoleOutput bool
	Verbose         bool
	logger          logging.Logger
	mgmtCluster     *k8s.Cluster
}

func NewCtlConfig() *CtlConfig {
	return &CtlConfig{
		NoConsoleOutput: false,
		Verbose:         false,
		logger:          logging.NewKlogLogger(nil),
		mgmtCluster:     nil,
	}
}

func NewCtlConfigWithContext(ctx context.Context) context.Context {
	ctlCfg := NewCtlConfig()
	ctx = context.WithValue(ctx, ctlConfigCtxKey, ctlCfg)
	return ctx
}

func CtlConfigFromContext(ctx context.Context) *CtlConfig {
	cfg, ok := ctx.Value(ctlConfigCtxKey).(*CtlConfig)
	if !ok || cfg == nil {
		panic("no config in context")
	}
	return cfg
}

func (cfg *CtlConfig) SetLogger(logger logging.Logger) {
	cfg.logger = logger
}

func (cfg *CtlConfig) Logger() logging.Logger {
	return cfg.logger
}

func (cfg *CtlConfig) SetMgmtCluster(mgmtCluster *k8s.Cluster) {
	cfg.mgmtCluster = mgmtCluster
}
func (cfg *CtlConfig) MgmtCluster() *k8s.Cluster {
	return cfg.mgmtCluster
}
