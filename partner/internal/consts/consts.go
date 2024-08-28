package consts

import (
	"errors"
	"os"
	"time"
)

const (
	//DatabaseType declared type of database
	DatabaseType = "postgres"
	// AppName represents the name of the application (e.g., "subscription").
	AppName = "partner"
	// AcceptedVersions represents the accepted API versions for the application (e.g., "v1.0").
	AcceptedVersions = "v1.0"
)

const (
	ContextEndPoints      = "context-endpoints"
	ContextErrorResponses = "context-error-response"
)

const (
	FailureKey = "failure"
	SuccessKey = "success"
)

// tablenames
const (
	PartnerTable                = "partner"
	PartnerOauthCredentialTable = "partner_oauth_credential"
	LanguageTable               = "language"
	CurrencyTable               = "currency"
	CountryTable                = "country"
	StateTable                  = "country_state"
	SubscriptioDurationTable    = "subscription_duration"
	PartnerPlanTable            = "partner_plan"
	CountryStateTable           = "country_state"
	LookupTable                 = "lookup"
	BusinessModelTable          = "lookup"
	PaymentGatewayTable         = "payment_gateway"
	PartnerPaymentGatewayTable  = "partner_payment_gateway"
	OauthProviderTable          = "oauth_provider"
)

// table fields
const (
	CodeField        = "code"
	IsoField         = "iso"
	IDKey            = "id"
	GracePeriodField = "member_grace_period"
	CountryCodeField = "country_code"
)

