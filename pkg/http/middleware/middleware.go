package middleware

var whitelists = []string{"10.10.10.1"}

type Authorization struct {
	Username   string
	Role       string
	Permission string
}

type Client struct {
	IPAddress string
	Authorization
}
