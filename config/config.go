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
}

func LoadConfig() (*Config, error) {
	cfg := &Config{
		RedisHost:          os.Getenv("REDIS_HOST"),
		RedisPort:          os.Getenv("REDIS_PORT"),
		RedisPassword:      os.Getenv("REDIS_PASSWORD"),
		FCMCredentialsFile: os.Getenv("FCM_CREDENTIALS_FILE"),
		// Nombres de servicio de docker.compose.yml: alcanzables por DNS interno
		// de Docker dentro de app_internal_net, tanto en dev como en despliegue
		// real (mismo compose stack). Si no estan seteadas, caen a esos defaults.
		ModeloReportesURL: getEnvOrDefault("MODELO_REPORTES_URL", "http://modelo_reportes:8000"),
		ClasificadorURL:   getEnvOrDefault("CLASIFICADOR_URL", "http://clasificador_reportes:8001"),
		// Webhook externo: servicio api-rutas (Node.js, equipo de rutas/AG) que
		// se notifica cuando se crea una anomalia con coordenadas. Verificado
		// en produccion (GET /health -> {"status":"ok"} en
		// https://api-rutas.practicasoftware.fun) -- ya no es localhost:8004.
		// api-rutas es quien internamente consulta sus rutas/puntos, llama al
		// algoritmo genetico de rutas (AG, https://ag.practicasoftware.fun) y
		// notifica a los conductores por su propio servicio de WebSocket
		// (wss://websocket.practicasoftware.fun); gin-backend no necesita
		// hablar con el AG ni con ese WebSocket directamente. Ver
		// docs/implementacion-fix-anomalias-y-ag.md.
		AnomaliaCreadaWebhookURL: getEnvOrDefault("ANOMALIA_CREADA_WEBHOOK_URL", "https://api-rutas.practicasoftware.fun/anomalia_creada"),
	}
	// Aqui se pueden anadir validaciones (ej. que no esten vacios)
	return cfg, nil
}

func getEnvOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
