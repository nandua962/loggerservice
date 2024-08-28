package entities

// general error response format retireved from error localization
type ErrorResponse struct {
	Message   string                 `json:"message"`
	ErrorCode interface{}            `json:"errorCode"`
	Errors    map[string]interface{} `json:"errors"`
}

// structure to get partner oauth credential details
type GetPartnerOauthCredential struct {
	ClientId            string   `json:"client_id"`
	ClientSecret        string   `json:"client_secret"`
	RedirectUri         string   `json:"redirect_uri"`
	Scope               []string `json:"scope"`
	AccessTokenEndpoint string   `json:"access_token_endpoint"`
	ProviderName        string   `json:"oauth_provider_name"`
}

// structure to get header data for  partner oauth credential
type PartnerOAuthHeader struct {
	ProviderName    string `header:"provider_name"`
	ClientID        string `header:"client_id"`
	ClientSecret    string `header:"client_secret"`
	OauthProviderID string `json:"oauth_provider_id"`
}

// structure to define the conditions for retrieving Partner Oauth credential
type PartneOauthConditions struct {
	PartnerId string
}

// structure to define partner details
type Partner struct {
	Name                        string                  `json:"name"`
	URL                         string                  `json:"url"`
	Logo                        string                  `json:"logo"`
	Favicon                     string                  `json:"favicon"`
	Loader                      string                  `json:"loader"`
	LoginPageLogo               string                  `json:"login_page_logo"`
	ContactDetails              ContactDetails          `json:"contact_details"`
	BackgroundColor             string                  `json:"background_color"`
	BackgroundImage             string                  `json:"background_image"`
	WebsiteURL                  string                  `json:"website_url"`
	Language                    string                  `json:"language,omitempty"`
	BrowserTitle                string                  `json:"browser_title"`
	ProfileURL                  string                  `json:"profile_url"`
	PaymentURL                  string                  `json:"payment_url"`
	LandingPage                 string                  `json:"landing_page"`
	MobileVerifyInterval        int                     `json:"mobile_verify_interval"`
	TermsAndConditionsVersionID int                     `json:"terms_and_conditions_version_id"`
	PayoutTargetCurrency        string                  `json:"payout_target_currency"`
	PayoutTargetCurrencyID      int                     `json:"payout_target_currency_id,omitempty"`
	ThemeID                     int                     `json:"theme_id"`
	EnableMail                  bool                    `json:"enable_mail"`
	LoginType                   string                  `json:"login_type"`
	LoginTypeID                 int                     `json:"login_type_id"`
	BusinessModel               string                  `json:"business_model"`
	BusinessModelID             int                     `json:"business_model_id"`
	MemberPayToPartner          bool                    `json:"member_pay_to_partner"`
	AddressDetails              AddressDetails          `json:"address_details"`
	PaymentGatewayDetails       PaymentGatewayDetails   `json:"payment,omitempty"`
	SubscriptionPlanDetails     SubscriptionPlanDetails `json:"subscription_details,omitempty"`
	MemberGracePeriod           string                  `json:"member_grace_period"`
	MemberGracePeriodID         int                     `json:"member_grace_period_id,omitempty"`
	ExpiryWarningCount          int                     `json:"expiry_warning_count"`
	AlbumReviewEmail            string                  `json:"album_review_email"`
	SiteInfo                    string                  `json:"site_info"`
	DefaultPriceCode            int                     `json:"default_price_code"`
	DefaultPriceCodeCurrency    Currency                `json:"default_price_code_currency"`
	DefaultPriceCodeCurrencyID  int                     `json:"default_price_code_currency_id,omitempty"`
	MusicLanguage               string                  `json:"music_language"`
	MemberDefaultCountry        string                  `json:"member_default_country"`
	OutletsProcessingDuration   int                     `json:"outlets_processing_duration"`
	FreePlanLimit               int                     `json:"free_plan_limit"`
	MemberDefaultCountryID      string                  `json:"member_default_country_id,omitempty"`
	ProductReview               string                  `json:"product_review"`
	ProductReviewID             int                     `json:"product_review_id,omitempty"`
	DefaultCurrencyID           int                     `json:"default_currency_id,,omitempty"`
	DefaultPaymentGatewayId     string                  `json:"default_payment_gateway_id,omitempty"`
}

