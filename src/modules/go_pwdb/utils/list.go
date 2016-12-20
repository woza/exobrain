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
	conf,err := config.CoreParser(fake_argv)
	if  err != nil {
		fmt.Println("Failed to load config: ",err)
		return
	}
	fmt.Print("Enter database password: ")
	raw_pw, err := terminal.ReadPassword(int(syscall.Stdin))
	fmt.Print("\n")
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
	tags := db.GetAll()
	fmt.Println("Known tags:")
	for _,t := range tags {
		fmt.Println(t)
	}
}
