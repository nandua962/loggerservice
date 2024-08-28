package entities

import (
	"encoding/json"
	"partner/internal/consts"

	"gitlab.com/tuneverse/toolkit/utils"
)

type UpdateTermsAndConditions map[string]interface{}

func (x UpdateTermsAndConditions) GetData(key string) string {
	if name, ok := x[key].(string); ok {
		return name
	}
	return ""
}

type PartnerProperties map[string]interface{}

// function to get name from a map of type PartnerProperties
func (x PartnerProperties) GetPartnerName() string {
	if name, ok := x["name"].(string); ok {
		return name
	}
	return ""
}

// function to get background color  from a map of type PartnerProperties
func (x PartnerProperties) GetBGColor() string {
	if bgColor, ok := x["background_color"].(string); ok {
		if utils.IsEmpty(bgColor) {
			bgColor = consts.BackgroundColorDefaultVal
			return bgColor
		}
		return bgColor
	}
	return ""
}

// This function retrieves the background image from a given map of type PartnerProperties.
// If the background image is empty, the function will set it to a default value.
func (x PartnerProperties) GetBGImage() string {
	if bgImage, ok := x["background_image"].(string); ok {
		if utils.IsEmpty(bgImage) {
			bgImage = consts.BackgroundImageDefaultVal
			return bgImage
		}
		return bgImage
	}
	return ""
}

// function to get website url from a map of type PartnerProperties
func (x PartnerProperties) GetWebsiteURL() string {
	if websiteURL, ok := x["website_url"].(string); ok {
		return websiteURL
	}
	return ""
}

//		Function to get browser title from a given map of type PartnerProperties.
//	 If the browser title  is empty, the function will set it to a default value.
func (x PartnerProperties) GetBrowserTitle() string {
	if browserTitle, ok := x["browser_title"].(string); ok {
		if utils.IsEmpty(browserTitle) {
			browserTitle = consts.BrowserTitleDefaultVal
			return browserTitle
		}
		return browserTitle
	}
	return ""
}

// function to get profile url from a map of type PartnerProperties
func (x PartnerProperties) GetProfileURL() string {
	if profileURL, ok := x["profile_url"].(string); ok {
		return profileURL
	}
	return ""
}

// function to get album review email from a map of type PartnerProperties
func (x PartnerProperties) GetAlbumEmail() string {
	if albumEmail, ok := x["album_review_email"].(string); ok {
		return albumEmail
	}
	return ""
}

// function to get  url from a map of type PartnerProperties
func (x PartnerProperties) GetURL() string {
	if url, ok := x["url"].(string); ok {
		return url
	}
	return ""
}

// function to get payment url from a map of type PartnerProperties
func (x PartnerProperties) GetPaymentURL() string {
	if paymentURL, ok := x["payment_url"].(string); ok {
		return paymentURL
	}
	return ""
}

// function to get landing page from a map of type PartnerProperties
func (x PartnerProperties) GetLandingPage() string {
	if landingPage, ok := x["landing_page"].(string); ok {
		return landingPage
	}
	return ""
}

// function to get free plan limit from a map of type PartnerProperties
// If the free plan limit  is empty, the function will set it to a default value.
func (x PartnerProperties) GetFreePlanLimit() float64 {
	if freePlanLimit, ok := x["free_plan_limit"].(float64); ok {
		if freePlanLimit == 0 {
			freePlanLimit = consts.FreePlanLimitDefaultVal
			return freePlanLimit
		}
		return freePlanLimit
	}

	return 0
}

// function to get expiry warning count from a map of type PartnerProperties
// If the expiry warning count  is empty, the function will set it to a default value.
func (x PartnerProperties) GetExpiryWarningCount() float64 {
	if expiryWarningCount, ok := x["expiry_warning_count"].(float64); ok {
		if expiryWarningCount == 0 {
			expiryWarningCount = consts.ExpiryWarningCountDefaultVal
			return expiryWarningCount
		}
		return expiryWarningCount
	}
	return 0
}

// function to get outlets processing duration from a map of type PartnerProperties
// If the outlets processing duration  is empty, the function will set it to a default value.
func (x PartnerProperties) GetOutletsProcessingDuration() float64 {
	if outletsProcessingDuration, ok := x["outlets_processing_duration"].(float64); ok {
		if outletsProcessingDuration == 0 {
			outletsProcessingDuration = consts.OutletsProcessingDurationDefaultVal
			return outletsProcessingDuration
		}
		return outletsProcessingDuration
	}
	return 0
}

