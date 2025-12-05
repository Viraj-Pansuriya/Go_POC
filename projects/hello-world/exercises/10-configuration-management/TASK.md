# Exercise 10: Configuration Management with Viper

## 🎯 Goal
Learn how to manage application configuration in Go using Viper - the Go equivalent of Spring Boot's `@ConfigurationProperties` and `application.yml`!

---

## 📚 Spring Boot vs Go (Viper)

### Spring Boot
```java
// application.yml
server:
  port: 8080
database:
  host: localhost
  port: 5432
  name: myapp

// Properties class with @ConfigurationProperties
@ConfigurationProperties(prefix = "database")
@Component
public class DatabaseConfig {
    private String host;
    private int port;
    private String name;
    // getters, setters...
}

// Usage
@Autowired
private DatabaseConfig dbConfig;
```

### Go with Viper
```go
// config.yaml
server:
  port: 8080
database:
  host: localhost
  port: 5432
  name: myapp

// Config struct (like @ConfigurationProperties)
type Config struct {
    Server   ServerConfig   `mapstructure:"server"`
    Database DatabaseConfig `mapstructure:"database"`
}

type DatabaseConfig struct {
    Host string `mapstructure:"host"`
    Port int    `mapstructure:"port"`
    Name string `mapstructure:"name"`
}

// Load config
func LoadConfig() (*Config, error) {
    viper.SetConfigName("config")
    viper.SetConfigType("yaml")
    viper.AddConfigPath(".")
    
    if err := viper.ReadInConfig(); err != nil {
        return nil, err
    }
    
    var config Config
    if err := viper.Unmarshal(&config); err != nil {
        return nil, err
    }
    return &config, nil
}
```

---

## 🏗️ What is Viper?

Viper is a complete configuration solution for Go applications. It supports:

