// @author zhuweitung 2023/1/24
package main

type Config struct {
	Cookie                 string   `yaml:"cookie"`
	Cron                   string   `yaml:"cron"`
	TopicReplyPool         []string `yaml:"topicReplyPool"`
	LatestTopicId          string   `yaml:"latestTopicId"`
	LatestTopicPublishTime string   `yaml:"latestTopicPublishTime"`
}
