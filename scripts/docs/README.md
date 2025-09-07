# Scripts 目录说明

本目录包含用于Windows系统的各种管理脚本，包括防火墙管理、数据库初始化和应用程序管理。

## 脚本分类

### SQL脚本

#### init.sql
完整的数据库初始化脚本，包含所有表结构和初始数据。

**使用方法：**
```cmd
# 使用SQL脚本创建数据库
mysql -u root -p < scripts/init.sql
```

**功能：**
- 创建数据库和所有表结构
- 插入默认用户、角色、设备、配置等初始数据
- 包含甲醛监测相关的表和配置
- 创建视图和索引
- 适用于新系统部署

### 数据库管理脚本

#### 1. init-database-sql.bat
使用SQL脚本初始化数据库，创建表结构和插入初始数据。

**使用方法：**
```cmd
init-database-sql.bat
```

**功能：**
- 使用SQL脚本创建数据库和所有表结构
- 插入默认用户、角色、配置等初始数据
- 包含甲醛监测相关的表和配置
- 交互式输入MySQL root密码
- 提供详细的执行反馈

#### 2. init-database.bat
使用Go迁移工具初始化数据库，创建表结构和插入初始数据。

**使用方法：**
```cmd
init-database.bat
```

**功能：**
- 使用Go迁移工具创建所有数据表
- 插入默认用户、角色、配置等初始数据
- 提供详细的执行反馈

#### 3. check-database.bat
检查数据库状态，显示表结构和数据统计。

**使用方法：**
```cmd
check-database.bat
```

**功能：**
- 显示所有表的存在状态
- 统计各表的数据量
- 验证数据库初始化结果

### 应用程序管理脚本

#### 3. build-app.bat
构建应用程序，生成可执行文件。

**使用方法：**
```cmd
build-app.bat
```

**功能：**
- 编译Go代码生成可执行文件
- 自动创建bin目录
- 提供构建状态反馈

#### 4. start-app.bat
以开发模式启动应用程序（直接运行Go代码）。

**使用方法：**
```cmd
start-app.bat
```

**功能：**
- 直接运行Go代码，无需构建
- 适合开发调试
- 显示访问地址信息

#### 5. run-app.bat
运行已构建的应用程序。

**使用方法：**
```cmd
run-app.bat
```

**功能：**
- 运行bin目录下的可执行文件
- 适合生产环境
- 检查可执行文件是否存在

#### 6. setup.bat
一键安装脚本，完成环境检查和初始化。

**使用方法：**
```cmd
setup.bat
```

**功能：**
- 检查Go环境
- 下载项目依赖
- 初始化数据库
- 构建应用程序

### 防火墙管理脚本

#### 7. configure_firewall.bat
配置Windows防火墙规则，允许应用程序和MQTT端口通过防火墙。

**使用方法：**
```cmd
configure_firewall.bat
```

**功能：**
- 为 `air-quality-server.exe` 添加入站和出站规则
- 为端口 1883 (MQTT) 添加入站和出站规则
- 自动检测可执行文件路径

#### 8. remove_firewall_rules.bat
移除之前添加的防火墙规则。

**使用方法：**
```cmd
remove_firewall_rules.bat
```

**功能：**
- 移除 `air-quality-server` 相关的防火墙规则
- 移除端口 1883 相关的防火墙规则

#### 9. check_firewall_rules.ps1
检查现有的防火墙规则。

**使用方法：**
```powershell
powershell -ExecutionPolicy Bypass -File check_firewall_rules.ps1
```

**功能：**
- 列出所有与应用程序相关的防火墙规则
- 显示规则的详细信息

#### 10. start_with_firewall.bat
一键配置防火墙并启动应用程序。

**使用方法：**
```cmd
start_with_firewall.bat
```

**功能：**
- 自动配置防火墙规则
- 启动应用程序
- 提供清理选项

### 测试脚本

#### 11. test-mqtt-storage.bat
测试MQTT服务器的数据存储功能。

**使用方法：**
```cmd
test-mqtt-storage.bat
```

**功能：**
- 连接到MQTT服务器
- 发送测试传感器数据
- 验证数据存储功能
- 检查告警功能

## 使用流程

### 首次部署

#### 方法1：使用SQL脚本（推荐）

1. **环境检查**
   ```cmd
   setup.bat
   ```

2. **数据库初始化（SQL脚本）**
   ```cmd
   init-database-sql.bat
   ```

3. **验证数据库**
   ```cmd
   check-database.bat
   ```

4. **配置防火墙**
   ```cmd
   configure_firewall.bat
   ```

5. **启动应用程序**
   ```cmd
   start-app.bat
   ```

#### 方法2：使用Go迁移工具

1. **环境检查**
   ```cmd
   setup.bat
   ```

2. **数据库初始化（Go工具）**
   ```cmd
   init-database.bat
   ```

3. **验证数据库**
   ```cmd
   check-database.bat
   ```

4. **配置防火墙**
   ```cmd
   configure_firewall.bat
   ```

5. **启动应用程序**
   ```cmd
   start-app.bat
   ```

### 日常开发

1. **启动应用程序**
   ```cmd
   start-app.bat
   ```

2. **或使用防火墙版本**
   ```cmd
   start_with_firewall.bat
   ```

### 生产部署

1. **构建应用程序**
   ```cmd
   build-app.bat
   ```

2. **运行应用程序**
   ```cmd
   run-app.bat
   ```

## 注意事项

1. **管理员权限**：防火墙相关脚本需要管理员权限
2. **数据库配置**：确保 `config/config.yaml` 中的数据库配置正确
3. **Go环境**：确保Go 1.21+已正确安装
4. **MySQL服务**：确保MySQL服务正在运行

## 故障排除

### 数据库连接失败
- 检查MySQL服务是否运行
- 验证数据库配置
- 确认用户权限

### 防火墙问题
- 以管理员身份运行脚本
- 检查防火墙是否启用
- 确认规则是否正确创建

### 应用程序启动失败
- 检查端口是否被占用
- 验证配置文件
- 查看错误日志

## 技术说明

### 脚本特性
- 支持中文显示（UTF-8编码）
- 提供详细的执行反馈
- 自动错误检测和处理
- 用户友好的交互界面

### 路径处理
- 自动检测项目根目录
- 支持相对路径和绝对路径
- 自动创建必要的目录

### 错误处理
- 检查命令执行结果
- 提供详细的错误信息
- 给出解决建议

---

**文档版本**: v1.0  
**创建日期**: 2024-09-07  
**最后更新**: 2024-09-07  
**作者**: 空气质量监测系统开发团队