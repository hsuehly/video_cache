package bootstrap

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
)

type Env struct {
	AppEnv         string `json:"APP_ENV"`
	PORT           string `json:"PORT"`
	DBUser         string `json:"MYSQL_USER"`
	DBPass         string `json:"MYSQL_PASS"`
	DBHost         string `json:"MYSQL_HOST"`
	DBPort         int    `json:"MYSQL_PORT"`
	DBDaTa         string `json:"MYSQL_DATA"`
	ContextTimeout int    `json:"CONTEXT_TIMEOUT"`
	CacheKey       string `json:"CACHE_KEY"`
	RedisAddr      string `json:"REDIS_ADDR"`
	RedisPass      string `json:"REDIS_PASS"`
	ToPic          string `json:"TOPIC"`
	ConsumerGroup  string `json:"CONSUMER_GROUP"`
	ConsumerID     string `json:"CONSUMER_ID"`
	Yun139         string `json:"YUN139"`
	PanQQ          string `json:"PAN_QQ"`
	PanQY          string `json:"PAN_QY"`
	PanYK          string `json:"PAN_YK"`
	PanMG          string `json:"PAN_MG"`
}

var defaultEnv = map[string]interface{}{
	"APP_ENV":         "development", // development production
	"PORT":            ":8080",
	"MYSQL_USER":      "root",
	"MYSQL_PASS":      "hsueh095.",
	"MYSQL_HOST":      "localhost",
	"MYSQL_PORT":      3306,
	"MYSQL_DATA":      "m3u8_cache",
	"CONTEXT_TIMEOUT": 2,
	"CACHE_KEY":       "HsuehLY",
	"REDIS_ADDR":      "localhost:6379",
	"REDIS_PASS":      "hsueh095",
	"TOPIC":           "m3u8_topic",
	"CONSUMER_GROUP":  "m3u8_group",
	"CONSUMER_ID":     "m3u8_consumer",
	"YUN139":          "bW9iaWxlOjE3NjI5ODg2NjA1OmdXRW1aRDd0fDF8UkNTfDE3MTc0OTYzMzEwNDV8QVEuWXguVERTTGdzdEVrQTU4Q2VSRFlwQmpOa1NHYTFCSjQ5X3h6VHVxM1BYZGxHX3lRbFptcC5WUWNpUWg1Z1puOWdaR2x0cE5Zd2ZFUXZpa29odEduNEJ3WFUxSzB0NG5Pd3BTa1JieUZJRmxOT0Y1UFZfS3UzMkM0U3phNndqRVU0cXRwazl6VzZRRW00bFdaZG94TlhaRGR3YU50Z016eW1fTTVLRDdVLQ==",
	"PAN_QQ":          "1L11qfFSrFNi29620240421190742ve4",
	"PAN_QY":          "1L11qfFSrFNi29820240421190846mwz",
	"PAN_YK":          "1L11qfFSrFNi30120240421190916yt0",
	"PAN_MG":          "1L11qfFSrFNi29920240421190932lrt",
}

func NewEnv() *Env {
	env := &Env{}

	elem := reflect.ValueOf(env).Elem()
	elemType := elem.Type()

	for i := 0; i < elem.NumField(); i++ {
		field := elem.Field(i)
		fieldType := elemType.Field(i)
		envVarName := fieldType.Tag.Get("json")

		envVarValue := os.Getenv(envVarName)
		if envVarValue == "" {
			// Use default value if environment variable is not set
			envVarValue = fmt.Sprintf("%v", defaultEnv[envVarName])
		}

		switch field.Kind() {
		case reflect.String:
			field.SetString(envVarValue)
		case reflect.Int:
			num, err := strconv.Atoi(envVarValue)
			if err != nil {
				panic(fmt.Sprintf("Error converting %s to int: %s", envVarName, err))
			}
			field.SetInt(int64(num))
		}
	}

	return env
}
