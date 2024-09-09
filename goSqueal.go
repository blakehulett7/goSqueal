package goSqueal

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
)

func CheckForTable(tableName string) {
	path := fmt.Sprintf("./database/%d.db", tableName)
	_, err := os.Stat("./database/users.db")
	if !errors.Is(err, fs.ErrNotExist) {
		fmt.Println("Db exists")
		return
	}
	fmt.Println("Db does not exist, creating db")
	command := "cat init/users.sql | sqlite3 database/users.db"
	err = exec.Command("bash", "-c", command).Run()
	if err != nil {
		fmt.Println(err)
	}
}
