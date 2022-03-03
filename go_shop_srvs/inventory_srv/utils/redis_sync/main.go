package main

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	goredislib "github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
)

func main() {
	// Create a pool with go-redis (or redigo) which is the pool redisync will
	// use while communicating with Redis. This can also be any pool that
	// implements the `redis.Pool` interface.

	// 哪些变量可以放到global中, redis的配置是否应该放在nacos中
	client := goredislib.NewClient(&goredislib.Options{
		Addr: "172.18.0.1:6379",
	})
	pool := goredis.NewPool(client) // or, pool := redigo.NewPool(...)

	// Create an instance of redisync to be used to obtain a mutual exclusion
	// lock.
	rs := redsync.New(pool)

	// Obtain a new mutex by using the same name for all instances wanting the
	// same lock.

	var wg sync.WaitGroup

	gNum := 2
	mutexname := "421"

	wg.Add(gNum)
	for i := 0; i < gNum; i++ {
		go func(i int) {
			defer wg.Done()
			// mutex := rs.NewMutex(mutexname)
			mutex := rs.NewMutex(mutexname + strconv.Itoa(i))

			fmt.Println("开始获得锁", i)
			if err := mutex.Lock(); err != nil {
				panic(err)
			}
			fmt.Println("获得锁成功")

			time.Sleep(time.Second)

			fmt.Println("开始释放锁")

			if ok, err := mutex.Unlock(); !ok || err != nil {
				panic("unlock failed")
			}
		}(i)
	}

	wg.Wait()
}
