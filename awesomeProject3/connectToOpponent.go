package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
)

type connectionInfo struct {
	ip, port string
}

func initUserMenu() {
	var command string

	for {
		fmt.Println(" -> Please enter 'run' for begin the game or '?' for help or 'exit' for out")
		fmt.Scan(&command)

		switch command {
		case "?":
			fmt.Println("Instruction:\n For begin the game follow instruction in the terminal.\n" +
				" SeaBattle is the classic game with one opponent.\n In first part you and your opponent place a ships on the board." +
				"\n Then You and your opponent mark a coordinate on the board, where you want to shoot." +
				"\n The winner is the one who destroy all enemy ships. If you and your opponent destroy all of each other's ships, there will be a draw." +
				"\n Notations:" +
				"\n Empty space: ~ " +
				"\n Damaged ship: ❌ " +
				"\n Destroyed ship: ❎ " +
				"\n Miss: ● " +
				"\n Your ship: ■ " +
				"\nPress Enter for continue")
			bufio.NewReader(os.Stdin).ReadBytes('\n')
		case "run":
			con, err := getConnectionInfo()
			if err != nil {
				fmt.Println(err)
				break
			}
			//Start the game here
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

func getConnectionInfo() (con connectionInfo, err error) {
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
