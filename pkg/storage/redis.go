package storage

import (
	"time"

	"github.com/adinandradrs/omni-service-sdk/pkg/domain"
	"github.com/adinandradrs/omni-service-sdk/pkg/service"
	"github.com/go-redis/redis"
	"go.uber.org/zap"
)

type (
	RedisOptions struct {
		Addr   string
		Addrs  []string
		Passwd string
		Index  int
		Pool   int
		Idle   int
		Logger *zap.Logger
	}

	singleRedis struct {
		logger *zap.Logger
		cache  *redis.Client
	}

	clusterRedis struct {
		logger *zap.Logger
		cache  *redis.ClusterClient
	}
)

type Cacher interface {
	Set(k string, p string, v interface{}, d time.Duration) *domain.TechnicalError
	Delete(k string, p string) *domain.TechnicalError
	Get(k string, p string) (v string, e *domain.TechnicalError)
	Ttl(k string, p string) (t time.Duration, e *domain.TechnicalError)
}

func NewRedis(o *RedisOptions) Cacher {
	return &singleRedis{
		logger: o.Logger,
		cache: redis.NewClient(&redis.Options{
			Addr:         o.Addr,
			Password:     o.Passwd,
			DB:           o.Index,
			PoolSize:     o.Pool,
			MinIdleConns: o.Idle,
		}),
	}
}

func NewClusterRedis(o *RedisOptions) Cacher {
	return &clusterRedis{
		logger: o.Logger,
		cache: redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:        o.Addrs,
			Password:     o.Passwd,
			PoolSize:     o.Pool,
			MinIdleConns: o.Idle,
		}),
	}
}

func (r clusterRedis) Set(k string, p string, v interface{}, d time.Duration) *domain.TechnicalError {
	r.cache.Del(k + ":" + p)
	if d != 0*time.Second {
		_, err := r.cache.SetNX(k+":"+p, v, d).Result()
		if err != nil {
			r.logger.Error("failed on cluster setnx ops property ", zap.String("key", k), zap.String("pair", p))
			return service.Exception("failed on cluster setnx ops", err, r.logger)
		}
	} else {
		_, err := r.cache.Set(k+":"+p, v, 0).Result()
		if err != nil {
			r.logger.Error("failed on cluster set ops property", zap.String("key", k), zap.String("pair", p))
			return service.Exception("failed on cluster set ops", err, r.logger)
		}
	}
	return nil
}

func (r clusterRedis) Delete(k string, p string) (out *domain.TechnicalError) {
	if cmd := r.cache.Del(k + ":" + p); cmd.Err() != nil {
		r.logger.Error("failed on cluster delete ops property", zap.String("key", k), zap.String("pair", p))
		return service.Exception("failed on cluster delete ops", cmd.Err(), r.logger)
	}
	return nil
}

func (r clusterRedis) Get(k string, p string) (v string, e *domain.TechnicalError) {
	v, err := r.cache.Get(k + ":" + p).Result()
	if err != nil {
		r.logger.Error("failed on cluster get ops property", zap.String("key", k), zap.String("pair", p))
		return v, service.Exception("failed on cluster get ops", err, r.logger)
	}
	return v, nil
}

func (r clusterRedis) Ttl(k string, p string) (t time.Duration, e *domain.TechnicalError) {
	if cmd := r.cache.TTL(k + ":" + p); cmd.Err() != nil {
		r.logger.Error("failed on cluster TTL ops property", zap.String("key", k), zap.String("pair", p))
		return t, service.Exception("failed on cluster TTL ops", cmd.Err(), r.logger)
	} else {
		return cmd.Val(), nil
	}
}

func (r singleRedis) Set(k string, p string, v interface{}, d time.Duration) *domain.TechnicalError {
	r.cache.Del(k + ":" + p)
	if d != 0*time.Second {
		_, err := r.cache.SetNX(k+":"+p, v, d).Result()
		if err != nil {
			r.logger.Error("failed on setnx ops property", zap.String("key", k), zap.String("pair", p))
			return service.Exception("failed on setnx ops", err, r.logger)
		}
	} else {
		_, err := r.cache.Set(k+":"+p, v, 0).Result()
		if err != nil {
			r.logger.Error("failed on set ops property", zap.String("key", k), zap.String("pair", p))
			return service.Exception("failed on set ops", err, r.logger)
		}
	}
	return nil
}

func (r singleRedis) Delete(k string, p string) *domain.TechnicalError {
	if cmd := r.cache.Del(k + ":" + p); cmd.Err() != nil {
		r.logger.Error("failed on delete ops property", zap.String("key", k), zap.String("pair", p))
		return service.Exception("failed on delete ops", cmd.Err(), r.logger)
	}
	return nil
}

func (r singleRedis) Get(k string, p string) (string, *domain.TechnicalError) {
	v, err := r.cache.Get(k + ":" + p).Result()
	if err != nil {
		r.logger.Error("failed on get ops property", zap.String("key", k), zap.String("pair", p))
		return v, service.Exception("failed on get ops", err, r.logger)
	}
	return v, nil
}

func (r singleRedis) Ttl(k string, p string) (t time.Duration, e *domain.TechnicalError) {
	cmd := r.cache.TTL(k + ":" + p)
	if cmd.Err() != nil {
		r.logger.Error("failed on TTL ops", zap.Error(cmd.Err()))
		return t, service.Exception("failed on TTL ops", cmd.Err(), r.logger)
	} else {
		return cmd.Val(), nil
	}
}
