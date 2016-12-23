package main

import(
	"db"
	"fmt"
	"input"
	"os"
	"config"
	"bufio"
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
	conf.Password, err = input.Password(bufio.NewReader(os.Stdin))
	fmt.Print("\n")
	if err != nil{
		fmt.Println("Failed to read password")
		return
	}
	
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
