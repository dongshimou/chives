package huobi

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"log"
)

var DBTableConfig = []string{
	"engine=InnoDB",
	"auto_increment=1",
	"charset=utf8",
}
var DBDatabaseConfig = []string{
	"set auto_increment_increment=1",
	"set auto_increment_offset=1;",
}
var (
	db       *gorm.DB
	dbConfig *DBConfig
)

type Tick struct {
	OrderID        int64
	TradeID        int64 `gorm:"index;unique"`
	TradeAmount    float64
	TradePrice     float64
	TradeDirection bool //buy==1 sell==0
	TradeTimeStamp int64
}
type Kline struct {
	Timestamp int64 `gorm:"index;unique"`
	Open      float64
	Close     float64
	High      float64
	Low       float64
	Amount    float64
	Volume    float64
}

func IsDBSave() bool {
	return dbConfig.DBSave
}
func setDBNotSave() {
	dbConfig.DBSave = false
}
func InitDB(market string) (err error) {

	dbConfig = &DBConfig{}

	err = readConfig("./database", &dbConfig)
	if err != nil {
		setDBNotSave()
		return err
	}
	func() {
		temp, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/",
			dbConfig.DBUser, dbConfig.DBPassword, dbConfig.DBHost, dbConfig.DBPort))
		defer temp.Close()
		if err != nil {
			log.Println(err.Error())
			setDBNotSave()
			return
			//panic(err)
		}
		_, err = temp.Exec("CREATE DATABASE IF NOT EXISTS " + dbConfig.DBDatabase + market)
		if err != nil {
			log.Println(err.Error())
			setDBNotSave()
			return
			//panic(err)
		}
	}()
	if !dbConfig.DBSave {
		return nil
	}
	if db != nil {
		db.Close()
	}

	args := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		dbConfig.DBUser, dbConfig.DBPassword, dbConfig.DBHost, dbConfig.DBPort, dbConfig.DBDatabase+market)

	db, err = gorm.Open("mysql", args)
	if err != nil {
		setDBNotSave()
		return err
	}

	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return defaultTableName
	}

	//读取数据库设置
	for _, v := range DBDatabaseConfig {
		db = db.Exec(v)
	}

	first_init := true
	init_options := func() {
		if first_init {
			//创建表的时候的配置
			for _, v := range DBTableConfig {
				db = db.Set("gorm:table_options", v)
			}
			first_init = false
		}
	}
	//检查表示否存在,然后创建表
	if db.HasTable(&Tick{}) == false {
		init_options()
		db = db.CreateTable(&Tick{})
	}
	if db.HasTable(&Kline{}) == false {
		init_options()
		db = db.CreateTable(&Kline{})
	}
	if db.Error != nil {
		setDBNotSave()
		return db.Error
	} else {
		return nil
	}
}
func Insert(tx *gorm.DB, value interface{}) (err error) {
	if tx == nil {
		tx = db
	}
	res := tx.Create(value)
	return res.Error
}
func Update(tx *gorm.DB, model, value interface{}) (err error) {
	if tx == nil {
		tx = db
	}
	res := tx.Model(model).Updates(value)
	return res.Error
}
func Delete(tx *gorm.DB, value interface{}) (err error) {
	if tx == nil {
		tx = db
	}
	return tx.Delete(value).Error
}
func CreateTX() *gorm.DB {
	return db.Begin()
}
func BreakTX(tx *gorm.DB) {
	tx.Rollback()
}
func FinishTX(tx *gorm.DB) {
	tx.Commit()
}
func GetDBByModel(value interface{}) *gorm.DB {
	return db.Model(value)
}
func saveKline(klines chan Kline) {
	for {
		data, ok := <-klines
		if !ok {
			return
		}
		tx := CreateTX()
		err := Insert(tx, data)
		if err != nil {
			BreakTX(tx)
		} else {
			FinishTX(tx)
		}
	}

}
func saveTick(ticks chan TradeTick) {
	for {
		data, ok := <-ticks
		if !ok {
			return
		}
		tx := CreateTX()
		insertStatus := true
		for _, v := range data.Data {
			err := Insert(tx, Tick{
				OrderID:     data.ID,
				TradeID:     v.ID,
				TradeAmount: v.Amount,
				TradeDirection: func() bool {
					if v.Direction == BUY {
						return true
					} else {
						return false
					}
				}(),
				TradeTimeStamp: v.TS,
				TradePrice:     v.Price,
			})
			if err != nil {
				insertStatus = false
				break
			}
		}
		if insertStatus {
			FinishTX(tx)
		} else {
			BreakTX(tx)
		}
	}
}
