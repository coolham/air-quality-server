#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
MQTT认证测试脚本
测试不同的认证场景
"""

import paho.mqtt.client as mqtt
import json
import time
from colorama import init, Fore, Style

# 初始化colorama
init(autoreset=True)

class MQTTAuthTester:
    def __init__(self):
        self.client = None
        self.connected = False
        self.test_results = []

    def on_connect(self, client, userdata, flags, rc):
        """连接回调"""
        if rc == 0:
            self.connected = True
            print(f"{Fore.GREEN}✅ 连接成功 (RC: {rc})")
        else:
            self.connected = False
            print(f"{Fore.RED}❌ 连接失败 (RC: {rc})")
            print(f"   错误原因: {self.get_connect_error(rc)}")

    def on_disconnect(self, client, userdata, rc):
        """断开连接回调"""
        self.connected = False
        if rc == 0:
            print(f"{Fore.YELLOW}🔌 正常断开连接")
        else:
            print(f"{Fore.RED}🔌 异常断开连接 (RC: {rc})")

    def on_publish(self, client, userdata, mid):
        """发布消息回调"""
        print(f"{Fore.BLUE}📤 消息发布成功 (ID: {mid})")

    def get_connect_error(self, rc):
        """获取连接错误描述"""
        errors = {
            1: "协议版本不正确",
            2: "客户端ID无效",
            3: "服务器不可用",
            4: "用户名或密码错误",
            5: "未授权"
        }
        return errors.get(rc, f"未知错误 ({rc})")

    def test_connection(self, test_name, username=None, password=None, client_id=None):
        """测试连接"""
        print(f"\n{Fore.CYAN}🧪 测试: {test_name}")
        print(f"   用户名: {username or '无'}")
        print(f"   密码: {'***' if password else '无'}")
        print(f"   客户端ID: {client_id or '自动生成'}")
        
        # 创建客户端
        self.client = mqtt.Client(client_id=client_id)
        self.client.on_connect = self.on_connect
        self.client.on_disconnect = self.on_disconnect
        self.client.on_publish = self.on_publish
        
        # 设置用户名和密码
        if username or password:
            self.client.username_pw_set(username, password)
        
        try:
            print(f"{Fore.YELLOW}🔗 正在连接到 localhost:1883...")
            self.client.connect("localhost", 1883, 60)
            
            # 等待连接结果
            self.client.loop_start()
            time.sleep(2)
            
            if self.connected:
                # 尝试发布消息
                topic = "air-quality/hcho/test_auth/data"
                message = json.dumps({
                    "device_id": "test_auth",
                    "timestamp": int(time.time()),
                    "test": test_name,
                    "formaldehyde": 0.05
                })
                
                result = self.client.publish(topic, message, qos=1)
                if result.rc == mqtt.MQTT_ERR_SUCCESS:
                    print(f"{Fore.GREEN}✅ 认证测试通过 - 可以发布消息")
                    self.test_results.append((test_name, True, "成功"))
                else:
                    print(f"{Fore.RED}❌ 认证测试失败 - 无法发布消息")
                    self.test_results.append((test_name, False, "发布失败"))
            else:
                print(f"{Fore.RED}❌ 认证测试失败 - 连接被拒绝")
                self.test_results.append((test_name, False, "连接失败"))
                
        except Exception as e:
            print(f"{Fore.RED}❌ 连接异常: {e}")
            self.test_results.append((test_name, False, f"异常: {e}"))
        finally:
            if self.client:
                self.client.loop_stop()
                self.client.disconnect()
                time.sleep(1)

    def run_all_tests(self):
        """运行所有认证测试"""
        print(f"{Fore.MAGENTA}{'='*60}")
        print(f"{Fore.MAGENTA}🔐 MQTT认证机制测试")
        print(f"{Fore.MAGENTA}{'='*60}")
        
        # 测试1: 无认证连接
        self.test_connection("无认证连接")
        
        # 测试2: 空用户名密码
        self.test_connection("空用户名密码", "", "")
        
        # 测试3: 只有用户名
        self.test_connection("只有用户名", "test_user")
        
        # 测试4: 用户名和密码
        self.test_connection("用户名和密码", "test_user", "test_pass")
        
        # 测试5: 配置中的用户名密码
        self.test_connection("配置用户名密码", "admin", "password")
        
        # 测试6: 错误的用户名密码
        self.test_connection("错误用户名密码", "wrong_user", "wrong_pass")
        
        # 测试7: 特殊字符用户名密码
        self.test_connection("特殊字符认证", "user@domain", "pass#123")
        
        # 测试8: 长用户名密码
        long_username = "a" * 100
        long_password = "b" * 100
        self.test_connection("长用户名密码", long_username, long_password)
        
        # 测试9: 中文用户名密码
        self.test_connection("中文认证", "测试用户", "测试密码")
        
        # 测试10: 自定义客户端ID
        self.test_connection("自定义客户端ID", client_id="custom_client_123")

    def print_summary(self):
        """打印测试总结"""
        print(f"\n{Fore.MAGENTA}{'='*60}")
        print(f"{Fore.MAGENTA}📊 测试结果总结")
        print(f"{Fore.MAGENTA}{'='*60}")
        
        success_count = 0
        total_count = len(self.test_results)
        
        for test_name, success, reason in self.test_results:
            status = f"{Fore.GREEN}✅ 通过" if success else f"{Fore.RED}❌ 失败"
            print(f"{status} {test_name}: {reason}")
            if success:
                success_count += 1
        
        print(f"\n{Fore.CYAN}📈 统计信息:")
        print(f"   总测试数: {total_count}")
        print(f"   成功数: {success_count}")
        print(f"   失败数: {total_count - success_count}")
        print(f"   成功率: {success_count/total_count*100:.1f}%")
        
        if success_count == total_count:
            print(f"\n{Fore.GREEN}🎉 所有测试都通过了！")
            print(f"{Fore.YELLOW}💡 结论: MQTT服务器允许所有连接，不需要认证")
        else:
            print(f"\n{Fore.YELLOW}⚠️  部分测试失败")
            print(f"{Fore.YELLOW}💡 结论: MQTT服务器可能有认证限制")

def main():
    """主函数"""
    tester = MQTTAuthTester()
    tester.run_all_tests()
    tester.print_summary()

if __name__ == "__main__":
    main()
