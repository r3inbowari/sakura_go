package Sakura

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"regexp"
	"strconv"
	"time"
)

/**
 * local config struct
 */
type LocalConfig struct {
	Finger      string    `json:"finger"`     // canvas指纹
	MysqlArgs   string    `json:"mysql_args"` // canvas指纹
	APIAddr     string    `json:"api_addr"`   // API服务ADDR
	CacheTime   time.Time `json:"-"`          // 缓存时间
	LoggerLevel *string   `json:"log_level"`  // 日志等级
	JwtSecret   string    `json:"jwt_secret"` // jwt private key
	RedisURL    string    `json:"redis_url"`  // redis
	RedisPass   string    `json:"redis_pass"` // redis pwd
}

var config = new(LocalConfig)
var configPath = "bili.json"

func InitConfig() {
	Info("[CONF] Loading Config")
	if !Exists("bili.json") {
		var config LocalConfig
		var l = "debug"
		config.Finger = ""
		config.LoggerLevel = &l
		config.APIAddr = ":9090"
		_ = config.SetConfig()
	}
}

func GetConfig() *LocalConfig {
	if config.CacheTime.Before(time.Now()) {
		if err := LoadConfig(configPath, config); err != nil {
			Info("loading file failed")
			return nil
		}
		config.CacheTime = time.Now().Add(time.Second * 60)
	}
	return config
}

/**
 * save cnf/conf.json
 */
func (lc *LocalConfig) SetConfig() error {
	fp, err := os.Create(configPath)
	if err != nil {
		Info("[CONF] loading file failed")
	}
	defer fp.Close()
	data, err := json.Marshal(lc)
	if err != nil {
		Info("[CONF] marshal file failed")
	}
	n, err := fp.Write(data)
	if err != nil {
		Info("[CONF] write file failed")
	}
	Info("[CONF] already update config file | SIZE " + strconv.Itoa(n))
	return nil
}

const configFileSizeLimit = 10 << 20

/**
 * Load File
 * @param path 文件路径
 * @param dist 存放目标
 */
func LoadConfig(path string, dist interface{}) error {
	configFile, err := os.Open(path)
	if err != nil {
		Info("[CONF] Failed to open config file.")
		return err
	}

	fi, _ := configFile.Stat()
	if size := fi.Size(); size > (configFileSizeLimit) {
		Info("[CONF] Config file size exceeds reasonable limited")
		return errors.New("limited")
	}

	if fi.Size() == 0 {
		Info("[CONF] Config file is empty, skipping")
		return errors.New("empty")
	}

	buffer := make([]byte, fi.Size())
	_, err = configFile.Read(buffer)
	buffer, err = StripComments(buffer)
	if err != nil {
		Info("[CONF] Failed to strip comments from json")
		return err
	}

	buffer = []byte(os.ExpandEnv(string(buffer)))

	err = json.Unmarshal(buffer, &dist)
	if err != nil {
		Info("Failed unmarshalling json")
		return err
	}
	return nil
}

/**
 * 注释清除
 */
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
