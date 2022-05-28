package api

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/Franklynoble/bankapp/db/token"
	"github.com/gin-gonic/gin"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadkey = "authorization_payload"
)

func authMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authauthorizationHeader := ctx.GetHeader(authorizationHeaderKey)

		if len(authauthorizationHeader) == 0 {
			err := errors.New("authorization header is not provided")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errResponse(err))

			return
		}

		fields := strings.Fields(authauthorizationHeader) // split  the authorization by space

		if len(fields) < 2 { //check the len
			err := errors.New("invalid authorization header format")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errResponse(err))
			return
		}
		authorizationType := strings.ToLower(fields[0]) // convert to lower for easey comparism
		if authorizationType != authorizationTypeBearer {
			err := fmt.Errorf("unsopported authorization  typ %s", authorizationType)

			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errResponse(err))
			return
		}
		accessToken := fields[1]
		payload, err := tokenMaker.VerifyToken(accessToken)

		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errResponse(err))
			return
		}

		//token is valid store the payload in the context and  pass in the  value
		ctx.Set(authorizationPayloadkey, payload)
	}
}
