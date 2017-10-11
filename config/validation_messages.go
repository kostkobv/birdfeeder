package config

// ValidationMessages vocabulary
var ValidationMessages = map[string]string{
	"required":              "must have a value",
	"msisdn":                "should be a valid MSISDN",
	"textoriginator|msisdn": "use valid MSISDN or alphanumeric value (max. 11 symbols long)",
}
