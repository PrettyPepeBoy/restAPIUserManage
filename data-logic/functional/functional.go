package functional

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"tstUser/data-logic/read"
)

const userLenByte = 52

type User struct {
	Name    string
	Surname string
	Id      uint32
	Cash    uint32
	Year    uint16
	Month   uint8
	Day     uint8
}

type Users []User

func CreateUser() ([]byte, bool) {
	userSlc := make([]byte, 0, userLenByte)
	id := randInt(100000, 999999)
	userSlc = binary.BigEndian.AppendUint32(userSlc, id)

	for {
		fmt.Println("please, input your user name")
		input, err := read.ReadStdin()
		if err != nil {
			panic(err)
		}
		if read.CheckString(input) {
			data := make([]byte, 20)
			for i := range input {
				data[i] = input[i]
			}
			userSlc = append(userSlc, data...)
			break
		}
		fmt.Println("incorrect name, try again")
	}

	for {
		fmt.Println("please, input you user surname")
		input, err := read.ReadStdin()
		if err != nil {
			panic(err)
		}
		if read.CheckString(input) {
			data := make([]byte, 20)
			for i := range input {
				data[i] = input[i]
			}
			userSlc = append(userSlc, data...)
			break
		}
		fmt.Println("incorrect surname, try again")
	}

	for {
		t, check := read.CorrectData()
		if check {
			userSlc = binary.BigEndian.AppendUint16(userSlc, uint16(t.Year()))
			userSlc = append(userSlc, byte(t.Month()))
			userSlc = append(userSlc, byte(t.Day()))
			break
		}
		fmt.Println("incorrect data, try again")
	}

	for {
		cash, check := read.ReadCash()
		if check {
			userSlc = binary.BigEndian.AppendUint32(userSlc, cash)
			break
		}
		fmt.Println("incorrect cash, try again")
	}

	return userSlc, true
}

func GetUserName(name, surname string) {
	Usr := getData("users.txt")
	for i := 0; i < len(Usr); i++ {
		if Usr[i].Name == name && Usr[i].Surname == surname {
			fmt.Println(Usr[i])
		}
	}
}

func randInt(i, j int) uint32 {
	return uint32(i + rand.Intn(j-i))
}

func CheckUsr() (User, bool) {
	id, ok := read.ReadId()
	if !ok {
		fmt.Println("incorrect id, try again")
		CheckUsr()
	}
	if num := findUser(id); num != -1 {
		file, err := os.OpenFile("users.txt", os.O_RDONLY, 0644)
		if err != nil {
			panic(err)
		}
		file.Seek(int64((num-1)*userLenByte), 0)
		buf := make([]byte, userLenByte)
		file.Read(buf)
		usr := User{
			Id:      binary.BigEndian.Uint32(buf),
			Name:    string(buf[4:24]),
			Surname: string(buf[24:44]),
			Year:    binary.BigEndian.Uint16(buf[44:46]),
			Month:   buf[46],
			Day:     buf[47],
			Cash:    binary.BigEndian.Uint32(buf[48:52]),
		}
		return usr, true
	}
	return User{}, false
}

func getData(fileName string) Users {
	var usr Users
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	buf := make([]byte, userLenByte)
	for {
		_, err := reader.Read(buf)
		if err != nil {
			if err != io.EOF {
				panic(err)
			}
			break
		}
		Usr := User{
			Id:      binary.BigEndian.Uint32(buf),
			Name:    string(buf[4:24]),
			Surname: string(buf[24:44]),
			Year:    binary.BigEndian.Uint16(buf[44:46]),
			Month:   buf[46],
			Day:     buf[47],
			Cash:    binary.BigEndian.Uint32(buf[48:52]),
		}
		usr = append(usr, Usr)
	}
	return usr
}

