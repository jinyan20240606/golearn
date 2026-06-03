# errors
基于 github.com/pkg/errors 包，增加对 error code 的支持，完全兼容 github.com/pkg/errors。

性能跟 github.com/pkg/errors 基本持平。

## CODE设计规范
### 错误描述规范
错误描述包括：对外的错误描述和对内的错误描述两部分。

对内使用code 对外使用error msg。


