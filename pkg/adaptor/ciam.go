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
	Register(req RegisterRequest) (res interface{}, e *domain.TechnicalError)
	ConfirmRegister(req ConfirmRegisterRequest) (res interface{}, e *domain.TechnicalError)
	Login(req LoginRequest) (res interface{}, e *domain.TechnicalError)
	Logout(req domain.SessionRequest) (res interface{}, e *domain.TechnicalError)
	JwtInfo(t string) (res map[string]interface{}, e *domain.TechnicalError)
	RefreshToken(req domain.SessionRequest) (res interface{}, e *domain.TechnicalError)
	ChangePassword(req ChangePasswordRequest) (res interface{}, e *domain.TechnicalError)
}
