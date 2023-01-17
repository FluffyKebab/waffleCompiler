package main

import (
	"compiler/parser"
	"compiler/wasmCompiler"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) <= 1 {
		fmt.Println("no file given")
		os.Exit(1)
	}

	fileData, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	syntaxTree, err := parser.Parse(string(fileData))
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	byteCode, err := wasmCompiler.Compile(syntaxTree)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	err = os.WriteFile("main.wasm", byteCode, 0644)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func Compile(input string) ([]byte, error) {
	syntaxTree, err := parser.Parse(input)
	if err != nil {
		return []byte{}, err
	}

	return wasmCompiler.Compile(syntaxTree)
}
