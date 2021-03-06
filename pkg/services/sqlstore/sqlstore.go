package sqlstore

import (
	"fmt"

	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"golang.org/x/net/context"
	"github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/mattn/go-sqlite3"

	"github.com/emretiryaki/merkut/pkg/bus"
	"github.com/emretiryaki/merkut/pkg/log"
	"github.com/emretiryaki/merkut/pkg/registry"
	"github.com/emretiryaki/merkut/pkg/services/sqlstore/migrator"
	"github.com/emretiryaki/merkut/pkg/setting"
	"github.com/emretiryaki/merkut/pkg/services/sqlstore/migrations"
	"github.com/emretiryaki/merkut/pkg/services/sqlstore/sqlutil"


)

var (
	x       *xorm.Engine
	dialect migrator.Dialect

	sqlog log.Logger = log.New("sqlstore")
)

const ContextSessionName = "db-session"

func init() {
	registry.Register(&registry.Descriptor{
		Name:         "SqlStore",
		Instance:     &SqlStore{},
		InitPriority: registry.High,
	})
}

type SqlStore struct {
	Cfg *setting.Cfg `inject:""`
	Bus bus.Bus      `inject:""`

	dbCfg           DatabaseConfig
	engine          *xorm.Engine
	log             log.Logger
	skipEnsureAdmin bool
}

// NewSession returns a new DBSession
func (ss *SqlStore) NewSession() *DBSession {
	return &DBSession{Session: ss.engine.NewSession()}
}

// WithDbSession calls the callback with an session attached to the context.
func (ss *SqlStore) WithDbSession(ctx context.Context, callback dbTransactionFunc) error {
	sess, err := startSession(ctx, ss.engine, false)
	if err != nil {
		return err
	}

	return callback(sess)
}

// WithTransactionalDbSession calls the callback with an session within a transaction
func (ss *SqlStore) WithTransactionalDbSession(ctx context.Context, callback dbTransactionFunc) error {
	return ss.inTransactionWithRetryCtx(ctx, callback, 0)
}

func (ss *SqlStore) inTransactionWithRetryCtx(ctx context.Context, callback dbTransactionFunc, retry int) error {
	sess, err := startSession(ctx, ss.engine, true)
	if err != nil {
		return err
	}

	defer sess.Close()

	err = callback(sess)

	// special handling of database locked errors for sqlite, then we can retry 3 times
	if sqlError, ok := err.(sqlite3.Error); ok && retry < 5 {
		if sqlError.Code == sqlite3.ErrLocked {
			sess.Rollback()
			time.Sleep(time.Millisecond * time.Duration(10))
			sqlog.Info("Database table locked, sleeping then retrying", "retry", retry)
			return ss.inTransactionWithRetryCtx(ctx, callback, retry+1)
		}
	}

	if err != nil {
		sess.Rollback()
		return err
	} else if err = sess.Commit(); err != nil {
		return err
	}

	if len(sess.events) > 0 {
		for _, e := range sess.events {
			if err = bus.Publish(e); err != nil {
				log.Error(3, "Failed to publish event after commit", err)
			}
		}
	}

	return nil
}

func (ss *SqlStore) Init() error {
	ss.log = log.New("sqlstore")
	ss.readConfig()

	engine, err := ss.getEngine()

	if err != nil {
		return fmt.Errorf("Fail to connect to database: %v", err)
	}

	ss.engine = engine

	// temporarily still set global var
	x = engine
	dialect = migrator.NewDialect(x)
	migrator := migrator.NewMigrator(x)
	migrations.AddMigrations(migrator)

	for _, descriptor := range registry.GetServices() {
		sc, ok := descriptor.Instance.(registry.DatabaseMigrator)
		if ok {
			sc.AddMigration(migrator)
		}
	}

	if err := migrator.Start(); err != nil {
		return fmt.Errorf("Migration failed err: %v", err)
	}

	return nil

}


