package db

import (
	"strings"

	go_ora "github.com/sijms/go-ora/v2"
)

type ConnectArgs struct {
	Username         string
	Password         string
	PasswordCommand  string
	Server           string
	Port             int
	Service          string
	SID              string
	ConnectionString string
	Opts             map[string]interface{}
}

func (args *ConnectArgs) String() string {
	cleanArgs := *args
	cleanArgs.Password = "XXXXXXXX"
	return (&cleanArgs).ToConnectionString()
}

func (args *ConnectArgs) ToConnectionString() string {
	urloptions := make(map[string]string)
	if args.SID != "" {
		urloptions["SID"] = args.SID
	}
	for key, val := range args.Opts {
		urloptions[strings.ToUpper(key)] = val.(string)
	}
	url := go_ora.BuildUrl(args.Server, args.Port, args.Service, args.Username, args.Password, urloptions)
	return url
}
