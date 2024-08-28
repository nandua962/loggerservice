package utilities

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"utility/internal/consts"
	"utility/internal/entities"

	"github.com/gin-gonic/gin"
	"github.com/labstack/gommon/log"
	"gitlab.com/tuneverse/toolkit/core/logger"
	"gitlab.com/tuneverse/toolkit/models"
	"gitlab.com/tuneverse/toolkit/models/api"
	"gitlab.com/tuneverse/toolkit/utils"
)

// IsEndpointExists to check endpoint exists or not
func IsEndpointExists(ctx *gin.Context, isEndpointExists bool, contextError map[string]any, helpLink string, contextStatus bool, funcName string) {

	if !isEndpointExists {
		val, _, errorCode, errDet := utils.ParseFields(ctx, consts.InternalServerErr,
			"", contextError, "", "", nil, helpLink)
		log.Errorf("%s Specific Permissions failed: invalid endpoint %s", funcName, val)
		result := api.Response{Status: consts.Failure, Message: errDet.Message, Code: int(errorCode), Data: nil, Errors: map[string]string{}}
		ctx.JSON(http.StatusBadRequest, result)
		return
	}
	if !contextStatus {
		log.Errorf("%s failed, Error while retrieving data from the context", funcName)
		return
	}
}

// HandleError handles errors in a Gin context and generates an appropriate JSON response.
func HandleNotFoundError(ctx *gin.Context, message string, code int, errors any) {

	result := ErrorResponseGenerator(message, code, errors)
	ctx.JSON(http.StatusNotFound, result)
}

// HandleError handles errors in a Gin context and generates an appropriate JSON response.
func HandleError(ctx *gin.Context, err error, validation entities.Validation) {

	var (
		val       any
		errorCode float64
		hasError  bool
		errDet    models.ErrorDetails
	)
	ctxt := ctx.Request.Context()
	switch {
	case errors.Is(err, consts.ErrNotFound):
		logger.Log().WithContext(ctxt).Info("HandleError: no data found matching the request")
		val, hasError, errorCode, errDet = utils.ParseFields(ctx, consts.NotFound,
			"", validation.ContextError, "", "", nil, validation.HelpLink)

	default:
		logger.Log().WithContext(ctxt).Info("HandleError: failed")

		val, hasError, errorCode, errDet = utils.ParseFields(ctx, consts.InternalServerErr,
			"", validation.ContextError, "", "", nil, validation.HelpLink)

	}

	if !hasError {
		logger.Log().WithContext(ctxt).Info("HandleError: Error while parsing data : %s", val)
	}

	result := api.Response{Status: consts.Failure, Message: errDet.Message, Code: int(errorCode), Data: map[string]string{}, Errors: val}
	ctx.JSON(result.Code,
		result,
	)

}

// HandleError handles errors in a Gin context and generates an appropriate JSON response.
func HandleInternalServerError(ctx *gin.Context, validation entities.Validation, err error) {
	var (
		val       any
		errorCode float64
		hasError  bool
		errDet    models.ErrorDetails
	)

	ctxt := ctx.Request.Context()

	val, hasError, errorCode, errDet = utils.ParseFields(ctx, consts.InternalServerErr,
		"", validation.ContextError, "", "", nil, validation.HelpLink)

	if !hasError {
		logger.Log().WithContext(ctxt).Info("[HandleInternalServerError]: Error while parsing data : %s", val)
	}

	result := api.Response{Status: consts.Failure, Message: errDet.Message, Code: int(errorCode), Data: map[string]string{}, Errors: val}
	ctx.JSON(result.Code, result)
}

// Order validates the order parameter for SQL queries.
func Order(ctx context.Context, order string, endpoint string, method string, errs map[string]models.ErrorResponse) {
	switch strings.ToUpper(order) {
	case "ASC", "DESC", "":
	default:

		code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, "order", "invalid")
		if err != nil {
			logger.Log().WithContext(ctx).Errorf("Order : Failed,%s", err)
		}

		errs["order"] = models.ErrorResponse{
			Code:    code,
			Message: []string{"invalid"},
		}

	}
}

// GroupBy constructs a SQL GROUP BY clause from the given columns.
func GroupBy(sql string, columns ...string) string {

	if len(columns) > 0 {
		sql += " GROUP BY "
	}
	for i, column := range columns {
		sql = fmt.Sprintf("%s %s", sql, column)
		if i < len(columns)-1 {
			sql += ","
		}
	}
	return sql
}

// Search constructs a SQL query for search functionality.
func Search(sql, value string, columns ...string) string {

	if strings.HasPrefix(value, "~") {
		sql = fmt.Sprintf("%s %s", AddCondition(sql), utils.BuildSearchQuery(value[1:], columns...))
	} else {
		sql = ExactMatch(sql, value, columns...)
	}
	return sql
}

