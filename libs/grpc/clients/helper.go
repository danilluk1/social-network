package clients

import "fmt"

func createClientAddress(env, service string, port int) string {
	ip := fmt.Sprintf("dns:///%s", service)
	if env != "production" {
		ip = "127.0.0.1"
	}

	return fmt.Sprintf("%s:%v", ip, port)
}
