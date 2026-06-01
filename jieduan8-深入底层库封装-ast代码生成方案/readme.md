# 阶段8 深入底层库封装-ast代码生成方案

## 24周 log日志包设计

### 1章 如何设计日志包

#### 1-1为什么需要自己去设计日志包
- 见`jieduan8-深入底层库封装-ast代码生成方案/log/main/main.go`
- 自己封装的日志包的重要性
#### 1-2 go-zero和kratos中日志的处理
- go-zero 日志：全局内置统一logger，强集成、高性能、开箱即用（基于 zap）
  - 核心组件
  - logx：底层日志库（封装 zap）
  - logc：带 context 的封装（自动注入 traceid）
- kratos 日志：接口化、可插拔、弱绑定（适配器模式）
  - 核心设计（四大角色）
  - Logger（接口）：底层适配（zap/slog/ 自定义）
    - 只定义行为，不绑定实现。依赖注入的logger机制
  - Helper：业务代码调用（log.Info）
  - Filter：脱敏、过滤
  - Valuer：动态注入字段（traceid、时间）

#### 1-3 全局logger和传递参数的logger的用法

- 主要是讲我们在设计logger时，我们可以考虑2种：
  - 全局内置统一logger
    - 如果我们不考虑某个模块拿出来做开源项目，全局大家都依赖同一个logger，项目启动时就初始化好logger，后续所有地方直接引用它即可
  - 依赖注入的logger机制
    - 如果想要开源一部分代码，就把logger这块设计成一个内部通用接口，由外界自己根据要求的接口类型实现并注入logger实例即可

#### 1-4 日志包的基本需求
- 见`jieduan8-深入底层库封装-ast代码生成方案/log/main/main.go`

## 25周 ast代码生成工具开发