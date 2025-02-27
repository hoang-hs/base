package controller

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	errorCustom "github.com/hoang-hs/base/src/common"
	"github.com/hoang-hs/base/src/common/log"
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

func (b *Controller) ErrorData(c *gin.Context, err *errorCustom.Error) {
	c.JSON(err.GetHttpStatus(), err)
}

func (b *Controller) BindAndValidateRequest(c *gin.Context, req interface{}) *errorCustom.Error {
	if err := c.BindUri(req); err != nil {
		log.WarnCtx(c, "bind request err", log.Err(err))
		return errorCustom.ErrBadRequest(c).SetDetail(err.Error())
	}
	if err := c.Bind(req); err != nil {
		log.WarnCtx(c, "bind request err", log.Err(err))
		return errorCustom.ErrBadRequest(c).SetDetail(err.Error())
	}
	return b.ValidateRequest(c, req)
}

func (b *Controller) ValidateRequest(ctx context.Context, req interface{}) *errorCustom.Error {
	err := b.validate.Struct(req)

	if err != nil {
		var errs validator.ValidationErrors
		if !errors.As(err, &errs) {
			log.ErrorCtx(ctx, "Cannot parse validate", log.Err(err))
			return errorCustom.ErrSystemError(ctx, "ValidateFailed").SetDetail(err.Error())
		}
		var filedErrors []string
		for _, errValidate := range errs {
			log.DebugCtx(ctx, "field invalid", log.String("err", errValidate.Field()))
			filedErrors = append(filedErrors, errValidate.Error())
		}
		str := strings.Join(filedErrors, ",")
		log.WarnCtx(ctx, "invalid request", log.Err(err))
		return errorCustom.ErrBadRequest(ctx).SetDetail(fmt.Sprintf("field invalidate [%s]", str))
	}
	return nil
}
