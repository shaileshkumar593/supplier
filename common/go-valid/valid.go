package valid

import (
	"errors"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	// URI regexp pattern for URI pattern
	URI = `(\/([a-zA-Z0-9_]+|\{[a-zA-Z0-9_]+\}))+`
	// Email regexp pattern for RFC 5322 (email) electronic mail address
	Email = `(?:[a-zA-Z0-9!#$%&'*+/=?^_` + "`" +
		`{|}~-]+(?:\.[a-zA-Z0-9!#$%&'*+/=?^_` + "`" +
		`{|}~-]+)*|"(?:[\x01-\x08\x0b\x0c\x0e-\x1f\x21\x23-\x5b\x5d-\x7f]|\\[\x01-\x09\x0b\x0c\x0e-\x7f])*")@(?:(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?|\[(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?|[a-z0-9-]*[a-z0-9]:(?:[\x01-\x08\x0b\x0c\x0e-\x1f\x21-\x5a\x53-\x7f]|\\[\x01-\x09\x0b\x0c\x0e-\x7f])+)\])`
	// IPv4 regexp pattern for RFC 791 (ipv4) Internet protocol version 4
	IPv4 = `(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)`
	// HashMD5 regexp for (md5) hashed URI
	HashMD5 = `[a-f0-9]{32}`
	// HTTPMethod regexp for RFC 7231 and RFC 5789 HTTP request methods
	HTTPMethod = `get|put|patch|post|delete|options|head`
	// CSAlphaNum regexp for Case-sensitive Alphanumeric
	CSAlphaNum = `[a-z0-9]+`
	// CIAlphaNum regexp for Case-insensitive Alphanumeric
	CIAlphaNum = `[a-zA-Z0-9]+`
	// CIAlphaNumSpace regexp for Case-insensitive Alphanumeric and Space
	CIAlphaNumSpace = `[a-zA-Z0-9 ]+`
	// CIAlpha regexp for Case-insensitive Alphabetical
	CIAlpha = `[a-zA-Z]+`
	//ParamRequired  Parameter is a required field
	ParamRequired = "Parameter is a required field"
	//ParamURLInvalid  Parameter is not a valid url
	ParamURLInvalid = "Parameter is not valid"
	//ParamGreaterThan ...
	ParamGreaterThan = "Parameter is invalid"
	//ParamGender Parameter must be male or female
	ParamGender = "Parameter must be male or female"
	// CalculationMode acceptable values for calcumation mode
	CalculationMode = `source|receive`
	// DeliveryMethod ...
	DeliveryMethod = `CP|MP|EW|BD|HD|IB`
	//Birthday ...
	Birthday = `[0-9]{4}-[0-9]{2}-[0-9]{2}`
	//AnnexTypes ...
	AnnexTypes = `receipt|additional_fields|error|confirm|status|cancel`
	//CountryCodeAlpha3 ISO 3166-1 alpha-3
	CountryCodeAlpha3 = `ABW|AFG|AGO|AIA|ALA|ALB|AND|ARE|ARG|ARM|ASM|ATA|ATF|ATG|AUS|AUT|AZE|BDI|BEL|BEN|BES|BFA|BGD|BGR|BHR|BHS|BIH|BLM|BLR|BLZ|BMU|BOL|BRA|BRB|BRN|BTN|BVT|BWA|CAF|CAN|CCK|CHE|CHL|CHN|CIV|CMR|COD|COG|COK|COL|COM|CPV|CRI|CUB|CUW|CXR|CYM|CYP|CZE|DEU|DJI|DMA|DNK|DOM|DZA|ECU|EGY|ERI|ESH|ESP|EST|ETH|FIN|FJI|FLK|FRA|FRO|FSM|GAB|GBR|GEO|GGY|GHA|GIB|GIN|GLP|GMB|GNB|GNQ|GRC|GRD|GRL|GTM|GUF|GUM|GUY|HKG|HMD|HND|HRV|HTI|HUN|IDN|IMN|IND|IOT|IRL|IRN|IRQ|ISL|ISR|ITA|JAM|JEY|JOR|JPN|KAZ|KEN|KGZ|KHM|KIR|KNA|KOR|KWT|LAO|LBN|LBR|LBY|LCA|LIE|LKA|LSO|LTU|LUX|LVA|MAC|MAF|MAR|MCO|MDA|MDG|MDV|MEX|MHL|MKD|MLI|MLT|MMR|MNE|MNG|MNP|MOZ|MRT|MSR|MTQ|MUS|MWI|MYS|MYT|NAM|NCL|NER|NFK|NGA|NIC|NIU|NLD|NOR|NPL|NRU|NZL|OMN|PAK|PAN|PCN|PER|PHL|PLW|PNG|POL|PRI|PRK|PRT|PRY|PSE|PYF|QAT|REU|ROU|RUS|RWA|SAU|SDN|SEN|SGP|SGS|SHN|SJM|SLB|SLE|SLV|SMR|SOM|SPM|SRB|SSD|STP|SUR|SVK|SVN|SWE|SWZ|SXM|SYC|SYR|TCA|TCD|TGO|THA|TJK|TKL|TKM|TLS|TON|TTO|TUN|TUR|TUV|TWN|TZA|UGA|UKR|UMI|URY|USA|UZB|VAT|VCT|VEN|VGB|VIR|VNM|VUT|WLF|WSM|YEM|ZAF|ZMB|ZW`
	// CurrencyCode ISO 4217 Currency Codes
	CurrencyCode = `AED|AFN|ALL|AMD|ANG|AOA|ARS|AUD|AWG|AZN|BAM|BBD|BDT|BGN|BHD|BIF|BMD|BND|BOB|BOV|BRL|BSD|BTN|BWP|BYR|BZD|CAD|CDF|CHE|CHF|CHW|CLF|CLP|CNY|COP|COU|CRC|CUC|CUP|CVE|CZK|DJF|DKK|DOP|DZD|EGP|ERN|ETB|EUR|FJD|FKP|GBP|GEL|GHS|GIP|GMD|GNF|GTQ|GYD|HKD|HNL|HRK|HTG|HUF|IDR|ILS|INR|IQD|IRR|ISK|JMD|JOD|JPY|KES|KGS|KHR|KMF|KPW|KRW|KWD|KYD|KZT|LAK|LBP|LKR|LRD|LSL|LTL|LVL|LYD|MAD|MDL|MGA|MKD|MMK|MNT|MOP|MRO|MUR|MVR|MWK|MXN|MXV|MYR|MZN|NAD|NGN|NIO|NOK|NPR|NZD|OMR|PAB|PEN|PGK|PHP|PKR|PLN|PYG|QAR|RON|RSD|RUB|RWF|SAR|SBD|SCR|SDG|SEK|SGD|SHP|SLL|SOS|SRD|SSP|STD|SVC|SYP|SZL|THB|TJS|TMT|TND|TOP|TRY|TTD|TWD|TZS|UAH|UGX|USD|USN|USS|UYI|UYU|UZS|VEF|VND|VUV|WST|XAF|XAG|XAU|XBA|XBB|XBC|XBD|XCD|XDR|XFU|XOF|XPD|XPF|XPT|XSU|XTS|XUA|XXX|YER|ZAR|ZMW|ZWL`
	// PersonType Person Type
	PersonType = `sender|receiver`
	//Decimal ...
	Decimal = `[+,-]{0,1}[0-9]{1,16}[.]{0,1}[0-9]{0,6}`
	//IdentificationType ...
	IdentificationType = `nric|voters|drivers|passport|pan|ration|bills|license|loi|aadhaar|epfin|spass|wp|cmnd|military|medicare|drivers_id|s_license|voters_id|medicare_card|u_bill|alien_registration|birth_certificate|driving_license|employer_id|health_card|national_id|other|pan_card|resident_card|senior_citizen_id|social_security|tax_id|village_elder_id`
	//Numeric ...
	Numeric = `^[0-9]*$`
	//MobileNumber ...
	MobileNumber = `^[+]?[0-9]*$`
	// TransactionType acceptable values for transaction type
	TransactionType = `C2C|B2B|B2C|C2B`
	// BeneficiaryType acceptable values for beneficiary type
	BeneficiaryType = `consumer|corporate`
	// String regexp for Case-insensitive Alphanumeric string with some special characters: # - _ ? / > . < , @ ^ = [ ] \ : ;
	String = `^[#.0-9a-zA-Z\s,-_]+$`
	// PatternString a-z, A-Z, 0-9, /-?:().,'+.#!@&
	PatternString = `^[#.0-9a-zA-Z ,\\/\\_:+?')(@#!&-]+$`
)

