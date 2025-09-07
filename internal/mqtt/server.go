package mqtt

import (
	"air-quality-server/internal/config"
	"air-quality-server/internal/utils"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"sync"
	"time"

	mqtt "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/hooks/auth"
	"github.com/mochi-mqtt/server/v2/hooks/storage"
	"github.com/mochi-mqtt/server/v2/listeners"
	"github.com/mochi-mqtt/server/v2/packets"
	"github.com/mochi-mqtt/server/v2/system"
)

// Server MQTT服务器（基于Mochi MQTT的嵌入式实现）
type Server struct {
	config            *config.MQTTConfig
	logger            utils.Logger
	mu                sync.RWMutex
	ctx               context.Context
	cancel            context.CancelFunc
	running           bool
	server            *mqtt.Server
	sensorDataHandler *SensorDataHandler
}

// NewServer 创建MQTT服务器
func NewServer(cfg *config.MQTTConfig, logger utils.Logger, sensorDataHandler *SensorDataHandler) *Server {
	ctx, cancel := context.WithCancel(context.Background())

	return &Server{
		config:            cfg,
		logger:            logger,
		ctx:               ctx,
		cancel:            cancel,
		sensorDataHandler: sensorDataHandler,
	}
}

// Start 启动MQTT服务器
func (s *Server) Start() error {
	s.logger.Info("🚀 开始启动MQTT服务器...",
		utils.String("config_broker", s.config.Broker),
		utils.String("config_client_id", s.config.ClientID),
		utils.Int("config_keep_alive", s.config.KeepAlive),
		utils.Bool("config_clean_session", s.config.CleanSession),
		utils.Int("config_qos", s.config.QoS),
		utils.Bool("config_auto_reconnect", s.config.AutoReconnect),
		utils.Int("config_connect_timeout", s.config.ConnectTimeout))

	// 创建Mochi MQTT服务器
	s.logger.Debug("📦 正在创建Mochi MQTT服务器实例...")
	s.server = mqtt.New(nil)
	s.logger.Info("✅ Mochi MQTT服务器实例已创建",
		utils.String("server_type", "mochi-mqtt"),
		utils.String("version", "v2"))

	// 添加认证钩子（允许所有连接，生产环境应使用自定义认证）
	s.logger.Debug("🔐 正在添加认证钩子...")
	if err := s.server.AddHook(new(auth.AllowHook), nil); err != nil {
		s.logger.Error("❌ 添加认证钩子失败", utils.ErrorField(err))
		return fmt.Errorf("添加认证钩子失败: %w", err)
	}
	s.logger.Info("✅ 认证钩子已添加",
		utils.String("hook_type", "AllowHook"),
		utils.String("description", "允许所有连接"))

	// 添加消息处理钩子
	s.logger.Debug("📨 正在添加消息处理钩子...")
	if err := s.server.AddHook(new(MessageHandlerHook), map[string]interface{}{
		"logger":            s.logger,
		"sensorDataHandler": s.sensorDataHandler,
	}); err != nil {
		s.logger.Error("❌ 添加消息处理钩子失败", utils.ErrorField(err))
		return fmt.Errorf("添加消息处理钩子失败: %w", err)
	}
	s.logger.Info("✅ 消息处理钩子已添加",
		utils.String("hook_type", "MessageHandlerHook"),
		utils.String("description", "处理MQTT消息和事件"))

	// 添加TCP监听器（端口1883）
	s.logger.Debug("🌐 正在创建TCP监听器...")
	tcp := listeners.NewTCP(listeners.Config{
		ID:      "tcp1",
		Address: ":1883",
	})
	s.logger.Info("✅ TCP监听器已创建",
		utils.String("listener_id", "tcp1"),
		utils.String("address", ":1883"),
		utils.String("protocol", "TCP"),
		utils.String("port", "1883"))

	if err := s.server.AddListener(tcp); err != nil {
		s.logger.Error("❌ 添加TCP监听器失败",
			utils.String("listener_id", "tcp1"),
			utils.String("address", ":1883"),
			utils.ErrorField(err))
		return fmt.Errorf("添加TCP监听器失败: %w", err)
	}
	s.logger.Info("✅ TCP监听器已添加到服务器",
		utils.String("listener_id", "tcp1"),
		utils.String("address", ":1883"),
		utils.String("status", "ready"))

	// 启动服务器
	s.logger.Info("🚀 正在启动MQTT服务器服务...")
	s.running = true

	// 在goroutine中运行服务器
	go func() {
		s.logger.Info("🔄 MQTT服务器服务正在启动...")
		// Serve()方法是阻塞的，会一直运行直到服务器停止
		// 如果Serve()返回，说明服务器已经停止
		if err := s.server.Serve(); err != nil {
			s.logger.Error("❌ MQTT服务器服务异常停止",
				utils.ErrorField(err),
				utils.String("reason", "server_error"))
		} else {
			s.logger.Info("🛑 MQTT服务器服务正常停止",
				utils.String("reason", "normal_shutdown"))
		}
		// 无论正常停止还是异常停止，都设置running为false
		s.running = false
		s.logger.Info("📊 服务器状态已更新", utils.Bool("running", s.running))
	}()

	// 等待一小段时间确保服务器启动
	s.logger.Debug("⏳ 等待服务器启动...")
	time.Sleep(200 * time.Millisecond)

	// 检查服务器是否真正启动
	if s.server == nil {
		s.logger.Error("❌ MQTT服务器实例为空")
		return fmt.Errorf("MQTT服务器实例为空")
	}

	// 检查监听器是否已添加
	s.logger.Debug("🔍 检查MQTT服务器监听器状态...")

	// 等待running状态被设置
	s.logger.Debug("⏳ 等待服务器运行状态确认...")
	for i := 0; i < 10; i++ {
		if s.running {
			s.logger.Debug("✅ 服务器运行状态已确认",
				utils.Int("attempt", i+1),
				utils.Bool("running", s.running))
			break
		}
		time.Sleep(50 * time.Millisecond)
	}

	s.logger.Info("🎉 MQTT服务器启动成功！",
		utils.String("address", ":1883"),
		utils.String("client_id", s.config.ClientID),
		utils.String("server_type", "mochi-mqtt"),
		utils.Bool("running", s.running),
		utils.String("status", "ready"),
		utils.String("listeners", "TCP:1883"))

	s.logger.Info("📋 MQTT服务器配置摘要",
		utils.String("broker", s.config.Broker),
		utils.Int("keep_alive", s.config.KeepAlive),
		utils.Int("qos", s.config.QoS),
		utils.Bool("clean_session", s.config.CleanSession),
		utils.Bool("auto_reconnect", s.config.AutoReconnect),
		utils.Int("connect_timeout", s.config.ConnectTimeout))

	return nil
}

