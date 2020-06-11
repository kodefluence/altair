// Code generated by MockGen. DO NOT EDIT.
// Source: core/loader.go

// Package mock is a generated GoMock package.
package mock

import (
	sql "database/sql"
	core "github.com/codefluence-x/altair/core"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
	time "time"
)

// MockAppLoader is a mock of AppLoader interface
type MockAppLoader struct {
	ctrl     *gomock.Controller
	recorder *MockAppLoaderMockRecorder
}

// MockAppLoaderMockRecorder is the mock recorder for MockAppLoader
type MockAppLoaderMockRecorder struct {
	mock *MockAppLoader
}

// NewMockAppLoader creates a new mock instance
func NewMockAppLoader(ctrl *gomock.Controller) *MockAppLoader {
	mock := &MockAppLoader{ctrl: ctrl}
	mock.recorder = &MockAppLoaderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockAppLoader) EXPECT() *MockAppLoaderMockRecorder {
	return m.recorder
}

// Compile mocks base method
func (m *MockAppLoader) Compile(configPath string) (core.AppConfig, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Compile", configPath)
	ret0, _ := ret[0].(core.AppConfig)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Compile indicates an expected call of Compile
func (mr *MockAppLoaderMockRecorder) Compile(configPath interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Compile", reflect.TypeOf((*MockAppLoader)(nil).Compile), configPath)
}

// MockDatabaseLoader is a mock of DatabaseLoader interface
type MockDatabaseLoader struct {
	ctrl     *gomock.Controller
	recorder *MockDatabaseLoaderMockRecorder
}

// MockDatabaseLoaderMockRecorder is the mock recorder for MockDatabaseLoader
type MockDatabaseLoaderMockRecorder struct {
	mock *MockDatabaseLoader
}

// NewMockDatabaseLoader creates a new mock instance
func NewMockDatabaseLoader(ctrl *gomock.Controller) *MockDatabaseLoader {
	mock := &MockDatabaseLoader{ctrl: ctrl}
	mock.recorder = &MockDatabaseLoaderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockDatabaseLoader) EXPECT() *MockDatabaseLoaderMockRecorder {
	return m.recorder
}

// Compile mocks base method
func (m *MockDatabaseLoader) Compile(configPath string) (map[string]core.DatabaseConfig, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Compile", configPath)
	ret0, _ := ret[0].(map[string]core.DatabaseConfig)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Compile indicates an expected call of Compile
func (mr *MockDatabaseLoaderMockRecorder) Compile(configPath interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Compile", reflect.TypeOf((*MockDatabaseLoader)(nil).Compile), configPath)
}

// MockPluginLoader is a mock of PluginLoader interface
type MockPluginLoader struct {
	ctrl     *gomock.Controller
	recorder *MockPluginLoaderMockRecorder
}

// MockPluginLoaderMockRecorder is the mock recorder for MockPluginLoader
type MockPluginLoaderMockRecorder struct {
	mock *MockPluginLoader
}

// NewMockPluginLoader creates a new mock instance
func NewMockPluginLoader(ctrl *gomock.Controller) *MockPluginLoader {
	mock := &MockPluginLoader{ctrl: ctrl}
	mock.recorder = &MockPluginLoaderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockPluginLoader) EXPECT() *MockPluginLoaderMockRecorder {
	return m.recorder
}

// Compile mocks base method
func (m *MockPluginLoader) Compile(pluginPath string) (core.PluginBearer, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Compile", pluginPath)
	ret0, _ := ret[0].(core.PluginBearer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Compile indicates an expected call of Compile
func (mr *MockPluginLoaderMockRecorder) Compile(pluginPath interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Compile", reflect.TypeOf((*MockPluginLoader)(nil).Compile), pluginPath)
}

// MockPluginBearer is a mock of PluginBearer interface
type MockPluginBearer struct {
	ctrl     *gomock.Controller
	recorder *MockPluginBearerMockRecorder
}

// MockPluginBearerMockRecorder is the mock recorder for MockPluginBearer
type MockPluginBearerMockRecorder struct {
	mock *MockPluginBearer
}

// NewMockPluginBearer creates a new mock instance
func NewMockPluginBearer(ctrl *gomock.Controller) *MockPluginBearer {
	mock := &MockPluginBearer{ctrl: ctrl}
	mock.recorder = &MockPluginBearerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockPluginBearer) EXPECT() *MockPluginBearerMockRecorder {
	return m.recorder
}

// ConfigExists mocks base method
func (m *MockPluginBearer) ConfigExists(pluginName string) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ConfigExists", pluginName)
	ret0, _ := ret[0].(bool)
	return ret0
}

// ConfigExists indicates an expected call of ConfigExists
func (mr *MockPluginBearerMockRecorder) ConfigExists(pluginName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ConfigExists", reflect.TypeOf((*MockPluginBearer)(nil).ConfigExists), pluginName)
}

// CompilePlugin mocks base method
func (m *MockPluginBearer) CompilePlugin(pluginName string, injectedStruct interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CompilePlugin", pluginName, injectedStruct)
	ret0, _ := ret[0].(error)
	return ret0
}

// CompilePlugin indicates an expected call of CompilePlugin
func (mr *MockPluginBearerMockRecorder) CompilePlugin(pluginName, injectedStruct interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CompilePlugin", reflect.TypeOf((*MockPluginBearer)(nil).CompilePlugin), pluginName, injectedStruct)
}

// ForEach mocks base method
func (m *MockPluginBearer) ForEach(callbackFunc func(string) error) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ForEach", callbackFunc)
}

// ForEach indicates an expected call of ForEach
func (mr *MockPluginBearerMockRecorder) ForEach(callbackFunc interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ForEach", reflect.TypeOf((*MockPluginBearer)(nil).ForEach), callbackFunc)
}

// Length mocks base method
func (m *MockPluginBearer) Length() int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Length")
	ret0, _ := ret[0].(int)
	return ret0
}

// Length indicates an expected call of Length
func (mr *MockPluginBearerMockRecorder) Length() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Length", reflect.TypeOf((*MockPluginBearer)(nil).Length))
}

