package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"net"
	"os"
	"time"
)

type connectionInfo struct {
	ip, port string
}

func main() {
	var command string

	for {
		fmt.Println(" -> Please enter 'run' for begin the game or '?' for help or 'exit' for out")
		fmt.Scan(&command)
		switch command {
		case "?":
			fmt.Println(" -> Some instructions")
			bufio.NewReader(os.Stdin).ReadBytes('\n')
		case "run":
			beginTheGame(nil)
			/*con, err := tryToConnect()
			if err != nil {
				fmt.Println(err)
				break
			}
			err = prepareTheGame(con)
			if err != nil {
				fmt.Println(err)
				break
			}*/
		case "exit":
			os.Exit(0)
		}
	}
}

func prepareTheGame(conInf connectionInfo) (err error) {
	var con net.Conn
	buf := bytes.Buffer{}
	buf.WriteString(conInf.ip)
	buf.WriteString(":")
	buf.WriteString(conInf.port)
	if conInf.ip == "" {
		listener, err := net.Listen("tcp", buf.String())
		if err != nil {
			return err
		}
		con, err = listener.Accept()
		if err != nil {
			return errors.New(" -> Cannot get connection")
		}
		defer listener.Close()
		defer con.Close()
	} else {
		con, err = net.DialTimeout("tcp", buf.String(), time.Second*3)
		if err != nil {
			return errors.New(" -> Cannot connect to " + buf.String())
		}
		defer con.Close()
	}
	beginTheGame(con)
	return
}

func beginTheGame(con net.Conn) {
	myMap := [][]rune{
		{'~', '~', '~', '~', '~', '~', '~', '~', '~'},
		{'~', '~', '~', '~', '~', '~', '~', '~', '~'},
		{'~', '~', '~', '~', '~', '~', '~', '~', '~'},
		{'~', '~', '~', '~', '~', '~', '~', '~', '~'},
		{'~', '~', '~', '~', '~', '~', '~', '~', '~'},
		{'~', '~', '~', '~', '~', '~', '~', '~', '~'},
		{'~', '~', '~', '~', '~', '~', '~', '~', '~'},
		{'~', '~', '~', '~', '~', '~', '~', '~', '~'},
		{'~', '~', '~', '~', '~', '~', '~', '~', '~'},
	}
	enemyMap := [][]rune{
		{'~', '~', '~', '~', '~', '~', '~', '~', '~'},
		{'~', '~', '~', '■', '~', '~', '~', '~', '~'},
		{'~', '~', '~', 'X', '~', '~', '~', '~', '~'},
		{'~', '~', '~', '~', '~', '~', '~', '~', '~'},
		{'~', '~', '~', '~', '~', '~', '~', '~', '~'},
		{'~', '~', '~', '~', '~', '~', '~', '~', '~'},
		{'~', '~', '~', '~', '~', '~', '~', '~', '~'},
		{'~', '~', '~', '~', '~', '~', '~', '~', '~'},
		{'~', '~', '~', '~', '~', '~', '~', '~', '~'},
	}

	// <--PLACE SHIPS-->

	placeSmallShips(&myMap)
	placeBigShip(2, 3, &myMap)
	placeBigShip(3, 2, &myMap)
	placeBigShip(4, 1, &myMap)
	for {
		buf := make([]byte, 100)
		n, err := con.Read(buf)
		if err != nil {
			fmt.Println(err)
			return
		}
		if string(buf[:n]) == "start" {
			fmt.Println("starting the game")
			time.Sleep(time.Second * 2)
			err = startGame(con)
			if err != nil
		}
		fmt.Println("Waiting...")
		time.Sleep(time.Second * 3)
		break
	}
	outMaps(myMap, enemyMap)
}

func startGame(con net.Conn) (err error) {

}

