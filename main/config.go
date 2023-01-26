// @author zhuweitung 2023/1/24
package main

type Config struct {
	Cookie                 string       `yaml:"cookie"`
	Cron                   string       `yaml:"cron"`
	TopicReplyPool         []string     `yaml:"topicReplyPool"`
	LatestTopicId          string       `yaml:"latestTopicId"`
	LatestTopicPublishTime string       `yaml:"latestTopicPublishTime"`
	Notification           Notification `yaml:"notification"`
}

type Notification struct {
	Telegram Telegram `yaml:"telegram"` // 电报提醒配置
}

type Telegram struct {
	BotToken  string `yaml:"botToken"`  // 机器人token
	ChatId    string `yaml:"chatId"`    // 消息接收方id
	HttpProxy string `yaml:"httpProxy"` // 消息发送代理
}
