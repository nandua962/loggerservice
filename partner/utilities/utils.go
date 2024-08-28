package utilities

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"partner/internal/consts"
	"partner/internal/entities"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/tidwall/gjson"
	cacheConf "gitlab.com/tuneverse/toolkit/core/cache"
	"gitlab.com/tuneverse/toolkit/core/logger"
	"gitlab.com/tuneverse/toolkit/models"
	"gitlab.com/tuneverse/toolkit/models/api"
	"gitlab.com/tuneverse/toolkit/utils"
)

// Function to check if an interface{} value is empty
func IsValueEmpty(val interface{}) bool {
	switch v := val.(type) {
	case int:
		return v == 0
	case string:
		return v == ""
	case []int:
		return len(v) == 0
	default:
		return false
	}
}

// function to generate generate client id and client secret
func GenerateClientCredentials() (string, error) {
	randomBytes := make([]byte, consts.ByteSize)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	randomString := base64.StdEncoding.EncodeToString(randomBytes)

	var cleanString string
	for _, char := range randomString {
		if (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9') {
			cleanString = fmt.Sprintf("%s%s", cleanString, string(char))
		}
	}

	timestamp := time.Now().UnixNano()
	cleanString = fmt.Sprintf("%s%d", cleanString, timestamp)
	result, err := HashString(cleanString)
	if err != nil {
		return "", err
	}
	return result, nil
}

// function to hash the given string
func HashString(input string) (string, error) {
	hasher := sha256.New()
	_, err := hasher.Write([]byte(input))
	if err != nil {
		return "", err
	}
	hashed := hasher.Sum(nil)
	result := hex.EncodeToString(hashed)
	return result, err
}

// SetDefault function to set default values if empty
func SetDefault(partnerData *entities.Partner) {

	if utils.IsEmpty(partnerData.Language) {
		partnerData.Language = consts.LanguageDefaultVal
	}

	if utils.IsEmpty(partnerData.MusicLanguage) {
		partnerData.MusicLanguage = consts.MusicLanguageDefaultVal
	}

	if utils.IsEmpty(partnerData.DefaultPriceCodeCurrency.Name) {
		partnerData.DefaultPriceCodeCurrency.Name = consts.DefaultPriceCodeCurrencyDefaultVal
	}

	if partnerData.FreePlanLimit == 0 {
		partnerData.FreePlanLimit = consts.FreePlanLimitDefaultVal
	}

	if partnerData.LoginType == "" {
		partnerData.LoginType = consts.LoginTypeIdDefaultVal
	}

	if utils.IsEmpty(partnerData.MemberDefaultCountry) {
		partnerData.MemberDefaultCountry = consts.MemberDefaultCountryDefaultVal
	}

	if utils.IsEmpty(partnerData.BrowserTitle) {
		partnerData.BrowserTitle = consts.BrowserTitleDefaultVal
	}

	if utils.IsEmpty(partnerData.LoginPageLogo) {
		partnerData.LoginPageLogo = consts.LoginPageLogoDefaultVal
	}

	if utils.IsEmpty(partnerData.Logo) {
		partnerData.Logo = consts.LogoDefaultVal
	}

	if utils.IsEmpty(partnerData.Loader) {
		partnerData.Loader = consts.LoaderDefaultVal
	}

	if utils.IsEmpty(partnerData.BackgroundImage) {
		partnerData.BackgroundImage = consts.BackgroundImageDefaultVal
	}

	if partnerData.ThemeID == 0 {
		partnerData.ThemeID = consts.ThemeIdDefaultVal
	}

	if partnerData.OutletsProcessingDuration == 0 {
		partnerData.OutletsProcessingDuration = consts.OutletsProcessingDurationDefaultVal
	}

	if partnerData.PaymentGatewayDetails.MaxRemittancePerMonth == 0 {
		partnerData.PaymentGatewayDetails.MaxRemittancePerMonth = consts.MaxRemittancePerMonthDefaultVal
	}

	if partnerData.ExpiryWarningCount == 0 {
		partnerData.ExpiryWarningCount = consts.ExpiryWarningCountDefaultVal
	}

	if partnerData.PaymentGatewayDetails.PayoutMinLimit == 0 {
		partnerData.PaymentGatewayDetails.PayoutMinLimit = consts.PayoutMinLimitDefaultVal
	}

	if utils.IsEmpty(partnerData.Favicon) {
		partnerData.Favicon = consts.FaviconDefaultVal
	}

	if utils.IsEmpty(partnerData.PayoutTargetCurrency) {
		partnerData.PayoutTargetCurrency = consts.PayoutTargetCurrencyDefaultVal
	}

	if utils.IsEmpty(partnerData.PaymentGatewayDetails.DefaultCurrency) {
		partnerData.PaymentGatewayDetails.DefaultCurrency = consts.PayoutCurrencyDefaultVal
	}
}

