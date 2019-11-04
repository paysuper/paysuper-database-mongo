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
	errorSessionNotInit = "database session not init"
)

type Options struct {
	Dsn         string   `envconfig:"MONGO_DSN" default:"mongodb://localhost:27017/test"`
	DialTimeout int64    `envconfig:"MONGO_DIAL_TIMEOUT" default:"10"`
	Mode        mgo.Mode `envconfig:"MONGO_MODE" default:"1"`

	modeChanged bool
}

type Option func(*Options)

type Source struct {
	name           string
	connection     *Options
	session        *mgo.Session
	collections    map[string]*mgo.Collection
	database       *mgo.Database
	repositoriesMu sync.Mutex
}

func (c Options) String() (s string) {
	var u *url.URL
	var err error

	u, err = url.ParseRequestURI(c.Dsn)

	if err != nil {
		return ""
	}

	return u.String()
}

func Dsn(dsn string) Option {
	return func(opts *Options) {
		opts.Dsn = dsn
	}
}

func DialTimeout(t int64) Option {
	return func(opts *Options) {
		opts.DialTimeout = t
	}
}

func Mode(t mgo.Mode) Option {
	return func(opts *Options) {
		opts.Mode = t
		opts.modeChanged = true
	}
}

func NewDatabase(options ...Option) (*Source, error) {
	opts := Options{}
	conn := &Options{}

	for _, opt := range options {
		opt(&opts)
	}

	if opts.Dsn == "" || opts.DialTimeout == 0 {
		err := envconfig.Process("", conn)

		if err != nil {
			return nil, err
		}
	}

	if opts.Dsn != "" {
		conn.Dsn = opts.Dsn
	}

	if opts.DialTimeout > 0 {
		conn.DialTimeout = opts.DialTimeout
	}

	if opts.modeChanged {
		conn.Mode = opts.Mode
	}

	d := &Source{}
	err := d.Open(conn)

	if err != nil {
		return nil, err
	}

	return d, nil
}

func (s *Source) Open(conn *Options) error {
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

	s.session.SetMode(s.connection.Mode, true)

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
