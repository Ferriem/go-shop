package repositories

import (
	"context"
	"database/sql"
	"errors"
	"strconv"
	"time"

	"github.com/Ferriem/go-web/Product/common"
	"github.com/Ferriem/go-web/Product/datamodels"
	"github.com/redis/go-redis/v9"
)

type IUser interface {
	Conn() error
	Select(string) (*datamodels.User, error)
	Insert(*datamodels.User) (int64, error)
	LogDefinition(string) (*datamodels.User, error)
	FailPlusOne(string) error
}

type UserManager struct {
	table     string
	mysqlConn *sql.DB
	rdb       *redis.Client
}

func NewUserManager(table string, db *sql.DB) IUser {
	return &UserManager{table: table, mysqlConn: db}
}

func (u *UserManager) Conn() error {
	if u.mysqlConn == nil {
		mysql, err := common.NewMysqlConn()
		if err != nil {
			return err
		}
		u.mysqlConn = mysql
	}
	if u.rdb == nil {
		redis := common.NewRedisConn()
		u.rdb = redis
	}
	if u.table == "" {
		u.table = "user"
	}
	return nil
}

func (u *UserManager) Select(username string) (*datamodels.User, error) {
	if username == "" {
		return &datamodels.User{}, errors.New("username can't be empty!")
	}
	if err := u.Conn(); err != nil {
		return &datamodels.User{}, err
	}

	sql := "SELECT * FROM " + u.table + " WHERE userName=?"
	row, err := u.mysqlConn.Query(sql, username)
	defer row.Close()
	if err != nil {
		return &datamodels.User{}, err
	}

	result := common.GetResultRow(row)
	if len(result) == 0 {
		return &datamodels.User{}, errors.New("User didn't exist")
	}

	user := &datamodels.User{}
	common.DataToStructByTagSql(result, user)
	return user, nil
}

// register a user
func (u *UserManager) Insert(user *datamodels.User) (int64, error) {
	if err := u.Conn(); err != nil {
		return 0, err
	}
	ctx := context.Background()
	exists, err := u.rdb.Exists(ctx, user.UserName).Result()
	if err != nil {
		return 0, err
	}
	if exists == 1 {
		return 0, errors.New("User already exists")
	}

	sql := "INSERT " + u.table + " SET nickName=?,userName=?,password=?"
	stmt, err := u.mysqlConn.Prepare(sql)
	if err != nil {
		return 0, err
	}
	result, err := stmt.Exec(user.NickName, user.UserName, user.HashPassword)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	redis_user := map[string]interface{}{
		"ID":       id,
		"password": user.HashPassword,
	}
	err = u.rdb.HMSet(ctx, user.UserName, redis_user).Err()
	if err != nil {
		u.mysqlConn.Exec("DELETE FROM "+u.table+" WHERE userName=?", user.UserName)
		return 0, err
	}
	return id, nil
}

func (u *UserManager) SelectByID(id int64) (*datamodels.User, error) {
	if err := u.Conn(); err != nil {
		return &datamodels.User{}, err
	}
	sql := "SELECT * FROM " + u.table + " WHERE ID = ?" + strconv.FormatInt(id, 10)
	row, err := u.mysqlConn.Query(sql)
	if err != nil {
		return &datamodels.User{}, err
	}

	result := common.GetResultRow(row)
	if len(result) == 0 {
		return &datamodels.User{}, errors.New("User didn't exist")
	}
	user := &datamodels.User{}
	common.DataToStructByTagSql(result, user)

	return user, nil
}

func (u *UserManager) SelectByName(name string) (*datamodels.User, error) {
	if err := u.Conn(); err != nil {
		return &datamodels.User{}, err
	}
	sql := "SELECT * FROM " + u.table + " WHERE userName = ?"
	row, err := u.mysqlConn.Query(sql, name)
	if err != nil {
		return &datamodels.User{}, err
	}
	result := common.GetResultRow(row)
	if len(result) == 0 {
		return &datamodels.User{}, errors.New("User didn't exist")
	}
	user := &datamodels.User{}
	common.DataToStructByTagSql(result, user)
	return user, nil
}

func (u *UserManager) LogDefinition(name string) (*datamodels.User, error) {
	if err := u.Conn(); err != nil {
		return &datamodels.User{}, err
	}
	ctx := context.Background()

	exist, _ := u.rdb.Exists(ctx, name+"failTime").Result()
	if exist == 1 {
		failTime, err := u.rdb.Get(ctx, name+"failTime").Result()
		if err != nil {
			return &datamodels.User{}, err
		}
		if failTime > "5" {
			return &datamodels.User{}, errors.New("The account is locked")
		}
	}

	user := &datamodels.User{}

	values, err := u.rdb.HMGet(ctx, name, "ID", "password").Result()
	if err != nil {
		return &datamodels.User{}, err
	}
	user.ID, _ = strconv.ParseInt(values[0].(string), 10, 64)
	user.UserName = name
	user.HashPassword = values[1].(string)
	return user, nil
}

func (u *UserManager) FailPlusOne(name string) error {
	if err := u.Conn(); err != nil {
		return err
	}
	ctx := context.Background()
	exist, err := u.rdb.Exists(ctx, name+"failTime").Result()
	if err != nil {
		return err
	}
	var times int
	if exist == 1 {
		result, err := u.rdb.Get(ctx, name+"failTime").Result()
		if err != nil {
			return err
		}
		times, _ = strconv.Atoi(result)
	}
	times++
	err = u.rdb.Set(ctx, name+"failTime", times, 0).Err()
	if err != nil {
		return err
	}
	err = u.rdb.Expire(ctx, name+"failTime", 5*time.Minute).Err()
	if err != nil {
		return err
	}
	return nil
}