// function to check its a valid hex code
func IsValidHexColorCode(colorCode string) bool {
	if len(colorCode) < 1 || colorCode[0] != '#' {
		return false
	}
	colorCode = colorCode[1:]
	_, err := strconv.ParseUint(colorCode, consts.HexadecimalBase, consts.Uint32BitSize)
	return err == nil && (len(colorCode) == consts.ValidColorCodeLength3 || len(colorCode) == consts.ValidColorCodeLength6)
}

// IsValidPartnerName function used to validate name
func IsValidPartnerName(name string) bool {
	// Define a regular expression pattern that allows A-Z, 0-9, '-', and spaces
	validPattern := "^[a-zA-Z0-9\\- '.,]+$"
	regex := regexp.MustCompile(validPattern)
	return regex.MatchString(name)
}

// to check the postal code is valid or not
func IsValidPostalCode(code string) bool {
	validPattern := `^[a-zA-Z0-9\s\-]{4,10}$`
	regex := regexp.MustCompile(validPattern)
	return regex.MatchString(code)
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
	result = api.Response{Status: consts.FailureKey, Message: message, Code: code, Data: map[string]string{}, Errors: errors}
	return result
}

// function to validate an url
func IsValidURL(str string) bool {
	var rxURL = regexp.MustCompile(consts.URLExp)

	if str == "" || utf8.RuneCountInString(str) >= consts.MaxURLRuneCount || len(str) <= consts.MinURLRuneCount || strings.HasPrefix(str, ".") {
		return false
	}
	strTemp := str
	if strings.Contains(str, ":") && !strings.Contains(str, "://") {
		strTemp = "http://" + str
	}
	u, err := url.Parse(strTemp)
	if err != nil {
		return false
	}
	if strings.HasPrefix(u.Host, ".") {
		return false
	}
	if u.Host == "" && (u.Path != "" && !strings.Contains(u.Path, ".")) {
		return false
	}

	return rxURL.MatchString(str)
}

// function to check the existence of a currency
func IsCurrencyIsoExists(ctx context.Context, cache cacheConf.Cache, currencyIso string, apiUrl string) (bool, error) {
	var cachekey = fmt.Sprintf("%s%s", consts.CurrencyExistsCacheKey, currencyIso)

	log := logger.Log().WithContext(ctx)
	headers := make(map[string]interface{})
	// check whether the data is present in cache using cache key
	exists, err := IsCached(cache, cachekey)
	if err != nil {
		return false, err
	}
	if exists != nil {
		return *exists, nil
	}
	// retieve data from url if data is not available in the cache
	apiURL := fmt.Sprintf("%s/currencies/exists/%s", apiUrl, currencyIso)
	fmt.Println("urlllllllllllllllllllllllll", apiURL)
	headers["Content-Type"] = "application/json"
	response, err := utils.APIRequest(http.MethodHead, apiURL, headers, nil)
	if err != nil {
		log.Printf("currency service failed  :failed to make API request: %v", err)
		return false, consts.ErrUtilityServiceConnectionLost
	}
	defer response.Body.Close()
	//set data to cache  using cache key
	_, err = cache.Set(cachekey, response.StatusCode == http.StatusOK, consts.CacheExpiryTime)
	if err != nil {
		log.Printf("unable to connect to cache service: %s\n", err)
	}
	return response.StatusCode == http.StatusOK, nil

}

// function to get currency name by its id
func GetCurrencyName(ctx context.Context, cache cacheConf.Cache, id int, apiUrl string) (string, error) {
	var cachekey = fmt.Sprintf("%s%d", consts.CurrencyIdCacheKey, id)

	log := logger.Log().WithContext(ctx)
	headers := make(map[string]interface{})
	// check whether the data is present in cache using cache key
	ttl, err := cache.TTL(cachekey)
	if err != nil {
		cacheErr := fmt.Errorf("unable to get TTL: %s", err)
		return "", cacheErr
	}
	if ttl > 0 {
		databytes, cacheErr := cache.Get(cachekey).Bytes()

		if cacheErr != nil {
			log.Printf("%v", err)
			return "", err
		}
		currency := gjson.GetBytes(databytes, "data.iso").String()
		return currency, nil
	}

	// retieve data from url if data is not available in the cache
	apiURL := fmt.Sprintf("%s/currencies/%d", apiUrl, id)
	headers["Content-Type"] = "application/json"
	response, err := utils.APIRequest(http.MethodGet, apiURL, headers, nil)
	fmt.Println("urlllllllllllllllllllllllll", apiURL)
	if err != nil {
		log.Printf("failed to make API request: %v", err)
		return "", consts.ErrUtilityServiceConnectionLost
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		log.Printf("currency service failed :no data")
		return "", nil
	}
	// Read response body
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Printf("currency service failed :failed to read response body: %v", err)
		return "", err
	}
	currency := gjson.GetBytes(body, "data.iso").String()
	//set data to cache  using cache key
	if response.StatusCode == http.StatusOK {
		_, err = cache.Set(cachekey, body, consts.CacheExpiryTime)
		if err != nil {
			log.Printf("unable to connect to cache service: %s\n", err)
		}
	}
	return currency, nil

}

