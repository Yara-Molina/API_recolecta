package infrastructure

import (
	"context"
	"fmt"
	"os"
	"sort"
	"strings"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"github.com/vicpoo/API_recolecta/src/notificacion/domain"
	"google.golang.org/api/option"
)

type FCMClient struct {
	client *messaging.Client
}

func NewFCMClient(credentialsFile string) (*FCMClient, error) {
	ctx := context.Background()
	credentialPath := strings.TrimSpace(credentialsFile)
	if credentialPath == "" {
		credentialPath = strings.TrimSpace(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"))
	}
	if credentialPath == "" {
		return nil, fmt.Errorf("missing FCM credentials path: set FCM_CREDENTIALS_FILE or GOOGLE_APPLICATION_CREDENTIALS")
	}

	info, err := os.Stat(credentialPath)
	if err != nil {
		return nil, fmt.Errorf("fcm credentials file not found at '%s': %w", credentialPath, err)
	}
	if info.IsDir() {
		return nil, fmt.Errorf("fcm credentials path points to a directory, expected file: '%s'", credentialPath)
	}

	app, err := firebase.NewApp(ctx, nil, option.WithCredentialsFile(credentialPath))
	if err != nil {
		return nil, fmt.Errorf("error initializing firebase app with credentials file '%s': %w", credentialPath, err)
	}

	client, err := app.Messaging(ctx)
	if err != nil {
		return nil, fmt.Errorf("error creating firebase messaging client: %w", err)
	}

	return &FCMClient{client: client}, nil
}

func (c *FCMClient) SendToTopic(ctx context.Context, topic string, notification *domain.PushNotification) error {
	dataPayload := make(map[string]string, len(notification.Data)+1)
	for key, value := range notification.Data {
		dataPayload[key] = value
	}
	if notification.Type != "" {
		dataPayload["notificationType"] = notification.Type
	}

	_, err := c.client.Send(ctx, &messaging.Message{
		Topic: topic,
		Notification: &messaging.Notification{
			Title: notification.Title,
			Body:  notification.Body,
		},
		Data: dataPayload,
	})
	return err
}

const fcmMaxTokensPerRequest = 500

func (c *FCMClient) Send(ctx context.Context, userTokens map[string]string, notification *domain.PushNotification) (map[string]domain.SendResult, error) {
	if len(userTokens) == 0 {
		return make(map[string]domain.SendResult), nil
	}
	userIDs := make([]string, 0, len(userTokens))
	for userID := range userTokens {
		userIDs = append(userIDs, userID)
	}
	sort.Strings(userIDs)

	dataPayload := make(map[string]string, len(notification.Data)+1)
	for key, value := range notification.Data {
		dataPayload[key] = value
	}
	if notification.Type != "" {
		dataPayload["notificationType"] = notification.Type
	}

	results := make(map[string]domain.SendResult, len(userIDs))

	for start := 0; start < len(userIDs); start += fcmMaxTokensPerRequest {
		end := start + fcmMaxTokensPerRequest
		if end > len(userIDs) {
			end = len(userIDs)
		}
		batchUserIDs := userIDs[start:end]

		tokens := make([]string, 0, len(batchUserIDs))
		for _, userID := range batchUserIDs {
			tokens = append(tokens, userTokens[userID])
		}

		message := &messaging.MulticastMessage{
			Tokens: tokens,
			Notification: &messaging.Notification{
				Title: notification.Title,
				Body:  notification.Body,
			},
			Data: dataPayload,
		}

		batchResponse, err := c.client.SendEachForMulticast(ctx, message)
		if err != nil {
			return nil, err
		}

		for index, response := range batchResponse.Responses {
			userID := batchUserIDs[index]
			if response.Success {
				results[userID] = domain.SendResult{Success: true}
				continue
			}

			errorMessage := "unknown error"
			if response.Error != nil {
				errorMessage = response.Error.Error()
			}

			results[userID] = domain.SendResult{
				Success: false,
				Error:   errorMessage,
			}
		}
	}

	return results, nil
}
