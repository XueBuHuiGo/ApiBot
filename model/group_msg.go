package model

import (
	"ApiBot/handler"
	"ApiBot/util"
	"fmt"
	"strings"
)

func GroupMsg(Data map[string]interface{}) (string, bool, interface{}) {
	message := Data["NewMessage"].(string)
	var (
		rsg string
		//AutoEscape bool
		err interface{}
	)
	// 艾特或者q xxx，就调用ChatGPT回复
	cqstr := fmt.Sprintf("[CQ:at,qq=%v]", util.Cfg.Bot.BotQq)
	if strings.HasPrefix(message, "q") || strings.Contains(Data["message"].(string), cqstr) {
		if strings.HasPrefix(message, "q") {
			message = strings.TrimSpace(message[1:])
		}
		rsg, err = handler.ChatGPT(Data, message)
		if err != nil {
			return "GPT 请求失败", false, err
		}
		return rsg, false, nil
	}
	return rsg, false, nil
}