func findUser(id uint32) int {
	file, err := os.OpenFile("users.txt", os.O_RDONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	buf := make([]byte, userLenByte)
	count := 1
	for {
		_, err := file.Read(buf)
		if err != nil {
			if err != io.EOF {
				panic(err)
			}
			break
		}
		if binary.BigEndian.Uint32(buf) == id {
			return count
		}
		count++
	}
	return -1
}

func SortByName() {
	Usr := getData("users.txt")
	sort(Usr, 0, len(Usr)-1, true)
	for i := 0; i < len(Usr); i++ {
		fmt.Println(Usr[i])
	}
}

func SortByData() {
	Usr := getData("users.txt")
	sort(Usr, 0, len(Usr)-1, false)
	for i := 0; i < len(Usr); i++ {
		fmt.Println(Usr[i])
	}
}

func sort(arr Users, low int, high int, name bool) {
	if name {
		if low < high {
			index := partitionName(arr, low, high)
			sort(arr, low, index-1, name)
			sort(arr, index+1, high, name)
		}
	} else {
		if low < high {
			index := partitionData(arr, low, high)
			sort(arr, low, index-1, name)
			sort(arr, index+1, high, name)
		}
	}

}

func partitionName(arr Users, low int, high int) int {
	i := low - 1
	pivot := arr[high]
	for j := low; j < high; j++ {
		if arr[j].Name < pivot.Name {
			i++
			arr[i], arr[j] = arr[j], arr[i]
		} else if arr[j].Name == pivot.Name {
			if arr[j].Surname < pivot.Surname {
				i++
				arr[i], arr[j] = arr[j], arr[i]
			}
		}
	}
	arr[high], arr[i+1] = arr[i+1], arr[high]
	return i + 1
}

func partitionData(arr Users, low int, high int) int {
	i := low - 1
	pivot := arr[high]
	for j := low; j < high; j++ {
		if arr[j].Year < pivot.Year {
			i++
			arr[i], arr[j] = arr[j], arr[i]
		} else if arr[j].Year == pivot.Year {
			if arr[j].Month < pivot.Month {
				i++
				arr[i], arr[j] = arr[j], arr[i]
			} else if arr[j].Month == pivot.Month {
				if arr[j].Day < pivot.Day {
					i++
					arr[i], arr[j] = arr[j], arr[i]
				}
			}
		}
	}
	arr[high], arr[i+1] = arr[i+1], arr[high]
	return i + 1
}

func DeleteUser() {
	id, check := read.ReadId()
	if !check {
		fmt.Println("incorrect id, try again")
		DeleteUser()
	}
	numUsr := findUser(id)
	if numUsr == -1 {
		fmt.Println("Usr with this id doesn't exist")
		return
	}
	file, err := os.OpenFile("users.txt", os.O_RDWR, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	buf := make([]byte, userLenByte)
	file.Seek(int64(numUsr*userLenByte), 0)
	var newContent []byte
	for {
		_, err := file.Read(buf)
		if err != nil {
			if err != io.EOF {
				panic(err)
			}
			break
		}
		newContent = append(newContent, buf...)
	}
	file.Seek(int64((numUsr-1)*userLenByte), 0)
	file.Write(newContent)
	fmt.Println(int64(numUsr-1)*userLenByte+int64(len(newContent)), "len delete")
	err = os.Truncate("users.txt", int64(numUsr-1)*userLenByte+int64(len(newContent)))
	if err != nil {
		panic(err)
	}
}

func BuyThing() {
	fmt.Println("how much money you want to spend")
	for {
		if spend, ok := read.ReadCash(); ok {
			if id, ok := read.ReadId(); ok {
				if buy(spend, id) {
					break
				}
			}
		}
	}
}

func SendCash(id1, id2, send uint32) {
	fmt.Println("how much money you want to send")
	for {
		if buy(send, id1) {
			file, err := os.OpenFile("users.txt", os.O_RDWR, 0644)
			if err != nil {
				panic(err)
			}
			defer file.Close()
			file.Seek(int64(userLenByte*(findUser(id2)-1)+48), 0)
			buf := make([]byte, 4)
			_, err = file.Read(buf)
			if err != nil {
				panic(err)
			}
			cash := binary.BigEndian.Uint32(buf)
			file.Seek(int64(userLenByte*(findUser(id2)-1)+48), 0)
			buf = make([]byte, 0, 4)
			buf = binary.BigEndian.AppendUint32(buf, cash+send)
			file.Write(buf)
			break
		}
	}
}

func buy(spend, id uint32) bool {
	file, err := os.OpenFile("users.txt", os.O_RDWR, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	file.Seek(int64(userLenByte*(findUser(id)-1)+48), 0)
	buf := make([]byte, 4)
	_, err = file.Read(buf)
	if err != nil {
		panic(err)
	}
	cash := binary.BigEndian.Uint32(buf)
	if cash < spend {
		fmt.Println("Not enough money")
		return false
	}
	file.Seek(int64(userLenByte*(findUser(id)-1)+48), 0)
	buf = make([]byte, 0, 4)
	buf = binary.BigEndian.AppendUint32(buf, cash-spend)
	file.Write(buf)
	return true
}
