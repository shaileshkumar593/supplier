package cronjob

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"swallow-supplier/config"
	svc "swallow-supplier/iface"
	"swallow-supplier/mongo/domain/trip"
	reqresp "swallow-supplier/request_response/trip"
	"swallow-supplier/utils"
	"swallow-supplier/utils/constant"
	"swallow-supplier/utils/gcp/geoCoder"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

// ProductSyncToTrip content sync for product
func ProductSyncToTrip(ctx context.Context, logger log.Logger, mrepo svc.MongoRepository) (err error) {

	fmt.Println("1")
	cf := config.Instance()

	// Defer when panic
	defer func(context.Context) {
		r := recover()
		if r == nil {
			return
		}

		level.Error(logger).Log("error", "processing request went into panic mode", "panic", r)
		err = fmt.Errorf("panic occurred: %v", r)

	}(ctx)

	geocoder, err := geoCoder.NewGoogleGeoCoder(cf.GooglePlaceIdKey, cf.GooglePlaceIdBaseUrl)
	if err != nil {
		level.Error(logger).Log("error", "failed to create google geocoder", "err", err)
		return err
	}

	fmt.Println("2")

	products, err := mrepo.FetchProductsWithContentScheduleStatusFalse(ctx)
	if err != nil {
		return err
	}
	fmt.Println("3   FetchProductsUpdatedToday")

	length := len(products)
	count := length

	placeIds := make([]trip.GooglePlaceIdOfProduct, length)
	productContents := make([]trip.ProuctContent, length)
	packageContents := make([]trip.PackageContent, length)
	var productId string

	fmt.Println("000000000000000000000000000000000 start 00000000000000000000000000000000000")
	for _, product := range products {
		if strings.EqualFold(product.ProductStatusCode, "IN_SALE") {
			fmt.Println("4")

			productId = strconv.FormatInt(product.ProductID, 10)
			fmt.Println("#############################  productId, count ###################### ", productId, count)

			categoryFilter := make([]reqresp.CategoryFilter, 0, length)
			fmt.Println("5 categoryFilter")

			for _, category := range product.Categories {
				for _, level1 := range category.SubCategories {
					for _, level2 := range level1.SubCategories {
						code, err := strconv.Atoi(level2.CategoryCode)
						if err != nil {
							fmt.Println("Conversion failed:", err) // Prints an error
							return err
						}
						categoryCodeAndLevel := reqresp.CategoryFilter{
							CategoryCode:  code,
							CategoryLevel: int(level2.CategoryLevel),
						}

						categoryFilter = append(categoryFilter, categoryCodeAndLevel)
						fmt.Println("6  ")

					}
				}
			}
			fmt.Println("6  categoryFilter ", categoryFilter)

			fmt.Println("7    tripCategory")

			tripCategory, err := mrepo.GetTripsByCategory(ctx, categoryFilter)
			if err != nil {
				return err
			}
			fmt.Println("8   category")

			var category = make([]trip.ProductCategory, 0, length)
			for index, categoryName := range tripCategory {
				category[index].Code = categoryName
			}
			fmt.Println("9   tags")

			tags := make([]trip.ProductTags, 0, length)
			for _, tagname := range product.ConvenienceTypeCode {
				value, exists := constant.Tag[tagname]
				if !exists {
					fmt.Printf("Key '%s' do not exists with value: %d\n", tagname, value)
					return fmt.Errorf("Key '%s' do not exists with value: %d\n", tagname, value)
				}
				tagVal := trip.ProductTags{
					Id:      strconv.Itoa(value),
					TagName: tagname,
				}

				tags = append(tags, tagVal)
			}
			fmt.Println("10   tags = append")

			pois := make([]trip.ProductPoi, 0, length)
			destinations := make([]trip.DestinationObj, 0, length)
			departures := make([]trip.Departure, 0, length)
			redeemInfos := make([]trip.RedemptionLocation, 0, length)

			var destination trip.DestinationObj
			var departure trip.Departure
			var redeemLocation trip.RedemptionLocation

			fmt.Println("11  redeemLocation")

			for _, val := range product.ProductInfo.FacilityInfos {
				var poiInfo trip.ProductPoi
				poiInfo.GooglePlaceId, err = geocoder.GetPlaceID(ctx, val.Location.Latitude, val.Location.Longitude, logger)
				if err != nil {
					return err
				}
				destination.GooglePlaceId = poiInfo.GooglePlaceId
				departure.GooglePlaceId = poiInfo.GooglePlaceId
				redeemLocation.GooglePlaceId = poiInfo.GooglePlaceId

				poiInfo.SupplierPOI.SupplierId = "YANOLJA"
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

			fmt.Println("12   ticketInfo")

			var ticketInfo trip.TicketInfoObj
			var redemptionInfo trip.RedemptionInfoObj
			howToUse := make([]string, 0, length)

			for _, pictogram := range product.Pictograms {
				if strings.Contains(pictogram.PictogramContent, "QR/barcode") {
					ticketInfo.DeliveryMethods = "DIGITAL"
					redemptionInfo.Description = pictogram.PictogramContent

				}
				howToUse = append(howToUse, pictogram.PictogramContent)
			}
			ticketInfo.DeliveryMethods = "PRINT"

			redemptionInfo.RedemptionType = "Direct_Entry"
			redemptionInfo.RedemptionLocation = redeemInfos

			servicelanguage := make([]trip.ServiceLanguageObj, 0, length)
			var slang trip.ServiceLanguageObj
			slang.LanguageCode = "ko-KR"
			servicelanguage = append(servicelanguage, slang)

			fmt.Println("13  servicelanguage")

			tripImageId, err := mrepo.FetchTripImageIdsByProductID(ctx, product.ProductID)
			if err != nil {
				return err
			}

			var galleryId = make([]trip.ProductImageObj, 0, length)

			for index, id := range tripImageId {
				galleryId[index].TripImageId = id
			}
			fmt.Println("14   ProductImageObj")

			higlights := make([]string, length)
			higlights = append(higlights, product.ProductBriefIntroduction, product.ProductInfo.ProductBasicInfo, product.ProductInfo.ProductUsageInfo)

			desribes := make([]string, length)
			desribes = append(desribes, higlights...)
			desribes = append(desribes, product.ProductInfo.RefundInfo, product.ProductInfo.ServiceCenterInfo, product.ProductInfo.VoucherUsageInfo)

			detailInfo := strings.Join(desribes, "\n")

			fmt.Println("15   desribes")

			var bookingInfo trip.BookingSettingsObj
			var optnAry = make([]trip.ProductOptionObj, 0, length)
			var optn trip.ProductOptionObj
			var cancellation trip.CancellationPolicyObj

			var option = make([]trip.PackageOption, 0, length)
			var pkgOption trip.PackageOption

			fmt.Println("16   PackageOption")

			extraInfo := make([]string, 0, length)
			for i, optiongroup := range product.ProductOptionGroups {
				extraInfo = append(extraInfo, optiongroup.ProductOptionGroupDescription)
				if optiongroup.IsSchedule == true {
					bookingInfo.BookingType.DateType = "DATE_REQUIRED"
				} else {
					bookingInfo.BookingType.DateType = "DATE_NOT_REQUIRED"
				}
				if optiongroup.IsRound {
					optn.OptionCode = "Time_Slot"
					pkgOption.OptionCode = "Time_Slot"
					pkgOption.ValueCode = "Time_Slot" + strconv.Itoa(i)
					for _, productOptn := range optiongroup.ProductOptions {
						for _, prproductOptnItem := range productOptn.ProductOptionItems {
							if prproductOptnItem.Rounds != nil {
								pkgOption.ValueName = strings.Join(prproductOptnItem.Rounds, "-")
							}
						}
					}

				} else {
					optn.OptionCode = "Option"
					pkgOption.OptionCode = "Option"
					pkgOption.ValueCode = "Option" + strconv.Itoa(i)
				}
				if cancellation.Type == "" {
					for _, variant := range optiongroup.Variants {
						if variant.IsRefundableAfterExpiration == true && (variant.RefundApprovalTypeCode == "DIRECT" || variant.RefundApprovalTypeCode == "ADMIN") {
							cancellation.Type = "Free_Cancel"
						} else {
							cancellation.Type = "Non_Cancellable"
						}
					}
				}

				optnAry = append(optnAry, optn)
				option = append(option, pkgOption)
			}
			fmt.Println("17    optnAry")

			bookingInfo.PaymentConfirmationTime = 5
			if FindDaysBetweenDate(product.SalePeriod.StartDateTime, product.SalePeriod.EndDateTime) > 1 {
				bookingInfo.BookingType.DateLimit.DateLimitType = "Customized"
				bookingInfo.BookingType.DateLimit.CustomizedDateRange.FromDate = product.SalePeriod.StartDateTime
				bookingInfo.BookingType.DateLimit.CustomizedDateRange.ToDate = product.SalePeriod.EndDateTime
			} else {
				bookingInfo.BookingType.DateLimit.DateLimitType = "Single_date"
			}
			fmt.Println("18   additionalInfo")

			additionalInfo := strings.Join(extraInfo, "\n")

			var guestInfo trip.GuestInformationObj
			guestInfo.Type = "PER_ORDER"
			guestInfo.Code = append(guestInfo.Code, "GUEST_NAME", "BIRTH_DATE")

			cancellation.ConfirmationTime = 24
			cancellation.RateList = nil
			customcode := utils.GenerateRandAlphaNumeric(8)

			ticketType := []trip.TicketTypeObj{
				{
					Code:        "Customized",
					CustomCode:  customcode,
					Description: product.ProductInfo.ProductUsageInfo,
				},
			}

			// package data
			var bookingCutOffTime trip.BookingCutOffTime
			bookingCutOffTime.DayBeforeVisitDate = "1"
			bookingCutOffTime.Time = "0:00"

			var optnList = make([]trip.PackageOptionList, 0, length)
			var pkgOptnList trip.PackageOptionList
			var unit trip.UnitObj
			var unitAry = make([]trip.UnitObj, 0, length)

			pkgOptnList.OptionStatus = "active"
			pkgOptnList.OptionDescription = product.ProductBriefIntroduction
			pkgOptnList.BookingCutOffTime = bookingCutOffTime

			pkgOptnList.BookingQuestions = nil

			plus, err := mrepo.FetchAllPluHashes(ctx)
			if err != nil {
				return err
			}
			fmt.Println("19     hashmap")

			for _, hashmap := range plus {
				for uid := range hashmap {
					unit.PLU = strconv.Itoa(uid)
					unit.Reference = product.ProductName
					unit.TicketTypeCode = "Customized"
					unit.CustomCode = customcode
					unit.Restrictions.MinUnits = "1"
					unit.Restrictions.MaxUnits = "30"
					unit.Restrictions.UnitPax = "1"
					unit.Restrictions.CompanionRequired = "Not_Required"

					unit.Currency.NetPriceCurrency = "KRW"
					unit.Currency.RetailPriceCurrency = "KRW"

				}
				unitAry = append(unitAry, unit)
			}
			pkgOptnList.Unit = unitAry
			optnList = append(optnList, pkgOptnList)
			fmt.Println("20    productContent")

			productContent := trip.ProuctContent{
				SupplierProductId:  productId,
				Reference:          product.ProductTypeCode,
				ContractId:         200952,
				PrimaryLanguage:    "ko-KR",
				Status:             "active",
				Category:           category,
				Tags:               tags,
				Title:              product.ProductName,
				Poi:                pois,
				Destination:        destinations,
				Departure:          departures,
				TicketInfo:         ticketInfo,
				RedemptionInfo:     redemptionInfo,
				ServiceLanguage:    servicelanguage,
				Gallery:            galleryId,
				Highlight:          higlights,
				Description:        detailInfo,
				HowToUse:           howToUse,
				AdditionalInfo:     additionalInfo,
				GuestInformation:   guestInfo,
				BookingSettings:    bookingInfo,
				CancellationPolicy: cancellation,
				TicketType:         ticketType,
				Option:             optnAry,
				MetaData:           additionalInfo,
				SyncStatus:         "NotSync",
			}
			fmt.Println("21   packageContent")

			packageContent := trip.PackageContent{
				SupplierProductId: productId,
				OptionList:        optnList,
				SyncStatus:        "NotSync",
			}

			productContents = append(productContents, productContent)
			packageContents = append(packageContents, packageContent)
			fmt.Println("22  final append")

		}
		fmt.Println("*********************count******************* : ", count)
		count = count - 1
		if count == 0 {
			fmt.Println("============================= content is created=================================")
			break
		}
	}
	fmt.Println("000000000000000000000000000000000 end 00000000000000000000000000000000000")

	fmt.Println("*****************BulkUpsertProductContent************************")
	err = mrepo.BulkUpsertProductContent(ctx, productContents)
	if err != nil {
		return err
	}

	time.Sleep(2 * time.Minute)

	fmt.Println("^^^^^^^^^^^^^^^^^^^^^^BulkUpsertPackageContent^^^^^^^^^^^^^^^^^^^^^^^^^^")
	err = mrepo.BulkUpsertPackageContent(ctx, packageContents)
	if err != nil {
		return err
	}

	level.Info(logger).Log("Info", "Successfully synced")
	return nil
}

// FindDaysBetweenDate find number of days between two dates
func FindDaysBetweenDate(startdatetime string, enddatetime string) int {
	// Parse timestamps
	layout := time.RFC3339

	startTime, err := time.Parse(layout, startdatetime)
	if err != nil {
		fmt.Println("Error parsing startDateTime:", err)
		return -1
	}

	endTime, err := time.Parse(layout, enddatetime)
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
