package consumer

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/twmb/franz-go/pkg/kgo"
	healthgrpc "google.golang.org/grpc/health/grpc_health_v1"

	"github.com/rodney-b/swish-test-consumer/internal/pkg/config"
	"github.com/rodney-b/swish-test-consumer/internal/pkg/healthcheck"
	"github.com/rodney-b/swish-test-consumer/internal/pkg/kafka"
	"github.com/rodney-b/swish-test-consumer/internal/pkg/logger"
	"github.com/rodney-b/swish-test-consumer/internal/pkg/telemetry"
)

func Run(cp config.ConfigProvider) error {
	log := logger.New("consumer")

	err := healthcheck.Start(cp)
	if err != nil {
		return err
	}

	ctx, ctxCancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer ctxCancel()

	tel, err := telemetry.NewTelemetry(ctx, cp, log)
	if err != nil {
		return errors.Join(err, errors.New("error initializing telemetry"))
	}
	defer tel.Shutdown()

	// Unnecessary for this app since it's not "serving" anything, but here for demonstration purposes
	healthcheck.SetAppReadinessStatus(healthgrpc.HealthCheckResponse_SERVING)

	err = consume(ctx, cp, log, tel)
	if err != nil {
		log.Error("error consuming from kafka", "error", err.Error())
		return errors.Join(errors.New("error consuming from message queue"), err)
	}

	return nil
}

func consume(ctx context.Context, cp config.ConfigProvider, log *slog.Logger, tel *telemetry.Telemetry) error {
	kafkaClient, err := kafka.NewClient(ctx, cp)
	if err != nil {
		return errors.Join(errors.New("error creating kafka client"), err)
	}
	defer kafkaClient.Close()

	if err := kafkaClient.Ping(ctx); err != nil {
		return errors.Join(errors.New("error pinging kafka client"), err)
	}

	log.Debug("message queue", "topics", cp.GetMessageQueueTopics()) // The publisher already logs all messages at info level

	for {
		fetches := kafkaClient.PollFetches(ctx)

		if err := ctx.Err(); err != nil {
			log.Info("consumer stopped - context cancelled")
			break
		}

		if errs := fetches.Errors(); len(errs) > 0 {
			for _, fErr := range errs {
				log.Error("fetch error", "topic", fErr.Topic, "partition", fErr.Partition, "error", fErr.Err)
			}
			continue
		}

		fetches.EachRecord(func(r *kgo.Record) {
			log.Info("message consumed",
				"topic", r.Topic,
				"msg", string(r.Value),
			)

			tel.IncrementMessageCounter(ctx, cp)
		})
	}

	return nil
}
