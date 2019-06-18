package cmd

import (
	"fmt"
	"syscall"

	"github.com/fvbock/endless"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/yitume/caller"
	"github.com/yitume/caller/ginsession"
	"github.com/yitume/caller/gorm"
	"github.com/yitume/caller/redigo"
	callerzap "github.com/yitume/caller/zap"
	"go.uber.org/zap"

	"git.yitum.com/saas/shop-admin/model"
	"git.yitum.com/saas/shop-admin/pkg/bootstrap"
	"git.yitum.com/saas/shop-admin/router"
	"git.yitum.com/saas/shop-admin/service"
)

// startCmd represents the hello command
var startCmd = &cobra.Command{
	Use:  "start",
	Long: `Starts shop-admin server`,
	Run:  startFn,
}

func init() {
	startCmd.PersistentFlags().StringVar(&bootstrap.Arg.CfgFile, "config", "conf/conf.toml", "config file (default is $HOME/.cobra-example.yaml)")
	RootCmd.AddCommand(startCmd)
	cobra.OnInitialize(initConfig)
}

func startFn(cmd *cobra.Command, args []string) {
	if err := caller.Init(
		bootstrap.Arg.CfgFile,
		callerzap.New,
		gorm.New,
		ginsession.New,
		redigo.New,
	); err != nil {
		panic(err)
	}

	model.Init()

	// 配置初始化
	if err := bootstrap.InitConfig(bootstrap.Arg.CfgFile); err != nil {
		model.Logger.Panic(err.Error())
	}
	service.Init()
	service.InitGen()
	handler := router.InitRouter()

	endless.DefaultReadTimeOut = bootstrap.Conf.Server.ReadTimeout.Duration
	endless.DefaultWriteTimeOut = bootstrap.Conf.Server.WriteTimeout.Duration
	endless.DefaultMaxHeaderBytes = bootstrap.Conf.Server.MaxHeaderBytes
	server := endless.NewServer(bootstrap.Conf.Server.Addr, handler)
	server.BeforeBegin = func(add string) {
		model.Logger.Info(fmt.Sprintf("Actual pid is %d", syscall.Getpid()))
	}

	err := server.ListenAndServe()
	if err != nil {
		model.Logger.Error("Server err", zap.String("err", err.Error()))
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.SetConfigFile(bootstrap.Arg.CfgFile)
	viper.AutomaticEnv() // read in environment variables that match
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
