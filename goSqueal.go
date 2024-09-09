package goSqueal

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
)

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

func GetTableFields(tableName string) {
	directory := "./database"
	path := directory + "/" + tableName + ".db"
	fmt.Println(path)
}

func CreateTableEntry(tableName string) {

}