// function to get theme from a map of type PartnerProperties
// If the theme  is empty, the function will set it to a default value.
func (x PartnerProperties) GetThemeID() float64 {
	if themeID, ok := x["theme_id"].(float64); ok {
		if themeID == 0 {
			themeID = consts.ThemeIdDefaultVal
			return themeID
		}
		return themeID
	}
	return 0
}

// function to get business model from a map of type PartnerProperties
// If the business model  is empty, the function will set it to a default value.
func (x PartnerProperties) GetBusinessModel() string {
	if businessModel, ok := x["business_model"].(string); ok {
		return businessModel
	}
	return ""
}

// function to get product review from a map of type PartnerProperties
// If the product review  is empty, the function will set it to a default value.
func (x PartnerProperties) GetProductReview() string {
	if productReview, ok := x["product_review"].(string); ok {
		if utils.IsEmpty(productReview) {
			productReview = consts.ProductReviewDefaultVal
			return productReview
		}
		return productReview
	}
	return ""
}

// function to get login type from a map of type PartnerProperties
// If the login type  is empty, the function will set it to a default value.
func (x PartnerProperties) GetLoginType() string {
	if loginType, ok := x["login_type"].(string); ok {
		if utils.IsEmpty(loginType) {
			loginType = consts.LoginTypeIdDefaultVal
			return loginType
		}

		return loginType
	}
	return ""
}

// function to get member grace period from a map of type PartnerProperties
// If the member grace  period  is empty, the function will set it to a default value.
func (x PartnerProperties) GetMemberGracePeriod() string {
	if memberGracePeriod, ok := x["member_grace_period"].(string); ok {
		if utils.IsEmpty(memberGracePeriod) {
			memberGracePeriod = consts.MemberGracePeriodDefaultVal
			return memberGracePeriod
		}
		return memberGracePeriod
	}
	return ""
}

// function to get payout currency from a map of type PartnerProperties
// If the get payout currency is empty, the function will set it to a default value.
func (x PartnerProperties) GetPayoutTargetCurren() string {
	if payoutTargetCurren, ok := x["payout_target_currency"].(string); ok {
		if utils.IsEmpty(payoutTargetCurren) {
			payoutTargetCurren = consts.PayoutTargetCurrencyDefaultVal
			return payoutTargetCurren
		}
		return payoutTargetCurren
	}
	return ""
}

// function to get member default country from a map of type PartnerProperties
// If the get member default country is empty, the function will set it to a default value.
func (x PartnerProperties) GetMemberDefaultCountry() string {
	if memberDefCountry, ok := x["member_default_country"].(string); ok {
		if utils.IsEmpty(memberDefCountry) {
			memberDefCountry = consts.MemberDefaultCountryDefaultVal
			return memberDefCountry
		}
		return memberDefCountry
	}
	return ""
}

// function to get music language from a map of type PartnerProperties
// If the get music language is empty, the function will set it to a default value.
func (x PartnerProperties) GetMusicLang() string {
	if musicLang, ok := x["music_language"].(string); ok {
		if utils.IsEmpty(musicLang) {
			musicLang = consts.MusicLanguageDefaultVal
			return musicLang
		}
		return musicLang
	}
	return ""
}

// function to get languagefrom a map of type PartnerProperties
// If the get language is empty, the function will set it to a default value.
func (x PartnerProperties) GetLang() string {
	if lang, ok := x["language"].(string); ok {
		if utils.IsEmpty(lang) {
			lang = consts.LanguageDefaultVal
			return lang
		}
		return lang
	}
	return ""
}

// function to get site info from a map of type PartnerProperties
func (x PartnerProperties) GetSiteInfo() string {
	if site, ok := x["site_info"].(string); ok {
		return site
	}
	return ""
}

// function to get logo from a map of type PartnerProperties
func (x PartnerProperties) GetLogo() string {
	if site, ok := x["logo"].(string); ok {
		return site
	}
	return ""
}

// function to get  member pay to partner from a map of type PartnerProperties
func (x PartnerProperties) GetMemberPayToPartner() bool {
	if memberPay, ok := x["member_pay_to_partner"].(bool); ok {
		return memberPay
	}
	return false
}

// function to get mobile verify interval from a map of type PartnerProperties
// If the get mobile verify interval is empty, the function will set it to a default value.
func (x PartnerProperties) GetMobileVerifyInterval() float64 {
	if mobileVerify, ok := x["mobile_verify_interval"].(float64); ok {
		if mobileVerify == 0 {
			mobileVerify = consts.MobileVerifyIntervalDefaultVal
			return mobileVerify
		}
		return mobileVerify
	}
	return 0
}

// function to get enable mail field from a map of type PartnerProperties
func (x PartnerProperties) GetEnableMail() bool {
	if enableMail, ok := x["enable_mail"].(bool); ok {
		return enableMail
	}
	return false
}

