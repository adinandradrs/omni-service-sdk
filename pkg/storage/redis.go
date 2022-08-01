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

	Redis struct {
		logger *zap.Logger
		Cache  *redis.Client
	}

	ClusterRedis struct {
		logger *zap.Logger
		Cache  *redis.ClusterClient
	}
)

type Cacher interface {
	Set(k string, p string, v interface{}, d time.Duration) *domain.TechnicalError
	Delete(k string, p string) *domain.TechnicalError
	Get(k string, p string) (v string, e *domain.TechnicalError)
	Ttl(k string, p string) (t time.Duration, e *domain.TechnicalError)
}

func NewRedis(o *RedisOptions) Cacher {
	return &Redis{
		logger: o.Logger,
		Cache: redis.NewClient(&redis.Options{
			Addr:         o.Addr,
			Password:     o.Passwd,
			DB:           o.Index,
			PoolSize:     o.Pool,
			MinIdleConns: o.Idle,
		}),
	}
}

func NewClusterRedis(o *RedisOptions) Cacher {
	return &ClusterRedis{
		logger: o.Logger,
		Cache: redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:        o.Addrs,
			Password:     o.Passwd,
			PoolSize:     o.Pool,
			MinIdleConns: o.Idle,
		}),
	}
}

func (r ClusterRedis) Set(k string, p string, v interface{}, d time.Duration) *domain.TechnicalError {
	r.Cache.Del(k + ":" + p)
	if d != 0*time.Second {
		_, err := r.Cache.SetNX(k+":"+p, v, d).Result()
		if err != nil {
			r.logger.Error("failed on cluster setnx ops property ", zap.String("key", k), zap.String("pair", p))
			return service.Exception("failed on cluster setnx ops", err, r.logger)
		}
	} else {
		_, err := r.Cache.Set(k+":"+p, v, 0).Result()
		if err != nil {
			r.logger.Error("failed on cluster set ops property", zap.String("key", k), zap.String("pair", p))
			return service.Exception("failed on cluster set ops", err, r.logger)
		}
	}
	return nil
}

func (r ClusterRedis) Delete(k string, p string) (out *domain.TechnicalError) {
	if cmd := r.Cache.Del(k + ":" + p); cmd.Err() != nil {
		r.logger.Error("failed on cluster delete ops property", zap.String("key", k), zap.String("pair", p))
		return service.Exception("failed on cluster delete ops", cmd.Err(), r.logger)
	}
	return nil
}

func (r ClusterRedis) Get(k string, p string) (v string, e *domain.TechnicalError) {
	v, err := r.Cache.Get(k + ":" + p).Result()
	if err != nil {
		r.logger.Error("failed on cluster get ops property", zap.String("key", k), zap.String("pair", p))
		return v, service.Exception("failed on cluster get ops", err, r.logger)
	}
	return v, nil
}

func (r ClusterRedis) Ttl(k string, p string) (t time.Duration, e *domain.TechnicalError) {
	if cmd := r.Cache.TTL(k + ":" + p); cmd.Err() != nil {
		r.logger.Error("failed on cluster TTL ops property", zap.String("key", k), zap.String("pair", p))
		return t, service.Exception("failed on cluster TTL ops", cmd.Err(), r.logger)
	} else {
		return cmd.Val(), nil
	}
}

func (r Redis) Set(k string, p string, v interface{}, d time.Duration) *domain.TechnicalError {
	r.Cache.Del(k + ":" + p)
	if d != 0*time.Second {
		_, err := r.Cache.SetNX(k+":"+p, v, d).Result()
		if err != nil {
			r.logger.Error("failed on setnx ops property", zap.String("key", k), zap.String("pair", p))
			return service.Exception("failed on setnx ops", err, r.logger)
		}
	} else {
		_, err := r.Cache.Set(k+":"+p, v, 0).Result()
		if err != nil {
			r.logger.Error("failed on set ops property", zap.String("key", k), zap.String("pair", p))
			return service.Exception("failed on set ops", err, r.logger)
		}
	}
	return nil
}

func (r Redis) Delete(k string, p string) *domain.TechnicalError {
	if cmd := r.Cache.Del(k + ":" + p); cmd.Err() != nil {
		r.logger.Error("failed on delete ops property", zap.String("key", k), zap.String("pair", p))
		return service.Exception("failed on delete ops", cmd.Err(), r.logger)
	}
	return nil
}

func (r Redis) Get(k string, p string) (string, *domain.TechnicalError) {
	v, err := r.Cache.Get(k + ":" + p).Result()
	if err != nil {
		r.logger.Error("failed on get ops property", zap.String("key", k), zap.String("pair", p))
		return v, service.Exception("failed on get ops", err, r.logger)
	}
	return v, nil
}

func (r Redis) Ttl(k string, p string) (t time.Duration, e *domain.TechnicalError) {
	cmd := r.Cache.TTL(k + ":" + p)
	if cmd.Err() != nil {
		r.logger.Error("failed on TTL ops", zap.Error(cmd.Err()))
		return t, service.Exception("failed on TTL ops", cmd.Err(), r.logger)
	} else {
		return cmd.Val(), nil
	}
}
