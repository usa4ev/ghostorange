package psqldb

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
	_ "github.com/jackc/pgx/stdlib"

	"ghostorange/internal/app/auth"
	"ghostorange/internal/app/auth/session"
	"ghostorange/internal/app/model"
	"ghostorange/internal/pkg/encryption"
)

type (
	Database struct {
		*sql.DB
	}
)

func New(dsn string) (*Database, error) {
	var (
		db  Database
		err error
	)

	db.DB, err = sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("cannot connect to Database: %w", err)
	}

	err = db.initDB()
	if err != nil {
		return nil, fmt.Errorf("cannot init Database: %w", err)
	}

	return &db, nil
}

func (db Database) initDB() error {
	query := `CREATE TABLE IF NOT EXISTS users (
				id VARCHAR(100) PRIMARY KEY,
				username VARCHAR(256) not null,
                pwdhash VARCHAR(256) not null);`

	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create table users, %v", err)
	}

	// ItemCredentials struct {
	// 	ID          string      `json:"id"`
	// 	Credentials Credentials `json:"credentials"`
	// 	Name        string      `json:"name"`
	// 	Comment     string      `json:"comment"`
	// }

	query = `CREATE TABLE IF NOT EXISTS credentials (
					id VARCHAR(100) PRIMARY KEY UNIQUE,
					user_id VARCHAR(100) not null,
					ts timestamptz not null,
					encrypted bytea not null,
					name varchar(100) not null,
					comment varchar(1000) not null,
					FOREIGN KEY (user_id)
				REFERENCES users (id));`

	_, err = db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create table credentials, %v", err)
	}

	// ItemText struct {
	// 	ID      string `json:"id"`
	// 	Text    string `json:"text"`
	// 	Name    string `json:"name"`
	// 	Comment string `json:"comment"`
	// }

	query = `CREATE TABLE IF NOT EXISTS text (
		id VARCHAR(100) PRIMARY KEY UNIQUE,
		user_id varchar(100) not null,
		ts timestamptz not null,
		text bytea not null,
		name varchar(100) not null,
		comment varchar(1000) not null,
		FOREIGN KEY (user_id)
	REFERENCES users (id));`

	_, err = db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create table text, %v", err)
	}

	// ItemBinary struct {
	// 	ID        string `json:"id"`
	// 	Size      int    `json:"size"`
	// 	Extention string `json:"extention"`
	// 	Data      string `json:"data"`
	// 	Name      string `json:"name"`
	// 	Comment   string `json:"comment"`
	// }

	query = `CREATE TABLE IF NOT EXISTS binarydata (
		id VARCHAR(100) PRIMARY KEY UNIQUE,
		user_id varchar(100) not null,
		ts timestamptz not null,
		data bytea not null,
		extention bytea not null,
		size int not null,
		name varchar(255) not null,
		comment varchar(1000) not null,
		FOREIGN KEY (user_id)
	REFERENCES users (id));`

	_, err = db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create table binary, %v", err)
	}

	// ItemCard struct {
	// 	ID                 string    `json:"id"`
	// 	Number             string    `json:"number"`
	// 	Exp                time.Time `json:"expiration_date"`
	// 	CardholderName     string    `json:"holder_name"`
	// 	CardholderSurename string    `json:"holder_surename"`
	// 	CVVHash            string    `json:"cvv_hash"`
	// 	Name               string    `json:"name"`
	// 	Comment            string    `json:"comment"`
	// }

	query = `CREATE TABLE IF NOT EXISTS cards (
		id VARCHAR(100) PRIMARY KEY UNIQUE,
		user_id varchar(100) not null,
		ts timestamptz not null,
		number char(8) not null,
		expires date not null,
		cardholderName varchar(255) not null,
		cardholderSurename varchar(255) not null,
		cvvhash varchar(255) not null,
		name varchar(255) not null,
		comment varchar(1000) not null,
		FOREIGN KEY (user_id)
	REFERENCES users (id));`

	_, err = db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create table cards, %v", err)
	}

	return err
}

