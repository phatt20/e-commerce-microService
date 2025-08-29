package server

import (
	"context"
	"log"
	sagaQueue "microService/modules/saga/sagaHandler"
	"microService/modules/saga/sagaRepository"
	"microService/modules/saga/sagaUsecase"
	"microService/pkg/queue"
)

func (s *server) sagaService() {
	stateRepo := sagaRepository.NewMemoryRepo()
	idemRepo := sagaRepository.NewMemoryIdemRepo() 
	uc := sagaUsecase.New(
		stateRepo,
		idemRepo,
		sagaUsecase.Topics{
			OrderCmd:     "order.commands",
			PaymentCmd:   "payment.commands",
			InventoryCmd: "inventory.commands",
			ShippingCmd:  "shipping.commands",
		},
		[]string{s.cfg.Kafka.Url},
		s.cfg.Kafka.ApiKey,
		s.cfg.Kafka.Secret,
	)

	h := sagaQueue.NewSagaQueueHandler(uc)

	// Use ConsumerGroupOption struct
	opt := queue.ConsumerGroupOption{
		Brokers:              []string{s.cfg.Kafka.Url},
		GroupID:              "saga-service-group",
		APIKey:               s.cfg.Kafka.ApiKey,
		Secret:               s.cfg.Kafka.Secret,
		EnableTLS:            true,
		InsecureSkipTLSCheck: true,
		Version:              "1.0.0", // or use appropriate version string
	}

	cg, err := queue.NewConsumerGroup(opt)
	if err != nil {
		log.Fatalf("Failed to create consumer group: %v", err)
		return
	}

	// Create context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	topics := []string{"order.events", "inventory.events", "payment.events", "shipping.events"}

	// Start consumer in goroutine
	go func() {
		defer func() {
			if err := cg.Close(); err != nil {
				log.Printf("Error closing consumer group: %v", err)
			}
		}()

		// Use the helper function for better error handling
		if err := queue.RunConsumerGroup(ctx, cg, topics, h); err != nil {
			log.Printf("Error in consumer group: %v", err)
		}
	}()

	log.Println("Saga service consumer started successfully")
}