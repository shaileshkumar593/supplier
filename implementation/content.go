package implementation

import (
	"context"
	"fmt"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	"swallow-supplier/config"
	customError "swallow-supplier/error"
	"swallow-supplier/mongo/domain/trip"
	reqresp "swallow-supplier/request_response/trip"
	"swallow-supplier/request_response/yanolja"
	"swallow-supplier/utils"
	"swallow-supplier/utils/constant"
	"swallow-supplier/utils/gcp/geoCoder"

	"github.com/go-kit/log/level"
)

// ProductSyncToTrip content sync for product
func (s *service) ProductSyncToTrip(ctx context.Context) (resp yanolja.Response, err error) {
	level.Info(s.logger).Log("method_name ", "ProductSyncToTrip")

	mrepo := s.mongoRepository[config.Instance().MongoDBName]
	logger := s.logger
	cf := config.Instance()

	defer func(ctx context.Context) {
		if r := recover(); r != nil {
			level.Error(logger).Log("error", "panic occurred", "panic", r, "stack", string(debug.Stack()))
			resp.Code = "500"
			err = fmt.Errorf("panic occurred: %v", r)
		}
	}(ctx)

	geocoder, err := geoCoder.NewGoogleGeoCoder(cf.GooglePlaceIdKey, cf.GooglePlaceIdBaseUrl)
	if err != nil {
		level.Error(logger).Log("error", "failed to create google geocoder", "err", err)
		return resp, err
	}

	products, err := mrepo.FetchProductsWithContentScheduleStatusFalse(ctx)
	if err != nil {
		level.Error(logger).Log("error", "request to FetchProductsUpdatedToday raised error ", err)
		resp.Code = "500"
		return resp, customError.NewError(ctx, "leisure-api-1015", fmt.Sprintf("repository error from FetchProductsUpdatedToday, %v", err), "FetchProductsUpdatedToday")
	}

	length := len(products)
	count := length

	placeIds := make([]trip.GooglePlaceIdOfProduct, 0, length)
	var productContents []trip.ProuctContent
	var packageContents []trip.PackageContent

	for _, product := range products {
		if !strings.EqualFold(product.ProductStatusCode, "IN_SALE") {
			continue
		}

		productId := strconv.FormatInt(product.ProductID, 10)

		// Category filters (dedupe)
		catSeen := map[string]struct{}{}
		var categoryFilters []reqresp.CategoryFilter
		for _, category := range product.Categories {
			for _, level1 := range category.SubCategories {
				for _, level2 := range level1.SubCategories {
					key := level2.CategoryCode + "|" + strconv.Itoa(int(level2.CategoryLevel))
					if _, ok := catSeen[key]; ok {
						continue
					}
					catSeen[key] = struct{}{}
					code, err := strconv.Atoi(level2.CategoryCode)
					if err != nil {
						level.Error(logger).Log("error", "category code conversion failed", "err", err)
						resp.Code = "500"
						return resp, err
					}
					categoryFilters = append(categoryFilters, reqresp.CategoryFilter{
						CategoryCode:  code,
						CategoryLevel: int(level2.CategoryLevel),
					})
				}
			}
		}

		tripCategories, err := mrepo.GetTripsByCategory(ctx, categoryFilters)
		if err != nil {
			level.Error(logger).Log("error", "GetTripsByCategory failed", "err", err)
			resp.Code = "500"
			return resp, fmt.Errorf("repository error: %v", err)
		}

		var categories []trip.ProductCategory
		for _, name := range tripCategories {
			categories = append(categories, trip.ProductCategory{Code: strings.ToUpper(name)})
		}
		categories = uniqCategories(categories)

		// Tags (dedupe)
		var tags []trip.ProductTags
		for _, tag := range product.ConvenienceTypeCode {
			if value, exists := constant.Tag[tag]; exists {
				tags = append(tags, trip.ProductTags{
					Id:      strconv.Itoa(value),
					TagName: tag,
				})
			} else {
				level.Info(logger).Log("info", "unrecognized tag", "tag", tag)
			}
		}
		tags = uniqTags(tags)

		var (
			pois         []trip.ProductPoi
			destinations []trip.DestinationObj
			departures   []trip.Departure
			redeemInfos  []trip.RedemptionLocation
		)

		// Facility → POI/destination/departure/redemption (dedupe by coords)
		seenCoords := map[string]struct{}{}
		for _, val := range product.ProductInfo.FacilityInfos {
			coordinateKey := fmt.Sprintf("%.6f|%.6f", val.Location.Latitude, val.Location.Longitude)
			if _, ok := seenCoords[coordinateKey]; ok {
				continue
			}
			seenCoords[coordinateKey] = struct{}{}

			var poiInfo trip.ProductPoi
			var destination trip.DestinationObj
			var departure trip.Departure
			var redeemLocation trip.RedemptionLocation

			poiInfo.GooglePlaceId, err = geocoder.GetPlaceID(ctx, val.Location.Latitude, val.Location.Longitude, logger)
			if err != nil {
				level.Error(logger).Log("error", err)
				resp.Code = "500"
				return resp, err
			}
			destination.GooglePlaceId = poiInfo.GooglePlaceId
			departure.GooglePlaceId = poiInfo.GooglePlaceId
			redeemLocation.GooglePlaceId = poiInfo.GooglePlaceId

			poiInfo.SupplierPOI.SupplierId = "YANOLJA" // remove const val
			destination.SupplierDestination.SupplierId = "YANOLJA"
			departure.SupplierDeparture.SupplierId = "YANOLJA"
			redeemLocation.SupplierLocation.SupplierId = "YANOLJA"

			poiInfo.SupplierPOI.MappingElements.Name = val.FacilityName
			destination.SupplierDestination.MappingElements.Name = val.FacilityName
			departure.SupplierDeparture.MappingElements.Name = val.FacilityName
			redeemLocation.SupplierLocation.MappingElements.Name = val.FacilityName

			poiInfo.SupplierPOI.MappingElements.AddressDetail = val.Location.Address
			destination.SupplierDestination.MappingElements.AddressDetail = val.Location.Address
			departure.SupplierDeparture.MappingElements.AddressDetail = val.Location.Address
			redeemLocation.SupplierLocation.MappingElements.AddressDetail = val.Location.Address

			poiInfo.SupplierPOI.MappingElements.Latitude = val.Location.Latitude
			destination.SupplierDestination.MappingElements.Latitude = val.Location.Latitude
			departure.SupplierDeparture.MappingElements.Latitude = val.Location.Latitude
			redeemLocation.SupplierLocation.MappingElements.Latitude = val.Location.Latitude

			poiInfo.SupplierPOI.MappingElements.Longitude = val.Location.Longitude
			destination.SupplierDestination.MappingElements.Longitude = val.Location.Longitude
			departure.SupplierDeparture.MappingElements.Longitude = val.Location.Longitude
			redeemLocation.SupplierLocation.MappingElements.Longitude = val.Location.Longitude

			placeInfo := trip.GooglePlaceIdOfProduct{
				Latitude:  val.Location.Latitude,
				Longitude: val.Location.Longitude,
				ProductID: product.ProductID,
				PlaceId:   poiInfo.GooglePlaceId,
			}

			placeIds = append(placeIds, placeInfo)
			pois = append(pois, poiInfo)
			destinations = append(destinations, destination)
			departures = append(departures, departure)
			redeemInfos = append(redeemInfos, redeemLocation)
		}

		pois = uniqPOI(pois)
		destinations = uniqDestination(destinations)
		departures = uniqDeparture(departures)
		redeemInfos = uniqRedemptionLoc(redeemInfos)

		// Ticket + HowToUse
		var ticketInfo trip.TicketInfoObj
		var redemptionInfo trip.RedemptionInfoObj
		var howToUse []string

		ticketInfo.DeliveryMethods = "PRINT"
		for _, pictogram := range product.Pictograms {
			if strings.Contains(strings.ToLower(pictogram.PictogramContent), "qr/barcode") {
				ticketInfo.DeliveryMethods = "DIGITAL"
				redemptionInfo.Description = pictogram.PictogramContent
			}
			howToUse = append(howToUse, pictogram.PictogramContent)
		}
		howToUse = uniqStrings(howToUse)

		redemptionInfo.RedemptionType = "Direct_Entry"
		redemptionInfo.RedemptionLocation = redeemInfos

		// Service languages
		var serviceLanguages []trip.ServiceLanguageObj
		serviceLanguages = append(serviceLanguages, trip.ServiceLanguageObj{LanguageCode: "en"})

		// Gallery from Trip image ids
		ImageIds, err := mrepo.FetchTripImageIdsByProductID(ctx, product.ProductID)
		if err != nil {
			level.Error(logger).Log("error", "request to FetchTripImageIdsByProductID raised error ", err)
			resp.Code = "500"
			return resp, customError.NewError(ctx, "leisure-api-1015", fmt.Sprintf("repository error from FetchTripImageIdsByProductID, %v", err), "FetchTripImageIdsByProductID")
		}
		var galleryId = make([]trip.ProductImageObj, len(ImageIds))
		for index, id := range ImageIds {
			galleryId[index].TripImageId = id
		}
		galleryId = uniqImages(galleryId)

		// Highlights + Description
		highlights := uniqStrings([]string{"", ""})

		var describes []string
		describes = append(describes,
			product.ProductInfo.ProductBasicInfo,
			product.ProductInfo.VoucherUsageInfo,
			product.ProductInfo.ProductUsageInfo,
			product.ProductInfo.RefundInfo,
		)
		detailInfo := strings.Join(describes, "\n")

		// Booking, options, cancellation, ticket types
		var bookingInfo trip.BookingSettingsObj
		var productOptions []trip.ProductOptionObj
		var cancellation trip.CancellationPolicyObj
		var ticketTypeAry []trip.TicketTypeObj

		var option []trip.PackageOption
		var extraInfo []string

		// per-variant TicketType meta for Unit construction
		type ttMeta struct {
			ticketCode string
			customCode string
		}
		variantMeta := map[int64]ttMeta{}

		for i, optionGroup := range product.ProductOptionGroups {
			extraInfo = append(extraInfo, optionGroup.ProductOptionGroupDescription)

			if optionGroup.IsSchedule {
				bookingInfo.BookingType.DateType = "DATE_REQUIRED"
			} else {
				bookingInfo.BookingType.DateType = "DATE_NOT_REQUIRED"
			}

			var productOption trip.ProductOptionObj
			var pkgOption trip.PackageOption

			if optionGroup.IsRound {
				productOption.OptionCode = "Time_Slot"
				pkgOption.OptionCode = "Time_Slot"
				pkgOption.ValueCode = "Time_Slot" + strconv.Itoa(i)
				for _, po := range optionGroup.ProductOptions {
					for _, item := range po.ProductOptionItems {
						if item.Rounds != nil {
							pkgOption.ValueName = strings.Join(item.Rounds, "-")
						}
					}
				}
			} else {
				productOption.OptionCode = "Option"
				pkgOption.OptionCode = "Option"
				pkgOption.ValueCode = "Option" + strconv.Itoa(i)
			}

			// cancellation and TicketType creation per-variant
			if cancellation.Type == "" {
				for _, v := range optionGroup.Variants {
					if v.IsRefundableAfterExpiration && (v.RefundApprovalTypeCode == "DIRECT" || v.RefundApprovalTypeCode == "ADMIN") {
						cancellation.Type = "Free_Cancel"
					} else {
						cancellation.Type = "Non_Cancellable"
					}

					nameLC := strings.ToLower(v.VariantName)
					ttCode := "Customized"
					switch {
					case strings.Contains(nameLC, "adult"):
						ttCode = "Adult"
					case strings.Contains(nameLC, "child"):
						ttCode = "Child"
					case strings.Contains(nameLC, "senior"):
						ttCode = "Senior"
					case strings.Contains(nameLC, "youth"):
						ttCode = "Youth"
					case strings.Contains(nameLC, "infant"):
						ttCode = "Infant"
					case strings.Contains(nameLC, "student"):
						ttCode = "Student"
					case strings.Contains(nameLC, "traveler"):
						ttCode = "Traveler"
					}

					vCustom := utils.GenerateRandAlphaNumeric(8) // per-variant custom code

					tt := trip.TicketTypeObj{
						Code:        ttCode,
						CustomCode:  vCustom,
						CustomName:  v.VariantName,
						Description: product.ProductInfo.ProductUsageInfo,
					}
					// optional age ranges
					switch ttCode {
					case "Adult":
						tt.Restrictions = trip.CrowdLimitObj{MinAge: 18, MaxAge: 60}
					case "Child":
						tt.Restrictions = trip.CrowdLimitObj{MinAge: 5, MaxAge: 12}
					case "Senior":
						tt.Restrictions = trip.CrowdLimitObj{MinAge: 60, MaxAge: 75}
					case "Youth":
						tt.Restrictions = trip.CrowdLimitObj{MinAge: 12, MaxAge: 18}
					case "Infant":
						tt.Restrictions = trip.CrowdLimitObj{MinAge: 2, MaxAge: 5}
					case "Student":
						tt.Restrictions = trip.CrowdLimitObj{MinAge: 8, MaxAge: 30}
					case "Traveler":
						tt.Restrictions = trip.CrowdLimitObj{MinAge: 5, MaxAge: 90}
					}

					ticketTypeAry = append(ticketTypeAry, tt)
					variantMeta[v.VariantID] = ttMeta{ticketCode: ttCode, customCode: vCustom}
				}
			}

			productOptions = append(productOptions, productOption)
			option = append(option, pkgOption)
		}

		// Booking window
		bookingInfo.PaymentConfirmationTime = 5
		if FindDaysBetweenDate(product.SalePeriod.StartDateTime, product.SalePeriod.EndDateTime) > 1 {
			bookingInfo.BookingType.DateLimit.DateLimitType = "Customized"

			t, err := time.Parse(time.RFC3339, product.SalePeriod.StartDateTime)
			if err != nil {
				level.Error(logger).Log("error ", "parsing issue in start date format")
				resp.Code = "500"
				return resp, err
			}
			bookingInfo.BookingType.DateLimit.CustomizedDateRange.FromDate = t.Format("2006-01-02")

			t, err = time.Parse(time.RFC3339, product.SalePeriod.EndDateTime)
			if err != nil {
				level.Error(logger).Log("error ", "parsing issue in end date format")
				resp.Code = "500"
				return resp, err
			}
			bookingInfo.BookingType.DateLimit.CustomizedDateRange.ToDate = t.Format("2006-01-02")
		} else {
			bookingInfo.BookingType.DateLimit.DateLimitType = "Single_date"
		}

		additionalInfo := strings.Join(extraInfo, "\n")

		guestInfo := trip.GuestInformationObj{
			Type: "PER_ORDER",
			Code: []string{"GUEST_NAME", "BIRTH_DATE"}, // struct declation
		}

		cancellation.ConfirmationTime = 24
		cancellation.RateList = nil

		// package data
		var bookingCutOffTime trip.BookingCutOffTime
		bookingCutOffTime.DayBeforeVisitDate = "1"
		bookingCutOffTime.Time = "0:00"

		var optionsList []trip.PackageOptionList
		var pkgOptionsList trip.PackageOptionList
		var unitAry []trip.UnitObj

		pkgOptionsList.OptionStatus = "active"
		pkgOptionsList.OptionDescription = product.ProductBriefIntroduction
		pkgOptionsList.BookingCutOffTime = bookingCutOffTime

		plus, err := mrepo.FetchPluHashesByProductID(ctx, product.ProductID)
		if err != nil {
			level.Error(logger).Log("error", "request to FetchAllPluHashes raised error ", err)
			resp.Code = "500"
			return resp, customError.NewError(ctx, "leisure-api-1015", fmt.Sprintf("repository error from FetchAllPluHashes, %v", err), "FetchAllPluHashes")
		}

		seenPLU := make(map[string]struct{}, len(plus))
		for key, val := range plus {
			if key == "" {
				continue
			}
			if _, dup := seenPLU[key]; dup {
				continue
			}
			seenPLU[key] = struct{}{}

			ttCode := "Customized"
			uCustom := "" // will fill from variantMeta if found

			if val != "" {
				detail := strings.Split(val, "|")
				if len(detail) >= 3 {
					if variantID, err := strconv.ParseInt(detail[2], 10, 64); err == nil {
						if meta, ok := variantMeta[variantID]; ok {
							ttCode = meta.ticketCode
							uCustom = meta.customCode
						}
					}
				}
			}
			if uCustom == "" {
				uCustom = utils.GenerateRandAlphaNumeric(8)
			}

			unitAry = append(unitAry, trip.UnitObj{
				PLU:            key,
				Reference:      product.ProductName,
				TicketTypeCode: ttCode,
				CustomCode:     uCustom,
				Restrictions: trip.Restrictions{
					MinUnits: "1", MaxUnits: "30", UnitPax: "1", CompanionRequired: "Not_Required",
				},
				Currency: trip.CurrencyObj{
					NetPriceCurrency: "KRW", RetailPriceCurrency: "KRW",
				},
			})
		}

		pkgOptionsList.Unit = unitAry
		optionsList = append(optionsList, pkgOptionsList)

		// Final dedupe before write
		ticketTypeAry = uniqTicketTypes(ticketTypeAry)
		highlights = uniqStrings(highlights)

		// Product content
		productContent := trip.ProuctContent{
			SupplierProductId:  productId,
			SupplierName:       product.SupplierName,
			Reference:          product.ProductTypeCode,
			ContractId:         200952,
			PrimaryLanguage:    "ko-KR",
			Status:             "active",
			Category:           categories,
			Tags:               tags,
			Title:              product.ProductName,
			Poi:                pois,
			Destination:        destinations,
			Departure:          departures,
			TicketInfo:         ticketInfo,
			RedemptionInfo:     redemptionInfo,
			ServiceLanguage:    serviceLanguages,
			Gallery:            galleryId,
			Highlight:          highlights,
			Description:        detailInfo,
			HowToUse:           howToUse,
			AdditionalInfo:     additionalInfo,
			GuestInformation:   guestInfo,
			BookingSettings:    bookingInfo,
			CancellationPolicy: cancellation,
			TicketType:         ticketTypeAry,
			Option:             productOptions,
			MetaData:           additionalInfo,
			SyncStatus:         "NotSync",
		}

		// Package content
		packageContent := trip.PackageContent{
			SupplierProductId: productId,
			SupplierName:      product.SupplierName,
			OptionList:        optionsList,
			SyncStatus:        "NotSync",
		}

		productContents = append(productContents, productContent)
		packageContents = append(packageContents, packageContent)

		count = count - 1
		if count == 0 {
			break
		}
	}

	if len(productContents) == 0 {
		level.Error(logger).Log("error", "no product content available for insertion")
		resp.Code = "500"
		resp.Body = fmt.Sprintln("no more product content for insertion")
		return resp, fmt.Errorf("no product content provided for upsert")
	}

	if err = mrepo.BulkUpsertProductContent(ctx, productContents); err != nil {
		level.Error(logger).Log("error", "request to BulkUpsertProductContent raised error ", err)
		resp.Code = "500"
		return resp, customError.NewError(ctx, "leisure-api-1015", fmt.Sprintf("repository error from BulkUpsertProductContent, %v", err), "BulkUpsertProductContent")
	}

	// Trip Sync
	if len(packageContents) == 0 { // fixed check
		level.Error(logger).Log("error", "no package content available for insertion")
		resp.Code = "500"
		resp.Body = fmt.Sprintln("no more package content for insertion")
		return resp, fmt.Errorf("no package content provided for upsert")
	}
	if err = mrepo.BulkUpsertPackageContent(ctx, packageContents); err != nil {
		level.Error(logger).Log("error", "request to BulkUpsertPackageContent raised error ", err)
		resp.Code = "500"
		return resp, customError.NewError(ctx, "leisure-api-1015", fmt.Sprintf("repository error from BulkUpsertPackageContent, %v", err), "BulkUpsertPackageContent")
	}

	placeIds = uniqPlaceIds(placeIds)
	if len(placeIds) == 0 {
		level.Error(logger).Log("error", "no PlaceId  available for insertion")
		resp.Code = "500"
		resp.Body = fmt.Sprintln("no more PlaceId for insertion")
		return resp, fmt.Errorf("no PlaceId provided for upsert")
	}
	if err = mrepo.BulkUpsertGooglePlaceIdOfProduct(ctx, placeIds); err != nil {
		level.Error(logger).Log("error", "request to BulkUpsertGooglePlaceIdOfProduct raised error ", err)
		resp.Code = "500"
		return resp, customError.NewError(ctx, "leisure-api-1015", fmt.Sprintf("repository error from BulkUpsertGooglePlaceIdOfProduct, %v", err), "BulkUpsertGooglePlaceIdOfProduct")
	}

	resp.Body = "Successfully inserted product and package"
	resp.Code = "200"
	level.Info(logger).Log("Info", "Successfully synced")
	return resp, nil
}

