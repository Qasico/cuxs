package cuxs

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/qasico/cuxs/log"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/mssql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var Orm  map[string]*gorm.DB

func init() {
	Orm = make(map[string]*gorm.DB)
}

func NewDB(name interface{}) {
	var orm *gorm.DB
	var err error

	c := Config.DatabaseConfig
	dbname := c.DBName

	if n, ok := name.(string); ok && n != "" {
		dbname = n
	}

	orm, err = gorm.Open(c.Engine, openConnection(c, dbname))
	if err != nil {
		log.Errorf("Cannot connect to database, %s", err.Error())
	}

	if Config.Runmode == "dev" {
		log.Infof("Connected database engine %s", log.Color.CyanBg(fmt.Sprintf(" %s on %s:%d ", c.Engine, c.ServerHost, c.ServerPort), "1"))
		orm.LogMode(true)
	}

	orm.DB().SetMaxIdleConns(c.IdleMax)
	orm.DB().SetMaxOpenConns(c.ConnMax)
	orm.SetLogger(log.OrmLogger{})

	Orm[dbname] = orm
}

func openConnection(c DatabaseConfig, dbname string) (conn string) {
	switch c.Engine {
	case "mysql":
		conn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8", c.DBUser, c.DBPassword, c.ServerHost, c.ServerPort, dbname)
	default:
		conn = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", c.ServerHost, c.ServerPort, c.DBUser, c.DBPassword, dbname, "disable")
	}

	return
}

func ORM() *gorm.DB {
	return Orm[Config.DatabaseConfig.DBName]
}