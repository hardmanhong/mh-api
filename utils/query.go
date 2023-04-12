package utils

import (
	"net/url"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func GetPaginationParams(ctx *gin.Context) (page int, pageSize int) {
	pageSizeStr := ctx.Query("pageSize")
	pageStr := ctx.Query("page")
	var err error
	pageSize, err = strconv.Atoi(pageSizeStr)
	if err != nil || pageSize <= 0 {
		pageSize = 10
	}
	page, err = strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1
	}
	return page, pageSize
}

func ParseDate(keys [2]string, values url.Values) (*time.Time, *time.Time) {
	result := []*time.Time{nil, nil}
	for i, key := range keys {
		strValue := values.Get(key)
		if strValue == "" {
			continue
		}
		t, err := time.Parse("2006-01-02", strValue)
		if err != nil {
			return nil, nil
		}
		if i == 0 {
			t = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
		} else {
			t = time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 0, t.Location())
		}
		result[i] = &t
	}
	return result[0], result[1]
}
