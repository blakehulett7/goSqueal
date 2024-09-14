package goSqueal

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"strings"
)

var defaultOpenPermissions int = 0777

func CheckForTable(tableName string) {
	command := fmt.Sprintf("cat init/%v.sql | sqlite3 database.db", tableName)
	err := exec.Command("bash", "-c", command).Run()
	if err != nil {
		fmt.Println("table creation error:", err)
	}
}

func DropTable(tableName string) {
	sqlQueryString := fmt.Sprintf("DROP TABLE %v;", tableName)
	os.WriteFile("query.sql", []byte(sqlQueryString), fs.FileMode(defaultOpenPermissions))
	defer exec.Command("rm", "query.sql").Run()
	command := "cat query.sql | sqlite3 database.db"
	err := exec.Command("bash", "-c", command).Run()
	if err != nil {
		fmt.Println("error:", err)
	}
}

func GetTableFields(tableName string) []string {
	sqlQueryString := fmt.Sprintf("SELECT name FROM pragma_table_info('%v');", tableName)
	os.WriteFile("query.sql", []byte(sqlQueryString), fs.FileMode(defaultOpenPermissions))
	defer exec.Command("rm", "query.sql").Run()

	command := "cat query.sql | sqlite3 database.db"
	fieldsData, err := exec.Command("bash", "-c", command).Output()
	if err != nil {
		fmt.Println("get table fields error:", err)
	}
	fields := string(fieldsData)
	fieldsSlice := strings.Split(fields, "\n")
	fieldsSlice = fieldsSlice[:len(fieldsSlice)-1]
	return fieldsSlice
}

func GetTableEntry(tableName string, id string) map[string]string {
	fields := GetTableFields(tableName)
	sqlQueryString := fmt.Sprintf("SELECT * FROM %v WHERE id = '%v';", tableName, id)
	os.WriteFile("query.sql", []byte(sqlQueryString), fs.FileMode(defaultOpenPermissions))
	defer exec.Command("rm", "query.sql").Run()
	command := "cat query.sql | sqlite3 database.db"
	entryData, err := exec.Command("bash", "-c", command).Output()
	if err != nil {
		fmt.Println("get table entry error:", err)
	}
	entry := string(entryData)
	entry = strings.Replace(entry, "\n", "", 1)
	values := strings.Split(entry, "|")
	if len(fields) != len(values) {
		fmt.Println("HEY: fields and values slices are of different lengths!")
	}
	result := map[string]string{}
	for i := 0; i < len(values); i++ {
		result[fields[i]] = values[i]
	}
	if result[fields[0]] == "" {
		return nil
	}
	return result
}

func CreateTableEntry(tableName string, params map[string]string) {
	sqlQueryString := fmt.Sprintf("INSERT INTO %v VALUES", tableName)
	fields := GetTableFields(tableName)
	valuesSlice := []string{}
	for _, field := range fields {
		valuesSlice = append(valuesSlice, "'"+params[field]+"'")
	}
	values := strings.Join(valuesSlice, ",")
	sqlQueryString = sqlQueryString + fmt.Sprintf("(%v);", values)
	os.WriteFile("query.sql", []byte(sqlQueryString), 0777)
	defer exec.Command("rm", "query.sql").Run()

	command := "cat query.sql | sqlite3 database.db"
	err := exec.Command("bash", "-c", command).Run()
	if err != nil {
		fmt.Println("create entry error:", err)
	}
}

func DeleteTableEntry(tableName string, id string) {
	sqlQueryString := fmt.Sprintf("DELETE FROM %v WHERE id = '%v';", tableName, id)
	os.WriteFile("query.sql", []byte(sqlQueryString), fs.FileMode(defaultOpenPermissions))
	defer exec.Command("rm", "query.sql").Run()
	command := "cat query.sql | sqlite3 database.db"
	err := exec.Command("bash", "-c", command).Run()
	if err != nil {
		fmt.Println("delete entry error:", err)
	}
}
