package sapi

import (
	"fmt"
	"github.com/s3ndd/sen-go/log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// defaultPageSize is the default page size for paginated requests
const defaultPageSize = 100

// RequestBody binds the request body to the given struct
func RequestBody(ctx *gin.Context, request interface{}) error {
	if err := ctx.ShouldBindJSON(request); err != nil {
		log.Global().WithField("req", request).Error(fmt.Sprintf("failed request, %s", err))
		return NewValidationError(err.Error(), http.StatusBadRequest)
	}
	return nil
}

// RequestPagination returns the page index and page size from the request
func RequestPagination(ctx *gin.Context) (int, int) {
	pageIndex, err := strconv.Atoi(ctx.Query("page_index"))
	if err != nil {
		pageIndex = 1
	}

	pageSize, err := strconv.Atoi(ctx.Query("page_size"))
	if err != nil {
		pageSize = defaultPageSize
	}

	return pageIndex, pageSize
}

// RequestIncludeDeleted returns true if the request includes deleted items
func RequestIncludeDeleted(ctx *gin.Context) bool {
	deleted, err := strconv.ParseBool(ctx.Query("deleted"))
	if err != nil {
		return false
	}
	return deleted
}
