package consts

// Add values to the map
var ErrorCodeMap = map[string]interface{}{
	"partners": map[string]interface{}{
		"post": map[string]interface{}{
			"name": map[string]interface{}{
				"Required":       "P1001",
				"already_exists": "P1002",
				"Limit_exceeds":  "P1003",
				"invalid":        "P1004",
			},
			"url": map[string]interface{}{
				"Required":       "P1005",
				"invalid":        "P1006",
				"already_exists": "P1007",
				"Limit_exceeds":  "P1154",
			},
			"logo": map[string]interface{}{
				"Required": "P1008",
				"invalid":  "P1009",
			},
			"business_model": map[string]interface{}{
				"Required": "P2001",
				"invalid":  "P1010",
			},
			"contact_person": map[string]interface{}{
				"Required":      "P1011",
				"Limit_exceeds": "P1012",
			},
			"email": map[string]interface{}{
				"Required":       "P1013",
				"invalid":        "P1014",
				"already_exists": "P1015",
				"Limit_exceeds":  "P1155",
			},
			"no_reply_email": map[string]interface{}{
				"Required":       "P1016",
				"invalid":        "P1017",
				"already_exists": "P1018",
				"Limit_exceeds":  "P1156",
			},
			"support_email": map[string]interface{}{
				"Required":       "P1019",
				"already_exists": "P1020",
				"invalid":        "P1021",
				"Limit_exceeds":  "P1157",
			},
			"feedback_email": map[string]interface{}{
				"Required":       "P1022",
				"already_exists": "P1023",
				"invalid":        "P1024",
			},
			"album_review_email": map[string]interface{}{
				"Required":       "P1025",
				"already_exists": "P1026",
				"invalid":        "P1027",
				"Limit_exceeds":  "P1158",
			},
			"address": map[string]interface{}{
				"Required":      "P1028",
				"Limit_exceeds": "P1029",
			},
			"country": map[string]interface{}{
				"Required": "P2003",
				"invalid":  "P1030",
			},
			"payment_gateway": map[string]interface{}{
				"invalid":  "P1031",
				"Required": "P1032",
			},
			"favicon": map[string]interface{}{
				"Required": "P1033",
			},
			"free_plan_limit": map[string]interface{}{
				"invalid": "P1034",
			},
			"expiry_warning_count": map[string]interface{}{
				"invalid": "P1035",
			},
			"max_remittance_per_month": map[string]interface{}{
				"invalid": "P1036",
			},
			"outlets_processing_duration": map[string]interface{}{
				"Limit_exceeds": "P1037",
			},
			"background_color": map[string]interface{}{
				"invalid": "P1038",
			},
			"background_image": map[string]interface{}{
				"invalid": "P2234",
			},
			"currency": map[string]interface{}{
				"invalid": "P1039",
			},
			"payout_target_currency": map[string]interface{}{
				"invalid": "P1040",
			},
			"state": map[string]interface{}{
				"invalid": "P1041",
			},
			"default_subscription_currency": map[string]interface{}{
				"invalid": "P1042",
			},
			"default_price_code_currency": map[string]interface{}{
				"invalid": "P1043",
			},
			"music_language": map[string]interface{}{
				"invalid": "P1044",
			},
			"member_default_country": map[string]interface{}{
				"invalid": "P1045",
			},
			"member_grace_period": map[string]interface{}{
				"invalid": "P1046",
			},
			"language": map[string]interface{}{
				"invalid": "P1047",
			},
			"product_review": map[string]interface{}{
				"invalid": "P1048",
			},
			"login_type": map[string]interface{}{
				"invalid": "P1049",
			},
			"noreply_email": map[string]interface{}{
				"Required":       "P1050",
				"already_exists": "P1051",
				"invalid":        "P1052",
				"Limit_exceeds":  "P1159",
			},
			"website_url": map[string]interface{}{
				"invalid":        "P1053",
				"already_exists": "P1054",
				"Limit_exceeds":  "P1160",
			},
			"landing_page": map[string]interface{}{
				"already_exists": "P1055",
				"invalid":        "P1056",
				"Limit_exceeds":  "P1161",
			},
			"profile_url": map[string]interface{}{
				"invalid":       "P1057",
				"Limit_exceeds": "P1162",
			},
			"payment_url": map[string]interface{}{
				"invalid":       "P1058",
				"Limit_exceeds": "P1163",
			},
			"city": map[string]interface{}{
				"Limit_exceeds": "P1059",
			},
			"street": map[string]interface{}{
				"Limit_exceeds": "P1060",
			},
			"browser_title": map[string]interface{}{
				"Limit_exceeds": "P1061",
			},
			"payout_currency": map[string]interface{}{
				"invalid": "P1062",
			},
			"postal_code": map[string]interface{}{
				"invalid": "P1063",
			},
			"payment_gateway_email": map[string]interface{}{
				"invalid":  "P1064",
				"Required": "P1065",
			},
			"default_payin_currency": map[string]interface{}{
				"invalid":  "P1066",
				"Required": "P1067",
			},
			"default_payout_currency": map[string]interface{}{
				"invalid":  "P1068",
				"Required": "P1069",
			},
			"client_id": map[string]interface{}{
				"Required": "P1070",
			},
			"client_secret": map[string]interface{}{
				"Required": "P1071",
			},
			"site_info": map[string]interface{}{
				"Limit_exceeds": "P1072",
			},
			"default_payment_gateway": map[string]interface{}{
				"Required":  "P1164",
				"invalid":   "P1165",
				"not_found": "P1675",
			},
			"mobile_verify_interval": map[string]interface{}{
				"invalid": "P3456",
			},
			"theme_id": map[string]interface{}{
				"invalid": "P3456",
			},
			"payout_min_limit": map[string]interface{}{
				"invalid": "P1822",
			},
		},
		"patch": map[string]interface{}{
			"email": map[string]interface{}{
				"invalid":        "P1073",
				"Required":       "P1074",
				"already_exists": "P1075",
				"Limit_exceeds":  "P1162",
			},
			"url": map[string]interface{}{
				"invalid":        "P1076",
				"Required":       "P1077",
				"already_exists": "P1078",
				"Limit_exceeds":  "P1163",
			},
			"website_url": map[string]interface{}{
				"invalid":        "P1079",
				"already_exists": "P1080",
				"Limit_exceeds":  "P1164",
			},
			"profile_url": map[string]interface{}{
				"invalid":       "P1081",
				"Limit_exceeds": "P1165",
			},
			"payment_url": map[string]interface{}{
				"invalid":       "P1082",
				"Limit_exceeds": "P1166",
			},
			"payout_min_limit": map[string]interface{}{
				"invalid": "P1822",
			},
			"landing_page": map[string]interface{}{
				"invalid":        "P1083",
				"already_exists": "P1084",
				"Limit_exceeds":  "P1167",
			},
			"no_reply_email": map[string]interface{}{
				"invalid":        "P1085",
				"Required":       "P1086",
				"already_exists": "P1087",
				"Limit_exceeds":  "P1168",
			},
			"feedback_email": map[string]interface{}{
				"invalid":        "P1088",
				"Required":       "P1089",
				"already_exists": "P1090",
				"Limit_exceeds":  "P1169",
			},
			"support_email": map[string]interface{}{
				"invalid":        "P1091",
				"already_exists": "P1092",
				"Required":       "P1093",
			},
			"album_review_email": map[string]interface{}{
				"invalid":        "P1094",
				"Required":       "P1095",
				"already_exists": "P1096",
				"Limit_exceeds":  "P1170",
			},
			"name": map[string]interface{}{
				"already_exists": "P1097",
				"invalid":        "P1098",
				"Limit_exceeds":  "P1099",
				"Required":       "P1101",
			},
			"language": map[string]interface{}{
				"invalid":  "P1102",
				"Required": "p0987",
			},
			"country": map[string]interface{}{
				"Required": "P2004",
				"invalid":  "P1103",
			},
			"state": map[string]interface{}{
				"invalid": "P1104",
			},
			"music_language": map[string]interface{}{
				"invalid": "P1105",
			},
			"member_default_country": map[string]interface{}{
				"invalid": "P1106",
			},
			"payout_currency": map[string]interface{}{
				"invalid": "P1107",
			},
			"default_price_code_currency": map[string]interface{}{
				"invalid": "P1108",
			},
			"default_subscription_currency": map[string]interface{}{
				"invalid": "P1109",
			},
			"payout_target_currency": map[string]interface{}{
				"invalid": "P1110",
			},
			"member_grace_period": map[string]interface{}{
				"invalid": "P1111",
			},
			"plan_id": map[string]interface{}{
				"invalid": "P1112",
			},
			"background_color": map[string]interface{}{
				"invalid": "P1113",
			},
			"background_image": map[string]interface{}{
				"invalid": "P1114",
			},
			"partner_id": map[string]interface{}{
				"not_found": "P1115",
				"invalid":   "P3452",
			},
			"business_model": map[string]interface{}{
				"Required": "P2002",
				"invalid":  "P1116",
			},
			"product_review": map[string]interface{}{
				"invalid": "P1117",
			},
			"login_type": map[string]interface{}{
				"invalid": "P1118",
			},
			"expiry_warning_count": map[string]interface{}{
				"invalid": "P1119",
			},
			"free_plan_limit": map[string]interface{}{
				"invalid": "P1120",
			},
			"max_remittance_per_month": map[string]interface{}{
				"invalid": "P1121",
			},
			"outlets_processing_duration": map[string]interface{}{
				"Limit_exceeds": "P1122",
			},
			"contact_person": map[string]interface{}{
				"invalid":       "P1123",
				"Limit_exceeds": "P1124",
			},
			"browser_title": map[string]interface{}{
				"Limit_exceeds": "P1125",
			},
			"address": map[string]interface{}{
				"Limit_exceeds": "P1126",
				"Required":      "P1127",
			},
			"street": map[string]interface{}{
				"Limit_exceeds": "P1128",
			},
			"city": map[string]interface{}{
				"Limit_exceeds": "P1129",
			},
			"logo": map[string]interface{}{
				"Required": "p1234",
				"invalid":  "P1130",
			},
			"postal_code": map[string]interface{}{
				"invalid": "P1131",
			},
			"noreply_email": map[string]interface{}{
				"Required":       "P1132",
				"already_exists": "P1133",
				"invalid":        "P1134",
				"Limit_exceeds":  "P1171",
			},
			"default_payin_currency": map[string]interface{}{
				"invalid":  "P1135",
				"Required": "P1136",
			},
			"default_payout_currency": map[string]interface{}{
				"invalid":  "P1137",
				"Required": "P1138",
			},
			"payment_gateway": map[string]interface{}{
				"Required": "P1139",
			},
			"payment_gateway_email": map[string]interface{}{
				"invalid":  "P1240",
				"Required": "P1140",
			},
			"client_id": map[string]interface{}{
				"Required": "P1141",
			},
			"client_secret": map[string]interface{}{
				"Required": "P1142",
			},
			"site_info": map[string]interface{}{
				"Limit_exceeds": "P1143",
			},
			"terms_and_conditions_description": map[string]interface{}{
				"Required": "P1167",
			},
			"terms_and_conditions_name": map[string]interface{}{
				"Limit_exceeds": "P1299",
				"Required":      "P1201",
			},
			"terms_and_conditions_language": map[string]interface{}{
				"Required": "P1202",
			},
			"default_payment_gateway": map[string]interface{}{
				"Required":  "P1166",
				"invalid":   "P1167",
				"not_found": "P1612",
			},
			"mobile_verify_interval": map[string]interface{}{
				"invalid": "P3456",
			},
			"theme_id": map[string]interface{}{
				"invalid": "P3456",
			},
			"member_id": map[string]interface{}{
				"Required":  "P1192",
				"not_found": "P1193",
				"invalid":   "P1194",
			},
		},
		"get": map[string]interface{}{
			"page": map[string]interface{}{
				"invalid": "P1144",
			},
			"partner_id": map[string]interface{}{
				"invalid":        "P2453",
				"already_exists": "P1145",
				"not_found":      "P2762",
			},
			"key": map[string]interface{}{
				"invalid": "P1146",
			},
			"sort": map[string]interface{}{
				"invalid":  "P1147",
				"Required": "P1148",
			},
			"order": map[string]interface{}{
				"invalid": "P1149",
			},
			"active": map[string]interface{}{
				"invalid": "P1150",
			},
			"limit": map[string]interface{}{
				"invalid": "P1151",
			},
			"status": map[string]interface{}{
				"invalid": "P1152",
			},
			"country": map[string]interface{}{
				"invalid": "P1153",
			},
			"encryption_key": map[string]interface{}{
				"invalid": "P2121",
			},
			"oauth_provider": map[string]interface{}{
				"invalid":  "P00345",
				"Required": "P12367",
			},
			"member_id": map[string]interface{}{
				"Required":  "P1192",
				"not_found": "P1193",
				"invalid":   "P1194",
			},
			"fields": map[string]interface{}{
				"invalid": "P8888",
			},
		},
		"delete": map[string]interface{}{
			"partner_id": map[string]interface{}{
				"invalid":   "P1144",
				"not_found": "P4563",
			},
			"role_id": map[string]interface{}{
				"not_found": "P4500",
			},
			"genre_id": map[string]interface{}{
				"not_found": "P4511",
			},
			"member_id": map[string]interface{}{
				"Required":  "P1192",
				"not_found": "P1193",
				"invalid":   "P1194",
			},
		},
	},
}
