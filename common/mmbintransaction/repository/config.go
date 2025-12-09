package repository

import (
	"context"
	"fmt"

	"swallow-supplier/common/mmbintransaction/constant"
	st "swallow-supplier/common/mmbintransaction/struct"

	"github.com/Masterminds/squirrel"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/log/level"
	"github.com/guregu/null"
	"github.com/jmoiron/sqlx"
)

// Configuration ...
type Configuration struct {
	ID           string      `db:"id"`
	HashID       string      `db:"hash_id"`
	Key          string      `db:"key"`
	ValueString  null.String `db:"value_string"`
	ValueJSON    null.String `db:"value_json"`
	ValueType    string      `db:"value_type"`
	IsActive     string      `db:"is_active"`
	ResourceType string      `db:"resource_type"`
	ProductCode  string      `db:"product_code"`
	IsDeleted    string      `db:"is_deleted"`
}

// ConfigList ...
type ConfigList []Configuration

// GetConfigurations ...
func (me *Configuration) GetConfigurations(ctx context.Context, dbconn *sqlx.DB, log log.Logger, req st.ConfigRequest, entityName string, excludeKeys []string) (list ConfigList, err error) {
	var (
		statement string
		args      []interface{}
		row       Configuration
	)

	query := squirrel.Select(
		"id",
		"hash_id",
		"key",
		"value_string",
		"value_json",
		"value_type",
		"is_active",
		"resource_type",
		"is_deleted",
	).From(
		entityName,
	)

	if len(excludeKeys) != 0 {
		query = query.Where(squirrel.NotEq{"key": excludeKeys})
	}

	if req.ID != "" {
		query = query.Where(squirrel.Eq{"hash_id": req.ID})
	} else {
		if req.Key != "" {
			query = query.Where(squirrel.Eq{"key": req.Key})
		}

		if req.ResourceType != "" {
			query = query.Where(squirrel.Eq{"resource_type": req.ResourceType})
		}

		if entityName == constant.TableProductConfigs {
			if req.ProductCode != "" {
				query = query.Where(squirrel.Eq{"product_code": req.ProductCode})
			}
		}
	}

	if req.IsActive != "" {
		query = query.Where(squirrel.Eq{"is_active": req.IsActive})
	}

	statement, args, err = query.OrderBy("id").ToSql()
	if err != nil {
		level.Error(log).Log("database error", err.Error())
		return list, err
	}

	statement = dbconn.Rebind(statement)
	level.Info(log).Log("sql query executed", statement)

	rows, err := dbconn.DB.QueryContext(ctx,
		statement,
		args...,
	)
	if err != nil {
		level.Error(log).Log("database error", err.Error())
		return list, err
	}

	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(
			&row.ID,
			&row.HashID,
			&row.Key,
			&row.ValueString,
			&row.ValueJSON,
			&row.ValueType,
			&row.IsActive,
			&row.ResourceType,
			&row.IsDeleted,
		)
		if err != nil {
			level.Error(log).Log("database error", err.Error())
			return list, err
		}
		list = append(list, row)
	}

	if len(list) <= 0 {
		list = make([]Configuration, 0)
	}
	return list, err
}

// GetAllConfigurations ...
func (me *Configuration) GetAllConfigurations(ctx context.Context, dbconn *sqlx.DB, log log.Logger, req st.ConfigRequest) (list ConfigList, err error) {
	var (
		excludeKeys []string
	)

	/*cacheID, getCache, cacheLayer, err := repo.RepositoryCache(req)
	if err != nil {
		level.Error(log).Log("message", "Error on RepositoryCache", "error", err, "method", "GetConfigurations", "cacheID", cacheID, "req", req)
	}
	if getCache != "" {
		err := json.Unmarshal([]byte(getCache), &list)
		if err != nil {
			level.Error(log).Log("message", "Error on unmarshal cache", "error", err, "method", "GetConfigurations", "cacheID", cacheID, "req", req)
		}
		return list, err
	}*/

	tableList := [2]string{constant.TableProductConfigs, constant.TableProductBinConfigs}

	for _, table := range tableList {
		fmt.Println(table)
		items, _ := me.GetConfigurations(ctx, dbconn, log, req, table, excludeKeys)
		for _, item := range items {
			if item.IsDeleted == "0" {
				if item.IsActive == "1" {
					list = append(list, item)
				}
				excludeKeys = append(excludeKeys, item.Key)
			}
		}
	}

	if len(list) <= 0 {
		list = make([]Configuration, 0)
	}

	// Cache the response
	/*rawData, err := json.Marshal(list)
	cacheData := rawData
	if err == nil {
		err = cacheLayer.SetTTL(context.Background(), cacheID, string(cacheData), (time.Hour)*(constant.CacheExpiryInHours))
		if err != nil {
			level.Error(repo.logger).Log("error", "Error writing to cache: %s", err)
		}
	}*/
	return list, err
}