func (ss *SqlStore) buildConnectionString() (string, error) {
	cnnstr := ss.dbCfg.ConnectionString

	// special case used by integration tests
	if cnnstr != "" {
		return cnnstr, nil
	}

	switch ss.dbCfg.Type {
	case migrator.MYSQL:
		protocol := "tcp"
		if strings.HasPrefix(ss.dbCfg.Host, "/") {
			protocol = "unix"
		}

		cnnstr = fmt.Sprintf("%makeCerts:%s@%s(%s)/%s?collation=utf8mb4_unicode_ci&allowNativePasswords=true",
			ss.dbCfg.User, ss.dbCfg.Pwd, protocol, ss.dbCfg.Host, ss.dbCfg.Name)

		if ss.dbCfg.SslMode == "true" || ss.dbCfg.SslMode == "skip-verify" {
			tlsCert, err := makeCert("custom", ss.dbCfg)
			if err != nil {
				return "", err
			}
			mysql.RegisterTLSConfig("custom", tlsCert)
			cnnstr += "&tls=custom"
		}
	case migrator.POSTGRES:
		var host, port = "127.0.0.1", "5432"
		fields := strings.Split(ss.dbCfg.Host, ":")
		if len(fields) > 0 && len(strings.TrimSpace(fields[0])) > 0 {
			host = fields[0]
		}
		if len(fields) > 1 && len(strings.TrimSpace(fields[1])) > 0 {
			port = fields[1]
		}
		if ss.dbCfg.Pwd == "" {
			ss.dbCfg.Pwd = "''"
		}
		if ss.dbCfg.User == "" {
			ss.dbCfg.User = "''"
		}
		cnnstr = fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=%s sslcert=%s sslkey=%s sslrootcert=%s", ss.dbCfg.User, ss.dbCfg.Pwd, host, port, ss.dbCfg.Name, ss.dbCfg.SslMode, ss.dbCfg.ClientCertPath, ss.dbCfg.ClientKeyPath, ss.dbCfg.CaCertPath)
	case migrator.SQLITE:
		// special case for tests
		if !filepath.IsAbs(ss.dbCfg.Path) {
			ss.dbCfg.Path = filepath.Join(setting.DataPath, ss.dbCfg.Path)
		}
		os.MkdirAll(path.Dir(ss.dbCfg.Path), os.ModePerm)
		cnnstr = "file:" + ss.dbCfg.Path + "?cache=shared&mode=rwc"
	default:
		return "", fmt.Errorf("Unknown database type: %s", ss.dbCfg.Type)
	}

	return cnnstr, nil
}

func (ss *SqlStore) getEngine() (*xorm.Engine, error) {
	connectionString, err := ss.buildConnectionString()

	if err != nil {
		return nil, err
	}

	sqlog.Info("Connecting to DB", "dbtype", ss.dbCfg.Type)
	engine, err := xorm.NewEngine(ss.dbCfg.Type, connectionString)
	if err != nil {
		return nil, err
	}

	engine.SetMaxOpenConns(ss.dbCfg.MaxOpenConn)
	engine.SetMaxIdleConns(ss.dbCfg.MaxIdleConn)
	engine.SetConnMaxLifetime(time.Second * time.Duration(ss.dbCfg.ConnMaxLifetime))

	// configure sql logging
	debugSql := ss.Cfg.Raw.Section("database").Key("log_queries").MustBool(false)
	if !debugSql {
		engine.SetLogger(&xorm.DiscardLogger{})
	} else {
		//engine.SetLogger(NewXormLogger(log.LvlInfo, log.New("sqlstore.xorm"))) TODO Loging
		engine.ShowSQL(true)
		engine.ShowExecTime(true)
	}

	return engine, nil
}

