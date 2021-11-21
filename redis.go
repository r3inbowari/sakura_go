package Sakura

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/go-redis/redis/v8"
	"github.com/robfig/cron"
	"golang.org/x/net/context"
	"strconv"
	"time"
)

type Snapshot struct {
	rdb *redis.Client
	ctx context.Context
}

func InitCacheService() *Snapshot {
	var sp = Snapshot{}
	sp.ctx = context.Background()
	sp.rdb = redis.NewClient(&redis.Options{
		Addr:     GetConfig(false).RedisURL,
		Password: GetConfig(false).RedisPass,
		DB:       0,
	})
	Log.Info("[Cache] Redis pre-connected -> " + GetConfig(false).RedisURL)
	return &sp
}

func (s *Snapshot) Set(key string, value interface{}) error {
	err := s.rdb.Set(s.ctx, key, value, 0).Err()
	return err
}

func (s *Snapshot) SetEx(key string, value interface{}, ex time.Duration) error {
	cacheStart = time.Now()
	err := s.rdb.Set(s.ctx, key, value, ex).Err()
	spend := time.Now().UnixNano() - cacheStart.UnixNano()
	Log.Info("[Cache] play link snapshot | " + strconv.Itoa(int(spend)/1e6) + "ms")
	return err
}

func (s *Snapshot) Get(key string) (string, error) {
	val2, err := s.rdb.Get(s.ctx, key).Result()
	if err == redis.Nil {
		return "", err
	} else if err != nil {
		return "", err
	} else {
		return val2, nil
	}
}

func (s *Snapshot) UseCache() {
	Log.Info("[Cache] Cache service OPEN")
	DMSakuraHomepageSnapshotCron()
	c := cron.New()
	_ = c.AddFunc("0 */10 * * * ?", DMSakuraHomepageSnapshotCron)
	c.Start()
}

func HomepageSnapshotCron() {
	cacheStart = time.Now()
	var err error
	LatestHome, err = goquery.NewDocument("http://www.yhdm.tv")
	if err != nil {
		Log.Warn("[Cache] cache pull failed")
	}
	spend := time.Now().UnixNano() - cacheStart.UnixNano()
	Log.Info("[Cache] Homepage snapshot | " + strconv.Itoa(int(spend)/1e6) + "ms")
}
