package mmbintransaction

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"swallow-supplier/common/mmbintransaction/constant"
	co "swallow-supplier/common/mmbintransaction/services/ggt"
	st "swallow-supplier/common/mmbintransaction/struct"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/log/level"

	"swallow-supplier/common/mmbintransaction/repository"

	"github.com/jmoiron/sqlx"
)

// OpenloopReversal function deals with reversal business logic of MM across different provider
func OpenloopReversal(ctx context.Context, dbconn *sqlx.DB, log log.Logger, rev st.ReversalRequest) (res ReversalResponse, err error) {
	//programCode, fundingType, action, hashID, transactionType, processType string
	var (
		repo                                                                            repository.Transaction
		pendingReversal, initiatedReversal, successReversal, failedReversal             []repository.Transaction
		successTransaction, failedTransaction, pendingTransaction, initiatedTransaction []repository.Transaction
		configReq                                                                       st.ConfigRequest
		stage, transactionIdelTime, finalLevel                                          int64
	)

	configReq.ProductCode = rev.ProgramCode
	configReq.ResourceType = strings.ToLower(rev.FundingType + "_" + rev.Action)
	configReq.Key = "config"
	configReq.IsActive = "1"

	dataConfig, err := GetConfigurations(ctx, dbconn, log, configReq)

	fmt.Println("get config data", dataConfig)

	if len(dataConfig.Records) > 0 {
		for _, configVal := range dataConfig.Records {
			if configVal.Type == constant.ValueTypeString {
				strtointstage, err := strconv.Atoi(configVal.Value.(string))
				if err != nil {
					return res, err
				}
				stage = int64(strtointstage)
			}
			if configVal.Type == constant.ValueTypeJSON {
				var customValidations = make(map[string]interface{})

				confByt, err := json.Marshal(configVal)
				if err != nil {
					return res, err
				}
				if len(confByt) > 0 {
					err := json.Unmarshal(confByt, &customValidations)
					if err != nil {
						return res, err
					}
				}

				fmt.Println("consfig customValidations data", customValidations)
				var configValue = customValidations["value"].(map[string]interface{})

				fmt.Println("configValue", configValue)
				fmt.Println("total_stage", configValue["total_stages"])
				if _, ok := configValue["total_stages"]; ok {
					stage = int64(configValue["total_stages"].(float64))
					fmt.Println("stage", stage)
				}

				if _, ok := configValue["transaction_idel_time"]; ok {
					transactionIdelTime = int64(configValue["transaction_idel_time"].(float64))
					fmt.Println("stage", stage)
				}

				if _, ok := configValue["final_level"]; ok {
					finalLevel = int64(configValue["final_level"].(float64))
					fmt.Println("stage", stage)
				}

			}

		}
	}

	if stage == 0 {
		return res, errors.New("process stage not defined")
	}

	// for old tranctions if binsponor is not enabled and if we recive an reversal then without transaction reverse the amount like CV
	if !rev.BinSponsorEnabled {
		if rev.RetryAttempt > 0 { // if we have the same hash with retry attempt is more then once then we have check the pervious status and decide and new
			res.CheckStatus = stage
			res.Type = constant.Reversal
			return res, err
		}
		res.IniateReversal = true
		res.ReversalLevel = stage
		res.Type = constant.Reversal
		return res, err
	}

	// if settlemnet and transaction type is credit vocucher then base trasnaction will not be present just oly reversal ex cash back
	if strings.ToLower(rev.ProcessType) == constant.Settelement && strings.ToLower(rev.TransactionType) == constant.CreditVoucher {
		if rev.RetryAttempt > 0 { // if we have the same hash with retry attempt is more then once then we have check the pervious status and decide and new
			res.CheckStatus = stage
			res.Type = constant.Reversal
			return res, err
		}
		res.IniateReversal = true
		res.ReversalLevel = stage
		res.Type = constant.Reversal
		return res, err
	}

	tranactionData, err := repo.GetTranctionByHashID(ctx, dbconn, log, rev.HashID)
	if err != nil {
		println(err.Error())
		return res, err
	}
	fmt.Println("tranactionData", tranactionData)

	if len(tranactionData) == 0 {
		fmt.Println("no transaction records")
		res.NoAction = true
		return res, err
	}
	for _, data := range tranactionData {
		switch data.Type {
		case constant.Reversal:
			switch data.Status {
			case constant.Pending:
				pendingReversal = append(pendingReversal, data)
			case constant.Initiated:
				initiatedReversal = append(initiatedReversal, data)
			case constant.Success:
				successReversal = append(successReversal, data)
			case constant.Failed:
				failedReversal = append(failedReversal, data)
			}
		case constant.Transaction:
			switch data.Status {
			case constant.Pending:
				pendingTransaction = append(pendingTransaction, data)
			case constant.Initiated:
				initiatedTransaction = append(initiatedTransaction, data)
			case constant.Success:
				successTransaction = append(successTransaction, data)
			case constant.Failed:
				failedTransaction = append(failedTransaction, data)
			}
		}
	}

	data := tranactionData[0]

	// check the transaction is ongoing within expected time as per config
	elaspedTime := GetTransactionElaspedTimeInMinutes(data.DateModified)
	fmt.Println("data", data)
	if data.Type == constant.Transaction && data.Stage <= stage && elaspedTime < float64(transactionIdelTime) && !strings.Contains("success,failed", data.Status) {
		fmt.Println("transaction in-progress within allowed time so push the message back to queue")
		res.PushToQueue = true
		return res, err
	}

	if data.Type == constant.Transaction && data.Stage <= stage && elaspedTime > float64(transactionIdelTime) && !strings.Contains("success,failed", data.Status) {
		fmt.Println("transaction in-progress beyound the allowed time so we willhave to check the last record status at bin levbel and do the reversal based on the stage")
		res.CheckStatus = data.Stage // check from current stage and below it
		res.Type = constant.Transaction
		return res, err
	}

	// For transaction ReversalLevel is set based on the level need to be revered if it is > 0 then oly that particular level and level below it will be reversed if it is 0 then all level will be reversed
	if data.Type == constant.Transaction && strings.Contains(constant.Success, data.Status) {
		fmt.Println("transaction last record is success so initiate the reversal based on the stage")
		res.IniateReversal = true
		res.ReversalLevel = data.Stage
		res.Type = constant.Transaction
		//for case if autoincrement has some issue and lower stage has greater autoincrement value
		if int64(len(successTransaction)) == stage {
			res.ReversalLevel = stage
		}
		return res, err
	}

	if data.Type == constant.Transaction && strings.Contains(constant.Failed, data.Status) {
		fmt.Println("last transaction data is failed check the stage and perform the reversal if the stage is greather than 1 then check the avaiable success transction and do the reversal")
		res.NoAction = true
		if data.Stage == stage {
			if stage > 1 {
				//checking the last record is success
				if len(successTransaction) > 0 {
					res.IniateReversal = true
					res.ReversalLevel = successTransaction[len(successTransaction)-1].Stage
					res.Type = constant.Transaction
					res.NoAction = false
					return res, err
				}
			}
		}
		return res, err
	}

	fmt.Println("reversal data", data)

	if data.Type == constant.Reversal && data.Stage <= stage && elaspedTime < float64(transactionIdelTime) && !strings.Contains("success,failed", data.Status) {
		fmt.Println("reversal in-progress within allowed time so push the message back to queue")
		res.PushToQueue = true
		return res, err
	}

	if data.Type == constant.Reversal && data.Stage <= stage && elaspedTime > float64(transactionIdelTime) && !strings.Contains("success,failed", data.Status) {
		fmt.Println("reversal in-progress beyound the allowed time so we will have to check the last record status at bin and do the reversal based on the stage status")
		res.CheckStatus = data.Stage // check the current status from the stage till total stages
		res.Type = constant.Reversal
		return res, err
	}

	// For Reversal ReversalLevel is set based on the level need to perform reversal for all the stages upto total stages from the stage mentioned
	// check this use case
	if data.Type == constant.Reversal && strings.Contains(constant.Success, data.Status) {
		fmt.Println("reversal last record is success so initiate the reversal based on the stage")
		if data.Stage == finalLevel { // since reversal level 2 first and level is final so if level 1 is done then reversal was successful so NO ACTION
			fmt.Println("if the last reversal record is the final stage then no action of reversal")
			res.NoAction = true
			return res, err
		}
		fmt.Println("if the last reversal is not the final stage then we have initiate reversal for that stage")
		res.IniateReversal = true
		res.NoAction = false
		res.ReversalLevel = data.Stage - 1 // reversal stage 1 is success then we have perform reversal from stage 2 till total stages
		res.Type = constant.Reversal
		return res, err
	}

	if data.Type == constant.Reversal && strings.Contains(constant.Failed, data.Status) {
		fmt.Println("reversal last record is in failed stage")
		res.IniateReversal = true
		res.Type = constant.Reversal
		res.ReversalLevel = data.Stage
		if stage > 1 {
			//checking the last record is success is success
			if len(successReversal) > 0 {
				fmt.Println("have to perform for the stage which is not successfull might occur due to auto increment value mismatch")
				//wrong implemenattion
				res.ReversalLevel = successReversal[len(successReversal)-1].Stage - 1
				if res.ReversalLevel == finalLevel {
					fmt.Println("if the success full reversal is final stage then no action to perform since all the level is done")
					res.NoAction = true
					res.IniateReversal = false
					res.Type = ""
					res.ReversalLevel = 0
					return res, err
				}
			}
		}
		return res, err
	}

	return res, err
}

