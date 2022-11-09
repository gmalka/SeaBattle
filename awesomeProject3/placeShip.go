package main

import "fmt"

func placeBigShip(shipSize int, shipCount int, myMap *[][]rune) {
	var readBuf [4]byte
	copyMap := make([][]rune, len(*myMap))
	copy(copyMap, *myMap)

OuterCycle:
	for i := 0; i < shipCount; i++ {
		outMaps(*myMap, nil, nil)
		fmt.Println("Enter first coordinates(example A1A2) for put", shipSize, "cub ship:")
		_, err := fmt.Scanf("%c%d%c%d", &readBuf[0], &readBuf[1], &readBuf[2], &readBuf[3])
		if err != nil {
			fmt.Println(err)
			i--
			continue
		}
		if readBuf[0] >= 'a' && readBuf[0] <= 'z' {
			readBuf[0] -= 32
		}
		if readBuf[2] >= 'a' && readBuf[2] <= 'z' {
			readBuf[2] -= 32
		}
		/*fmt.Println("Enter second coordinates(example A2) for put", shipSize, "cub ship:")
		_, err = fmt.Scanf("%c%d", &readBuf[2], &readBuf[3])*/
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

func placeSmallShips(myMap *[][]rune, shipCount int) {
	var readBuf [2]byte

	for i := 0; i < shipCount; i++ {
		outMaps(*myMap, nil, nil)
		fmt.Println("Enter coordinates(example A1) for put small ship(1 cub):")
		_, err := fmt.Scanf("%c%d", &readBuf[0], &readBuf[1])
		if readBuf[0] >= 'a' && readBuf[0] <= 'z' {
			readBuf[0] -= 32
		}
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
	if copyMap[y][x] == '■' {
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
	fmt.Print("    YOUR MAP")
	if enemyMap != nil {
		fmt.Print("\t\t\tENEMY MAP")
	}
	fmt.Println()
	//TOP OUT ALPHABET
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
	if logs != nil {
		fmt.Print("\t\t\tFIGHT LOGS")
	}
	fmt.Println()
	for s := 0; s < cap(myMap); s++ {
		//OUT MAP AND NUMBERS
		fmt.Print(s+1, "|")
		for _, t := range myMap[s] {
			fmt.Printf("%c|", t)
		}
		fmt.Print(s + 1)

		if enemyMap != nil {
			fmt.Print("\t")
			fmt.Print(s+1, "|")
			for _, t := range enemyMap[s] {
				fmt.Printf("%c|", t)
			}
			fmt.Print(s + 1)
		}

		if logs != nil {
			fmt.Printf("\t|\t%s", logs.get(s))
		}
		fmt.Println()
	}

	//BOTTOM OUT ALPHABET
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
}