// SearchIso constructs a SQL query for exact match searches for ISO values.
func SearchIso(sql, value string, columns ...string) string {

	sql = ExactMatch(sql, value, columns...)

	return sql
}

// SearchIdList constructs a SQL query to search for a list of IDs in the specified column.
func SearchIdList(sql string, id string, columnName string) string {
	sql = AddCondition(sql)
	idList := strings.Split(id, ",")
	var query []string
	for _, value := range idList {
		query = append(query, fmt.Sprintf("%s = '%s'", columnName, value))

	}

	result := strings.Join(query, " OR ")

	return sql + result
}

// SearchById constructs a SQL query for exact match searches for an ID in the specified columns.
func SearchById(sql, value string, columns ...string) string {

	return exactMatchs(sql, value, columns...)
}

// exactMatchs constructs a SQL query for exact match searches for a value in the specified columns.
func exactMatchs(sql string, value string, columns ...string) string {
	sql = AddCondition(sql)
	var searchQuery strings.Builder
	for i, column := range columns {
		_, err := searchQuery.WriteString(fmt.Sprintf(" %s %s = '%s' ", sql, column, value))
		_ = err
		if i < len(columns)-1 {
			_, err := searchQuery.WriteString("OR ")
			_ = err
		}
	}
	return searchQuery.String()
}

// AddCondition adds a WHERE or AND clause to the SQL query based on its current state.
func AddCondition(query string) string {
	if strings.Contains(query, "WHERE") {
		return query + " AND "
	}
	return query + " WHERE "
}

// ExactMatch constructs a SQL query for case-insensitive exact match searches for a value in the specified columns.
func ExactMatch(sql string, value string, columns ...string) string {
	sql = AddCondition(sql)
	var searchQuery strings.Builder
	for i, column := range columns {
		_, err := searchQuery.WriteString(fmt.Sprintf(" %s LOWER(%s) = LOWER('%s') ", sql, column, value))
		_ = err
		if i < len(columns)-1 {
			_, err := searchQuery.WriteString("OR ")
			_ = err
		}
	}
	return searchQuery.String()
}

// OrderBy constructs a SQL ORDER BY clause from the given column and order parameters.
func OrderBy(ctx context.Context, column, order string, allowedFilters map[consts.Sort]consts.Field, endpoint string, method string, errs map[string]models.ErrorResponse) (string, error) {

	var columns, orders []string
	var log = logger.Log().WithContext(ctx)
	//if both sorting options and sorting order are empty set the value to default
	if column == "" {
		// columns = append(columns, defaultColumn)
		columns = append(columns, consts.NameKey)
	} else {
		columns = strings.Split(column, ",")
	}

	if order == "" {
		orders = append(orders, consts.DefaultOrder)
	} else {
		orders = strings.Split(order, ",")
	}

	columnsLength := len(columns)

	if columnsLength > 0 && len(orders) == 0 {
		for _, v := range columns {
			if filter, ok := allowedFilters[consts.Sort(v)]; ok {
				orders = append(orders, filter.Order)
			}
		}
	}

	orderLength := len(orders)

	if columnsLength != orderLength {

		code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.SortField, consts.ArgumentKey)

		if err != nil {
			log.Errorf("OrderBy : Failed, Error in sort arguments %s", err)
		}
		errs[consts.SortField] = models.ErrorResponse{
			Code:    code,
			Message: []string{consts.ArgumentKey},
		}

		return "", nil
	}

	sql := "ORDER BY"
	for i, c := range columns {
		Order(ctx, orders[i], endpoint, method, errs)

		if len(errs) > 0 {
			break
		}
		v, ok := consts.IsValidSort(allowedFilters, c)
		if !ok {

			code, err := utils.GetErrorCode(consts.ErrorCodeMap, endpoint, method, consts.SortField, consts.InvalidKey)
			if err != nil {
				log.Errorf("OrderBy : Failed, Error in invalid sort  %s", err)
			}
			errs[consts.SortField] = models.ErrorResponse{
				Code:    code,
				Message: []string{consts.InvalidKey},
			}

			break
		}

		sql += fmt.Sprintf(" %s %s", v.Column, orders[i])
		if i < columnsLength-1 {
			sql += consts.CommaDelimiter
		}

	}

	return sql, nil
}

// For success response
func SuccessResponseGenerator(message string, code int, data any) api.Response {
	var result api.Response
	if data == "" {
		data = map[string]interface{}{}
	}
	result = api.Response{Status: consts.SuccessKey, Message: message, Code: code, Data: data, Errors: map[string]string{}}
	return result
}

