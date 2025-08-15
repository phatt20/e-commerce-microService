package userHandler

import (
	"context"
	"log"
	"microService/config"
	user "microService/modules/user"
	userUsecase "microService/modules/user/userUsecase"
	"microService/pkg/queue"
	"os"
	"os/signal"
	"syscall"

	"github.com/IBM/sarama"
)

type (
	UserQueueHandlerService interface {
		DockedUserMoney()
		AddUserMoney()
		RollbackUserTransaction()
	}

	userQueueHandler struct {
		cfg         *config.Config
		userUsecase userUsecase.UserUsecase
	}
)

func NewUserQueueHandler(cfg *config.Config, userUsecase userUsecase.UserUsecase) UserQueueHandlerService {
	return &userQueueHandler{
		cfg:         cfg,
		userUsecase: userUsecase,
	}
}

func (h *userQueueHandler) UserConsumer(pctx context.Context) (sarama.PartitionConsumer, error) {
	worker, err := queue.ConnectConsumer([]string{h.cfg.Kafka.Url}, h.cfg.Kafka.ApiKey, h.cfg.Kafka.Secret)
	if err != nil {
		return nil, err
	}

	offset, err := h.userUsecase.GetOffset(pctx)
	if err != nil {
		return nil, err
	}

	consumer, err := worker.ConsumePartition("user", 0, offset)
	if err != nil {
		log.Println("Trying to set offset as 0")
		consumer, err = worker.ConsumePartition("user", 0, 0)
		if err != nil {
			log.Println("Error: PaymentConsumer failed: ", err.Error())
			return nil, err
		}
	}

	return consumer, nil
}

func (h *userQueueHandler) DockedUserMoney() {
	ctx := context.Background()

	consumer, err := h.UserConsumer(ctx)
	if err != nil {
		return
	}
	defer consumer.Close()

	log.Println("Start DockedUserMoney ...")

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case err := <-consumer.Errors():
			log.Println("Error: DockedUserMoney failed: ", err.Error())
			continue
		case msg := <-consumer.Messages():
			if string(msg.Key) == "buy" {
				h.userUsecase.UpserOffset(ctx, msg.Offset+1)

				req := new(user.CreateUserTransactionReq)

				if err := queue.DecodeMessage(req, msg.Value); err != nil {
					continue
				}

				h.userUsecase.DockedUserMoneyRes(ctx, h.cfg, req)

				log.Printf("DockedUserMoney | Topic(%s)| Offset(%d) Message(%s) \n", msg.Topic, msg.Offset, string(msg.Value))
			}
		case <-sigchan:
			log.Println("Stop DockedUserMoney...")
			return
		}
	}
}

func (h *userQueueHandler) AddUserMoney() {
	ctx := context.Background()

	consumer, err := h.UserConsumer(ctx)
	if err != nil {
		return
	}
	defer consumer.Close()

	log.Println("Start AddUserMoney ...")

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case err := <-consumer.Errors():
			log.Println("Error: AddUserMoney failed: ", err.Error())
			continue
		case msg := <-consumer.Messages():
			if string(msg.Key) == "sell" {
				h.userUsecase.UpserOffset(ctx, msg.Offset+1)

				req := new(user.CreateUserTransactionReq)

				if err := queue.DecodeMessage(req, msg.Value); err != nil {
					continue
				}

				h.userUsecase.AddUserMoneyRes(ctx, h.cfg, req)

				log.Printf("AddUserMoney | Topic(%s)| Offset(%d) Message(%s) \n", msg.Topic, msg.Offset, string(msg.Value))
			}
		case <-sigchan:
			log.Println("Stop AddUserMoney...")
			return
		}
	}
}

func (h *userQueueHandler) RollbackUserTransaction() {
	ctx := context.Background()

	consumer, err := h.UserConsumer(ctx)
	if err != nil {
		return
	}
	defer consumer.Close()

	log.Println("Start RollbackUserTransaction ...")

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case err := <-consumer.Errors():
			log.Println("Error: RollbackUserTransaction failed: ", err.Error())
			continue
		case msg := <-consumer.Messages():
			if string(msg.Key) == "rtransaction" {
				h.userUsecase.UpserOffset(ctx, msg.Offset+1)

				req := new(user.RollbackUserTransactionReq)

				if err := queue.DecodeMessage(req, msg.Value); err != nil {
					continue
				}

				h.userUsecase.RollbackUserTransaction(ctx, req)

				log.Printf("RollbackUserTransaction | Topic(%s)| Offset(%d) Message(%s) \n", msg.Topic, msg.Offset, string(msg.Value))
			}
		case <-sigchan:
			log.Println("Stop RollbackUserTransaction...")
			return
		}
	}
}