func (db Database) execInsUpdStatement(ctx context.Context, query string, args ...interface{}) (int64, error) {
	tx, err := db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	res, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, fmt.Errorf("error when executing query context %w", err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("error when finding rows affected %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return 0, err
	}

	return affected, nil
}

// AddUser adds new row to Database and return new user ID or error if addition failed
func (db Database) AddUser(ctx context.Context, username, hash string) (string, error) {
	id := uuid.New().String()

	query := `INSERT INTO users(id, username, pwdhash) VALUES ($1, $2, $3) ON CONFLICT (id) DO NOTHING;`

	rowsAffected, err := db.execInsUpdStatement(ctx, query, id, username, hash)

	if err != nil {
		return "", err
	} else if rowsAffected == 0 {
		return "", auth.ErrUserAlreadyExists
	}

	return id, nil
}

// UserExists returns true if user found by given userName or false otherwise
func (db Database) UserExists(ctx context.Context, userName string) (bool, error) {
	var exists bool

	query := "SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)"

	err := db.QueryRowContext(ctx, query, userName).Scan(&exists)

	if !exists || errors.Is(err, sql.ErrNoRows) {
		return false, nil
	} else if err != nil {
		return false, fmt.Errorf("failed to check if user exists: %w", err)
	}

	return true, nil
}

// GetPasswordHash returns user ID and pwd hash found by given userName or empty string as user ID if user not found
func (db Database) GetPasswordHash(ctx context.Context, userName string) (string, string, error) {
	var userID, hash string

	query := "SELECT id, pwdhash FROM users WHERE username = $1"

	err := db.QueryRowContext(ctx, query, userName).Scan(&userID, &hash)

	if errors.Is(err, sql.ErrNoRows) {
		return "", "", nil
	} else if err != nil {
		return "", "", fmt.Errorf("failed to get a password hash from Database: %w", err)
	}

	return userID, hash, nil
}

func (db *Database) Count(ctx context.Context, dataType int, userID string) (int, error) {
	stmt, err := db.prepCountStmnt(dataType)
	if err != nil {
		return 0,
			fmt.Errorf(
				"failed to prepare db statement for datatype %v: %w",
				model.GetItemTitle(dataType),
				err)
	}

	var res int

	err = stmt.QueryRowContext(ctx, userID).Scan(&res)
	if err != nil {
		return 0,
			fmt.Errorf("failed to execute db statement for datatype %v: %w",
				model.GetItemTitle(dataType),
				err)
	}

	return res, nil
}

func (db *Database) GetData(ctx context.Context, dataType int) (any, error) {
	stmt, err := db.prepLoadStmnt(dataType)
	if err != nil {
		return nil,
			fmt.Errorf("failed to prepare db statement for datatype %v: %w",
				model.GetItemTitle(dataType),
				err)
	}

	return execLoad(ctx, stmt, dataType)
}

func execLoad(ctx context.Context, stmt *sql.Stmt, dataType int) (any, error) {

	userID, ok := ctx.Value(session.CtxKeyUserID).(string)
	if !ok {
		return nil,
			fmt.Errorf("context is missing user ID")
	}

	rows, err := stmt.QueryContext(ctx, userID)
	if err != nil {
		return nil,
			fmt.Errorf("failed to execute db statement: %w", err)
	}

	defer rows.Close()

	switch dataType {
	case model.KeyCredentials:
		res := make([]model.ItemCredentials, 0)

		for rows.Next() {
			item, err := itemCredsFromRow(rows)
			if err != nil {
				return nil, err
			}

			res = append(res, item)
		}

		if rows.Err() != nil {
			return nil,
				fmt.Errorf("failed to scan values from database result: %w", err)
		}

		return res, nil

	case model.KeyText:
		res := make([]model.ItemText, 0)

		for rows.Next() {
			item, err := itemTextFromRow(rows)
			if err != nil {
				return nil, err
			}

			res = append(res, item)
		}

		if rows.Err() != nil {
			return nil,
				fmt.Errorf("failed to scan values from database result: %w", err)
		}

		return res, nil
	case model.KeyCards:
		res := make([]model.ItemCard, 0)

		for rows.Next() {
			item, err := itemCardFromRow(rows)
			if err != nil {
				return nil, err
			}

			res = append(res, item)
		}

		if rows.Err() != nil {
			return nil,
				fmt.Errorf("failed to scan values from database result: %w", err)
		}

		return res, nil
	case model.KeyBinary:
		res := make([]model.ItemBinary, 0)

		for rows.Next() {
			item, err := itemBinaryFromRow(rows)
			if err != nil {
				return nil, err
			}

			res = append(res, item)
		}

		if rows.Err() != nil {
			return nil,
				fmt.Errorf("failed to scan values from database result: %w", err)
		}

		return res, nil
	}

	return nil, fmt.Errorf("attempted tp load an unknown data type")
}

func itemCredsFromRow(rows *sql.Rows) (model.ItemCredentials, error) {
	var encrypted []byte

	item := model.ItemCredentials{}

	// fields: id, encrypted, name, comment
	err := rows.Scan(&item.ID,
		&encrypted,
		&item.Name,
		&item.Comment)

	if err != nil {
		return model.ItemCredentials{},
			fmt.Errorf("failed to scan values from database result: %w", err)
	}

	item.Credentials, err = decryptCred(encrypted)
	if err != nil {
		return model.ItemCredentials{},
			fmt.Errorf("failed to decrypt credentials result: %w", err)
	}

	return item, nil
}

func decryptCred(encrypted []byte) (model.Credentials, error) {
	var res model.Credentials

	b, err := encryption.Decrypt(encrypted)
	if err != nil {
		return res, nil
	}

	buf := bytes.NewBuffer(b)
	dec := json.NewDecoder(buf)
	err = dec.Decode(&res)

	return res, err
}

func encryptCred(item model.Credentials) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	enc := json.NewEncoder(buf)
	err := enc.Encode(item)
	if err != nil {
		return nil, err
	}

	return encryption.Encrypt(buf.Bytes())
}

