// @author zhuweitung 2023/1/22
package main

import (
	"encoding/json"
	"github.com/gocolly/colly/v2"
	"github.com/robfig/cron/v3"
	"log"
	"math/rand"
	"regexp"
	"runtime/debug"
	"strconv"
	"strings"
	"time"
)

func main() {

	defer func() {
		if err := recover(); err != nil {
			log.Fatalf("err=%v, stack=%s\n", err, string(debug.Stack()))
		}
	}()

	LoadConfig()
	if CONFIG.Cookie == "" {
		log.Fatalln("cookie为空")
		return
	}

	if CONFIG.Cron == "" {
		log.Fatalln("定时任务表达式为空")
		return
	}

	c := cron.New(cron.WithSeconds())
	c.AddFunc(CONFIG.Cron, run)
	c.Start()

	// 阻塞，让main函数不退出，保持程序运行
	select {}

}

func run() {

	// 设置随机种子
	rand.Seed(time.Now().UnixNano())

	// 评论池
	topicReplies := []string{"", "喵喵喵", "早上/中午/晚上好", "每日一水", "Have a nice day guys", "打个卡", "打卡", "签个到", "签到"}
	if len(CONFIG.TopicReplyPool) > 0 {
		topicReplies = CONFIG.TopicReplyPool
	}

	// nya表情
	emojis := []string{
		"{:10_737:}", "{:10_729:}", "{:10_730:}", "{:10_731:}", "{:10_732:}", "{:10_733:}", "{:10_734:}", "{:10_735:}",
		"{:10_736:}", "{:10_728:}", "{:10_738:}", "{:10_739:}", "{:10_740:}", "{:10_741:}", "{:10_742:}", "{:10_743:}",
		"{:10_744:}", "{:10_720:}", "{:10_712:}", "{:10_713:}", "{:10_714:}", "{:10_715:}", "{:10_716:}", "{:10_717:}",
		"{:10_718:}", "{:10_719:}", "{:10_711:}", "{:10_721:}", "{:10_722:}", "{:10_723:}", "{:10_724:}", "{:10_725:}",
		"{:10_726:}", "{:10_727:}", "",
	}

	c := colly.NewCollector()

	c.OnRequest(func(r *colly.Request) {
		// 设置cookie
		r.Headers.Add("cookie", CONFIG.Cookie)
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Fatalln("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	c.OnResponse(func(r *colly.Response) {
		if r.StatusCode == 200 && r.Request.URL.String() == "https://www.skyey2.com/login.php" {
			msg := "cookie失效，请获取最新的cookie并设置"
			log.Fatalln(msg)
			Send(msg)
		}
	})

	// 获取发帖中符合要求的第一个帖子
	c.OnHTML("#delform > table > tbody", func(e *colly.HTMLElement) {
		topicLinks := []string{}
		e.ForEachWithBreak("tr:not(:nth-child(1)) th", func(i int, th *colly.HTMLElement) bool {
			a := th.DOM.Find("a:first-child")
			topicTitle := a.Text()
			if strings.Index(topicTitle, "月水楼]") != -1 {
				href, _ := a.Attr("href")
				topicLinks = append(topicLinks, href)
				return false
			}
			return true
		})
		if len(topicLinks) > 0 {
			href := topicLinks[0]
			log.Println("访问水楼：" + href)
			e.Request.Visit(href)
		}
	})

	// 获取帖子发表时间
	c.OnHTML("#postlist > div:nth-child(3) em[id*=authorposton]", func(e *colly.HTMLElement) {
		// 获取帖子发布时间
		reg := regexp.MustCompile(`(\d{4}-\d{1,2}-\d{1,2}\s\d{2}:\d{2}:\d{2})`)
		match := reg.FindStringSubmatch(e.Text)
		if len(match) > 0 {
			publishTime := match[0]
			// 保存配置文件
			SaveConfig("latestTopicPublishTime", publishTime)
			log.Printf("水楼发布时间: %s\n", publishTime)
		}
	})

	// 获取帖子详情页面的formhash
	c.OnHTML("#fastpostform input[name=formhash]", func(e *colly.HTMLElement) {
		// 正则获取tid
		reg := regexp.MustCompile(`\&?tid=(\d+)\&?`)
		match := reg.FindStringSubmatch(e.Request.URL.String())
		if len(match) > 1 {
			id := match[1]
			// 保存配置文件
			SaveConfig("latestTopicId", id)
			formHash := e.Attr("value")
			postUrl := "https://www.skyey2.com/forum.php?mod=post&action=reply&fid=8&tid=" + id + "&extra=&replysubmit=yes&infloat=yes&handlekey=fastpost&inajax=1"
			// 随机评论
			rand.Shuffle(len(topicReplies), func(i, j int) { topicReplies[i], topicReplies[j] = topicReplies[j], topicReplies[i] })
			// 随机表情
			rand.Shuffle(len(emojis), func(i, j int) { emojis[i], emojis[j] = emojis[j], emojis[i] })
			message := topicReplies[0] + emojis[0]
			requestData := map[string]string{
				"message":  message,
				"posttime": strconv.FormatInt(time.Now().Unix(), 10),
				"formhash": formHash,
			}
			json, _ := json.Marshal(requestData)
			log.Println("灌水: " + string(json))
			e.Request.Post(postUrl, requestData)
		}
	})

	c.OnScraped(func(r *colly.Response) {
		//log.Println("结束", r.Request.URL)
	})

	// 读取配置文件，判断最新水楼发布时间是否当月。若配置项空或非当月 检索获取最新水楼；否则直接访问最新水楼页面
	if CONFIG.LatestTopicId == "" ||
		CONFIG.LatestTopicPublishTime == "" ||
		!strings.Contains(CONFIG.LatestTopicPublishTime, time.Now().Format("2006-1-")) {
		log.Println("检索最新水楼")
		c.Visit("https://www.skyey2.com/home.php?mod=space&uid=12617&do=thread&view=me&from=space")
	} else {
		href := "https://www.skyey2.com/forum.php?mod=viewthread&tid=" + CONFIG.LatestTopicId
		log.Println("访问水楼：" + href)
		c.Visit(href)
	}
}
