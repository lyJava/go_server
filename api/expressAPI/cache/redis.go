package cache

import (
	"apiProject/api/expressAPI/config"
	cfg "apiProject/api/expressAPI/types/config"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/mitchellh/mapstructure"
	"log"
	"reflect"
	"time"
)

var RedisCache = &redis.Client{}

func init() {

	viperConfig := config.ReadConfig("api/expressApi/config", "application", "yml")
	//redisMap := viperConfig.Get("redis").(map[string]string)
	var redisConfig cfg.RedisConfig
	/*err := mapstructure.Decode(redisMap, &redisConfig)
	if err != nil {
		fmt.Println("failed to decode Redis config:", err)
		return
	}*/
	viperConfig.Unmarshal(&redisConfig)

	if (redisConfig == cfg.RedisConfig{}) {
		fmt.Println("no Redis configuration found, skipping initialization")
		return
	}

	/*if err := DecodeConfig(redisMap, &redisConfig); err != nil {
		fmt.Println("failed to decode Redis config:", err)
		return
	}

	if len(redisMap) == 0 || redisConfig.Address == "" {
		fmt.Println("no Redis configuration found, skipping initialization")
		return
	}*/
	redisConfigItem := redisConfig.Redis
	if redisConfigItem.Address == "" {
		fmt.Println("no Redis configuration found, skipping initialization")
		return
	}

	client, err := initRedisClient(redisConfigItem)
	if err != nil {
		log.Printf("failed to initialize Redis client:%+v", err)
		return
	}
	RedisCache = client
	log.Println("Redis initialized successfully")
}

func initRedisClient(config cfg.RedisConfigItem) (*redis.Client, error) {
	// net.JoinHostPort(config.Host, config.Port),
	client := redis.NewClient(&redis.Options{
		Addr:     config.Address,
		Password: config.Password,
		DB:       config.Db,
	})

	_, err := client.Ping().Result()
	if err != nil {
		log.Printf("failed to ping Redis error:%+v", err)
		return nil, fmt.Errorf("failed to ping Redis server: %w", err)
	}

	return client, nil
}

// todo 这里可以避免出现
// * 'Password' expected type 'string', got unconvertible type 'int', value: '244112311'
// * 'Port' expected type 'string', got unconvertible type 'int', value: '6379'
// todo 或者在redis配置文件中密码与端口号加上双引号，这样viper解析的时候不会将他们解析为整数了
// stringFromIntHook 定义一个自定义的 DecodeHook 函数
func stringFromIntHook(f reflect.Kind, t reflect.Kind, data interface{}) (interface{}, error) {
	if f == reflect.Int && t == reflect.String {
		return fmt.Sprintf("%d", data), nil
	}
	return data, nil
}

// DecodeConfig 在DecodeConfig函数中使用 DecodeHook
func DecodeConfig(redisMap map[string]interface{}, redisConfig *cfg.RedisConfig) error {
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		DecodeHook:       mapstructure.ComposeDecodeHookFunc(stringFromIntHook),
		WeaklyTypedInput: true,
		Result:           redisConfig,
	})
	if err != nil {
		return err
	}
	if err := decoder.Decode(redisMap); err != nil {
		return err
	}
	return nil
}

// Get 获取缓存数据
func Get(key string) (string, error) {
	result, err := RedisCache.Get(key).Result()
	return result, err
}

// withRedisCache 使用闭包方式
func withRedisCache(fn func() error) {
	if RedisCache == nil || reflect.DeepEqual(RedisCache, &redis.Client{}) {
		log.Println("RedisCache 未初始化，请检查 Redis 初始化过程")
		return
	}
	err := fn()
	if err != nil {
		log.Printf("操作失败: %s", err.Error())
	}
}

func checkRedisCache() bool {
	if RedisCache == nil || reflect.DeepEqual(RedisCache, &redis.Client{}) {
		// RedisCache 未初始化，可能发生了错误
		// 进行适当的错误处理
		log.Println("RedisCache 未初始化，请检查 Redis 初始化过程")
		return false
	}
	return true
}

// Set 设置数据到缓存
// 参数
//
//		key 存储的键
//		value 存储的值
//	 timeout 缓存时间(秒)
func Set(key string, value interface{}, timeout int64) {
	withRedisCache(func() error {
		// 将 timeout 转换为持续时间
		duration := time.Duration(timeout) * time.Second
		return RedisCache.Set(key, value, duration).Err()
	})
}

