package main

import (
	"flag"
	"fmt"
	"os"
	"passlocker/internal/locker"
)

func main() {
	locker := locker.Locker{
		Key:      "test",
		Locked:   false,
		Elements: []locker.Element{},
	}
	locker.Connect()
	defer locker.Disconnect()
	fmt.Println("Hello World!")
	setPassFlag := flag.Bool("set", false, "Add new Password")
	getPassFlag := flag.Bool("get", false, "Get Password")
	flag.Parse()
	passName := ""
	if *setPassFlag && !*getPassFlag {
		passValue := ""
		fmt.Println("We are setting Password")
		fmt.Print("Password Name: ")
		fmt.Scan(&passName)
		fmt.Print("Password Value: ")
		fmt.Scan(&passValue)
		if passName != "" && passValue != "" {
			locker.AddElement(passName, passValue)
			fmt.Printf("Added new Password %s\n", passName)
		} else {
			fmt.Println("You need to add both password name '-n' and value '-p'")
			os.Exit(1)
		}
	} else if *getPassFlag && !*setPassFlag {
		fmt.Print("Password Name: ")
		fmt.Scan(&passName)
		fmt.Println("We are getting Password")
		data := locker.GetElement(passName)
		fmt.Printf("Your Password: %s\n", data)
	} else {
		fmt.Println("Please use only one of 'get' or 'set' flags")
		os.Exit(1)
	}
}
