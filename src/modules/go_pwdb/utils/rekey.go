package main

import(
	"db"
	"fmt"
	"input"
	"os"
	"config"
	"errors"
	"bufio"
)

func main(){
	fake_argv := []string{
		"-conf", os.Args[1],
	}
	src := bufio.NewReader(os.Stdin)
	pw,err := get_user_pw(src)
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
	conf.Password, err = input.Password(src)
	if err != nil{
		fmt.Println("Failed to read password")
		return
	}
	
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
	fmt.Println("Password changed")
}

func get_user_pw( src *bufio.Reader) (string,error){
	fmt.Print("Enter new password: ")
	pw, err := input.Password(src)
	if err != nil{
		return "", errors.New("Failed to read password")
	}
	fmt.Print("\nConfirm new password: ")
	confirm, err := input.Password(src)
	if err != nil{
		return "", errors.New("Failed to read password")
	}
	fmt.Print("\n")
	if (pw != confirm){
		return "", errors.New("Passwords don't match")
	}
	return pw,nil
}
