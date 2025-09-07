# 空气质量监测系统文档索引

## 文档概览

本文档目录包含了空气质量监测系统的完整技术文档，按功能模块分类整理。

## 核心文档

### 1. 系统架构设计
- **[系统架构设计指南](system_architecture_guide.md)** - 完整的系统架构设计，包含微服务模块、接口定义、部署架构等

### 2. 数据库设计与管理
- **[数据库设计与管理指南](database_guide.md)** - 数据库设计、设置、数据模型和管理的完整指南

### 3. MQTT功能实现
- **[MQTT功能设计与实现指南](mqtt_guide.md)** - MQTT协议设计、数据存储、消息处理和系统实现

### 4. Web功能开发
- **[Web功能开发指南](web_features_guide.md)** - Web数据查看、传感器字段管理、界面开发和功能实现

### 5. Golang Web开发
- **[Golang Web开发指南](golang_web_development_guide.md)** - 基于实际项目经验的Gin框架Web开发最佳实践

## 部署文档

### 6. 部署指南
- **[部署指南](deployment_guide.md)** - 系统部署架构、配置和运维指南

### 7. Windows设置
- **[Windows系统设置指南](windows_setup.md)** - Windows系统下的数据库初始化和环境配置

## 文档使用建议

### 新用户入门
1. 首先阅读 **[系统架构设计指南](system_architecture_guide.md)** 了解整体架构
2. 然后阅读 **[数据库设计与管理指南](database_guide.md)** 进行数据库初始化
3. 参考 **[Windows系统设置指南](windows_setup.md)** 完成环境配置
4. 最后阅读 **[部署指南](deployment_guide.md)** 进行系统部署

### 开发者参考
1. **[Golang Web开发指南](golang_web_development_guide.md)** - Web开发最佳实践
2. **[MQTT功能设计与实现指南](mqtt_guide.md)** - MQTT功能开发
3. **[Web功能开发指南](web_features_guide.md)** - Web界面开发

### 运维人员参考
1. **[部署指南](deployment_guide.md)** - 生产环境部署
2. **[数据库设计与管理指南](database_guide.md)** - 数据库运维
3. **[系统架构设计指南](system_architecture_guide.md)** - 系统监控

## 文档更新记录

- **2024-09-07**: 完成文档整理，将重复文档合并为单一指南
- **2024-09-07**: 新增Golang Web开发指南，总结实际开发经验
- **2024-09-07**: 优化文档结构，按功能模块分类

## 贡献指南

如需更新文档，请遵循以下原则：
1. 保持文档的完整性和准确性
2. 使用清晰的中文表达
3. 提供完整的代码示例
4. 包含必要的配置说明
5. 更新文档版本和日期

---

**文档版本**: v1.0  
**创建日期**: 2024-09-07  
**最后更新**: 2024-09-07  
**维护团队**: 空气质量监测系统开发团队
