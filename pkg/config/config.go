package config

import (
	"time"
)

const (
	Db_host                   = "localhost"
	Db_port                   = 5432
	Db_name                   = "production"
	Db_username               = "app_go"
	Db_password               = "app12345"
	APP_NAME                  = "http-server-go"
	JWT_SIGNATURE_KEY         = "kambing goes over the fence with machine gun"
	LOGIN_EXPIRATION_DURATION = time.Duration(1) * time.Hour
	Jwt_signing_method        = "HS256"
)
