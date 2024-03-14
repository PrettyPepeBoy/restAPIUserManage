package read

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func ReadStdin() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	input = strings.TrimSpace(input)
	return input, nil
}

func ReadButton() int {
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	input = strings.TrimSpace(input)
	i, err := strconv.Atoi(input)
	if err != nil {
		log.Fatal(err)
	}
	return i
}

func ReadId() (uint32, bool) {
	fmt.Println("input user's id")
	id, ok := readUint32()
	return id, ok
}

func ReadCash() (uint32, bool) {
	fmt.Println("input your cash")
	cash, ok := readUint32()
	return cash, ok
}

func readUint32() (uint32, bool) {
	object, err := ReadStdin()
	if err != nil {
		fmt.Println("incorrect input")
		return 0, false
	}
	identificator, err := strconv.Atoi(object)
	if err != nil {
		panic(err)
	}
	return uint32(identificator), true
}

func CheckString(input string) bool {
	if len(input) > 20 {
		fmt.Println("your name is to big")
		return false
	}
	for _, a := range input {
		if (a < 'A' || a > 'Z') && (a < 'a' || a > 'z') {
			return false
		}
	}
	return true
}

func addZero(input string) string {
	placeHolder := []byte(input)
	number, err := strconv.Atoi(input)
	if err != nil {
		log.Fatal(err)
	}
	if placeHolder[0] != 48 && number < 10 {
		return "0" + input
	}
	return input
}

func CorrectData() (time.Time, bool) {
	var data string
	fmt.Println("please, input the year of your birth")
	input, err := ReadStdin()
	if err != nil {
		log.Fatal(err)
	}
	input = addZero(input)
	data = data + input

	fmt.Println("please, input the month of your birth")
	input, err = ReadStdin()
	if err != nil {
		log.Fatal(err)
	}
	input = addZero(input)
	data = data + input

	fmt.Println("please, input the day of your birth")
	input, err = ReadStdin()
	if err != nil {
		log.Fatal(err)
	}
	input = addZero(input)
	data = data + input

	layout := "20060102"
	t, err := time.Parse(layout, data)
	if err != nil {
		return t, false
	}
	return t, true
}
