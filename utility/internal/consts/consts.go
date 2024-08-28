package consts

import (
	"errors"
	"os"
)

// Constants defining fundamental properties and settings of the application.
const (
	DatabaseType = "postgres"
	AppName      = "utility"
)

// Context setting keys used to manage and access specific values within the application context.
const (
	ContextEndPoints = "context-endpoints"
	ContextErr       = "Error occured while loading error from service"
)

// Cache Keys
const (
	CacheErrorKey     = "ERROR_CACHE_KEY_LABEL"
	CacheEndpointsKey = "endpoints"
	PaginationKey     = "pagination"
)

// logger informations
const (
	LogMaxAge    = 7
	LogMaxSize   = 1024 * 1024 * 10
	LogMaxBackup = 5
)

// Default identifier for generating language labels.
const (
	DefaultIdentifier = "Label"
	DefaultLimit      = 10
	MaxLimit          = 50
	DefaultPage       = 1
	DefaultOrder      = "asc"
	PaymentIdSize     = 32
	MaximumLength     = 60
	MinimumLength     = 3
)

// pagination
const (
	DefaultPageStr  = "1"
	DefaultLimitStr = "10"
)

// key names for retrieving errors.
const (
	NameKey             = "name"
	Required            = "required"
	GenreExists         = "genre_exists"
	Genre               = "Genre"
	Role                = "Role"
	Deleted             = "deleted"
	CommonMsg           = "common_message"
	InvalidRole         = "invalid_role"
	InvalidGenre        = "invalid_genre"
	RoleExists          = "role_exists"
	LanguageLabelExists = "language_label_exists"
	LanguageLabel       = "language_label"
)

// db error constants
const (
	Country         = "country_id"
	CountryCode     = "country_code"
	NotFound        = "not_found"
	CountryNotFound = "Please provide country id"
	IDKey           = "id"
)

// errors
var (
	ErrNotFound        = errors.New("no records")
	ErrCountryNotFound = errors.New("country id not found")
	ErrNotExist        = errors.New("not exists")
)

type (
	Status string
	Sort   string
	Search string
	State  string
)

const (
	Active   Status = "active"
	Inactive Status = "inactive"
	All      Status = "all"
	Empty    Status = ""
)

const (
	Title Sort = "name"
)

var LangStatus = map[Status]bool{
	Active:   true,
	Inactive: true,
	All:      true,
	Empty:    true,
}

type Field struct {
	Column, Order string
}

var (
	SortOptns     = map[Sort]Field{Title: {Column: "name", Order: "asc"}}
	StateSortOpns = map[Sort]Field{Title: {Column: "cs.name", Order: "asc"}}
)

func IsValidSort(filterOptns map[Sort]Field, key string) (Field, bool) {
	v, ok := filterOptns[Sort(key)]
	return v, ok
}

// Key Names
const (
	InternalServerErr   = "internal_server_error"
	ValidationErr       = "validation_error"
	SuccessKey          = "Success"
	Failure             = "Failure"
	BindingError        = "Binding error"
	MaximumRequestError = "error: Cannot exceed maximum page limit"
	GenreNotExist       = "Genre doesn't exist"
	GenreNotFound       = "Please provide genre id"
	RoleNotExist        = "Role doesn't exist"
	RoleNotFound        = "Please provide role id"
	NotExist            = "not_exist"
)

const (
	LabelIdentifier   = "Label"
	LengthKey         = "length"
	IsoLength         = 3
	CodeLength        = 2
	StateIsoLength    = 2
	CodeField         = "code"
	IsoField          = "iso"
	InvalidKey        = "invalid"
	CurrencyIdField   = "currency_id"
	LanguageIdField   = "language_id"
	GenreIdField      = "genre_id"
	RoleIdField       = "role_id"
	ThemeIdField      = "theme_id"
	LookupTypeIdField = "lookup_type_id"
	SortField         = "sort"
	ArgumentKey       = "arguments"
	PaymentIdField    = "payment_gateway_id"
)

const (
	CommaDelimiter = ","
	DecimalBase    = 10
	BitSize32      = 32
	BitSize64      = 64
)

// API endpoint URL for other services.
var (

	// LoggerServiceURL is the URL of the Logger Service.
	LoggerServiceURL = os.Getenv("UTILITY_LOGGER_SERVICE_URL")

	// LoggerSecret is the secret key for the Logger Service.
	LoggerSecret = os.Getenv("UTILITY_LOGGER_SECRET")

	// LocalisationServiceURL is the URL of the Localisation Service.
	LocalisationServiceURL = os.Getenv("UTILITY_LOCALISATION_SERVICE_URL")

	// ErrorHelpLink is the URL for the error help link.
	ErrorHelpLink = os.Getenv("UTILITY_ERROR_HELP_LINK")
)

// function names
const (
	CheckStateExistsIdentifier      = "CheckStateExists"
	GetCountriesIdentifier          = "GetCountries"
	GetStatesOfCountryIdentifier    = "GetStatesOfCountry"
	CheckCountryExistsIdentifier    = "CheckCountryExists"
	GetAllCountryCodesIdentifier    = "GetAllCountryCodes"
	GetAllCurrencyIdentifier        = "GetAllCurrency"
	GetCurrencyByIDIdentifier       = "GetCurrencyByID"
	GetCurrencyByISOIdentifier      = "GetCurrencyByISO"
	CreateGenreIdentifier           = "CreateGenre"
	GetGenresIdentifier             = "GetGenres"
	DeleteGenreIdentifier           = "DeleteGenre"
	UpdateGenreIdentifier           = "UpdateGenre"
	GetGenresByIDIdentifier         = "GetGenresByID"
	GetLanguagesIdentifier          = "GetLanguages"
	GetLanguageCodeExistsIdentifier = "GetLanguageCodeExists"
	GetLookupByTypeNameIdentifier   = "GetLookupByTypeName"
	GetLookupByIdListIdentifier     = "GetLookupByIdList"
	GetPaymentGatewayByIDIdentifier = "GetPaymentGatewayByID"
	GetRoleByIDIdentifier           = "GetRoleByID"
	GetRolesIdentifier              = "GetRoles"
	DeleteRolesIdentifier           = "DeleteRoles"
	CreateRoleIdentifier            = "CreateRole"
	UpdateRoleIdentifier            = "UpdateRole"
	GetThemeByIDIdentifier          = "GetThemeByID"
	GetAllPaymentGatewayIdentifier  = "GetAllPaymentGateway"
)