// For error response
func ErrorResponseGenerator(message string, code int, errors any) api.Response {
	var result api.Response
	if errors == "" {
		errors = map[string]interface{}{}
	}
	result = api.Response{Status: consts.Failure, Message: message, Code: code, Data: map[string]string{}, Errors: errors}
	return result
}

// Extracts service codes from the ErrorResponse objects and returns a map with the service codes.
func ExtractServicecode(errMap map[string]models.ErrorResponse) map[string]string {

	serviceCode := make(map[string]string)
	for key, value := range errMap {
		serviceCode[key] = value.Code
	}

	return serviceCode
}

// Validates if a string contains only alphabetic characters (A-Z, a-z).
func IsValidValue(value string) bool {

	// Define a regular expression pattern that allows A-Z,a-z .
	validPattern := "^[a-zA-Z]+$"

	regex := regexp.MustCompile(validPattern)

	return regex.MatchString(value)
}

// Validates if a string contains only alphabetic characters (A-Z, a-z and  _ underscore).
func IsValidValueWithUnderscore(value string) bool {

	// Define a regular expression pattern that allows A-Z,a-z .
	validPattern := "^[a-zA-Z_]+$"

	regex := regexp.MustCompile(validPattern)

	return regex.MatchString(value)
}

// Validates if a string contains only numeric characters (0-9).
func IsValidID(value string) bool {

	// Define a regular expression pattern that allows only numeric value.
	validPattern := "^[0-9]+$"

	regex := regexp.MustCompile(validPattern)

	return regex.MatchString(value)
}

// Validatecode validates specified parameters.
func Validatecode(ctx context.Context, validation entities.Validation, length int16, field string, errs map[string]models.ErrorResponse) map[string]models.ErrorResponse {

	var log = logger.Log().WithContext(ctx)
	code := validation.ID
	if len(code) != int(length) {

		code, err := utils.GetErrorCode(consts.ErrorCodeMap, validation.Endpoint, validation.Method, field, consts.LengthKey)

		if err != nil {
			log.Errorf(" : Failed, Error in iso length %s", err)
		}
		errs[field] = models.ErrorResponse{
			Code:    code,
			Message: []string{"length"},
		}

	} else {

		isISOValid := IsValidValue(code)

		if !isISOValid {
			code, err := utils.GetErrorCode(consts.ErrorCodeMap, validation.Endpoint, validation.Method, field, consts.InvalidKey)

			if err != nil {
				log.Errorf(" : Failed, Error in iso validation %s", err)
			}
			errs[field] = models.ErrorResponse{
				Code:    code,
				Message: []string{consts.InvalidKey},
			}
		}

	}

	return errs
}

func HandleValidationError(ctx *gin.Context, errMap map[string]models.ErrorResponse, validation entities.Validation) {

	log := logger.Log().WithContext(ctx)
	// service based codes are extracted here
	serviceCode := ExtractServicecode(errMap)
	fields := utils.FieldMapping(errMap)
	val, hasError, errorCode, errDet := utils.ParseFields(ctx, consts.ValidationErr,
		fields, validation.ContextError, validation.Endpoint, validation.Method, serviceCode, validation.HelpLink)
	if hasError {
		log.Errorf("[HandleValidationError] Error : %s", val)
	}
	result := api.Response{Status: consts.Failure, Message: errDet.Message, Code: int(errorCode), Data: map[string]string{}, Errors: val}
	ctx.JSON(http.StatusBadRequest, result)
}

// IDValidation validates ID parameters.
func IDValidation(ctx context.Context, validation entities.Validation, errMap map[string]models.ErrorResponse, field string) map[string]models.ErrorResponse {

	var log = logger.Log().WithContext(ctx)

	isIDValid := IsValidID(validation.ID)

	if !isIDValid {
		code, err := utils.GetErrorCode(consts.ErrorCodeMap, validation.Endpoint, validation.Method, field, consts.InvalidKey)
		if err != nil {
			log.Errorf("[IsValidID] Error while loading service code Error : %s", err.Error())
		}
		errMap[field] = models.ErrorResponse{
			Code:    code,
			Message: []string{consts.InvalidKey},
		}
	}
	return errMap
}

// LookUPIDValidation validates Lookup Ids parameters.
func LookUPIDValidation(ctx context.Context, validation entities.Validation, errMap map[string]models.ErrorResponse, field string) map[string]models.ErrorResponse {

	var log = logger.Log().WithContext(ctx)

	isIDValid := IsValidID(validation.ID)

	if !isIDValid {
		code, err := utils.GetErrorCode(consts.ErrorCodeMap, validation.Endpoint, validation.Method, field, consts.InvalidKey)
		if err != nil {
			log.Errorf("[IsValidID] Error while loading service code Error : %s", err.Error())
		}
		errMap[field] = models.ErrorResponse{
			Code:    code,
			Message: []string{consts.InvalidKey},
		}
	}
	return errMap
}