// Stop 停止MQTT服务器
func (s *Server) Stop() {
	s.logger.Info("🛑 开始停止MQTT服务器...")

	s.cancel()
	s.running = false
	s.logger.Debug("📊 服务器状态已更新", utils.Bool("running", s.running))

	if s.server != nil {
		s.logger.Debug("🔌 正在关闭MQTT服务器连接...")
		s.server.Close()
		s.logger.Info("✅ MQTT服务器连接已关闭")
	} else {
		s.logger.Warn("⚠️ MQTT服务器实例为空，无需关闭")
	}

	s.logger.Info("🎯 MQTT服务器已完全停止",
		utils.String("status", "stopped"),
		utils.Bool("running", s.running))
}

// Publish 发布消息到MQTT服务器
func (s *Server) Publish(topic string, payload interface{}) error {
	if !s.running || s.server == nil {
		s.logger.Error("❌ 无法发布消息：MQTT服务器未运行",
			utils.Bool("running", s.running),
			utils.Bool("server_exists", s.server != nil))
		return fmt.Errorf("MQTT服务器未运行")
	}

	s.logger.Debug("📤 开始发布MQTT消息",
		utils.String("topic", topic),
		utils.String("payload_type", fmt.Sprintf("%T", payload)))

	var data []byte
	var err error

	switch v := payload.(type) {
	case []byte:
		data = v
		s.logger.Debug("📦 使用字节数组作为消息载荷", utils.Int("size", len(data)))
	case string:
		data = []byte(v)
		s.logger.Debug("📝 使用字符串作为消息载荷",
			utils.Int("size", len(data)),
			utils.String("content", v))
	default:
		// 尝试JSON序列化
		s.logger.Debug("🔄 尝试JSON序列化消息载荷", utils.String("type", fmt.Sprintf("%T", payload)))
		data, err = json.Marshal(payload)
		if err != nil {
			s.logger.Error("❌ 序列化消息失败",
				utils.ErrorField(err),
				utils.String("payload_type", fmt.Sprintf("%T", payload)))
			return err
		}
		s.logger.Debug("✅ JSON序列化成功",
			utils.Int("size", len(data)),
			utils.String("json", string(data)))
	}

	// 使用Mochi MQTT的Publish方法
	s.logger.Debug("🚀 正在发布消息到MQTT服务器",
		utils.String("topic", topic),
		utils.Int("payload_size", len(data)),
		utils.Bool("retain", false),
		utils.Int("qos", 1))

	if err := s.server.Publish(topic, data, false, 1); err != nil {
		s.logger.Error("❌ 发布消息失败",
			utils.String("topic", topic),
			utils.Int("payload_size", len(data)),
			utils.ErrorField(err))
		return err
	}

	s.logger.Info("✅ 发布消息成功",
		utils.String("topic", topic),
		utils.Int("payload_size", len(data)),
		utils.String("status", "published"))

	return nil
}

