package main

import(
	"crypto/rand"
	"os"
	"fmt"
	"strings"
	"bufio"
	"config"
	"db"
	"golang.org/x/crypto/ssh/terminal"
	"syscall"
	"errors"
)

type Param struct{
	prompt string
	key string
}

func main(){
	conf_path,err := get_input("Enter location for new configuration file: ")
	if err != nil{
		fmt.Println("Failed to fetch validated input: ",err)
		return
	}
	handle,err := os.Create(conf_path)
	if err != nil{
		fmt.Println("Failed to open configuration file: ",err)
		return
	}
	defer handle.Close()

	params := []Param{
		{"Enter path for database: ", "db"},
		{"Enter key for talking to UI component: ", "ui_key"},
		{"Enter certificate for talking to assert identity to UI component: ", "ui_cert"},
		{"Enter certificate to validate identity of UI component: ", "ui_ca"},
		{"Enter key for talking to display component: ", "display_key"},
		{"Enter certificate for talking to assert identity to display component: ",
			"display_cert"},
		{"Enter certificate to validate identity of display component: ", "display_ca"},
		{"Enter address:port of display server: ", "display_address"},
		{"Enter adddress:port on which to accept connections from UI component: ", "accept"},
		{"Enter server name of the display server: ", "display_name"},
	}
	var val	= ""
	var db_path = "";
	for _,p := range params{
		val,err = get_input(p.prompt)
		if err != nil{
			fmt.Println("Failed to fetch validated input: ",err)
			return
		}
		_, err = handle.WriteString(p.key + "=" + val + "\n")
		if err != nil{
			fmt.Println("Failed to write configuration to file: ",err)
			return
		}
		if p.key == "db"{
			db_path = val
		}
	}

	salt := make( []byte, 32 )
	n,err := rand.Read( salt )	
	if n != len(salt) || err != nil{
		fmt.Println("Failed to generate whole random salt: ",err)
		return
	}
	db_pw,err := get_pw()
	if err != nil{
		fmt.Println("Failed to read password: ",err)
		return
	}
	conf := config.Config{db_path, salt, db_pw,
		config.Credentials{"", "", ""},
		config.Credentials{"", "", ""},
		"", "", ""}
	/* Save a new database file, empty apart from metadata */
	db.Save(conf)
	
	fmt.Println("Configuration file written to " + conf_path)
}

func get_input( prompt string ) (string, error){
	fmt.Print(prompt)
	src := bufio.NewReader(os.Stdin)
	ret,err := src.ReadString('\n')
	return strings.TrimSpace(ret),err
}
				
func get_pw() (string, error){
	fmt.Print("Enter database password: ")
	raw_pw, err := terminal.ReadPassword(int(syscall.Stdin))
	fmt.Print("\n")
	if err != nil{
		return "",err
	}
	fmt.Print("Confirm database password: ")	
	raw_pw2, err := terminal.ReadPassword(int(syscall.Stdin))
	fmt.Print("\n")
	if err != nil{
		return "",err
	}

	pw := strings.TrimSpace(string(raw_pw))
	pw2 := strings.TrimSpace(string(raw_pw2))
	if pw != pw2 {
		return "", errors.New("Passwords did not match")
	}
	return pw,nil
}
