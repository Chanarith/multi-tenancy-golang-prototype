package middleware

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	DEFAULT_PAGE_TEXT    = "page"
	DEFAULT_SIZE_TEXT    = "size"
	DEFAULT_PAGE         = "1"
	DEFAULT_PAGE_SIZE    = "10"
	DEFAULT_MIN_PAGESIZE = 10
	DEFAULT_MAX_PAGESIZE = 100
)

type Paginator struct {
	pageText    string
	sizeText    string
	page        string
	pageSize    string
	minPageSize int
	maxPageSize int
}

func DefaultPaginator() *Paginator {
	return &Paginator{
		pageText:    DEFAULT_PAGE_TEXT,
		sizeText:    DEFAULT_SIZE_TEXT,
		page:        DEFAULT_PAGE,
		pageSize:    DEFAULT_PAGE_SIZE,
		minPageSize: DEFAULT_MIN_PAGESIZE,
		maxPageSize: DEFAULT_MAX_PAGESIZE,
	}
}

func (p *Paginator) Handle(c *gin.Context) {
	pageStr := c.DefaultQuery(p.pageText, DEFAULT_PAGE_TEXT)
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "page number must be an integer"})
		return
	}
	if page < 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "page number must be positive"})
		return
	}
	sizeStr := c.DefaultQuery(p.sizeText, DEFAULT_PAGE_SIZE)
	size, err := strconv.Atoi(sizeStr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "page size must be an integer"})
		return
	}
	c.Set(p.pageText, page)
	c.Set(p.sizeText, size)
	c.Next()
}
