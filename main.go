package main

import (
	"encoding/json"
	"fmt"
	"time"

	"goora_date_test/db"

	"github.com/spf13/viper"
)

type TimeDebug struct {
	CurrentDate          string    `db:"CURRENT_DATE"`
	CurrentTimestamp     string    `db:"CURRENT_TIMESTAMP"`
	Sysdate              string    `db:"SYSDATE"`
	Systimestamp         string    `db:"SYSTIMESTAMP"`
	TimeCurrentDate      time.Time `db:"T_CURRENT_DATE"`
	TimeCurrentTimestamp time.Time `db:"T_CURRENT_TIMESTAMP"`
	TimeSysdate          time.Time `db:"T_SYSDATE"`
	TimeSystimestamp     time.Time `db:"T_SYSTIMESTAMP"`
}

var (
	query string = `
		SELECT
			to_char(CURRENT_DATE, 'YYYY-MM-DD HH24:MI:SS') "CURRENT_DATE",
			to_char(CURRENT_TIMESTAMP, 'YYYY-MM-DD HH24:MI:SS TZR') "CURRENT_TIMESTAMP",
			to_char(SYSDATE, 'YYYY-MM-DD HH24:MI:SS') "SYSDATE",
			to_char(SYSTIMESTAMP, 'YYYY-MM-DD HH24:MI:SS TZR') "SYSTIMESTAMP",
			CURRENT_DATE as T_CURRENT_DATE,
			CURRENT_TIMESTAMP as T_CURRENT_TIMESTAMP,
			SYSDATE as T_SYSDATE,
			SYSTIMESTAMP as T_SYSTIMESTAMP
		 FROM DUAL`
)

func main() {
	dbc, err := db.FromViper(viper.Sub(*dbConfKey))
	if err != nil {
		panic(err)
	}
	defer func() {
		err := dbc.Close()
		if err != nil {
			panic(fmt.Errorf("Failed to close connections: %w", err))
		}
	}()
	dbcx, err := db.SqlxFromViper(viper.Sub(*dbConfKey))
	if err != nil {
		panic(err)
	}
	defer func() {
		err := dbcx.Close()
		if err != nil {
			panic(fmt.Errorf("Failed to close connections: %w", err))
		}
	}()
	gdbc, err := db.GodrorFromViper(viper.Sub(*dbConfKey))
	if err != nil {
		panic(err)
	}
	defer func() {
		err := gdbc.Close()
		if err != nil {
			panic(fmt.Errorf("Failed to close connections: %w", err))
		}
	}()

	fmt.Print("sqlx go-ora:\n\n")
	TestDbcX(dbcx)
	fmt.Print("\ngo-ora:\n\n")
	TestDbc(dbc)
	fmt.Print("\ngodror:\n\n")
	TestDbc(gdbc)
}

func TestDbc(dbc *db.Connection) {
	r := TimeDebug{}
	row := dbc.QueryRow(query)
	err := row.Scan(
		&r.CurrentDate, &r.CurrentTimestamp, &r.Sysdate, &r.Systimestamp,
		&r.TimeCurrentDate, &r.TimeCurrentTimestamp, &r.TimeSysdate, &r.TimeSystimestamp,
	)
	if err != nil {
		fmt.Printf("scan error: %#v", err)
	}
	DumpJson(r)
	DumpTime(r)
}

func TestDbcX(dbc *db.SqlxConnection) {
	data := []TimeDebug{}

	err := dbc.Select(&data, query)
	if err != nil {
		fmt.Printf("Error: %#v", err)
	}

	r := data[0]
	DumpJson(r)
	DumpTime(r)
}

func DumpJson(d any) {
	b, _ := json.MarshalIndent(d, "", "  ")
	fmt.Printf("%s\n", string(b))
}

func DumpTime(r TimeDebug) {
	//fmt.Printf("%16s: %s\n", "CurrentDate", r.TimeCurrentDate.Format(time.UnixDate))
	//fmt.Printf("%16s: %s\n", "CurrentTimestamp", r.TimeCurrentTimestamp.Format(time.UnixDate))
	fmt.Printf("%16s: %s\n", "Sysdate", r.TimeSysdate.Format(time.UnixDate))
	fmt.Printf("%16s: %s\n", "Systimestamp", r.TimeSystimestamp.Format(time.UnixDate))
	fmt.Printf("%16s: %s\n", "Real Time", time.Now().Format(time.UnixDate))
}
