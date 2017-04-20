package database

import (
	mgo "gopkg.in/mgo.v2"
)

const (
	DBURL string = "localhost"
)

type DBSession struct {
	Session *mgo.Session
}

func InitDatabase() *mgo.Session {
	//DBSession * session
	session, err := mgo.Dial(DBURL)
	if err != nil {
		panic(err)
	}
	session.SetMode(mgo.Monotonic, true)
	return session
}

func CloseDatabase(session *mgo.Session) {
	session.Close()
}
