package handler

import (
	"ApiBot/util"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/goccy/go-json"
	"golang.org/x/net/context"
	"log"
	"strconv"
	"strings"
	"time"
)

const OpenaiApiUrl = "https://agent-openai.ccrui.dev/v1/chat/completions"

type OpenAiRcv struct {
	Id      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		//Text         string      `json:"text"`
		//Logprobs     interface{} `json:"logprobs"`
		FinishReason string `json:"finish_reason"`
		Index        int    `json:"index"`
	} `json:"choices"`
	Usage struct {
		PromptTokens    int `json:"prompt_tokens"`
		CompletionTokes int `json:"completion_tokens"`
		TotalTokens     int `json:"total_tokens"`
	} `json:"usage"`
}

// ChatGPT 调用openai的API生成文本
func ChatGPT(Data map[string]interface{}, text string) (string, error) {
	var GPTHF []map[string]interface{}
	var GPTCache string
	var client *redis.Client
	if util.Cfg.Redis.Switch == true {
		err := util.ConnectRedis()
		if err != nil {
			return "", err
		}
		client = util.GetRedisClient()
		GPTCache, err = client.Get(context.Background(), "GPTCache_"+strconv.FormatInt(int64(Data["user_id"].(float64)), 10)).Result()
		if err != nil {
			log.Println("获取不到数据/数据过期,直接进行请求", err)
		}
	}
	log.Println("正在调用OpenAI API生成文本...", text)
	GPTHFYS := map[string]interface{}{"role": "user", "content": fmt.Sprintf(GPTCache+"ME: %v \nAI: ", text)}
	GPTHF = append(GPTHF, GPTHFYS)
	postData := map[string]interface{}{
		"model":             util.Cfg.OpenAi.Model,
		"messages":          GPTHF,
		"max_tokens":        util.Cfg.OpenAi.MaxTokens,
		"temperature":       util.Cfg.OpenAi.Temperature,
		"top_p":             util.Cfg.OpenAi.Top_p,
		"frequency_penalty": util.Cfg.OpenAi.Frequency_penalty,
		"presence_penalty":  util.Cfg.OpenAi.Presence_penalty,
	}
	header := []string{
		"Content-Type: application/json; charset=utf-8",
		"Authorization: " + "Bearer " + util.Cfg.OpenAi.ApiKey,
	}
	jsonData, _ := json.Marshal(postData)
	log.Println(string(jsonData))
	resp, err := util.MyRequest("POST", OpenaiApiUrl, header, string(jsonData))
	var openAiRcv OpenAiRcv
	err = json.Unmarshal(resp, &openAiRcv)
	if len(openAiRcv.Choices) == 0 || err != nil {
		return "ChatGPT 请求失败", err
	}
	openAiRcv.Choices[0].Message.Content = strings.Trim(openAiRcv.Choices[0].Message.Content, "\n\n")
	if util.Cfg.Redis.Switch == true {
		client.Set(context.Background(), "GPTCache_"+strconv.FormatInt(int64(Data["user_id"].(float64)), 10), GPTCache+fmt.Sprintf("ME: %v AI:%v ;\n", text, openAiRcv.Choices[0].Message.Content), 10*time.Minute)
		closeerr := client.Close()
		if closeerr != nil {
			log.Println("关闭Redis客户端失败")
		}
	}
	if Data["message_type"].(string) == "private" {
		rsg := openAiRcv.Choices[0].Message.Content
		return rsg, nil
	}
	rsg := fmt.Sprintf("[CQ:reply,id=%v]%v", int64(Data["message_id"].(float64)), openAiRcv.Choices[0].Message.Content)
	return rsg, nil
}

func GPTMoneyQuery() (string, bool, error) {
	var data map[string]interface{}
	header := []string{
		"Content-Type: application/json; charset=utf-8",
		"Authorization: " + "Bearer " + util.Cfg.OpenAi.ApiKey,
	}
	rsp, err := util.MyRequest("GET", "https://agent-openai.ccrui.dev/dashboard/billing/credit_grants", header, "")
	err = json.Unmarshal(rsp, &data)
	if err != nil {
		return "GPT余额 查询失败", true, err
	}
	rsg := fmt.Sprintf("一共：$%v\n已用：$%.3f\n剩余：$%.3f", data["total_granted"].(float64), data["total_used"].(float64), data["total_available"].(float64))
	return rsg, true, nil
}

func GPTCacheClear(Data map[string]interface{}) (string, bool, error) {
	var rsg string
	err := util.ConnectRedis()
	if err != nil {
		return "Redis连接失败,清除GPT缓存失败", true, err
	}
	client := util.GetRedisClient()
	_, err = client.Del(context.Background(), "GPTCache_"+strconv.FormatInt(int64(Data["user_id"].(float64)), 10)).Result()
	// 关闭 Redis 客户端
	closeerr := client.Close()
	if closeerr != nil {
		log.Println("关闭Redis客户端失败")
	}
	if err != nil {
		switch Data["message_type"].(string) {
		case "private":
			return "清除GPT缓存失败", true, err
		case "group":
			return fmt.Sprintf("[CQ:reply,id=%v]清除GPT缓存失败", int64(Data["message_id"].(float64))), false, err
		}
	} else {
		switch Data["message_type"].(string) {
		case "private":
			rsg = "清除GPT缓存成功"
		case "group":
			rsg = fmt.Sprintf("[CQ:reply,id=%v]清除GPT缓存成功", int64(Data["message_id"].(float64)))
		}
	}
	return rsg, false, nil
}