// keys
const (
	NameKey                                       = "name"
	URLKey                                        = "url"
	LogoKey                                       = "logo"
	ThemeKey                                      = "theme_id"
	WebsiteURLKey                                 = "website_url"
	BrowserTitleKey                               = "browser_title"
	BusinessModelKey                              = "business_model"
	ContactPersonKey                              = "contact_person"
	EmailKey                                      = "email"
	NoReplyEmailKey                               = "noreply_email"
	SupportEmailKey                               = "support_email"
	FeedbackEmailKey                              = "feedback_email"
	AlbumReviewEmailKey                           = "album_review_email"
	AddressKey                                    = "address"
	CountryKey                                    = "country"
	FaviconKey                                    = "favicon"
	PaymentGateWayIDKey                           = "payment_gateway_id"
	RequiredKey                                   = "Required"
	ProfileURLKey                                 = "profile_url"
	PaymentURLKey                                 = "payment_url"
	LandingPageKey                                = "landing_page"
	BackgroundImageKey                            = "background_image"
	PayoutCurrencyKey                             = "payout_currency"
	PlanIDKey                                     = "plan_id"
	FreePlanLimitKey                              = "free_plan_limit"
	ExpiryWarningCountKey                         = "expiry_warning_count"
	MaxRemittancePerMonthKey                      = "max_remittance_per_month"
	OutletsProcessingDurationKey                  = "outlets_processing_duration"
	BackgroundColorKey                            = "background_color"
	LanguageKey                                   = "language"
	PayoutTargetCurrencyKey                       = "payout_target_currency"
	StateKey                                      = "state"
	DefaultSubscriptionCurrencyKey                = "default_subscription_currency"
	DefaultPriceCodeCurrencyKey                   = "default_price_code_currency"
	MusicLanguageKey                              = "music_language"
	MemberDefaultCountryKey                       = "member_default_country"
	MemberGracePeriodKey                          = "member_grace_period"
	ProductReviewKey                              = "product_review"
	LoginTypeKey                                  = "login_type"
	StreetKey                                     = "street"
	CityKey                                       = "city"
	PostalCodeKey                                 = "postal_code"
	GatewayKey                                    = "payment_gateway"
	PaymentGatewayEmailKey                        = "payment_gateway_email"
	DefaultPayoutCurrencyKey                      = "default_payout_currency"
	DefaultPayinCurrencyKey                       = "default_payin_currency"
	SiteInfoKey                                   = "site_info"
	TermsAndConditionsNameKey                     = "terms_and_conditions_name"
	TermsAndConditionsDescriptionKey              = "terms_and_conditions_description"
	TermsAndConditionsLanguageKey                 = "terms_and_conditions_language"
	DefaultPaymentGatewayKey                      = "default_payment_gateway"
	MobileVerifyIntervalKey                       = "mobile_verify_interval"
	PartnerIDKey                                  = "partner_id"
	MemberIdKey                                   = "member_id"
	PageKey                                       = "page"
	PaymentGatewayKey                             = "payment_gateway"
	ClientIdKey                                   = "client_id"
	ClientSecretKey                               = "client_secret"
	EncryptionKey                                 = "encryption_key"
	SortKey                                       = "sort"
	OrderKey                                      = "order"
	StatusKey                                     = "status"
	LimitKey                                      = "limit"
	PayoutMinLimitKey                             = "payout_min_limit"
	ValidationErr                                 = "validation_error"
	ForbiddenErr                                  = "forbidden"
	UnauthorisedErr                               = "unauthorized"
	NotFoundKey                                   = "not_found"
	InternalServerErr                             = "internal_server_error"
	InvalidKey                                    = "invalid"
	AlreadyExistsKey                              = "already_exists"
	LimitExceedsKey                               = "Limit_exceeds"
	GenreIdKey                                    = "genre_id"
	RoleIdKey                                     = "role_id"
	OauthProviderKey                              = "oauth_provider"
	PartnerBaseUrlKey                             = "base_url_partner"
	PartnerNameKey                                = "partner_name"
	GenreNameKey                                  = "genre_name"
	RoleNameKey                                   = "role_name"
	DateTimeKey                                   = "date_time"
	StoresKey                                     = "stores"
	HeadingKey                                    = "heading"
	ContentKey                                    = "content"
	PaymentDetailsKey                             = "payment_details"
	LanguageCodeKey                               = "language_code"
	IsActiveKey                                   = "is_active"
	UpdatedByKey                                  = "updated_by"
	UpdatedOnKey                                  = "updated_on"
	ContactDetailsKey                             = "contact_details"
	SubscriptionDetailsKey                        = "subscription_details"
	PlanStartDateKey                              = "plan_start_date"
	PlanLaunchDateKey                             = "plan_launch_date"
	PlanIdKey                                     = "plan_id"
	AddressDetailsKey                             = "address_details"
	PaymentKey                                    = "payment"
	DefaultCurrencyKey                            = "default_currency"
	DescriptionKey                                = "description"
	FieldKey                                      = "fields"
	ProductTypeKey                                = "product_type"
	TrackFileQualityKey                           = "track_file_quality"
	IsPartnerExistMethod                          = "get"
	IsPartnerExistEndpoint                        = "partners"
	PartnerTermsAndConditionsUpdatedActvityLogKey = "partner_terms_and_conditions_updated"
	PartnerUpdatedActivityLogKey                  = "partner_updated"
	PartnerGenreDeletedActivityLog                = "partner_genre_deleted"
	PartnerRoleDeletedActivityLog                 = "partner_role_deleted"
	PartnerStoreCreatedActivityLog                = "partner_store_created"
	ProductReviewLookupTypeName                   = "product_review"
	LoginLookupTypeName                           = "partner_login_type"
	BusinessLookupTypeName                        = "business_model"
)

const (
	ClientIdLength               = 32
	ClientSecretLength           = 32
	ByteSize                     = 32
	LogMaxAge                    = 7
	LogMaxSize                   = 1024 * 1024 * 10
	LogMaxBackup                 = 5
	MaxExpiryWarningCount        = 20
	MaxFreePlanLimit             = 100
	MaxRemittancePerMonth        = 20
	MinOutletsProcessingDuration = 7
	MaxOutletsProcessingDuration = 30
	PartnerNameMaxLength         = 120
	ContactPersonMaxLength       = 60
	BrowserTitleMaxLength        = 500
	AddressMaxLength             = 250
	StreetMaxLength              = 120
	CityMaxLength                = 120
	SiteInfoMaxLength            = 2040
	URLMaxLength                 = 500
	EmailMaxLength               = 120
	GetPartnerByIDGoroutinesNum  = 9
	CacheExpiryTime              = 1 * time.Minute
)

