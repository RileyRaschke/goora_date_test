package db

import (
	"fmt"

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
	args.GetPass()
	urloptions := make(map[string]string)
	if args.SID != "" {
		urloptions["SID"] = args.SID
	}
	for key, val := range args.Opts {
		urloptions[key] = val.(string)
	}
	url := go_ora.BuildUrl(args.Server, args.Port, args.Service, args.Username, args.Password, urloptions)
	return url
}

func (args *ConnectArgs) GodrorConnectString() string {
	args.GetPass()
	opts := ""
	for key, val := range args.Opts {
		switch val.(type) {
		case int:
			opts = opts + fmt.Sprintf(`%s=%d`, key, val.(int))
			break
		case bool:
			opts = opts + fmt.Sprintf(`%s=%b`, key, val.(string))
			break
		default: // string
			opts = opts + fmt.Sprintf(`%s="%s"`, key, val.(string))
		}
	}
	connStr := fmt.Sprintf(`user="%s" password="%s" connectString="%s:%d/%s" %s`,
		args.Username,
		args.Password,
		args.Server,
		args.Port,
		args.Service,
		opts,
	)
	return connStr
}

func (args *ConnectArgs) GetPass() {
	var err error
	if args.Password == "" {
		if args.PasswordCommand != "" {
			args.Password, err = passwordFromCommand(args.PasswordCommand)
			if err != nil {
				panic(err)
			}
		}
	}
	if args.Password == "" {
		args.Password, err = passwordFromShell()
		if err != nil {
			panic(err)
		}
	}
}
