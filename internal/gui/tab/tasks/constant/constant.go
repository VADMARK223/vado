package constant

import "fmt"

const TasksFilePath = "./data/tasks.json"

//var TasksDataSourceName = util.Tpl(
//	"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
//	"127.0.0.1", 5432, "vadmark", "5125341", "vadodb",
//)

const (
	Host     = "127.0.0.1"
	Port     = 5432
	User     = "vadmark"
	Password = "5125341"
	DBName   = "vadodb"
)

var TasksDataSourceName = fmt.Sprintf(
	"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
	Host, Port, User, Password, DBName,
)

//const TasksBaseURL = "http://localhost:8080/v1"
