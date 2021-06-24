package cmd

type Auth struct {
	Host     string
	Port     int
	UserName string
	Password string
}
type TcpRequest struct {
	RequestType string `json:"request_type"`           //请求类型 status start stop restart reload
	ProcessName string `json:"process_name,omitempty"` // 进程名称
	Auth        Auth   `json:"auth"`
}

func NewAuth(host string, port int, username, password string) *Auth {
	return &Auth{
		Host:     host,
		Port:     port,
		UserName: username,
		Password: password,
	}
}
func DefaultAuth() Auth {
	return Auth{
		Host:     "",
		Port:     0,
		UserName: "",
		Password: "",
	}
}
func NewTcpRequest(requestType, processName string) *TcpRequest {
	return &TcpRequest{
		RequestType: requestType,
		ProcessName: processName,
	}
}
