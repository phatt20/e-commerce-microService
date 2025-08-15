// package queue

// import (
// 	"crypto/tls"
// 	"encoding/json"
// 	"errors"
// 	"log"

// 	"github.com/IBM/sarama"
// 	"github.com/go-playground/validator/v10"
// )

// func ConnectProducer(brokerUrls []string, apiKey, secret string) (sarama.SyncProducer, error) {
// 	config := sarama.NewConfig()
// 	if apiKey != "" && secret != "" {
// 		config.Net.SASL.Enable = true
// 		config.Net.SASL.User = apiKey
// 		config.Net.SASL.Password = secret
// 		config.Net.SASL.Mechanism = "PLAIN"
// 		config.Net.SASL.Handshake = true
// 		config.Net.SASL.Version = sarama.SASLHandshakeV1
// 		config.Net.TLS.Enable = true
// 		config.Net.TLS.Config = &tls.Config{
// 			InsecureSkipVerify: true,
// 			ClientAuth:         tls.NoClientCert,
// 		}
// 	}
// 	config.Producer.Return.Successes = true
// 	config.Producer.RequiredAcks = sarama.WaitForAll
// 	config.Producer.Retry.Max = 3

// 	producer, err := sarama.NewSyncProducer(brokerUrls, config)
// 	if err != nil {
// 		log.Printf("Error: Failed to connect to producer: %s", err.Error())
// 		return nil, errors.New("error: failed to connect to producer")
// 	}
// 	return producer, nil
// }

// func PushMessageWithKeyToQueue(brokerUrls []string, apiKey, secret, topic, key string, message []byte) error {
// 	producer, err := ConnectProducer(brokerUrls, apiKey, secret)
// 	if err != nil {
// 		log.Printf("Error: Failed to connect to producer: %s", err.Error())
// 		return errors.New("error: failed to connect to producer")
// 	}
// 	defer producer.Close()

// 	msg := &sarama.ProducerMessage{
// 		Topic: topic,
// 		Value: sarama.StringEncoder(message),
// 		Key:   sarama.StringEncoder(key),
// 	}

// 	partition, offset, err := producer.SendMessage(msg)
// 	if err != nil {
// 		log.Printf("Error: Failed to send message: %s", err.Error())
// 		return errors.New("error: failed to send message")
// 	}
// 	log.Printf("Message is stored in topic(%s)/partition(%d)/offset(%d)\n", topic, partition, offset)

// 	return nil
// }

// func ConnectConsumer(brokerUrls []string, apiKey, secret string) (sarama.Consumer, error) {
// 	config := sarama.NewConfig()
// 	if apiKey != "" && secret != "" {
// 		config.Net.SASL.Enable = true
// 		config.Net.SASL.User = apiKey
// 		config.Net.SASL.Password = secret
// 		config.Net.SASL.Mechanism = "PLAIN"
// 		config.Net.SASL.Handshake = true
// 		config.Net.SASL.Version = sarama.SASLHandshakeV1
// 		config.Net.TLS.Enable = true
// 		config.Net.TLS.Config = &tls.Config{
// 			InsecureSkipVerify: true,
// 			ClientAuth:         tls.NoClientCert,
// 		}
// 	}
// 	config.Consumer.Return.Errors = true
// 	config.Consumer.Fetch.Max = 3

// 	consumer, err := sarama.NewConsumer(brokerUrls, config)
// 	if err != nil {
// 		log.Printf("Error: Failed to connect to consumer: %s", err.Error())
// 		return nil, errors.New("error: failed to connect to consumer")
// 	}

// 	return consumer, nil
// }

// func DecodeMessage(obj any, value []byte) error {
// 	if err := json.Unmarshal(value, &obj); err != nil {
// 		log.Printf("Error: Failed to decode message: %s", err.Error())
// 		return errors.New("error: failed to decode message")
// 	}

// 	validate := validator.New()
// 	if err := validate.Struct(obj); err != nil {
// 		log.Printf("Error: Failed to validate message: %s", err.Error())
// 		return errors.New("error: failed to validate message")
// 	}

//		return nil
//	}
package queue

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/IBM/sarama"
	"github.com/go-playground/validator"
)