// OpenloopAuthorize ...
func OpenloopAuthorize(ctx context.Context, dbconn *sqlx.DB, log log.Logger, hashID string) (res AuthorizeResponse, err error) {
	var (
		repo repository.Transaction
	)
	tranactionData, err := repo.GetTranctionByHashID(ctx, dbconn, log, hashID)
	if err != nil {
		return res, err
	}
	if len(tranactionData) == 0 {
		res.IniateAuthorize = true
	} else {
		res.NoAction = true
	}

	return res, err
}

// GetTranctionData ...
func GetTranctionData(ctx context.Context, dbconn *sqlx.DB, log log.Logger, transaction repository.Transaction) (res []repository.Transaction, err error) {

	level.Info(log).Log("info", "GetTranctionData")
	var (
		repo repository.Transaction
	)

	tranactionData, err := repo.GetTranctionByFilters(ctx, dbconn, log, transaction)
	if err != nil {
		println(err.Error())
		return res, err
	}
	return tranactionData, err
}

// GetConfigurations ...
func GetConfigurations(ctx context.Context, dbconn *sqlx.DB, log log.Logger, data interface{}) (res st.ConfigResponse, err error) {
	//logger := log.With(log, "method", "GetConfigurations")
	level.Info(log).Log("info", "Retrieve Configurations")

	var (
		repo repository.Configuration
		item st.Configs
		req  st.ConfigRequest
	)

	//var queueUnmarshalData map[string]interface{}
	queueMarshalData, err := json.Marshal(data)
	if err != nil {
		level.Info(log).Log("call", "Publish to webhook service kafka",
			"error", err)
	}
	fmt.Println("packagequeueUnmarshalData", string(queueMarshalData))
	json.Unmarshal(queueMarshalData, &req)
	fmt.Println("packagequeueUnmarshalData", req)

	if req.Version == "" {
		req.Version = "v1"
	}

	// Retrieve the Configurations
	listConfigs, err := repo.GetAllConfigurations(ctx, dbconn, log, req)
	if err != nil {
		level.Error(log).Log("error", "GetConfigurations: ", err)
		return res, err
	}

	for _, record := range listConfigs {
		item.Key = record.Key
		item.Value = record.ValueString.String
		if record.ValueType == constant.ValueTypeJSON {
			item.Value = json.RawMessage(record.ValueJSON.String)
		}

		item.Type = record.ValueType
		res.Records = append(res.Records, item)
	}

	res.ProductCode = req.ProductCode
	res.ResourceType = req.ResourceType
	res.Version = req.Version

	if len(res.Records) <= 0 {
		res.Records = make([]st.Configs, 0)
	}

	// Cache the response
	/*cacheData, err := json.Marshal(res)
	if err == nil {
		err = cacheLayer.SetTTL(ctx, cacheID, string(cacheData), (time.Hour)*(constant.CacheExpiryInHours))
		if err != nil {
			level.Error(logger).Log("error", "Error writing to cache: %s", err)
		}
	}*/

	return res, err
}