// IsRunning 检查服务器是否运行中
func (s *Server) IsRunning() bool {
	running := s.running && s.server != nil
	s.logger.Debug("🔍 检查服务器运行状态",
		utils.Bool("running", s.running),
		utils.Bool("server_exists", s.server != nil),
		utils.Bool("is_running", running))
	return running
}

// GetStatus 获取服务器状态
func (s *Server) GetStatus() map[string]interface{} {
	s.logger.Debug("📊 获取MQTT服务器状态...")

	status := map[string]interface{}{
		"running":   s.running,
		"connected": s.server != nil,
		"address":   ":1883",
		"client_id": s.config.ClientID,
	}

	if s.server != nil {
		status["server_type"] = "mochi-mqtt"
		status["listeners"] = s.server.Listeners
	}

	s.logger.Debug("📋 服务器状态信息",
		utils.Bool("running", s.running),
		utils.Bool("connected", s.server != nil),
		utils.String("address", ":1883"),
		utils.String("client_id", s.config.ClientID))

	return status
}

// MessageHandlerHook MQTT消息处理钩子
type MessageHandlerHook struct {
	logger            utils.Logger
	sensorDataHandler *SensorDataHandler
}

// ID 返回钩子ID
func (h *MessageHandlerHook) ID() string {
	return "message-handler"
}

// Provides 返回钩子提供的事件
func (h *MessageHandlerHook) Provides(b byte) bool {
	// 使用bytes.Contains方法检查事件是否在支持的事件列表中
	// 这种方式更优雅，避免了硬编码的switch语句
	return bytes.Contains([]byte{
		mqtt.OnConnect,             // 客户端连接
		mqtt.OnDisconnect,          // 客户端断开
		mqtt.OnConnectAuthenticate, // 连接认证
		mqtt.OnACLCheck,            // ACL检查
		mqtt.OnSubscribe,           // 订阅
		mqtt.OnSubscribed,          // 已订阅
		mqtt.OnUnsubscribe,         // 取消订阅
		mqtt.OnUnsubscribed,        // 已取消订阅
		mqtt.OnPublish,             // 发布消息
		mqtt.OnPublished,           // 消息已发布
		mqtt.OnPublishDropped,      // 发布丢弃
		mqtt.OnSysInfoTick,         // 系统信息定时器
		mqtt.OnSessionEstablish,    // 会话建立
		mqtt.OnSessionEstablished,  // 会话已建立
		mqtt.OnQosPublish,          // QoS发布
		mqtt.OnQosComplete,         // QoS完成
		mqtt.OnQosDropped,          // QoS丢弃
		mqtt.OnPacketIDExhausted,   // 包ID耗尽
		mqtt.OnClientExpired,       // 客户端过期
	}, []byte{b})
}

// Init 初始化钩子
func (h *MessageHandlerHook) Init(config interface{}) error {
	if configMap, ok := config.(map[string]interface{}); ok {
		// 初始化logger
		if logger, ok := configMap["logger"]; ok {
			if l, ok := logger.(utils.Logger); ok {
				h.logger = l
			} else {
				return fmt.Errorf("logger类型转换失败")
			}
		} else {
			return fmt.Errorf("未找到logger配置")
		}

		// 初始化数据处理器
		if sensorDataHandler, ok := configMap["sensorDataHandler"]; ok {
			if handler, ok := sensorDataHandler.(*SensorDataHandler); ok {
				h.sensorDataHandler = handler
			} else {
				return fmt.Errorf("sensorDataHandler类型转换失败")
			}
		} else {
			return fmt.Errorf("未找到sensorDataHandler配置")
		}

		h.logger.Info("🔧 MQTT消息处理钩子已初始化",
			utils.String("hook_id", h.ID()),
			utils.String("description", "处理MQTT消息和事件"),
			utils.Bool("has_sensor_handler", h.sensorDataHandler != nil))
	} else {
		return fmt.Errorf("配置类型转换失败")
	}
	return nil
}