// function to get payment details from a map of type PartnerProperties
func (x PartnerProperties) GetPaymentDetails() (map[string]interface{}, error) {
	var data map[string]interface{}
	if paymentDetails, ok := x["payment"]; ok {
		d, err := json.Marshal(paymentDetails)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(d, &data)
		if err != nil {
			return nil, err
		}
	}
	return data, nil
}

// function to get payout min limit from a map of type PartnerProperties
// If the get payout min limit is empty, the function will set it to a default value.
func (x PartnerProperties) GetPayoutMinLmit() (float64, error) {
	data, err := x.GetPaymentDetails()
	if err != nil {
		return 0, err
	}

	if payoutMinLimit, ok := data["payout_min_limit"]; ok {
		if payoutMinLimit == 0 {
			payoutMinLimit = consts.PayoutMinLimitDefaultVal
			return payoutMinLimit.(float64), nil
		}
		return payoutMinLimit.(float64), nil

	}

	return 0, nil
}

// function to get default currency from a map of type PartnerProperties
func (x PartnerProperties) GetDefaultCurrency() (string, error) {
	data, err := x.GetPaymentDetails()
	if err != nil {
		return "", err
	}

	if defaultCurrency, ok := data["default_currency"]; ok {
		return defaultCurrency.(string), nil

	}

	return "", nil
}

// function to get default gateway from a map of type PartnerProperties
func (x PartnerProperties) GetDefaultGateway() (string, error) {
	data, err := x.GetPaymentDetails()
	if err != nil {
		return "", err
	}

	if defaultGateway, ok := data["default_payment_gateway"]; ok {
		return defaultGateway.(string), nil

	}

	return "", nil
}

// function to get max remittance per month from a map of type PartnerProperties
// If the get max remittance per month is empty, the function will set it to a default value.
func (x PartnerProperties) GetMaxRemittance() (float64, error) {
	data, err := x.GetPaymentDetails()
	if err != nil {
		return 0, err
	}

	if maxRemittance, ok := data["payout_max_remittance_per_month"]; ok {
		if maxRemittance == 0 {
			maxRemittance = consts.MaxRemittancePerMonthDefaultVal
			return maxRemittance.(float64), nil
		}
		return maxRemittance.(float64), nil

	}

	return 0, nil
}

// function to get payment gateways from a map of type PartnerProperties
func (x PartnerProperties) GetPaymentGateways() ([]PaymentGateways, error) {
	var data []PaymentGateways

	paymentDetails, err := x.GetPaymentDetails()
	if err != nil {
		return []PaymentGateways{}, err
	}
	if paymentGateways, ok := paymentDetails["payment_gateways"]; ok {
		d, err := json.Marshal(paymentGateways)
		if err != nil {
			return []PaymentGateways{}, err
		}

		err = json.Unmarshal(d, &data)
		if err != nil {
			return []PaymentGateways{}, err
		}
	}
	return data, nil
}

// function to get contact details from a map of type PartnerProperties
func (x PartnerProperties) GetContactDetails() (map[string]interface{}, error) {

	var data map[string]interface{}
	if contactDetails, ok := x["contact_details"]; ok {
		d, err := json.Marshal(contactDetails)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(d, &data)
		if err != nil {
			return nil, err
		}
	}
	return data, nil
}

// function to get email from a map of type PartnerProperties
func (x PartnerProperties) GetEmail() (string, error) {
	data, err := x.GetContactDetails()
	if err != nil {
		return "", err
	}

	if email, ok := data["email"]; ok {
		return email.(string), nil

	}

	return "", nil
}

// function to get support email from a map of type PartnerProperties
func (x PartnerProperties) GetSupportEmail() (string, error) {
	data, err := x.GetContactDetails()
	if err != nil {
		return "", err
	}

	if email, ok := data["support_email"]; ok {
		return email.(string), nil

	}

	return "", nil
}

// function to get feedback email from a map of type PartnerProperties
func (x PartnerProperties) GetFeedBackEmail() (string, error) {
	data, err := x.GetContactDetails()
	if err != nil {
		return "", err
	}

	if email, ok := data["feedback_email"]; ok {
		return email.(string), nil

	}

	return "", nil
}

// function to get no reply email from a map of type PartnerProperties
func (x PartnerProperties) GetNoReplyEmail() (string, error) {
	data, err := x.GetContactDetails()
	if err != nil {
		return "", err
	}
	if email, ok := data["noreply_email"]; ok {
		return email.(string), nil

	}

	return "", nil
}

