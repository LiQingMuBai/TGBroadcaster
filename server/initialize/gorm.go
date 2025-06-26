package initialize

import (
	"gorm.io/gorm"
)

func Gorm_A() *gorm.DB {
	return GormMysql_A()
}
func Gorm_F() *gorm.DB {
	return GormMysql_F()
}

func Gorm_B() *gorm.DB {
	return GormMysql_B()
}
func Gorm_C() *gorm.DB {
	return GormMysql_C()
}
func Gorm_D() *gorm.DB {
	return GormMysql_D()
}
func Gorm_E() *gorm.DB {
	return GormMysql_E()
}
func Gorm_Local() *gorm.DB {
	return GormMysql_Local()
}