// function to get currency id based on its name
func GetCurrencyId(ctx context.Context, cache cacheConf.Cache, iso string, apiUrl string) (int, error) {
	var cachekey = fmt.Sprintf("%s%s", consts.CurrencyNameCacheKey, iso)

	log := logger.Log().WithContext(ctx)
	headers := make(map[string]interface{})
	// check whether the data is present in cache using cache key
	ttl, err := cache.TTL(cachekey)
	if err != nil {
		cacheErr := fmt.Errorf("unable to get TTL: %s", err)
		return 0, cacheErr
	}
	if ttl > 0 {
		databytes, cacheErr := cache.Get(cachekey).Bytes()

		if cacheErr != nil {
			log.Printf("%v", err)
			return 0, err
		}
		currencyId := gjson.GetBytes(databytes, "data.id").Int()
		return int(currencyId), nil
	}

	// retieve data from url if data is not available in the cache
	apiURL := fmt.Sprintf("%s/currencies/exists/%s", apiUrl, iso)
	fmt.Println("urlllllllllllllllllllllllll", apiURL)
	headers["Content-Type"] = "application/json"

	response, err := utils.APIRequest(http.MethodGet, apiURL, headers, nil)
	if err != nil {
		log.Printf("failed to connect currency service: %v", err)
		return 0, consts.ErrUtilityServiceConnectionLost
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		log.Printf("currency service failed :no data")
		return 0, nil
	}
	// Read response body
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Printf("currency service failed :failed to read response body: %v", err)
		return 0, err
	}
	currencyId := gjson.GetBytes(body, "data.id").Int()
	//set data to cache  using cache key
	if response.StatusCode == http.StatusOK {
		_, err = cache.Set(cachekey, body, consts.CacheExpiryTime)
		if err != nil {
			log.Printf("unable to connect to cache service: %s\n", err)
		}
	}
	return int(currencyId), nil

}

// function to check the existence of a country
func IsCountryExists(ctx context.Context, cache cacheConf.Cache, country string, apiUrl string) (bool, error) {
	var cachekey = fmt.Sprintf("%s%s", consts.CountryExistsCacheKey, country)

	log := logger.Log().WithContext(ctx)
	// check whether the data is present in cache using cache key
	exists, err := IsCached(cache, cachekey)
	if err != nil {
		return false, err
	}
	if exists != nil {
		return *exists, nil
	}
	// Prepare API URL with query parameters
	// retieve data from url if data is not available in the cache
	apiURL := fmt.Sprintf("%s/countries/exists?iso=%s", apiUrl, country)
	fmt.Println("urlllllllllllllllllllllllll", apiURL)
	// Make the API request
	headers := make(map[string]interface{})
	headers["Content-Type"] = "application/json"

	response, err := utils.APIRequest(http.MethodGet, apiURL, headers, nil)
	if err != nil {
		log.Printf("failed to connect country service: %v", err)
		return false, consts.ErrUtilityServiceConnectionLost
	}
	defer response.Body.Close()

	// Read response body
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Printf("country service failed :failed to read response body: %v", err)
		return false, err
	}
	countryExists := gjson.GetBytes(body, "data.exists").Bool()

	_, err = cache.Set(cachekey, response.StatusCode == http.StatusOK && countryExists, consts.CacheExpiryTime)
	//set data to cache  using cache key
	if err != nil {
		log.Printf("unable to connect to cache service: %s\n", err)
	}
	return response.StatusCode == http.StatusOK && countryExists, nil
}

// function to check the existence of a language
func IsLanguageIsoExists(ctx context.Context, cache cacheConf.Cache, languageIso string, apiUrl string) (bool, error) {
	var cachekey = fmt.Sprintf("%s%s", consts.LanguageExistsCacheKey, languageIso)

	log := logger.Log().WithContext(ctx)
	headers := make(map[string]interface{})
	// check whether the data is present in cache using cache key
	exists, err := IsCached(cache, cachekey)
	if err != nil {
		return false, err
	}
	if exists != nil {
		return *exists, nil
	}
	// retieve data from url if data is not available in the cache
	apiURL := fmt.Sprintf("%s/languages/exists/%s", apiUrl, languageIso)
	fmt.Println("urlllllllllllllllllllllllll", apiURL)
	headers["Content-Type"] = "application/json"

	response, err := utils.APIRequest(http.MethodHead, apiURL, headers, nil)
	if err != nil {
		log.Printf("failed to connect language service :failed to make API request: %v", err)
		return false, consts.ErrUtilityServiceConnectionLost
	}
	defer response.Body.Close()
	//set data to cache  using cache key
	_, err = cache.Set(cachekey, response.StatusCode == http.StatusOK, consts.CacheExpiryTime)
	if err != nil {
		log.Printf("unable to connect to cache service: %s\n", err)
	}

	return response.StatusCode == http.StatusOK, nil

}

