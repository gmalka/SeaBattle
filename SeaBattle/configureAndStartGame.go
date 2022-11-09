package main

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"os"
	"time"
)

const (
	DeadlineForOpponentReadyToStart = 180
	DeadlineForRemakeGame           = 30
	DeadlineForEnemyTurn            = 120
)

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
	pregameConfigureAndStart(con)
	return
}

func pregameConfigureAndStart(con net.Conn) {
	for {
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
			{'~', '~', '~', '~', '~', '~', '~', '~', '~'},
			{'~', '~', '~', '~', '~', '~', '~', '~', '~'},
			{'~', '~', '~', '~', '~', '~', '~', '~', '~'},
			{'~', '~', '~', '~', '~', '~', '~', '~', '~'},
			{'~', '~', '~', '~', '~', '~', '~', '~', '~'},
			{'~', '~', '~', '~', '~', '~', '~', '~', '~'},
			{'~', '~', '~', '~', '~', '~', '~', '~', '~'},
			{'~', '~', '~', '~', '~', '~', '~', '~', '~'},
		}

		// <--PLACE SHIPS-->
		placeSmallShips(&myMap, 4)
		placeBigShip(2, 3, &myMap)
		placeBigShip(3, 2, &myMap)
		placeBigShip(4, 1, &myMap)
		outMaps(myMap, nil, nil)

		_, err := con.Write([]byte("start"))
		if err != nil {
			fmt.Println(err)
			return
		}
		buf := make([]byte, 100)
		fmt.Println("Waiting for answer from opponent...")
		err = con.SetReadDeadline(time.Now().Add(time.Second * DeadlineForOpponentReadyToStart))
		if err != nil {
			fmt.Println(err)
			return
		}
		n, err := con.Read(buf)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Waiting answer from opponent...")
		if string(buf[:n]) == "start" {
			fmt.Println("starting the game")
			time.Sleep(time.Second * 2)
			err = startGame(con, myMap, enemyMap) // Start the game
			if err != nil {
				fmt.Println(err)
				return
			}
			for {
				// Remake game after end or return to menu
				fmt.Println("Play Again? Y/N")
				n, err = os.Stdin.Read(buf)
				if err != nil {
					fmt.Println(err)
					return
				}
				if buf[0] == 'Y' {
					_, err = con.Write([]byte{1})
					if err != nil {
						fmt.Println(err)
						return
					}
					err = con.SetReadDeadline(time.Now().Add(time.Second * DeadlineForRemakeGame))
					if err != nil {
						fmt.Println(err)
						return
					}
					fmt.Println("Waiting answer from opponent...")
					_, err = con.Read(buf)
					if err != nil {
						fmt.Println("Opponent dont answer")
						return
					}
					if buf[0] == 1 {
						break
					}
				}
				if buf[0] == 'N' {
					return
				}
			}
		}
	}
}

func startGame(con net.Conn, myMap [][]rune, enemyMap [][]rune) (err error) {
	var x, y byte
	shipCount := countFoShips(&myMap)
	readBuf := make([]byte, 2)
	log := new(log)

	for {
		outMaps(myMap, enemyMap, log)
		err = con.SetReadDeadline(time.Now().Add(time.Second * DeadlineForEnemyTurn))
		if err != nil {
			return
		}
		fmt.Print("Coordinate(like A2) to shoot:")
		_, err = fmt.Scanf("%c%d", &readBuf[0], &readBuf[1])
		if readBuf[0] >= 'a' && readBuf[0] <= 'z' {
			readBuf[0] -= 32
		}
		if err != nil {
			log.add("Incorrect coordinate, try again")
			continue
		}
		x = readBuf[0] - 'A'
		y = readBuf[1] - 1
		if y < 0 || y > 9 || x < 0 || x > 9 || enemyMap[y][x] == '●' || enemyMap[y][x] == '❌' || enemyMap[y][x] == '❎' {
			log.add("Incorrect coordinate, try again")
			continue
		}
		_, err = con.Write([]byte{x, y})
		if err != nil {
			log.add("Incorrect coordinate, try again")
			continue
		}
		_, err = con.Read(readBuf)
		if err != nil {
			return err
		}
		_, err = con.Write(checkForHit(readBuf[0], readBuf[1], &myMap, &shipCount, log))
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

func countFoShips(myMap *[][]rune) int {
	count := 0
	for _, j := range *myMap {
		for _, i := range j {
			if i == '■' {
				count++
			}
		}
	}
	return count
}

// Mark hit or miss on an opponent
func markTheHit(x byte, y byte, b byte, enemyMap *[][]rune, logs *log) bool {
	switch b {
	case 0:
		logs.add(fmt.Sprintf("You -> Miss by coordinate %c%d", x+'A', y+1))
		(*enemyMap)[y][x] = '●'
	case 1:
		logs.add(fmt.Sprintf("You -> Destroy the ship by coordinate %c%d", x+'A', y+1))
		// Change symbol on enemy Map when ship destroy
		(*enemyMap)[y][x] = '❎'
		for y != 0 && (*enemyMap)[y-1][x] == '❌' {
			(*enemyMap)[y-1][x] = '❎'
			y--
		}
		for y != byte(len(*enemyMap)-1) && (*enemyMap)[y+1][x] == '❌' {
			(*enemyMap)[y+1][x] = '❎'
			y++
		}
		for x != 0 && (*enemyMap)[y][x-1] == '❌' {
			(*enemyMap)[y][x-1] = '❎'
			x--
		}
		for x != byte(len((*enemyMap)[y])-1) && (*enemyMap)[y][x+1] == '❌' {
			(*enemyMap)[y][x+1] = '❎'
			x++
		}
	case 2:
		logs.add(fmt.Sprintf("You -> Hit the ship by coordinate %c%d", x+'A', y+1))
		(*enemyMap)[y][x] = '❌'
	case 4:
		return true
	}
	return false
}

// Check is enemy hit your ship and return code for opponent
// 0 - miss, 1 - destroy, 2 - hit, 4 - destroy all ships
func checkForHit(x byte, y byte, myMap *[][]rune, shipCount *int, logs *log) []byte {
	switch (*myMap)[y][x] {
	case '~', '❌':
		logs.add(fmt.Sprintf("Enemy -> Missed your ship by coordinate %c%d", x+'A', y+1))
		return []byte{0}
	case '■':
		*shipCount--
		(*myMap)[y][x] = '❌'
		if *shipCount == 0 {
			return []byte{4}
		} else if (x != 0 && (*myMap)[y][x-1] == '■') || (x != byte(len((*myMap)[0])-1) && (*myMap)[y][x+1] == '■') || (y != 0 && (*myMap)[y-1][x] == '■') || (y != byte(len(*myMap)-1) && (*myMap)[y+1][x] == '■') {
			logs.add(fmt.Sprintf("Enemy -> Hit your ship by coordinate %c%d", x+'A', y+1))
			return []byte{2}
		} else {
			logs.add(fmt.Sprintf("Enemy -> Destroy your ship by coordinate %c%d", x+'A', y+1))
			return []byte{1}
		}
	default:
		return []byte{4}
	}
}