// structure to define payment gateway details
type PaymentGatewayDetails struct {
	PaymentGateWayID       string            `json:"payment_gateway_id,omitempty"`
	DefaultCurrency        string            `json:"default_currency,omitempty"`
	DefaultCurrencyDetails Currency          `json:"default_currency_details,omitempty"`
	PaymentGateways        []PaymentGateways `json:"payment_gateways,omitempty"`
	PayoutMinLimit         int               `json:"payout_min_limit,omitempty"`
	MaxRemittancePerMonth  int               `json:"payout_max_remittance_per_month,omitempty"`
	DefaultPaymentGateway  string            `json:"default_payment_gateway,omitempty"`
}

// structure to define configuration details for various payment gateways.
type PaymentGateways struct {
	Gateway               string `json:"gateway"`
	GatewayId             string `json:"gateway_id,omitempty"`
	Email                 string `json:"email"`
	ClientId              string `json:"client_id"`
	ClientSecret          string `json:"client_secret"`
	Payin                 bool   `json:"payin"`
	Payout                bool   `json:"payout"`
	DefaultPayinCurrency  string `json:"default_payin_currency"`
	DefaultPayoutCurrency string `json:"default_payout_currency"`
}

// structure to get payment gateway details
type GetPaymentGateways struct {
	PaymentGatewayDetails PaymentGatewayDetails `json:"payment"`
}

// structure to define address details of a partner
type AddressDetails struct {
	Address    string `json:"address,omitempty"`
	Street     string `json:"street,omitempty"`
	Country    string `json:"country,omitempty"`
	State      string `json:"state,omitempty"`
	City       string `json:"city,omitempty"`
	PostalCode string `json:"postal_code,omitempty"`
}

// structure to define subscription plan details of a partner
type SubscriptionPlanDetails struct {
	ID         string `json:"plan_id,omitempty"`
	StartDate  string `json:"plan_start_date,omitempty"`
	LaunchDate string `json:"plan_launch_date,omitempty"`
}

// structure define currency details
type Currency struct {
	ID   interface{} `json:"id,omitempty"`
	Name string      `json:"name,omitempty"`
}

// structure to define payout curreny details
type PayoutCurrencyDetails struct {
	ID   interface{} `json:"id"`
	Name string      `json:"name"`
}

// structure to define contact details of a partner
type ContactDetails struct {
	ContactPerson string `json:"contact_person,omitempty"`
	Email         string `json:"email,omitempty"`
	NoReplyEmail  string `json:"noreply_email,omitempty"`
	FeedbackEmail string `json:"feedback_email,omitempty"`
	SupportEmail  string `json:"support_email,omitempty"`
}

// structure to define the response format for any get all endpoints
type ResponseData struct {
	Metadata any `json:"metadata"`
	Data     any `json:"records"`
}

// structure to define request format for creating partner stores
type PartnerStores struct {
	Stores []string `json:"stores"`
}

// structure to define response format for retrieving store details
type GetPartnerStores struct {
	Stores []Store `json:"stores"`
}

// structure to define store details
type Store struct {
	Id      string `json:"id"`
	StoreId string `json:"store_id"`
}

// the response structure of store data from store service
type StoreResponse struct {
	StoreData []StoreData `json:"data"`
}

// the store deatils retrieved from store service
type StoreData struct {
	Id   string `json:"store_id"`
	Name string `json:"name"`
}

// structure to define terms and conditions for partner
type TermsAndConditions struct {
	Name        string
	Description string
	Language    string
}

// structure to define ouath credentials of a partner
type PartnerOauthCredential struct {
	ClientId     string
	ClientSecret string
	ProviderId   string
	PartnerId    string
}

// structure to define query params for retrieving all partners
type Params struct {
	Limit   int32  `form:"limit" default:"10"`
	Page    int32  `form:"page" default:"1"`
	Status  string `form:"status"`
	Order   string `form:"order"`
	Sort    string `form:"sort"`
	Name    string `form:"name"`
	Country string `form:"country"`
}

// structure to define error map result of get partner by id
type GetPartnerByIdMapResult struct {
	Value interface{}
	Err   error
}

