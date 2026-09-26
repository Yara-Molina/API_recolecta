package config

import "os"

type Config struct {
	RedisHost                string
	RedisPort                string
	RedisPassword            string
	FCMCredentialsFile       string
	ModeloReportesURL        string
	ClasificadorURL          string
	AnomaliaCreadaWebhookURL string
	APIRutasURL              string
}

func LoadConfig() (*Config, error) {
	cfg := &Config{
		RedisHost:          os.Getenv("REDIS_HOST"),
		RedisPort:          os.Getenv("REDIS_PORT"),
		RedisPassword:      os.Getenv("REDIS_PASSWORD"),
		FCMCredentialsFile: os.Getenv("FCM_CREDENTIALS_FILE"),
		ModeloReportesURL: getEnvOrDefault("MODELO_REPORTES_URL"),
		ClasificadorURL:   getEnvOrDefault("CLASIFICADOR_URL"),
		AnomaliaCreadaWebhookURL: getEnvOrDefault("ANOMALIA_CREADA_WEBHOOK_URL"),

		APIRutasURL: os.Getenv("API_RUTAS_URL"),
	}
	return cfg, nil
}

func getEnvOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
