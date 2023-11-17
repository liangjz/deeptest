package agentExec

import (
	"database/sql"
	"encoding/json"
	"fmt"
	queryUtils "github.com/aaronchen2k/deeptest/internal/agent/exec/utils/query"
	"github.com/aaronchen2k/deeptest/internal/pkg/consts"
	"github.com/aaronchen2k/deeptest/internal/pkg/domain"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

func ExecDbOpt(opt *domain.DatabaseOptBase) (err error) {
	if opt.Type == "" {
		return
	}

	opt.ResultStatus = consts.Pass

	if opt.Type == consts.DbTypeOracle {
		orclDb, err1 := OpenOracle(opt)
		if err1 != nil {
			err = err1
			return
		}

		err1 = queryOracle(orclDb, opt)
		err = err1

		return
	}

	var db *gorm.DB

	if opt.Type == consts.DbTypeMySql {
		db, err = OpenMySqlDb(opt)
	} else if opt.Type == consts.DbTypeSqlServer {
		db, err = OpenSqlServer(opt)
	} else if opt.Type == consts.DbTypePostgreSql {
		db, err = OpenPostgreSQL(opt)
	} else if opt.Type == consts.DbTypeOracle {
		db, err = OpenPostgreSQL(opt)
	}

	if err != nil {
		return
	}

	queryResult, err := query(db, opt)
	opt.Result = queryUtils.JsonPath(string(queryResult), opt.JsonPath)

	return
}

func OpenMySqlDb(opt *domain.DatabaseOptBase) (db *gorm.DB, err error) {
	params := "charset=utf8mb4&parseTime=True&loc=Local"

	connStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s", opt.Username, opt.Password,
		opt.Host, opt.Port, opt.DbName, params)

	config := mysql.Config{
		DSN: connStr,
	}

	db, err = gorm.Open(mysql.New(config), &gorm.Config{})

	return
}

func OpenSqlServer(opt *domain.DatabaseOptBase) (db *gorm.DB, err error) {
	connStr := fmt.Sprintf("sqlserver://%s:%s@%s:%s?database=%s",
		opt.Username, opt.Password,
		opt.Host, opt.Port, opt.DbName)

	db, err = gorm.Open(sqlserver.Open(connStr), &gorm.Config{})

	return
}

func OpenPostgreSQL(opt *domain.DatabaseOptBase) (db *gorm.DB, err error) {
	params := "sslmode=disable TimeZone=Asia/Shanghai"

	connStr := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s %s",
		opt.Username, opt.Password,
		opt.Host, opt.Port, opt.DbName, params)

	db, err = gorm.Open(postgres.Open(connStr), &gorm.Config{})

	return
}

func OpenOracle(opt *domain.DatabaseOptBase) (db *sql.DB, err error) {
	connStr := fmt.Sprintf("%s/%s@%s:%s/%s",
		opt.Username, opt.Password,
		opt.Host, opt.Port, opt.DbName)

	db, err = sql.Open("goracle", connStr)

	return
}

func query(db *gorm.DB, opt *domain.DatabaseOptBase) (result []byte, err error) {
	data := []map[string]interface{}{}

	err = db.Raw(opt.Sql).
		Scan(&data).Error

	result, err = json.Marshal(data)

	return
}

func queryOracle(db *sql.DB, opt *domain.DatabaseOptBase) (err error) {
	rows, err := db.Query(opt.Sql, 100)
	if err != nil {
		return
	}
	defer rows.Close()

	cols, _ := rows.Columns()

	data := make([]map[string]interface{}, 0)
	for rows.Next() {
		columns := make([]interface{}, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i, _ := range columns {
			columnPointers[i] = &columns[i]
		}

		if err = rows.Scan(columnPointers...); err != nil {
			return
		}

		mp := make(map[string]interface{})
		for i, colName := range cols {
			val := columnPointers[i].(*interface{})
			mp[colName] = *val
		}

		data = append(data, mp)
	}

	bytes, err := json.Marshal(data)
	opt.ResultMsg = string(bytes)

	return
}
