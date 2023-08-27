package configs

import (
	"bytes"

	"github.com/c0dered273/go-musthave-diploma-tpl/internal/validators"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	configFileType = "yaml"

	envVars = []string{
		"RUN_ADDRESS",
		"DATABASE_URI",
		"ACCRUAL_SYSTEM_ADDRESS",
		"API_SECRET",
	}
)

type ServerConfig struct {
	RunAddress           string   `mapstructure:"run_address" validate:"required"`
	DatabaseURI          string   `mapstructure:"database_uri" validate:"required"`
	AccrualSystemAddress string   `mapstructure:"accrual_system_address" validate:"required"`
	APISecret            string   `mapstructure:"api_secret" validate:"required"`
	Database             Database `mapstructure:"database"`
	Server               Server   `mapstructure:"server"`
}

type Database struct {
	LoggerLevel string     `mapstructure:"logger_level"`
	Connection  Connection `mapstructure:"connection"`
}

type Connection struct {
	Options map[string]string `mapstructure:"options"`
}

type Server struct {
	Name        string `mapstructure:"name"`
	Logger      Logger `mapstructure:"logger"`
	DebugConfig bool   `mapstructure:"debug_config"`
	PprofEnable bool   `mapstructure:"pprof_enable"`
}

type Logger struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
	Caller bool   `mapstructure:"caller"`
}

func bindPFlags() error {
	pflag.StringP("run_address", "a", viper.GetString("run_address"), "Server address:port")
	pflag.StringP("database_uri", "d", viper.GetString("database_uri"), "Database connection uri")
	pflag.StringP("accrual_system_address", "r", viper.GetString("accrual_system_address"), "Accrual server address:port")
	pflag.Parse()
	err := viper.BindPFlags(pflag.CommandLine)
	if err != nil {
		return err
	}
	return nil
}

func setDefaults() {
	viper.SetDefault("api_secret", uuid.NewString())
	viper.SetDefault("database.logger_level", "info")
	viper.SetDefault("server.logger.level", "info")
}

func NewServerConfig(filename string, configPath []string, logger zerolog.Logger, validator *validator.Validate) (*ServerConfig, error) {
	setDefaults()

	err := bindConfigFile(filename, configPath, logger)
	if err != nil {
		return nil, err
	}

	err = bindEnvVars()
	if err != nil {
		return nil, err
	}

	err = bindPFlags()
	if err != nil {
		return nil, err
	}

	cfg, nErr := newConfig()
	if nErr != nil {
		return nil, nErr
	}

	err = validators.ValidateStructWithLogger(cfg, logger, validator)
	if err != nil {
		return nil, err
	}

	logDebugInfo(cfg, logger)

	return cfg, nil
}

func bindConfigFile(filename string, configPath []string, logger zerolog.Logger) error {
	viper.SetConfigName(filename)
	viper.SetConfigType(configFileType)
	for _, path := range configPath {
		viper.AddConfigPath(path)
	}
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			logger.Error().Msg("config: file not found")
		} else {
			return err
		}
	}
	return nil
}

func bindEnvVars() error {
	for _, env := range envVars {
		err := viper.BindEnv(env)
		if err != nil {
			return err
		}
	}
	return nil
}

func newConfig() (*ServerConfig, error) {
	cfg := &ServerConfig{}
	err := viper.Unmarshal(cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

func logDebugInfo(cfg *ServerConfig, logger zerolog.Logger) {
	var debug bytes.Buffer
	viper.DebugTo(&debug)
	if cfg.Server.DebugConfig {
		logger.Debug().Msgf("config: debug info \n%s", debug.String())
	}
}