// function to update terms and conditions of a member by using partner_id and member id
func UpdateMemberTermsAndConditions(ctx *sql.Tx, termsAndConditionsId int, ischecked bool, partnerID string, apiUrl string) (bool, error) {

	apiURL := fmt.Sprintf("%s/members/terms-and-conditions", apiUrl)
	fmt.Println("urlllllllllllllllllllllllll", apiURL)
	headers := make(map[string]interface{})
	headers["Content-Type"] = "application/json"
	// Create the request body as a map
	body := map[string]interface{}{
		"terms_and_condition_checked": ischecked,
		"terms_and_conditions_id":     termsAndConditionsId,
		"partner_id":                  partnerID,
	}

	response, err := utils.APIRequest(http.MethodPatch, apiURL, headers, body)

	if err != nil {
		log.Print("failed to connect member service", err)
		return false, consts.ErrMemberServiceConnectionLost
	}
	defer response.Body.Close()

	return response.StatusCode == http.StatusOK, nil
}

// function to get lookup id based on its lookup type and value
func GetLookupId(ctx context.Context, cache cacheConf.Cache, lookupType string, value string, apiUrl string) (int, error) {
	var cachekey = fmt.Sprintf("%s%s%s", consts.LookupNameCacheKey, lookupType, value)

	log := logger.Log().WithContext(ctx)
	headers := make(map[string]interface{})
	// check whether the data is present in cache using cache key
	ttl, err := cache.TTL(cachekey)
	if err != nil {
		cacheErr := fmt.Errorf("unable to get TTL: %s", err)
		return 0, cacheErr
	}
	if ttl > 0 {
		databytes, cacheErr := cache.Get(cachekey).Bytes()

		if cacheErr != nil {
			log.Printf("%v", err)
			return 0, err
		}
		lookupId := gjson.GetBytes(databytes, "data.0.id").Int()
		return int(lookupId), nil
	}

	// Prepare API URL with query parameters
	// retieve data from url if data is not available in the cache
	apiURL := fmt.Sprintf("%s/lookup/type/%s?value=%s", apiUrl, url.QueryEscape(lookupType), url.QueryEscape(value))
	fmt.Println("urlllllllllllllllllllllllll", apiURL)
	// Make the API request
	headers["Content-Type"] = "application/json"

	response, err := utils.APIRequest(http.MethodGet, apiURL, headers, nil)
	if err != nil {
		log.Printf("failed to connect lookup service ,err=%s", err)
		return 0, consts.ErrUtilityServiceConnectionLost
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		log.Printf("lookup service failed :no data")
		return 0, nil
	}
	// Read response body
	body, err := io.ReadAll(response.Body)

	if err != nil {
		log.Printf("lookup service failed :failed to read response body: %v", err)
		return 0, err
	}
	lookupId := gjson.GetBytes(body, "data.0.id").Int()
	if response.StatusCode == http.StatusOK {
		_, err = cache.Set(cachekey, body, consts.CacheExpiryTime)
		if err != nil {
			log.Printf("unable to connect to cache service: %s\n", err)
		}
	}
	return int(lookupId), nil
}

// function to get lookup name based on its id
func GetLookupName(ctx context.Context, cache cacheConf.Cache, id int, apiUrl string) (string, error) {
	var cachekey = fmt.Sprintf("%s%d", consts.LookupIdCacheKey, id)

	log := logger.Log().WithContext(ctx)
	headers := make(map[string]interface{})
	// check whether the data is present in cache using cache key
	ttl, err := cache.TTL(cachekey)
	if err != nil {
		cacheErr := fmt.Errorf("unable to get TTL: %s", err)
		return "", cacheErr
	}
	if ttl > 0 {
		databytes, cacheErr := cache.Get(cachekey).Bytes()

		if cacheErr != nil {
			log.Printf("%v", err)
			return "", err
		}
		name := gjson.GetBytes(databytes, "data.Name").String()
		return name, nil
	}

	// Prepare API URL with query parameters
	// retieve data from url if data is not available in the cache
	apiURL := fmt.Sprintf("%s/lookup/%d", apiUrl, id)
	fmt.Println("urlllllllllllllllllllllllll", apiURL)
	// Make the API request
	headers["Content-Type"] = "application/json"

	response, err := utils.APIRequest(http.MethodGet, apiURL, headers, nil)
	if err != nil {
		log.Printf("lookup service failed :failed to make API request: %v", err)
		return "", consts.ErrUtilityServiceConnectionLost
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		log.Printf("lookup service failed :no data")
		return "", nil
	}
	// Read response body
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Printf("lookup service failed :failed to read response body: %v", err)
		return "", err
	}
	name := gjson.GetBytes(body, "data.Name").String()

	//set data to cache  using cache key
	if response.StatusCode == http.StatusOK {
		_, err = cache.Set(cachekey, body, consts.CacheExpiryTime)
		if err != nil {
			log.Printf("unable to connect to cache service: %s\n", err)
		}
	}
	return name, nil
}

