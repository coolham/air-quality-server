package mqtt

import (
	"air-quality-server/internal/config"
	"air-quality-server/internal/utils"
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// Client MQTT客户端
type Client struct {
	client   mqtt.Client
	config   *config.MQTTConfig
	logger   utils.Logger
	handlers map[string]MessageHandler
	mu       sync.RWMutex
	ctx      context.Context
	cancel   context.CancelFunc
}

// MessageHandler 消息处理器接口
type MessageHandler interface {
	HandleMessage(topic string, payload []byte) error
}

// NewClient 创建MQTT客户端
func NewClient(cfg *config.MQTTConfig, logger utils.Logger) *Client {
	ctx, cancel := context.WithCancel(context.Background())

	return &Client{
		config:   cfg,
		logger:   logger,
		handlers: make(map[string]MessageHandler),
		ctx:      ctx,
		cancel:   cancel,
	}
}

// Connect 连接到MQTT Broker
func (c *Client) Connect() error {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(c.config.Broker)
	opts.SetClientID(c.config.ClientID)
	opts.SetUsername(c.config.Username)
	opts.SetPassword(c.config.Password)
	opts.SetKeepAlive(time.Duration(c.config.KeepAlive) * time.Second)
	opts.SetCleanSession(c.config.CleanSession)
	opts.SetAutoReconnect(c.config.AutoReconnect)
	opts.SetMaxReconnectInterval(time.Duration(c.config.MaxReconnectInterval) * time.Second)
	opts.SetConnectTimeout(time.Duration(c.config.ConnectTimeout) * time.Second)
	opts.SetWriteTimeout(time.Duration(c.config.WriteTimeout) * time.Second)
	// SetReadTimeout 方法在 paho.mqtt.golang 中不存在，移除这行

	// 设置连接回调
	opts.SetOnConnectHandler(c.onConnect)
	opts.SetConnectionLostHandler(c.onConnectionLost)
	opts.SetReconnectingHandler(c.onReconnecting)

	// 创建客户端
	c.client = mqtt.NewClient(opts)

	// 连接
	if token := c.client.Connect(); token.Wait() && token.Error() != nil {
		c.logger.Error("MQTT连接失败", utils.ErrorField(token.Error()))
		return token.Error()
	}

	c.logger.Info("MQTT客户端连接成功",
		utils.String("broker", c.config.Broker),
		utils.String("client_id", c.config.ClientID))

	return nil
}

// Disconnect 断开连接
func (c *Client) Disconnect() {
	c.cancel()
	if c.client != nil && c.client.IsConnected() {
		c.client.Disconnect(250)
		c.logger.Info("MQTT客户端已断开连接")
	}
}

// Subscribe 订阅主题
func (c *Client) Subscribe(topic string, handler MessageHandler) error {
	c.mu.Lock()
	c.handlers[topic] = handler
	c.mu.Unlock()

	if !c.client.IsConnected() {
		return fmt.Errorf("MQTT客户端未连接")
	}

	token := c.client.Subscribe(topic, byte(c.config.QoS), c.messageHandler)
	if token.Wait() && token.Error() != nil {
		c.logger.Error("订阅主题失败",
			utils.String("topic", topic),
			utils.ErrorField(token.Error()))
		return token.Error()
	}

	c.logger.Info("订阅主题成功", utils.String("topic", topic))
	return nil
}

// Unsubscribe 取消订阅主题
func (c *Client) Unsubscribe(topic string) error {
	if !c.client.IsConnected() {
		return fmt.Errorf("MQTT客户端未连接")
	}

	token := c.client.Unsubscribe(topic)
	if token.Wait() && token.Error() != nil {
		c.logger.Error("取消订阅主题失败",
			utils.String("topic", topic),
			utils.ErrorField(token.Error()))
		return token.Error()
	}

	c.mu.Lock()
	delete(c.handlers, topic)
	c.mu.Unlock()

	c.logger.Info("取消订阅主题成功", utils.String("topic", topic))
	return nil
}

// Publish 发布消息
func (c *Client) Publish(topic string, payload interface{}) error {
	if !c.client.IsConnected() {
		return fmt.Errorf("MQTT客户端未连接")
	}

	var data []byte
	var err error

	switch v := payload.(type) {
	case []byte:
		data = v
	case string:
		data = []byte(v)
	default:
		data, err = json.Marshal(payload)
		if err != nil {
			c.logger.Error("序列化消息失败", utils.ErrorField(err))
			return err
		}
	}

	token := c.client.Publish(topic, byte(c.config.QoS), false, data)
	if token.Wait() && token.Error() != nil {
		c.logger.Error("发布消息失败",
			utils.String("topic", topic),
			utils.ErrorField(token.Error()))
		return token.Error()
	}

	c.logger.Debug("发布消息成功",
		utils.String("topic", topic),
		utils.Int("size", len(data)))

	return nil
}

// IsConnected 检查连接状态
func (c *Client) IsConnected() bool {
	return c.client != nil && c.client.IsConnected()
}

// messageHandler 消息处理回调
func (c *Client) messageHandler(client mqtt.Client, msg mqtt.Message) {
	topic := msg.Topic()
	payload := msg.Payload()

	c.logger.Debug("收到MQTT消息",
		utils.String("topic", topic),
		utils.Int("size", len(payload)))

	// 查找匹配的处理器
	c.mu.RLock()
	var handler MessageHandler
	for pattern, h := range c.handlers {
		if matchTopic(pattern, topic) {
			handler = h
			break
		}
	}
	c.mu.RUnlock()

	if handler != nil {
		if err := handler.HandleMessage(topic, payload); err != nil {
			c.logger.Error("处理MQTT消息失败",
				utils.String("topic", topic),
				utils.ErrorField(err))
		}
	} else {
		c.logger.Warn("未找到消息处理器", utils.String("topic", topic))
	}
}

// onConnect 连接成功回调
func (c *Client) onConnect(client mqtt.Client) {
	c.logger.Info("MQTT连接已建立")

	// 重新订阅所有主题
	c.mu.RLock()
	for topic := range c.handlers {
		if token := client.Subscribe(topic, byte(c.config.QoS), c.messageHandler); token.Wait() && token.Error() != nil {
			c.logger.Error("重新订阅主题失败",
				utils.String("topic", topic),
				utils.ErrorField(token.Error()))
		}
	}
	c.mu.RUnlock()
}

// onConnectionLost 连接丢失回调
func (c *Client) onConnectionLost(client mqtt.Client, err error) {
	c.logger.Error("MQTT连接丢失", utils.ErrorField(err))
}

// onReconnecting 重连回调
func (c *Client) onReconnecting(client mqtt.Client, options *mqtt.ClientOptions) {
	c.logger.Info("MQTT正在重连...")
}

// matchTopic 匹配主题模式
func matchTopic(pattern, topic string) bool {
	// 简单的通配符匹配，支持 + 和 #
	// + 匹配单级， # 匹配多级
	if pattern == topic {
		return true
	}

	// 这里可以实现更复杂的主题匹配逻辑
	// 暂时使用简单的字符串匹配
	return pattern == topic
}

// GetConnectionStatus 获取连接状态信息
func (c *Client) GetConnectionStatus() map[string]interface{} {
	status := map[string]interface{}{
		"connected": c.IsConnected(),
		"broker":    c.config.Broker,
		"client_id": c.config.ClientID,
		"username":  c.config.Username,
	}

	if c.client != nil {
		status["server_uri"] = c.config.Broker
	}

	return status
}
