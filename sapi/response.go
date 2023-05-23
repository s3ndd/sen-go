package sapi

import (
	"github.com/s3ndd/sen-go/log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// RespondWithData responds with the given data
func RespondWithData(ctx *gin.Context, statusCode int, data interface{}) {
	ctx.JSON(statusCode, data)
}

// RespondWithStatusCode responds with the given status code
func RespondWithStatusCode(ctx *gin.Context, statusCode int) {
	RespondWithData(ctx, statusCode, nil)
}

// RespondWithErrorAndLog responds with the given error and logs it
func RespondWithErrorAndLog(ctx *gin.Context, respErr error) {
	if statusCodeErr, ok := respErr.(StatusCode); ok && statusCodeErr.StatusCode() >= http.StatusInternalServerError {
		log.ContextLogger(ctx).WithFields(log.Fields{
			"request": extractRequestData(ctx),
			"error":   respErr.Error(),
		}).Error("Internal Server Error")
	}

	RespondWithError(ctx, respErr)
}

// RespondWithError responds with the given error
func RespondWithError(ctx *gin.Context, respErr error) {
	statusCode := http.StatusInternalServerError
	if statusCodeErr, ok := respErr.(StatusCode); ok {
		statusCode = statusCodeErr.StatusCode()
	}

	response := ErrorResponse{}

	if errorCodeErr, ok := respErr.(ErrorCode); ok {
		response.Code = errorCodeErr.ErrorCode()
	} else {
		respErr = NewPrivateError(respErr)
	}

	if errorFields, ok := respErr.(ErrorFields); ok {
		response.Fields = errorFields.ErrorFields()
	}

	response.Error = respErr.Error()
	ctx.AbortWithStatusJSON(statusCode, response)
}

// extractRequestData extracts the request data from the gin context
func extractRequestData(ctx *gin.Context) *RequestForLog {
	return &RequestForLog{
		RequestID: ctx.GetString(log.RequestIDHeaderName),
		Method:    ctx.Request.Method,
		URL:       ctx.Request.URL.String(),
		Host:      ctx.Request.Host,
		ClientIP:  ctx.ClientIP(),
	}
}
