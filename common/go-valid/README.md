# valid
--
    import "bitbucket.org/matchmove/go-valid"


## Usage

```go
const (
	// URI regexp pattern for URI pattern
	URI = `(\/([a-zA-Z0-9_]+|\{[a-zA-Z0-9_]+\}))+`
	// Email regexp pattern for RFC 5322 (email) electronic mail address
	Email = `(?:[a-z0-9!#$%&'*+/=?^_` + "`" +
		`{|}~-]+(?:\.[a-z0-9!#$%&'*+/=?^_` + "`" +
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
	IdentificationType = `nric|voters|drivers|passport|pan|ration|bills|license|loi|aadhaar|epfin|spass|wp|cmnd|military|medicare`
)
```

#### func  Gender

```go
func Gender(v interface{}, param string) error
```
Gender validate if field is male or female

#### func  GreaterThan

```go
func GreaterThan(v interface{}, param string) error
```
GreaterThan ...

#### func  Pattern

```go
func Pattern(v interface{}, param string) error
```
Pattern validates if the value follows a valid regex pattern

#### func  Required

```go
func Required(v interface{}, param string) error
```
Required validate if field exist

#### func  SQLValue

```go
func SQLValue(v interface{}, param string) error
```
SQLValue validates if the value follows a valid sql value

#### func  URL

```go
func URL(v interface{}, param string) error
```
URL validate if field exist
