package sqlutil

type TestDB struct {
	DriverName string
	ConnStr    string
}

var TestDB_Sqlite3 = TestDB{DriverName: "sqlite3", ConnStr: ":memory:"}
var TestDB_Mysql = TestDB{DriverName: "mysql", ConnStr: "merkut:password@tcp(localhost:3306)/merkut_tests?collation=utf8mb4_unicode_ci"}
var TestDB_Postgres = TestDB{DriverName: "postgres", ConnStr: "user=merkuttest password=merkuttest host=localhost port=5432 dbname=merkuttest sslmode=disable"}
var TestDB_Mssql = TestDB{DriverName: "mssql", ConnStr: "server=localhost;port=1433;database=merkuttest;user id=merkut;password=Password!"}

