package config

// Application 应用配置
type Application struct {
	Processes int `yaml:"processes"`
}

// ParseApplication 解析Application配置
func ParseApplication(in []byte) (app *Application, err error) {

	return
}