// MockDatabaseConfig is a mock of DatabaseConfig interface
type MockDatabaseConfig struct {
	ctrl     *gomock.Controller
	recorder *MockDatabaseConfigMockRecorder
}

// MockDatabaseConfigMockRecorder is the mock recorder for MockDatabaseConfig
type MockDatabaseConfigMockRecorder struct {
	mock *MockDatabaseConfig
}

// NewMockDatabaseConfig creates a new mock instance
func NewMockDatabaseConfig(ctrl *gomock.Controller) *MockDatabaseConfig {
	mock := &MockDatabaseConfig{ctrl: ctrl}
	mock.recorder = &MockDatabaseConfigMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockDatabaseConfig) EXPECT() *MockDatabaseConfigMockRecorder {
	return m.recorder
}

// Driver mocks base method
func (m *MockDatabaseConfig) Driver() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Driver")
	ret0, _ := ret[0].(string)
	return ret0
}

// Driver indicates an expected call of Driver
func (mr *MockDatabaseConfigMockRecorder) Driver() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Driver", reflect.TypeOf((*MockDatabaseConfig)(nil).Driver))
}

// DBMigrationSource mocks base method
func (m *MockDatabaseConfig) DBMigrationSource() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DBMigrationSource")
	ret0, _ := ret[0].(string)
	return ret0
}

// DBMigrationSource indicates an expected call of DBMigrationSource
func (mr *MockDatabaseConfigMockRecorder) DBMigrationSource() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DBMigrationSource", reflect.TypeOf((*MockDatabaseConfig)(nil).DBMigrationSource))
}

// DBHost mocks base method
func (m *MockDatabaseConfig) DBHost() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DBHost")
	ret0, _ := ret[0].(string)
	return ret0
}

