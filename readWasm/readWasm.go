package readWasm

import (
	"compiler/leb128"
	"compiler/wasmCompiler/code"
	"fmt"
	"os"
)

func GetFuncFromFile(fileName string, funcIndex int) ([]byte, error) {
	funcSection, err := GetAllFuncsFromFile(fileName)
	if err != nil {
		return []byte{}, err
	}

	numFuncsInFile, err := leb128.LEB128ToInt32(funcSection[0:])
	if err != nil {
		return []uint8{}, err
	}
	numBytesStoringNumFunctions := len(leb128.Int32ToULEB128(int32(numFuncsInFile)))

	if funcIndex >= numFuncsInFile {
		return []uint8{}, fmt.Errorf("Error getting wasm func %v from file %s: Function index does not exist in file given", funcIndex, fileName)
	}

	curIndex := numBytesStoringNumFunctions

	for i := 0; i < funcIndex; i++ {
		if curIndex >= len(funcSection) {
			return []uint8{}, fmt.Errorf("Error getting wasm func %v from file %s: Function skipping failed", funcIndex, fileName)
		}

		curIndex, err = skipFunction(funcSection, curIndex)
		if err != nil {
			return []byte{}, err
		}
	}

	functionLen, err := leb128.LEB128ToInt32(funcSection[curIndex:])
	if err != nil {
		return []byte{}, err
	}
	numBytesStoringFunctionLen := len(leb128.Int32ToULEB128(int32(functionLen)))

	if curIndex+functionLen-1+numBytesStoringFunctionLen >= len(funcSection) {
		return []uint8{}, fmt.Errorf("Error getting wasm func from file theIndex of the start the function plus the functionLen-1 plus numBytesStoringFunctionLen is greater than the length of the func section")
	}

	return funcSection[curIndex+numBytesStoringFunctionLen : curIndex+functionLen+numBytesStoringFunctionLen], nil
}

func GetAllFuncsFromFile(fileName string) ([]byte, error) {
	fileContent, err := os.ReadFile(fileName)
	if err != nil {
		return []uint8{}, err
	}

	return getSectionFromFile(fileContent, code.SECTION_CODE)
}

func getSectionFromFile(fileContent []byte, sectionCode byte) ([]byte, error) {
	var err error
	curIndex := skipMagic()

	for {
		if curIndex >= len(fileContent) {
			return []uint8{}, fmt.Errorf("Error getting wasm func: skipping sections failed")
		}

		if fileContent[curIndex] == sectionCode {
			break
		}

		curIndex, err = skipSection(fileContent, curIndex)
		if err != nil {
			return []byte{}, err
		}
	}

	numBytesInSection, err := leb128.LEB128ToInt32(fileContent[curIndex+1:])
	if err != nil {
		return []byte{}, err
	}

	numBytesStoringSectionLen := len(leb128.Int32ToULEB128(int32(numBytesInSection)))

	return fileContent[curIndex+numBytesStoringSectionLen+1 : curIndex+numBytesInSection+numBytesStoringSectionLen+1], nil
}

func skipMagic() int {
	return 8
}

func skipSection(funcSection []byte, curIndex int) (int, error) {
	sectionSize, err := leb128.LEB128ToInt32(funcSection[curIndex+1:])
	return curIndex + 1 + len(leb128.Int32ToULEB128(int32(sectionSize))) + sectionSize, err
}

//Assumes the curIndex is on function size part of code section
func skipFunction(funcSection []byte, curIndex int) (int, error) {
	functionSize, err := leb128.LEB128ToInt32(funcSection[curIndex:])
	return curIndex + len(leb128.Int32ToULEB128(int32(functionSize))) + functionSize, err
}
