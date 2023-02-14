package psqldb

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"strings"

	"ghostorange/internal/app/model"
)

func (db *Database) prepCountStmnt(dataType int) (*sql.Stmt, error) {
	var table string

	switch dataType {
	case model.KeyCredentials:
		table = "credentials"
	case model.KeyText:
		table = "text"
	case model.KeyBinary:
		table = "binarydata"
	case model.KeyCards:
		table = "cards"
	}

	query := fmt.Sprintf("SELECT COUNT(id) FROM %v WHERE user_id = $1",
	 table)

	return db.Prepare(query)
}

func (db *Database) prepLoadStmnt(dataType int) (*sql.Stmt, error) {
	var query string

	switch dataType {
	case model.KeyCredentials:
		query = selCredentials()
	case model.KeyText:
		query = selText()
	case model.KeyBinary:
		query = selBinary()
	case model.KeyCards:
		query = selCards()
	}

	return db.Prepare(query)
}

func selCredentials() string {
	return `SELECT id, encrypted, name, comment
		FROM credentials
		WHERE user_id = $1`
}

func selText() string {
	return `SELECT id, text, name, comment
		FROM text
		WHERE user_id = $1`
}

func selBinary() string {
	return `SELECT id, data, extention, size, name, comment
		FROM binarydata
		WHERE user_id = $1`
}

func selCards() string {
	return `SELECT id, number, name, comment
		FROM cards
		WHERE user_id = $1`
}

func insCredentials()string{
	return `INSERT INTO credentials(
		id, user_id, ts, encrypted, name, comment
		) 
		VALUES (
			$1, $2, now()::timestamptz, $3, $4, $5
			) 
			ON CONFLICT (id) DO UPDATE SET
			encrypted=$3, 
			name=$4, 
			comment=$5`
}

// argsCredentials returns slice of args required 
// by query. See insCredentials.
func argsCredentials(id, userID string, item model.ItemCredentials)([]any,error){
	encrypted, err := encryptCred(item.Credentials)
	if err != nil{
		return nil, err
	}

	if item.ID != ""{
		id = item.ID
	}
	
	return []any{
		id,
		userID,
		encrypted,
		item.Name,
		item.Comment,
	}, nil
}

func insText()string{
	return `INSERT INTO text(
		id, user_id, ts, text, name, comment
		) 
		VALUES (
			$1, $2, now()::timestamptz, $3, $4, $5
			) 
			ON CONFLICT (id) DO UPDATE SET
			text=$3, 
			name=$4, 
			comment=$5`
}

// argsText returns slice of args required 
// by query. See insText.
func argsText(id, userID string, item model.ItemText)([]any,error){
	if item.ID != ""{
		id = item.ID
	}

	return []any{
		id,
		userID,
		[]byte(item.Text),
		item.Name,
		item.Comment,
	}, nil
}

func insCard()string{
	return `INSERT INTO cards(
		id, user_id, ts, number, full_number, cvvhash, expires, 
		name, comment,
		cardholdername,cardholdersurename
		) 
		VALUES (
			$1, $2, now()::timestamptz, $3, $4, $5, $6, $7, $8, $9, $10
			) 
			ON CONFLICT (id) DO UPDATE SET
			number = $3, 
			full_number = $4, 
			cvvhash=$5, 
			expires=$6, 
			name=$7, 
			comment=$8,
			cardholdername=$9,
			cardholdersurename=$10`

}

// argsCard returns slice of args required 
// by query. See insCard.
func argsCard(id, userID string, item model.ItemCard)([]any,error){
	// Replace 8 middle charachters with *
	number := item.Number[:4] + strings.Repeat("*", 8) + item.Number[12:]

	if item.ID != ""{
		id = item.ID
	}

	return []any{
		id,
		userID,
		number,
		item.Number,
		item.CVVHash,
		item.Exp,
		item.Name,
		item.Comment,
		item.CardholderName,
		item.CardholderSurename,
	}, nil
}

func insBinary()string{
	return `INSERT INTO binarydata(
		id, user_id, ts, data, extention, size, name, comment
		) 
		VALUES (
			$1, $2, now()::timestamptz, $3, $4, $5, $6, $7
			) 
			ON CONFLICT (id) DO UPDATE SET
			data=$3, 
			extention=$4, 
			size=$5, 
			name=$6, 
			comment=$7`
}

// argsBinary returns slice of args required 
// by query. See insBinary.
func argsBinary(id, userID string, item model.ItemBinary)([]any,error){
	data,err := base64.StdEncoding.DecodeString(item.Data)
	if err !=nil{
		return nil,err
	}

	if item.ID != ""{
		id = item.ID
	}
	
	return []any{
		id,
		userID,
		data,
		item.Extention,
		item.Size,
		item.Name,
		item.Comment,
	}, nil
}
