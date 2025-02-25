package shell

import "github.com/Ali-Farhadnia/goshell/internal/service/user"

type Session struct {
	User       *user.User
	WorkingDir string
}