// Stop 停止钩子
func (h *MessageHandlerHook) Stop() error {
	return nil
}

// SetOpts 设置选项
func (h *MessageHandlerHook) SetOpts(l *slog.Logger, o *mqtt.HookOptions) {
	// 实现空方法
}

// OnStarted 服务器启动时调用
func (h *MessageHandlerHook) OnStarted() {
	if h.logger != nil {
		h.logger.Info("🎯 MQTT服务器已启动，开始接受客户端连接",
			utils.String("status", "running"),
			utils.String("port", "1883"),
			utils.String("protocol", "TCP"))
	}
}

// OnStopped 服务器停止时调用
func (h *MessageHandlerHook) OnStopped() {
	if h.logger != nil {
		h.logger.Info("🛑 MQTT服务器已停止，不再接受新连接",
			utils.String("status", "stopped"))
	}
}

// OnConnectAuthenticate 连接认证
func (h *MessageHandlerHook) OnConnectAuthenticate(cl *mqtt.Client, pk packets.Packet) bool {
	if h.logger != nil {
		h.logger.Info("🔐 客户端认证请求",
			utils.String("client_id", cl.ID),
			utils.String("remote_addr", cl.Net.Remote),
			utils.String("protocol_version", fmt.Sprintf("MQTT v%d", pk.ProtocolVersion)),
			utils.Bool("clean_session", pk.Connect.Clean),
			utils.Int("keep_alive", int(pk.Connect.Keepalive)))

		// 记录连接包详细信息
		if len(pk.Connect.Username) > 0 {
			h.logger.Info("👤 客户端提供用户名",
				utils.String("client_id", cl.ID),
				utils.String("username", string(pk.Connect.Username)))
		} else {
			h.logger.Info("👤 客户端未提供用户名",
				utils.String("client_id", cl.ID))
		}
		if len(pk.Connect.Password) > 0 {
			h.logger.Info("🔑 客户端提供密码",
				utils.String("client_id", cl.ID),
				utils.Bool("has_password", true))
		} else {
			h.logger.Info("🔑 客户端未提供密码",
				utils.String("client_id", cl.ID))
		}
	}

	// 允许所有连接
	result := true
	if h.logger != nil {
		h.logger.Info("✅ 客户端认证通过",
			utils.String("client_id", cl.ID),
			utils.Bool("authenticated", result),
			utils.String("reason", "AllowHook - 允许所有连接"))
	}
	return result
}

// OnACLCheck ACL检查（允许所有访问）
func (h *MessageHandlerHook) OnACLCheck(cl *mqtt.Client, topic string, write bool) bool {
	if h.logger != nil {
		h.logger.Debug("MQTT ACL检查",
			utils.String("client_id", cl.ID),
			utils.String("topic", topic),
			utils.Bool("write", write),
			utils.String("action", map[bool]string{true: "发布", false: "订阅"}[write]))
	}
	return true
}

// OnSysInfoTick 系统信息更新
func (h *MessageHandlerHook) OnSysInfoTick(info *system.Info) {
	if h.logger != nil {
		// 每60秒记录一次系统信息（避免日志过多）
		if info.Uptime%60 == 0 {
			h.logger.Info("MQTT服务器系统信息",
				utils.Int64("uptime_seconds", info.Uptime),
				utils.Int64("connected_clients", info.ClientsConnected),
				utils.Int64("total_clients", info.ClientsTotal),
				utils.Int64("subscriptions", info.Subscriptions),
				utils.Int64("retained_messages", info.Retained),
				utils.Int64("messages_received", info.MessagesReceived),
				utils.Int64("messages_sent", info.MessagesSent),
				utils.Int64("bytes_received", info.BytesReceived),
				utils.Int64("bytes_sent", info.BytesSent))
		}
	}
}

