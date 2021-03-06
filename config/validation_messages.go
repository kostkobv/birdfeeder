package config

// ValidationMessages vocabulary
var ValidationMessages = map[string]string{
	"required":              "must have a value",
	"msisdn":                "should be a valid MSISDN",
	"textoriginator|msisdn": "use valid MSISDN or alphanumeric value (max. 11 symbols long)",
	"textoriginator":        "use alphanumeric value (max. 11 symbols long)",
	"max":                   "outreached limit for characters amount (max. 1377 for plain and 603 for unicode)",
}