// function to get contact person from a map of type PartnerProperties
func (x PartnerProperties) GetContactPerson() (string, error) {
	data, err := x.GetContactDetails()
	if err != nil {
		return "", err
	}

	if contactPerson, ok := data["contact_person"]; ok {
		return contactPerson.(string), nil

	}

	return "", nil
}

// function to get address details from a map of type PartnerProperties
func (x PartnerProperties) GetAddressDetails() (map[string]interface{}, error) {

	var data map[string]interface{}
	if addressDetails, ok := x["address_details"]; ok {
		d, err := json.Marshal(addressDetails)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(d, &data)
		if err != nil {
			return nil, err
		}
	}
	return data, nil
}

// function to get address from address details within a map of type PartnerProperties
func (x PartnerProperties) GetAddress() (string, error) {
	data, err := x.GetAddressDetails()
	if err != nil {
		return "", err
	}

	if address, ok := data["address"]; ok {
		return address.(string), nil

	}

	return "", nil
}

// function to get street from address details within a map of type PartnerProperties
func (x PartnerProperties) GetStreet() (string, error) {
	data, err := x.GetAddressDetails()
	if err != nil {
		return "", err
	}

	if street, ok := data["street"]; ok {
		return street.(string), nil

	}

	return "", nil
}

// function to get country from address details within a map of type PartnerProperties
func (x PartnerProperties) GetCountry() (string, error) {
	data, err := x.GetAddressDetails()
	if err != nil {
		return "", err
	}

	if country, ok := data["country"]; ok {
		return country.(string), nil

	}

	return "", nil
}

// function to get state from address details within a map of type PartnerProperties
func (x PartnerProperties) GetState() (string, error) {
	data, err := x.GetAddressDetails()
	if err != nil {
		return "", err
	}

	if state, ok := data["state"]; ok {
		return state.(string), nil

	}

	return "", nil
}

// function to get city from address details within a map of type PartnerProperties
func (x PartnerProperties) GetCity() (string, error) {
	data, err := x.GetAddressDetails()
	if err != nil {
		return "", err
	}

	if city, ok := data["city"]; ok {
		return city.(string), nil

	}

	return "", nil
}

// function to get postal code from address details  within a map of type PartnerProperties
func (x PartnerProperties) GetPostalCode() (string, error) {
	data, err := x.GetAddressDetails()
	if err != nil {
		return "", err
	}

	if postalCode, ok := data["postal_code"]; ok {
		return postalCode.(string), nil

	}

	return "", nil
}

// function to get default price code currency from a map of type PartnerProperties
// If the get default price code currency is empty, the function will set it to a default value.
func (x PartnerProperties) GetDefaultPriceCodeCurrency() (Currency, error) {

	var priceCodeCurren Currency
	if Currency, ok := x["default_price_code_currency"]; ok {
		d, err := json.Marshal(Currency)
		if err != nil {
			return priceCodeCurren, err
		}

		err = json.Unmarshal(d, &priceCodeCurren)
		if err != nil {
			return priceCodeCurren, err
		}
		if utils.IsEmpty(priceCodeCurren.Name) {
			priceCodeCurren.Name = consts.DefaultPriceCodeCurrencyDefaultVal
			return priceCodeCurren, nil
		}
	}
	return priceCodeCurren, nil
}

// function to get subscription plan details from a map of type PartnerProperties
func (x PartnerProperties) GetSubscriptionPlanDetails() (map[string]interface{}, error) {

	var data map[string]interface{}
	if SubscriptionPlanDetails, ok := x["subscription_details"]; ok {
		d, err := json.Marshal(SubscriptionPlanDetails)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(d, &data)
		if err != nil {
			return nil, err
		}
	}
	return data, nil
}

// function to get subscription Id from subscription plan details within a map of type PartnerProperties
func (x PartnerProperties) GetSubscriptionId() (string, error) {
	data, err := x.GetSubscriptionPlanDetails()
	if err != nil {
		return "", err
	}

	if planId, ok := data["plan_id"]; ok {
		return planId.(string), nil

	}

	return "", nil
}

// function to get start date from subscription plan details within a map of type PartnerProperties
func (x PartnerProperties) GetStartDate() (string, error) {
	data, err := x.GetSubscriptionPlanDetails()
	if err != nil {
		return "", err
	}

	if planStartDate, ok := data["plan_start_date"]; ok {
		return planStartDate.(string), nil

	}

	return "", nil
}

// function to get launch date from subscription plan details within a map of type PartnerProperties
func (x PartnerProperties) GetLaunchDate() (string, error) {
	data, err := x.GetSubscriptionPlanDetails()
	if err != nil {
		return "", err
	}

	if planLaunchDate, ok := data["plan_launch_date"]; ok {
		return planLaunchDate.(string), nil

	}

	return "", nil
}
