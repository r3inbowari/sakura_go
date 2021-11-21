package Sakura

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/sirupsen/logrus"
	"os"
	"regexp"
	"time"
)

// LocalConfig 配置
type LocalConfig struct {
	Finger      string    `json:"finger"`      // canvas指纹
	APIAddr     string    `json:"api_addr"`    // API服务ADDR
	CacheTime   time.Time `json:"-"`           // 缓存时间
	LoggerLevel *string   `json:"log_level"`   // 日志等级
	CheckLink   string    `json:"check_link"`  // 检查更新地址
	AutoUpdate  bool      `json:"auto_update"` // 自动更新
	RedisURL    string    `json:"redis_url"`
	RedisPass   string    `json:"redis_pass"`
}

func InitConfig() {
	if !Exists("bili.json") {
		Log.Info("[FILE] Init user configuration")
		var config LocalConfig
		var l = "debug"
		config.Finger = "532ca3dd-2104-4799-8fd2-9f4b16e5cdc9"
		config.LoggerLevel = &l
		config.APIAddr = ":9090"
		config.AutoUpdate = true
		config.RedisPass = "123456789"
		config.RedisURL = "123.org"
		_ = config.SetConfig()
	}
}

var config = new(LocalConfig)
var configPath = "bili.json"

// GetConfig 返回配置文件
// imm 立即返回
func GetConfig(imm bool) *LocalConfig {
	if config.CacheTime.Before(time.Now()) || imm {
		if err := LoadConfig(configPath, config); err != nil {
			Log.Error("loading file failed")
			time.Sleep(time.Second * 5)
			os.Exit(76)
			return nil
		}
		config.CacheTime = time.Now().Add(time.Second * 60)
	}
	return config
}

// SetConfig 更新
func (lc *LocalConfig) SetConfig() error {
	fp, err := os.Create(configPath)
	if err != nil {
		Log.WithFields(logrus.Fields{"err": err}).Error("loading file failed")
	}
	defer fp.Close()
	data, err := json.Marshal(lc)
	if err != nil {
		Log.WithFields(logrus.Fields{"err": err}).Error("marshal file failed")
	}
	n, err := fp.Write(data)
	if err != nil {
		Log.WithFields(logrus.Fields{"err": err}).Error("write file failed")
	}
	Log.WithFields(logrus.Fields{"size": n}).Info("[FILE] Update user configuration")
	return nil
}

const configFileSizeLimit = 10 << 20

// LoadConfig path 文件路径 dist 存放目标
func LoadConfig(path string, dist interface{}) error {
	configFile, err := os.Open(path)
	if err != nil {
		Log.WithFields(logrus.Fields{"path": path, "err": err}).Error("Failed to open config file.")
		return err
	}

	fi, _ := configFile.Stat()
	if size := fi.Size(); size > (configFileSizeLimit) {
		Log.WithFields(logrus.Fields{"path": path, "size": size}).Error("Config file size exceeds reasonable limited")
		return errors.New("limited")
	}

	if fi.Size() == 0 {
		Log.WithFields(logrus.Fields{"path": path, "size": 0}).Error("Config file is empty, skipping")
		return errors.New("empty")
	}

	buffer := make([]byte, fi.Size())
	_, err = configFile.Read(buffer)
	buffer, err = StripComments(buffer)
	if err != nil {
		Log.WithFields(logrus.Fields{"err": err}).Error("Failed to strip comments from json")
		return err
	}

	buffer = []byte(os.ExpandEnv(string(buffer)))

	err = json.Unmarshal(buffer, &dist)
	if err != nil {
		Log.WithFields(logrus.Fields{"err": err}).Error("Failed unmarshalling json")
		return err
	}
	return nil
}

// StripComments 注释清除
func StripComments(data []byte) ([]byte, error) {
	data = bytes.Replace(data, []byte("\r"), []byte(""), 0)
	lines := bytes.Split(data, []byte("\n"))
	filtered := make([][]byte, 0)

	for _, line := range lines {
		match, err := regexp.Match(`^\s*#`, line)
		if err != nil {
			return nil, err
		}
		if !match {
			filtered = append(filtered, line)
		}
	}
	return bytes.Join(filtered, []byte("\n")), nil
}
