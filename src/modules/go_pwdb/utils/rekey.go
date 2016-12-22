package main

import(
	"db"
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	"os"
	"strings"
	"syscall"
	"config"
	"errors"
)

func main(){
	fake_argv := []string{
		"-conf", os.Args[1],
	}

	pw,err := get_user_pw()
	if  err != nil {
		fmt.Println("Failed to get passwords: ",err)
		return
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
	
	conf,err = db.Load(conf)
	if  err != nil {
		fmt.Println("Failed to load database: ",err)
		return
	}
	conf.Password = pw;
	err = db.Save(conf)
	if  err != nil {
		fmt.Println("Failed to save database: ",err)
		return
	}
	fmt.Println("Password added")
}

func get_user_pw() (string,error){
	fmt.Print("Enter new password: ")
	raw_pw, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil{
		return "", errors.New("Failed to read password")
	}
	fmt.Print("\nConfirm new password: ")
	raw_confirm, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil{
		return "", errors.New("Failed to read password")
	}
	fmt.Print("\n")
	pw := strings.TrimSpace(string(raw_pw))
	confirm := strings.TrimSpace(string(raw_confirm))
	if (pw != confirm){
		return "", errors.New("Passwords don't match")
	}
	return pw,nil
}