// structure to retrieve all partners
type ListAllPartners struct {
	UUID                 string         `json:"uuid"`
	Name                 string         `json:"name"`
	URL                  string         `json:"url"`
	Logo                 string         `json:"logo"`
	Language             string         `json:"language"`
	BrowserTitle         string         `json:"browser_title"`
	Active               bool           `json:"active"`
	PaymentURL           string         `json:"payment_url"`
	ProfileURL           string         `json:"profile_url"`
	LandingPage          string         `json:"landing_page"`
	MobileVerifyInterval int            `json:"mobile_verify_interval"`
	PayoutTargetCurrency string         `json:"payout_target_currency"`
	EnableMail           bool           `json:"enable_mail"`
	BusinessModel        string         `json:"business_model"`
	MemberPayToPartner   bool           `json:"member_pay_to_partner"`
	ContactDetails       ContactDetails `json:"contact_details"`
	Theme                string         `json:"theme"`
	LoginType            string         `json:"login_type"`
	Users                int            `json:"users"`
}

// strcture to get a single partner
type GetPartner struct {
	Name                        string                  `json:"name"`
	URL                         string                  `json:"url"`
	Logo                        string                  `json:"logo"`
	Favicon                     string                  `json:"favicon"`
	Loader                      string                  `json:"loader"`
	LoginPageLogo               string                  `json:"login_page_logo"`
	ContactDetails              ContactDetails          `json:"contact_details"`
	BackgroundColor             string                  `json:"background_color"`
	BackgroundImage             string                  `json:"background_image"`
	Active                      bool                    `json:"active"`
	WebsiteURL                  string                  `json:"website_url"`
	Language                    string                  `json:"language,omitempty"`
	BrowserTitle                string                  `json:"browser_title"`
	ProfileURL                  string                  `json:"profile_url"`
	PaymentURL                  string                  `json:"payment_url"`
	LandingPage                 string                  `json:"landing_page"`
	MobileVerifyInterval        int                     `json:"mobile_verify_interval"`
	TermsAndConditionsVersionID int                     `json:"terms_and_conditions_version_id"`
	PayoutTargetCurrency        string                  `json:"payout_target_currency"`
	Theme                       string                  `json:"theme"`
	EnableMail                  bool                    `json:"enable_mail"`
	LoginType                   string                  `json:"login_type"`
	BusinessModel               string                  `json:"business_model"`
	MemberPayToPartner          bool                    `json:"member_pay_to_partner"`
	AddressDetails              AddressDetails          `json:"address_details"`
	SubscriptionPlanDetails     SubscriptionPlanDetails `json:"subscription_details,omitempty"`
	MemberGracePeriod           string                  `json:"member_grace_period"`
	ExpiryWarningCount          int                     `json:"expiry_warning_count"`
	AlbumReviewEmail            string                  `json:"album_review_email"`
	SiteInfo                    string                  `json:"site_info"`
	DefaultPriceCodeCurrency    Currency                `json:"default_price_code_currency"`
	MusicLanguage               string                  `json:"music_language"`
	MemberDefaultCountry        string                  `json:"member_default_country"`
	OutletsProcessingDuration   int                     `json:"outlets_processing_duration"`
	FreePlanLimit               int                     `json:"free_plan_limit"`
	ProductReview               string                  `json:"product_review"`
	DefaultCurrency             string                  `json:"default_currency"`
	DefaultPaymentGateway       string                  `json:"default_payment_gateway"`
}

// structure to define the request format to update partner status
type UpdatePartnerStatus struct {
	Active bool `json:"active"`
}

// structure to define params to get all product type or track file quality of a particular partner

type QueryParams struct {
	Limit  int32  `form:"limit" default:"10"`
	Page   int32  `form:"page" default:"1"`
	Order  string `form:"order"`
	Sort   string `form:"sort"`
	Fields string `form:"fields"`
}

// structure to define partner product types
type GetPartnerProdTypesAndTrackQuality struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

// response structure of get all product types
type ProductTypes struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Struct for the  DefaultApiResponse
type DefaultApiResponse struct {
	Status  string                 `json:"status"`
	Message string                 `json:"message"`
	Code    int                    `json:"code"`
	Data    Data                   `json:"data"`
	Errors  map[string]interface{} `json:"errors"`
}

// Struct for the response data
type Data struct {
	Records []ProductTypes `json:"records"`
}
