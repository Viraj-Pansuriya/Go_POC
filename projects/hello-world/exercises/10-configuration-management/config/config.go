package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// Config holds all application configuration
// Similar to Spring's @ConfigurationProperties
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	App      AppConfig      `mapstructure:"app"`
}

type ServerConfig struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	ReadTimeout  int    `mapstructure:"read_timeout"`
	WriteTimeout int    `mapstructure:"write_timeout"`
}

type DatabaseConfig struct {
	Driver   string `mapstructure:"driver"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	SSLMode  string `mapstructure:"sslmode"`
}

type JWTConfig struct {
	Secret     string `mapstructure:"secret"`
	Expiration int    `mapstructure:"expiration"` // hours
}

type AppConfig struct {
	Name    string `mapstructure:"name"`
	Version string `mapstructure:"version"`
	Debug   bool   `mapstructure:"debug"`
}

// setDefaults sets default values for configuration
// Like Spring's @Value("${property:defaultValue}")
func setDefaults() {
	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.read_timeout", 30)
	viper.SetDefault("server.write_timeout", 30)

	viper.SetDefault("database.driver", "postgres")
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 5432)
	viper.SetDefault("database.sslmode", "disable")

	viper.SetDefault("jwt.expiration", 24)
	viper.SetDefault("jwt.secret", []byte("verebhfevegreethergewergerwgewrwnhtgerfdsv"))

	viper.SetDefault("app.name", "Go App")
	viper.SetDefault("app.version", "1.0.0")
	viper.SetDefault("app.debug", false)
}

// LoadConfig loads configuration from file and environment
// profile: "dev", "prod", etc. (like Spring profiles)
func LoadConfig(configPath string, profile string) (*Config, error) {
	// Set defaults first (lowest priority)
	setDefaults()

	// Configure viper to read config files
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configPath)

	// Read base config
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read base config: %w", err)
	}

	// Merge profile-specific config if provided
	// Like Spring's application-{profile}.yml
	if profile != "" {
		viper.SetConfigName("config-" + profile)
		if err := viper.MergeInConfig(); err != nil {
			// Profile config is optional, don't fail if not found
			fmt.Printf("Note: No config-%s.yaml found, using base config\n", profile)
		}
	}

	// Enable environment variable override
	// Like Spring's SPRING_DATASOURCE_URL -> spring.datasource.url
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Unmarshal into struct (like @ConfigurationProperties)
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil
}

// GetDSN returns database connection string
func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode,
	)
}
