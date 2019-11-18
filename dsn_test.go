package database

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type DsnTestSuite struct {
	suite.Suite
}

func Test_Dsn(t *testing.T) {
	suite.Run(t, new(DsnTestSuite))
}

func (suite *DsnTestSuite) SetupTest() {}

func (suite *DsnTestSuite) TearDownTest() {}

func (suite *DsnTestSuite) TestDsn_NewMongoDBDsnManyHosts_Ok() {
	source := "mongodb://admin:123456@host1:27017,host2:27018,host3:27019/test?readPreference=primary&ssl=true&w=majority"
	dsn, err := NewDSN(source)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), dsn)
	assert.Equal(suite.T(), source, dsn.Dsn)
	assert.Equal(suite.T(), "mongodb", dsn.Protocol)
	assert.NotNil(suite.T(), dsn.Auth)
	assert.Equal(suite.T(), "admin", dsn.Auth.User)
	assert.Equal(suite.T(), "123456", dsn.Auth.Password)
	assert.Len(suite.T(), dsn.Hosts, 3)
	assert.Equal(suite.T(), "host1", dsn.Hosts[0].Host)
	assert.Equal(suite.T(), "27017", dsn.Hosts[0].Port)
	assert.Equal(suite.T(), "host2", dsn.Hosts[1].Host)
	assert.Equal(suite.T(), "27018", dsn.Hosts[1].Port)
	assert.Equal(suite.T(), "host3", dsn.Hosts[2].Host)
	assert.Equal(suite.T(), "27019", dsn.Hosts[2].Port)
	assert.Equal(suite.T(), "test", dsn.Database)
	assert.NotEmpty(suite.T(), dsn.Options)
	assert.Contains(suite.T(), dsn.Options, "readPreference")
	assert.Contains(suite.T(), dsn.Options, "ssl")
	assert.Contains(suite.T(), dsn.Options, "w")
	assert.Equal(suite.T(), "primary", dsn.Options["readPreference"])
	assert.Equal(suite.T(), "true", dsn.Options["ssl"])
	assert.Equal(suite.T(), "majority", dsn.Options["w"])
}

func (suite *DsnTestSuite) TestDsn_NewMongoDBDsnOneHost_Ok() {
	source := "mongodb://admin:123456@host1:27017/test?readPreference=primary&ssl=true&w=majority"
	dsn, err := NewDSN(source)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), dsn)
	assert.Equal(suite.T(), source, dsn.Dsn)
	assert.Equal(suite.T(), "mongodb", dsn.Protocol)
	assert.NotNil(suite.T(), dsn.Auth)
	assert.Equal(suite.T(), "admin", dsn.Auth.User)
	assert.Equal(suite.T(), "123456", dsn.Auth.Password)
	assert.Len(suite.T(), dsn.Hosts, 1)
	assert.Equal(suite.T(), "host1", dsn.Hosts[0].Host)
	assert.Equal(suite.T(), "27017", dsn.Hosts[0].Port)
	assert.Equal(suite.T(), "test", dsn.Database)
	assert.NotEmpty(suite.T(), dsn.Options)
	assert.Contains(suite.T(), dsn.Options, "readPreference")
	assert.Contains(suite.T(), dsn.Options, "ssl")
	assert.Contains(suite.T(), dsn.Options, "w")
	assert.Equal(suite.T(), "primary", dsn.Options["readPreference"])
	assert.Equal(suite.T(), "true", dsn.Options["ssl"])
	assert.Equal(suite.T(), "majority", dsn.Options["w"])
}

func (suite *DsnTestSuite) TestDsn_NewPostgreSQLDsn_Ok() {
	source := "postgres://postgres:postgres123456@localhost:5432/pg_test?sslmode=disable"
	dsn, err := NewDSN(source)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), dsn)
	assert.Equal(suite.T(), source, dsn.Dsn)
	assert.Equal(suite.T(), "postgres", dsn.Protocol)
	assert.NotNil(suite.T(), dsn.Auth)
	assert.Equal(suite.T(), "postgres", dsn.Auth.User)
	assert.Equal(suite.T(), "postgres123456", dsn.Auth.Password)
	assert.Len(suite.T(), dsn.Hosts, 1)
	assert.Equal(suite.T(), "localhost", dsn.Hosts[0].Host)
	assert.Equal(suite.T(), "5432", dsn.Hosts[0].Port)
	assert.Equal(suite.T(), "pg_test", dsn.Database)
	assert.Len(suite.T(), dsn.Options, 1)
	assert.Contains(suite.T(), dsn.Options, "sslmode")
	assert.Equal(suite.T(), "disable", dsn.Options["sslmode"])
}

func (suite *DsnTestSuite) TestDsn_NewDsn_ProtocolNotFound() {
	dsn, err := NewDSN("admin:123456@host1:27017/test?readPreference=primary&ssl=true&w=majority")
	assert.Equal(suite.T(), err, ErrorProtocolNotFound)
	assert.Nil(suite.T(), dsn)
}

