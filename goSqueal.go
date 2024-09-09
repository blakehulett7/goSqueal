package goSqueal

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"strings"
)

var defaultOpenPermissions int = 0777

func CheckForTable(tableName string) {
	path := fmt.Sprintf("./database/%v.db", tableName)
	_, err := os.Stat(path)
	if !errors.Is(err, fs.ErrNotExist) {
		fmt.Println(tableName, "table exists")
		return
	}
	fmt.Println(tableName, "table does not exist, creating it now...")
	command := fmt.Sprintf("cat init/%v.sql | sqlite3 database/%v.db", tableName, tableName)
	err = exec.Command("bash", "-c", command).Run()
	if err != nil {
		fmt.Println(err)
	}
}

func GetTableFields(tableName string) []string {
	sqlQueryString := fmt.Sprintf("SELECT name FROM pragma_table_info('%v')", tableName)
	os.WriteFile("query.sql", []byte(sqlQueryString), fs.FileMode(defaultOpenPermissions))
	defer exec.Command("rm", "query.sql").Run()

	path := "database/" + tableName + ".db"

	command := fmt.Sprintf("cat query.sql | sqlite3 %v", path)
	fieldsData, err := exec.Command("bash", "-c", command).Output()
	if err != nil {
		fmt.Println("error:", err)
	}
	fields := string(fieldsData)
	fieldsSlice := strings.Split(fields, "\n")
	fieldsSlice = fieldsSlice[:len(fieldsSlice)-1]
	return fieldsSlice
}

func CreateTableEntry(tableName string, params map[string]string) {
	sqlQueryString := fmt.Sprintf("INSERT INTO %v VALUES", tableName)
	fields := GetTableFields(tableName)
	valuesSlice := []string{}
	for _, field := range fields {
		valuesSlice = append(valuesSlice, "'"+params[field]+"'")
	}
	values := strings.Join(valuesSlice, ",")
	sqlQueryString = sqlQueryString + fmt.Sprintf("(%v)", values)
	os.WriteFile("query.sql", []byte(sqlQueryString), 0777)
	defer exec.Command("rm", "query.sql").Run()

	command := fmt.Sprintf("cat query.sql | sqlite3 database/%v.db", tableName)
	err := exec.Command("bash", "-c", command).Run()
	if err != nil {
		fmt.Println("error:", err)
	}
}
