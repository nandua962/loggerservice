package consts

var ErrorCodeMap = map[string]interface{}{
	"genre": map[string]interface{}{
		"post": map[string]interface{}{
			"name": map[string]interface{}{
				"genre_exists": "G1001",
				"required":     "G1002",
				"length":       "G1003",
			},
		},
		"delete": map[string]interface{}{
			"common_message": map[string]interface{}{
				"delete_genre":  "G1004",
				"invalid_genre": "G1013",
			},
		},
		"patch": map[string]interface{}{
			"name": map[string]interface{}{
				"genre_exists": "G1005",
				"required":     "G1006",
				"length":       "G1007",
			},
			"common_message": map[string]interface{}{
				"invalid_genre": "G1014",
			},
		},
		"get": map[string]interface{}{
			"order": map[string]interface{}{
				"invalid": "G1008",
			},
			"sort": map[string]interface{}{
				"invalid":   "G1009",
				"arguments": "G1010",
			},
			"limit": map[string]interface{}{
				"invalid": "G1011",
			},
			"page": map[string]interface{}{
				"invalid": "G1012",
			},
			"genre_id": map[string]interface{}{
				"invalid": "G1013",
			},
		},
	},
	"roles": map[string]interface{}{
		"post": map[string]interface{}{
			"name": map[string]interface{}{
				"role_exists": "R1001",
				"required":    "R1002",
				"length":      "R1003",
			},
			"language_label": map[string]interface{}{
				"language_exists": "R1004",
			},
		},
		"delete": map[string]interface{}{
			"common_message": map[string]interface{}{
				"deleted":      "R1005",
				"invalid_role": "R1006",
				"not_found":    "R1007",
			},
			"role": map[string]interface{}{
				"deleted":   "R1008",
				"invalid":   "R1009",
				"not_found": "R1010",
			},
		},
		"patch": map[string]interface{}{
			"name": map[string]interface{}{
				"role_exists": "R1011",
				"required":    "R1012",
				"length":      "R1013",
			},
			"common_message": map[string]interface{}{
				"invalid_role": "R1014",
			},
			"language_label": map[string]interface{}{
				"language_exists": "R1015",
			},
		},
		"get": map[string]interface{}{
			"order": map[string]interface{}{
				"invalid": "R1016",
			},
			"sort": map[string]interface{}{
				"invalid":   "R1017",
				"arguments": "R1018",
			},
			"limit": map[string]interface{}{
				"invalid": "R1019",
			},
			"page": map[string]interface{}{
				"invalid": "R1020",
			},
			"role_id": map[string]interface{}{
				"invalid": "R1021",
			},
		},
	},
	"currencies": map[string]interface{}{
		"get": map[string]interface{}{
			"order": map[string]interface{}{
				"invalid": "C1000",
			},
			"sort": map[string]interface{}{
				"invalid":   "C1002",
				"arguments": "C1003",
			},
			"limit": map[string]interface{}{
				"invalid": "C1004",
			},
			"page": map[string]interface{}{
				"invalid": "C1005",
			},
			"currency_id": map[string]interface{}{
				"invalid": "C1006",
			},
			"iso": map[string]interface{}{
				"invalid": "C1007",
				"length":  "C1008",
			},
		},
	},
	"countries": map[string]interface{}{
		"get": map[string]interface{}{
			"order": map[string]interface{}{
				"invalid": "CO1000",
			},
			"sort": map[string]interface{}{
				"invalid":   "CO1002",
				"arguments": "CO1003",
			},
			"limit": map[string]interface{}{
				"invalid": "CO1004",
			},
			"page": map[string]interface{}{
				"invalid": "CO1005",
			},
			"iso": map[string]interface{}{
				"invalid":  "CO1006",
				"length":   "CO1007",
				"required": "CO1008",
			},
		},
	},
	"states": map[string]interface{}{
		"get": map[string]interface{}{
			"order": map[string]interface{}{
				"invalid": "S1000",
			},
			"sort": map[string]interface{}{
				"invalid":   "S1002",
				"arguments": "S1003",
			},
			"limit": map[string]interface{}{
				"invalid": "S1004",
			},
			"page": map[string]interface{}{
				"invalid": "S1005",
			},
			"country_id": map[string]interface{}{
				"invalid": "S1006",
			},
		},
	},
	"languages": map[string]interface{}{
		"get": map[string]interface{}{
			"order": map[string]interface{}{
				"invalid": "L1000",
			},
			"sort": map[string]interface{}{
				"invalid":   "L1002",
				"arguments": "L1003",
			},
			"limit": map[string]interface{}{
				"invalid": "L1004",
			},
			"page": map[string]interface{}{
				"invalid": "L1005",
			},
			"status": map[string]interface{}{
				"invalid": "L1005",
			},
			"language_id": map[string]interface{}{
				"invalid": "L1006",
			},
			"code": map[string]interface{}{
				"invalid": "L1007",
				"length":  "L1008",
			},
		},
	},
	"theme": map[string]interface{}{
		"get": map[string]interface{}{
			"theme_id": map[string]interface{}{
				"invalid": "T1000",
			},
		},
	},
	"lookup": map[string]interface{}{
		"get": map[string]interface{}{
			"order": map[string]interface{}{
				"invalid": "LO1001",
			},
			"sort": map[string]interface{}{
				"invalid":   "LO1002",
				"arguments": "LO1003",
			},
			"limit": map[string]interface{}{
				"invalid": "LO1004",
			},
			"page": map[string]interface{}{
				"invalid": "LO1005",
			},
			"lookup_id": map[string]interface{}{
				"invalid": "LO1006",
			},
			"lookup_type_id": map[string]interface{}{
				"invalid": "LO1007",
			},
			"name": map[string]interface{}{
				"invalid": "LO1008",
			},
		},
	},
	"gateway": map[string]interface{}{
		"get": map[string]interface{}{
			"payment_gateway_id": map[string]interface{}{
				"invalid": "G1000",
			},
		},
	},
}
