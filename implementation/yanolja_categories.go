package implementation

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"swallow-supplier/config"
	customError "swallow-supplier/error"
	model "swallow-supplier/mongo/domain/yanolja"
	"swallow-supplier/request_response/yanolja"
	yanoljasvc "swallow-supplier/services/suppliers/yanolja"
	"swallow-supplier/utils"
	"swallow-supplier/utils/constant"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/xuri/excelize/v2"
)

// GetCategories get all products categories from yanolja
func (s *service) GetCategories(ctx context.Context) (resp yanolja.Response, err error) {
	var requestID = utils.GenerateUUID("GGT", true)

	logger := log.With(
		s.logger,
		"method", "GetCategories",
		"Request ID", requestID,
	)

	// Defer when panic
	defer func(context.Context) {
		r := recover()
		if r == nil {
			return
		}

		level.Info(logger).Log("info", "processing request went into panic mode")

		resp.Code = "500"

	}(ctx)

	var getsvc, _ = yanoljasvc.New(ctx)
	resp, err = getsvc.GetAllCategories(ctx)

	if err != nil {
		level.Error(logger).Log("error", "request to yanolja client raise error ", err)
		resp.Code = "500"
		resp.Body = "yanolja error"
		return resp, err
	}
	level.Info(logger).Log("response from yanolja", resp)

	doc, err := json.Marshal(resp.Body)
	if err != nil {
		resp.Code = "500"
		return resp, fmt.Errorf("marshal error ")
	}

	// Unmarshal JSON response
	var rec []model.Category
	err = json.Unmarshal([]byte(doc), &rec)
	if err != nil {
		level.Error(logger).Log("error", "failed to unmarshal JSON response", err)
		resp.Code = "500"
		return resp, fmt.Errorf("json unmarshal error: %w", err)
	}

	err = s.mongoRepository[config.Instance().MongoDBName].InsertCategories(ctx, rec)
	if err != nil {
		level.Error(logger).Log("error", "request to InsertCategories raised error ", err)
		resp.Code = "500"
		return resp, customError.NewError(ctx, "leisure-api-1015", fmt.Sprintf("repository error from InsertCategories, %v", err), "InsertCategories")
	}

	return resp, nil

}

// InsertAllCategories insert all product's categories from yanolja to ggt
// This api reads from yanolja serve and write it to GGT
func (s *service) InsertAllCategories(ctx context.Context) (resp yanolja.Response, err error) {
	var requestID = utils.GenerateUUID("GGT", true)

	logger := log.With(
		s.logger,
		"method", "InsertAllCategories",
		"Request ID", requestID,
	)

	// Defer when panic
	defer func(context.Context) {
		r := recover()
		if r == nil {
			return
		}

		level.Info(logger).Log("info", "processing request went into panic mode")

		resp.Code = "500"

	}(ctx)

	var getsvc, _ = yanoljasvc.New(ctx)
	resp, err = getsvc.GetAllCategories(ctx)

	if err != nil {
		level.Error(logger).Log("error", "request to yanolja client raise error ", err)
		resp.Code = "500"
		resp.Body = "yanolja error"
		return resp, err
	}
	level.Info(logger).Log("response from yanolja", resp)

	doc, err := json.Marshal(resp.Body)
	if err != nil {
		resp.Code = "500"
		return resp, fmt.Errorf("marshal error ")
	}

	// Unmarshal JSON response
	var rec []model.Category
	err = json.Unmarshal([]byte(doc), &rec)
	if err != nil {
		level.Error(logger).Log("error", "failed to unmarshal JSON response", err)
		resp.Code = "500"
		return resp, fmt.Errorf("json unmarshal error: %w", err)
	}

	for i := 0; i < len(rec); i++ {
		rec[i].SupplierGGTChannel = constant.YANOLJAGGTTRIP
	}

	err = s.mongoRepository[config.Instance().MongoDBName].InsertCategories(ctx, rec)
	if err != nil {
		level.Error(logger).Log("error", "request to InsertCategories raised error ", err)
		resp.Code = "500"
		return resp, customError.NewError(ctx, "leisure-api-1015", fmt.Sprintf("repository error from InsertCategories, %v", err), "InsertCategories")
	}

	return resp, nil

}

