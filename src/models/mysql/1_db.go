package mysql

import (
	"database/sql"
	"libs/utils"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

var (
	pool *sql.DB
)

// InitPool 开启连接池
func InitPool() {
	// 获取参数
	mysqlUser := utils.GetConfigStr("conf", "app.conf", "mysql", "mysql_user") + utils.GetConfigStr("conf", "app.conf", "server", "http_port")
	mysqlPwd := utils.GetConfigStr("conf", "app.conf", "mysql", "mysql_pwd")
	mysqlAddress := utils.GetConfigStr("conf", "app.conf", "mysql", "mysql_address")
	mysqlDB := utils.GetConfigStr("conf", "app.conf", "mysql", "mysql_db")
	dataSourceName := mysqlUser + ":" + mysqlPwd + "@tcp(" + mysqlAddress + ")/" + mysqlDB
	// 开始连接
	var err error
	pool, err = sql.Open("mysql", dataSourceName+"?charset=utf8")
	utils.LogFatal("mysql", err)
	// 连接池设置
	//maxOpenConn := utils.GetConfigInt("conf", "app.conf", "mysql", "mysql_max_open_conn")
	//maxIdleConn := utils.GetConfigInt("conf", "app.conf", "mysql", "mysql_max_idle_conn")
	//connMaxLifeMin := utils.GetConfigInt64("conf", "app.conf", "mysql", "mysql_conn_max_life_min")
	//pool.SetMaxOpenConns(maxOpenConn) // 最大连接数
	//pool.SetMaxIdleConns(maxIdleConn) // 设置闲置连接数
	//pool.SetConnMaxLifetime(time.Minute * time.Duration(connMaxLifeMin))
	err = pool.Ping() // 调用完毕后会马上把连接返回给连接池。
	utils.LogFatal("mysql", err)
}

// ClosePool
func ClosePool() {
	if pool != nil {
		pool.Close()
	}
}

//************************************************************************************//
//************************************ mysql工具类 ************************************//
//************************************************************************************//

type SqlCURD struct {
	conn *sql.DB
	tx   *sql.Tx
	// curd
	curd       string
	table      string
	selectArgs string
	set        string
	join       string
	on         string
	where      string
	whereArgs  string
	group      string
	having     string
	order      string
	limit      string
	// result
	err    error
	rows   *sql.Rows
	stmt   *sql.Stmt
	result sql.Result
}

// 普通连接
func mysqlDB() *SqlCURD {
	db := &SqlCURD{conn: pool}
	return db
}

// Close 手动close
func (m *SqlCURD) Close() *SqlCURD {
	if m.rows != nil {
		defer m.rows.Close()
	}
	if m.stmt != nil {
		defer m.stmt.Close()
	}
	return m
}

// 事务
func mysqlTX() *SqlCURD {
	tx, _ := pool.Begin()
	db := &SqlCURD{tx: tx}
	return db
}

// 事务(已有的)
func mysqlTX2(tx *sql.Tx) *SqlCURD {
	db := &SqlCURD{tx: tx}
	return db
}

// 事务提交
func (m *SqlCURD) Commit() error {
	err := m.tx.Commit()
	if err != nil {
		utils.LogErr("sql", err)
		m.tx.Rollback()
	}
	return err
}

//************************************** SELECT **************************************//
// Select
func (m *SqlCURD) Select(selectArgs string) *SqlCURD {
	m.curd = "SELECT "
	m.selectArgs = selectArgs + " "
	return m
}

// Form
func (m *SqlCURD) Form(table string) *SqlCURD {
	m.table = table + " "
	return m
}

// Query
func (m *SqlCURD) Query(args ...interface{}) *SqlCURD {
	query := m.curd + m.selectArgs + "FROM " + m.table + m.join + m.on + m.where + m.group + m.having + m.order + m.limit
	query = strings.TrimSpace(query) + ";"
	utils.LogDebug("sql", query)
	if m.err != nil {
		if m.tx != nil {
			m.tx.Rollback()
		}
		return m
	} else if m.tx != nil {
		m.tx.Exec("SET NAMES utf8mb4")
		m.rows, m.err = m.tx.Query(query, args...)
		if m.err != nil {
			m.tx.Rollback()
		}
	} else {
		m.conn.Exec("SET NAMES utf8mb4")
		m.rows, m.err = m.conn.Query(query, args...)
	}
	utils.LogErr("sql", m.err)
	return m
}

// NextScan 内部自动close
func (m *SqlCURD) NextScan(dest ...interface{}) *SqlCURD {
	if m.err != nil {
		return m
	} else if m.Next() {
		m.Scan(dest...)
	}
	return m
}

func (m *SqlCURD) Next() bool {
	if m.err == nil && m.rows != nil {
		return m.rows.Next()
	}
	return false
}

func (m *SqlCURD) Scan(dest ...interface{}) error {
	if m.err == nil && m.rows != nil {
		m.err = m.rows.Scan(dest...)
		utils.LogErr("sql", m.err)
	}
	return m.err
}

//******************************** INSERT/UPDATE/DELETE ********************************//
// Insert
func (m *SqlCURD) Insert(table string) *SqlCURD {
	m.curd = "INSERT "
	m.table = table + " "
	return m
}

func (m *SqlCURD) Update(table string) *SqlCURD {
	m.curd = "UPDATE "
	m.table = table + " "
	return m
}

//func (m *SqlCURD) Delete(table string) *SqlCURD {
//	m.curd = "DELETE "
//	m.table = table + " "
//	return m
//}

// Set
func (m *SqlCURD) Set(set string) *SqlCURD {
	m.set = "SET " + set + " "
	return m
}

func (m *SqlCURD) Exec(args ...interface{}) *SqlCURD {
	query := m.curd + m.table + m.set + m.join + m.on + m.where + m.group + m.having + m.order + m.limit
	query = strings.TrimSpace(query) + ";"
	utils.LogDebug("sql", query)
	if m.err != nil {
		if m.tx != nil {
			m.tx.Rollback()
		}
		return m
	} else if m.tx != nil {
		m.tx.Exec("SET NAMES utf8mb4")
		m.stmt, m.err = m.tx.Prepare(query)
		if m.err == nil {
			m.result, m.err = m.stmt.Exec(args...)
		} else {
			m.err = m.tx.Rollback()
		}
	} else {
		m.conn.Exec("SET NAMES utf8mb4")
		m.stmt, m.err = m.conn.Prepare(query)
		if m.err == nil {
			m.result, m.err = m.stmt.Exec(args...)
		}
	}
	utils.LogErr("sql", m.err)
	return m
}

func (m *SqlCURD) Result() sql.Result {
	return m.result
}

//************************************** COMMON **************************************//
// Err
func (m *SqlCURD) Err() error {
	return m.err
}

func (m *SqlCURD) Join(table string) *SqlCURD {
	m.join = "JOIN " + table + " "
	return m
}

func (m *SqlCURD) LeftJoin(table string) *SqlCURD {
	if len(strings.TrimSpace(table)) < 0 {
		return m
	}
	m.join = "LEFT JOIN " + table + " "
	return m
}

func (m *SqlCURD) RightJoin(table string) *SqlCURD {
	if len(strings.TrimSpace(table)) <= 0 {
		return m
	}
	m.join = "RIGHT JOIN " + table + " "
	return m
}

func (m *SqlCURD) On(on string) *SqlCURD {
	m.on = "ON " + on + " "
	return m
}

// Where
func (m *SqlCURD) Where(where string) *SqlCURD {
	m.where = "WHERE " + where + " "
	return m
}

// Where
func (m *SqlCURD) Group(group string) *SqlCURD {
	m.group = "GROUP BY " + group + " "
	return m
}

// Having
func (m *SqlCURD) Having(having string) *SqlCURD {
	m.having = "HAVING " + having + " "
	return m
}

func (m *SqlCURD) OrderUp(order string) *SqlCURD {
	m.Order(order + " ASC")
	return m
}

func (m *SqlCURD) OrderDown(order string) *SqlCURD {
	m.Order(order + " DESC")
	return m
}

func (m *SqlCURD) Order(order string) *SqlCURD {
	m.order = "ORDER BY " + order + " "
	return m
}

func (m *SqlCURD) Limit(offset int, limit int) *SqlCURD {
	if limit <= 0 {
		return m
	}
	m.limit = "LIMIT " + strconv.Itoa(offset) + "," + strconv.Itoa(limit) + " "
	return m
}
