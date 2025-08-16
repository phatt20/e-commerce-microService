package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type (
	Config struct {
		App      App
		Mongo    *MongoConfig
		Postgres *PostgresConfig
		Jwt      Jwt
		Kafka    Kafka
		Grpc     Grpc
		Paginate Paginate
	}

	App struct {
		Name  string
		Url   string
		Stage string
	}

	MongoConfig struct {
		Url string
	}

	PostgresConfig struct {
		Host     string
		Port     int
		User     string
		Password string
		DBName   string
		SSLMode  string
		Schema   string
	}

	Jwt struct {
		AccessSecretKey  string
		RefreshSecretKey string
		ApiSecretKey     string
		AccessDuration   int64
		RefreshDuration  int64
		ApiDuration      int64
	}

	Kafka struct {
		Url    string
		ApiKey string
		Secret string
	}

	Grpc struct {
		AuthUrl      string
		UserUrl      string
		InventoryUrl string
		ItemUrl      string
		PaymentUrl   string
	}

	Paginate struct {
		ItemNextPageBasedUrl      string
		InventoryNextPageBasedUrl string
	}
)

// LoadConfig อ่านค่า config จาก .env + ENV
func LoadConfig(path string) Config {
	if err := godotenv.Load(path); err != nil {
		log.Println("⚠️ Warning: .env not found, using system ENV only")
	}

	// Mongo (optional)
	var mongo *MongoConfig
	if os.Getenv("MONGO_URL") != "" {
		mongo = &MongoConfig{
			Url: os.Getenv("MONGO_URL"),
		}
	}

	// Postgres (optional)
	var pg *PostgresConfig
	if os.Getenv("PG_HOST") != "" {
		port, _ := strconv.Atoi(os.Getenv("PG_PORT"))
		pg = &PostgresConfig{
			Host:     os.Getenv("PG_HOST"),
			Port:     port,
			User:     os.Getenv("PG_USER"),
			Password: os.Getenv("PG_PASSWORD"),
			DBName:   os.Getenv("PG_DBNAME"),
			SSLMode:  os.Getenv("PG_SSLMODE"),
			Schema:   os.Getenv("PG_SCHEMA"),
		}
	}

	return Config{
		App: App{
			Name:  os.Getenv("APP_NAME"),
			Url:   os.Getenv("APP_URL"),
			Stage: os.Getenv("APP_STAGE"),
		},
		Mongo:    mongo,
		Postgres: pg,
		Jwt: Jwt{
			AccessSecretKey:  os.Getenv("JWT_ACCESS_SECRET_KEY"),
			RefreshSecretKey: os.Getenv("JWT_REFRESH_SECRET_KEY"),
			ApiSecretKey:     os.Getenv("JWT_API_SECRET_KEY"),
			AccessDuration:   mustParseInt(os.Getenv("JWT_ACCESS_DURATION")),
			RefreshDuration:  mustParseInt(os.Getenv("JWT_REFRESH_DURATION")),
			ApiDuration:      mustParseInt(os.Getenv("JWT_API_DURATION")),
		},
		Kafka: Kafka{
			Url:    os.Getenv("KAFKA_URL"),
			ApiKey: os.Getenv("KAFKA_API_KEY"),
			Secret: os.Getenv("KAFKA_SECRET"),
		},
		Grpc: Grpc{
			AuthUrl:      os.Getenv("GRPC_AUTH_URL"),
			UserUrl:      os.Getenv("GRPC_USER_URL"),
			ItemUrl:      os.Getenv("GRPC_ITEM_URL"),
			InventoryUrl: os.Getenv("GRPC_INVENTORY_URL"),
			PaymentUrl:   os.Getenv("GRPC_PAYMENT_URL"),
		},
		Paginate: Paginate{
			ItemNextPageBasedUrl:      os.Getenv("PAGINATE_ITEM_NEXT_PAGE_BASED_URL"),
			InventoryNextPageBasedUrl: os.Getenv("PAGINATE_INVENTORY_NEXT_PAGE_BASED_URL"),
		},
	}
}

func mustParseInt(s string) int64 {
	if s == "" {
		return 0
	}
	result, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		log.Fatalf("❌ Error parsing int from string: %s", s)
	}
	return result
}
