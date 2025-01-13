package auth

import (
	"fmt"

	"github.com/casbin/casbin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Authorizer struct {
    enforcer *casbin.Enforcer
}

func NewAuthorizer(model string, policy string) *Authorizer {
    enforcer := casbin.NewEnforcer(model, policy)
    return &Authorizer{
        enforcer: enforcer,
    }
}

func (a *Authorizer) Authorize(subject string, object string, action string) error {
    if !a.enforcer.Enforce(subject, object, action) {
        msg := fmt.Sprintf(
            "%s not permitted for %s on %s",
            subject,
            action,
            object,
        )
        st := status.New(codes.PermissionDenied, msg)
        return st.Err()
    }
    return nil
}

