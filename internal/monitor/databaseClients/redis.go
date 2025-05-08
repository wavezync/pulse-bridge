package databaseClients

import (
	"context"
	"fmt"
	"strconv"
	"wavezync/pulse-bridge/internal/types"

	"github.com/redis/go-redis/v9"
)

func ExecRedisQuery(useConnString bool, config DatabaseClientConfig, command ...string) *types.MonitorError {
	var options *redis.Options
	var err error

	if useConnString {
		options, err = redis.ParseURL(config.ConnString)
		if err != nil {
			return types.NewConfigError(fmt.Errorf("failed to parse Redis connection string: %w", err))
		}
	} else {
		options = &redis.Options{
			Addr:           fmt.Sprintf("%s:%s", config.Host, config.Port),
			Password:       config.Password,
			MaxIdleConns:   config.MaxIdleConns,
			MaxActiveConns: config.MaxOpenConns,
		}

		dbIndex := 0
		if config.Dbname != "" {
			if parsedIndex, err := strconv.Atoi(config.Dbname); err == nil {
				dbIndex = parsedIndex
			}
		}
		options.DB = dbIndex
	}

	redisClient := redis.NewClient(options)
	defer redisClient.Close()

	ctx, cancel := context.WithTimeout(context.Background(), config.Timeout)
	defer cancel()

	_, err = redisClient.Ping(ctx).Result()
	if err != nil {
		return types.NewClientError(fmt.Errorf("redis connection test failed: %w", err))
	}

	if len(command) > 0 && command[0] != "" {
		cmd := redisClient.Do(ctx, command)
		if cmd.Err() != nil {
			return types.NewClientError(fmt.Errorf("redis command execution failed: %w", cmd.Err()))
		}
	}

	return nil
}
