package main

import 
(
	"fmt"
	"flag"
	"os"
	"github.com/syndtr/goleveldb/leveldb"
	"log"
)

func main(){
	db, err := leveldb.OpenFile("tmp/test.db", nil)
	if err != nil {
		log.Fatal("Yikes!")
	}
	defer db.Close()
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
		if passName != "" && passValue != ""{
			err = db.Put([]byte(passName), []byte(passValue), nil)
			if err != nil{
				log.Fatal("Yikes!")
			}
			fmt.Printf("Added new Password %s\n", passName)
		}else{
			fmt.Println("You need to add both password name '-n' and value '-p'")
			os.Exit(1)
		}
	}else if *getPassFlag && !*setPassFlag {
		fmt.Print("Password Name: ")
		fmt.Scan(&passName)
		fmt.Println("We are getting Password")
		data, err := db.Get([]byte(passName), nil)
		if err != nil{
			log.Fatal("Yikes!")
		}

		fmt.Printf("Your Password: %s\n", string(data))
	}else{
		fmt.Println("Please use only one of 'get' or 'set' flags")
		os.Exit(1)
	}
}
