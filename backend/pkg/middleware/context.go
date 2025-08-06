package middleware

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"social-app/internal/models"
)

type Context struct {
	Writer  gin.ResponseWriter
	root    *gin.Context
	Request *http.Request
	User    models.User
}

type HandlerFunc func(*Context)

func Verified(handler HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, done := setContext(c)
		if done {
			return
		}

		c.Set("context", ctx)

		if !ctx.User.Verified {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "user not verified"})
			return
		}

		handler(ctx)
	}
}

func Profile(handler HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := &Context{
			root:    c,
			Writer:  c.Writer,
			Request: c.Request,
		}

		c.Set("context", ctx)

		handler(ctx)
	}
}

func setContext(c *gin.Context) (*Context, bool) {
	user, done := currentAuth(c)
	if done {
		c.AbortWithStatusJSON(401, gin.H{"error": "unauthorized"})
		return nil, true
	}

	ctx := &Context{
		User:    user,
		root:    c,
		Writer:  c.Writer,
		Request: c.Request,
	}
	return ctx, false
}

func currentAuth(c *gin.Context) (models.User, bool) {
	user := models.User{
		ID:       c.GetUint64("user_id"),
		Role:     c.GetString("role"),
		Email:    c.GetString("email"),
		Username: c.GetString("username"),
		Verified: c.GetBool("verified"),
	}

	if user.ID == 0 {
		return user, true
	}

	return user, false
}

func (c *Context) JSON(code int, obj any) {
	c.root.JSON(code, obj)
}

func (c *Context) BindJSON(obj any) error {
	err := c.root.BindJSON(obj)
	if err != nil {
		return fmt.Errorf("failed to bind JSON: %w", err)
	}
	return nil
}

func (c *Context) Query(key string) string {
	return c.root.Query(key)
}

func (c *Context) Param(key string) string {
	return c.root.Param(key)
}

func (c *Context) GetUint(key string) uint {
	return c.root.GetUint(key)
}

func (c *Context) GetUint64(key string) (uint64, error) {
	id, err := strconv.Atoi(c.Param(key))
	if err != nil {
		return 0, fmt.Errorf("invalid ID format: %w", err)
	}
	if id < 0 {
		return 0, fmt.Errorf("ID must be a positive integer")
	}
	return uint64(id), nil
}

func (c *Context) FormFile(key string) (*multipart.FileHeader, error) {
	header, err := c.root.FormFile(key)
	if err != nil {
		return nil, fmt.Errorf("failed to get form file for key '%s': %w", key, err)
	}
	return header, nil
}

func (c *Context) MultipartForm() (*multipart.Form, error) {
	form, err := c.root.MultipartForm()
	if err != nil {
		return nil, fmt.Errorf("failed to get multipart form: %w", err)
	}
	return form, nil
}

func (c *Context) Next() {
	c.root.Next()
}

func (c *Context) NextEr() error {
	c.Next()
	if len(c.root.Errors) > 0 {
		return fmt.Errorf("errors occurred during request processing: %v", c.root.Errors)
	}
	return nil
}

func (c *Context) ShouldBindJSON(u any) error {
	err := c.root.ShouldBindJSON(u)
	if err != nil {
		return fmt.Errorf("failed to bind JSON: %w", err)
	}
	return nil
}

func (c *Context) QueryArray(s string) []string {
	return c.root.QueryArray(s)
}
