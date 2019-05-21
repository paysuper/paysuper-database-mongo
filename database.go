package database

import (
	"errors"
	"github.com/globalsign/mgo"
	"github.com/kelseyhightower/envconfig"
	"net/url"
	"sync"
	"time"
)

const (
	connectionScheme     = "mongodb"
	errorSessionNotInit  = "database session not init"
	errorConfigIncorrect = "env variables to config connection is incorrect"
)

type Connection struct {
	Host        string `envconfig:"MONGO_HOST"`
	Database    string `envconfig:"MONGO_DB"`
	User        string `envconfig:"MONGO_USER" default:""`
	Password    string `envconfig:"MONGO_PASSWORD" default:""`
	DialTimeout int64  `envconfig:"MONGO_DIAL_TIMEOUT" default:"10"`

	Dns string `envconfig:"MONGO_DNS"`
}

type Source struct {
	name           string
	connection     *Connection
	session        *mgo.Session
	collections    map[string]*mgo.Collection
	database       *mgo.Database
	repositoriesMu sync.Mutex
}

func (c Connection) String() (s string) {
	var u *url.URL
	var err error

	if c.Dns != "" {
		u, err = url.ParseRequestURI(c.Dns)

		if err != nil {
			return ""
		}
	} else {
		if c.Database == "" {
			return ""
		}

		var userInfo *url.Userinfo

		if c.User != "" {
			if c.Password == "" {
				userInfo = url.User(c.User)
			} else {
				userInfo = url.UserPassword(c.User, c.Password)
			}
		}

		u = &url.URL{
			Scheme: connectionScheme,
			Path:   c.Database,
			Host:   c.Host,
			User:   userInfo,
		}
	}

	return u.String()
}

func NewDatabase() (*Source, error) {
	conn := &Connection{}
	err := envconfig.Process("", conn)

	if err != nil {
		return nil, err
	}

	if conn.Dns == "" && (conn.Host == "" || conn.Database == "") {
		return nil, errors.New(errorConfigIncorrect)
	}

	d := &Source{}
	err = d.Open(conn)

	if err != nil {
		return nil, err
	}

	return d, nil
}

func (s *Source) Open(conn *Connection) error {
	s.connection = conn
	return s.open()
}

func (s *Source) open() error {
	var err error

	u := s.connection.String()
	s.session, err = mgo.DialWithTimeout(u, time.Duration(s.connection.DialTimeout)*time.Second)

	if err != nil {
		return err
	}

	s.session.SetMode(mgo.Monotonic, true)

	s.collections = map[string]*mgo.Collection{}
	s.database = s.session.DB("")

	return nil
}

func (s *Source) Close() {
	if s.session != nil {
		s.session.Close()
	}
}

func (s *Source) Ping() error {
	if s.session == nil {
		return errors.New(errorSessionNotInit)
	}

	return s.session.Ping()
}

func (s *Source) Clone() *Source {
	newSession := s.session.Copy()

	clone := &Source{
		name:        s.name,
		connection:  s.connection,
		session:     newSession,
		database:    newSession.DB(s.database.Name),
		collections: map[string]*mgo.Collection{},
	}

	return clone
}

func (s *Source) Drop() error {
	return s.database.DropDatabase()
}

func (s *Source) Collection(name string) *mgo.Collection {
	s.repositoriesMu.Lock()
	defer s.repositoriesMu.Unlock()

	var col *mgo.Collection
	var ok bool

	if col, ok = s.collections[name]; !ok {
		col = s.Clone().database.C(name)
		s.collections[name] = col
	}

	return col
}