// OnConnect 处理客户端连接
func (h *MessageHandlerHook) OnConnect(cl *mqtt.Client, pk packets.Packet) error {
	if h.logger != nil {
		h.logger.Info("🔗 客户端连接成功！",
			utils.String("client_id", cl.ID),
			utils.String("remote_addr", cl.Net.Remote),
			utils.String("protocol_version", fmt.Sprintf("MQTT v%d", pk.ProtocolVersion)),
			utils.Bool("clean_session", pk.Connect.Clean),
			utils.Int("keep_alive", int(pk.Connect.Keepalive)),
			utils.String("will_topic", pk.Connect.WillTopic),
			utils.Bool("has_will", pk.Connect.WillTopic != ""))

		// 记录客户端详细信息
		h.logger.Info("📋 客户端连接详细信息",
			utils.String("client_id", cl.ID),
			utils.String("remote_addr", cl.Net.Remote),
			utils.String("username", string(pk.Connect.Username)),
			utils.String("session_type", map[bool]string{true: "新会话", false: "持久会话"}[pk.Connect.Clean]))
	}
	return nil
}

// OnSessionEstablish 会话建立
func (h *MessageHandlerHook) OnSessionEstablish(cl *mqtt.Client, pk packets.Packet) {
	if h.logger != nil {
		h.logger.Info("🤝 会话建立中",
			utils.String("client_id", cl.ID),
			utils.Bool("clean_session", pk.Connect.Clean),
			utils.String("session_type", map[bool]string{true: "新会话", false: "持久会话"}[pk.Connect.Clean]))
	}
}

// OnSessionEstablished 会话已建立
func (h *MessageHandlerHook) OnSessionEstablished(cl *mqtt.Client, pk packets.Packet) {
	if h.logger != nil {
		h.logger.Info("✅ 会话已建立，客户端就绪",
			utils.String("client_id", cl.ID),
			utils.String("status", "ready"),
			utils.String("message", "可以开始发送和接收消息"))
	}
}

// OnDisconnect 处理客户端断开连接
func (h *MessageHandlerHook) OnDisconnect(cl *mqtt.Client, err error, expire bool) {
	if h.logger != nil {
		if err != nil {
			h.logger.Error("❌ 客户端异常断开连接",
				utils.String("client_id", cl.ID),
				utils.String("remote_addr", cl.Net.Remote),
				utils.Bool("session_expired", expire),
				utils.String("disconnect_reason", "网络错误或协议错误"),
				utils.ErrorField(err))
		} else {
			h.logger.Info("🔌 客户端正常断开连接",
				utils.String("client_id", cl.ID),
				utils.String("remote_addr", cl.Net.Remote),
				utils.Bool("session_expired", expire),
				utils.String("disconnect_reason", "客户端主动断开"))
		}

		// 记录断开连接的详细信息
		h.logger.Info("📊 客户端断开连接详细信息",
			utils.String("client_id", cl.ID),
			utils.String("remote_addr", cl.Net.Remote),
			utils.Bool("session_expired", expire),
			utils.String("error_type", fmt.Sprintf("%T", err)))
	}
}

// OnAuthPacket 认证包处理
func (h *MessageHandlerHook) OnAuthPacket(cl *mqtt.Client, pk packets.Packet) (packets.Packet, error) {
	return pk, nil
}

// OnPacketRead 包读取
func (h *MessageHandlerHook) OnPacketRead(cl *mqtt.Client, pk packets.Packet) (packets.Packet, error) {
	return pk, nil
}

// OnPacketEncode 包编码
func (h *MessageHandlerHook) OnPacketEncode(cl *mqtt.Client, pk packets.Packet) packets.Packet {
	return pk
}

// OnPacketSent 包发送
func (h *MessageHandlerHook) OnPacketSent(cl *mqtt.Client, pk packets.Packet, b []byte) {
	// 实现空方法
}

// OnPacketProcessed 包处理完成
func (h *MessageHandlerHook) OnPacketProcessed(cl *mqtt.Client, pk packets.Packet, err error) {
	// 实现空方法
}

// OnSubscribe 订阅
func (h *MessageHandlerHook) OnSubscribe(cl *mqtt.Client, pk packets.Packet) packets.Packet {
	if h.logger != nil {
		// 记录订阅的主题列表
		topics := make([]string, len(pk.Filters))
		for i, filter := range pk.Filters {
			topics[i] = filter.Filter
		}

		h.logger.Info("MQTT客户端订阅请求",
			utils.String("client_id", cl.ID),
			utils.String("topics", fmt.Sprintf("%v", topics)),
			utils.Int("topic_count", len(pk.Filters)))

		// 记录每个主题的QoS级别
		for _, filter := range pk.Filters {
			h.logger.Debug("订阅主题详情",
				utils.String("client_id", cl.ID),
				utils.String("topic", filter.Filter),
				utils.Int("qos", int(filter.Qos)))
		}
	}
	return pk
}