// LPush RPush 使用RPush命令往队列右边加入
func LPush(key string, value ...interface{}) error {
	if checkRedisCache() {
		err := RedisCache.LPush(key, value).Err()
		return err
	}
	return nil
}

// RPop LPop 取出并移除左边第一个元素
func RPop(key string) (interface{}, error) {
	var result interface{}
	var err error
	withRedisCache(func() error {
		result, err = RedisCache.RPop(key).Result()
		return err
	})
	return result, err
}

// BRPop BLPop 取出并移除左边第一个元素， 如果列表没有元素会阻塞列表直到等待超时或发现可弹出元素为止。
func BRPop(timeout time.Duration, key string) (interface{}, error) {
	if checkRedisCache() {
		result, err := RedisCache.BRPop(timeout, key).Result()
		return result, err
	}
	return nil, nil
}

// LLen 获取数据长度
func LLen(key string) (int64, error) {
	result, err := RedisCache.LLen(key).Result()
	return result, err
}

// LRange 获取数据列表
func LRange(key string, start, end int64) ([]string, error) {
	result, err := RedisCache.LRange(key, start, end).Result()
	return result, err
}

// HSet hash相关操作
// set hash 适合存储结构
func HSet(hashKey, key string, value interface{}) error {
	err := RedisCache.HSet(hashKey, key, value).Err()
	return err
}

// HGet get Hash
func HGet(hashKey, key string) (interface{}, error) {
	result, err := RedisCache.HGet(hashKey, key).Result()
	return result, err
}

// HGetAll 获取所以hash ,返回map
func HGetAll(hashKey string) (map[string]string, error) {
	result, err := RedisCache.HGetAll(hashKey).Result()
	return result, err
}

// HDel 删除一个或多个哈希表字段
func HDel(hashKey string, key ...string) error {
	err := RedisCache.HDel(hashKey, key...).Err()
	return err
}

// HExists 查看哈希表的指定字段是否存在
func HExists(hashKey, key string) (bool, error) {
	result, err := RedisCache.HExists(hashKey, key).Result()
	return result, err
}

// SAdd -----------------Set------------------------
// 添加Set
func SAdd(key string, values ...interface{}) error {
	err := RedisCache.SAdd(key, values).Err()
	return err
}

// SCard 获取集合的成员数
func SCard(key string) (int64, error) {
	result, err := RedisCache.SCard(key).Result()
	return result, err
}

// SMembers 获取集合的所有成员
func SMembers(key string) ([]string, error) {
	result, err := RedisCache.SMembers(key).Result()
	return result, err
}

// SRem 移除集合里的某个元素
func SRem(key string, value interface{}) error {
	err := RedisCache.SRem(key, value).Err()
	return err
}

// SPop 移除并返回set的一个随机元素(SET是无序的)
func SPop(key string) (interface{}, error) {
	result, err := RedisCache.SPop(key).Result()
	return result, err
}

// ZAdd ------------------ZSet-------------------------
func ZAdd(key string, values []redis.Z) error {
	err := RedisCache.ZAdd(key, values...).Err()
	return err
}

// ZIncrBy 给指定的key和值加分
func ZIncrBy(key string, score float64, value string) error {
	err := RedisCache.ZIncrBy(key, score, value).Err()
	return err
}

// ZRevRangeWithScores 取zSet里的前n名热度的数据
func ZRevRangeWithScores(key string, start, end int64) ([]redis.Z, error) {
	result, err := RedisCache.ZRevRangeWithScores(key, start, end).Result()
	return result, err
}

// Expire 给指定key 设置过期时间
func Expire(key string, duration time.Duration) error {
	err := RedisCache.Expire(key, duration).Err()
	return err
}

// ExpireAt 给指定Key 设置过期时间，时间格式为UNIX时间
func ExpireAt(key string, duration time.Time) error {
	err := RedisCache.ExpireAt(key, duration).Err()
	return err
}

// TTL 获取key的生存时间
func TTL(key string) (time.Duration, error) {
	result, err := RedisCache.TTL(key).Result()
	return result, err
}
