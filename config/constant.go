package config

// App const
const (
	// AppName 应用名称
	AppName = "webrouting"
	// HTTPStatusBadGateway HTTP状态码 Bad Gateway
	HTTPStatusBadGateway = 502
)

// Default const
const (
	// DefaultTCPTimeout 后端服务器连接超时时间
	DefaultTCPTimeout int64 = 3000
	// DefaultTCPRetries 后端服务器失败重试次数
	DefaultTCPRetries int = 0
	// DefaultCheckInterval tcp恢复检查间隔
	DefaultCheckInterval int64 = 3000
)
