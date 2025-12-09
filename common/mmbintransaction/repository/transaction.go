package repository

import (
	"context"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/log/level"
	"github.com/jmoiron/sqlx"
)

// Transaction ...
type Transaction struct {
	ID                  int64     `db:"id"`
	ProgramCode         string    `db:"program_code"`
	RequestID           string    `db:"request_id"`
	FundingType         string    `db:"funding_type"`
	ProcessType         string    `db:"process_type"`
	Type                string    `db:"type"`
	TranactionType      string    `db:"transaction_type"`
	TranactionReference string    `db:"transaction_reference"`
	Action              string    `db:"action"`
	Stage               int64     `db:"stage"`
	Status              string    `db:"status"`
	DateAdded           time.Time `db:"date_addded"`
	DateModified        time.Time `db:"date_modified"`
}

// GetTranctionByHashID retrieve trasnaction by hash id
func (repo *Transaction) GetTranctionByHashID(ctx context.Context, conn *sqlx.DB, log log.Logger, hashID string) (p []Transaction, err error) {
	var (
		statement string
		args      []interface{}
		row       Transaction
	)

	statement, args, err = squirrel.Select(
		"id",
		"program_code",
		"request_id",
		"funding_type",
		"process_type",
		"type",
		"transaction_type",
		"transaction_reference",
		"action",
		"stage",
		"status",
		"date_added",
		"date_modified",
	).From(
		"transactions",
	).Where(
		squirrel.Eq{"request_id": hashID},
	).OrderBy(
		"id DESC",
	).ToSql()
	if err != nil {
		return p, err
	}

	statement = conn.Rebind(statement)
	level.Info(log).Log("sql query executed", statement)

	rows, err := conn.QueryContext(ctx,
		statement,
		args...,
	)
	if err != nil {
		level.Error(log).Log("database error", err.Error())
		return p, err
	}

	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(
			&row.ID,
			&row.ProgramCode,
			&row.RequestID,
			&row.FundingType,
			&row.ProcessType,
			&row.Type,
			&row.TranactionType,
			&row.TranactionReference,
			&row.Action,
			&row.Stage,
			&row.Status,
			&row.DateAdded,
			&row.DateModified,
		)
		if err != nil {
			level.Error(log).Log("database error", err.Error())
			return p, err
		}
		p = append(p, row)
	}

	if len(p) <= 0 {
		p = make([]Transaction, 0)
	}

	return p, err
}

// GetTranctionByFilters retrieve trasnaction using filters
func (repo *Transaction) GetTranctionByFilters(ctx context.Context, conn *sqlx.DB, log log.Logger, filters Transaction) (p []Transaction, err error) {
	var (
		statement string
		args      []interface{}
		row       Transaction
	)

	q := squirrel.Select(
		"id",
		"program_code",
		"request_id",
		"funding_type",
		"process_type",
		"type",
		"transaction_type",
		"transaction_reference",
		"action",
		"stage",
		"status",
		"date_added",
		"date_modified",
	).From(
		"transactions",
	)

	if filters.Stage > 0 {
		q = q.Where(squirrel.Eq{"stage": filters.Stage})
	}

	if filters.FundingType != "" {
		q = q.Where(squirrel.Eq{"funding_type": filters.FundingType})
	}

	if filters.ProcessType != "" {
		q = q.Where(squirrel.Eq{"process_type": filters.ProcessType})
	}

	if filters.ProgramCode != "" {
		q = q.Where(squirrel.Eq{"program_code": filters.ProgramCode})
	}

	if filters.Action != "" {
		q = q.Where(squirrel.Eq{"action": filters.Action})
	}

	if filters.TranactionType != "" {
		q = q.Where(squirrel.Eq{"transaction_type": filters.TranactionType})
	}

	if filters.TranactionReference != "" {
		q = q.Where(squirrel.Eq{"transaction_reference": filters.TranactionReference})
	}

	if filters.Type != "" {
		q = q.Where(squirrel.Eq{"type": filters.Type})
	}

	if filters.RequestID != "" {
		q = q.Where(squirrel.Eq{"request_id": filters.RequestID})
	}

	if filters.Status != "" {
		q = q.Where(squirrel.Eq{"status": filters.Status})
	}

	statement, args, err = q.OrderBy(
		"id DESC",
	).ToSql()
	if err != nil {
		return p, err
	}

	statement = conn.Rebind(statement)
	level.Info(log).Log("sql query executed", statement)

	rows, err := conn.QueryContext(ctx,
		statement,
		args...,
	)
	if err != nil {
		level.Error(log).Log("database error", err.Error())
		return p, err
	}

	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(
			&row.ID,
			&row.ProgramCode,
			&row.RequestID,
			&row.FundingType,
			&row.ProcessType,
			&row.Type,
			&row.TranactionType,
			&row.TranactionReference,
			&row.Action,
			&row.Stage,
			&row.Status,
			&row.DateAdded,
			&row.DateModified,
		)
		if err != nil {
			level.Error(log).Log("database error", err.Error())
			return p, err
		}
		p = append(p, row)
	}

	if len(p) <= 0 {
		p = make([]Transaction, 0)
	}

	return p, err
}

// InsertTransaction ...
func (repo *Transaction) InsertTransaction(ctx context.Context, conn *sqlx.DB, log log.Logger, data Transaction) (lastID int64, err error) {
	var (
		createdDate = time.Now().UTC().Format("2006-01-02 15:04:05")
	)

	err = conn.QueryRow("INSERT INTO transactions(program_code,request_id,funding_type,process_type,type,transaction_type,transaction_reference,action,stage,status,date_modified) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11) RETURNING ID", data.ProgramCode, data.RequestID, data.FundingType, data.ProcessType, data.Type, data.TranactionType, data.TranactionReference, data.Action, data.Stage, data.Status, createdDate).Scan(&lastID)

	if err != nil {
		level.Error(log).Log("database error on getting last inserted id", err.Error())
		return lastID, err
	}

	return lastID, err

	/*q := squirrel.Insert(
		"transactions",
	).Columns(
		"program_code",
		"request_id",
		"funding_type",
		"type",
		"transaction_type",
		"transaction_reference",
		"action",
		"stage",
		"status",
		"date_modified",
	).Values(
		data.ProgramCode,
		data.RequestID,
		data.FundingType,
		data.Type,
		data.TranactionType,
		data.TranactionReference,
		data.Action,
		data.Stage,
		data.Status,
		createdDate,
	)

	statement, args, _ := q.ToSql()
	level.Info(log).Log("sql query executed", statement)
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

// UpdateTransaction ...
func (repo *Transaction) UpdateTransaction(ctx context.Context, conn *sqlx.DB, log log.Logger, data Transaction) (err error) {
	var (
		currentDate = time.Now().UTC().Format("2006-01-02 15:04:05")
		query       squirrel.UpdateBuilder
	)

	if data.ID > 0 {
		query = squirrel.Update("transactions").Set("date_modified", currentDate).
			Where(squirrel.Eq{"id": data.ID})
	} else {
		query = squirrel.Update("transactions").Set("date_modified", currentDate).
			Where(squirrel.Eq{"request_id": data.RequestID}).
			Where(squirrel.Eq{"funding_type": data.FundingType}).
			Where(squirrel.Eq{"type": data.Type}).
			Where(squirrel.Eq{"stage": data.Stage})
	}

	if data.Status != "" {
		query = query.Set("status", data.Status)
	}

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