// FindDaysBetweenDate find number of days between two dates
func FindDaysBetweenDate(startDateTime string, endDateTime string) int {
	// Parse timestamps
	layout := time.RFC3339

	startTime, err := time.Parse(layout, startDateTime)
	if err != nil {
		fmt.Println("Error parsing startDateTime:", err)
		return -1
	}

	endTime, err := time.Parse(layout, endDateTime)
	if err != nil {
		fmt.Println("Error parsing endDateTime:", err)
		return -1
	}

	// Calculate duration in days
	duration := endTime.Sub(startTime)
	daysDifference := int(duration.Hours() / 24)

	fmt.Printf("Difference in days: %d days\n", daysDifference)

	return daysDifference
}

// --- helpers (put near this file) ---
func uniqStrings(in []string) []string {
	m := make(map[string]struct{}, len(in))
	out := make([]string, 0, len(in))
	for _, s := range in {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		if _, ok := m[s]; ok {
			continue
		}
		m[s] = struct{}{}
		out = append(out, s)
	}
	return out
}

func uniqCategories(in []trip.ProductCategory) []trip.ProductCategory {
	m := make(map[string]struct{}, len(in))
	out := make([]trip.ProductCategory, 0, len(in))
	for _, c := range in {
		k := strings.ToUpper(c.Code)
		if _, ok := m[k]; ok {
			continue
		}
		m[k] = struct{}{}
		out = append(out, trip.ProductCategory{Code: k})
	}
	return out
}

