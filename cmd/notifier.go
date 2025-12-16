package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-aws/sqs"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/goccy/go-json"
	"github.com/saleh-ghazimoradi/Cartopher/config"
	"github.com/saleh-ghazimoradi/Cartopher/internal/domain"
	"github.com/saleh-ghazimoradi/Cartopher/internal/service"
	"github.com/saleh-ghazimoradi/Cartopher/pkg/uploadProvider"

	"github.com/spf13/cobra"
)

// notifierCmd represents the notifier command
var notifierCmd = &cobra.Command{
	Use:   "notifier",
	Short: "It's responsible for notification purposes",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("notifier called")

		log.Println("Starting notification service...")

		ctx := context.Background()

		cfg, err := config.GetInstance()
		if err != nil {
			log.Fatalf("failed to get config: %v", err)
		}

		emailNotifier := service.NewEmailNotifier(cfg)

		awsCfg, err := uploadProvider.CreateAWSConfig(ctx, cfg.AWS.S3Endpoint, cfg.AWS.Region)
		if err != nil {
			log.Fatalf("failed to create aws config: %v", err)
		}

		logger := watermill.NewStdLogger(false, false)

		subscriber, err := sqs.NewSubscriber(sqs.SubscriberConfig{
			AWSConfig: awsCfg,
		}, logger)

		if err != nil {
			log.Fatalf("failed to create aws sqs subscriber: %v", err)
		}

		messages, err := subscriber.Subscribe(ctx, cfg.AWS.EventQueueName)
		if err != nil {
			subscriber.Close()
			log.Fatalf("failed to subscribe to queue: %v", err)
		}

		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

		log.Println("Notification service started. Waiting for messages...")

		for {
			select {
			case msg := <-messages:
				if err := processMessage(msg, emailNotifier); err != nil {
					log.Printf("failed to process message: %v", err)
					msg.Nack()
				} else {
					msg.Ack()
				}
			case <-sigChan:
				log.Println("Notification service shutting down...")
				subscriber.Close()
				return
			}
		}
	},
}

func processMessage(msg *message.Message, emailNotifier service.Notifier) error {
	eventType := msg.Metadata.Get("event_type")

	switch eventType {
	case service.UserLoggedIn:
		return handleUserLoggedIn(msg, emailNotifier)
	default:
		log.Printf("Unknown event type: %s", eventType)
		return nil
	}
}

func handleUserLoggedIn(msg *message.Message, emailNotifier service.Notifier) error {
	var user domain.User

	if err := json.Unmarshal(msg.Payload, &user); err != nil {
		return err
	}

	userName := user.FirstName + " " + user.LastName
	if userName == " " {
		userName = "User"
	}

	log.Printf("Sending login notification for: %s", user.Email)

	return emailNotifier.SendLoginNotification(user.Email, userName)
}

func init() {
	rootCmd.AddCommand(notifierCmd)
}
