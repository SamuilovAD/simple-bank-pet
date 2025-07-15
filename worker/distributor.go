package worker

import (
	"context"
	"github.com/hibiken/asynq"
)

type TaskDistributorInterface interface {
	Distribute(
		ctx context.Context,
		payload PayloadSendVerifyEmail,
		opts ...asynq.Option,
	) error
}

type RedisTaskDistributor struct {
	client *asynq.Client
}

func NewRedisTaskDistributor(redisOpt asynq.RedisClientOpt) TaskDistributorInterface {
	client := asynq.NewClient(redisOpt)
	return &RedisTaskDistributor{client}
}