// Pattern validates if the value follows a valid regex pattern
func Pattern(v interface{}, param string) error {
	val := reflect.ValueOf(v).String()

	pat := map[string]string{
		"uri":                 URI,
		"email":               Email,
		"ipv4":                IPv4,
		"md5":                 HashMD5,
		"httpmethod":          HTTPMethod,
		"calculation_mode":    CalculationMode,
		"channel_code":        DeliveryMethod,
		"csalphanum":          CSAlphaNum,
		"cialphanum":          CIAlphaNum,
		"cialphanumspace":     CIAlphaNumSpace,
		"cialpha":             CIAlpha,
		"birthday":            Birthday,
		"currency":            CurrencyCode,
		"country":             CountryCodeAlpha3,
		"person_type":         PersonType,
		"decimal":             Decimal,
		"annex_type":          AnnexTypes,
		"identification_type": IdentificationType,
		"numeric":             Numeric,
		"mobile_number":       MobileNumber,
		"transaction_type":    TransactionType,
		"beneficiary_type":    BeneficiaryType,
		"string":              String,
		"name":                PatternString,
	}

	if _, ok := pat[param]; !ok {
		return errors.New("missing pattern=" + param)
	}

	if val == "" {
		return nil
	}

	if param == "identification_type" || param == "beneficiary_type" {
		val = strings.ToLower(val)
	}

	pattern := `^(` + pat[param] + `)$`
	if param == "cialpha" || param == "string" {
		pattern = pat[param]
	}

	if param == "birthday" {
		_, err := time.Parse("2006-01-02", val)
		if err != nil {
			return errors.New("invalid pattern=" + param)
		}
	}

	if param == "name" {
		var allSupported = regexp.MustCompile(pat[param]).MatchString(val)
		if !allSupported {
			return errors.New("invalid pattern=" + param)
		}
	}

	if ok, _ := regexp.MatchString(pattern, val); !ok {
		if param == "string" {
			return errors.New("request type invalid")
		}

		return errors.New("invalid pattern=" + param)
	}

	return nil
}

