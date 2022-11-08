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
			con, err := tryToConnect()
			if err != nil {
				fmt.Println(err)
				break
			}
			err = prepareTheGame(con)
			if err != nil {
				fmt.Println(err)
				break
			}
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
	_ = enemyMap

	// <--PLACE SHIPS-->

	placeSmallShips(&myMap)
	placeBigShip(2, 3, &myMap)
	placeBigShip(3, 2, &myMap)
	placeBigShip(4, 1, &myMap)
	for {
		_, err := con.Write([]byte("start"))
		if err != nil {
			fmt.Println(err)
			return
		}
		buf := make([]byte, 100)
		fmt.Println("Waiting for answer from opponent...")
		err = con.SetReadDeadline(time.Now().Add(time.Second * 1080))
		if err != nil {
			fmt.Println(err)
			return
		}
		n, err := con.Read(buf)
		if err != nil {
			fmt.Println(err)
			return
		}
		if string(buf[:n]) == "start" {
			fmt.Println("starting the game")
			time.Sleep(time.Second * 2)
			err = startGame(con, myMap, enemyMap)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}

func startGame(con net.Conn, myMap [][]rune, enemyMap [][]rune) (err error) {
	var x, y byte
	shipCount := 19
	readBuf := make([]byte, 2)
	log := new(log)

	for {
		outMaps(myMap, enemyMap, log)
		err = con.SetReadDeadline(time.Now().Add(time.Second * 90))
		if err != nil {
			return
		}
		fmt.Println("Enter coordinate(like A2) to shoot")
		_, err = fmt.Scanf("%c%d", &readBuf[0], &readBuf[1])
		if err != nil {
			fmt.Println("Incorrect coordinate, try again")
			continue
		}
		x = readBuf[0] - 'A'
		y = readBuf[1] - 1
		if y < 0 || y > 9 || x < 0 || x > 9 {
			fmt.Println("Incorrect coordinate, try again")
			continue
		}
		_, err = con.Write([]byte{x, y})
		if err != nil {
			fmt.Println("Incorrect coordinate, try again")
			continue
		}
		_, err = con.Read(readBuf)
		if err != nil {
			return err
		}
		_, err = con.Write(checkForHit(readBuf[0], readBuf[1], &myMap, &shipCount))
		if err != nil {
			return err
		}
		_, err = con.Read(readBuf)
		if err != nil {
			return err
		}
		if markTheHit(x, y, readBuf[0], &enemyMap, log) == true {
			if shipCount == 0 {
				fmt.Println("Draw in the game. No winners;)")
				return err
			} else {
				fmt.Println("Congratulations, you are win!!!")
				return err
			}
		}
		if shipCount == 0 {
			fmt.Println("Sorry, you lose(((")
			return
		}
	}
}

func markTheHit(x byte, y byte, b byte, enemyMap *[][]rune, logs *log) bool {
	switch b {
	case 0:
		//fmt.Printf("Miss for coordinate %c %d\n", x+'A', y+1)
		logs.add(fmt.Sprintf("Miss for coordinate %c %d", x+'A', y+1))
		//fmt.Sprintf()
		(*enemyMap)[y][x] = '○'
	case 1:
		//fmt.Printf("Destryed ship for coordinate %c %d\n", x+'A', y+1)
		logs.add(fmt.Sprintf("Destryed ship for coordinate %c %d", x+'A', y+1))
		(*enemyMap)[y][x] = 'X'
	case 2:
		//fmt.Printf("Destryed ship for coordinate %c %d\n", x+'A', y+1)
		logs.add(fmt.Sprintf("Destryed ship for coordinate %c %d", x+'A', y+1))
		(*enemyMap)[y][x] = 'X'
	case 4:
		return true
	}
	return false
}

// 0 - miss, 1 - destroy, 2 - hit, 4 - destroy all ships
func checkForHit(x byte, y byte, myMap *[][]rune, shipCount *int) []byte {
	switch (*myMap)[x][y] {
	case '~', 'X':
		return []byte{0}
	case '■':
		*shipCount--
		(*myMap)[y][x] = 'X'
		if *shipCount == 0 {
			return []byte{4}
		} else if (x != 0 && (*myMap)[y][x-1] != '■') && (x != byte(len((*myMap)[0])-1) && (*myMap)[y][x+1] != '■') && (y != 0 && (*myMap)[y-1][x] != '■') && (y != byte(len(*myMap)-1) && (*myMap)[y+1][x] != '■') {
			return []byte{2}
		} else {
			return []byte{1}
		}
	default:
		return []byte{4}
	}
}

func placeBigShip(shipSize int, shipCount int, myMap *[][]rune) {
	var readBuf [4]byte
	copyMap := make([][]rune, len(*myMap))
	copy(copyMap, *myMap)

OuterCycle:
	for i := 0; i < shipCount; i++ {
		outMaps(*myMap, nil, nil)
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
					continue OuterCycle
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
					continue OuterCycle
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
					continue OuterCycle
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
					continue OuterCycle
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
		outMaps(*myMap, nil, nil)
		fmt.Println("Enter coordinates(example A1) for put small ship(1 cub):")
		_, err := fmt.Scanf("%c%d", &readBuf[0], &readBuf[1])
		if err != nil {
			fmt.Println(err)
			i--
			continue
		}
		if checkShipForValid(readBuf[1]-1, readBuf[0]-'A', *myMap) == false {
			fmt.Println("Ship cant be placed here, try again!")
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

func outMaps(myMap [][]rune, enemyMap [][]rune, logs *log) {
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

		if logs != nil {
			fmt.Printf("\t|\t%s", logs.get(s))
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
