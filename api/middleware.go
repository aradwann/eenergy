package api

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/aradwann/eenergy/token"

	"github.com/gin-gonic/gin"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "authorization_payload"
)

// authMiddleware is a higher order function that returns the auth middleware function
func authMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// extract HTTP authorization header
		authorizationHeader := ctx.GetHeader(authorizationHeaderKey)
		// check if the header is provided
		if len(authorizationHeader) == 0 {
			err := errors.New("authorization header is not provided")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		// separate authorization header strings by space
		fields := strings.Fields(authorizationHeader)
		// check if it is two fields
		if len(fields) < 2 {
			err := errors.New("invalid authorization header format")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		// get the authorization header type
		authorizationType := strings.ToLower(fields[0])
		// check if the type is bearer
		if authorizationType != authorizationTypeBearer {
			err := fmt.Errorf("unsupported authorization type %s", authorizationType)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		// get the acccess token
		accessToken := fields[1]
		// verify the access token and extract the payload
		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		// store the token payload in the context
		ctx.Set(authorizationPayloadKey, payload)
		// forward to the next middleware
		ctx.Next()

	}
}