func (suite *DsnTestSuite) TestDsn_NewDsn_AuthWithoutPassword_Ok() {
	dsn, err := NewDSN("mongodb://admin:@host1:27017/test?readPreference=primary&ssl=true&w=majority")
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), dsn)
	assert.Equal(suite.T(), "mongodb", dsn.Protocol)
	assert.NotNil(suite.T(), dsn.Auth)
	assert.Equal(suite.T(), "admin", dsn.Auth.User)
	assert.Equal(suite.T(), "", dsn.Auth.Password)
	assert.Len(suite.T(), dsn.Hosts, 1)
	assert.Equal(suite.T(), "host1", dsn.Hosts[0].Host)
	assert.Equal(suite.T(), "27017", dsn.Hosts[0].Port)
	assert.Equal(suite.T(), "test", dsn.Database)
	assert.NotEmpty(suite.T(), dsn.Options)
	assert.Contains(suite.T(), dsn.Options, "readPreference")
	assert.Contains(suite.T(), dsn.Options, "ssl")
	assert.Contains(suite.T(), dsn.Options, "w")
	assert.Equal(suite.T(), "primary", dsn.Options["readPreference"])
	assert.Equal(suite.T(), "true", dsn.Options["ssl"])
	assert.Equal(suite.T(), "majority", dsn.Options["w"])
}

func (suite *DsnTestSuite) TestDsn_NewDsn_AuthWithoutPassword_2_Ok() {
	dsn, err := NewDSN("mongodb://admin@host1:27017/test?readPreference=primary&ssl=true&w=majority")
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), dsn)
	assert.Equal(suite.T(), "mongodb", dsn.Protocol)
	assert.NotNil(suite.T(), dsn.Auth)
	assert.Equal(suite.T(), "admin", dsn.Auth.User)
	assert.Equal(suite.T(), "", dsn.Auth.Password)
	assert.Len(suite.T(), dsn.Hosts, 1)
	assert.Equal(suite.T(), "host1", dsn.Hosts[0].Host)
	assert.Equal(suite.T(), "27017", dsn.Hosts[0].Port)
	assert.Equal(suite.T(), "test", dsn.Database)
	assert.NotEmpty(suite.T(), dsn.Options)
	assert.Contains(suite.T(), dsn.Options, "readPreference")
	assert.Contains(suite.T(), dsn.Options, "ssl")
	assert.Contains(suite.T(), dsn.Options, "w")
	assert.Equal(suite.T(), "primary", dsn.Options["readPreference"])
	assert.Equal(suite.T(), "true", dsn.Options["ssl"])
	assert.Equal(suite.T(), "majority", dsn.Options["w"])
}

func (suite *DsnTestSuite) TestDsn_NewDsn_HostsNotFound_Error() {
	dsn, err := NewDSN("mongodb://admin@/test?readPreference=primary&ssl=true&w=majority")
	assert.Equal(suite.T(), err, ErrorHostsNotFound)
	assert.Nil(suite.T(), dsn)
}

func (suite *DsnTestSuite) TestDsn_NewDsn_DsnWithoutDatabase_Ok() {
	dsn, err := NewDSN("mongodb://admin:123456@host1:27017?readPreference=primary&ssl=true&w=majority")
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), dsn)
	assert.Equal(suite.T(), "mongodb", dsn.Protocol)
	assert.NotNil(suite.T(), dsn.Auth)
	assert.Equal(suite.T(), "admin", dsn.Auth.User)
	assert.Equal(suite.T(), "123456", dsn.Auth.Password)
	assert.Len(suite.T(), dsn.Hosts, 1)
	assert.Equal(suite.T(), "host1", dsn.Hosts[0].Host)
	assert.Equal(suite.T(), "27017", dsn.Hosts[0].Port)
	assert.Zero(suite.T(), dsn.Database)
	assert.NotEmpty(suite.T(), dsn.Options)
	assert.Contains(suite.T(), dsn.Options, "readPreference")
	assert.Contains(suite.T(), dsn.Options, "ssl")
	assert.Contains(suite.T(), dsn.Options, "w")
	assert.Equal(suite.T(), "primary", dsn.Options["readPreference"])
	assert.Equal(suite.T(), "true", dsn.Options["ssl"])
	assert.Equal(suite.T(), "majority", dsn.Options["w"])
}

