package middlewares

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"utility/internal/consts"
	"utility/internal/entities"
	"utility/utilities"

	"github.com/gin-gonic/gin"
	"github.com/labstack/gommon/log"
	constants "gitlab.com/tuneverse/toolkit/consts"
	"gitlab.com/tuneverse/toolkit/core/logger"
	"gitlab.com/tuneverse/toolkit/models"
	"gitlab.com/tuneverse/toolkit/models/api"
	"gitlab.com/tuneverse/toolkit/utils"
)

type middlewares struct {
	Cfg *entities.EnvConfig
}

func NewMiddlewares(cfg *entities.EnvConfig) *middlewares {
	return &middlewares{
		Cfg: cfg,
	}
}

// QueryOptions defines optional parameters for the QueryParams middleware.
type QueryOptions struct {
	Key                                        string
	DefaultLimit, DefaultPage, MaxAllowedLimit int32
}

// QueryParams is a middleware for handling query parameters in HTTP GET requests.
func (m *middlewares) QueryParams(opts ...QueryOptions) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if ctx.Request.Method == http.MethodGet {
			var (
				params  entities.Pagination
				lg      = logger.Log().WithContext(ctx.Request.Context())
				apiResp = api.Response{
					Status:  "failure",
					Message: "validation error",
					Code:    http.StatusBadRequest,
					Data:    struct{}{},
					Errors:  struct{}{},
				}
			)
			opt := QueryOptions{}

			if len(opts) > 0 {
				opt = opts[0]
			}

			paginationKey := "pagination"
			if opt.Key != "" {
				paginationKey = opt.Key
			}

			if opt.DefaultLimit <= 0 {
				apiResp.Errors = errors.New("set default value for the page limit; should be greater than 0").Error()
				ctx.AbortWithStatusJSON(http.StatusBadRequest, apiResp)
				return
			}

			if opt.DefaultPage <= 0 {
				apiResp.Errors = errors.New("set default value for the page; should be greater than 0").Error()
				ctx.AbortWithStatusJSON(http.StatusBadRequest, apiResp)
				return
			}

			if opt.MaxAllowedLimit <= 0 {
				apiResp.Errors = errors.New("set default value for the maximum allowed limit; should be greater than 0").Error()
				ctx.AbortWithStatusJSON(http.StatusBadRequest, apiResp)
				return
			}

			limitStr := ctx.DefaultQuery("limit", consts.DefaultLimitStr)
			pageStr := ctx.DefaultQuery("page", consts.DefaultPageStr)

			helpLink := consts.ErrorHelpLink
			httpMethod, endpointURL := strings.ToLower(ctx.Request.Method), ctx.FullPath()
			contextEndpoints, isEndpointExists := utils.GetContext[[]models.DataItem](ctx, consts.ContextEndPoints)
			contextError, contextStatus := utils.GetContext[map[string]any](ctx, constants.ContextErrorResponses)
			if !isEndpointExists {
				val, _, errorCode, errDet := utils.ParseFields(ctx, consts.InternalServerErr,
					"", contextError, "", "", nil, helpLink)
				log.Errorf("Invalid endpoint %s", val)
				result := api.Response{Status: consts.Failure, Message: errDet.Message, Code: int(errorCode), Data: nil, Errors: map[string]string{}}
				ctx.JSON(http.StatusBadRequest, result)
				return
			}

			endpoint := utils.GetEndPoints(contextEndpoints, endpointURL, httpMethod)
			if !contextStatus {
				lg.Error("QueryParams failed, Error while retrieving data from the context")
				return
			}

			errMap := utilities.NewErrorMap()

			page, err := strconv.Atoi(pageStr)
			if err != nil || page <= 0 {
				//utils.AppendValuesToMap(validationErrs, "page", "invalid")

				code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, httpMethod, "page", "invalid")
				if err != nil {
					logger.Log().WithContext(ctx).Errorf("ListActivityLogs: Failed,error while loading service code")
				}

				errMap["page"] = models.ErrorResponse{
					Code:    code,
					Message: []string{"invalid"},
				}

			}

			limit, err := strconv.Atoi(limitStr)
			if err != nil || limit <= 0 {
				//utils.AppendValuesToMap(validationErrs, "limit", "invalid")

				code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, httpMethod, "limit", "invalid")
				if err != nil {
					logger.Log().WithContext(ctx).Errorf("ListActivityLogs: Failed,error while loading service code")
				}

				errMap["limit"] = models.ErrorResponse{
					Code:    code,
					Message: []string{"invalid"},
				}
			}

			params.Limit = int32(limit)
			params.Page = int32(page)

			// service based codes are extracted here
			serviceCode := utilities.ExtractServicecode(errMap)

			if len(errMap) > 0 {
				fields := utils.FieldMapping(errMap)
				val, _, errorCode, _ := utils.ParseFields(ctx, constants.ValidationErr,
					fields, contextError, endpoint, httpMethod, serviceCode, helpLink)

				apiResp.Code = int(errorCode)
				apiResp.Errors = val
				ctx.AbortWithStatusJSON(http.StatusBadRequest, apiResp)
				return
			}

			if params.Limit > opt.MaxAllowedLimit {

				result := utilities.ErrorResponseGenerator("Too many request", http.StatusTooManyRequests, consts.MaximumRequestError)

				ctx.AbortWithStatusJSON(http.StatusTooManyRequests, result)
				return
			}
			params.Page, params.Limit = utils.Paginate(params.Page, params.Limit, opt.DefaultLimit)
			ctx.Set(paginationKey, params)
		}

		ctx.Next()
	}
}
