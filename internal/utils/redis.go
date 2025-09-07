package utils

import (
	"context"
	"fmt"
	"time"

	"air-quality-server/internal/config"

	"github.com/go-redis/redis/v8"
)

// Redis Redis连接管理器
type Redis struct {
	Client *redis.Client
	config *config.RedisConfig
	logger Logger
}

// NewRedis 创建新的Redis连接
func NewRedis(cfg *config.RedisConfig, log Logger) (*Redis, error) {
	// 创建Redis客户端
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
		PoolSize: cfg.PoolSize,
	})

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("Redis连接测试失败: %w", err)
	}

	redisClient := &Redis{
		Client: client,
		config: cfg,
		logger: log,
	}

	if log != nil {
		log.Info("Redis连接成功",
			String("host", cfg.Host),
			Int("port", cfg.Port),
			Int("db", cfg.DB),
		)
	}

	return redisClient, nil
}

// Close 关闭Redis连接
func (r *Redis) Close() error {
	if r.Client != nil {
		return r.Client.Close()
	}
	return nil
}

// Ping 测试Redis连接
func (r *Redis) Ping(ctx context.Context) error {
	if r.Client != nil {
		return r.Client.Ping(ctx).Err()
	}
	return fmt.Errorf("Redis连接未初始化")
}

// GetStats 获取Redis连接统计信息
func (r *Redis) GetStats() *redis.PoolStats {
	if r.Client != nil {
		return r.Client.PoolStats()
	}
	return nil
}

// PubSub Redis发布订阅管理器
type PubSub struct {
	redis *Redis
}

// NewPubSub 创建发布订阅管理器
func NewPubSub(redis *Redis) *PubSub {
	return &PubSub{redis: redis}
}

// Publish 发布消息
func (ps *PubSub) Publish(ctx context.Context, channel string, message interface{}) error {
	return ps.redis.Client.Publish(ctx, channel, message).Err()
}

// Subscribe 订阅频道
func (ps *PubSub) Subscribe(ctx context.Context, channels ...string) *redis.PubSub {
	return ps.redis.Client.Subscribe(ctx, channels...)
}

// PSubscribe 模式订阅
func (ps *PubSub) PSubscribe(ctx context.Context, patterns ...string) *redis.PubSub {
	return ps.redis.Client.PSubscribe(ctx, patterns...)
}

// Cache Redis缓存管理器
type Cache struct {
	redis *Redis
}

// NewCache 创建缓存管理器
func NewCache(redis *Redis) *Cache {
	return &Cache{redis: redis}
}

// Set 设置缓存
func (c *Cache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return c.redis.Client.Set(ctx, key, value, expiration).Err()
}

// Get 获取缓存
func (c *Cache) Get(ctx context.Context, key string) (string, error) {
	return c.redis.Client.Get(ctx, key).Result()
}

// GetInt 获取整数缓存
func (c *Cache) GetInt(ctx context.Context, key string) (int64, error) {
	return c.redis.Client.Get(ctx, key).Int64()
}

// GetFloat 获取浮点数缓存
func (c *Cache) GetFloat(ctx context.Context, key string) (float64, error) {
	return c.redis.Client.Get(ctx, key).Float64()
}

// Del 删除缓存
func (c *Cache) Del(ctx context.Context, keys ...string) error {
	return c.redis.Client.Del(ctx, keys...).Err()
}

// Exists 检查键是否存在
func (c *Cache) Exists(ctx context.Context, keys ...string) (int64, error) {
	return c.redis.Client.Exists(ctx, keys...).Result()
}

// Expire 设置过期时间
func (c *Cache) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return c.redis.Client.Expire(ctx, key, expiration).Err()
}

// TTL 获取剩余过期时间
func (c *Cache) TTL(ctx context.Context, key string) (time.Duration, error) {
	return c.redis.Client.TTL(ctx, key).Result()
}

// Incr 自增
func (c *Cache) Incr(ctx context.Context, key string) (int64, error) {
	return c.redis.Client.Incr(ctx, key).Result()
}

// Decr 自减
func (c *Cache) Decr(ctx context.Context, key string) (int64, error) {
	return c.redis.Client.Decr(ctx, key).Result()
}

// HSet 设置哈希字段
func (c *Cache) HSet(ctx context.Context, key string, values ...interface{}) error {
	return c.redis.Client.HSet(ctx, key, values...).Err()
}

// HGet 获取哈希字段
func (c *Cache) HGet(ctx context.Context, key, field string) (string, error) {
	return c.redis.Client.HGet(ctx, key, field).Result()
}

// HGetAll 获取所有哈希字段
func (c *Cache) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	return c.redis.Client.HGetAll(ctx, key).Result()
}

// HDel 删除哈希字段
func (c *Cache) HDel(ctx context.Context, key string, fields ...string) error {
	return c.redis.Client.HDel(ctx, key, fields...).Err()
}

// SAdd 添加到集合
func (c *Cache) SAdd(ctx context.Context, key string, members ...interface{}) error {
	return c.redis.Client.SAdd(ctx, key, members...).Err()
}

// SMembers 获取集合成员
func (c *Cache) SMembers(ctx context.Context, key string) ([]string, error) {
	return c.redis.Client.SMembers(ctx, key).Result()
}

// SRem 从集合中移除
func (c *Cache) SRem(ctx context.Context, key string, members ...interface{}) error {
	return c.redis.Client.SRem(ctx, key, members...).Err()
}

// ZAdd 添加到有序集合
func (c *Cache) ZAdd(ctx context.Context, key string, members ...*redis.Z) error {
	return c.redis.Client.ZAdd(ctx, key, members...).Err()
}

// ZRange 获取有序集合范围
func (c *Cache) ZRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return c.redis.Client.ZRange(ctx, key, start, stop).Result()
}

// ZRem 从有序集合中移除
func (c *Cache) ZRem(ctx context.Context, key string, members ...interface{}) error {
	return c.redis.Client.ZRem(ctx, key, members...).Err()
}
