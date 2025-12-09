package repository

import (
	"context"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/log/level"
	"github.com/jmoiron/sqlx"
)

// TransactionLog ...
type TransactionLog struct {
	ID               int64     `db:"id"`
	TransactionID    int64     `db:"transaction_id"`
	StageReferenceID string    `db:"stage_reference_id"`
	Request          string    `db:"request"`
	Response         string    `db:"response"`
	DateAdded        time.Time `db:"date_addded"`
	DateModified     time.Time `db:"date_modified"`
}

// InsertTransactionLog ...
func (repo *TransactionLog) InsertTransactionLog(ctx context.Context, conn *sqlx.DB, log log.Logger, data TransactionLog) (lastID int64, err error) {
	var (
	//createdDate = time.Now().UTC().Format("2006-01-02 15:04:05")
	)

	err = conn.QueryRow("INSERT INTO transaction_logs(transaction_id,stage_reference_id,request) VALUES($1,$2,$3) RETURNING ID", data.TransactionID, data.StageReferenceID, data.Request).Scan(&lastID)

	if err != nil {
		level.Error(log).Log("database error on getting last inserted id", err.Error())
		return lastID, err
	}

	return lastID, err

	/*q := squirrel.Insert(
		"transaction_logs",
	).Columns(
		"transaction_id",
		"stage_reference_id",
		"request",
		"response",
		"date_modified",
	).Values(
		data.TransactionID,
		data.StageReferenceID,
		data.Request,
		data.Response,
		createdDate,
	)
	statement, args, _ := q.ToSql()
	statement = conn.Rebind(statement)
	level.Info(log).Log("sql query executed", statement)

	queryResult, err := conn.Exec(
		statement,
		args...,
	)

	if err != nil {
		level.Error(log).Log("database error", err.Error())
		return lastID, err
	}

	lastID, err = queryResult.LastInsertId()

	if err != nil {
		level.Error(log).Log("database error on getting last inserted id", err.Error())
		return lastID, err
	}

	return lastID, err*/
}

// UpdateTransactionLog ...
func (repo *TransactionLog) UpdateTransactionLog(ctx context.Context, conn *sqlx.DB, log log.Logger, data TransactionLog) (err error) {
	var (
		currentDate = time.Now().UTC().Format("2006-01-02 15:04:05")
		query       squirrel.UpdateBuilder
	)

	if data.ID > 0 {
		query = squirrel.Update("transaction_logs").Set("date_modified", currentDate).Set("response", data.Response).
			Where(squirrel.Eq{"id": data.ID})
		statement, args, err := query.ToSql()
		if err != nil {
			level.Error(log).Log("database error", err.Error())
		}

		statement = conn.Rebind(statement)
		level.Info(log).Log("sql query executed", statement)

		_, err = conn.ExecContext(ctx,
			statement,
			args...,
		)

		if err != nil {
			level.Error(log).Log("database error", err.Error())
		}
		return err
	}
	return err
}
