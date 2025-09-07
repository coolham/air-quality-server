#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
MQTT测试工具演示脚本
展示如何使用各种测试工具
"""

import time
import subprocess
import sys
import os

def print_banner():
    """打印横幅"""
    print("=" * 60)
    print("    MQTT测试工具演示")
    print("=" * 60)
    print()

def print_section(title):
    """打印章节标题"""
    print(f"\n📋 {title}")
    print("-" * 40)

def run_command(command, description):
    """运行命令并显示结果"""
    print(f"\n🚀 {description}")
    print(f"命令: {command}")
    print()
    
    try:
        result = subprocess.run(command, shell=True, capture_output=True, text=True)
        if result.returncode == 0:
            print("✅ 命令执行成功")
            if result.stdout:
                print("输出:")
                print(result.stdout)
        else:
            print("❌ 命令执行失败")
            if result.stderr:
                print("错误:")
                print(result.stderr)
    except Exception as e:
        print(f"❌ 执行命令时出错: {e}")

def check_dependencies():
    """检查依赖"""
    print_section("检查依赖")
    
    # 检查Python
    try:
        result = subprocess.run([sys.executable, "--version"], capture_output=True, text=True)
        if result.returncode == 0:
            print(f"✅ Python版本: {result.stdout.strip()}")
        else:
            print("❌ Python未正确安装")
            return False
    except Exception as e:
        print(f"❌ 检查Python时出错: {e}")
        return False
    
    # 检查paho-mqtt
    try:
        import paho.mqtt.client
        print("✅ paho-mqtt已安装")
    except ImportError:
        print("❌ paho-mqtt未安装，正在安装...")
        run_command("pip install paho-mqtt", "安装paho-mqtt")
    
    return True

def demo_basic_test():
    """演示基础测试"""
    print_section("基础测试演示")
    
    print("这个演示将发送3条测试消息到MQTT Broker")
    print("请确保MQTT Broker正在运行 (localhost:1883)")
    
    input("\n按回车键开始测试...")
    
    # 运行基础测试
    run_command(
        f"{sys.executable} basic_test.py --count 3 --device-id demo_device_001",
        "发送3条测试消息"
    )

def demo_advanced_test():
    """演示高级测试"""
    print_section("高级测试演示")
    
    print("这个演示将测试配置下发和命令控制功能")
    print("请确保MQTT Broker正在运行")
    
    input("\n按回车键开始测试...")
    
    # 运行高级测试
    run_command(
        f"{sys.executable} advanced_test.py --test-type all --device-id demo_device_002",
        "运行完整的高级测试"
    )

def demo_config_driven_test():
    """演示配置驱动测试"""
    print_section("配置驱动测试演示")
    
    print("这个演示将展示配置驱动的测试功能")
    print("配置文件: test_config.json")
    
    # 显示配置文件内容
    if os.path.exists("test_config.json"):
        print("\n📄 配置文件内容预览:")
        with open("test_config.json", 'r', encoding='utf-8') as f:
            import json
            config = json.load(f)
            print(f"  - MQTT Broker数量: {len(config['mqtt_brokers'])}")
            print(f"  - 测试场景数量: {len(config['test_scenarios'])}")
            print(f"  - 设备模板数量: {len(config['device_templates'])}")
    
    input("\n按回车键开始测试...")
    
    # 运行配置驱动测试
    run_command(
        f"{sys.executable} config_driven_test.py --broker 1 --scenario 1",
        "运行配置驱动的测试"
    )

def demo_simulation():
    """演示持续模拟"""
    print_section("持续模拟演示")
    
    print("这个演示将启动持续模拟模式")
    print("设备将每10秒发送一次数据，持续30秒")
    print("按Ctrl+C可以提前停止")
    
    input("\n按回车键开始模拟...")
    
    # 运行持续模拟
    run_command(
        f"{sys.executable} basic_test.py --simulate --interval 10 --device-id demo_device_003",
        "启动持续模拟"
    )

def show_help():
    """显示帮助信息"""
    print_section("帮助信息")
    
    print("可用的测试工具:")
    print("1. basic_test.py - 基础MQTT测试")
    print("2. advanced_test.py - 高级MQTT测试")
    print("3. config_driven_test.py - 配置驱动测试")
    print("4. run_test.bat/run_test.sh - 便捷脚本")
    
    print("\n常用命令:")
    print("• 快速测试: python basic_test.py --count 10")
    print("• 持续模拟: python basic_test.py --simulate")
    print("• 高级测试: python advanced_test.py --test-type all")
    print("• 交互模式: python config_driven_test.py --interactive")
    
    print("\nMQTT主题格式:")
    print("• 数据主题: air-quality/hcho/{device_id}/data")
    print("• 状态主题: air-quality/hcho/{device_id}/status")
    print("• 配置主题: air-quality/hcho/{device_id}/config")
    print("• 命令主题: air-quality/hcho/{device_id}/command")
    print("• 响应主题: air-quality/hcho/{device_id}/response")

def main():
    """主函数"""
    print_banner()
    
    # 检查依赖
    if not check_dependencies():
        print("❌ 依赖检查失败，请先解决依赖问题")
        return
    
    while True:
        print("\n" + "=" * 60)
        print("请选择演示内容:")
        print("1. 基础测试演示")
        print("2. 高级测试演示")
        print("3. 配置驱动测试演示")
        print("4. 持续模拟演示")
        print("5. 显示帮助信息")
        print("6. 退出")
        print("=" * 60)
        
        choice = input("\n请输入选择 (1-6): ").strip()
        
        if choice == "1":
            demo_basic_test()
        elif choice == "2":
            demo_advanced_test()
        elif choice == "3":
            demo_config_driven_test()
        elif choice == "4":
            demo_simulation()
        elif choice == "5":
            show_help()
        elif choice == "6":
            print("\n👋 感谢使用MQTT测试工具演示!")
            break
        else:
            print("❌ 无效选择，请输入1-6之间的数字")
        
        input("\n按回车键继续...")

if __name__ == "__main__":
    try:
        main()
    except KeyboardInterrupt:
        print("\n\n🛑 演示被用户中断")
    except Exception as e:
        print(f"\n❌ 演示过程中出错: {e}")
