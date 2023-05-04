package apm_server

import (
	"APM-server/internal/pkg/known"
	"APM-server/internal/pkg/log"
	"APM-server/internal/pkg/middleware"
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
		Use:   "APM-server",
		Short: "A brief description of your application",
		Long: `A longer description that spans multiple lines and likely contains
		examples and usage of using your application. For example:
		
		Cobra is a CLI library for Go that empowers applications.
		This application is a tool to generate the needed files
		to quickly create a Cobra application.`,
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
	cmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "The path to the miniblog configuration file. Empty string for no configuration file.")

	// Cobra 也支持本地标志，本地标志只能在其所绑定的命令上使用
	cmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	return cmd
}

func run() error {
	// 初始化 store 层
	if err := initStore(); err != nil {
		return err
	}

	// 设置 token 包的签发密钥，用于 token 包 token 的签发和解析
	token.Init(os.Getenv("jwtSecret"), known.XEmailKey)

	gin.SetMode(viper.GetString("runmode"))

	g := gin.New()

	mws := []gin.HandlerFunc{gin.Recovery(), middleware.NoCache, middleware.Cors, middleware.Secure, middleware.RequestID()}

	g.Use(mws...)

	if err := installRouters(g); err != nil {
		return err
	}

	httpsrv := &http.Server{Addr: viper.GetString("addr"), Handler: g}
	log.Infow("logger is running", "level:", viper.GetString("log.level"))
	log.Infow("start to listening the incoming requests on http address", "address:", viper.GetString("addr"))
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