// DBHost indicates an expected call of DBHost
func (mr *MockDatabaseConfigMockRecorder) DBHost() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DBHost", reflect.TypeOf((*MockDatabaseConfig)(nil).DBHost))
}

// DBPort mocks base method
func (m *MockDatabaseConfig) DBPort() (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DBPort")
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DBPort indicates an expected call of DBPort
func (mr *MockDatabaseConfigMockRecorder) DBPort() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DBPort", reflect.TypeOf((*MockDatabaseConfig)(nil).DBPort))
}

// DBUsername mocks base method
func (m *MockDatabaseConfig) DBUsername() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DBUsername")
	ret0, _ := ret[0].(string)
	return ret0
}

// DBUsername indicates an expected call of DBUsername
func (mr *MockDatabaseConfigMockRecorder) DBUsername() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DBUsername", reflect.TypeOf((*MockDatabaseConfig)(nil).DBUsername))
}

// DBPassword mocks base method
func (m *MockDatabaseConfig) DBPassword() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DBPassword")
	ret0, _ := ret[0].(string)
	return ret0
}

// DBPassword indicates an expected call of DBPassword
func (mr *MockDatabaseConfigMockRecorder) DBPassword() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DBPassword", reflect.TypeOf((*MockDatabaseConfig)(nil).DBPassword))
}

// DBDatabase mocks base method
func (m *MockDatabaseConfig) DBDatabase() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DBDatabase")
	ret0, _ := ret[0].(string)
	return ret0
}

// DBDatabase indicates an expected call of DBDatabase
func (mr *MockDatabaseConfigMockRecorder) DBDatabase() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DBDatabase", reflect.TypeOf((*MockDatabaseConfig)(nil).DBDatabase))
}

// DBConnectionMaxLifetime mocks base method
func (m *MockDatabaseConfig) DBConnectionMaxLifetime() (time.Duration, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DBConnectionMaxLifetime")
	ret0, _ := ret[0].(time.Duration)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DBConnectionMaxLifetime indicates an expected call of DBConnectionMaxLifetime
func (mr *MockDatabaseConfigMockRecorder) DBConnectionMaxLifetime() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DBConnectionMaxLifetime", reflect.TypeOf((*MockDatabaseConfig)(nil).DBConnectionMaxLifetime))
}

// DBMaxIddleConn mocks base method
func (m *MockDatabaseConfig) DBMaxIddleConn() (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DBMaxIddleConn")
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DBMaxIddleConn indicates an expected call of DBMaxIddleConn
func (mr *MockDatabaseConfigMockRecorder) DBMaxIddleConn() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DBMaxIddleConn", reflect.TypeOf((*MockDatabaseConfig)(nil).DBMaxIddleConn))
}

// DBMaxOpenConn mocks base method
func (m *MockDatabaseConfig) DBMaxOpenConn() (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DBMaxOpenConn")
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DBMaxOpenConn indicates an expected call of DBMaxOpenConn
func (mr *MockDatabaseConfigMockRecorder) DBMaxOpenConn() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DBMaxOpenConn", reflect.TypeOf((*MockDatabaseConfig)(nil).DBMaxOpenConn))
}

// Dump mocks base method
func (m *MockDatabaseConfig) Dump() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Dump")
	ret0, _ := ret[0].(string)
	return ret0
}

// Dump indicates an expected call of Dump
func (mr *MockDatabaseConfigMockRecorder) Dump() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Dump", reflect.TypeOf((*MockDatabaseConfig)(nil).Dump))
}

// MockDatabaseBearer is a mock of DatabaseBearer interface
type MockDatabaseBearer struct {
	ctrl     *gomock.Controller
	recorder *MockDatabaseBearerMockRecorder
}

// MockDatabaseBearerMockRecorder is the mock recorder for MockDatabaseBearer
type MockDatabaseBearerMockRecorder struct {
	mock *MockDatabaseBearer
}

