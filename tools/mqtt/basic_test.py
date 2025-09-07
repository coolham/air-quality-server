#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
MQTT测试程序
用于测试空气质量监测系统的MQTT功能

支持的功能：
1. 发布甲醛传感器数据
2. 发布设备状态信息
3. 订阅设备响应
4. 模拟多个设备
5. 数据验证和错误处理
"""

import json
import time
import random
import argparse
import threading
from datetime import datetime
from typing import Dict, Any, Optional

try:
    import paho.mqtt.client as mqtt
except ImportError:
    print("请安装paho-mqtt库: pip install paho-mqtt")
    exit(1)


class MQTTTester:
    """MQTT测试器"""
    
    def __init__(self, broker_host: str = "localhost", broker_port: int = 1883,
                 username: str = "admin", password: str = "password"):
        self.broker_host = broker_host
        self.broker_port = broker_port
        self.username = username
        self.password = password
        
        # 创建MQTT客户端
        self.client = mqtt.Client()
        self.client.username_pw_set(username, password)
        
        # 设置回调函数
        self.client.on_connect = self.on_connect
        self.client.on_disconnect = self.on_disconnect
        self.client.on_message = self.on_message
        self.client.on_publish = self.on_publish
        
        # 连接状态
        self.connected = False
        self.message_count = 0
        
    def on_connect(self, client, userdata, flags, rc):
        """连接回调"""
        if rc == 0:
            self.connected = True
            print(f"✅ 成功连接到MQTT Broker: {self.broker_host}:{self.broker_port}")
        else:
            print(f"❌ 连接失败，错误代码: {rc}")
            
    def on_disconnect(self, client, userdata, rc):
        """断开连接回调"""
        self.connected = False
        print(f"🔌 与MQTT Broker断开连接，错误代码: {rc}")
        
    def on_message(self, client, userdata, msg):
        """消息接收回调"""
        try:
            payload = json.loads(msg.payload.decode('utf-8'))
            print(f"📨 收到消息 - 主题: {msg.topic}")
            print(f"   内容: {json.dumps(payload, indent=2, ensure_ascii=False)}")
            self.message_count += 1
        except Exception as e:
            print(f"❌ 解析消息失败: {e}")
            
    def on_publish(self, client, userdata, mid):
        """发布消息回调"""
        print(f"📤 消息发布成功 (ID: {mid})")
        
    def connect(self):
        """连接到MQTT Broker"""
        try:
            print(f"🔗 正在连接到 {self.broker_host}:{self.broker_port}...")
            self.client.connect(self.broker_host, self.broker_port, 60)
            self.client.loop_start()
            
            # 等待连接建立
            timeout = 10
            while not self.connected and timeout > 0:
                time.sleep(0.1)
                timeout -= 0.1
                
            if not self.connected:
                raise Exception("连接超时")
                
        except Exception as e:
            print(f"❌ 连接失败: {e}")
            return False
        return True
        
    def disconnect(self):
        """断开连接"""
        if self.connected:
            self.client.loop_stop()
            self.client.disconnect()
            print("🔌 已断开连接")
            
    def subscribe(self, topic: str):
        """订阅主题"""
        if not self.connected:
            print("❌ 未连接到MQTT Broker")
            return False
            
        result = self.client.subscribe(topic)
        if result[0] == mqtt.MQTT_ERR_SUCCESS:
            print(f"📡 成功订阅主题: {topic}")
            return True
        else:
            print(f"❌ 订阅失败: {result}")
            return False
            
    def publish(self, topic: str, payload: Dict[str, Any], qos: int = 1):
        """发布消息"""
        if not self.connected:
            print("❌ 未连接到MQTT Broker")
            return False
            
        try:
            message = json.dumps(payload, ensure_ascii=False)
            result = self.client.publish(topic, message, qos=qos)
            
            if result.rc == mqtt.MQTT_ERR_SUCCESS:
                print(f"📤 发布消息到主题: {topic}")
                return True
            else:
                print(f"❌ 发布失败: {result}")
                return False
                
        except Exception as e:
            print(f"❌ 发布消息时出错: {e}")
            return False


class FormaldehydeSensorSimulator:
    """甲醛传感器模拟器"""
    
    def __init__(self, device_id: str, mqtt_tester: MQTTTester):
        self.device_id = device_id
        self.mqtt_tester = mqtt_tester
        self.running = False
        self.thread = None
        
    def generate_sensor_data(self) -> Dict[str, Any]:
        """生成传感器数据"""
        # 模拟甲醛浓度 (mg/m³)
        formaldehyde = round(random.uniform(0.01, 0.15), 3)
        
        # 模拟温度 (°C)
        temperature = round(random.uniform(18.0, 30.0), 1)
        
        # 模拟湿度 (%)
        humidity = round(random.uniform(40.0, 80.0), 1)
        
        # 模拟电池电量 (%)
        battery = random.randint(20, 100)
        
        # 模拟信号强度 (dBm)
        signal_strength = random.randint(-90, -30)
        
        return {
            "device_id": self.device_id,
            "device_type": "hcho",
            "sensor_id": f"sensor_{self.device_id}_01",
            "sensor_type": "hcho",
            "timestamp": int(time.time()),
            "data": {
                "formaldehyde": formaldehyde,
                "temperature": temperature,
                "humidity": humidity,
                "battery": battery
            },
            "location": {
                "latitude": 39.9042 + random.uniform(-0.01, 0.01),
                "longitude": 116.4074 + random.uniform(-0.01, 0.01),
                "address": "北京市朝阳区测试位置"
            },
            "quality": {
                "signal_strength": signal_strength,
                "data_quality": "good" if formaldehyde < 0.1 else "poor"
            }
        }
        
    def generate_status_data(self) -> Dict[str, Any]:
        """生成设备状态数据"""
        return {
            "device_id": self.device_id,
            "device_type": "hcho",
            "timestamp": int(time.time()),
            "status": {
                "online": True,
                "battery_level": random.randint(20, 100),
                "signal_strength": random.randint(-90, -30),
                "last_data_time": int(time.time()) - random.randint(0, 60),
                "error_code": 0,
                "error_message": ""
            },
            "firmware": {
                "version": "1.2.3",
                "build_date": "2024-01-15"
            }
        }
        
    def start_simulation(self, interval: int = 30):
        """开始模拟"""
        if self.running:
            print(f"⚠️  设备 {self.device_id} 已在运行")
            return
            
        self.running = True
        self.thread = threading.Thread(target=self._simulation_loop, args=(interval,))
        self.thread.daemon = True
        self.thread.start()
        print(f"🚀 开始模拟设备 {self.device_id}，数据上报间隔: {interval}秒")
        
    def stop_simulation(self):
        """停止模拟"""
        self.running = False
        if self.thread:
            self.thread.join()
        print(f"🛑 停止模拟设备 {self.device_id}")
        
    def _simulation_loop(self, interval: int):
        """模拟循环"""
        while self.running:
            try:
                # 发布传感器数据
                data_topic = f"air-quality/hcho/{self.device_id}/data"
                sensor_data = self.generate_sensor_data()
                self.mqtt_tester.publish(data_topic, sensor_data)
                
                # 随机发布状态信息
                if random.random() < 0.3:  # 30%概率发布状态
                    status_topic = f"air-quality/hcho/{self.device_id}/status"
                    status_data = self.generate_status_data()
                    self.mqtt_tester.publish(status_topic, status_data)
                
                time.sleep(interval)
                
            except Exception as e:
                print(f"❌ 模拟设备 {self.device_id} 出错: {e}")
                time.sleep(interval)


def main():
    """主函数"""
    parser = argparse.ArgumentParser(description="MQTT测试程序")
    parser.add_argument("--host", default="localhost", help="MQTT Broker地址")
    parser.add_argument("--port", type=int, default=1883, help="MQTT Broker端口")
    parser.add_argument("--username", default="admin", help="MQTT用户名")
    parser.add_argument("--password", default="password", help="MQTT密码")
    parser.add_argument("--device-id", default="hcho_001", help="设备ID")
    parser.add_argument("--interval", type=int, default=30, help="数据上报间隔(秒)")
    parser.add_argument("--count", type=int, default=10, help="发送消息数量")
    parser.add_argument("--subscribe", action="store_true", help="订阅响应主题")
    parser.add_argument("--simulate", action="store_true", help="持续模拟模式")
    
    args = parser.parse_args()
    
    # 创建MQTT测试器
    tester = MQTTTester(args.host, args.port, args.username, args.password)
    
    try:
        # 连接到MQTT Broker
        if not tester.connect():
            return
            
        # 订阅响应主题
        if args.subscribe:
            response_topic = f"air-quality/hcho/{args.device_id}/response"
            tester.subscribe(response_topic)
            
        if args.simulate:
            # 持续模拟模式
            simulator = FormaldehydeSensorSimulator(args.device_id, tester)
            simulator.start_simulation(args.interval)
            
            print("\n📋 模拟模式运行中...")
            print("按 Ctrl+C 停止模拟")
            
            try:
                while True:
                    time.sleep(1)
            except KeyboardInterrupt:
                print("\n🛑 收到停止信号")
                simulator.stop_simulation()
        else:
            # 单次测试模式
            print(f"\n🧪 开始单次测试，发送 {args.count} 条消息...")
            
            for i in range(args.count):
                # 生成测试数据
                simulator = FormaldehydeSensorSimulator(args.device_id, tester)
                sensor_data = simulator.generate_sensor_data()
                
                # 发布数据
                data_topic = f"air-quality/hcho/{args.device_id}/data"
                tester.publish(data_topic, sensor_data)
                
                # 随机发布状态
                if random.random() < 0.5:
                    status_data = simulator.generate_status_data()
                    status_topic = f"air-quality/hcho/{args.device_id}/status"
                    tester.publish(status_topic, status_data)
                
                time.sleep(1)
                
            print(f"\n✅ 测试完成，共发送 {args.count} 条消息")
            
    except KeyboardInterrupt:
        print("\n🛑 程序被用户中断")
    except Exception as e:
        print(f"\n❌ 程序出错: {e}")
    finally:
        tester.disconnect()


if __name__ == "__main__":
    main()
