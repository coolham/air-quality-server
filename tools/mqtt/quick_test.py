#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
MQTT数据存储快速测试脚本
"""

import json
import time
import paho.mqtt.client as mqtt

def test_mqtt_data_storage():
    """测试MQTT数据存储"""
    
    # 创建MQTT客户端
    client = mqtt.Client()
    
    def on_connect(client, userdata, flags, rc):
        if rc == 0:
            print("✅ 连接到MQTT服务器成功")
        else:
            print(f"❌ 连接失败，错误代码: {rc}")
    
    def on_publish(client, userdata, mid):
        print(f"📤 消息发布成功，消息ID: {mid}")
    
    # 设置回调
    client.on_connect = on_connect
    client.on_publish = on_publish
    
    try:
        # 连接到MQTT服务器
        print("🔗 正在连接到MQTT服务器...")
        client.connect("localhost", 1883, 60)
        client.loop_start()
        
        # 等待连接建立
        time.sleep(1)
        
        # 准备测试数据
        test_data = {
            "device_id": "hcho_001",
            "device_type": "hcho",
            "sensor_id": "sensor_hcho_001_01",
            "sensor_type": "hcho",
            "timestamp": int(time.time()),
            "data": {
                "formaldehyde": 0.05,
                "temperature": 22.5,
                "humidity": 45.0,
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
        
        # 发布消息
        topic = "air-quality/hcho/hcho_001/data"
        payload = json.dumps(test_data, ensure_ascii=False, indent=2)
        
        print(f"📊 发布测试数据到主题: {topic}")
        print(f"📋 数据内容:")
        print(f"   设备ID: {test_data['device_id']}")
        print(f"   传感器ID: {test_data['sensor_id']}")
        print(f"   甲醛浓度: {test_data['data']['formaldehyde']} mg/m³")
        print(f"   温度: {test_data['data']['temperature']}°C")
        print(f"   湿度: {test_data['data']['humidity']}%")
        print(f"   电池电量: {test_data['data']['battery']}%")
        
        # 发布消息
        result = client.publish(topic, payload, qos=1)
        
        if result.rc == mqtt.MQTT_ERR_SUCCESS:
            print("✅ 消息发布成功")
        else:
            print(f"❌ 消息发布失败，错误代码: {result.rc}")
        
        # 等待消息处理
        time.sleep(2)
        
        print("\n📋 请检查以下内容:")
        print("   1. MQTT服务器日志中是否显示消息接收")
        print("   2. 数据库中是否存储了传感器数据")
        print("   3. 设备状态是否更新")
        
    except Exception as e:
        print(f"❌ 测试异常: {e}")
    finally:
        # 断开连接
        client.loop_stop()
        client.disconnect()
        print("🔌 已断开MQTT连接")

if __name__ == "__main__":
    print("🧪 MQTT数据存储快速测试")
    print("=" * 40)
    test_mqtt_data_storage()