// InsertTransaction ...
func InsertTransaction(ctx context.Context, dbcon *sqlx.DB, log log.Logger, data repository.Transaction) (lastID int64, err error) {
	var repo repository.Transaction
	return repo.InsertTransaction(ctx, dbcon, log, data)
}

// UpdateTransaction ...
func UpdateTransaction(ctx context.Context, dbcon *sqlx.DB, log log.Logger, data repository.Transaction) (err error) {

	var repo repository.Transaction
	return repo.UpdateTransaction(ctx, dbcon, log, data)
}

// InsertTransactionLog ...
func InsertTransactionLog(ctx context.Context, dbcon *sqlx.DB, log log.Logger, data repository.TransactionLog) (lastID int64, err error) {

	var repo repository.TransactionLog
	return repo.InsertTransactionLog(ctx, dbcon, log, data)
}

// UpdateTransactionLog ...
func UpdateTransactionLog(ctx context.Context, dbcon *sqlx.DB, log log.Logger, data repository.TransactionLog) (err error) {

	var repo repository.TransactionLog
	return repo.UpdateTransactionLog(ctx, dbcon, log, data)
}

// GetTransactionElaspedTimeInMinutes ...
func GetTransactionElaspedTimeInMinutes(transactiomTime time.Time) float64 {
	currentTime := time.Now()
	return currentTime.Sub(transactiomTime).Minutes()
}

// GetUserDetails ...
func GetUserDetails(log log.Logger, customerServiceHost, programCode, userHashID, opKey, opSecret string) (map[string]interface{}, error) {
	return co.GetUserDetails(log, customerServiceHost, programCode, userHashID, opKey, opSecret)
}
