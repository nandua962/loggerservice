package utilities

import (
	"partner/internal/entities"
	"time"

	"github.com/google/uuid"
)

// function to return conditions and columns that are needed for partner oauth select query
func PartnerOauthConditions(conditions entities.PartneOauthConditions) ([]string, map[string]interface{}) {
	columns := []string{
		"client_id",
		"client_secret",
		"redirect_uri",
		"scope",
		"oauth_provider_id",
		"access_token_endpoint",
	}
	condition := map[string]interface{}{
		"partner_id": conditions.PartnerId,
	}
	return columns, condition
}

// function to map  partner data
func GeneratePartnerData(partner entities.Partner) map[string]interface{} {
	data := map[string]interface{}{
		"name": partner.Name,
		"logo": partner.Logo,
		"url":  partner.URL,
		// "favicon":                          partner.Favicon,
		// "loader":                           partner.Loader,
		// "login_page_logo":                  partner.LoginPageLogo,
		"contact_person":                 partner.ContactDetails.ContactPerson,
		"email":                          partner.ContactDetails.Email,
		"noreply_email":                  partner.ContactDetails.NoReplyEmail,
		"feedback_email":                 partner.ContactDetails.FeedbackEmail,
		"support_email":                  partner.ContactDetails.SupportEmail,
		"background_color":               partner.BackgroundColor,
		"background_image":               partner.BackgroundImage,
		"website_url":                    partner.WebsiteURL,
		"browser_title":                  partner.BrowserTitle,
		"profile_url":                    partner.ProfileURL,
		"payment_url":                    partner.PaymentURL,
		"landing_page":                   partner.LandingPage,
		"theme_id":                       partner.ThemeID,
		"enable_mail":                    partner.EnableMail,
		"login_type_id":                  partner.LoginTypeID,
		"member_pay_to_partner":          partner.MemberPayToPartner,
		"address":                        partner.AddressDetails.Address,
		"city":                           partner.AddressDetails.City,
		"street":                         partner.AddressDetails.Street,
		"postal_code":                    partner.AddressDetails.PostalCode,
		"payout_min_limit":               partner.PaymentGatewayDetails.PayoutMinLimit,
		"max_remittance_per_month":       partner.PaymentGatewayDetails.MaxRemittancePerMonth,
		"member_grace_period":            partner.MemberGracePeriodID,
		"expiry_warning_count":           partner.ExpiryWarningCount,
		"album_review_email":             partner.AlbumReviewEmail,
		"site_info":                      partner.SiteInfo,
		"outlets_processing_duration":    partner.OutletsProcessingDuration,
		"free_plan_limit":                partner.FreePlanLimit,
		"product_review":                 partner.ProductReviewID,
		"language_code":                  partner.Language,
		"payout_target_currency_id":      partner.PayoutTargetCurrencyID,
		"business_model_id":              partner.BusinessModelID,
		"country_code":                   partner.AddressDetails.Country,
		"state_code":                     partner.AddressDetails.State,
		"default_currency_id":            partner.DefaultCurrencyID,
		"default_price_code_currency_id": partner.DefaultPriceCodeCurrencyID,
		"music_language_code":            partner.MusicLanguage,
		"member_default_country_code":    partner.MemberDefaultCountry,
		"mobile_verify_interval":         partner.MobileVerifyInterval,
		"default_payment_gateway_id":     partner.DefaultPaymentGatewayId,
	}
	return data
}

// function to map  partner data
func UpdatePartnerData(partner entities.Partner, partnerId string, memberID uuid.UUID) map[string]interface{} {
	data := map[string]interface{}{
		"name": partner.Name,
		"logo": partner.Logo,
		"url":  partner.URL,
		// "favicon":                          partner.Favicon,
		// "loader":                           partner.Loader,
		// "login_page_logo":                  partner.LoginPageLogo,
		"contact_person":                 partner.ContactDetails.ContactPerson,
		"email":                          partner.ContactDetails.Email,
		"noreply_email":                  partner.ContactDetails.NoReplyEmail,
		"feedback_email":                 partner.ContactDetails.FeedbackEmail,
		"support_email":                  partner.ContactDetails.SupportEmail,
		"background_color":               partner.BackgroundColor,
		"background_image":               partner.BackgroundImage,
		"website_url":                    partner.WebsiteURL,
		"browser_title":                  partner.BrowserTitle,
		"profile_url":                    partner.ProfileURL,
		"payment_url":                    partner.PaymentURL,
		"landing_page":                   partner.LandingPage,
		"theme_id":                       partner.ThemeID,
		"enable_mail":                    partner.EnableMail,
		"login_type_id":                  partner.LoginTypeID,
		"member_pay_to_partner":          partner.MemberPayToPartner,
		"address":                        partner.AddressDetails.Address,
		"city":                           partner.AddressDetails.City,
		"street":                         partner.AddressDetails.Street,
		"postal_code":                    partner.AddressDetails.PostalCode,
		"payout_min_limit":               partner.PaymentGatewayDetails.PayoutMinLimit,
		"max_remittance_per_month":       partner.PaymentGatewayDetails.MaxRemittancePerMonth,
		"member_grace_period":            partner.MemberGracePeriodID,
		"expiry_warning_count":           partner.ExpiryWarningCount,
		"album_review_email":             partner.AlbumReviewEmail,
		"site_info":                      partner.SiteInfo,
		"outlets_processing_duration":    partner.OutletsProcessingDuration,
		"free_plan_limit":                partner.FreePlanLimit,
		"product_review":                 partner.ProductReviewID,
		"language_code":                  partner.Language,
		"payout_target_currency_id":      partner.PayoutTargetCurrencyID,
		"business_model_id":              partner.BusinessModelID,
		"country_code":                   partner.AddressDetails.Country,
		"state_code":                     partner.AddressDetails.State,
		"default_currency_id":            partner.DefaultCurrencyID,
		"default_price_code_currency_id": partner.DefaultPriceCodeCurrencyID,
		"music_language_code":            partner.MusicLanguage,
		"member_default_country_code":    partner.MemberDefaultCountry,
		"mobile_verify_interval":         partner.MobileVerifyInterval,
		"partner_plan_id":                partner.SubscriptionPlanDetails.ID,
		"payment_gateway_id":             partner.PaymentGatewayDetails.PaymentGateWayID,
		"plan_start_date":                partner.SubscriptionPlanDetails.StartDate,
		"plan_launch_date":               partner.SubscriptionPlanDetails.LaunchDate,
		"default_payment_gateway_id":     partner.DefaultPaymentGatewayId,
		"updated_on":                     time.Now(),
		"updated_by":                     memberID,
	}
	return data
}

func OauthCredentialData(oauthCredential entities.PartnerOauthCredential) map[string]interface{} {
	data := map[string]interface{}{
		"client_id":         oauthCredential.ClientId,
		"client_secret":     oauthCredential.ClientSecret,
		"partner_id":        oauthCredential.PartnerId,
		"oauth_provider_id": oauthCredential.ProviderId,
	}
	return data
}
