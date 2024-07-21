package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
)

type VerifySession struct {
	options *sessmodels.VerifySessionOptions
}

func NewVerifySession(opt *sessmodels.VerifySessionOptions) *VerifySession {
	return &VerifySession{opt}
}

func (v *VerifySession) Handle(c *gin.Context) {
	session.VerifySession(v.options, func(rw http.ResponseWriter, r *http.Request) {
		c.Request = c.Request.WithContext(r.Context())
		c.Next()
	})(c.Writer, c.Request)
	c.Abort()
}
