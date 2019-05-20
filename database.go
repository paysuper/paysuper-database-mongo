package database

import (
	"errors"
	"github.com/globalsign/mgo"
	"github.com/kelseyhightower/envconfig"
	"net/url"
	"sync"
)

const connectionScheme = "mongodb"

type Connection struct {
	Host     string `envconfig:"MONGO_HOST" required:"true"`
	Database string `envconfig:"MONGO_DB" required:"true"`
	User     string `envconfig:"MONGO_USER" default:""`
	Password string `envconfig:"MONGO_PASSWORD" default:""`
}

type Source struct {
	name           string
	connection     Connection
	session        *mgo.Session
	collections    map[string]*mgo.Collection
	database       *mgo.Database
	repositoriesMu sync.Mutex
}

func (c Connection) String() (s string) {
	if c.Database == "" {
		return ""
	}

	vv := url.Values{}

	var userInfo *url.Userinfo

	if c.User != "" {
		if c.Password == "" {
			userInfo = url.User(c.User)
		} else {
			userInfo = url.UserPassword(c.User, c.Password)
		}
	}

	u := url.URL{
		Scheme:   connectionScheme,
		Path:     c.Database,
		Host:     c.Host,
		User:     userInfo,
		RawQuery: vv.Encode(),
	}

	return u.String()
}

func NewDatabase() (*Source, error) {
	conn := Connection{}
	err := envconfig.Process("", conn)

	if err != nil {
		return nil, err
	}

	d := &Source{}

	if err := d.Open(conn); err != nil {
		return nil, err
	}

	return d, nil
}

func (s *Source) Open(conn Connection) error {
	s.connection = conn
	return s.open()
}

func (s *Source) open() error {
	var err error

	s.session, err = mgo.Dial(s.connection.String())

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
		return errors.New("No db session")
	}
	err := s.session.Ping()
	return err
}

func (s *Source) Clone() (*Source, error) {
	newSession := s.session.Copy()

	clone := &Source{
		name:        s.name,
		connection:  s.connection,
		session:     newSession,
		database:    newSession.DB(s.database.Name),
		collections: map[string]*mgo.Collection{},
	}

	return clone, nil
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
		c, err := s.Clone()

		if err != nil {
			col = s.database.C(name)
		} else {
			col = c.database.C(name)
		}

		s.collections[name] = col
	}

	return col
}
