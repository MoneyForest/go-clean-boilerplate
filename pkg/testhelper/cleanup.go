package testhelper

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

func Cleanup(ctx context.Context, g *Gateway) {
	if g.MySQLClient != nil {
		rows, err := g.MySQLClient.Query("SHOW TABLES")
		if err != nil {
			log.Printf("failed to get tables: %v", err)
		} else {
			defer rows.Close()
			if _, err := g.MySQLClient.Exec("SET FOREIGN_KEY_CHECKS = 0"); err != nil {
				log.Printf("failed to disable foreign key checks: %v", err)
			}
			for rows.Next() {
				var tableName string
				if err := rows.Scan(&tableName); err != nil {
					log.Printf("failed to scan table name: %v", err)
					continue
				}
				if _, err := g.MySQLClient.Exec("TRUNCATE TABLE " + tableName); err != nil {
					log.Printf("failed to truncate table %s: %v", tableName, err)
				}
			}
			if _, err := g.MySQLClient.Exec("SET FOREIGN_KEY_CHECKS = 1"); err != nil {
				log.Printf("failed to enable foreign key checks: %v", err)
			}
		}
		if err := g.MySQLClient.Close(); err != nil {
			log.Printf("failed to close MySQL connection: %v", err)
		}
	}

	if g.RedisClient != nil {
		if err := g.RedisClient.FlushAll(ctx).Err(); err != nil {
			log.Printf("failed to flush Redis: %v", err)
		}
		if err := g.RedisClient.Close(); err != nil {
			log.Printf("failed to close Redis connection: %v", err)
		}
	}

	if g.SQSClient != nil {
		result, err := g.SQSClient.Client.ListQueues(ctx, &sqs.ListQueuesInput{})
		if err != nil {
			log.Printf("failed to list SQS queues: %v", err)
		} else {
			for _, queueURL := range result.QueueUrls {
				_, err := g.SQSClient.Client.PurgeQueue(ctx, &sqs.PurgeQueueInput{
					QueueUrl: aws.String(queueURL),
				})
				if err != nil {
					log.Printf("failed to purge queue %s: %v", queueURL, err)
				}
			}
		}
	}
}