package write

import (
	"log"
	"os"
)

//func WriteInFile(str []byte) {
//	err := os.WriteFile("users.txt", str, 0777)
//	if err != nil {
//		panic(err.Error())
//	}
//}

func WriteInFile(str []byte) {
	file, err := OpenFile("users.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer CloseFile(file)
	file.Write(str)
}

func OpenFile(fileName string) (*os.File, error) {
	return os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
}

func CloseFile(file *os.File) {
	file.Close()
}
