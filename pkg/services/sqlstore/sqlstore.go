package sqlstore

import (
	"github.com/emretiryaki/merkut/pkg/bus"
	"github.com/emretiryaki/merkut/pkg/log"
	"github.com/emretiryaki/merkut/pkg/registry"
	"github.com/emretiryaki/merkut/pkg/services/sqlstore/migrator"
	"github.com/emretiryaki/merkut/pkg/setting"
	"github.com/go-xorm/xorm"
	"github.com/mattn/go-sqlite3"
	"golang.org/x/net/context"
	"time"
)

var (
	x       *xorm.Engine
	dialect migrator.Dialect

	sqlog log.Logger = log.New("sqlstore")
)

const ContextSessionName = "db-session"

func init()  {
	registry.Register(&registry.Descriptor{
		Name:"SqlStore",
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

func (ss *SqlStore) WithDbSession(ctx context.Context,callback dbTransactionFunc)  error {

	sess, err := startSession(ctx,ss.engine,false)
	if err != nil {
		return err
	}
	return callback(sess)

}

func (ss *SqlStore) WithTransactionalDbSession(ctx context.Context, callback dbTransactionFunc) error {
	return ss.inTransactionWithRetryCtx(ctx, callback, 0)
}

func (ss *SqlStore) inTransactionWithRetryCtx (ctx context.Context, callback dbTransactionFunc, retry int) error {

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
