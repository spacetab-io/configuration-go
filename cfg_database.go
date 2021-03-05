package config

import (
	"fmt"
	"net"
	"strings"
	"time"
)

const (
	connMaxLifeTime = 30 * time.Minute
	maxOpenConn     = 4
	defaultSchema   = "public"
	cfgDBErrStr     = "config/database: %s"
)

type Database struct {
	Driver             string        `yaml:"driver"`
	Host               string        `yaml:"host"`
	Port               uint          `yaml:"port"`
	User               string        `yaml:"user"`
	Password           string        `yaml:"password"`
	Schema             string        `yaml:"schema"`
	Name               string        `yaml:"name"`
	SSLMode            string        `yaml:"ssl_mode"`
	Extensions         []string      `yaml:"extensions"`
	MaxIdleConnections int           `yaml:"max_idle_connections"`
	MaxOpenConnections int           `yaml:"max_open_connections"`
	ConnectionLifeTime time.Duration `yaml:"connection_lifetime"`

	SeedsPath       string `yaml:"seeds_path"`
	MigrationsPath  string `yaml:"migrations_path"`
	MigrationsTable string `yaml:"migrations_table"`
	MigrateOnStart  bool   `yaml:"migrate_on_start"`
	Debug           bool   `yaml:"debug"`
}

func (d *Database) ConnectionURL() string {
	return fmt.Sprintf(
		"%s://%s:%s@%s:%d/%s?sslmode=%s",
		d.Driver,
		d.User,
		d.Password,
		d.Host,
		d.Port,
		d.Name,
		d.SSLMode,
	)
}

func (d *Database) DSN() string {
	if d.Driver != "postgres" {
		return d.ConnectionURL()
	}

	return fmt.Sprintf(
		"host=%s port=%d dbname=%s user=%s password=%s sslmode=%s",
		d.Host,
		d.Port,
		d.Name,
		d.User,
		d.Password,
		d.SSLMode,
	)
}

func (d *Database) MigrationDSN() string {
	return fmt.Sprintf("%s&x-migrations-table=%s.schema_migrations", d.DSN(), d.Schema)
}

func (d *Database) Validate() []error {
	var errList = make([]error, 0)

	d.Driver = strings.TrimSpace(d.Driver)
	switch d.Driver {
	case "postgres", "mysql":
		break
	case "":
		errList = append(errList, fmt.Errorf(cfgDBErrStr, "driver is empty"))
	default:
		errList = append(errList, fmt.Errorf(cfgDBErrStr, "driver is unknown. only 'postgres' and 'mysql' are well-known"))
	}

	d.User = strings.TrimSpace(d.User)
	if d.User == "" {
		errList = append(errList, fmt.Errorf(cfgDBErrStr, "user is empty"))
	}
	//
	//d.Password = strings.TrimSpace(d.Password)
	//if d.Password == "" {
	//	errList = append(errList, fmt.Errorf(cfgDBErrStr, "password is empty"))
	//}

	d.Host = strings.TrimSpace(d.Host)
	if d.Host == "" {
		errList = append(errList, fmt.Errorf(cfgDBErrStr, "host is empty"))
	} else if _, err := net.LookupHost(d.Host); err != nil {
		errList = append(errList, fmt.Errorf(fmt.Sprintf(cfgDBErrStr, "host resolve err: %s"), err.Error()))
	}

	if d.Port == 0 || d.Port > 70000 {
		errList = append(errList, fmt.Errorf(cfgDBErrStr, "port is invalid"))
	}

	d.Schema = strings.TrimSpace(d.Schema)
	if d.Schema == "" {
		d.Schema = defaultSchema
	}

	d.Name = strings.TrimSpace(d.Name)
	if d.Name == "" {
		errList = append(errList, fmt.Errorf(cfgDBErrStr, "db name is empty"))
	}

	if d.ConnectionLifeTime == 0 {
		d.ConnectionLifeTime = connMaxLifeTime
	}

	if d.MaxOpenConnections == 0 {
		d.MaxOpenConnections = maxOpenConn
	}

	d.MigrationsPath = strings.TrimSpace(d.MigrationsPath)
	if len(d.MigrationsPath) == 0 {
		errList = append(errList, fmt.Errorf(cfgDBErrStr, "migrations_path is empty"))
	}

	return errList
}
