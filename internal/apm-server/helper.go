package apm_server

import (
	"APM-server/internal/apm-server/store"
	"APM-server/internal/pkg/log"
	"APM-server/internal/pkg/model"
	"APM-server/pkg/db"
	"APM-server/pkg/kafka"
	"fmt"
	"github.com/spf13/viper"
	"os"
	"strings"
)

const (
	recommendedHomeDir = ".APM-server"

	defaultConfigName = "apm-server.yaml"
)

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		//home, err := os.UserHomeDir()
		//cobra.CheckErr(err)

		// Search config in home directory with name ".APM-server" (without extension).
		//viper.AddConfigPath(filepath.Join(home, recommendedHomeDir))
		viper.AddConfigPath(".")
		viper.SetConfigType("yaml")
		viper.SetConfigName(defaultConfigName)
	}

	viper.AutomaticEnv() // read in environment variables that match

	viper.SetEnvPrefix("APM-SERVER")

	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}

	//if err := godotenv.Load(".env"); err != nil {
	//	log.Fatalw("Error loading .env file")
	//}
}

// logOptions 从 viper 中读取日志配置，构建 `*log.Options` 并返回.
// 注意：`viper.Get<Type>()` 中 key 的名字需要使用 `.` 分割，以跟 YAML 中保持相同的缩进.
func logOptions() *log.Options {
	return &log.Options{
		DisableCaller:     viper.GetBool("log.disable-caller"),
		DisableStacktrace: viper.GetBool("log.disable-stacktrace"),
		Level:             viper.GetString("log.level"),
		Format:            viper.GetString("log.format"),
		OutputPaths:       viper.GetStringSlice("log.output-paths"),
	}
}

func initStore() error {
	dbOptions := &db.MySQLOptions{
		Host:                  viper.GetString("db.host"),
		Username:              viper.GetString("db.username"),
		Password:              viper.GetString("db.password"),
		Database:              viper.GetString("db.database"),
		MaxIdleConnections:    viper.GetInt("db.max-idle-connections"),
		MaxOpenConnections:    viper.GetInt("db.max-open-connections"),
		MaxConnectionLifeTime: viper.GetDuration("db.max-connection-life-time"),
		LogLevel:              viper.GetInt("db.log-level"),
	}

	ins, err := db.NewMySQL(dbOptions)
	if err != nil {
		return err
	}

	if err := ins.Set("gorm:table_options", "ENGINE=MyISAM").AutoMigrate(&model.UserM{}); err != nil {
		return err
	}

	_ = store.NewStore(ins)

	return nil
}

func initConsumerGroup() error {
	cgOptions := &kafka.KafkaOptions{
		ConsumerReturnErr: viper.GetBool("kafka.consumer-return-err"),
		GroupID:           viper.GetString("kafka.group-id"),
		Brokers:           viper.GetStringSlice("kafka.brokers"),
	}

	client, err := kafka.NewKafkaClient(cgOptions)
	if err != nil {
		return err
	}
	kafka.StoreClient(client)

	return nil
}
