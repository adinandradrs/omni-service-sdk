package adaptor

import "github.com/adinandradrs/omni-service-sdk/pkg/domain"

type (
	LoginRequest struct {
		Identifier string
		Password   string
	}

	RegisterRequest struct {
		Identifier       string
		Password         string
		Email            string
		PhoneNo          string
		Fullname         string
		AdditionalFields []string
		AdditionalValues []string
	}

	ChangePasswordRequest struct {
		Identifier  string
		NewPassword string
		OldPassword string
		domain.SessionRequest
	}

	ConfirmRegisterRequest struct {
		UserId           string
		ConfirmationCode string
		PhoneNo          string
	}
)

type CiamWatcher interface {
	GetSecret(u string) string
	Register(req RegisterRequest) (result interface{}, err error)
	ConfirmRegister(req ConfirmRegisterRequest) (result interface{}, err error)
	Login(req LoginRequest) (result interface{}, err error)
	Logout(req domain.SessionRequest) (result interface{}, err error)
	JwtInfo(t string) (jwtdata map[string]interface{}, err error)
	RefreshToken(req domain.SessionRequest) (result interface{}, err error)
	ChangePassword(req ChangePasswordRequest) (out interface{}, err error)
}
