package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)


type Column struct{
	Name string `yaml:"name"`
	Type string `yaml:"type"`
	Null string `yaml:"null"`
	Key string `yaml:"key"`
	Default string `yaml:"default"`
	Extra string `yaml:"extra"`
}

type Table struct{
	Name string	`yaml:"name"`
	Columns []Column `yaml:"columns"`
}

type database struct{
	Name string `yaml:"name"`
	Tables []Table `yaml:"tables"`
}


func (db *database)createDB ()(ret string){
	ret = fmt.Sprint("CREATE DATABASE " + db.Name + ";")
	return ret
}

func (table *Table)getColumns ()(res string){
	res = ""
	c := " , "
	for i,column := range table.Columns{
		temp := ""
		s := " "
		temp += column.Name + s +
			 column.Type + s +
			 column.Null + s +
			 column.Key + s +
			 column.Default + s +
			 column.Extra
		if i != 0 {
			res += c
		}
		res += temp
	}
	return res
}

func (db *database)createTables ()(ret string){
	for _,table := range db.Tables{
		ret = fmt.Sprint("CREATE TABLE " + table.Name+ "("+ table.getColumns() +")")
	}
	return ret
}

func main() {

	// Create the database handle, confirm driver is present
	constSTR := "root:oQ=VZ9yq&ziPh!bT@tcp(localhost:8080)"
	dbname := ""
	db, _ := sql.Open("mysql", constSTR+"/"+dbname)
	defer db.Close()

	//// Connect and check the server versionuse
	var version string
	db.QueryRow("SELECT VERSION()").Scan(&version)
	fmt.Println("Connected to:", version)

	b,_ := ioutil.ReadFile("db.yaml")
	var f database
	erro := yaml.Unmarshal(b,&f)
	if erro != nil {
		fmt.Println(erro)
	}

	res := f.createDB()
	_ , err := db.Exec(res)
	if err != nil {
		fmt.Println(err)
	}
	db.Exec("USE "+ f.Name)
	res = f.createTables()
	_ , err = db.Exec(res)
	if err != nil {
		fmt.Println(err)
	}
}

