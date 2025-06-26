package initialize

import (
	"boost/data/server/config"
	"boost/data/server/global"
	"boost/data/server/initialize/internal"
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func GormMysql_F() *gorm.DB {
	m := global.GVA_CONFIG.Mysql_F
	return initMysqlDatabase(m)
}
func GormMysql_A() *gorm.DB {
	m := global.GVA_CONFIG.Mysql_A
	return initMysqlDatabase(m)
}
func GormMysql_B() *gorm.DB {
	m := global.GVA_CONFIG.Mysql_B
	return initMysqlDatabase(m)
}
func GormMysql_C() *gorm.DB {
	m := global.GVA_CONFIG.Mysql_C
	return initMysqlDatabase(m)
}
func GormMysql_D() *gorm.DB {
	m := global.GVA_CONFIG.Mysql_D
	return initMysqlDatabase(m)
}
func GormMysql_E() *gorm.DB {
	m := global.GVA_CONFIG.Mysql_E
	return initMysqlDatabase(m)
}
func GormMysql_Local() *gorm.DB {
	m := global.GVA_CONFIG.Mysql_Local
	return initMysqlDatabase(m)
}

// initMysqlDatabase 初始化Mysql数据库的辅助函数
func initMysqlDatabase(m config.Mysql) *gorm.DB {
	if m.Dbname == "" {
		return nil
	}

	mysqlConfig := mysql.Config{
		DSN:                       m.Dsn(), // DSN data source name
		DefaultStringSize:         191,     // string 类型字段的默认长度
		SkipInitializeWithVersion: false,   // 根据版本自动配置
	}

	if db, err := gorm.Open(mysql.New(mysqlConfig), internal.Gorm.Config(m.Prefix, m.Singular)); err != nil {
		panic(err)
	} else {
		db.InstanceSet("gorm:table_options", "ENGINE="+m.Engine)
		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(m.MaxIdleConns)
		sqlDB.SetMaxOpenConns(m.MaxOpenConns)
		return db
	}
}
