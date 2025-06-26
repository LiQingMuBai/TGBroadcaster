package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/auth"
	"github.com/gotd/td/tg"
	"github.com/robfig/cron/v3"
	"gopkg.in/yaml.v3"
)

// 配置结构体
type Config struct {
	APIID         int    `yaml:"api_id"`
	APIHash       string `yaml:"api_hash"`
	PhoneNumber   string `yaml:"phone_number"`
	MessageGroups []struct {
		GroupIDs      []int64  `yaml:"group_ids"`
		Messages      []string `yaml:"messages"`
		IntervalMs    int      `yaml:"interval_ms"`    // 消息间隔毫秒数
		DailySchedule []string `yaml:"daily_schedule"` // 每天发送时间点，格式"15:04"
	} `yaml:"message_groups"`
}

var (
	config Config
	client *telegram.Client
	ctx    = context.Background()
)

// 获取用户加入的所有群组/频道
func getAllGroups(ctx context.Context, api *tg.Client) ([]tg.ChatClass, error) {
	var (
		groups []tg.ChatClass
		offset int
		limit  = 100
	)

	for {
		// 分页获取对话列表
		dialogs, err := api.MessagesGetDialogs(ctx, &tg.MessagesGetDialogsRequest{
			OffsetPeer: &tg.InputPeerEmpty{},
			Limit:      limit,
			OffsetDate: offset,
		})
		if err != nil {
			return nil, fmt.Errorf("获取对话失败: %w", err)
		}

		switch d := dialogs.(type) {
		case *tg.MessagesDialogs:
			// 提取群组
			for _, chat := range d.Chats {
				if isGroup(chat) {
					groups = append(groups, chat)
				}
			}
			return groups, nil

		case *tg.MessagesDialogsSlice:
			// 处理分页数据
			for _, chat := range d.Chats {
				if isGroup(chat) {
					groups = append(groups, chat)
				}
			}

			if len(d.Dialogs) < limit {
				return groups, nil
			}
			offset = d.Dialogs[len(d.Dialogs)-1].GetTopMessage()

		default:
			return nil, fmt.Errorf("未知的响应类型: %T", dialogs)
		}
	}
}

// 判断是否是群组/频道
func isGroup(chat tg.ChatClass) bool {
	switch chat.(type) {
	case *tg.Chat, *tg.Channel:
		return true
	default:
		return false
	}
}

// 获取真实ID (处理超级群组的-100前缀)
func getRealID(chat tg.ChatClass) int64 {
	switch c := chat.(type) {
	case *tg.Chat:
		return c.ID
	case *tg.Channel:
		return c.ID
	default:
		return 0
	}
}
func main() {
	// 加载配置文件
	loadConfig("config.yaml")

	// 初始化Telegram客户端
	initClient()

	// 设置定时任务
	setupScheduler()

	// 保持程序运行
	select {}
}

// 加载配置文件
func loadConfig(filename string) {
	data, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalf("读取配置文件失败: %v", err)
	}

	if err := yaml.Unmarshal(data, &config); err != nil {
		log.Fatalf("解析配置文件失败: %v", err)
	}
}

// 初始化Telegram客户端
func initClient() {
	// 创建客户端选项
	options := telegram.Options{
		Device: telegram.DeviceConfig{
			DeviceModel:   "MyUserBot",
			SystemVersion: "1.0",
		},
	}

	// 创建客户端
	client = telegram.NewClient(config.APIID, config.APIHash, options)

	// 认证流程
	authFlow := auth.NewFlow(
		TermAuth{phone: config.PhoneNumber},
		auth.SendCodeOptions{},
	)

	// 运行客户端
	if err := client.Run(ctx, func(ctx context.Context) error {
		// 认证
		if err := client.Auth().IfNecessary(ctx, authFlow); err != nil {
			return fmt.Errorf("认证失败: %w", err)
		}

		// 获取自身信息
		self, err := client.Self(ctx)
		if err != nil {
			return fmt.Errorf("获取自身信息失败: %w", err)
		}

		log.Printf("登录成功: %s (%s)", self.Username, self.Phone)
		return nil
	}); err != nil {
		log.Fatalf("客户端运行失败: %v", err)
	}
}

// 设置定时任务
func setupScheduler() {
	c := cron.New(cron.WithLocation(time.Local))

	for i, group := range config.MessageGroups {
		for _, schedule := range group.DailySchedule {
			// 为每个群组和每个时间点创建定时任务
			_, err := c.AddFunc(fmt.Sprintf("0 %s", schedule), func(groupIndex int) func() {
				return func() {
					sendMessagesToGroup(groupIndex)
				}
			}(i))
			if err != nil {
				log.Printf("创建定时任务失败: %v", err)
			}
		}
	}

	c.Start()
	log.Println("定时任务已启动")
}

// 向群组发送消息
func sendMessagesToGroup(groupIndex int) {
	group := config.MessageGroups[groupIndex]
	var wg sync.WaitGroup

	for _, groupID := range group.GroupIDs {
		wg.Add(1)
		go func(chatID int64) {
			defer wg.Done()
			sendMessagesWithInterval(chatID, group.Messages, group.IntervalMs)
		}(groupID)
	}

	wg.Wait()
	log.Printf("群组消息发送完成: 组%d", groupIndex)
}

// 按间隔发送多条消息
func sendMessagesWithInterval(chatID int64, messages []string, intervalMs int) {
	api := client.API()

	// 获取所有群组
	groups, err := getAllGroups(ctx, client.API())
	if err != nil {

	}

	// 打印结果
	for _, g := range groups {
		fmt.Printf("群组:  (ID: %d)\n", getRealID(g))
	}

	for _, msg := range messages {
		// 构造输入Peer
		inputPeer := &tg.InputPeerChannel{
			ChannelID:  chatID,
			AccessHash: 0, // 如果知道access hash可以设置
		}

		// 发送消息
		_, err := api.MessagesSendMessage(ctx, &tg.MessagesSendMessageRequest{
			Peer:    inputPeer,
			Message: msg,
		})
		if err != nil {
			log.Printf("发送消息到群组%d失败: %v", chatID, err)
			continue
		}

		log.Printf("已发送消息到群组%d: %s", chatID, msg)

		// 按配置间隔等待
		time.Sleep(time.Duration(intervalMs) * time.Millisecond)
	}
}

// 终端认证实现
type TermAuth struct {
	phone string
}

func (t TermAuth) Phone(_ context.Context) (string, error) {
	return t.phone, nil
}

func (t TermAuth) Password(_ context.Context) (string, error) {
	fmt.Print("请输入密码: ")
	var password string
	_, err := fmt.Scanln(&password)
	return password, err
}

func (t TermAuth) Code(_ context.Context, _ *tg.AuthSentCode) (string, error) {
	fmt.Print("请输入验证码: ")
	var code string
	_, err := fmt.Scanln(&code)
	return code, err
}

func (t TermAuth) SignUp(_ context.Context) (auth.UserInfo, error) {
	return auth.UserInfo{}, fmt.Errorf("不支持注册")
}

func (t TermAuth) AcceptTermsOfService(_ context.Context, tos tg.HelpTermsOfService) error {
	return nil
}