func itemTextFromRow(rows *sql.Rows) (model.ItemText, error) {
	item := model.ItemText{}

	// fields: id, text, name, comment
	err := rows.Scan(&item.ID,
		&item.Text,
		&item.Name,
		&item.Comment)

	if err != nil {
		return model.ItemText{},
			fmt.Errorf("failed to scan values from database result: %w", err)
	}

	return item, nil
}

func itemCardFromRow(rows *sql.Rows) (model.ItemCard, error) {
	item := model.ItemCard{}

	// fields: id, number, cvvhash, name, comment
	err := rows.Scan(&item.ID,
		&item.Number,
		&item.CVVHash,
		&item.Name,
		&item.Comment)

	if err != nil {
		return model.ItemCard{},
			fmt.Errorf("failed to scan values from database result: %w", err)
	}

	return item, nil
}

func itemBinaryFromRow(rows *sql.Rows) (model.ItemBinary, error) {
	item := model.ItemBinary{}

	// fields: id, data, extention, size, name, comment
	b := make([]byte, 0)

	err := rows.Scan(&item.ID,
		&b,
		&item.Extention,
		&item.Size,
		&item.Name,
		&item.Comment)

	if err != nil {
		return model.ItemBinary{},
			fmt.Errorf("failed to scan values from database result: %w", err)
	}

	item.Data = base64.StdEncoding.EncodeToString(b)

	return item, nil
}

func (db *Database) AddData(ctx context.Context, dataType int, userID string, data any) error {
	// Create new item ID using UUID
	id := uuid.NewString()
	query := itemInsQuery(dataType)
	args, err := itemInsArgs(dataType, id, userID, data)

	if err != nil {
		return fmt.Errorf("failed to compose args for db query: %w", err)
	}

	rowsAffected, err := db.execInsUpdStatement(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("data addition query failed: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("data addition query failed: no rows were added")
	}

	return nil
}

func itemInsQuery(datatype int) string {
	switch datatype {
	case model.KeyCredentials:
		return insCredentials()
	case model.KeyText:
		return insText()
	case model.KeyBinary:
		return insBinary()
	case model.KeyCards:
		return insCard()
	}

	return ""
}

// itemInsArgs returns slice of arguments that matches
// datatype-specific query
func itemInsArgs(datatype int, id, userID string, data any) ([]any, error) {
	switch datatype {
	case model.KeyCredentials:
		item, err := assertItem[model.ItemCredentials](data)
		if err != nil {
			return nil, err
		}

		return argsCredentials(id, userID, item)
	case model.KeyText:
		item, err := assertItem[model.ItemText](data)
		if err != nil {
			return nil, err
		}

		return argsText(id, userID, item)
	case model.KeyBinary:
		item, err := assertItem[model.ItemBinary](data)
		if err != nil {
			return nil, err
		}

		return argsBinary(id, userID, item)
	case model.KeyCards:
		item, err := assertItem[model.ItemCard](data)
		if err != nil {
			return nil, err
		}

		return argsCard(id, userID, item)
	}

	return nil, nil
}

func assertItem[T model.Item](data any) (T, error) {
	val, ok := data.(T)
	if !ok {
		return val, fmt.Errorf("data type mismatch actual data type")
	}

	return val, nil
}

func (db *Database) UpdateData(ctx context.Context, dataType int, userID string, data any) error {

	return nil
}
