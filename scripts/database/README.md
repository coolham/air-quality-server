# 数据库管理脚本

本目录包含数据库初始化、检查和管理的脚本。

## 脚本说明

- **init-database-sql.bat** - 使用SQL脚本初始化数据库（推荐）
- **init-database.bat** - 使用Go迁移工具初始化数据库
- **check-database.bat** - 检查数据库状态和表结构
- **init.sql** - 完整的数据库初始化SQL脚本

## 使用流程

### 初始化数据库
```cmd
# 方法1：SQL脚本（推荐）
init-database-sql.bat

# 方法2：Go迁移工具
init-database.bat
```

### 检查数据库
```cmd
check-database.bat
```