// OnSubscribed 已订阅
func (h *MessageHandlerHook) OnSubscribed(cl *mqtt.Client, pk packets.Packet, reasonCodes []byte) {
	if h.logger != nil {
		// 记录订阅结果
		topics := make([]string, len(pk.Filters))
		for i, filter := range pk.Filters {
			topics[i] = filter.Filter
		}

		h.logger.Info("MQTT客户端订阅成功",
			utils.String("client_id", cl.ID),
			utils.String("topics", fmt.Sprintf("%v", topics)),
			utils.Int("topic_count", len(pk.Filters)))

		// 记录每个主题的订阅结果
		for i, filter := range pk.Filters {
			reasonCode := "成功"
			if i < len(reasonCodes) {
				switch reasonCodes[i] {
				case 0x00:
					reasonCode = "成功"
				case 0x80:
					reasonCode = "失败"
				case 0x01, 0x02:
					reasonCode = fmt.Sprintf("QoS %d", reasonCodes[i])
				default:
					reasonCode = fmt.Sprintf("未知(%d)", reasonCodes[i])
				}
			}

			h.logger.Debug("订阅结果详情",
				utils.String("client_id", cl.ID),
				utils.String("topic", filter.Filter),
				utils.String("result", reasonCode))
		}
	}
}

// OnSelectSubscribers 选择订阅者
func (h *MessageHandlerHook) OnSelectSubscribers(subs *mqtt.Subscribers, pk packets.Packet) *mqtt.Subscribers {
	return subs
}

// OnUnsubscribe 取消订阅
func (h *MessageHandlerHook) OnUnsubscribe(cl *mqtt.Client, pk packets.Packet) packets.Packet {
	return pk
}

// OnUnsubscribed 已取消订阅
func (h *MessageHandlerHook) OnUnsubscribed(cl *mqtt.Client, pk packets.Packet) {
	// 实现空方法
}

// OnPublish 发布
func (h *MessageHandlerHook) OnPublish(cl *mqtt.Client, pk packets.Packet) (packets.Packet, error) {
	if h.logger != nil {
		// 记录发布请求基本信息 - 这是最重要的日志
		h.logger.Info("🎯 收到客户端消息！",
			utils.String("client_id", cl.ID),
			utils.String("topic", pk.TopicName),
			utils.Int("qos", int(pk.FixedHeader.Qos)),
			utils.Bool("retain", pk.FixedHeader.Retain),
			utils.Bool("duplicate", pk.FixedHeader.Dup),
			utils.Int("payload_size", len(pk.Payload)),
			utils.String("remote_addr", cl.Net.Remote))

		// 记录消息内容预览
		if len(pk.Payload) > 0 {
			payloadStr := string(pk.Payload)

			// 尝试解析JSON格式的消息
			var jsonData map[string]interface{}
			if err := json.Unmarshal(pk.Payload, &jsonData); err == nil {
				// JSON格式消息 - 提取关键字段
				h.logger.Info("📋 消息内容 (JSON格式)",
					utils.String("client_id", cl.ID),
					utils.String("topic", pk.TopicName),
					utils.String("json_data", payloadStr))

				// 尝试提取设备ID和时间戳
				if deviceID, ok := jsonData["device_id"].(string); ok {
					h.logger.Info("🏷️ 设备信息",
						utils.String("client_id", cl.ID),
						utils.String("device_id", deviceID),
						utils.String("topic", pk.TopicName))
				}
				if timestamp, ok := jsonData["timestamp"]; ok {
					h.logger.Info("⏰ 时间戳信息",
						utils.String("client_id", cl.ID),
						utils.String("timestamp", fmt.Sprintf("%v", timestamp)),
						utils.String("topic", pk.TopicName))
				}
				// 尝试提取甲醛浓度数据
				if formaldehyde, ok := jsonData["formaldehyde"]; ok {
					h.logger.Info("🌡️ 甲醛浓度数据",
						utils.String("client_id", cl.ID),
						utils.String("formaldehyde", fmt.Sprintf("%v", formaldehyde)),
						utils.String("topic", pk.TopicName))
				}
			} else {
				// 非JSON格式消息
				if len(pk.Payload) <= 200 {
					h.logger.Info("📋 发布消息内容预览 (文本格式)",
						utils.String("client_id", cl.ID),
						utils.String("topic", pk.TopicName),
						utils.String("payload_preview", payloadStr))
				} else {
					h.logger.Info("📋 发布消息内容预览 (文本格式-截断)",
						utils.String("client_id", cl.ID),
						utils.String("topic", pk.TopicName),
						utils.String("payload_preview", payloadStr[:100]+"..."),
						utils.Int("total_size", len(pk.Payload)))
				}
			}
		} else {
			h.logger.Info("📋 发布消息内容预览 (空消息)",
				utils.String("client_id", cl.ID),
				utils.String("topic", pk.TopicName))
		}

		// 记录QoS相关信息
		if pk.FixedHeader.Qos > 0 {
			h.logger.Debug("🔢 QoS消息详情",
				utils.String("client_id", cl.ID),
				utils.String("topic", pk.TopicName),
				utils.Int("packet_id", int(pk.PacketID)),
				utils.Int("qos", int(pk.FixedHeader.Qos)))
		}
	}

	// 处理传感器数据消息
	if h.sensorDataHandler != nil {
		// 检查是否是传感器数据主题
		if isSensorDataTopic(pk.TopicName) {
			h.logger.Info("🔧 开始处理传感器数据",
				utils.String("client_id", cl.ID),
				utils.String("topic", pk.TopicName))

			// 调用数据处理器处理消息
			if err := h.sensorDataHandler.HandleMessage(pk.TopicName, pk.Payload); err != nil {
				h.logger.Error("❌ 处理传感器数据失败",
					utils.String("client_id", cl.ID),
					utils.String("topic", pk.TopicName),
					utils.ErrorField(err))
			} else {
				h.logger.Info("✅ 传感器数据处理成功",
					utils.String("client_id", cl.ID),
					utils.String("topic", pk.TopicName))
			}
		}
	} else {
		h.logger.Warn("⚠️ 数据处理器未初始化，跳过数据处理",
			utils.String("client_id", cl.ID),
			utils.String("topic", pk.TopicName))
	}

	return pk, nil
}

