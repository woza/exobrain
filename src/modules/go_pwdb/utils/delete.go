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
	fake_argv := []string{
		"-conf", os.Args[1],
	}
	tag := strings.TrimSpace(os.Args[2])

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
	
	conf,err = db.Load(conf)
	if  err != nil {
		fmt.Println("Failed to load database: ",err)
		return
	}

	db.Remove(tag)
	err = db.Save(conf)
	if  err != nil {
		fmt.Println("Failed to save database: ",err)
		return
	}
	fmt.Println("Password added")
}