func (suite *DsnTestSuite) TestDsn_NewDsn_DsnWithoutDatabase_2_Ok() {
	dsn, err := NewDSN("mongodb://admin:123456@host1:27017/?readPreference=primary&ssl=true&w=majority")
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), dsn)
	assert.Equal(suite.T(), "mongodb", dsn.Protocol)
	assert.NotNil(suite.T(), dsn.Auth)
	assert.Equal(suite.T(), "admin", dsn.Auth.User)
	assert.Equal(suite.T(), "123456", dsn.Auth.Password)
	assert.Len(suite.T(), dsn.Hosts, 1)
	assert.Equal(suite.T(), "host1", dsn.Hosts[0].Host)
	assert.Equal(suite.T(), "27017", dsn.Hosts[0].Port)
	assert.Zero(suite.T(), dsn.Database)
	assert.NotEmpty(suite.T(), dsn.Options)
	assert.Contains(suite.T(), dsn.Options, "readPreference")
	assert.Contains(suite.T(), dsn.Options, "ssl")
	assert.Contains(suite.T(), dsn.Options, "w")
	assert.Equal(suite.T(), "primary", dsn.Options["readPreference"])
	assert.Equal(suite.T(), "true", dsn.Options["ssl"])
	assert.Equal(suite.T(), "majority", dsn.Options["w"])
}

func (suite *DsnTestSuite) TestDsn_NewDsn_HostWithoutPort_Ok() {
	dsn, err := NewDSN("mongodb://admin:123456@host1:/test?readPreference=primary&ssl=true&w=majority")
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), dsn)
	assert.Equal(suite.T(), "mongodb", dsn.Protocol)
	assert.NotNil(suite.T(), dsn.Auth)
	assert.Equal(suite.T(), "admin", dsn.Auth.User)
	assert.Equal(suite.T(), "123456", dsn.Auth.Password)
	assert.Len(suite.T(), dsn.Hosts, 1)
	assert.Equal(suite.T(), "host1", dsn.Hosts[0].Host)
	assert.Zero(suite.T(), dsn.Hosts[0].Port)
	assert.Equal(suite.T(), "test", dsn.Database)
	assert.NotEmpty(suite.T(), dsn.Options)
	assert.Contains(suite.T(), dsn.Options, "readPreference")
	assert.Contains(suite.T(), dsn.Options, "ssl")
	assert.Contains(suite.T(), dsn.Options, "w")
	assert.Equal(suite.T(), "primary", dsn.Options["readPreference"])
	assert.Equal(suite.T(), "true", dsn.Options["ssl"])
	assert.Equal(suite.T(), "majority", dsn.Options["w"])
}

func (suite *DsnTestSuite) TestDsn_NewDsn_HostWithoutPort_2_Ok() {
	dsn, err := NewDSN("mongodb://admin:123456@host1/test?readPreference=primary&ssl=true&w=majority")
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), dsn)
	assert.Equal(suite.T(), "mongodb", dsn.Protocol)
	assert.NotNil(suite.T(), dsn.Auth)
	assert.Equal(suite.T(), "admin", dsn.Auth.User)
	assert.Equal(suite.T(), "123456", dsn.Auth.Password)
	assert.Len(suite.T(), dsn.Hosts, 1)
	assert.Equal(suite.T(), "host1", dsn.Hosts[0].Host)
	assert.Zero(suite.T(), dsn.Hosts[0].Port)
	assert.Equal(suite.T(), "test", dsn.Database)
	assert.NotEmpty(suite.T(), dsn.Options)
	assert.Contains(suite.T(), dsn.Options, "readPreference")
	assert.Contains(suite.T(), dsn.Options, "ssl")
	assert.Contains(suite.T(), dsn.Options, "w")
	assert.Equal(suite.T(), "primary", dsn.Options["readPreference"])
	assert.Equal(suite.T(), "true", dsn.Options["ssl"])
	assert.Equal(suite.T(), "majority", dsn.Options["w"])
}

func (suite *DsnTestSuite) TestDsn_NewDsn_OptionWithoutValue_Ok() {
	dsn, err := NewDSN("mongodb://admin:123456@host1/test?readPreference=primary&ssl=true&w=majority&someOptionWithoutValue")
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), dsn)
	assert.Equal(suite.T(), "mongodb", dsn.Protocol)
	assert.NotNil(suite.T(), dsn.Auth)
	assert.Equal(suite.T(), "admin", dsn.Auth.User)
	assert.Equal(suite.T(), "123456", dsn.Auth.Password)
	assert.Len(suite.T(), dsn.Hosts, 1)
	assert.Equal(suite.T(), "host1", dsn.Hosts[0].Host)
	assert.Zero(suite.T(), dsn.Hosts[0].Port)
	assert.Equal(suite.T(), "test", dsn.Database)
	assert.NotEmpty(suite.T(), dsn.Options)
	assert.Contains(suite.T(), dsn.Options, "readPreference")
	assert.Contains(suite.T(), dsn.Options, "ssl")
	assert.Contains(suite.T(), dsn.Options, "w")
	assert.Contains(suite.T(), dsn.Options, "someOptionWithoutValue")
	assert.Equal(suite.T(), "primary", dsn.Options["readPreference"])
	assert.Equal(suite.T(), "true", dsn.Options["ssl"])
	assert.Equal(suite.T(), "majority", dsn.Options["w"])
	assert.Zero(suite.T(), dsn.Options["someOptionWithoutValue"])
}