// default values for partner data
const (
	LoginPageLogoDefaultVal               = "LoginPageLogoDefaultvalue"
	LoaderDefaultVal                      = "LoaderDefaultvalue"
	BackgroundColorDefaultVal             = "#ffffff"
	BackgroundImageDefaultVal             = "https://background.com/image"
	LanguageDefaultVal                    = "en"
	BrowserTitleDefaultVal                = "browser title"
	MobileVerifyIntervalDefaultVal        = 1
	PayoutTargetCurrencyDefaultVal        = "INR"
	ThemeIdDefaultVal                     = 1
	LoginTypeIdDefaultVal                 = "normal"
	PayoutMinLimitDefaultVal              = 1
	PayoutCurrencyDefaultVal              = "USD"
	MaxRemittancePerMonthDefaultVal       = 2
	MemberGracePeriodDefaultVal           = "Quarterly"
	ExpiryWarningCountDefaultVal          = 3
	DefaultSubscriptionCurrencyDefaultVal = "INR"
	DefaultPriceCodeCurrencyDefaultVal    = "INR"
	MusicLanguageDefaultVal               = "en"
	MemberDefaultCountryDefaultVal        = "IN"
	OutletsProcessingDurationDefaultVal   = 10
	FreePlanLimitDefaultVal               = 0
	ProductReviewDefaultVal               = "Both"
	LogoDefaultVal                        = "https://logo.com/image"
	FaviconDefaultVal                     = "FaviconDefaultVal"
	Key                                   = "key"
	InternalVal                           = "internal"
	LimitDefaultVal                       = 10
	PageDefaultVal                        = 1
)

// Cache Keys
const (
	CacheErrorKey     = "ERROR_CACHE_KEY_LABEL"
	CacheEndpointsKey = "endpoints"
)

// error message
const (
	UpdatePartnerErrMsg                   = "UpdatePartner failed err=%s"
	PartnerOauthCredentialErrMsg          = "CreatePartner PartnerOauthCredentialGeneration failed  err=%s"
	GetPartnerOauthCredentialErrMsg       = "GetPartnerByOauthCredential failed, err = %s"
	CreatePartnerErrMsg                   = "CreatePartner failed err=%s"
	GetPartnerPaymentGatewaysErrMsg       = "GetPartnerPaymentGateways failed ,err=%s"
	GetAllPartnersErrMsg                  = "Get all Partners failed, err = %s"
	UpdateTermsAndConditionsErrMsg        = "Update Terms and conditions failed, err = %s"
	DeletePartnerErrMsg                   = "Delete partner failed,error=%s"
	GetAllTermsAndConditionErrMsg         = "Get all terms and conditions failed, err = %s"
	GetPartnerByIdErrMsg                  = "Get partner by ID failed, err = %s"
	InvalidEndpointErrMsg                 = "invalid endpoint %s"
	LocalizationModuleErrMsg              = "error occured while retrieving errors from localization module %s"
	LogErrMsg                             = "failed to retrieve logs, err=%s"
	ContextErrMsg                         = "failed to fetch error values from context"
	PartnerNotFoundErrMsg                 = "partner not found"
	PathParameterErrMsg                   = "path parameter error"
	GetPartnerStoresErrMsg                = "Get partner stores failed err=%s"
	CreatePartnerStoresErrMsg             = "CreatePartnerStores failed err=%s"
	GetPartnerProductTypesErrMsg          = "Get partner product types failed err=%s"
	GetPartnerTrackQualityErrMsg          = "Get partner track quality failed err=%s"
	DeletePartnerArtistRoleLanguageErrMsg = "DeletePartnerArtistRoleLanguage failed err=%s"
	DeletePartnerGenreLanguageErrMsg      = "DeletePartnerGenreLanguage failed err=%s"
	ActivityLogErrMsg                     = "Activity log failed: err=%s"
	UpdatePartnerStatusErrMsg             = "UpdatePartnerStatus failed err =%s "
	GetPartnerNameErrMsg                  = "GetPartnerName failed err=%s"
	IsPartnerExistErrMsg                  = "Is partner exists failed err=%s"
	InvalidMemberErrMsg                   = "Invalid member id"
	BindingErrorErrMsg                    = "binding error"
	QueryBindingErrorErrMsg               = "query parameter binding error"
	ParsingErrMsg                         = "error occured while parsing"
	QueryParsingError                     = "unable to bind query parameters"
	InvalidHeaderData                     = "invalid header data"
	NotFoundErrMsg                        = "Data not found"
	ServiceUnavailableErrMsg              = "service unvailable"
	MaximumRequestErrMsg                  = "cannot exceed maximum page limit"
)

