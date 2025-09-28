package main

import (
	"fmt"
	"github.com/eatmoreapple/openwechat"
	"gmodUtils/service"
	"log"
	"time"
)

func main() {
	service.Pic()
}

func main2() {
	bot := openwechat.DefaultBot()

	// 设置消息处理函数
	bot.MessageHandler = func(msg *openwechat.Message) {
		// 空处理，仅用于保持登录状态
	}

	// 设置登录回调
	bot.UUIDCallback = openwechat.PrintlnQrcodeUrl

	// 创建热登录存储
	storage := openwechat.NewFileHotReloadStorage("storage.json")
	defer storage.Close()

	// 执行热登录
	if err := bot.HotLogin(storage); err != nil {
		fmt.Println("热登录失败，尝试普通登录")
		if err = bot.Login(); err != nil {
			log.Fatal(err)
		}
	}

	// 获取所有群组
	self, err := bot.GetCurrentUser()
	if err != nil {
	}
	groups, err := self.Groups()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("找到 %d 个群组\n", len(groups))

	// 遍历所有群组
	for _, group := range groups {
		fmt.Printf("\n群组: %s\n", group.NickName)

		// 获取群成员
		members, err := group.Members()
		if err != nil {
			log.Printf("获取群成员失败: %v", err)
			continue
		}

		fmt.Printf("成员数量: %d\n", len(members))

		// 打印每个成员的昵称
		for _, member := range members {
			fmt.Printf(" - %s\n", member.NickName)
		}

		// 添加延迟避免请求过于频繁
		time.Sleep(1 * time.Second)
	}

	// 阻塞主goroutine
	bot.Block()
}
