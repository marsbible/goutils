package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"goutils"
	"log"
	"reflect"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var (
	db *sql.DB
)

type AllDataType struct {
	ID         int64     `gorm:"column:id;primary_key" json:"id"`
	Varchar    string    `gorm:"column:varchar" json:"varchar"`
	Tinyint    int       `gorm:"column:tinyint" json:"tinyint"`
	Text       string    `gorm:"column:text" json:"text"`
	Date       time.Time `gorm:"column:date" json:"date"`
	Smallint   int       `gorm:"column:smallint" json:"smallint"`
	Mediumint  int       `gorm:"column:mediumint" json:"mediumint"`
	Int        int       `gorm:"column:int" json:"int"`
	Bigint     int64     `gorm:"column:bigint" json:"bigint"`
	Float      float32   `gorm:"column:float" json:"float"`
	Double     float64   `gorm:"column:double" json:"double"`
	Decimal    float64   `gorm:"column:decimal" json:"decimal"`
	Datetime   time.Time `gorm:"column:datetime" json:"datetime"`
	Timestamp  time.Time `gorm:"column:timestamp" json:"timestamp"`
	Time       time.Time `gorm:"column:time" json:"time"`
	Char       string    `gorm:"column:char" json:"char"`
	Tinyblob   []byte    `gorm:"column:tinyblob" json:"tinyblob"`
	Tinytext   string    `gorm:"column:tinytext" json:"tinytext"`
	Blob       []byte    `gorm:"column:blob" json:"blob"`
	Mediumblob []byte    `gorm:"column:mediumblob" json:"mediumblob"`
	Mediumtext string    `gorm:"column:mediumtext" json:"mediumtext"`
	Longblob   []byte    `gorm:"column:longblob" json:"longblob"`
	Longtext   string    `gorm:"column:longtext" json:"longtext"`
	Enum       string    `gorm:"column:enum" json:"enum"`
	Set        string    `gorm:"column:set" json:"set"`
	Bool       int       `gorm:"column:bool" json:"bool"`
	Binary     []byte    `gorm:"column:binary" json:"binary"`
	Varbinary  []byte    `gorm:"column:varbinary" json:"varbinary"`
}

// TableName sets the insert table name for this struct type
func (a *AllDataType) TableName() string {
	return "all_data_types"
}

type Test1 struct {
	a int
	b string
	c float64
	d []int
	e [8]string
}

func init() {
	goutils.RegisterType((*Test1)(nil))
	url := fmt.Sprintf("%s:%s@(%s:%v)/%s?charset=utf8&parseTime=True&loc=Local",
		"jbex", "jbex", "127.0.0.1", 3306, "jbex_com")

	dd, err := sql.Open("mysql", url)
	if err != nil {
		log.Fatal(err.Error())
	}
	db = dd
}

func testReflectNew() {
	t := goutils.MakeInstance("Test1").(Test1)
	tp := goutils.MakeInstancePtr("Test1").(*Test1)

	t.a = 1
	t.b = "b"
	t.c = 3.14
	t.d = []int{1, 2, 3}
	t.e = [8]string{"a", "b", "c", "d", "e", "f", "g", "h"}
	tp.a = 1
	tp.b = "b"
	tp.c = 3.14
	tp.d = []int{1, 2, 3}
	tp.e = [8]string{"a", "b", "c", "d", "e", "f", "g", "h"}

	if fmt.Sprint(t) != fmt.Sprint(*tp) {
		panic("Reflect new test failed")
	}
}

func testStructScan() {
	rows, err := db.Query("select * from all_data_types")
	if err != nil {
		panic("query db  failed")
	}

	var v interface{}
	for rows.Next() {
		v, err = goutils.StructScan(rows, reflect.TypeOf(AllDataType{}))
		if err != nil {
			panic(err)
		}
		x := v.(AllDataType)
		bb, _ := json.Marshal(x)
		fmt.Println(string(bb))
	}
}

func main() {
	testReflectNew()
	testStructScan()
}
