#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
配置驱动的MQTT测试程序
根据配置文件执行不同的测试场景
"""

import json
import time
import random
import argparse
import threading
from datetime import datetime
from typing import Dict, Any, List

try:
    import paho.mqtt.client as mqtt
except ImportError:
    print("请安装paho-mqtt库: pip install paho-mqtt")
    exit(1)


class ConfigDrivenTester:
    """配置驱动的测试器"""
    
    def __init__(self, config_file: str = "test_config.json"):
        self.config_file = config_file
        self.config = self.load_config()
        self.clients = {}
        self.running = False
        
    def load_config(self) -> Dict[str, Any]:
        """加载配置文件"""
        try:
            with open(self.config_file, 'r', encoding='utf-8') as f:
                return json.load(f)
        except FileNotFoundError:
            print(f"❌ 配置文件 {self.config_file} 不存在")
            exit(1)
        except json.JSONDecodeError as e:
            print(f"❌ 配置文件格式错误: {e}")
            exit(1)
            
    def create_mqtt_client(self, broker_config: Dict[str, Any]) -> mqtt.Client:
        """创建MQTT客户端"""
        client = mqtt.Client()
        client.username_pw_set(broker_config["username"], broker_config["password"])
        
        def on_connect(client, userdata, flags, rc):
            if rc == 0:
                print(f"✅ 连接到 {broker_config['name']}: {broker_config['host']}:{broker_config['port']}")
            else:
                print(f"❌ 连接失败: {broker_config['name']}, 错误代码: {rc}")
                
        def on_disconnect(client, userdata, rc):
            print(f"🔌 断开连接: {broker_config['name']}")
            
        def on_message(client, userdata, msg):
            try:
                payload = json.loads(msg.payload.decode('utf-8'))
                print(f"📨 收到消息 - Broker: {broker_config['name']}, 主题: {msg.topic}")
            except Exception as e:
                print(f"❌ 解析消息失败: {e}")
                
        client.on_connect = on_connect
        client.on_disconnect = on_disconnect
        client.on_message = on_message
        
        return client
        
    def connect_to_broker(self, broker_config: Dict[str, Any]) -> bool:
        """连接到MQTT Broker"""
        client = self.create_mqtt_client(broker_config)
        
        try:
            client.connect(broker_config["host"], broker_config["port"], 60)
            client.loop_start()
            self.clients[broker_config["name"]] = client
            time.sleep(1)  # 等待连接建立
            return True
        except Exception as e:
            print(f"❌ 连接失败 {broker_config['name']}: {e}")
            return False
            
    def disconnect_all(self):
        """断开所有连接"""
        for name, client in self.clients.items():
            try:
                client.loop_stop()
                client.disconnect()
                print(f"🔌 断开连接: {name}")
            except Exception as e:
                print(f"❌ 断开连接失败 {name}: {e}")
                
    def generate_sensor_data(self, device_template: Dict[str, Any]) -> Dict[str, Any]:
        """生成传感器数据"""
        formaldehyde = round(random.uniform(0.01, 0.15), 3)
        temperature = round(random.uniform(18.0, 30.0), 1)
        humidity = round(random.uniform(40.0, 80.0), 1)
        battery = random.randint(20, 100)
        signal_strength = random.randint(-90, -30)
        
        return {
            "device_id": device_template["device_id"],
            "device_type": device_template["device_type"],
            "sensor_id": device_template.get("sensor_id", f"sensor_{device_template['device_id']}_01"),
            "sensor_type": device_template.get("sensor_type", "hcho"),
            "timestamp": int(time.time()),
            "data": {
                "formaldehyde": formaldehyde,
                "temperature": temperature,
                "humidity": humidity,
                "battery": battery
            },
            "location": device_template["location"],
            "quality": {
                "signal_strength": signal_strength,
                "data_quality": "good" if formaldehyde < 0.1 else "poor"
            }
        }
        
    def publish_sensor_data(self, broker_name: str, device_template: Dict[str, Any]):
        """发布传感器数据"""
        if broker_name not in self.clients:
            print(f"❌ Broker {broker_name} 未连接")
            return False
            
        client = self.clients[broker_name]
        sensor_data = self.generate_sensor_data(device_template)
        
        topic = f"air-quality/hcho/{device_template['device_id']}/data"
        message = json.dumps(sensor_data, ensure_ascii=False)
        
        try:
            result = client.publish(topic, message, qos=1)
            if result.rc == mqtt.MQTT_ERR_SUCCESS:
                print(f"📤 发布数据 - Broker: {broker_name}, 设备: {device_template['device_id']}")
                return True
            else:
                print(f"❌ 发布失败 - Broker: {broker_name}, 错误: {result}")
                return False
        except Exception as e:
            print(f"❌ 发布数据时出错 - Broker: {broker_name}: {e}")
            return False
            
    def run_scenario(self, scenario: Dict[str, Any], broker_name: str):
        """运行测试场景"""
        print(f"\n🧪 开始测试场景: {scenario['name']}")
        print(f"   描述: {scenario['description']}")
        print(f"   设备数量: {scenario['device_count']}")
        print(f"   消息间隔: {scenario['message_interval']}秒")
        print(f"   测试时长: {scenario['test_duration']}秒")
        
        if broker_name not in self.clients:
            print(f"❌ Broker {broker_name} 未连接")
            return
            
        # 选择设备模板
        device_templates = self.config["device_templates"][:scenario["device_count"]]
        
        start_time = time.time()
        message_count = 0
        
        try:
            while time.time() - start_time < scenario["test_duration"]:
                for device_template in device_templates:
                    if self.publish_sensor_data(broker_name, device_template):
                        message_count += 1
                        
                if scenario["message_interval"] > 0:
                    time.sleep(scenario["message_interval"])
                else:
                    break  # 一次性测试
                    
        except KeyboardInterrupt:
            print(f"\n🛑 测试被用户中断")
            
        print(f"\n✅ 测试场景完成: {scenario['name']}")
        print(f"   总消息数: {message_count}")
        print(f"   实际时长: {time.time() - start_time:.1f}秒")
        
    def list_brokers(self):
        """列出可用的MQTT Broker"""
        print("\n📋 可用的MQTT Broker:")
        for i, broker in enumerate(self.config["mqtt_brokers"], 1):
            print(f"  {i}. {broker['name']}")
            print(f"     地址: {broker['host']}:{broker['port']}")
            print(f"     描述: {broker['description']}")
            print()
            
    def list_scenarios(self):
        """列出可用的测试场景"""
        print("\n📋 可用的测试场景:")
        for i, scenario in enumerate(self.config["test_scenarios"], 1):
            print(f"  {i}. {scenario['name']}")
            print(f"     描述: {scenario['description']}")
            print(f"     设备数: {scenario['device_count']}, 间隔: {scenario['message_interval']}秒")
            print()
            
    def interactive_mode(self):
        """交互模式"""
        print("🎮 进入交互模式")
        print("输入 'help' 查看可用命令")
        
        while True:
            try:
                command = input("\n> ").strip().lower()
                
                if command == "help":
                    print("\n可用命令:")
                    print("  brokers - 列出MQTT Broker")
                    print("  scenarios - 列出测试场景")
                    print("  connect <broker_index> - 连接到指定Broker")
                    print("  run <scenario_index> - 运行测试场景")
                    print("  disconnect - 断开所有连接")
                    print("  quit - 退出程序")
                    
                elif command == "brokers":
                    self.list_brokers()
                    
                elif command == "scenarios":
                    self.list_scenarios()
                    
                elif command.startswith("connect "):
                    try:
                        broker_index = int(command.split()[1]) - 1
                        if 0 <= broker_index < len(self.config["mqtt_brokers"]):
                            broker = self.config["mqtt_brokers"][broker_index]
                            self.connect_to_broker(broker)
                        else:
                            print("❌ 无效的Broker索引")
                    except (ValueError, IndexError):
                        print("❌ 请提供有效的Broker索引")
                        
                elif command.startswith("run "):
                    try:
                        scenario_index = int(command.split()[1]) - 1
                        if 0 <= scenario_index < len(self.config["test_scenarios"]):
                            scenario = self.config["test_scenarios"][scenario_index]
                            # 使用第一个连接的Broker
                            if self.clients:
                                broker_name = list(self.clients.keys())[0]
                                self.run_scenario(scenario, broker_name)
                            else:
                                print("❌ 请先连接到MQTT Broker")
                        else:
                            print("❌ 无效的场景索引")
                    except (ValueError, IndexError):
                        print("❌ 请提供有效的场景索引")
                        
                elif command == "disconnect":
                    self.disconnect_all()
                    
                elif command == "quit":
                    break
                    
                else:
                    print("❌ 未知命令，输入 'help' 查看可用命令")
                    
            except KeyboardInterrupt:
                print("\n🛑 程序被用户中断")
                break
            except Exception as e:
                print(f"❌ 命令执行出错: {e}")
                
        self.disconnect_all()


def main():
    """主函数"""
    parser = argparse.ArgumentParser(description="配置驱动的MQTT测试程序")
    parser.add_argument("--config", default="test_config.json", help="配置文件路径")
    parser.add_argument("--broker", type=int, help="Broker索引")
    parser.add_argument("--scenario", type=int, help="测试场景索引")
    parser.add_argument("--interactive", action="store_true", help="交互模式")
    
    args = parser.parse_args()
    
    # 创建测试器
    tester = ConfigDrivenTester(args.config)
    
    try:
        if args.interactive:
            # 交互模式
            tester.interactive_mode()
        else:
            # 命令行模式
            if args.broker is None or args.scenario is None:
                print("❌ 请指定 --broker 和 --scenario 参数，或使用 --interactive 进入交互模式")
                return
                
            # 连接到指定的Broker
            if 0 <= args.broker - 1 < len(tester.config["mqtt_brokers"]):
                broker = tester.config["mqtt_brokers"][args.broker - 1]
                if tester.connect_to_broker(broker):
                    # 运行指定的测试场景
                    if 0 <= args.scenario - 1 < len(tester.config["test_scenarios"]):
                        scenario = tester.config["test_scenarios"][args.scenario - 1]
                        tester.run_scenario(scenario, broker["name"])
                    else:
                        print("❌ 无效的测试场景索引")
                else:
                    print("❌ 连接MQTT Broker失败")
            else:
                print("❌ 无效的Broker索引")
                
    except KeyboardInterrupt:
        print("\n🛑 程序被用户中断")
    except Exception as e:
        print(f"\n❌ 程序出错: {e}")
    finally:
        tester.disconnect_all()


if __name__ == "__main__":
    main()
