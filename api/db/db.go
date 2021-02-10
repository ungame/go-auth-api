package db

import mgo "gopkg.in/mgo.v2"

type Connection interface {
	Close()
	DB() *mgo.Database
}

type conn struct {
	session  *mgo.Session
	database *mgo.Database
}

func NewConnection(cfg Config) (Connection, error) {
	session, err := mgo.Dial(cfg.Dsn())
	if err != nil {
		return nil, err
	}
	database := session.DB(cfg.DbName())
	return &conn{session, database}, nil
}

func (c *conn) Close() {
	c.session.Close()
}

func (c *conn) DB() *mgo.Database {
	return c.database
}
