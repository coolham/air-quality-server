#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
MQTT数据存储测试脚本
测试MQTT服务器是否能正确接收和存储传感器数据
"""

import json
import time
import paho.mqtt.client as mqtt
from datetime import datetime
import random

class MQTTDataStorageTester:
    def __init__(self, broker_host="localhost", broker_port=1883):
        self.broker_host = broker_host
        self.broker_port = broker_port
        self.client = None
        self.connected = False
        
    def on_connect(self, client, userdata, flags, rc):
        """连接回调"""
        if rc == 0:
            print(f"✅ 成功连接到MQTT服务器 {self.broker_host}:{self.broker_port}")
            self.connected = True
        else:
            print(f"❌ 连接失败，错误代码: {rc}")
            self.connected = False
    
    def on_disconnect(self, client, userdata, rc):
        """断开连接回调"""
        print(f"🔌 与MQTT服务器断开连接，错误代码: {rc}")
        self.connected = False
    
    def on_publish(self, client, userdata, mid):
        """发布消息回调"""
        print(f"📤 消息发布成功，消息ID: {mid}")
    
    def connect(self):
        """连接到MQTT服务器"""
        try:
            self.client = mqtt.Client()
            self.client.on_connect = self.on_connect
            self.client.on_disconnect = self.on_disconnect
            self.client.on_publish = self.on_publish
            
            print(f"🔗 正在连接到MQTT服务器 {self.broker_host}:{self.broker_port}...")
            self.client.connect(self.broker_host, self.broker_port, 60)
            self.client.loop_start()
            
            # 等待连接建立
            timeout = 10
            while not self.connected and timeout > 0:
                time.sleep(0.1)
                timeout -= 0.1
            
            if not self.connected:
                print("❌ 连接超时")
                return False
            
            return True
            
        except Exception as e:
            print(f"❌ 连接异常: {e}")
            return False
    
    def disconnect(self):
        """断开连接"""
        if self.client:
            self.client.loop_stop()
            self.client.disconnect()
            print("🔌 已断开MQTT连接")
    
    def generate_sensor_data(self, device_id="hcho_001"):
        """生成传感器数据"""
        # 生成模拟的传感器数据
        formaldehyde = round(random.uniform(0.01, 0.15), 3)  # 甲醛浓度 0.01-0.15 mg/m³
        temperature = round(random.uniform(18.0, 28.0), 1)   # 温度 18-28°C
        humidity = round(random.uniform(30.0, 70.0), 1)     # 湿度 30-70%
        battery = random.randint(20, 100)                   # 电池电量 20-100%
        
        # 生成位置信息
        latitude = round(random.uniform(39.9, 40.0), 6)     # 北京纬度范围
        longitude = round(random.uniform(116.3, 116.5), 6)  # 北京经度范围
        
        data = {
            "device_id": device_id,
            "device_type": "hcho",
            "sensor_id": f"sensor_{device_id}_01",
            "sensor_type": "hcho",
            "timestamp": int(time.time()),
            "data": {
                "formaldehyde": formaldehyde,
                "temperature": temperature,
                "humidity": humidity,
                "battery": battery
            },
            "location": {
                "latitude": latitude,
                "longitude": longitude,
                "address": "北京市朝阳区测试位置"
            },
            "quality": {
                "signal_strength": random.randint(-80, -40),  # 信号强度 -80到-40 dBm
                "data_quality": "good"
            }
        }
        
        return data
    
    def publish_sensor_data(self, device_id="hcho_001", count=5):
        """发布传感器数据"""
        if not self.connected:
            print("❌ 未连接到MQTT服务器")
            return False
        
        topic = f"air-quality/hcho/{device_id}/data"
        print(f"📡 开始发布传感器数据到主题: {topic}")
        
        for i in range(count):
            # 生成数据
            data = self.generate_sensor_data(device_id)
            
            # 转换为JSON
            payload = json.dumps(data, ensure_ascii=False, indent=2)
            
            print(f"\n📊 第 {i+1} 条数据:")
            print(f"   设备ID: {data['device_id']}")
            print(f"   传感器ID: {data['sensor_id']}")
            print(f"   甲醛浓度: {data['data']['formaldehyde']} mg/m³")
            print(f"   温度: {data['data']['temperature']}°C")
            print(f"   湿度: {data['data']['humidity']}%")
            print(f"   电池电量: {data['data']['battery']}%")
            print(f"   信号强度: {data['quality']['signal_strength']} dBm")
            
            # 发布消息
            try:
                result = self.client.publish(topic, payload, qos=1)
                if result.rc == mqtt.MQTT_ERR_SUCCESS:
                    print(f"   ✅ 发布成功")
                else:
                    print(f"   ❌ 发布失败，错误代码: {result.rc}")
            except Exception as e:
                print(f"   ❌ 发布异常: {e}")
            
            # 等待一段时间再发送下一条
            if i < count - 1:
                time.sleep(2)
        
        print(f"\n🎯 已完成 {count} 条数据的发布")
        return True
    
    def test_multiple_devices(self):
        """测试多个设备的数据发布"""
        devices = ["hcho_001", "hcho_002", "hcho_003"]
        
        print(f"🔄 开始测试多个设备的数据发布...")
        
        for device_id in devices:
            print(f"\n📱 测试设备: {device_id}")
            self.publish_sensor_data(device_id, 3)
            time.sleep(1)
        
        print(f"\n✅ 多设备测试完成")
    
    def test_alert_conditions(self):
        """测试告警条件"""
        print(f"🚨 开始测试告警条件...")
        
        # 生成高甲醛浓度数据
        high_formaldehyde_data = {
            "device_id": "hcho_001",
            "device_type": "hcho",
            "sensor_id": "sensor_hcho_001_01",
            "sensor_type": "hcho",
            "timestamp": int(time.time()),
            "data": {
                "formaldehyde": 0.12,  # 超过0.08的警告阈值
                "temperature": 25.0,
                "humidity": 50.0,
                "battery": 85
            },
            "location": {
                "latitude": 39.9042,
                "longitude": 116.4074,
                "address": "北京市朝阳区"
            },
            "quality": {
                "signal_strength": -65,
                "data_quality": "good"
            }
        }
        
        topic = "air-quality/hcho/hcho_001/data"
        payload = json.dumps(high_formaldehyde_data, ensure_ascii=False, indent=2)
        
        print(f"📊 发布高甲醛浓度数据:")
        print(f"   甲醛浓度: {high_formaldehyde_data['data']['formaldehyde']} mg/m³ (超过警告阈值)")
        
        try:
            result = self.client.publish(topic, payload, qos=1)
            if result.rc == mqtt.MQTT_ERR_SUCCESS:
                print(f"   ✅ 高浓度数据发布成功")
            else:
                print(f"   ❌ 高浓度数据发布失败，错误代码: {result.rc}")
        except Exception as e:
            print(f"   ❌ 高浓度数据发布异常: {e}")

def main():
    print("🧪 MQTT数据存储测试工具")
    print("=" * 50)
    
    # 创建测试器
    tester = MQTTDataStorageTester()
    
    try:
        # 连接到MQTT服务器
        if not tester.connect():
            print("❌ 无法连接到MQTT服务器，请检查服务器是否运行")
            return
        
        # 等待连接稳定
        time.sleep(1)
        
        # 测试基本数据发布
        print("\n📊 测试1: 基本传感器数据发布")
        tester.publish_sensor_data("hcho_001", 5)
        
        # 等待一段时间
        time.sleep(3)
        
        # 测试多设备
        print("\n📱 测试2: 多设备数据发布")
        tester.test_multiple_devices()
        
        # 等待一段时间
        time.sleep(3)
        
        # 测试告警条件
        print("\n🚨 测试3: 告警条件测试")
        tester.test_alert_conditions()
        
        print("\n🎉 所有测试完成！")
        print("\n📋 请检查以下内容:")
        print("   1. MQTT服务器日志中是否显示消息接收")
        print("   2. 数据库中是否存储了传感器数据")
        print("   3. 是否生成了告警记录")
        print("   4. 设备状态是否更新")
        
    except KeyboardInterrupt:
        print("\n⏹️ 测试被用户中断")
    except Exception as e:
        print(f"\n❌ 测试异常: {e}")
    finally:
        # 断开连接
        tester.disconnect()

if __name__ == "__main__":
    main()
