package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/auth"
	"github.com/gotd/td/tg"
	"github.com/sashabaranov/go-openai"
	"gopkg.in/yaml.v3"
)

// Config 配置结构
type Config struct {
	APIID       int    `yaml:"api_id"`
	APIHash     string `yaml:"api_hash"`
	PhoneNumber string `yaml:"phone_number"`
	OpenAIKey   string `yaml:"openai_key"`
	MinMessages int    `yaml:"min_messages"`
	Cooldown    int    `yaml:"cooldown"`
}

// GroupTracker 群组状态跟踪
type GroupTracker struct {
	sync.Mutex
	Counts    map[int64]int
	Histories map[int64][]string
	LastReply map[int64]time.Time
}

var (
	ctx          = context.Background()
	config       Config
	client       *telegram.Client
	aiClient     *openai.Client
	groupTracker = GroupTracker{
		Counts:    make(map[int64]int),
		Histories: make(map[int64][]string),
		LastReply: make(map[int64]time.Time),
	}
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
	if err := loadConfig("config.yaml"); err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	aiClient = openai.NewClient(config.OpenAIKey)

	client = telegram.NewClient(config.APIID, config.APIHash, telegram.Options{})

	ctx := context.Background()
	if err := client.Run(ctx, func(ctx context.Context) error {
		if err := authFlow(ctx); err != nil {
			return err
		}
		return runMonitor(ctx)
	}); err != nil {
		log.Fatalf("运行失败: %v", err)
	}
}

func loadConfig(file string) error {
	data, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(data, &config)
}

func authFlow(ctx context.Context) error {
	flow := auth.NewFlow(
		TermAuth{phone: config.PhoneNumber},
		auth.SendCodeOptions{},
	)

	if err := client.Auth().IfNecessary(ctx, flow); err != nil {
		return fmt.Errorf("认证失败: %w", err)
	}

	self, err := client.Self(ctx)
	if err != nil {
		return fmt.Errorf("获取用户信息失败: %w", err)
	}

	log.Printf("登录成功: %s (%s)", self.Username, self.Phone)
	return nil
}

func runMonitor(ctx context.Context) error {
	// 初始化状态
	state, err := client.API().UpdatesGetState(ctx)
	if err != nil {
		return fmt.Errorf("获取初始状态失败: %w", err)
	}

	// 创建消息处理器
	msgHandler := func(ctx context.Context, ent tg.Entities, msg *tg.Message) error {
		if msg.Out {
			return nil // 忽略自己发送的消息
		}

		var chatID int64
		switch peer := msg.PeerID.(type) {
		case *tg.PeerChat:
			chatID = peer.ChatID
		case *tg.PeerChannel:
			if ch, ok := ent.Channels[peer.ChannelID]; ok {
				chatID = ch.ID
			}
		default:
			return nil
		}

		if chatID == 0 {
			return nil
		}

		return processMessage(ctx, chatID, msg.Message)
	}

	// 主循环
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			// 获取更新差异
			diff, err := client.API().UpdatesGetDifference(ctx, &tg.UpdatesGetDifferenceRequest{
				Pts:  state.Pts,
				Date: int(time.Now().Unix()),
			})
			if err != nil {
				log.Printf("获取更新差异失败: %v", err)
				time.Sleep(5 * time.Second)
				continue
			}

			// 处理更新
			switch d := diff.(type) {
			case *tg.UpdatesDifference:
				state.Pts = d.State.Pts
				for _, update := range d.OtherUpdates {
					if err := handleUpdate(ctx, update, msgHandler); err != nil {
						log.Printf("处理更新失败: %v", err)
					}
				}
			case *tg.UpdatesDifferenceSlice:
				state.Pts = d.IntermediateState.Pts
				for _, update := range d.OtherUpdates {
					if err := handleUpdate(ctx, update, msgHandler); err != nil {
						log.Printf("处理更新失败: %v", err)
					}
				}
			}

			time.Sleep(10 * time.Second)
		}
	}
}

// 处理单个更新
func handleUpdate(ctx context.Context, update tg.UpdateClass, handler func(context.Context, tg.Entities, *tg.Message) error) error {
	switch u := update.(type) {
	case *tg.UpdateNewMessage:
		if msg, ok := u.Message.(*tg.Message); ok {
			return handler(ctx, tg.Entities{}, msg)
		}
	case *tg.UpdateNewChannelMessage:
		if msg, ok := u.Message.(*tg.Message); ok {
			return handler(ctx, tg.Entities{}, msg)
		}
	}
	return nil
}

func processMessage(ctx context.Context, chatID int64, text string) error {
	groupTracker.Lock()
	defer groupTracker.Unlock()

	// 初始化群组记录
	if _, ok := groupTracker.Histories[chatID]; !ok {
		groupTracker.Histories[chatID] = []string{}
	}

	// 更新消息历史
	groupTracker.Histories[chatID] = append(groupTracker.Histories[chatID], text)
	if len(groupTracker.Histories[chatID]) > config.MinMessages {
		groupTracker.Histories[chatID] = groupTracker.Histories[chatID][1:]
	}

	// 更新计数器
	groupTracker.Counts[chatID]++

	// 检查回复条件
	if groupTracker.Counts[chatID] >= config.MinMessages &&
		time.Since(groupTracker.LastReply[chatID]) > time.Duration(config.Cooldown)*time.Minute {

		go func() {
			if err := sendAIResponse(context.Background(), chatID); err != nil {
				log.Printf("回复失败: %v", err)
			} else {
				groupTracker.Lock()
				groupTracker.Counts[chatID] = 0
				groupTracker.LastReply[chatID] = time.Now()
				groupTracker.Unlock()
			}
		}()
	}

	return nil
}

func sendAIResponse(ctx context.Context, chatID int64) error {
	groupTracker.Lock()
	messages := make([]string, len(groupTracker.Histories[chatID]))
	copy(messages, groupTracker.Histories[chatID])
	groupTracker.Unlock()

	// 构建AI提示
	prompt := fmt.Sprintf(`你正在参与一个中文Telegram群组聊天，以下是最近的%d条消息：
-----
%s
-----
请生成一个自然、友好的中文回复，要求：
1. 符合上下文语境
2. 不超过2句话
3. 可以提问或回应特定消息`, config.MinMessages, strings.Join(messages, "\n"))

	// 调用ChatGPT
	resp, err := aiClient.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: openai.GPT3Dot5Turbo,
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleUser, Content: prompt},
		},
		Temperature: 0.7,
		MaxTokens:   60,
	})
	if err != nil || len(resp.Choices) == 0 {
		return fmt.Errorf("AI生成失败: %v", err)
	}

	// 发送回复
	_, err = client.API().MessagesSendMessage(ctx, &tg.MessagesSendMessageRequest{
		Peer:    &tg.InputPeerChat{ChatID: chatID},
		Message: resp.Choices[0].Message.Content,
	})
	return err
}

// TermAuth 认证实现
type TermAuth struct{ phone string }

func (t TermAuth) Phone(_ context.Context) (string, error)    { return t.phone, nil }
func (t TermAuth) Password(_ context.Context) (string, error) { return "cywhoyi1989", nil }
func (t TermAuth) Code(_ context.Context, _ *tg.AuthSentCode) (string, error) {
	fmt.Print("请输入验证码: ")
	var code string
	_, err := fmt.Scanln(&code)
	return code, err
}
func (t TermAuth) SignUp(_ context.Context) (auth.UserInfo, error) {
	return auth.UserInfo{}, nil
}
func (t TermAuth) AcceptTermsOfService(_ context.Context, tos tg.HelpTermsOfService) error {
	return nil
}