// function to get subscription duration id based on its name
func GetSubscriptionDurationId(ctx context.Context, cache cacheConf.Cache, name string, apiUrl string) (int, error) {
	var cachekey = fmt.Sprintf("%s%s", consts.SubDurationNameCacheKey, name)
	log := logger.Log().WithContext(ctx)
	headers := make(map[string]interface{})
	// check whether the data is present in cache using cache key
	ttl, err := cache.TTL(cachekey)
	if err != nil {
		cacheErr := fmt.Errorf("unable to get TTL: %s", err)
		return 0, cacheErr
	}
	if ttl > 0 {
		databytes, cacheErr := cache.Get(cachekey).Bytes()

		if cacheErr != nil {
			log.Printf("%v", err)
			return 0, err
		}
		subcriptionDurationId := gjson.GetBytes(databytes, "data.records.0.id").Int()
		return int(subcriptionDurationId), nil
	}

	// Prepare API URL with query parameters
	// retieve data from url if data is not available in the cache
	apiURL := fmt.Sprintf("%s/subscriptions/durations?name=%s", apiUrl, url.QueryEscape(name))
	fmt.Println("urlllllllllllllllllllllllll", apiURL)
	// Make the API request
	headers["Content-Type"] = "application/json"

	response, err := utils.APIRequest(http.MethodGet, apiURL, headers, nil)
	if err != nil {
		log.Printf("subscription service failed :failed to make API request: %v", err)
		return 0, consts.ErrSubscriptionServiceConnectionLost
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		log.Printf("lookup service failed :no data")
		return 0, nil
	}
	// Read response body
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Printf("subscription service failed :failed to read response body: %v", err)
		return 0, err
	}
	subcriptionDurationId := gjson.GetBytes(body, "data.records.0.id").Int()
	//set data to cache  using cache key
	if response.StatusCode == http.StatusOK {
		_, err = cache.Set(cachekey, body, consts.CacheExpiryTime)
		if err != nil {
			log.Printf("unable to connect to cache service: %s\n", err)
		}
	}
	return int(subcriptionDurationId), nil

}

// function to get subscription duration name based on duration id
func GetSubcriptionDurationName(ctx context.Context, cache cacheConf.Cache, durationId int, apiUrl string) (string, error) {
	log := logger.Log().WithContext(ctx)
	var cachekey = fmt.Sprintf("%s%d", consts.SubDurationIdCacheKey, durationId)

	headers := make(map[string]interface{})
	// check whether the data is present in cache using cache key
	ttl, err := cache.TTL(cachekey)
	if err != nil {
		cacheErr := fmt.Errorf("unable to get TTL: %s", err)
		return "", cacheErr
	}
	if ttl > 0 {
		databytes, cacheErr := cache.Get(cachekey).Bytes()

		if cacheErr != nil {
			log.Printf("%v", err)
			return "", err
		}
		subcriptionDurationName := gjson.GetBytes(databytes, "data.name").String()
		return subcriptionDurationName, nil
	}

	// retieve data from url if data is not available in the cache
	apiURL := fmt.Sprintf("%s/subscriptions/durations/%d", apiUrl, durationId)
	fmt.Println("urlllllllllllllllllllllllll", apiURL)
	// Make the API request
	headers["Content-Type"] = "application/json"

	response, err := utils.APIRequest(http.MethodGet, apiURL, headers, nil)
	if err != nil {
		log.Printf("subscription service failed :failed to make API request: %v", err)
		return "", consts.ErrSubscriptionServiceConnectionLost
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		log.Printf("subscription service failed :no data")
		return "", nil
	}
	// Read response body
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Printf("subscription service failed :failed to read response body: %v", err)
		return "", err
	}
	subcriptionDurationName := gjson.GetBytes(body, "data.name").String()

	//set data to cache  using cache key
	if response.StatusCode == http.StatusOK {
		_, err = cache.Set(cachekey, body, consts.CacheExpiryTime)
		if err != nil {
			log.Printf("unable to connect to cache service: %s\n", err)
		}
	}
	return subcriptionDurationName, nil

}

// function to check the existence of member
func IsMemberExists(ctx context.Context, memberId string, partnerId string, apiUrl string) (bool, error) {

	apiURL := fmt.Sprintf("%s/members/%s?partner_id=%s", apiUrl, memberId, partnerId)
	fmt.Println("urlllllllllllllllllllllllll", apiURL)
	headers := make(map[string]interface{})
	headers["Content-Type"] = "application/json"

	response, err := utils.APIRequest(http.MethodHead, apiURL, headers, nil)
	if err != nil {
		log.Printf("failed to connect member service:failed to make API request: %v", err)
		return false, consts.ErrMemberServiceConnectionLost
	}
	defer response.Body.Close()
	return response.StatusCode == http.StatusOK, nil
}

