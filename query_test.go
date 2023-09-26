package query

import (
	"database/sql"
	_ "embed"
	"sync"
	"testing"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type User struct {
	BaseModel
	Name string
}

var (
	_db  *gorm.DB
	once sync.Once
)

//go:embed dsn.env
var dsn string

func init() {
	initMysqlDB(dsn, true)
}

func initMysqlDB(dsn string, debug bool) {
	once.Do(func() {
		sqlDB, err := sql.Open("mysql", dsn)
		if err != nil {
			panic(err)
		}
		db, err := gorm.Open(mysql.New(mysql.Config{
			Conn: sqlDB,
		}))
		if err != nil {
			panic(err)
		}

		if debug {
			db = db.Debug()
		}

		// 注册opentracing插件
		if err := db.Use(NewOpentracingPlugin()); err != nil {
			panic(err)
		}

		_db = db
	})
}

func TestNewOrder(t *testing.T) {
	NewAction[User]().WithDB(_db).Order("name").Desc().First()
	NewAction[User]().WithDB(_db).Order("name").Asc().First()
	NewAction[User]().WithDB(_db).Order("id").Desc().First()
	NewAction[User]().WithDB(_db).Order("id").Asc().First()
}

func TestFirst(t *testing.T) {
	first, err := NewAction[User]().WithDB(_db).First()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(first)
}

func TestLast(t *testing.T) {
	last, err := NewAction[User]().WithDB(_db).Last()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(last)
}

func TestFind(t *testing.T) {
	pageInfo := NewPage(1, 10)
	users, err := NewAction[User]().WithDB(_db).List(pageInfo)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(users, pageInfo)
}

func TestCrud(t *testing.T) {
	user := &User{
		Name: "test",
	}
	if err := NewAction[User]().WithDB(_db).Create(user); err != nil {
		t.Fatal(err)
	}
	t.Log(user)

	// 更新
	user.Name = "test2"
	if err := NewAction[User]().WithDB(_db).Update(user); err != nil {
		t.Fatal(err)
	}

	// 查询
	newUser, err := NewAction[User]().WithDB(_db).First(WhereID(user.ID))
	if err != nil {
		t.Fatal(err)
	}

	if newUser.Name != user.Name {
		t.Fatal("update failed")
	}

	if newUser.Name != "test2" {
		t.Fatal("update failed")
	}

	// 软删除
	if err := NewAction[User]().WithDB(_db).Delete(WhereID(user.ID)); err != nil {
		t.Fatal(err)
	}

	// 硬删除
	if err := NewAction[User]().WithDB(_db).ForcedDelete(WhereID(user.ID)); err != nil {
		t.Fatal(err)
	}
}
