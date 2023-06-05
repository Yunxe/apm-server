package apm_server

import (
	"APM-server/internal/pkg/known"
	"APM-server/internal/pkg/log"
	"APM-server/internal/pkg/middleware"
	"APM-server/pkg/kafka"
	"APM-server/pkg/token"
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/spf13/cobra"
)

var cfgFile string
var Logger *zap.SugaredLogger

func NewApmServerCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "apm-server",
		Short: "A little monitor server for Specialize Design",
		//SilenceUsage: true,

		RunE: func(cmd *cobra.Command, args []string) error {
			log.Init(logOptions())

			defer log.Sync() // Sync 将缓存中的日志刷新到磁盘文件中
			return run()
		},
		Args: func(cmd *cobra.Command, args []string) error {
			for _, arg := range args {
				if len(arg) > 0 {
					return fmt.Errorf("%q does not take any arguments, got %q", cmd.CommandPath(), args)
				}
			}

			return nil
		},
	}
	// 以下设置，使得 initConfig 函数在每个命令运行时都会被调用以读取配置
	cobra.OnInitialize(initConfig)

	// 在这里您将定义标志和配置设置。

	// Cobra 支持持久性标志(PersistentFlag)，该标志可用于它所分配的命令以及该命令下的每个子命令
	cmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "The path to the apm-server configuration file. Empty string for no configuration file.")

	// Cobra 也支持本地标志，本地标志只能在其所绑定的命令上使用
	cmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	var initTopicFlag bool
	initTopicCmd := &cobra.Command{
		Use:   "init-topic",
		Short: "init topics and partitions",
		Run: func(cmd *cobra.Command, args []string) {
			if initTopicFlag{
				err:=initConsumerGroup()
				if err != nil {
					log.Errorw("init cg err","err",err)
				}
				kafka.InitTopics()
				kafka.KS().Client().Close()
			}
		},
	}
	initTopicCmd.Flags().BoolVarP(&initTopicFlag,"topic","t",false,"true to init topics and partitions")
	cmd.AddCommand(initTopicCmd)

	return cmd
}

func run() error {
	//var wg sync.WaitGroup

	// 初始化 store 层
	if err := initStore(); err != nil {
		return err
	}


	//初始化kafka消费组
	err:=initConsumerGroup();
	if err != nil {
		log.Errorw("init cg err","err",err)
	}
	defer kafka.KS().Client().Close()

	// 设置 token 包的签发密钥，用于 token 包 token 的签发和解析
	token.Init(viper.GetString("jwtSecret"), known.XEmailKey)

	gin.SetMode(viper.GetString("runmode"))

	g := gin.New()

	mws := []gin.HandlerFunc{gin.Recovery(), middleware.NoCache, middleware.Cors, middleware.Secure, middleware.RequestID()}

	g.Use(mws...)

	if err := installRouters(g); err != nil {
		return err
	}

	httpsrv := &http.Server{Addr: viper.GetString("port"), Handler: g}
	log.Infow("logger is running", "level:", viper.GetString("log.level"))
	log.Infow("start to listening the incoming requests on http address", "address:", viper.GetString("port"))
	go func() {
		if err := httpsrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalw(err.Error())
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Infow("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := httpsrv.Shutdown(ctx); err != nil {
		log.Errorw("Insecure Server forced to shutdown", "err", err)
		return err
	}
	log.Infow("Server exiting")

	return nil
}
