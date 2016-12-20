package config
import (
	"flag"
	"os"
	"bufio"
	"strings"
	"golang.org/x/crypto/ssh/terminal"
	"syscall"
	"fmt"
)

type Credentials struct{
	TLS_key_path string
	TLS_cert_path string
	TLS_ca_path string
}

type Config struct{
	Path string
	Salt []byte
	Password string
	From_ui Credentials
	To_display Credentials
	Display_Address string
	Accept_Address string
	Display_Hostname string
}

func New() (Config, error){
	ret,err := CoreParser(os.Args[1:])
	if err != nil{
		return Config{},err
	}
	fmt.Print("Enter database password: ")
	raw_pw, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil{
		return Config{},err
	}

	ret.Password = strings.TrimSpace(string(raw_pw))
	return ret,nil
}

/* Core function, exported for automation of tests */
func CoreParser( src []string ) (Config, error){
	ret := Config{"", []byte{}, "",
		Credentials{"", "", ""},
		Credentials{"", "", ""},
		"", "", ""}
	conf_path := ""
	args := flag.NewFlagSet("gopwdb", flag.ContinueOnError)
	args.StringVar(&conf_path, "conf", "/usr/local/etc/gopwdb.conf",
		"Path to config file")
	args.Parse(src)
	return parse_config_file( ret, conf_path )
}
			
func parse_config_file( conf Config, path string )(Config, error){
	
	fmt.Println("Parsing config file ",path);
	handle,err := os.Open(path)
	if err == nil {
		defer handle.Close()
		scanner := bufio.NewScanner(handle)
		for scanner.Scan(){
			bits := strings.Split(scanner.Text(),"=")
			if len(bits) < 2{
				continue
			}
			key := strings.TrimSpace(bits[0])
			value := strings.TrimSpace(bits[1])
			if key == "db" && conf.Path == ""{
				conf.Path = value
			}
			if key == "ui_key" && conf.From_ui.TLS_key_path == ""{
				conf.From_ui.TLS_key_path = value
			}
			if key == "ui_cert" && conf.From_ui.TLS_cert_path == ""{
				conf.From_ui.TLS_cert_path = value
			}
			if key == "ui_ca" && conf.From_ui.TLS_ca_path == ""{
				conf.From_ui.TLS_ca_path = value
			}
			if key == "display_key" && conf.To_display.TLS_key_path == ""{
				conf.To_display.TLS_key_path = value
			}
			if key == "display_cert" && conf.To_display.TLS_cert_path == ""{
				conf.To_display.TLS_cert_path = value
			}
			if key == "display_ca" && conf.To_display.TLS_ca_path == ""{
				conf.To_display.TLS_ca_path = value
			}
			if key == "display_address" &&
				conf.Display_Address == ""{
				conf.Display_Address = value
			}
			if key == "accept" &&
				conf.Accept_Address == ""{
				conf.Accept_Address = value
			}
			if key == "display_name" &&
				conf.Display_Hostname == ""{
				conf.Display_Hostname = value
			}
			
		}
	}
	return conf,nil
}
	
