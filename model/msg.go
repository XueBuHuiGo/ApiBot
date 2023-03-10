package model

import (
	"ApiBot/util"
	"fmt"
	"github.com/goccy/go-json"
	"log"
	"regexp"
	"strings"
)

// HandleMessage 处理消息
func HandleMessage(Data map[string]interface{}) (string, bool, error) {
	// 定义返回值
	var rsg string
	var AutoEscape bool
	var err interface{}
	// 处理接收到的数据：去除首尾空格，去除CQ码
	message, _ := Data["message"].(string)
	message = regexp.MustCompile(`\[CQ:.*?]`).ReplaceAllString(message, "")
	message = strings.TrimSpace(message)
	Data["NewMessage"] = message
	if message != "" {
		//检索自定义词库
		for word, reply := range util.WordMap {
			if strings.Contains(strings.ToLower(message), word) {
				// 包含词库中的单词或短语，输出对应的回复内容
				AutoEscape = true
				rsg = reply
				break
			}
		}
		if rsg == "" {
			rsg, AutoEscape, err = AllMsg(Data)
		}

		if rsg == "" {
			switch Data["message_type"].(string) {
			case "private":
				rsg, AutoEscape, err = PrivateMsg(Data)
			case "group":
				rsg, AutoEscape, err = GroupMsg(Data)
			}
		}

		if err != nil {
			return rsg, AutoEscape, fmt.Errorf("HandleMessage failed: %v", err)
		}
	}
	return rsg, AutoEscape, nil
}

// SendMessage 发送消息
func SendMessage(msg string, AutoEscape bool, Data map[string]interface{}) ([]byte, error) {
	if len(msg) > 5000 {
		msg = msg[:4900] + "\n（由于QQ限制，最多发送5k字）"
	}
	postData := map[string]interface{}{
		"message":     msg,
		"auto_escape": AutoEscape,
	}
	//判断，并往postData里增添数据
	if Data["message_type"] != nil {
		postData["message_type"] = Data["message_type"].(string)
	}
	if Data["user_id"] != nil {
		postData["user_id"] = int64(Data["user_id"].(float64))
	}
	if Data["group_id"] != nil {
		postData["group_id"] = int64(Data["group_id"].(float64))
	}
	jsonData, err := json.Marshal(postData)
	if err != nil {
		return nil, fmt.Errorf("SendMessage failed to parse jsonData: %v", err)
	}
	header := []string{
		"Content-Type: application/json; charset=utf-8",
		"Referer: never",
		"User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.102 Safari/537.36",
	}
	log.Println("Post数据：", string(jsonData))
	resp, err := util.MyRequest("POST", util.Cfg.Bot.BotApi+"send_msg", header, string(jsonData))
	if err != nil {
		return nil, fmt.Errorf("SendMessage failed to MyPOST: %v", err)
	}
	return resp, nil
}

func sendMsgByOwn(messageType string, msg interface{}, sendID int64, autoEscape bool) ([]byte, error) {
	postData := map[string]interface{}{
		"message_type": messageType,
		"auto_escape":  autoEscape,
	}

	switch messageType {
	case "private":
		postData["user_id"] = sendID
	case "group":
		postData["group_id"] = sendID
	}

	var message string
	switch m := msg.(type) {
	case []byte:
		message = string(m[:5000])
	case string:
		message = m[:5000]
	default:
		return nil, fmt.Errorf("unsupported message type: %T", msg)
	}
	postData["message"] = message

	jsonData, err := json.Marshal(postData)
	if err != nil {
		return nil, err
	}

	header := []string{
		"Content-Type: application/json; charset=utf-8",
		"Referer: never",
		"User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.102 Safari/537.36",
	}
	resp, err := util.MyRequest("POST", util.Cfg.Bot.BotApi+"send_msg", header, string(jsonData))
	if err != nil {
		return nil, fmt.Errorf("sendMsgByOwn failed to MyPOST: %v", err)
	}
	return resp, nil
}

// GetMessage 获取回复消息内容
func GetMessage(Data map[string]interface{}) ([]byte, error) {
	postData := map[string]interface{}{
		"message_id": int64(Data["message_id"].(float64)),
	}
	jsonData, err := json.Marshal(postData)
	if err != nil {
		log.Fatal(err)
	}
	header := []string{
		"Content-Type: application/json; charset=utf-8",
		"Referer: never",
		"User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.102 Safari/537.36",
	}
	resp, err := util.MyRequest("POST", util.Cfg.Bot.BotApi+"get_msg", header, string(jsonData))
	if err != nil {
		return nil, fmt.Errorf("GetMessage failed to MyPOST: %v", err)
	}
	return resp, nil
}
