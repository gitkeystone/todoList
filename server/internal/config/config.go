// Package config 负责从环境变量加载服务配置（PRD §9.3）。
package config

import (
	"os"
	"strings"
)

// Config 服务运行配置。
type Config struct {
	Port           string   // 后端监听端口
	DBPath         string   // SQLite 数据库文件路径
	GinMode        string   // Gin 运行模式：debug / release
	AllowedOrigins []string // CORS 白名单
	WebDist        string   // 生产模式下前端静态资源目录
}

// Load 从环境变量加载配置，缺失项回退到默认值（对应 .env.example）。
func Load() *Config {
	return &Config{
		Port:           getEnv("PORT", "8080"),
		DBPath:         getEnv("DB_PATH", "data/todolist.db"),
		GinMode:        getEnv("GIN_MODE", "debug"),
		AllowedOrigins: splitEnv("ALLOWED_ORIGINS", "http://localhost:5173"),
		WebDist:        getEnv("WEB_DIST", "web/dist"),
	}
}

// DSN 返回 SQLite 连接串（含 PRAGMA，PRD §7.1）。
func (c *Config) DSN() string {
	return c.DBPath +
		"?_pragma=journal_mode(WAL)" +
		"&_pragma=busy_timeout(5000)" +
		"&_pragma=foreign_keys(ON)" +
		"&_pragma=synchronous(NORMAL)"
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

// splitEnv 解析逗号分隔的环境变量为切片。
func splitEnv(key, def string) []string {
	v := os.Getenv(key)
	if v == "" {
		v = def
	}
	parts := strings.Split(v, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if p = strings.TrimSpace(p); p != "" {
			out = append(out, p)
		}
	}
	return out
}