// NewMockDatabaseBearer creates a new mock instance
func NewMockDatabaseBearer(ctrl *gomock.Controller) *MockDatabaseBearer {
	mock := &MockDatabaseBearer{ctrl: ctrl}
	mock.recorder = &MockDatabaseBearerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockDatabaseBearer) EXPECT() *MockDatabaseBearerMockRecorder {
	return m.recorder
}

// Database mocks base method
func (m *MockDatabaseBearer) Database(dbName string) (*sql.DB, core.DatabaseConfig, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Database", dbName)
	ret0, _ := ret[0].(*sql.DB)
	ret1, _ := ret[1].(core.DatabaseConfig)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Database indicates an expected call of Database
func (mr *MockDatabaseBearerMockRecorder) Database(dbName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Database", reflect.TypeOf((*MockDatabaseBearer)(nil).Database), dbName)
}

// MockAppConfig is a mock of AppConfig interface
type MockAppConfig struct {
	ctrl     *gomock.Controller
	recorder *MockAppConfigMockRecorder
}

// MockAppConfigMockRecorder is the mock recorder for MockAppConfig
type MockAppConfigMockRecorder struct {
	mock *MockAppConfig
}

// NewMockAppConfig creates a new mock instance
func NewMockAppConfig(ctrl *gomock.Controller) *MockAppConfig {
	mock := &MockAppConfig{ctrl: ctrl}
	mock.recorder = &MockAppConfigMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockAppConfig) EXPECT() *MockAppConfigMockRecorder {
	return m.recorder
}

// Port mocks base method
func (m *MockAppConfig) Port() int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Port")
	ret0, _ := ret[0].(int)
	return ret0
}

// Port indicates an expected call of Port
func (mr *MockAppConfigMockRecorder) Port() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Port", reflect.TypeOf((*MockAppConfig)(nil).Port))
}

// BasicAuthUsername mocks base method
func (m *MockAppConfig) BasicAuthUsername() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BasicAuthUsername")
	ret0, _ := ret[0].(string)
	return ret0
}

// BasicAuthUsername indicates an expected call of BasicAuthUsername
func (mr *MockAppConfigMockRecorder) BasicAuthUsername() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BasicAuthUsername", reflect.TypeOf((*MockAppConfig)(nil).BasicAuthUsername))
}

// BasicAuthPassword mocks base method
func (m *MockAppConfig) BasicAuthPassword() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BasicAuthPassword")
	ret0, _ := ret[0].(string)
	return ret0
}

// BasicAuthPassword indicates an expected call of BasicAuthPassword
func (mr *MockAppConfigMockRecorder) BasicAuthPassword() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BasicAuthPassword", reflect.TypeOf((*MockAppConfig)(nil).BasicAuthPassword))
}

// ProxyHost mocks base method
func (m *MockAppConfig) ProxyHost() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ProxyHost")
	ret0, _ := ret[0].(string)
	return ret0
}

// ProxyHost indicates an expected call of ProxyHost
func (mr *MockAppConfigMockRecorder) ProxyHost() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ProxyHost", reflect.TypeOf((*MockAppConfig)(nil).ProxyHost))
}

// PluginExists mocks base method
func (m *MockAppConfig) PluginExists(pluginName string) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PluginExists", pluginName)
	ret0, _ := ret[0].(bool)
	return ret0
}

// PluginExists indicates an expected call of PluginExists
func (mr *MockAppConfigMockRecorder) PluginExists(pluginName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PluginExists", reflect.TypeOf((*MockAppConfig)(nil).PluginExists), pluginName)
}

// Plugins mocks base method
func (m *MockAppConfig) Plugins() []string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Plugins")
	ret0, _ := ret[0].([]string)
	return ret0
}

// Plugins indicates an expected call of Plugins
func (mr *MockAppConfigMockRecorder) Plugins() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Plugins", reflect.TypeOf((*MockAppConfig)(nil).Plugins))
}

// Dump mocks base method
func (m *MockAppConfig) Dump() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Dump")
	ret0, _ := ret[0].(string)
	return ret0
}

// Dump indicates an expected call of Dump
func (mr *MockAppConfigMockRecorder) Dump() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Dump", reflect.TypeOf((*MockAppConfig)(nil).Dump))
}

