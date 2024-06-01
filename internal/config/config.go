package config

import (
	"fmt"

	"go-admin/pkg/encoding/json"
	"go-admin/pkg/logging"
)

type Config struct {
	Logger     logging.LoggerConfig
	General    General
	Storage    Storage
	Middleware Middleware
	Util       Util
	Dictionary Dictionary
}

type General struct {
	AppName            string `default:"go-admin"`
	Version            string `default:"v1.0.0"`
	Debug              bool
	PprofAddr          string
	DisableSwagger     bool
	DisablePrintConfig bool
	DefaultLoginPwd    string `default:"6351623c8cef86fefabfa7da046fc619"` // MD5(abc-123)
	WorkDir            string // From command arguments
	ThirdPartyDir      string `default:"./third_party"`
	MenuFile           string // From schema.Menus (JSON/YAML)
	DenyDeleteMenu     bool
	HTTP               struct {
		Addr            string `default:":5001"`
		ShutdownTimeout int    `default:"10"` // seconds
		ReadTimeout     int    `default:"60"` // seconds
		WriteTimeout    int    `default:"60"` // seconds
		IdleTimeout     int    `default:"10"` // seconds
		CertFile        string
		KeyFile         string
	}
	Root struct {
		ID       string `default:"root"`
		Email    string `default:"admin@admin.com"`
		FullName string `default:"Admin"`
		Password string
	}
}

type Storage struct {
	Cache struct {
		Type      string `default:"memory"` // memory/badger/redis
		Delimiter string `default:":"`      // delimiter for key
		Memory    struct {
			CleanupInterval int `default:"60"` // seconds
		}
		Badger struct {
			Path string `default:"data/cache"`
		}
		Redis struct {
			Addr     string
			Username string
			Password string
			DB       int
		}
	}
	DB struct {
		Debug        bool
		Type         string `default:"postgres"`                                                                              // sqlite3/mysql/postgres
		DSN          string `default:"host=postgres user=postgres password=123456 dbname=go-admin port=5432 sslmode=disable"` // database source name
		MaxLifetime  int    `default:"86400"`                                                                                 // seconds
		MaxIdleTime  int    `default:"3600"`                                                                                  // seconds
		MaxOpenConns int    `default:"100"`                                                                                   // connections
		MaxIdleConns int    `default:"50"`                                                                                    // connections
		TablePrefix  string `default:""`
		AutoMigrate  bool
		Resolver     []struct {
			DBType   string   // sqlite3/mysql/postgres
			Sources  []string // DSN
			Replicas []string // DSN
			Tables   []string
		}
	}
}

type Util struct {
	Captcha struct {
		Length    int    `default:"4"`
		Width     int    `default:"400"`
		Height    int    `default:"160"`
		CacheType string `default:"memory"` // memory/redis
		Redis     struct {
			Addr      string
			Username  string
			Password  string
			DB        int
			KeyPrefix string `default:"captcha:"`
		}
	}
	Prometheus struct {
		Enable         bool
		Port           int    `default:"9100"`
		BasicUsername  string `default:"admin"`
		BasicPassword  string `default:"admin"`
		LogApis        []string
		LogMethods     []string
		DefaultCollect bool
	}
}

type Dictionary struct {
	UserCacheExp int `default:"4"` // hours
}

func (c *Config) IsDebug() bool {
	return c.General.Debug
}

func (c *Config) String() string {
	b, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		panic("Failed to marshal config: " + err.Error())
	}
	return string(b)
}

func (c *Config) PreLoad() {
	if addr := c.Storage.Cache.Redis.Addr; addr != "" {
		username := c.Storage.Cache.Redis.Username
		password := c.Storage.Cache.Redis.Password
		if c.Util.Captcha.CacheType == "redis" &&
			c.Util.Captcha.Redis.Addr == "" {
			c.Util.Captcha.Redis.Addr = addr
			c.Util.Captcha.Redis.Username = username
			c.Util.Captcha.Redis.Password = password
		}
		if c.Middleware.RateLimiter.Store.Type == "redis" &&
			c.Middleware.RateLimiter.Store.Redis.Addr == "" {
			c.Middleware.RateLimiter.Store.Redis.Addr = addr
			c.Middleware.RateLimiter.Store.Redis.Username = username
			c.Middleware.RateLimiter.Store.Redis.Password = password
		}
		if c.Middleware.Auth.Store.Type == "redis" &&
			c.Middleware.Auth.Store.Redis.Addr == "" {
			c.Middleware.Auth.Store.Redis.Addr = addr
			c.Middleware.Auth.Store.Redis.Username = username
			c.Middleware.Auth.Store.Redis.Password = password
		}
	}
}

func (c *Config) Print() {
	if c.General.DisablePrintConfig {
		return
	}
	fmt.Println("// ----------------------- Load configurations start ------------------------")
	fmt.Println(c.String())
	fmt.Println("// ----------------------- Load configurations end --------------------------")
}

func (c *Config) FormatTableName(name string) string {
	return c.Storage.DB.TablePrefix + name
}
