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
	// DefaultTCPTimeout 后端服务器连接超时时间/ms
	DefaultRequestTimeout int64 = 10000

	// DefaultClientConnCount 代理客户端最大连接数
	DefaultClientMaxConnCount int = 1024
)