// UpsertCategoryMapping upsert all  products categories from yanolja- ggt - trip mapping
// reading from excel provided by @james
func (s *service) UpsertCategoryMapping(ctx context.Context) (resp yanolja.Response, err error) {
	var requestID = utils.GenerateUUID("GGT", true)

	logger := log.With(
		s.logger,
		"method", "UpsertCategoryMapping",
		"Request ID", requestID,
	)

	// Defer when panic
	defer func(context.Context) {
		r := recover()
		if r == nil {
			return
		}

		level.Info(logger).Log("info", "processing request went into panic mode")

		resp.Code = "500"

	}(ctx)
	fmt.Println("1")
	rows, err := readExcelFile(config.Instance().YGTFilePath, logger)
	if err != nil {
		level.Error(logger).Log("msg", "Error reading Excel file", "err", err)
		resp.Code = "500"
		return resp, fmt.Errorf("Error reading Excel file")
	}

	fmt.Println("2")

	records, headers, err := cleanData(rows, logger)
	fmt.Println("Headers:", headers)
	if err != nil {
		level.Error(logger).Log("msg", "Error cleaning data", "err", err)
		resp.Code = "500"
		return resp, fmt.Errorf("Error cleaning data")
	}
	fmt.Println("3")

	//
	/* err = deleteFile(cleanedFilePath, logger)
	if err != nil {
		level.Error(logger).Log("msg", "Error deleting cleaned_data.xlsx file ", "err", err)
		resp.Code = "500"
		return resp, fmt.Errorf("Error deleting cleaned_data.xlsx file")
	} */

	modifiedRecordCount, err := s.mongoRepository[config.Instance().MongoDBName].BulkUpsertCategoryMapping(ctx, records)
	if err != nil {
		level.Error(logger).Log("msg", "Error storing data in MongoDB", "err", err)
		resp.Code = "500"
		return resp, customError.NewError(ctx, "leisure-api-1015", fmt.Sprintf("repository error from BulkUpsertCategoryMapping, %v", err), "BulkUpsertCategoryMapping")
	}

	fmt.Println("4")

	level.Info(logger).Log("msg", "Data processing and MongoDB upload completed successfully")
	resp.Code = "200"
	resp.Body = modifiedRecordCount

	return
}

// deleteFile checks if the file exists before attempting to delete it.
func deleteFile(path string, logger log.Logger) error {
	if _, err := os.Stat(path); err == nil {
		if err := os.Remove(path); err != nil {
			level.Error(logger).Log("msg", "Failed to delete file", "file", path, "err", err)
			return err
		}
		level.Info(logger).Log("msg", "Deleted existing file", "file", path)
	}
	return nil
}

// readExcelFile opens the Excel file and returns all rows from the first sheet.
func readExcelFile(filePath string, logger log.Logger) ([][]string, error) {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open excel file: %w", err)
	}
	defer f.Close()

	sheetName := f.GetSheetName(0)
	if sheetName == "" {
		return nil, errors.New("no sheets found in excel file")
	}

	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, fmt.Errorf("failed to get rows: %w", err)
	}
	level.Info(logger).Log("msg", "Excel file loaded", "rows", len(rows), "sheet", sheetName)
	return rows, nil
}

// cleanData processes the Excel rows and ensures data integrity.
func cleanData(rows [][]string, logger log.Logger) ([]map[string]interface{}, []string, error) {
	if len(rows) < 2 {
		return nil, nil, errors.New("not enough rows in excel file")
	}

	headers := rows[1] // Use the second row as headers
	validCols := []int{}
	validHeaders := []string{}

	for i, header := range headers {
		if strings.TrimSpace(header) != "" {
			validCols = append(validCols, i)
			validHeaders = append(validHeaders, header)
		}
	}
	if len(validCols) == 0 {
		return nil, nil, errors.New("no valid headers found")
	}

	records := []map[string]interface{}{}
	colNumeric := make(map[int]bool)
	for _, col := range validCols {
		colNumeric[col] = true
	}
	for r := 2; r < len(rows); r++ {
		for _, col := range validCols {
			if col < len(rows[r]) && strings.TrimSpace(rows[r][col]) != "" {
				if _, err := strconv.ParseFloat(rows[r][col], 64); err != nil {
					colNumeric[col] = false
				}
			}
		}
	}

	for r := 2; r < len(rows); r++ {
		record := make(map[string]interface{})
		for _, col := range validCols {
			header := headers[col]
			var value interface{}
			if col < len(rows[r]) {
				cellValue := strings.TrimSpace(rows[r][col])
				if colNumeric[col] {
					if cellValue == "" {
						value = 0
					} else {
						fVal, err := strconv.ParseFloat(cellValue, 64)
						if err != nil {
							value = 0
						} else {
							value = int(fVal)
						}
					}
				} else {
					value = cellValue
				}
			}
			record[header] = value
			record["supplier_GGT_Channel"] = constant.YANOLJAGGTTRIP
		}
		records = append(records, record)
	}
	level.Info(logger).Log("msg", "Data cleaning completed", "records", len(records))
	return records, validHeaders, nil
}