// SQLValue validates if the value follows a valid sql value
func SQLValue(v interface{}, param string) error {
	val := reflect.ValueOf(v).String()

	pat := map[string]string{
		"bool":     `[0-1]`,
		"date":     `[0-9]{4}\-[0-9]{2}\-[0-9]{2}`,
		"datetime": `[0-9]{4}\-[0-9]{2}\-[0-9]{2}\ [0-9]{2}:[0-9]{2}:[0-9]{2}`,
	}

	if _, ok := pat[param]; !ok {
		return errors.New("missing sqlvalue=" + param)
	}

	if val == "" {
		return nil
	}

	if ok, _ := regexp.MatchString(`^(`+pat[param]+`)$`, val); !ok {
		return errors.New("invalid sqlvalue=" + param)
	}

	return nil
}

// Gender validate if field is male or female
func Gender(v interface{}, param string) error {
	var (
		gender string
	)

	st := reflect.ValueOf(v)
	if st.String() != "" {
		gender = strings.ToLower(st.String())
		if gender != "male" && gender != "female" {
			return errors.New(ParamGender)
		}
	}

	return nil
}

// GreaterThan ...
func GreaterThan(v interface{}, param string) error {
	st := reflect.ValueOf(v)
	if st.String() == "" {
		return nil
	}

	val, err := strconv.ParseFloat(st.String(), 64)
	if err != nil {
		return errors.New(ParamGreaterThan)
	}

	if val >= 0 {
		return nil
	}

	return errors.New(ParamGreaterThan)
}

// Required validate if field exist
func Required(v interface{}, param string) error {
	st := reflect.ValueOf(v)
	if st.String() == "" {
		return errors.New(ParamRequired)
	}

	return nil
}

// URL validate if field exist
func URL(v interface{}, param string) error {
	st := reflect.ValueOf(v)
	if st.String() == "" {
		return nil
	}

	u, err := url.Parse(st.String())
	if strings.HasPrefix(u.Host, ".") || err != nil {
		return errors.New(ParamURLInvalid)
	}
	if u.Host == "" && (u.Path != "" && !strings.Contains(u.Path, ".")) {
		return errors.New(ParamURLInvalid)
	}

	return nil
}