// Metric mocks base method
func (m *MockAppConfig) Metric() core.MetricConfig {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Metric")
	ret0, _ := ret[0].(core.MetricConfig)
	return ret0
}

// Metric indicates an expected call of Metric
func (mr *MockAppConfigMockRecorder) Metric() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Metric", reflect.TypeOf((*MockAppConfig)(nil).Metric))
}

// MockMetricConfig is a mock of MetricConfig interface
type MockMetricConfig struct {
	ctrl     *gomock.Controller
	recorder *MockMetricConfigMockRecorder
}

// MockMetricConfigMockRecorder is the mock recorder for MockMetricConfig
type MockMetricConfigMockRecorder struct {
	mock *MockMetricConfig
}

// NewMockMetricConfig creates a new mock instance
func NewMockMetricConfig(ctrl *gomock.Controller) *MockMetricConfig {
	mock := &MockMetricConfig{ctrl: ctrl}
	mock.recorder = &MockMetricConfigMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockMetricConfig) EXPECT() *MockMetricConfigMockRecorder {
	return m.recorder
}

// Interface mocks base method
func (m *MockMetricConfig) Interface() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Interface")
	ret0, _ := ret[0].(string)
	return ret0
}

// Interface indicates an expected call of Interface
func (mr *MockMetricConfigMockRecorder) Interface() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Interface", reflect.TypeOf((*MockMetricConfig)(nil).Interface))
}

// MockAppBearer is a mock of AppBearer interface
type MockAppBearer struct {
	ctrl     *gomock.Controller
	recorder *MockAppBearerMockRecorder
}

// MockAppBearerMockRecorder is the mock recorder for MockAppBearer
type MockAppBearerMockRecorder struct {
	mock *MockAppBearer
}

// NewMockAppBearer creates a new mock instance
func NewMockAppBearer(ctrl *gomock.Controller) *MockAppBearer {
	mock := &MockAppBearer{ctrl: ctrl}
	mock.recorder = &MockAppBearerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockAppBearer) EXPECT() *MockAppBearerMockRecorder {
	return m.recorder
}

// Config mocks base method
func (m *MockAppBearer) Config() core.AppConfig {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Config")
	ret0, _ := ret[0].(core.AppConfig)
	return ret0
}

// Config indicates an expected call of Config
func (mr *MockAppBearerMockRecorder) Config() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Config", reflect.TypeOf((*MockAppBearer)(nil).Config))
}

// DownStreamPlugins mocks base method
func (m *MockAppBearer) DownStreamPlugins() []core.DownStreamPlugin {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DownStreamPlugins")
	ret0, _ := ret[0].([]core.DownStreamPlugin)
	return ret0
}

// DownStreamPlugins indicates an expected call of DownStreamPlugins
func (mr *MockAppBearerMockRecorder) DownStreamPlugins() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DownStreamPlugins", reflect.TypeOf((*MockAppBearer)(nil).DownStreamPlugins))
}

// InjectDownStreamPlugin mocks base method
func (m *MockAppBearer) InjectDownStreamPlugin(InjectedDownStreamPlugin core.DownStreamPlugin) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "InjectDownStreamPlugin", InjectedDownStreamPlugin)
}

// InjectDownStreamPlugin indicates an expected call of InjectDownStreamPlugin
func (mr *MockAppBearerMockRecorder) InjectDownStreamPlugin(InjectedDownStreamPlugin interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InjectDownStreamPlugin", reflect.TypeOf((*MockAppBearer)(nil).InjectDownStreamPlugin), InjectedDownStreamPlugin)
}

// InjectController mocks base method
func (m *MockAppBearer) InjectController(injectedController core.Controller) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "InjectController", injectedController)
}

// InjectController indicates an expected call of InjectController
func (mr *MockAppBearerMockRecorder) InjectController(injectedController interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InjectController", reflect.TypeOf((*MockAppBearer)(nil).InjectController), injectedController)
}