func placeBigShip(shipSize int, shipCount int, myMap *[][]rune) {
	var readBuf [4]byte
	copyMap := make([][]rune, len(*myMap))
	copy(copyMap, *myMap)

Label:
	for i := 0; i < shipCount; i++ {
		outMaps(*myMap, nil)
		fmt.Println("Enter first coordinates(example A1) for put", shipSize, "cub ship:")
		_, err := fmt.Scanf("%c%d", &readBuf[0], &readBuf[1])
		if err != nil {
			fmt.Println(err)
			i--
			continue
		}
		fmt.Println("Enter second coordinates(example A2) for put", shipSize, "cub ship:")
		_, err = fmt.Scanf("%c%d", &readBuf[2], &readBuf[3])
		if err != nil {
			fmt.Println(err)
			i--
			continue
		}
		if readBuf[0] != readBuf[2] && readBuf[1] != readBuf[3] {
			fmt.Println("Incorrect data, try again")
			i--
			continue
		}
		if abs(int(readBuf[1])-int(readBuf[3])) != byte(shipSize-1) && abs(int(readBuf[0])-int(readBuf[2])) != byte(shipSize-1) {
			fmt.Println("Incorrect ship size!")
			i--
			continue
		}

		if readBuf[0] == readBuf[2] && readBuf[1] < readBuf[3] {
			for j := 0; j < shipSize; j++ {
				if checkShipForValid(readBuf[1]+byte(j)-1, readBuf[0]-'A', *myMap) == false {
					fmt.Println("Ship cant be near other ship, try again!")
					i--
					continue Label
				}
			}
			for j := 0; j < shipSize; j++ {
				(*myMap)[readBuf[1]-1+byte(j)][readBuf[0]-'A'] = '■'
			}
		} else if readBuf[0] == readBuf[2] && readBuf[1] > readBuf[3] {
			for j := 0; j < shipSize; j++ {
				if checkShipForValid(readBuf[1]-byte(j)-1, readBuf[0]-'A', *myMap) == false {
					fmt.Println("Ship cant be near other ship, try again!")
					i--
					continue Label
				}
			}
			for j := 0; j < shipSize; j++ {
				(*myMap)[readBuf[1]-1-byte(j)][readBuf[0]-'A'] = '■'
			}
		} else if readBuf[1] == readBuf[3] && readBuf[0] < readBuf[2] {
			for j := 0; j < shipSize; j++ {
				if checkShipForValid(readBuf[1]-1, readBuf[0]-'A'+byte(j), *myMap) == false {
					fmt.Println("Ship cant be near other ship, try again!")
					i--
					continue Label
				}
			}
			for j := 0; j < shipSize; j++ {
				(*myMap)[readBuf[1]-1][readBuf[0]-'A'+byte(j)] = '■'
			}
		} else if readBuf[1] == readBuf[3] && readBuf[0] > readBuf[2] {
			for j := 0; j < shipSize; j++ {
				if checkShipForValid(readBuf[1]-1, readBuf[0]-'A'-byte(j), *myMap) == false {
					fmt.Println("Ship cant be near other ship, try again!")
					i--
					continue Label
				}
			}
			for j := 0; j < shipSize; j++ {
				(*myMap)[readBuf[1]-1][readBuf[0]-'A'-byte(j)] = '■'
			}
		}
	}
}

func abs(num int) byte {
	if num < 0 {
		return byte(-num)
	}
	return byte(num)
}

func placeSmallShips(myMap *[][]rune) {
	var readBuf [2]byte

	for i := 0; i < 4; i++ {
		outMaps(*myMap, nil)
		fmt.Println("Enter coordinates(example A1) for put small ship(1 cub):")
		_, err := fmt.Scanf("%c%d", &readBuf[0], &readBuf[1])
		if err != nil {
			fmt.Println(err)
			i--
			continue
		}
		if checkShipForValid(readBuf[1]-1, readBuf[0]-'A', *myMap) == false {
			fmt.Println("Ship cant be near other ship, try again!")
			i--
			continue
		}
		(*myMap)[readBuf[1]-1][readBuf[0]-'A'] = '■'
	}
}

func checkShipForValid(y byte, x byte, copyMap [][]rune) bool {
	if y < 0 || y > 9 || x < 0 || x > 9 {
		return false
	}
	if y != 0 && copyMap[y-1][x] == '■' {
		return false
	}
	if y != byte(len(copyMap)-1) && copyMap[y+1][x] == '■' {
		return false
	}
	if x != 0 && copyMap[y][x-1] == '■' {
		return false
	}
	if x != byte(len(copyMap[y])-1) && copyMap[y][x+1] == '■' {
		return false
	}
	return true
}

func outMaps(myMap [][]rune, enemyMap [][]rune) {
	fmt.Print("  ")
	for i := 0; i < 9; i++ {
		fmt.Printf("%c ", 'A'+i)
	}
	if enemyMap != nil {
		fmt.Print("\t  ")
		for i := 0; i < 9; i++ {
			fmt.Printf("%c ", 'A'+i)
		}
	}
	fmt.Println()
	for s := 0; s < cap(myMap); s++ {
		fmt.Print(s+1, "|")
		for _, t := range myMap[s] {
			fmt.Printf("%c|", t)
		}

		if enemyMap != nil {
			fmt.Print("\t")
			fmt.Print(s+1, "|")
			for _, t := range enemyMap[s] {
				fmt.Printf("%c|", t)
			}
		}
		fmt.Println()
	}
}

func tryToConnect() (con connectionInfo, err error) {
	fmt.Println(" -> Enter 'create' for create a game or 'connect' for connect to other player")
	if str, _, _ := bufio.NewReader(os.Stdin).ReadLine(); string(str) == "connect" {
		fmt.Println(" -> Please enter ip to connect or 'default' for connect this machine\n -> Enter return for return")
		fmt.Scan(&con.ip)
		if con.ip == "return" {
			err = errors.New(" -> returning...")
			return
		}
		if con.ip == "default" {
			con.ip = "127.0.0.1"
		}
	} else if string(str) != "create" {
		err = errors.New(" -> returning...")
		return
	}
	fmt.Println(" -> Please enter port or 'default' for use default port\n -> Enter return for return")
	fmt.Scan(&con.port)
	if con.port == "return" {
		err = errors.New(" -> returning...")
		return
	}
	if con.port == "default" {
		con.port = "8080"
	}
	return
}
