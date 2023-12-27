package cache

import (
	"context"
	"github.com/redis/go-redis/v9"
	"log"
	"testing"
	"time"
)

func TestCacheable(t *testing.T) {
	type User struct {
		ID    int     `json:"id"`
		Name  string  `json:"name"`
		Money float64 `json:"money"`
		OK    bool    `json:"ok"`
	}
	// 初始化
	Init(&redis.Options{Addr: "127.0.0.1:6379"})

	ctx := context.Background()
	var user User

	// 若key不存在，回调函数执行，并将结果放入缓存，下次再执行key已存在，不执行回调，而是将数据映射到user指针
	if err := Cacheable(ctx, "test", &user, func() (data interface{}, err error) {
		log.Println("数据库查询等操作，执行了...")
		return User{
			ID:   1,
			Name: "test",
		}, nil
	}, time.Minute*10); err != nil {
		log.Fatal(err)
	}

	log.Println("res", user)
}

func TestEvict(t *testing.T) {
	Init(&redis.Options{Addr: "127.0.0.1:6379"})
	ctx := context.Background()

	// 执行回调后，删除key对应的缓存
	if err := Evict(ctx, "test", func() error {
		log.Println("执行了某些操作后，删除key对应的缓存")
		return nil
	}); err != nil {
		log.Fatal(err)
	}
}

func TestEvictByDynamicKey(t *testing.T) {
	Init(&redis.Options{Addr: "127.0.0.1:6379"})
	ctx := context.Background()

	// 执行回调后，删除key对应的缓存，key取决于回调函数的返回值
	if err := EvictByDynamicKey(ctx, func() (key string, err error) {
		log.Println("回调函数执行并设置key为test，而后删除...")
		key = "test"
		return key, nil
	}); err != nil {
		log.Fatal(err)
	}
}

func TestPut(t *testing.T) {
	type User struct {
		ID    int     `json:"id"`
		Name  string  `json:"name"`
		Money float64 `json:"money"`
		OK    bool    `json:"ok"`
	}
	Init(&redis.Options{Addr: "127.0.0.1:6379"})
	ctx := context.Background()

	// 执行回调函数后，将回调执行结果放入缓存，与Cacheable不同的是Put不会进行前置查询，常用于更新操作
	if err := Put(ctx, "test", func() (data interface{}, err error) {
		log.Println("执行回调后，更新缓存...")
		return &User{ID: 100, Name: "updated by 'Put'"}, nil
	}, time.Minute*10); err != nil {
		log.Fatal(err)
	}
}

func TestPutByDynamicKey(t *testing.T) {
	type User struct {
		ID    int     `json:"id"`
		Name  string  `json:"name"`
		Money float64 `json:"money"`
		OK    bool    `json:"ok"`
	}
	Init(&redis.Options{Addr: "127.0.0.1:6379"})
	ctx := context.Background()

	// 执行回调函数后，将回调执行结果放入缓存，key取决于回调函数的返回值，与Cacheable不同的是Put不会进行前置查询，常用于更新操作
	if err := PutByDynamicKey(ctx, func() (key string, data interface{}, err error) {
		log.Println("执行回调后，更新缓存...")
		return "test", &User{ID: 100, Name: "updated by 'PutByDynamicKey'"}, nil
	}, time.Minute*10); err != nil {
		log.Fatal(err)
	}
}
