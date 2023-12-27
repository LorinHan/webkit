package cache

import (
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"time"
)

/**
* @description: 缓存代理，以aop方式加入缓存切面
	Cacheable：执行回调函数前会查询缓存，若key不存在则执行回调，将回调执行结果放入缓存，若key存在，将数据映射到参数v(应传入指针)且不执行回调函数
	Put：执行回调函数后，将回调执行结果放入缓存，与Cacheable不同的是Put不会进行前置查询，常用于更新操作
	Evict：执行回调函数后，删除该key的缓存
	Put和Evict有以ByDynamicKey为后缀的扩展，可通过回调函数的返回值来设定key
*/

type Func func() (data interface{}, err error)
type KeyFunc func() (key string, data interface{}, err error)
type KeyDelFunc func() (key string, err error)
type DelFunc func() error

func get(ctx context.Context, key string, v interface{}) (bool, error) {
	data, err := rdb.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return false, nil
		}
		return false, err
	}

	if err = json.Unmarshal([]byte(data), v); err != nil {
		return false, err
	}

	return true, nil
}

func set(ctx context.Context, key string, data, container interface{}, expiration ...time.Duration) error {
	if data == nil {
		return nil
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	var ex time.Duration
	if len(expiration) > 0 {
		ex = expiration[0]
	}
	if err = rdb.Set(ctx, key, jsonData, ex).Err(); err != nil {
		return err
	}

	if container == nil {
		return nil
	}
	return json.Unmarshal(jsonData, container)
}

// Cacheable 执行回调函数前会查询缓存，若key不存在则执行回调，将回调执行结果放入缓存，若key存在，将数据映射到参数v(应传入指针)且不执行回调函数
func Cacheable(ctx context.Context, key string, v interface{}, f Func, expiration ...time.Duration) error {
	exist, err := get(ctx, key, v)
	if err != nil {
		return err
	}
	if exist {
		return nil
	}

	res, err := f()
	if err != nil {
		return err
	}

	return set(ctx, key, res, v, expiration...)
}

// Evict 执行回调函数后，删除该key的缓存
func Evict(ctx context.Context, key string, f DelFunc) error {
	if err := f(); err != nil {
		return err
	}

	return rdb.Del(ctx, key).Err()
}

// EvictByDynamicKey 执行回调函数后，删除缓存，key取决于回调函数的返回值
func EvictByDynamicKey(ctx context.Context, f KeyDelFunc) error {
	key, err := f()
	if err != nil {
		return err
	}

	return rdb.Del(ctx, key).Err()
}

// Put 执行回调函数后，将回调执行结果放入缓存，与Cacheable不同的是Put不会进行前置查询，常用于更新操作
func Put(ctx context.Context, key string, f Func, expiration ...time.Duration) error {
	res, err := f()
	if err != nil {
		return err
	}

	return set(ctx, key, res, nil, expiration...)
}

// PutByDynamicKey 执行回调函数后，将回调执行结果放入缓存，key取决于回调函数的返回值，与Cacheable不同的是PutByDynamicKey不会进行前置查询，常用于更新操作
func PutByDynamicKey(ctx context.Context, f KeyFunc, expiration ...time.Duration) error {
	key, res, err := f()
	if err != nil {
		return err
	}

	return set(ctx, key, res, nil, expiration...)
}
