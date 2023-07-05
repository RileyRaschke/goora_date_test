package db

import (
	"database/sql"
	"fmt"

	_ "github.com/godror/godror"
	"github.com/jmoiron/sqlx"
	_ "github.com/sijms/go-ora/v2"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const DB_TYPE = "oracle"

type Connection struct {
	*sql.DB
}

type SqlxConnection struct {
	*sqlx.DB
}

func Connect(args *ConnectArgs) (*Connection, error) {
	c, err := sql.Open(DB_TYPE, args.ToConnectionString())
	log.Tracef("%s", args)
	if err != nil {
		panic(fmt.Errorf("error in sql.Open: %w", err))
	}
	db := &Connection{c}
	err = db.Ping()
	if err != nil {
		log.Fatalf("Error pinging db: %w", err)
	}
	//db.Exec("ALTER SESSION SET TIME_ZONE = 'America/Chicago'")
	return db, nil
}

func SqlxConnect(args *ConnectArgs) (*SqlxConnection, error) {
	sqlx.BindDriver(DB_TYPE, sqlx.NAMED)
	c, err := sqlx.Open(DB_TYPE, args.ToConnectionString())
	log.Tracef("%s", args)
	if err != nil {
		panic(fmt.Errorf("error in sql.Open: %w", err))
	}
	db := &SqlxConnection{c}
	err = db.Ping()
	if err != nil {
		log.Fatalf("Error pinging db: %w", err)
	}
	//db.MustExec("ALTER SESSION SET TIME_ZONE = 'America/Chicago'")
	return db, nil
}

func GodrorConnect(args *ConnectArgs) (*Connection, error) {
	c, err := sql.Open("godror", args.GodrorConnectString())
	log.Tracef("%s", args)
	if err != nil {
		panic(fmt.Errorf("error in sql.Open: %w", err))
	}
	db := &Connection{c}
	err = db.Ping()
	if err != nil {
		log.Fatalf("Error pinging db: %w", err)
	}
	//db.Exec("ALTER SESSION SET TIME_ZONE = 'America/Chicago'")
	return db, nil
}

func GodrorFromViper(v *viper.Viper) (*Connection, error) {
	return GodrorConnect(
		&ConnectArgs{
			Username:        v.GetString("Username"),
			Password:        v.GetString("Password"),
			PasswordCommand: v.GetString("PasswordCommand"),
			Server:          v.GetString("Server"),
			Port:            v.GetInt("Port"),
			Service:         v.GetString("Service"),
			SID:             v.GetString("SID"),
			Opts:            v.Get("Options").(map[string]any),
		})
}

func FromViper(v *viper.Viper) (*Connection, error) {
	return Connect(
		&ConnectArgs{
			Username:        v.GetString("Username"),
			Password:        v.GetString("Password"),
			PasswordCommand: v.GetString("PasswordCommand"),
			Server:          v.GetString("Server"),
			Port:            v.GetInt("Port"),
			Service:         v.GetString("Service"),
			SID:             v.GetString("SID"),
			Opts:            v.Get("Options").(map[string]any),
		})
}

func SqlxFromViper(v *viper.Viper) (*SqlxConnection, error) {
	return SqlxConnect(
		&ConnectArgs{
			Username:        v.GetString("Username"),
			Password:        v.GetString("Password"),
			PasswordCommand: v.GetString("PasswordCommand"),
			Server:          v.GetString("Server"),
			Port:            v.GetInt("Port"),
			Service:         v.GetString("Service"),
			SID:             v.GetString("SID"),
			Opts:            v.Get("Options").(map[string]any),
		})
}
