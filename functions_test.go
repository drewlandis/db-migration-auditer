package main

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

func TestDatabaseFunctions(t *testing.T) {
	fmt.Println("test")

	// Read all filenames in ./migrations. (ioutil sorts by filename)
	migrationsDir := "/d3/testHD/hyperdark/appliance/uidb/migrations"
	files, err := ioutil.ReadDir(migrationsDir)
	if err != nil {
		panic(err)
	}

	// print filenames
	for _, v := range files {
		fmt.Println(v.Name())
	}

	// Read contents of a file
	fmt.Println("\n")
	//fileNumber := 21 // 11 up  -> 1 func with var names
	//fileNumber := 95 // 48 up  -> 1 func with NO var names
	fileNumber := 161 // 81 up  -> 2 funcs with var names
	//fileNumber := 181 // 91 up   -> lots of funcs with var names
	fileName := files[fileNumber].Name()
	fmt.Println(fileName)

	direction := getDirection(fileName)
	fmt.Println(direction)

	migrationNumber, err := getMigrationNumber(fileName)
	if err != nil {
		panic(err)
	}
	fmt.Println(migrationNumber)

	contents, err := ioutil.ReadFile(migrationsDir + "/" + files[fileNumber].Name())
	if err != nil {
		panic(err)
	}
	//fmt.Println(string(contents))

	// Find any updated functions...
	//fmt.Println("\n")
	// CREATE OR REPLACE FUNCTION seen(threat TEXT, epoch BIGINT)
	//
	//r := strings.NewReplacer("<h1>", "", "</h1>", "")
	//fmt.Println(r.Replace("<h1>Hello World!</h1>"))
	//

	//funcName, funcDef := FindFunc(string(contents), "CREATE OR REPLACE FUNCTION", "LANGUAGE plpgsql;")
	//
	continueLoop := true
	fileStr := string(contents)
	var funcName string
	//var funcDef string
	//upList := make(map[int]map[string]string) // map[migrationNumber][funcName][funcDef]
	for continueLoop {
		funcName, _, fileStr = FindFunc(fileStr, "CREATE OR REPLACE FUNCTION", "LANGUAGE plpgsql;")
		//fmt.Println(funcDef)
		//fmt.Println("\n")
		fmt.Println(funcName)
		//fmt.Println("\n")

		if funcName == "" {
			continueLoop = false
		}
	}

	// TODO: write values to upList map
}

// FindFunc searches for the first instance of a function and returns the
// functionName, functionDefinition, and everything in the original string AFTER
// the function definition
func FindFunc(str string, start string, end string) (string, string, string) {
	// Find location of function start
	s := strings.Index(str, start)
	if s == -1 {
		return "", "", ""
	}

	// Find location of function end
	e := strings.Index(str, end)
	e += len(end)

	// Using start and end locations, get function definition
	funcDef := str[s:e]

	// Get function name with full signature
	fNameStart := s + len(start) + 1
	fNameEnd := fNameStart + strings.Index(str[fNameStart:], ")") + 1
	funcNameWithVarNames := str[fNameStart:fNameEnd]

	// If present, trim extra spaces
	funcNameWithVarNames = strings.TrimSpace(funcNameWithVarNames)

	// If present, remove carriage returns
	funcNameWithVarNames = regexp.MustCompile(`\r?\n`).ReplaceAllString(funcNameWithVarNames, " ")

	// If present, remove space before closing paren
	funcNameWithVarNames = regexp.MustCompile(` \)`).ReplaceAllString(funcNameWithVarNames, ")")

	// If present, remove named variables from signature
	fOpenParen := strings.Split(funcNameWithVarNames, "(")
	fSpace := strings.Split(fOpenParen[1], " ")
	funcNameWithOutVarNames := fOpenParen[0] + "("
	for _, v := range fSpace {
		if strings.Index(v, ",") != -1 || strings.Index(v, ")") != -1 {
			funcNameWithOutVarNames += v
		}
	}

	// Remove function definition from original string
	fileStringWithoutFuncDef := str[e:]

	return strings.ToLower(funcNameWithOutVarNames), funcDef, fileStringWithoutFuncDef
}

func getDirection(fileName string) string {
	if strings.Index(fileName, "up") != -1 {
		return "up"
	} else {
		return "down"
	}
}

func getMigrationNumber(fileName string) (int, error) {
	// example fileName: 0081_create_threat_traffic_tables_and_functions.up.sql

	// Remove leading zeros
	fileName = strings.TrimLeft(fileName, "0")

	// Remove everything after first "_"
	migrationStr := fileName[:strings.Index(fileName, "_")]

	// Convert string to int
	migrationInt, err := strconv.Atoi(migrationStr)

	return migrationInt, err
}
