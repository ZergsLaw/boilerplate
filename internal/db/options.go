package db

import (
	"strconv"
	"strings"
)

type (
	// Option for building config.
	Option func(*config)
	// SSLMode determines whether or with what priority a secure SSL TCP/IP.
	SSLMode uint8
	// For connect to database.
	config struct {
		dbName string
		dbUser string
		dbPass string
		dbHost string
		dbPort int
	}
)

// FormatDSN returns dataSourceName string with properly escaped
// connection parameters suitable for sql.Open.
func (o config) FormatDSN() string {
	// Borrowed from pq.ParseURL.
	var kvs []string

	replacer := strings.NewReplacer(` `, `\ `, `'`, `\'`, `\`, `\\`)
	accrue := func(k, v string) {
		if v != "" {
			kvs = append(kvs, k+"="+replacer.Replace(v))
		}
	}

	accrue("dbname", o.dbName)
	accrue("user", o.dbUser)
	accrue("password", o.dbPass)
	accrue("host", o.dbHost)
	accrue("port", strconv.Itoa(o.dbPort))

	return strings.Join(kvs, " ")
}

func defaultConfig() *config {
	return &config{
		dbName: "postgres",
		dbUser: "postgres",
		dbPass: "postgres",
		dbHost: "localhost",
		dbPort: 5432,
	}
}

// Name sets the connection parameters.
func Name(name string) Option {
	return func(config *config) {
		config.dbName = name
	}
}

// User sets the connection parameters.
func User(user string) Option {
	return func(config *config) {
		config.dbUser = user
	}
}

// Pass sets the connection parameters.
func Pass(pass string) Option {
	return func(config *config) {
		config.dbPass = pass
	}
}

// Host sets the connection parameters.
func Host(host string) Option {
	return func(config *config) {
		config.dbHost = host
	}
}

// Port sets the connection parameters.
func Port(port int) Option {
	return func(config *config) {
		config.dbPort = port
	}
}
