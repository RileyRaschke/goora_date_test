package db

import (
	"database/sql"
	"fmt"

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
	var err error
	if args.Password == "" {
		if args.PasswordCommand != "" {
			args.Password, err = passwordFromCommand(args.PasswordCommand)
			if err != nil {
				return nil, err
			}
		}
	}
	if args.Password == "" {
		args.Password, err = passwordFromShell()
		if err != nil {
			return nil, err
		}
	}

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
	//db.MustExec("ALTER SESSION SET TIME_ZONE = 'America/Chicago'")
	db.Exec("ALTER SESSION SET TIME_ZONE = 'America/Chicago'")
	return db, nil
}

func SqlxConnect(args *ConnectArgs) (*SqlxConnection, error) {
	sqlx.BindDriver(DB_TYPE, sqlx.NAMED)
	var err error
	if args.Password == "" {
		if args.PasswordCommand != "" {
			args.Password, err = passwordFromCommand(args.PasswordCommand)
			if err != nil {
				return nil, err
			}
		}
	}
	if args.Password == "" {
		args.Password, err = passwordFromShell()
		if err != nil {
			return nil, err
		}
	}

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
	db.MustExec("ALTER SESSION SET TIME_ZONE = 'America/Chicago'")
	return db, nil
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
