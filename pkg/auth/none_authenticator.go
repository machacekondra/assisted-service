package auth

import (
	"net/http"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/security"
	"github.com/openshift/assisted-service/pkg/ocm"
	"github.com/sirupsen/logrus"
)

type NoneAuthenticator struct {
	log logrus.FieldLogger
}

func NewNoneAuthenticator(log logrus.FieldLogger) *NoneAuthenticator {
	return &NoneAuthenticator{log: log}
}

func (a *NoneAuthenticator) AuthType() AuthType {
	return TypeNone
}

func (a *NoneAuthenticator) AuthAgentAuth(_ string) (interface{}, error) {
	a.log.Debug("Agent Authentication Disabled")
	return ocm.AdminPayload(), nil
}

func (a *NoneAuthenticator) AuthUserAuth(_ string) (interface{}, error) {
	a.log.Debug("User Authentication Disabled")
	return ocm.AdminPayload(), nil
}

func (a *NoneAuthenticator) CreateAuthenticator() func(_, _ string, authenticate security.TokenAuthentication) runtime.Authenticator {
	return func(_ string, _ string, authenticate security.TokenAuthentication) runtime.Authenticator {
		return security.HttpAuthenticator(func(_ *http.Request) (bool, interface{}, error) {
			p, _ := authenticate("")
			return true, p, nil
		})
	}
}
