package database

import (
	"context"
	"errors"
	"github.com/kelseyhightower/envconfig"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"net/url"
	"sync"
	"time"
)

var (
	ErrorSessionNotInit   = errors.New("database session not init")
	ErrorDatabaseNotFound = errors.New("database name not found in DSN connection string")
)

type SourceInterface interface {
	Close() error
	Ping() error
	Drop() error
	Collection(name string) CollectionInterface
}

type Options struct {
	Dsn     string `envconfig:"MONGO_DSN" default:"mongodb://localhost:27017/test"`
	Mode    string `envconfig:"MONGO_MODE" default:"primary"`
	Context context.Context
}

type Option func(*Options)

type Source struct {
	name           string
	connection     *Options
	repositoriesMu sync.Mutex

	client      *mongo.Client
	database    *mongo.Database
	collections map[string]*Collection
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

func Mode(mode string) Option {
	return func(opts *Options) {
		opts.Mode = mode
	}
}

func Context(ctx context.Context) Option {
	return func(opts *Options) {
		opts.Context = ctx
	}
}

func NewDatabase(options ...Option) (SourceInterface, error) {
	opts := Options{}
	conn := &Options{}

	for _, opt := range options {
		opt(&opts)
	}

	if opts.Dsn == "" || opts.Mode == "" {
		err := envconfig.Process("", conn)

		if err != nil {
			return nil, err
		}
	}

	if opts.Dsn != "" {
		conn.Dsn = opts.Dsn
	}

	if opts.Mode != "" {
		conn.Mode = opts.Mode
	}

	if opts.Context == nil {
		opts.Context, _ = context.WithTimeout(context.Background(), 5*time.Second)
	}

	conn.Context = opts.Context

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
	dsn, err := NewDSN(s.connection.String())

	if err != nil {
		return err
	}

	if dsn.Database == "" {
		return ErrorDatabaseNotFound
	}

	mode, err := readpref.ModeFromString(s.connection.Mode)

	if err != nil {
		return err
	}

	readPref, err := readpref.New(mode)

	if err != nil {
		return err
	}

	opts := options.Client().
		ApplyURI(dsn.Dsn).
		SetReadPreference(readPref)
	s.client, err = mongo.Connect(s.connection.Context, opts)

	if err != nil {
		return err
	}

	s.collections = make(map[string]*Collection)
	s.database = s.client.Database(dsn.Database)
	return nil
}

func (s *Source) Close() error {
	if s.client != nil {
		return s.client.Disconnect(s.connection.Context)
	}

	return nil
}

func (s *Source) Ping() error {
	if s.client == nil {
		return ErrorSessionNotInit
	}

	return s.client.Ping(s.connection.Context, readpref.Primary())
}

func (s *Source) Drop() error {
	return s.database.Drop(s.connection.Context)
}

func (s *Source) Collection(name string) CollectionInterface {
	s.repositoriesMu.Lock()
	col, ok := s.collections[name]

	if !ok {
		col = &Collection{
			collection: s.database.Collection(name),
		}
		s.collections[name] = col
	}
	s.repositoriesMu.Unlock()
	return col
}
