package config

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta2"

	configv1beta1 "github.com/projectsveltos/addon-controller/api/v1beta1"
	eventv1beta1 "github.com/projectsveltos/event-manager/api/v1beta1"
	libsveltosv1beta1 "github.com/projectsveltos/libsveltos/api/v1beta1"
	utilsv1beta1 "github.com/projectsveltos/sveltosctl/api/v1beta1"
	"github.com/projectsveltos/sveltosctl/internal/k8s"
)

const envPrefix = "SVELTOSCTL"
const cfgPath = ".sveltosctl"

func (cfg *CtlConfig) Initialize(cmd *cobra.Command) error {
	if cfg.Verbose {
		cfg.logger.Logger().Info("Initializing Viper config")
	}
	err := cfg.initViper()
	if err != nil {
		return err
	}
	return nil
}

func (cfg *CtlConfig) initViper() error {
	viper.SetEnvPrefix(envPrefix)
	// Allow for nested keys in environment variables (e.g. `MYAPP_DATABASE_HOST`)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "*", "-", "*"))
	viper.AutomaticEnv()
	if cfg.CfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfg.CfgFile)
	} else {
		// Search for a config file in default locations and panic if it can not be found
		home := GetUserHomeDir()

		// Search for a config file with the name "config" (without extension).
		viper.AddConfigPath(".")
		viper.AddConfigPath(home + "/" + cfgPath)
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
	}
	return nil
}

func (cfg *CtlConfig) createManagementClusterAccess() {
	var apiSchemasList []k8s.ApiSchemaFunc = []k8s.ApiSchemaFunc{
		corev1.AddToScheme,
		appsv1.AddToScheme,
		configv1beta1.AddToScheme,
		utilsv1beta1.AddToScheme,
		clusterv1.AddToScheme,
		libsveltosv1beta1.AddToScheme,
		eventv1beta1.AddToScheme,
		rbacv1.AddToScheme,
		apiextensionsv1.AddToScheme,
	}
	cfg.mgmtCluster = k8s.NewCluster(apiSchemasList...)
}
