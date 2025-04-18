package service

type HelloService struct{}

func (h HelloService) Hello() string {
	return "Hello Go"
}