// OnPublished 已发布
func (h *MessageHandlerHook) OnPublished(cl *mqtt.Client, pk packets.Packet) {
	if h.logger != nil {
		// 记录消息基本信息
		h.logger.Info("✅ 消息处理完成",
			utils.String("client_id", cl.ID),
			utils.String("topic", pk.TopicName),
			utils.Int("qos", int(pk.FixedHeader.Qos)),
			utils.Bool("retain", pk.FixedHeader.Retain),
			utils.Bool("duplicate", pk.FixedHeader.Dup),
			utils.Int("payload_size", len(pk.Payload)),
			utils.Int("packet_id", int(pk.PacketID)),
			utils.String("status", "processed"))

		// 记录消息内容
		if len(pk.Payload) > 0 {
			payloadStr := string(pk.Payload)

			// 尝试解析JSON格式的消息
			var jsonData map[string]interface{}
			if err := json.Unmarshal(pk.Payload, &jsonData); err == nil {
				// JSON格式消息
				h.logger.Info("📋 消息内容 (JSON格式)",
					utils.String("client_id", cl.ID),
					utils.String("topic", pk.TopicName),
					utils.String("json_data", payloadStr))
			} else {
				// 非JSON格式消息
				if len(pk.Payload) <= 500 {
					// 小消息直接显示
					h.logger.Info("📋 消息内容 (文本格式)",
						utils.String("client_id", cl.ID),
						utils.String("topic", pk.TopicName),
						utils.String("payload", payloadStr))
				} else {
					// 大消息截断显示
					h.logger.Info("📋 消息内容 (文本格式-截断)",
						utils.String("client_id", cl.ID),
						utils.String("topic", pk.TopicName),
						utils.String("payload_preview", payloadStr[:200]+"..."),
						utils.Int("total_size", len(pk.Payload)))
				}
			}
		} else {
			h.logger.Info("📋 消息内容 (空消息)",
				utils.String("client_id", cl.ID),
				utils.String("topic", pk.TopicName))
		}

		// 记录消息来源信息
		h.logger.Debug("🔍 消息来源信息",
			utils.String("client_id", cl.ID),
			utils.String("remote_addr", cl.Net.Remote),
			utils.String("protocol_version", fmt.Sprintf("%d", cl.Properties.ProtocolVersion)),
			utils.String("topic", pk.TopicName))
	}
}

// OnPublishDropped 发布丢弃
func (h *MessageHandlerHook) OnPublishDropped(cl *mqtt.Client, pk packets.Packet) {
	if h.logger != nil {
		h.logger.Warn("MQTT消息发布被丢弃",
			utils.String("client_id", cl.ID),
			utils.String("topic", pk.TopicName),
			utils.Int("qos", int(pk.FixedHeader.Qos)),
			utils.String("reason", "可能是客户端断开或QoS处理失败"))
	}
}

