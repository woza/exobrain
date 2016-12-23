package main

import(
	"db"
	"fmt"
	"os"
	"strings"
	"config"
	"errors"
	"input"
	"bufio"
)

func main(){
	fake_argv := []string{
		"-conf", os.Args[1],
	}
	tag := strings.TrimSpace(os.Args[2])
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

func get_user_pw( src *bufio.Reader) (string,error){
	fmt.Print("Enter password: ")
	pw, err := input.Password(src)
	if err != nil{
		return "", errors.New("Failed to read password")
	}
	fmt.Print("\nConfirm password: ")
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
