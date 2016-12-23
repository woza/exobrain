package input

import(
	"golang.org/x/crypto/ssh/terminal"
	"syscall"
	"bufio"
	"strings"
)

func Password( src *bufio.Reader ) (string, error){
	stdin_fd := int(syscall.Stdin)
	if terminal.IsTerminal(stdin_fd){
		raw,err := terminal.ReadPassword(stdin_fd)
		if err != nil{
			return "",err
		}
		return strings.TrimSpace(string(raw)),err
	}
	// Stdin piped in from file or other mechanism
	raw,err := src.ReadString('\n')
	if err != nil{
		return "",err
	}
	return strings.TrimSpace(raw),nil	
}