// OnRetainMessage 保留消息
func (h *MessageHandlerHook) OnRetainMessage(cl *mqtt.Client, pk packets.Packet, r int64) {
	// 实现空方法
}

// OnRetainPublished 保留消息已发布
func (h *MessageHandlerHook) OnRetainPublished(cl *mqtt.Client, pk packets.Packet) {
	// 实现空方法
}

// OnQosPublish QoS发布
func (h *MessageHandlerHook) OnQosPublish(cl *mqtt.Client, pk packets.Packet, sent int64, resends int) {
	// 实现空方法
}

// OnQosComplete QoS完成
func (h *MessageHandlerHook) OnQosComplete(cl *mqtt.Client, pk packets.Packet) {
	// 实现空方法
}

// OnQosDropped QoS丢弃
func (h *MessageHandlerHook) OnQosDropped(cl *mqtt.Client, pk packets.Packet) {
	if h.logger != nil {
		h.logger.Warn("MQTT QoS消息被丢弃",
			utils.String("client_id", cl.ID),
			utils.String("topic", pk.TopicName),
			utils.Int("qos", int(pk.FixedHeader.Qos)),
			utils.Int("packet_id", int(pk.PacketID)),
			utils.String("reason", "可能是重试次数超限或客户端断开"))
	}
}

// OnPacketIDExhausted 包ID耗尽
func (h *MessageHandlerHook) OnPacketIDExhausted(cl *mqtt.Client, pk packets.Packet) {
	if h.logger != nil {
		h.logger.Error("MQTT包ID已耗尽",
			utils.String("client_id", cl.ID),
			utils.String("topic", pk.TopicName),
			utils.String("reason", "客户端有太多未确认的QoS消息，可能需要检查网络连接或增加包ID范围"))
	}
}

// OnWill 遗嘱
func (h *MessageHandlerHook) OnWill(cl *mqtt.Client, will mqtt.Will) (mqtt.Will, error) {
	return will, nil
}

// OnWillSent 遗嘱已发送
func (h *MessageHandlerHook) OnWillSent(cl *mqtt.Client, pk packets.Packet) {
	// 实现空方法
}

// OnClientExpired 客户端过期
func (h *MessageHandlerHook) OnClientExpired(cl *mqtt.Client) {
	if h.logger != nil {
		h.logger.Warn("MQTT客户端会话过期",
			utils.String("client_id", cl.ID),
			utils.String("remote_addr", cl.Net.Remote),
			utils.String("reason", "客户端长时间未连接，会话已清理"))
	}
}

// OnRetainedExpired 保留消息过期
func (h *MessageHandlerHook) OnRetainedExpired(filter string) {
	// 实现空方法
}

// StoredClients 存储的客户端
func (h *MessageHandlerHook) StoredClients() ([]storage.Client, error) {
	return nil, nil
}

// StoredSubscriptions 存储的订阅
func (h *MessageHandlerHook) StoredSubscriptions() ([]storage.Subscription, error) {
	return nil, nil
}

// StoredInflightMessages 存储的飞行消息
func (h *MessageHandlerHook) StoredInflightMessages() ([]storage.Message, error) {
	return nil, nil
}

// StoredRetainedMessages 存储的保留消息
func (h *MessageHandlerHook) StoredRetainedMessages() ([]storage.Message, error) {
	return nil, nil
}

// StoredSysInfo 存储的系统信息
func (h *MessageHandlerHook) StoredSysInfo() (storage.SystemInfo, error) {
	return storage.SystemInfo{}, nil
}

// isSensorDataTopic 判断是否是传感器数据主题
func isSensorDataTopic(topic string) bool {
	// 检查主题格式: air-quality/hcho/{device_id}/data
	// 或者: air-quality/{device_type}/{device_id}/data
	parts := strings.Split(topic, "/")
	if len(parts) != 4 {
		return false
	}

	// 检查前缀
	if parts[0] != "air-quality" {
		return false
	}

	// 检查后缀
	if parts[3] != "data" {
		return false
	}

	// 检查设备类型（第二部分）
	deviceType := parts[1]
	if deviceType != "hcho" && deviceType != "esp32" && deviceType != "sensor" {
		return false
	}

	// 检查设备ID（第三部分）
	deviceID := parts[2]
	return deviceID != ""
}