// connection failure errors
var (
	ErrStoreServiceConnectionLost        = errors.New("failed to connect store service")
	ErrMemberServiceConnectionLost       = errors.New("failed to connect member service")
	ErrSubscriptionServiceConnectionLost = errors.New("failed to connect subscription service")
	ErrUtilityServiceConnectionLost      = errors.New("failed to connect utility service")
	ErrOauthServiceConnectionLost        = errors.New("failed to connect oauth service")
	ErrMaximumRequest                    = errors.New("cannot exceed maximum page limit")
)

// success message
const (
	GetPartnerPaymentGatewaysSuccessMsg       = "Partner's payment details retrieved successfully"
	GetAllPartnersSuccessMsg                  = "Partners listed successfully"
	GetPartnerOauthCredentialSuccessMsg       = "Partner oauth credentials retrieved successfully"
	CreatePartnerSuccessMsg                   = "Partner created successfully"
	GetPartnerByIdSuccessMsg                  = "Partner data retrieved successfully"
	GetAllTermsAndConditionSuccessMsg         = "Terms and conditions retrieved successfully"
	UpdateTermsAndConditionsSuccessMsg        = "Partner terms and conditions updated successfully"
	UpdatePartnerSuccessMsg                   = "Successfully updated the partner details"
	DeletePartnerSuccessMsg                   = "Partner deleted successfully"
	GetPartnerStoresSuccessMsg                = "Partner stores retrieved successfully"
	CreatePartnerStoreSuccessMsg              = "Partner stores created successfully"
	DeletePartnerArtistRoleLanguageSuccessMsg = "Partner artist role language deleted successfully"
	DeletePartnerGenreLanguageSuccessMsg      = "Partner genre language deleted successfully"
	UpdatePartnerStatusSuccessMsg             = "Partner status updated successfully"
	IsPartnerExistSuccessMsg                  = "partner exist"
)

const (
	TermsAndConditionsNameMaxLength = 250
	MaxLimit                        = 50
	ProcessingDuration              = 10
	ReviewProcessingDuration        = 10
)

const (
	HexadecimalBase       = 16
	Uint32BitSize         = 32
	ValidColorCodeLength3 = 3
	ValidColorCodeLength6 = 6
)

const (
	Seperator      = ","
	StatusActive   = "active"
	StatusInActive = "inactive"
	StatusAll      = "all"
	Ascending      = "asc"
	Descending     = "desc"
)

var (
	SupportedCurrencies = []string{"INR", "USD", "JPY"}
)

