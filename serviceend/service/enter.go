package service

// 所以得service 都要在这里注册
type ServiceGroup struct {
	HelloService HelloService
}

var ServiceGroupApp = new(ServiceGroup)