func (ss *SqlStore) readConfig() {
	sec := ss.Cfg.Raw.Section("database")

	cfgURL := sec.Key("url").String()
	if len(cfgURL) != 0 {
		dbURL, _ := url.Parse(cfgURL)
		ss.dbCfg.Type = dbURL.Scheme
		ss.dbCfg.Host = dbURL.Host

		pathSplit := strings.Split(dbURL.Path, "/")
		if len(pathSplit) > 1 {
			ss.dbCfg.Name = pathSplit[1]
		}

		userInfo := dbURL.User
		if userInfo != nil {
			ss.dbCfg.User = userInfo.Username()
			ss.dbCfg.Pwd, _ = userInfo.Password()
		}
	} else {
		ss.dbCfg.Type = sec.Key("type").String()
		ss.dbCfg.Host = sec.Key("host").String()
		ss.dbCfg.Name = sec.Key("name").String()
		ss.dbCfg.User = sec.Key("user").String()
		ss.dbCfg.ConnectionString = sec.Key("connection_string").String()
		ss.dbCfg.Pwd = sec.Key("password").String()
	}

	ss.dbCfg.MaxOpenConn = sec.Key("max_open_conn").MustInt(0)
	ss.dbCfg.MaxIdleConn = sec.Key("max_idle_conn").MustInt(2)
	ss.dbCfg.ConnMaxLifetime = sec.Key("conn_max_lifetime").MustInt(14400)

	ss.dbCfg.SslMode = sec.Key("ssl_mode").String()
	ss.dbCfg.CaCertPath = sec.Key("ca_cert_path").String()
	ss.dbCfg.ClientKeyPath = sec.Key("client_key_path").String()
	ss.dbCfg.ClientCertPath = sec.Key("client_cert_path").String()
	ss.dbCfg.ServerCertName = sec.Key("server_cert_name").String()
	ss.dbCfg.Path = sec.Key("path").MustString("data/merkut.db")
}

func InitTestDB(t *testing.T) *SqlStore {
	t.Helper()
	sqlstore := &SqlStore{}
	sqlstore.skipEnsureAdmin = true
	sqlstore.Bus = bus.New()

	dbType := migrator.SQLITE

	// environment variable present for test db?
	if db, present := os.LookupEnv("MERKUT_TEST_DB"); present {
		dbType = db
	}

	// set test db config
	sqlstore.Cfg = setting.NewCfg()
	sec, _ := sqlstore.Cfg.Raw.NewSection("database")
	sec.NewKey("type", dbType)

	switch dbType {
	case "mysql":
		sec.NewKey("connection_string", sqlutil.TestDB_Mysql.ConnStr)
	case "postgres":
		sec.NewKey("connection_string", sqlutil.TestDB_Postgres.ConnStr)
	default:
		sec.NewKey("connection_string", sqlutil.TestDB_Sqlite3.ConnStr)
	}

	// need to get engine to clean db before we init
	engine, err := xorm.NewEngine(dbType, sec.Key("connection_string").String())
	if err != nil {
		t.Fatalf("Failed to init test database: %v", err)
	}

	dialect = migrator.NewDialect(engine)
	if err := dialect.CleanDB(); err != nil {
		t.Fatalf("Failed to clean test db %v", err)
	}

	if err := sqlstore.Init(); err != nil {
		t.Fatalf("Failed to init test database: %v", err)
	}

	sqlstore.engine.DatabaseTZ = time.UTC
	sqlstore.engine.TZLocation = time.UTC

	return sqlstore
}

func IsTestDbMySql() bool {
	if db, present := os.LookupEnv("MERKUT_TEST_DB"); present {
		return db == migrator.MYSQL
	}

	return false
}

func IsTestDbPostgres() bool {
	if db, present := os.LookupEnv("MERKUT_TEST_DB"); present {
		return db == migrator.POSTGRES
	}

	return false
}

type DatabaseConfig struct {
	Type, Host, Name, User, Pwd, Path, SslMode string
	CaCertPath                                 string
	ClientKeyPath                              string
	ClientCertPath                             string
	ServerCertName                             string
	ConnectionString                           string
	MaxOpenConn                                int
	MaxIdleConn                                int
	ConnMaxLifetime                            int
}