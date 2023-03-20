package cmd

import (
	"accountbalance/internal"
	"accountbalance/pkg/setting"
	"fmt"
	"os"
	"path"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var confStruct = &setting.Config{}

var rootCmd = &cobra.Command{
	Use:   "ab",
	Short: "etc",
	Long:  `etc`,
	Run: func(cmd *cobra.Command, args []string) {

		go internal.Balance(confStruct)
		internal.Send(confStruct)
		// time.Sleep(5 * time.Second)
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func initConfig() {
	var vp = viper.New()
	if cfgFile != "" {
		vp.SetConfigFile(cfgFile)
	} else {
		localPath, err := os.Getwd()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		vp.AddConfigPath(path.Join(localPath, "config"))
		vp.SetConfigName("config")
		vp.SetConfigType("toml")
	}
	if err := vp.ReadInConfig(); err != nil {
		panic(err)
	}
	conf := setting.NewSetting(vp)
	if err := conf.ReadConfig(confStruct); err != nil {
		panic(err)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is ./config/config.toml)")
}
