package storage

import (
	"time"

	"github.com/adinandradrs/omni-service-sdk/pkg/domain"
	"github.com/adinandradrs/omni-service-sdk/pkg/service"
	"github.com/go-redis/redis"
	"go.uber.org/zap"
)

type RedisConn struct {
	Addr   string
	Addrs  []string
	Passwd string
	Index  int
	Pool   int
	Idle   int
}

func NewRedisClient(rc RedisConn) *redis.Client {
	r := redis.NewClient(&redis.Options{
		Addr:         rc.Addr,
		Password:     rc.Passwd,
		DB:           rc.Index,
		PoolSize:     rc.Pool,
		MinIdleConns: rc.Idle,
	})
	return r
}

func NewClusterRedisClient(rc RedisConn) *redis.ClusterClient {
	r := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:        rc.Addrs,
		Password:     rc.Passwd,
		PoolSize:     rc.Pool,
		MinIdleConns: rc.Idle,
	})
	return r
}

type Cacher interface {
	Set(k string, p string, v interface{}, d time.Duration) *domain.TechnicalError
	Delete(k string, p string) *domain.TechnicalError
	Get(k string, p string) (v string, e *domain.TechnicalError)
	Ttl(k string, p string) (t time.Duration, e *domain.TechnicalError)
}

type RedisCapsule struct {
	Logger *zap.Logger
	Cache  *redis.Client
}

type ClusterRedisCapsule struct {
	Logger *zap.Logger
	Cache  *redis.ClusterClient
}

func NewRedisOps(rc RedisCapsule) Cacher {
	return &rc
}

func NewClusterRedisOps(rc ClusterRedisCapsule) Cacher {
	return &rc
}

func (r ClusterRedisCapsule) Set(k string, p string, v interface{}, d time.Duration) *domain.TechnicalError {
	r.Cache.Del(k + ":" + p)
	if d != 0*time.Second {
		_, err := r.Cache.SetNX(k+":"+p, v, d).Result()
		if err != nil {
			r.Logger.Error("failed on cluster setnx ops property ", zap.String("key", k), zap.String("pair", p))
			return service.Exception("failed on cluster setnx ops", err, r.Logger)
		}
	} else {
		_, err := r.Cache.Set(k+":"+p, v, 0).Result()
		if err != nil {
			r.Logger.Error("failed on cluster set ops property", zap.String("key", k), zap.String("pair", p))
			return service.Exception("failed on cluster set ops", err, r.Logger)
		}
	}
	return nil
}

func (r ClusterRedisCapsule) Delete(k string, p string) (out *domain.TechnicalError) {
	if cmd := r.Cache.Del(k + ":" + p); cmd.Err() != nil {
		r.Logger.Error("failed on cluster delete ops property", zap.String("key", k), zap.String("pair", p))
		return service.Exception("failed on cluster delete ops", cmd.Err(), r.Logger)
	}
	return nil
}

func (r ClusterRedisCapsule) Get(k string, p string) (v string, e *domain.TechnicalError) {
	v, err := r.Cache.Get(k + ":" + p).Result()
	if err != nil {
		r.Logger.Error("failed on cluster get ops property", zap.String("key", k), zap.String("pair", p))
		return v, service.Exception("failed on cluster get ops", err, r.Logger)
	}
	return v, nil
}

func (r ClusterRedisCapsule) Ttl(k string, p string) (t time.Duration, e *domain.TechnicalError) {
	if cmd := r.Cache.TTL(k + ":" + p); cmd.Err() != nil {
		r.Logger.Error("failed on cluster TTL ops property", zap.String("key", k), zap.String("pair", p))
		return t, service.Exception("failed on cluster TTL ops", cmd.Err(), r.Logger)
	} else {
		return cmd.Val(), nil
	}
}

func (r RedisCapsule) Set(k string, p string, v interface{}, d time.Duration) *domain.TechnicalError {
	r.Cache.Del(k + ":" + p)
	if d != 0*time.Second {
		_, err := r.Cache.SetNX(k+":"+p, v, d).Result()
		if err != nil {
			r.Logger.Error("failed on setnx ops property", zap.String("key", k), zap.String("pair", p))
			return service.Exception("failed on setnx ops", err, r.Logger)
		}
	} else {
		_, err := r.Cache.Set(k+":"+p, v, 0).Result()
		if err != nil {
			r.Logger.Error("failed on set ops property", zap.String("key", k), zap.String("pair", p))
			return service.Exception("failed on set ops", err, r.Logger)
		}
	}
	return nil
}

func (r RedisCapsule) Delete(k string, p string) *domain.TechnicalError {
	if cmd := r.Cache.Del(k + ":" + p); cmd.Err() != nil {
		r.Logger.Error("failed on delete ops property", zap.String("key", k), zap.String("pair", p))
		return service.Exception("failed on delete ops", cmd.Err(), r.Logger)
	}
	return nil
}

func (r RedisCapsule) Get(k string, p string) (string, *domain.TechnicalError) {
	v, err := r.Cache.Get(k + ":" + p).Result()
	if err != nil {
		r.Logger.Error("failed on get ops property", zap.String("key", k), zap.String("pair", p))
		return v, service.Exception("failed on get ops", err, r.Logger)
	}
	return v, nil
}

func (r RedisCapsule) Ttl(k string, p string) (t time.Duration, e *domain.TechnicalError) {
	cmd := r.Cache.TTL(k + ":" + p)
	if cmd.Err() != nil {
		r.Logger.Error("failed on TTL ops", zap.Error(cmd.Err()))
		return t, service.Exception("failed on TTL ops", cmd.Err(), r.Logger)
	} else {
		return cmd.Val(), nil
	}
}
