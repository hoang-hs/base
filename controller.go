package base

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/hoang-hs/base/log"
	"net/http"
	"strings"
)

type Controller struct {
	validate *validator.Validate
}

func NewBaseController(validate *validator.Validate) *Controller {
	return &Controller{
		validate: validate,
	}
}

func (b *Controller) Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, data)
}

func (b *Controller) ErrorData(c *gin.Context, err *Error) {
	c.JSON(err.GetHttpStatus(), err)
}

func (b *Controller) BindAndValidateRequest(c *gin.Context, req interface{}) *Error {
	if err := c.BindUri(req); err != nil {
		log.Warn(c, "bind request err, err:[%s]", err)
		return ErrBadRequest(c).SetDetail(err.Error())
	}
	if err := c.Bind(req); err != nil {
		log.Warn(c, "bind request err, err:[%s]", err)
		return ErrBadRequest(c).SetDetail(err.Error())
	}
	return b.ValidateRequest(c, req)
}

func (b *Controller) ValidateRequest(ctx context.Context, req interface{}) *Error {
	err := b.validate.Struct(req)

	if err != nil {
		var errs validator.ValidationErrors
		if !errors.As(err, &errs) {
			log.Error(ctx, "Cannot parse validate error: %+v", err)
			return ErrSystemError(ctx, "ValidateFailed").SetDetail(err.Error())
		}
		var filedErrors []string
		for _, errValidate := range errs {
			log.Debug(ctx, "field invalid, err:[%s]", errValidate.Field())
			filedErrors = append(filedErrors, errValidate.Error())
		}
		str := strings.Join(filedErrors, ",")
		log.Warn(ctx, "invalid request, err:[%s]", err.Error())
		return ErrBadRequest(ctx).SetDetail(fmt.Sprintf("field invalidate [%s]", str))
	}
	return nil
}
