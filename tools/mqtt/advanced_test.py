#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
高级MQTT测试程序
支持配置下发、命令控制等高级功能测试
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


class AdvancedMQTTTester:
    """高级MQTT测试器"""
    
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
        self.received_messages = []
        
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
            message_info = {
                'topic': msg.topic,
                'payload': payload,
                'timestamp': datetime.now().isoformat()
            }
            self.received_messages.append(message_info)
            
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


class DeviceSimulator:
    """设备模拟器"""
    
    def __init__(self, device_id: str, tester: AdvancedMQTTTester):
        self.device_id = device_id
        self.tester = tester
        self.running = False
        self.thread = None
        self.config = {
            "report_interval": 30,
            "formaldehyde_warning": 0.08,
            "formaldehyde_critical": 0.1,
            "calibration_enabled": True,
            "calibration_interval": 86400
        }
        
    def generate_sensor_data(self) -> Dict[str, Any]:
        """生成传感器数据"""
        formaldehyde = round(random.uniform(0.01, 0.15), 3)
        temperature = round(random.uniform(18.0, 30.0), 1)
        humidity = round(random.uniform(40.0, 80.0), 1)
        battery = random.randint(20, 100)
        signal_strength = random.randint(-90, -30)
        
        return {
            "device_id": self.device_id,
            "device_type": "hcho",
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
        
    def generate_response(self, command: str) -> Dict[str, Any]:
        """生成命令响应"""
        responses = {
            "calibrate": {
                "status": "success",
                "message": "校准完成",
                "calibration_time": int(time.time()),
                "calibration_duration": 300
            },
            "reboot": {
                "status": "success",
                "message": "设备重启中",
                "reboot_time": int(time.time())
            },
            "get_config": {
                "status": "success",
                "config": self.config
            },
            "update_config": {
                "status": "success",
                "message": "配置更新成功",
                "updated_at": int(time.time())
            }
        }
        
        return {
            "device_id": self.device_id,
            "timestamp": int(time.time()),
            "command": command,
            "response": responses.get(command, {
                "status": "error",
                "message": f"未知命令: {command}"
            })
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
                self.tester.publish(data_topic, sensor_data)
                
                # 随机发布状态信息
                if random.random() < 0.3:
                    status_topic = f"air-quality/hcho/{self.device_id}/status"
                    status_data = self.generate_status_data()
                    self.tester.publish(status_topic, status_data)
                
                time.sleep(interval)
                
            except Exception as e:
                print(f"❌ 模拟设备 {self.device_id} 出错: {e}")
                time.sleep(interval)


def test_config_publish(tester: AdvancedMQTTTester, device_id: str):
    """测试配置下发"""
    print(f"\n🔧 测试配置下发到设备 {device_id}")
    
    config_data = {
        "device_id": device_id,
        "timestamp": int(time.time()),
        "config": {
            "report_interval": 60,
            "thresholds": {
                "formaldehyde_warning": 0.08,
                "formaldehyde_critical": 0.1
            },
            "calibration": {
                "enabled": True,
                "interval": 86400
            }
        }
    }
    
    config_topic = f"air-quality/hcho/{device_id}/config"
    return tester.publish(config_topic, config_data)


def test_command_publish(tester: AdvancedMQTTTester, device_id: str, command: str):
    """测试命令下发"""
    print(f"\n⚡ 测试命令下发到设备 {device_id}: {command}")
    
    command_data = {
        "device_id": device_id,
        "timestamp": int(time.time()),
        "command": {
            "action": command,
            "parameters": {
                "duration": 300 if command == "calibrate" else None
            }
        }
    }
    
    command_topic = f"air-quality/hcho/{device_id}/command"
    return tester.publish(command_topic, command_data)


def test_device_response(tester: AdvancedMQTTTester, device_id: str, command: str):
    """测试设备响应"""
    print(f"\n📨 模拟设备 {device_id} 响应命令: {command}")
    
    simulator = DeviceSimulator(device_id, tester)
    response_data = simulator.generate_response(command)
    
    response_topic = f"air-quality/hcho/{device_id}/response"
    return tester.publish(response_topic, response_data)


def main():
    """主函数"""
    parser = argparse.ArgumentParser(description="高级MQTT测试程序")
    parser.add_argument("--host", default="localhost", help="MQTT Broker地址")
    parser.add_argument("--port", type=int, default=1883, help="MQTT Broker端口")
    parser.add_argument("--username", default="admin", help="MQTT用户名")
    parser.add_argument("--password", default="password", help="MQTT密码")
    parser.add_argument("--device-id", default="hcho_001", help="设备ID")
    parser.add_argument("--test-type", choices=["config", "command", "response", "all"], 
                       default="all", help="测试类型")
    parser.add_argument("--command", choices=["calibrate", "reboot", "get_config", "update_config"],
                       default="calibrate", help="命令类型")
    
    args = parser.parse_args()
    
    # 创建MQTT测试器
    tester = AdvancedMQTTTester(args.host, args.port, args.username, args.password)
    
    try:
        # 连接到MQTT Broker
        if not tester.connect():
            return
            
        # 订阅响应主题
        response_topic = f"air-quality/hcho/{args.device_id}/response"
        tester.subscribe(response_topic)
        
        print(f"\n🧪 开始高级测试 - 设备ID: {args.device_id}")
        
        if args.test_type in ["config", "all"]:
            # 测试配置下发
            test_config_publish(tester, args.device_id)
            time.sleep(2)
            
        if args.test_type in ["command", "all"]:
            # 测试命令下发
            test_command_publish(tester, args.device_id, args.command)
            time.sleep(2)
            
        if args.test_type in ["response", "all"]:
            # 测试设备响应
            test_device_response(tester, args.device_id, args.command)
            time.sleep(2)
            
        # 等待响应消息
        print(f"\n⏳ 等待响应消息...")
        time.sleep(5)
        
        # 显示接收到的消息
        if tester.received_messages:
            print(f"\n📋 接收到 {len(tester.received_messages)} 条消息:")
            for i, msg in enumerate(tester.received_messages, 1):
                print(f"\n消息 {i}:")
                print(f"  主题: {msg['topic']}")
                print(f"  时间: {msg['timestamp']}")
                print(f"  内容: {json.dumps(msg['payload'], indent=2, ensure_ascii=False)}")
        else:
            print("\n⚠️  未接收到任何响应消息")
            
        print(f"\n✅ 测试完成")
            
    except KeyboardInterrupt:
        print("\n🛑 程序被用户中断")
    except Exception as e:
        print(f"\n❌ 程序出错: {e}")
    finally:
        tester.disconnect()


if __name__ == "__main__":
    main()