// function to check the existence of a state
func IsStateIsoExists(ctx context.Context, cache cacheConf.Cache, countryCode string, stateIso string, apiUrl string) (bool, error) {
	var cachekey = fmt.Sprintf("%s%s", consts.CountryStateExistsCacheKey, stateIso)

	log := logger.Log().WithContext(ctx)
	headers := make(map[string]interface{})
	// check whether the data is present in cache using cache key
	exists, err := IsCached(cache, cachekey)
	if err != nil {
		return false, err
	}
	if exists != nil {
		return *exists, nil
	}
	// retieve data from url if data is not available in the cache
	apiURL := fmt.Sprintf("%s/countries/%s/states/%s", apiUrl, countryCode, stateIso)
	fmt.Println("urlllllllllllllllllllllllll", apiURL)
	headers["Content-Type"] = "application/json"

	response, err := utils.APIRequest(http.MethodHead, apiURL, headers, nil)
	if err != nil {
		log.Printf("failed to connect country state service:failed to make API request: %v", err)
		return false, consts.ErrUtilityServiceConnectionLost
	}
	defer response.Body.Close()
	//set data to cache  using cache key
	_, err = cache.Set(cachekey, response.StatusCode == http.StatusOK, consts.CacheExpiryTime)
	if err != nil {
		log.Printf("unable to connect to cache service: %s\n", err)
	}
	return response.StatusCode == http.StatusOK, nil
}

// function to get theme based on it id
func GetTheme(ctx context.Context, cache cacheConf.Cache, id int, apiUrl string) (string, error) {
	var cachekey = fmt.Sprintf("%s%d", consts.ThemeCacheKey, id)

	log := logger.Log().WithContext(ctx)
	headers := make(map[string]interface{})
	// check whether the data is present in cache using cache key
	ttl, err := cache.TTL(cachekey)
	if err != nil {
		cacheErr := fmt.Errorf("unable to get TTL: %s", err)
		return "", cacheErr
	}
	if ttl > 0 {
		databytes, cacheErr := cache.Get(cachekey).Bytes()

		if cacheErr != nil {
			log.Printf("%v", err)
			return "", err
		}
		themeName := gjson.GetBytes(databytes, "data.name").String()
		return themeName, nil
	}

	// retieve data from url if data is not available in the cache
	apiURL := fmt.Sprintf("%s/theme/%d", apiUrl, id)
	fmt.Println("urlllllllllllllllllllllllll", apiURL)
	// Make the API request
	headers["Content-Type"] = "application/json"

	response, err := utils.APIRequest(http.MethodGet, apiURL, headers, nil)
	if err != nil {
		log.Printf("theme service failed :failed to make API request: %v", err)
		return "", consts.ErrUtilityServiceConnectionLost
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		log.Printf("theme service failed :no data")
		return "", nil
	}
	// Read response body
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Printf("theme service failed :failed to read response body: %v", err)
		return "", err
	}
	themeName := gjson.GetBytes(body, "data.name").String()
	//set data to cache  using cache key
	if response.StatusCode == http.StatusOK {
		_, err = cache.Set(cachekey, body, consts.CacheExpiryTime)
		if err != nil {
			log.Printf("unable to connect to cache service: %s\n", err)
		}
	}
	return themeName, nil

}

// function to get genre based on its id
func GetGenre(ctx context.Context, id string, apiUrl string) (string, error) {

	log := logger.Log().WithContext(ctx)
	// Prepare API URL with query parameters
	apiURL := fmt.Sprintf("%s/genres/%s", apiUrl, id)
	fmt.Println("urlllllllllllllllllllllllll", apiURL)

	headers := make(map[string]interface{})
	headers["Content-Type"] = "application/json"
	response, err := utils.APIRequest(http.MethodGet, apiURL, headers, nil)
	if err != nil {
		log.Printf("genre service failed :failed to make API request: %v", err)
		return "", consts.ErrUtilityServiceConnectionLost
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		log.Printf("utility service failed :no data")
		return "", nil
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Printf("utility service failed :failed to read response body: %v", err)
		return "", err
	}
	genreName := gjson.GetBytes(body, "data.name").String()

	return genreName, nil

}

// function to get payment gateway name based on its id
func GetPaymentGatewayName(ctx context.Context, cache cacheConf.Cache, id string, apiUrl string) (string, error) {
	var cachekey = fmt.Sprintf("%s%s", consts.PaymentGatewayIdCacheKey, id)

	log := logger.Log().WithContext(ctx)
	// check whether the data is present in cache using cache key
	ttl, err := cache.TTL(cachekey)
	if err != nil {
		cacheErr := fmt.Errorf("unable to get TTL: %s", err)
		return "", cacheErr
	}
	if ttl > 0 {
		databytes, cacheErr := cache.Get(cachekey).Bytes()
		if cacheErr != nil {
			log.Printf("%v", err)
			return "", err
		}
		paymentGatewayName := gjson.GetBytes(databytes, "data.name").String()
		return paymentGatewayName, nil
	}

	// retieve data from url if data is not available in the cache
	apiURL := fmt.Sprintf("%s/payment_gateway/%s", apiUrl, id)
	fmt.Println("urlllllllllllllllllllllllll", apiURL)

	headers := make(map[string]interface{})
	headers["Content-Type"] = "application/json"
	response, err := utils.APIRequest(http.MethodGet, apiURL, headers, nil)
	if err != nil {
		log.Printf("payment gateway service failed :failed to make API request: %v", err)
		return "", consts.ErrUtilityServiceConnectionLost
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		log.Printf("utility service failed :no data")
		return "", nil
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Printf("utility service failed :failed to read response body: %v", err)
		return "", err
	}
	// Extract the payment gateway name using gjson
	paymentGatewayName := gjson.GetBytes(body, "data.name").String()
	//set data to cache  using cache key
	if response.StatusCode == http.StatusOK {
		_, err = cache.Set(cachekey, body, consts.CacheExpiryTime)
		if err != nil {
			log.Printf("unable to connect to cache service: %s\n", err)
		}
	}
	return paymentGatewayName, nil

}

