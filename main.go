package main

import (
	"ApiBot/model"
	"ApiBot/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func main() {
	// 创建 gin 实例
	// 正式发布模式
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	// 注册路由
	r.POST("/", func(c *gin.Context) {
		var Data map[string]interface{}
		err := c.ShouldBindJSON(&Data)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"retcode": http.StatusBadRequest, "msg": "请求出错"})
			return
		}
		go Message(Data)
		c.JSON(http.StatusOK, gin.H{"retcode": http.StatusOK, "msg": "请求成功"})
	})
	// 打印启动消息
	log.Println("[INFO]: 启动成功，当前版本1.2")
	log.Printf("[INFO]: ApiBot 服务器已启动: %v", util.Cfg.Servers.Address)
	// 启动服务
	err := r.Run(util.Cfg.Servers.Address)
	if err != nil {
		fmt.Println("r.Run error: ", err)
		return
	}
}

func Message(Data map[string]interface{}) {
	rsg, AutoEscape, err := model.HandleMessage(Data)
	if err != nil {
		log.Println("HandleMessage error: ", err)
	}
	if rsg != "" {
		_, err = model.SendMessage(rsg, AutoEscape, Data)
		if err != nil {
			log.Println("SendMessage error: ", err)
			return
		}
	}
}
