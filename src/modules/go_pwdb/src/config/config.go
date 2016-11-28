package config
import (
	"flag"
	"os"
	"bufio"
	"strings"
	"golang.org/x/crypto/ssh/terminal"
	"syscall"
	"fmt"
	"encoding/hex"
)

type Config struct{
	Path string
	Salt []byte
	Password string
	TLS_key_path string
	TLS_cert_path string
	TLS_ca_path string
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
		"", "", "",
		"", "", ""}
	hex_salt := ""
	conf_path := ""
	args := flag.NewFlagSet("gopwdb", flag.ContinueOnError)
	args.StringVar(&conf_path, "conf", "/usr/local/etc/gopwdb.conf",
		"Path to config file")
	args.StringVar(&ret.Path, "db", "", "Path to database")
	args.StringVar(&hex_salt, "salt", "", "Salt for database password")
	// But password cannot be passed via command line	
	args.StringVar(&ret.TLS_key_path, "key", "",
		"TLS key file")
	args.StringVar(&ret.TLS_cert_path, "cert", "",
		"TLS cert file")
	args.StringVar(&ret.TLS_ca_path, "ca", "",
		"TLS CA file")
	args.StringVar(&ret.Display_Address, "display", "",
		"Display Address")
	args.StringVar(&ret.Accept_Address, "accept", "",
		"Address to accept connections on")
	args.StringVar(&ret.Display_Hostname, "dhost", "",
		"Hostname to validate display host")
	args.Parse(src)

	handle,err := os.Open(conf_path)
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
			if key == "db" && ret.Path == ""{
				ret.Path = value
			}
			if key == "salt" && hex_salt == ""{
				hex_salt = value
			}
			if key == "key" && ret.TLS_key_path == ""{
				ret.TLS_key_path = value
			}
			if key == "cert" && ret.TLS_cert_path == ""{
				ret.TLS_cert_path = value
			}
			if key == "ca" && ret.TLS_ca_path == ""{
				ret.TLS_ca_path = value
			}
			if key == "display_address" &&
				ret.Display_Address == ""{
				ret.Display_Address = value
			}
			if key == "accept" &&
				ret.Accept_Address == ""{
				ret.Accept_Address = value
			}
			if key == "display_name" &&
				ret.Display_Hostname == ""{
				ret.Display_Hostname = value
			}
			
		}
	}
	ret.Salt,err = hex.DecodeString(hex_salt)
	return ret,err
}
			