// function too get payment gateway id based on its name
func GetPaymentGatewayId(ctx context.Context, cache cacheConf.Cache, name string, apiUrl string) (string, error) {
	var cachekey = fmt.Sprintf("%s%s", consts.PaymentGatewayNameCacheKey, name)

	log := logger.Log().WithContext(ctx)
	// check whether the data is present in cache using cache key
	ttl, err := cache.TTL(cachekey)
	if err != nil {
		cacheErr := fmt.Errorf("unable to get TTL: %s", err)
		return "", cacheErr
	}
	if ttl > 0 {
		databytes, cacheErr := cache.Get(cachekey).Bytes()

		if cacheErr != nil {
			log.Printf("%v", err)
			return "", err
		}
		paymentGatewayId := gjson.GetBytes(databytes, "data.records.0.id").String()
		return paymentGatewayId, nil
	}

	// Prepare API URL with query parameters
	// retieve data from url if data is not available in the cache
	apiURL := fmt.Sprintf("%s/payment_gateway/all?name=%s", apiUrl, name)
	fmt.Println("urlllllllllllllllllllllllll", apiURL)

	headers := make(map[string]interface{})
	headers["Content-Type"] = "application/json"
	response, err := utils.APIRequest(http.MethodGet, apiURL, headers, nil)
	if err != nil {
		log.Printf("payment gateway service failed :failed to make API request: %v", err)
		return "", consts.ErrUtilityServiceConnectionLost
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		log.Printf("utility service failed :no data")
		return "", nil
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Printf("utility service failed :failed to read response body: %v", err)
		return "", err
	}
	paymentGatewayId := gjson.GetBytes(body, "data.records.0.id").String()

	//set data to cache  using cache key
	if response.StatusCode == http.StatusOK {
		_, err = cache.Set(cachekey, body, consts.CacheExpiryTime)
		if err != nil {
			log.Printf("unable to connect to cache service: %s\n", err)
		}
	}
	return paymentGatewayId, nil

}

// function to get artist role name based on its id
func GetAristRole(ctx context.Context, id string, apiUrl string) (string, error) {
	log := logger.Log().WithContext(ctx)
	// Prepare API URL with query parameters
	apiURL := fmt.Sprintf("%s/roles/%s", apiUrl, id)
	fmt.Println("urlllllllllllllllllllllllll", apiURL)

	headers := make(map[string]interface{})
	headers["Content-Type"] = "application/json"
	response, err := utils.APIRequest(http.MethodGet, apiURL, headers, nil)
	if err != nil {
		log.Printf("artist role service failed :failed to make API request: %v", err)
		return "", consts.ErrUtilityServiceConnectionLost
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		log.Printf("utility service failed :no data")
		return "", nil
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Printf("utility service failed :failed to read response body: %v", err)
		return "", err
	}

	roleName := gjson.GetBytes(body, "data.name").String()
	return roleName, nil

}

// function to get oauth provider id based name
func GetOauthProviderId(ctx context.Context, cache cacheConf.Cache, name string, apiUrl string) (string, error) {
	var cachekey = fmt.Sprintf("%s%s", consts.OauthProviderNameCacheKey, name)

	log := logger.Log().WithContext(ctx)
	headers := make(map[string]interface{})
	// check whether the data is present in cache using cache key
	ttl, err := cache.TTL(cachekey)
	if err != nil {
		cacheErr := fmt.Errorf("unable to get TTL: %s", err)
		return "", cacheErr
	}
	if ttl > 0 {
		databytes, cacheErr := cache.Get(cachekey).Bytes()
		if cacheErr != nil {
			log.Printf("%v", err)
			return "", err
		}
		providerId := gjson.GetBytes(databytes, "data.provider_id").String()
		return providerId, nil
	}

	// Prepare API URL with query parameters
	// retieve data from url if data is not available in the cache
	apiURL := fmt.Sprintf("%s/oauth/partner?provider=%s", apiUrl, name)
	fmt.Println("urlllllllllllllllllllllllll", apiURL)
	// Make the API request
	headers["Content-Type"] = "application/json"

	response, err := utils.APIRequest(http.MethodGet, apiURL, headers, nil)
	if err != nil {
		log.Printf("oauth service failed :failed to make API request: %v", err)
		return "", consts.ErrOauthServiceConnectionLost
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		log.Printf("oauth service failed :no data")
		return "", nil
	}
	// Read response body
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Printf("oauth service failed :failed to read response body: %v", err)
		return "", err
	}
	providerId := gjson.GetBytes(body, "data.provider_id").String()
	//set data to cache  using cache key
	if response.StatusCode == http.StatusOK {
		_, err = cache.Set(cachekey, body, consts.CacheExpiryTime)
		if err != nil {
			log.Printf("unable to connect to cache service: %s\n", err)
		}
	}
	return providerId, nil

}

