package constant

var PRODUCTVARIANTSTATUSCODE = []string{"WAITING_FOR_SALE", "IN_SALE", "SOLD_OUT", "END_OF_SALE"}
var RECONCILESTATUSUNKNOWN string = "UNKNOWN" // this is for GGT use not for yanolja
var RECONCILESTATUSCREATED string = "CREATED"
var RECONCILESTATUSUSED string = "USED"
var RECONCILESTATUSCANCELED string = "CANCELED"
var RECONCILESTATUSRESTORED string = "RESTORED"
var FIXVALIDITYPERIODTYPECODE string = "FIX"
var BUYVALIDITYPERIODTYPECODE string = "BUY"
var CONSTANTCONSUME string = "CONSUME"
var CONSTANTRESTORE string = "RESTORE"
var ORDERVARIANTUNKNOWNSTATUS string = "UNKNOWN"     // Unconfirmed
var ORDERVARIANTUSEDSTATUS string = "USED"           // used
var ORDERVARIANTNOTUSEDSTATUS string = "NOT_USED"    // unused
var ORDERVARIANTCANCELEDSTATUS string = "CANCELED"   // order variant canceled
var ORDERVARIANTCANCELINGSTATUS string = "CANCELING" // Requesting cancellation

//var POLLINTERVAL = 100 * time.Second                 // Poll every 10 seconds

//var API_KEY string = "xOsUvg9SEjlp-W7e0fY1H9cWh-8POGRIE8BqdiQANCo="

const MAXORDERQUANTITYONETIME = 30

// VoucherSendType defines an enumeration for the voucher sending type
type VoucherSendType int

const (
	// Do not send (default)
	VoucherSendTypeDoNotSend VoucherSendType = iota

	// Send via email
	VoucherSendTypeSendViaEmail

	// Send via Kakao Talk Biz Messaging
	VoucherSendTypeSendViaKakaoTalkBizMsg

	// Send via email & Kakao Talk Biz Messaging
	VoucherSendTypeSendViaBoth
)

// String converts the VoucherSendType enum to a human-readable string
func (v VoucherSendType) String() string {
	switch v {
	case VoucherSendTypeDoNotSend:
		return "Do not Send"
	case VoucherSendTypeSendViaEmail:
		return "Send via Email"
	case VoucherSendTypeSendViaKakaoTalkBizMsg:
		return "Send via Kakao Talk Biz Messaging"
	case VoucherSendTypeSendViaBoth:
		return "Send via Email & Kakao Talk Biz Messaging"
	default:
		return "Unknown Voucher Send Type"
	}
}

var Tag = map[string]int{
	"TICKET_BOX":         1,
	"FREE_WIFI":          2,
	"LOCKER":             3,
	"PET":                4,
	"CHILD":              5,
	"SAFETY":             6,
	"MEDICAL_OFFICE":     7,
	"SMOKING_ROOM":       8,
	"PARKING_LOT":        9,
	"FREE_SHUTTLE":       10,
	"FREE_PICK_UP":       11,
	"FNB":                12,
	"ALLOW_OUTSIDE_FOOD": 13,
	"FITTING_ROOM":       14,
	"DEHYDRATOR":         15,
	"SHOWER":             16,
	"FREE_LESSON":        17,
	"DOCENT":             18,
	"AUDIO_GUIDE":        19,
}

const SUPPLIERYANOLJA = "Yanolja"
const SUPPLIERTRAVOLUTION = "Travolution"

const YANOLJAGGTTRIP = "Yanolja-GGT-Trip"

const REFUNDDIRECT = "DIRECT"
const REFUNDADMIN = "ADMIN"

const TRIPPREORDERREQUEST = "PreOrder"
const TRIPPAYMENTREQUEST = "Payment"

const TRIPFULLORDERCANCELREQUEST = "FullCancel"

const VALIDSTAUS = "VALID"
const EXPIREDSTAUS = "EXPIRED"

var MARKUPPERCENTAGE float32 = 3.0

const PERCENTAGE = "PERCENTAGE"
const FLATVALUE = "FLATVALUE"
const DEFAULTCURRENCY = "KRW"
const ORDERCOMPLETE = "CONFIRMED"
const ORDERDONE = "DONE"
const ORDERPREPARE = "PREPARE"
const ORDERNOTCOMPLETE = "NOT_CONFIRMED"
const STATUSNOTSYNC = "NotSync"

// Travolution
const ORDERAVAILABLE = "Available"
const ORDERCANCELREQUEST = "Cancel Request"
const ORDERCANCELED = "Canceled"
const ORDERAPPROVED = "Approved" // used
const ORDEREXPIRED = "Expired"

const BOOKINGPENDING = "Pending"
const BOOKINGAPPROVED = "Approved" // confirm
const BOOKINGREJECTED = "Rejected"

type Margine struct {
	ProductId   int64
	MargineType string
	Value       float32
}

// EverLand record for margine  dev : {10014508, 10012712} prod :{10240621,10240654}
var MARGINEDETAIL = []Margine{
	{
		ProductId:   10240621,
		MargineType: "FLATVALUE",
		Value:       1900,
	}, {
		ProductId:   10240654,
		MargineType: "FLATVALUE",
		Value:       1900,
	}, {
		ProductId:   10248417,
		MargineType: "PERCENTAGE",
		Value:       1.6,
	}, {
		ProductId:   10012383,
		MargineType: "PERCENTAGE",
		Value:       1.6,
	}, {
		ProductId:   10014508,
		MargineType: "FLATVALUE",
		Value:       1900,
	}, {
		ProductId:   10012712,
		MargineType: "FLATVALUE",
		Value:       1900,
	},
}

var ProductForGlobaltix = [...]int64{
	10017468,
	10014508,
	10012719,
	10012717,
	10012712,
	10012289,
	10011435,
}

// VoucherType enum for travolution
const (
	VoucherTypeOrder = 1 // One code per order
	VoucherTypeUnit  = 2 // One code per unit
	VoucherTypeFile  = 3 // File
)

// OrderStatus enum  for travolution
const (
	OrderStatusAvailable     = "AV"
	OrderStatusApproved      = "AP"
	OrderStatusCancelRequest = "CR"
	OrderStatusCanceled      = "CL"
	OrderStatusExpired       = "EP"
	OrderStatusRejected      = "RJ"
)

// ProductType enum
const (
	ProductTypeTicket  = "TK"
	ProductTypeBooking = "BK"
	ProductTypePass    = "PAS"
	ProductTypePackage = "PKG"
)

// BookingStatus enum
const (
	BookingStatusPending  = "PD"
	BookingStatusApproved = "AC"
	BookingStatusRejected = "RJ"
)

// CodeType enum
const (
	PIN        = "PIN"
	Barcode    = "Barcode"
	QR         = "2D"
	ImageOrPdf = "File"
)
