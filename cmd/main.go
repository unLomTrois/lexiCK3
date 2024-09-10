package main

import (
	// "ck3-parser/internal/app/linter"

	"ck3-parser/internal/app/lexer"
	"ck3-parser/internal/app/linter"
	"ck3-parser/internal/app/parser"
	"encoding/json"
	"io"
	"log"
	"os"
	"time"
)

func main() {
	start := time.Now()

	// Open file
	filepath := "data/0_elementary.txt"
	file, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	filecontent, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}
	log.Printf("Read took %s", time.Since((start)))
	read := time.Now()

	lexer := lexer.New(filecontent)
	tokenstream, err := lexer.Scan()
	if err != nil {
		panic(err)
	}
	log.Printf("Scan took %s", time.Since((read)))

	scan := time.Now()

	// err = SaveJSON(tokenstream, "tokenstream.json")
	// if err != nil {
	// 	panic(err)
	// }
	// log.Println("Tokenstream saved to tmp/tokenstream.json")

	parser := parser.New(tokenstream)
	parsetree := parser.Parse()
	log.Printf("Parse  took %s", time.Since((scan)))

	if err = SaveJSON(parsetree, "parsetree.json"); err != nil {
		panic(err)
	}
	log.Println("Parsed data saved to tmp/parsetree.json")

	// Lint file
	lintconfig := linter.LintConfig{
		IntendStyle:            linter.TABS,
		IntendSize:             4,
		TrimTrailingWhitespace: true,
		InsertFinalNewline:     true,
		CharSet:                "utf-8-bom",
		EndOfLine:              []byte("\r\n"),
	}

	linter := linter.New(parsetree, lintconfig)
	linter.Lint()
	if err = linter.Save("tmp/linted.txt"); err != nil {
		panic(err)
	}
	log.Println("Linted file saved to tmp/linted.txt")

	log.Printf("All took %s", time.Since(start))
}

func SaveJSON(data interface{}, filename string) error {
	err := os.MkdirAll("tmp", 0755)
	if err != nil {
		return err
	}

	filepath := "tmp/" + filename
	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	enc := json.NewEncoder(file)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", " ")
	if err := enc.Encode(data); err != nil {
		return err
	}

	return nil
}