func uniqTags(in []trip.ProductTags) []trip.ProductTags {
	m := make(map[string]struct{}, len(in))
	out := make([]trip.ProductTags, 0, len(in))
	for _, t := range in {
		k := strings.ToUpper(t.TagName)
		if _, ok := m[k]; ok {
			continue
		}
		m[k] = struct{}{}
		out = append(out, t)
	}
	return out
}

func uniqImages(in []trip.ProductImageObj) []trip.ProductImageObj {
	m := make(map[string]struct{}, len(in))
	out := make([]trip.ProductImageObj, 0, len(in))
	for _, im := range in {
		if im.TripImageId == "" {
			continue
		}
		if _, ok := m[im.TripImageId]; ok {
			continue
		}
		m[im.TripImageId] = struct{}{}
		out = append(out, im)
	}
	return out
}

func uniqPOI(in []trip.ProductPoi) []trip.ProductPoi {
	m := make(map[string]struct{}, len(in))
	out := make([]trip.ProductPoi, 0, len(in))
	for _, p := range in {
		k := p.GooglePlaceId
		if k == "" {
			continue
		}
		if _, ok := m[k]; ok {
			continue
		}
		m[k] = struct{}{}
		out = append(out, p)
	}
	return out
}

func uniqDestination(in []trip.DestinationObj) []trip.DestinationObj {
	m := make(map[string]struct{}, len(in))
	out := make([]trip.DestinationObj, 0, len(in))
	for _, d := range in {
		k := d.GooglePlaceId
		if k == "" {
			continue
		}
		if _, ok := m[k]; ok {
			continue
		}
		m[k] = struct{}{}
		out = append(out, d)
	}
	return out
}