type ProducerOption struct {
	Brokers              []string
	APIKey               string // SASL/PLAIN user (optional)
	Secret               string // SASL/PLAIN password (optional)
	EnableTLS            bool   // true เมื่อใช้ TLS จริง (อย่า skip verify ในโปรดักชัน)
	InsecureSkipTLSCheck bool   // ใช้แค่ตอนทดสอบเท่านั้น
	Acks                 sarama.RequiredAcks
	RetryMax             int
	Version              string // e.g. "3.5.0"
}

func defaultProducerOption() ProducerOption {
	return ProducerOption{
		Acks:     sarama.WaitForAll,
		RetryMax: 5,
		Version:  "3.5.0",
	}
}

func newSaramaVersion(v string) sarama.KafkaVersion {
	ver, err := sarama.ParseKafkaVersion(v)
	if err != nil {
		return sarama.V3_5_0_0
	}
	return ver
}

func ConnectProducer(brokerUrls []string, apiKey, secret string) (sarama.SyncProducer, error) {
	// Backward-compatible wrapper (เดิม)
	opt := defaultProducerOption()
	opt.Brokers = brokerUrls
	opt.APIKey = apiKey
	opt.Secret = secret
	opt.EnableTLS = (apiKey != "" && secret != "") // เดิมคุณเปิด TLS เมื่อตั้ง SASL
	return NewProducer(opt)
}

// แนะนำให้ใช้ฟังก์ชันนี้แทน (ปรับ option ได้ละเอียดกว่า)
func NewProducer(opt ProducerOption) (sarama.SyncProducer, error) {
	cfg := sarama.NewConfig()
	cfg.Version = newSaramaVersion(opt.Version)

	// Producer reliability
	cfg.Producer.Return.Successes = true
	cfg.Producer.RequiredAcks = opt.Acks
	cfg.Producer.Retry.Max = opt.RetryMax
	cfg.Producer.Idempotent = true                    // เปิด idempotent producer
	cfg.Producer.Compression = sarama.CompressionZSTD // ช่วยลด bandwidth/latency (ปรับได้)

	// Partitioning: ถ้ามี Key Sarama จะใช้ hash partitioner ให้เองอยู่แล้ว

	// Security (SASL/TLS)
	if opt.APIKey != "" && opt.Secret != "" {
		cfg.Net.SASL.Enable = true
		cfg.Net.SASL.User = opt.APIKey
		cfg.Net.SASL.Password = opt.Secret
		cfg.Net.SASL.Mechanism = sarama.SASLTypePlaintext
		cfg.Net.SASL.Handshake = true
		cfg.Net.SASL.Version = sarama.SASLHandshakeV1
	}

	if opt.EnableTLS {
		cfg.Net.TLS.Enable = true
		cfg.Net.TLS.Config = &tls.Config{
			// โปรดตั้งเป็น false ในโปรดักชันและใส่ RootCAs ให้ถูกต้อง
			InsecureSkipVerify: opt.InsecureSkipTLSCheck,
			ClientAuth:         tls.NoClientCert,
			MinVersion:         tls.VersionTLS12,
		}
	}

	p, err := sarama.NewSyncProducer(opt.Brokers, cfg)
	if err != nil {
		log.Printf("Error: Failed to create producer: %v", err)
		return nil, errors.New("error: failed to connect to producer")
	}
	return p, nil
}

func PushMessageWithKeyToQueue(brokerUrls []string, apiKey, secret, topic, key string, message []byte) error {
	// Backward-compatible (แต่ไม่แนะนำเพราะสร้าง producer ทุกครั้ง)
	p, err := ConnectProducer(brokerUrls, apiKey, secret)
	if err != nil {
		log.Printf("Error: Failed to connect to producer: %s", err.Error())
		return errors.New("error: failed to connect to producer")
	}
	defer p.Close()

	msg := &sarama.ProducerMessage{
		Topic:     topic,
		Value:     sarama.ByteEncoder(message),
		Key:       sarama.StringEncoder(key),
		Timestamp: time.Now(),
	}

	partition, offset, err := p.SendMessage(msg)
	if err != nil {
		log.Printf("Error: Failed to send message: %s", err.Error())
		return errors.New("error: failed to send message")
	}
	log.Printf("Message stored topic(%s)/partition(%d)/offset(%d)", topic, partition, offset)
	return nil
}

