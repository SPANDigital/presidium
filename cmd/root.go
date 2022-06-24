package cmd

import (
	"fmt"
	"github.com/SPANDigital/presidium-hugo/pkg/config"
	"github.com/SPANDigital/presidium-hugo/pkg/log"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"strings"

	"github.com/spf13/viper"
)

var (
	cfgFile string
	debug   bool
	rootCmd = &cobra.Command{
		Use:   "presidium",
		Short: "A brief description of your application",
		Long:  "CLI tools for managing Presidium Hugo content",
	}
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.FatalWithFields("error executing command", log.Fields{"error": err})
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().BoolVar(&debug, config.DebugKey, false, "enables debug logs")
	rootCmd.PersistentFlags().StringVar(&cfgFile, config.ConfigFileKey, "", "config file (default is $HOME/.presidium.yaml)")
	viper.BindPFlag(config.DebugKey, rootCmd.PersistentFlags().Lookup(config.DebugKey))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		viper.SetEnvPrefix(config.ApplicationName)
		viper.AddConfigPath(".")
		viper.SetConfigName(fmt.Sprintf(".%s", config.ApplicationName))
	}

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_", ".", "_"))

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.InfoWithFields("using config file", log.Fields{
			"file": viper.ConfigFileUsed(),
		})
	}

	if viper.GetBool(config.DebugKey) {
		log.SetLogLevel(logrus.DebugLevel)
	} else {
		logLevel := viper.GetString(config.LoggingLevelKey)
		if len(logLevel) == 0 {
			log.SetLogLevel(logrus.InfoLevel)
		} else {
			level, err := logrus.ParseLevel(logLevel)
			if err != nil {
				log.SetLogLevel(logrus.InfoLevel)
				log.WarnWithFields("unable to parse logging level. Defaulting to info", log.Fields{"level": logLevel})
			} else {
				log.SetLogLevel(level)
			}
		}
	}
}