func uniqDeparture(in []trip.Departure) []trip.Departure {
	m := make(map[string]struct{}, len(in))
	out := make([]trip.Departure, 0, len(in))
	for _, d := range in {
		k := d.GooglePlaceId
		if k == "" {
			continue
		}
		if _, ok := m[k]; ok {
			continue
		}
		m[k] = struct{}{}
		out = append(out, d)
	}
	return out
}

func uniqRedemptionLoc(in []trip.RedemptionLocation) []trip.RedemptionLocation {
	m := make(map[string]struct{}, len(in))
	out := make([]trip.RedemptionLocation, 0, len(in))
	for _, r := range in {
		k := r.GooglePlaceId
		if k == "" {
			continue
		}
		if _, ok := m[k]; ok {
			continue
		}
		m[k] = struct{}{}
		out = append(out, r)
	}
	return out
}

func uniqTicketTypes(in []trip.TicketTypeObj) []trip.TicketTypeObj {
	m := make(map[string]struct{}, len(in))
	out := make([]trip.TicketTypeObj, 0, len(in))
	for _, tt := range in {
		k := strings.ToUpper(tt.Code) + "|" + strings.ToUpper(tt.CustomName)
		if _, ok := m[k]; ok {
			continue
		}
		m[k] = struct{}{}
		out = append(out, tt)
	}
	return out
}

func uniqPlaceIds(in []trip.GooglePlaceIdOfProduct) []trip.GooglePlaceIdOfProduct {
	m := make(map[string]struct{}, len(in))
	out := make([]trip.GooglePlaceIdOfProduct, 0, len(in))
	for _, p := range in {
		k := fmt.Sprintf("%d|%s", p.ProductID, p.PlaceId)
		if p.PlaceId == "" {
			continue
		}
		if _, ok := m[k]; ok {
			continue
		}
		m[k] = struct{}{}
		out = append(out, p)
	}
	return out
}