const (
	MaxURLRuneCount   = 2083
	MinURLRuneCount   = 3
	URLPath           = `((\/|\?|#)[^\s]*)`
	IP                = `(([0-9a-fA-F]{1,4}:){7,7}[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,7}:|([0-9a-fA-F]{1,4}:){1,6}:[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,5}(:[0-9a-fA-F]{1,4}){1,2}|([0-9a-fA-F]{1,4}:){1,4}(:[0-9a-fA-F]{1,4}){1,3}|([0-9a-fA-F]{1,4}:){1,3}(:[0-9a-fA-F]{1,4}){1,4}|([0-9a-fA-F]{1,4}:){1,2}(:[0-9a-fA-F]{1,4}){1,5}|[0-9a-fA-F]{1,4}:((:[0-9a-fA-F]{1,4}){1,6})|:((:[0-9a-fA-F]{1,4}){1,7}|:)|fe80:(:[0-9a-fA-F]{0,4}){0,4}%[0-9a-zA-Z]{1,}|::(ffff(:0{1,4}){0,1}:){0,1}((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])|([0-9a-fA-F]{1,4}:){1,4}:((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9]))`
	URLPort           = `(:(\d{1,5}))`
	URLSchema         = `((ftp|tcp|udp|wss?|https?):\/\/)`
	URLUsername       = `(\S+(:\S*)?@)`
	URLIP             = `([1-9]\d?|1\d\d|2[01]\d|22[0-3]|24\d|25[0-5])(\.(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])){2}(?:\.([0-9]\d?|1\d\d|2[0-4]\d|25[0-5]))`
	URLSubdomain      = `((www\.)|([a-zA-Z0-9]+([-_\.]?[a-zA-Z0-9])*[a-zA-Z0-9]\.[a-zA-Z0-9]+))`
	URLExp            = `^` + URLSchema + `?` + URLUsername + `?` + `((` + URLIP + `|(\[` + IP + `\])|(([a-zA-Z0-9]([a-zA-Z0-9-_]+)?[a-zA-Z0-9]([-\.][a-zA-Z0-9]+)*)|(` + URLSubdomain + `?))?(([a-zA-Z\x{00a1}-\x{ffff}0-9]+-?-?)*[a-zA-Z\x{00a1}-\x{ffff}0-9]+)(?:\.([a-zA-Z\x{00a1}-\x{ffff}]{1,}))?))\.?` + URLPort + `?` + URLPath + `?$`
	ValidTitlePattern = "^[a-zA-Z0-9_-]+$"
)

var ValidTLDs = []string{".com", ".gov", ".edu"}

const (
	InvalidPartner = "invalid partner"
	InvalidMember  = "invalid member"
	RoleNotExist   = "Role_id doesn't exist"
	GenreNotExist  = "Genre_id doesn't exist"
)

var (
	UtilityServiceURL                  = os.Getenv("UTILITY_SERVICE_URL")
	SubcriptionServiceApiURL           = os.Getenv("SUBSCRIPTION_SERVICE_URL")
	MemberServiceURL                   = os.Getenv("MEMBER_SERVICE_URL")
	PartnerServiceURL                  = os.Getenv("PARTNER_SERVICE_URL")
	OAuthServiceURL                    = os.Getenv("OAUTH_SERVICE_URL")
	StoreServiceURL                    = os.Getenv("STORE_SERVICE_URL")
	ActivityLogServiceURL              = os.Getenv("ACTIVITY_LOG_SERVICE_URL")
	ErrorLocalizationURL               = os.Getenv("LOCALISATION_SERVICE_URL")
	LoggerServiceURL                   = os.Getenv("LOGGER_SERVICE_URL")
	LoggerSecret                       = os.Getenv("LOGGER_SECRET")
	ErrorHelpLink                      = os.Getenv("ERROR_HELP_LINK")
	ClientCredentialEncryptionKey      = os.Getenv("CLIENT_CREDENTIAL_ENCRYPTION_KEY")
	PartnerPaymentGatewayEncryptionKey = os.Getenv("PARTNER_PAYMENT_GATEWAY_ENCRYPTION_KEY")
)

// cache labels
const (
	SubDurationNameCacheKey    = "subscription_name_"
	SubDurationIdCacheKey      = "subscription_id_"
	ThemeCacheKey              = "theme_id_"
	PaymentGatewayIdCacheKey   = "payment_gateway_id_"
	PaymentGatewayNameCacheKey = "payment_gateway_name_"
	OauthProviderNameCacheKey  = "oauth_provider_name_"
	StoreCacheKey              = "stores"
	LookupIdCacheKey           = "lookup_id_"
	LookupNameCacheKey         = "lookup_name_"
	CurrencyIdCacheKey         = "currency_id_"
	CurrencyNameCacheKey       = "currency_name_"
	CurrencyExistsCacheKey     = "currency_exists_"
	CountryExistsCacheKey      = "country_exists_"
	LanguageExistsCacheKey     = "language_exists_"
	CountryStateExistsCacheKey = "state_exists_"
)
