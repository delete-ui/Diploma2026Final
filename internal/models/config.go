package models

import "time"

type Config struct {
	Env        string
	LogLevel   string
	LogFormat  string
	HTTPServer HTTPServerConfig
	Database   DatabaseConfig
	JWT        JWTConfig
	SMTP       SMTPConfig
}

type HTTPServerConfig struct {
	Address      string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type JWTConfig struct {
	Secret    string
	AccessTTL time.Duration
}

type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
	FromName string
}
