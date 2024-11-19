package sqs

import (
	"context"
	"testing"
	"time"
)

func TestInitSQS(t *testing.T) {
	tests := []struct {
		name        string
		config      SQSConfig
		wantErr     bool
		errContains string
	}{
		{
			name: "OK",
			config: SQSConfig{
				Environment: "test",
				Region:      "ap-northeast-1",
				Endpoint:    "http://localhost:4566",
				QueueNames: map[Key]string{
					SQSKeySample: "sample-queue",
				},
			},
			wantErr: false,
		},
		{
			name: "NG",
			config: SQSConfig{
				Environment: "invalid",
				Region:      "ap-northeast-1",
				Endpoint:    "http://localhost:4566",
				QueueNames: map[Key]string{
					SQSKeySample: "sample-queue",
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			client, err := InitSQS(ctx, tt.config)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error but got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if client == nil {
				t.Fatal("client is nil")
			}

			// Verify QueueURLs
			for key, queueName := range tt.config.QueueNames {
				url, exists := client.QueueURLs[key]
				if !exists {
					t.Errorf("QueueURL not found for key %v", key)
				}
				expectedURL := tt.config.Endpoint + "/000000000000/" + queueName
				if url != expectedURL {
					t.Errorf("incorrect QueueURL. expected: %s, got: %s", expectedURL, url)
				}
			}
		})
	}
}