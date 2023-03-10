package util

import (
	"bufio"
	"fmt"
	"github.com/go-redis/redis/v8"
	"golang.org/x/net/context"
	"gopkg.in/yaml.v3"
	"os"
	"strings"
)

type Config struct {
	Servers struct {
		Address string `yaml:"Address"`
	} `yaml:"Servers"`
	Bot struct {
		BotApi string `yaml:"BotApi"`
		BotQq  int64  `yaml:"BotQq"`
		OwnQq  int64  `yaml:"OwnQq"`
	} `yaml:"Bot"`
	OpenAi struct {
		ApiKey            string  `yaml:"ApiKey"`
		Model             string  `yaml:"Model"`
		MaxTokens         int32   `yaml:"MaxTokens"`
		Temperature       float32 `yaml:"Temperature"`
		Top_p             float32 `yaml:"Top_p"`
		Frequency_penalty float32 `yaml:"Frequency_penalty"`
		Presence_penalty  float32 `yaml:"Presence_penalty"`
	} `yaml:"OpenAi"`
	Redis struct {
		Switch   bool   `yaml:"Switch"`
		Addr     string `yaml:"Addr"`
		Password string `yaml:"Password"`
		DB       int    `yaml:"DB"`
	} `yaml:"Redis"`
}

// LoadConfig 加载配置
func LoadConfig() Config {
	yamlFile, err := os.ReadFile("config.yml")
	if err != nil {
		panic(err)
	}
	var config Config
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		panic(err)
	}
	return config
}

func LoadWordlist() map[string]string {
	// 加载词库
	file, err := os.Open("wordlist.txt")
	if err != nil {
		fmt.Println("Open wordlist error: ", err)
		return nil
	}
	defer func() {
		err := file.Close()
		if err != nil {
			fmt.Println("Close wordlist error: ", err)
		}
	}()
	// 读取词库文件中的所有行，并将单词或短语和对应的回复内容存储在map中
	scanner := bufio.NewScanner(file)
	wordMap := make(map[string]string)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, "|")
		if len(parts) == 2 {
			word := strings.TrimSpace(parts[0])
			reply := strings.TrimSpace(parts[1])
			wordMap[word] = reply
		}
	}
	return wordMap
}

// Cfg 加载配置文件
var Cfg = LoadConfig()

// WordMap 加载自定义词库
var WordMap = LoadWordlist()

// Rdb Redis配置
var Rdb *redis.Client

func ConnectRedis() error {
	Rdb = redis.NewClient(&redis.Options{
		Addr:     Cfg.Redis.Addr,
		Password: Cfg.Redis.Password,
		DB:       Cfg.Redis.DB,
	})
	_, err := Rdb.Ping(context.Background()).Result()
	if err != nil {
		fmt.Printf("连接redis出错，错误信息：%v", err)
		return err
	}
	return nil
}
func GetRedisClient() *redis.Client {
	return Rdb
}
