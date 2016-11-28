package main

import(
	"db"
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	"os"
	"strings"
	"syscall"
	"config"
)

func main(){
	fake_argv := []string{"-db", os.Args[1]}
	tag := strings.TrimSpace(os.Args[2])
	var pw = ""
	if len(os.Args) > 3 {
		pw = strings.TrimSpace(os.Args[3])
	}else{
		fmt.Print("Enter password: ")
		raw_pw, err := terminal.ReadPassword(int(syscall.Stdin))
		if err != nil{
			fmt.Println("Failed to read password")
			return
		}
		pw = strings.TrimSpace(string(raw_pw))
	}
	conf,err := config.CoreParser(fake_argv)
	if  err != nil {
		fmt.Println("Failed to load config: ",err)
		return
	}
	fmt.Print("Enter database password: ")
	raw_pw, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil{
		fmt.Println("Failed to read password")
		return
	}
	conf.Password = strings.TrimSpace(string(raw_pw))
	
	err = db.Load(conf)
	if  err != nil {
		fmt.Println("Failed to load database: ",err)
		return
	}

	_,absent := db.Get(tag)
	if absent == nil {
		fmt.Println("Tag already in database: ",tag)
		return
	}

	db.Put(tag, pw)
	err = db.Save(conf)
	if  err != nil {
		fmt.Println("Failed to save database: ",err)
		return
	}
	fmt.Println("Password added")
}