- **Setting defaults**
- **Reading from YAML, JSON, TOML, HCL, envfile, Java properties**
- **Reading from environment variables** (like Spring's `SPRING_DATASOURCE_URL`)
- **Reading from command line flags**
- **Live watching and hot-reloading config** (like Spring Cloud Config)
- **Reading from remote config systems** (etcd, Consul)

```
Priority (highest to lowest):
1. explicit call to Set()
2. flags
3. environment variables
4. config file
5. key/value store
6. default
```

---

## 🚀 Your Tasks

### Task 1: Basic Config Loading

**Install Viper:**
```bash
go get github.com/spf13/viper
```

Create a simple config loader:

```go
// config/config.go
package config

import "github.com/spf13/viper"

type Config struct {
    Server   ServerConfig   `mapstructure:"server"`
    Database DatabaseConfig `mapstructure:"database"`
    JWT      JWTConfig      `mapstructure:"jwt"`
}

type ServerConfig struct {
    Port         int    `mapstructure:"port"`
    Host         string `mapstructure:"host"`
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

func LoadConfig(path string) (*Config, error) {
    // TODO: Set config file path
    // TODO: Read config file
    // TODO: Unmarshal into Config struct
    // TODO: Return config
}
```

### Task 2: Environment Variable Override

Like Spring's `SPRING_DATASOURCE_URL`, Viper can read from env vars:

```go
// config/config.go
func LoadConfig(path string) (*Config, error) {
    viper.SetConfigFile(path)
    
    // Enable env variable override
    viper.AutomaticEnv()
    
    // Replace . with _ for nested keys
    // database.host -> DATABASE_HOST
    viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
    
    // TODO: Allow overriding any config with env vars
    // Example: DATABASE_HOST=prod-db.aws.com overrides database.host
}
```

### Task 3: Profile-Based Configuration

Like Spring profiles (`application-dev.yml`, `application-prod.yml`):

```go
// config/config.go
func LoadConfigWithProfile(profile string) (*Config, error) {
    // Load base config first
    viper.SetConfigName("config")
    viper.SetConfigType("yaml")
    viper.AddConfigPath("./configs")
    viper.ReadInConfig()
    
    // Then merge profile-specific config
    // configs/config-dev.yaml or configs/config-prod.yaml
    if profile != "" {
        viper.SetConfigName("config-" + profile)
        viper.MergeInConfig()  // Merge, don't replace!
    }
    
    // TODO: Implement profile-based loading
}
```

### Task 4: Set Defaults

Like Spring's default values:

```go
func setDefaults() {
    viper.SetDefault("server.port", 8080)
    viper.SetDefault("server.host", "localhost")
    viper.SetDefault("server.read_timeout", 30)
    viper.SetDefault("server.write_timeout", 30)
    viper.SetDefault("database.driver", "postgres")
    viper.SetDefault("database.sslmode", "disable")
    viper.SetDefault("jwt.expiration", 24)
}
```

### Task 5: Config Validation

Validate config at startup (fail fast like Spring):

```go
// config/validator.go
func (c *Config) Validate() error {
    if c.Server.Port <= 0 || c.Server.Port > 65535 {
        return fmt.Errorf("invalid server port: %d", c.Server.Port)
    }
    if c.Database.Host == "" {
        return fmt.Errorf("database host is required")
    }
    if c.JWT.Secret == "" {
        return fmt.Errorf("JWT secret is required")
    }
    if len(c.JWT.Secret) < 32 {
        return fmt.Errorf("JWT secret must be at least 32 characters")
    }
    return nil
}
```

### Task 6: Wire Into Application

```go
// main.go
func main() {
    // Get profile from env or flag
    profile := os.Getenv("APP_PROFILE")
    if profile == "" {
        profile = "dev" // default
    }
    
    // Load config
    cfg, err := config.LoadConfigWithProfile(profile)
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }
    
    // Validate
    if err := cfg.Validate(); err != nil {
        log.Fatalf("Invalid config: %v", err)
    }
    
    // Use config
    fmt.Printf("Starting server on %s:%d\n", cfg.Server.Host, cfg.Server.Port)
    
    // Pass config to services
    r := gin.Default()
    
    // Example: Pass to handler
    userHandler := handler.NewUserHandler(cfg)
    r.GET("/users", userHandler.List)
    
    r.Run(fmt.Sprintf(":%d", cfg.Server.Port))
}
```

---

## 📁 Project Structure

```
10-configuration-management/
├── main.go
├── go.mod
├── configs/
│   ├── config.yaml          # Base config
│   ├── config-dev.yaml      # Dev overrides
│   └── config-prod.yaml     # Prod overrides
├── config/
│   ├── config.go            # Config structs & loader
│   └── validator.go         # Config validation
├── handler/
│   └── user_handler.go      # Example handler using config
└── TASK.md
```

---

## 📄 Config Files to Create

### configs/config.yaml (Base)
```yaml
server:
  host: "0.0.0.0"
  port: 8080
  read_timeout: 30
  write_timeout: 30

database:
  driver: postgres
  host: localhost
  port: 5432
  user: app
  password: ""  # Override via env!
  dbname: myapp
  sslmode: disable

jwt:
  secret: ""    # Override via env!
  expiration: 24

app:
  name: "My Go App"
  version: "1.0.0"
  debug: false
```

### configs/config-dev.yaml
```yaml
server:
  port: 3000

database:
  host: localhost
  dbname: myapp_dev

app:
  debug: true
```

### configs/config-prod.yaml
```yaml
server:
  port: 8080

database:
  host: ${DATABASE_HOST}  # From env
  sslmode: require

app:
  debug: false
```

---

## 💡 Key Patterns

### 1. Environment Variable Naming
```bash
# Viper convention (like Spring)
DATABASE_HOST=prod-db.aws.com      # -> database.host
JWT_SECRET=my-super-secret-key     # -> jwt.secret
SERVER_PORT=9000                   # -> server.port
```

### 2. Nested Struct with mapstructure
```go
type Config struct {
    Database struct {
        Connection struct {
            MaxOpen int `mapstructure:"max_open"`
            MaxIdle int `mapstructure:"max_idle"`
        } `mapstructure:"connection"`
    } `mapstructure:"database"`
}

// Corresponds to:
// database:
//   connection:
//     max_open: 10
//     max_idle: 5
```

### 3. Reading Individual Values
```go
// Direct access (less type-safe)
port := viper.GetInt("server.port")
host := viper.GetString("database.host")
debug := viper.GetBool("app.debug")
timeout := viper.GetDuration("server.timeout")
```

### 4. Hot Reload (Advanced)
```go
viper.WatchConfig()
viper.OnConfigChange(func(e fsnotify.Event) {
    log.Printf("Config changed: %s", e.Name)
    // Reload your services...
})
```

---

## 🆚 Spring Boot Comparison

| Spring Boot | Viper (Go) | Notes |
|-------------|------------|-------|
| `application.yml` | `config.yaml` | Same format! |
| `application-dev.yml` | `config-dev.yaml` | Profile-specific |
| `@ConfigurationProperties` | `mapstructure` tags | Struct mapping |
| `@Value("${db.host}")` | `viper.GetString("db.host")` | Direct access |
| `SPRING_DATASOURCE_URL` | `DATABASE_URL` | Env override |
| `@Validated` | Custom `Validate()` | Config validation |
| Spring Cloud Config | `viper.WatchConfig()` | Hot reload |

---

## ✅ Expected Behavior

```bash
# Run with default (dev) profile
go run main.go
# Output: Starting server on 0.0.0.0:3000 (dev port)

# Run with prod profile
APP_PROFILE=prod go run main.go
# Output: Starting server on 0.0.0.0:8080 (prod port)

# Override with env vars
DATABASE_HOST=mydb.aws.com SERVER_PORT=9000 go run main.go
# Output: Starting server on 0.0.0.0:9000
# Database host: mydb.aws.com

# Missing required config
JWT_SECRET="" go run main.go
# Output: Invalid config: JWT secret is required (exits with error)
```

---

## 🎓 What You'll Learn

1. **Viper basics** - Loading YAML/JSON config files
2. **Struct mapping** - Using `mapstructure` tags
3. **Environment variables** - Override config without code changes
4. **Profiles** - Dev/Prod/Staging configurations
5. **Validation** - Fail fast on invalid config
6. **Defaults** - Sensible fallback values
7. **12-Factor App principles** - Config from environment

---

## 🔐 Best Practices

1. **Never commit secrets** - Use env vars for passwords, API keys
2. **Fail fast** - Validate config at startup
3. **Use defaults wisely** - Good for dev, explicit for prod
4. **Document your config** - Comment each field
5. **Type-safe structs** - Prefer unmarshaling to `viper.GetString()`

---

## ⏱️ Estimated Time: 25-35 minutes

Configuration done right saves debugging time later! 🔧

