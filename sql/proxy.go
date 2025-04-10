package sql

import (
	"database/sql"
)

// This file proxies from database/sql to monarch/sql

var (
	ErrNoRows   = sql.ErrNoRows
	ErrConnDone = sql.ErrConnDone
	ErrTxDone   = sql.ErrTxDone
)

type IsolationLevel sql.IsolationLevel

const (
	LevelDefault         IsolationLevel = IsolationLevel(sql.LevelDefault)
	LevelReadUncommitted IsolationLevel = IsolationLevel(sql.LevelReadUncommitted)
	LevelReadCommitted   IsolationLevel = IsolationLevel(sql.LevelReadCommitted)
	LevelWriteCommitted  IsolationLevel = IsolationLevel(sql.LevelWriteCommitted)
	LevelRepeatableRead  IsolationLevel = IsolationLevel(sql.LevelRepeatableRead)
	LevelSnapshot        IsolationLevel = IsolationLevel(sql.LevelSnapshot)
	LevelSerializable    IsolationLevel = IsolationLevel(sql.LevelSerializable)
	LevelLinearizable    IsolationLevel = IsolationLevel(sql.LevelLinearizable)
)

type TxOptions struct {
	Isolation IsolationLevel
	ReadOnly  bool
}

type NullBool sql.NullBool
type NullByte sql.NullByte
type NullFloat64 sql.NullFloat64
type Null16 sql.NullInt16
type NullInt32 sql.NullInt32
type NullInt64 sql.NullInt64
type NullString sql.NullString
type NullTime sql.NullTime

type Result sql.Result