// ---------- Consumer (legacy) ----------
// คงฟังก์ชันเดิมไว้ แต่แนะนำใช้ ConsumerGroup ด้านล่างแทน

func ConnectConsumer(brokerUrls []string, apiKey, secret string) (sarama.Consumer, error) {
	cfg := sarama.NewConfig()
	cfg.Version = sarama.V3_5_0_0

	if apiKey != "" && secret != "" {
		cfg.Net.SASL.Enable = true
		cfg.Net.SASL.User = apiKey
		cfg.Net.SASL.Password = secret
		cfg.Net.SASL.Mechanism = sarama.SASLTypePlaintext
		cfg.Net.SASL.Handshake = true
		cfg.Net.SASL.Version = sarama.SASLHandshakeV1
		cfg.Net.TLS.Enable = true
		cfg.Net.TLS.Config = &tls.Config{
			InsecureSkipVerify: true, // ใช้เทสเท่านั้น
			ClientAuth:         tls.NoClientCert,
			MinVersion:         tls.VersionTLS12,
		}
	}

	// ค่า Fetch เดิมคุณตั้งผิดบริบทนิดหน่อย: เป็นขนาด bytes ไม่ใช่จำนวนครั้ง
	cfg.Consumer.Fetch.Default = 1 << 20 // 1MB
	cfg.Consumer.Fetch.Max = 10 << 20    // 10MB
	cfg.Consumer.Return.Errors = true

	c, err := sarama.NewConsumer(brokerUrls, cfg)
	if err != nil {
		log.Printf("Error: Failed to connect to consumer: %s", err.Error())
		return nil, errors.New("error: failed to connect to consumer")
	}
	return c, nil
}

// ---------- ConsumerGroup (recommended) ----------

type ConsumerGroupOption struct {
	Brokers              []string
	GroupID              string
	APIKey               string
	Secret               string
	EnableTLS            bool
	InsecureSkipTLSCheck bool
	Version              string
}

func NewConsumerGroup(opt ConsumerGroupOption) (sarama.ConsumerGroup, error) {
	cfg := sarama.NewConfig()
	cfg.Version = newSaramaVersion(opt.Version)
	cfg.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRange
	cfg.Consumer.Offsets.Initial = sarama.OffsetNewest
	cfg.Consumer.Return.Errors = true
	if opt.APIKey != "" && opt.Secret != "" {
		cfg.Net.SASL.Enable = true
		cfg.Net.SASL.User = opt.APIKey
		cfg.Net.SASL.Password = opt.Secret
		cfg.Net.SASL.Mechanism = sarama.SASLTypePlaintext
		cfg.Net.SASL.Handshake = true
		cfg.Net.SASL.Version = sarama.SASLHandshakeV1
	}
	if opt.EnableTLS {
		cfg.Net.TLS.Enable = true
		cfg.Net.TLS.Config = &tls.Config{
			InsecureSkipVerify: opt.InsecureSkipTLSCheck,
			ClientAuth:         tls.NoClientCert,
			MinVersion:         tls.VersionTLS12,
		}
	}
	return sarama.NewConsumerGroup(opt.Brokers, opt.GroupID, cfg)
}

// Helper run loop
func RunConsumerGroup(ctx context.Context, cg sarama.ConsumerGroup, topics []string, handler sarama.ConsumerGroupHandler) error {
	for {
		if err := cg.Consume(ctx, topics, handler); err != nil {
			return err
		}
		if ctx.Err() != nil {
			return ctx.Err()
		}
	}
}


var validate = validator.New()

// NOTE: obj ต้องเป็น "pointer to struct" เสมอ เช่น &PaymentSucceeded{}
// เดิมคุณทำ json.Unmarshal(value, &obj) ซึ่งผิด (กลายเป็น **interface{})
func DecodeMessage(obj any, value []byte) error {
	if obj == nil {
		return errors.New("error: obj must be a non-nil pointer")
	}
	if err := json.Unmarshal(value, obj); err != nil {
		log.Printf("Error: Failed to decode message: %s", err.Error())
		return errors.New("error: failed to decode message")
	}
	if err := validate.Struct(obj); err != nil {
		log.Printf("Error: Failed to validate message: %s", err.Error())
		return errors.New("error: failed to validate message")
	}
	return nil
}
