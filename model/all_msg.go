package model

import (
	"ApiBot/handler"
	"ApiBot/util"
)

func AllMsg(Data map[string]interface{}) (string, bool, interface{}) {
	message := Data["NewMessage"].(string)
	var (
		rsg        string
		AutoEscape bool
		err        interface{}
	)
	if message == "菜单" {
		switch Data["message_type"] {
		case "group":
			rsg = "奖励：小姐姐 看美女 cos\n小功能：啾啾 一言 舔狗日记\nGPT：艾特我 问题 或者 q 问题，就能直接向我提问喽\n发送gptcls可清除对话缓存\n发送gpt-money可查询余额"
		case "private":
			rsg = "奖励：小姐姐 看美女 cos\n小功能：啾啾 一言 舔狗日记\nGPT：直接发送问题即可\n发送gptcls可清除对话缓存\n发送gpt-money可查询余额"
		}
		return rsg, true, nil
	}

	// 小功能：啾啾 一言 舔狗日记
	switch message {
	case "啾啾":
		rsg, AutoEscape, err = handler.Jiujiu(int64(Data["user_id"].(float64)))
		return rsg, AutoEscape, err
	case "一言":
		rsg, AutoEscape, err = handler.Yiyan()
		return rsg, AutoEscape, err
	case "舔狗日记":
		rsg, AutoEscape, err = handler.Dogdiary()
		return rsg, AutoEscape, err
	}

	// 奖励：小姐姐 看美女 cos
	switch message {
	case "cos":
		rsg, AutoEscape, err = handler.Cos()
		return rsg, AutoEscape, err
	case "看美女":
		rsg, AutoEscape, err = handler.Kanmn()
		return rsg, AutoEscape, err
	case "小姐姐":
		rsg, AutoEscape, err = handler.Xjj()
		return rsg, AutoEscape, err
	}

	// GPT扩展功能
	switch message {
	case "gpt-money":
		rsg, AutoEscape, err = handler.GPTMoneyQuery()
		return rsg, AutoEscape, err
	case "gptcls":
		if util.Cfg.Redis.Switch == true {
			rsg, AutoEscape, err = handler.GPTCacheClear(Data)
			return rsg, AutoEscape, err
		} else {
			return "未开启GPT连续对话功能", true, nil
		}
	}
	return "", true, nil
}