// function to get partner user count based on partner id
func GetPartnerUserCount(ctx context.Context, partnerId string, apiUrl string) (int, error) {
	log := logger.Log().WithContext(ctx)
	headers := make(map[string]interface{})

	// Prepare API URL with query parameters
	// retieve data from url if data is not available in the cache
	apiURL := fmt.Sprintf("%s/members/count?partner_id=%s", apiUrl, partnerId)
	fmt.Println("urlllllllllllllllllllllllll", apiURL)

	// Make the API request
	headers["Content-Type"] = "application/json"

	response, err := utils.APIRequest(http.MethodGet, apiURL, headers, nil)
	if err != nil {
		log.Printf("member service failed :failed to make API request: %v", err)
		return 0, consts.ErrMemberServiceConnectionLost
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		log.Printf("member service failed :no data")
		return 0, nil
	}
	// Read response body
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Printf("member service failed :failed to read response body: %v", err)
		return 0, err
	}
	userCount := gjson.GetBytes(body, "data.active_member_count").Int()
	return int(userCount), nil

}

// function to get all stores
func GetStores(ctx context.Context, cache cacheConf.Cache, apiUrl string) ([]entities.StoreData, error) {
	var (
		jsonResponse entities.StoreResponse
		cachekey     = consts.StoreCacheKey
	)
	log := logger.Log().WithContext(ctx)
	headers := make(map[string]interface{})
	payload := map[string]interface{}{}
	// check whether the data is present in cache using cache key
	ttl, err := cache.TTL(cachekey)
	if err != nil {
		cacheErr := fmt.Errorf("unable to get TTL: %s", err)
		return nil, cacheErr
	}
	if ttl > 0 {
		databytes, cacheErr := cache.Get(cachekey).Bytes()
		if cacheErr != nil {
			log.Printf("%v", err)
			return nil, err
		}
		err = json.Unmarshal(databytes, &jsonResponse)
		if err != nil {
			log.Printf("%v", err)
			return nil, err
		}
		return jsonResponse.StoreData, nil
	}

	// retieve data from url if data is not available in the cache
	headers["Content-Type"] = "application/json"
	apiURL := fmt.Sprintf("%s/stores", apiUrl)
	response, err := utils.APIRequest(http.MethodGet, apiURL, headers, payload)
	if err != nil {
		log.Printf("failed to connect store service :%v", err)
		return nil, consts.ErrStoreServiceConnectionLost
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		log.Printf("store service failed :no data")
		return nil, err
	}
	// Read response body
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Printf("store service failed :failed to read response body: %v", err)
		return nil, err
	}

	// Unmarshal JSON response
	err = json.Unmarshal(body, &jsonResponse)
	if err != nil {
		log.Printf("store service failed :failed to unmarshal JSON response: %v", err)
		return nil, err
	}
	if response.StatusCode == http.StatusOK {
		_, err = cache.Set(cachekey, body, consts.CacheExpiryTime)
		if err != nil {
			log.Printf("unable to connect to cache service: %s\n", err)
		}
	}
	//set data to cache  using cache key
	return jsonResponse.StoreData, nil
}

// function to check whether data is already cached or not
func IsCached(cache cacheConf.Cache, cacheKey string) (*bool, error) {

	ttl, err := cache.TTL(cacheKey)
	if err != nil {
		cacheErr := fmt.Errorf("unable to get TTL: %s", err)
		return nil, cacheErr
	}
	if ttl > 0 {
		exists, cacheErr := cache.Get(cacheKey).Bool()
		if cacheErr != nil {
			log.Printf("%v", err)
			return nil, err
		}
		return &exists, nil
	}
	return nil, nil
}

// Function to get all genres
func GetAllProductTypes(ctxt context.Context, cache cacheConf.Cache, apiUrl string) ([]entities.ProductTypes, error) {
	var (
		allRecords []entities.ProductTypes
		// records    []entities.ProductTypes
		// response entities.ProductTypeResponseMetaData
	)
	apiURL := "http://10.1.0.90:8080/api/v1.0/products/types"
	err := utils.ApiRequestWithPages(context.Background(), *utils.NewRequest(
		utils.WithHost(apiURL),
		utils.WithMethod(http.MethodGet),
	), func(ao *utils.APIResponseOutput, lastPage bool) bool {
		dataStr, ok := ao.Data.(string)
		if !ok {
			return !lastPage
		}
		var response entities.DefaultApiResponse
		if err := json.Unmarshal([]byte(dataStr), &response); err != nil {
			fmt.Println("Error unmarshalling data:", err)
			return !lastPage
		}
		allRecords = append(allRecords, response.Data.Records...)

		return !lastPage

	})

	if err != nil {
		return nil, err
	}
	fmt.Println("recordssssssssss", allRecords)
	return allRecords, nil
}

func NewActivityLog(memberID, action string, data map[string]interface{}) models.ActivityLog {
	return models.ActivityLog{
		MemberID: memberID,
		Action:   action,
		Data:     data,
	}
}
