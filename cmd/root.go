package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/xboshy/pure/internal/vault"
	"github.com/xboshy/pure/internal/vault/algorithms"
	"go.bryk.io/pkg/errors"
	xLog "go.bryk.io/pkg/log"
)

var algorithm = algorithms.Ed25519()

var rootCmd = &cobra.Command{
	Use:           "pure",
	Short:         "Pure Methods: Client Showcase",
	SilenceErrors: true,
	SilenceUsage:  true,
	Long: `Pure

For more information:
https://github.com/xboshy/pure`,
}

var (
	log      xLog.Logger
	cfgVault vault.Config
	homeDir  = ""
)

// Execute will process the CLI invocation.
func Execute() {
	// catch any panics
	defer func() {
		if err := errors.FromRecover(recover()); err != nil {
			log.Warning("recovered panic")
			fmt.Printf("%+v", err)
			os.Exit(1)
		}
	}()
	// execute command
	if err := rootCmd.Execute(); err != nil {
		if pe := new(errors.Error); errors.Is(err, pe) {
			log.WithField("error", err).Error("command failed")
		} else {
			log.Error(err.Error())
		}
		os.Exit(1)
	}
}

func init() {
	log = xLog.WithZero(xLog.ZeroOptions{PrettyPrint: true})
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	owDir, err := os.Getwd()
	if err != nil {
		log.Errorf("%s", err)
		os.Exit(1)
	}

	ex, err := os.Executable()
	if err != nil {
		log.Errorf("%s", err)
		os.Exit(1)
	}
	exPath := filepath.Dir(ex)
	homeDir = exPath

	// Set default values
	viper.SetDefault("address", "http://localhost:8200")
	viper.SetDefault("token", "please-update-me")
	viper.SetDefault("transit-path", "transit")

	// Set configuration file
	viper.AddConfigPath(owDir)
	viper.AddConfigPath(homeDir)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	// ENV
	viper.SetEnvPrefix("did")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// Read configuration file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Errorf("%s", err)

			viper.SafeWriteConfig()
			log.Infof("please update the generated configuration file: %s", "")

			os.Exit(1)
		}

		if viper.ConfigFileUsed() != "" {
			log.Errorf("failed to load configuration file: %s", viper.ConfigFileUsed())
			os.Exit(1)
		}
	}

	cfgVault.Address = viper.GetString("address")
	cfgVault.Token = viper.GetString("token")
	cfgVault.TransitPath = viper.GetString("transit-path")

	log.Infof("using configuration file: %s", viper.ConfigFileUsed())
}
