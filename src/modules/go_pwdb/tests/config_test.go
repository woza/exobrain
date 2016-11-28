package tests

import (
	"testing"
	"config"
	"os"
)

func TestConfigCmdLine( t *testing.T ){
	fake_argv := []string{
		"-db", "/path/to/database",
		"-salt", "deadbeef",
		"-key", "a_key_file.pem",
		"-cert", "a_cert_file.pem",
		"-ca", "ca_root_file.pem",
		"-display", "127.0.0.1",
		"-accept", "0.0.0.0",
		"-dhost", "nosuchhost",
	}

	conf,err := config.CoreParser( fake_argv )
	if err != nil{
		t.Log("Failed to parse inputs:",err)
		t.FailNow()
	}
	eval_config( t, conf )
}

func TestConfigFile( t *testing.T ){
	handle,err := os.OpenFile("config_test.in",
		os.O_RDWR | os.O_CREATE | os.O_TRUNC, 0600 )
	if err != nil{
		t.Log("Failed to open test file:",err)
		t.FailNow()
	}

	handle.Write([]byte("db=/path/to/database\n"))
	handle.Write([]byte("salt=deadbeef\n"))
	handle.Write([]byte("key=a_key_file.pem\n"))
	handle.Write([]byte("cert=a_cert_file.pem\n"))
	handle.Write([]byte("ca=ca_root_file.pem\n"))
	handle.Write([]byte("display_address=127.0.0.1\n"))
	handle.Write([]byte("accept=0.0.0.0\n"))
	handle.Write([]byte("display_name=nosuchhost\n"))
	handle.Close()

	fake_argv := []string{"-conf", "config_test.in"}
	conf,err := config.CoreParser( fake_argv )
	if err != nil{
		t.Log("Failed to parse inputs:",err)
		t.FailNow()
	}

	eval_config( t, conf )
}

func TestPrecedenceA( t *testing.T ){
	handle,err := os.OpenFile("config_test.in",
		os.O_RDWR | os.O_CREATE | os.O_TRUNC, 0600 )
	if err != nil{
		t.Log("Failed to open test file:",err)
		t.FailNow()
	}

	handle.Write([]byte("db=FILE/path/to/database\n"))
	handle.Write([]byte("salt=deadbeef\n"))
	handle.Write([]byte("key=FILEa_key_file.pem\n"))
	handle.Write([]byte("cert=a_cert_file.pem\n"))
	handle.Write([]byte("ca=FILEca_root_file.pem\n"))
	handle.Write([]byte("display_address=127.0.0.1\n"))
	handle.Write([]byte("accept=FILE0.0.0.0\n"))
	handle.Write([]byte("display_name=nosuchhost\n"))
	handle.Close()
	
	fake_argv := []string{
		"-conf", "config_test.in",
		"-db", "/path/to/database",
		"-key", "a_key_file.pem",
		"-ca", "ca_root_file.pem",
		"-accept", "0.0.0.0",
	}
	conf,err := config.CoreParser( fake_argv )
	if err != nil{
		t.Log("Failed to parse inputs:",err)
		t.FailNow()
	}

	eval_config( t, conf )
}

func TestPrecedenceB( t *testing.T ){
	handle,err := os.OpenFile("config_test.in",
		os.O_RDWR | os.O_CREATE | os.O_TRUNC, 0600 )
	if err != nil{
		t.Log("Failed to open test file:",err)
		t.FailNow()
	}

	handle.Write([]byte("db=/path/to/database\n"))
	handle.Write([]byte("salt=FILEdeadbeef\n"))
	handle.Write([]byte("key=a_key_file.pem\n"))
	handle.Write([]byte("cert=FILEa_cert_file.pem\n"))
	handle.Write([]byte("ca=ca_root_file.pem\n"))
	handle.Write([]byte("display_address=FILE127.0.0.1\n"))
	handle.Write([]byte("accept=0.0.0.0\n"))
	handle.Write([]byte("display_name=FILEnosuchhost\n"))
	handle.Close()
	
	fake_argv := []string{
		"-conf", "config_test.in",
		"-salt", "deadbeef",
		"-cert", "a_cert_file.pem",
		"-display", "127.0.0.1",
		"-dhost", "nosuchhost",
	}

	conf,err := config.CoreParser( fake_argv )
	if err != nil{
		t.Log("Failed to parse inputs:",err)
		t.FailNow()
	}

	eval_config( t, conf )
}

func eval_config( t *testing.T, conf config.Config ){
	if conf.Path != "/path/to/database" {
		t.Log("conf.Path was '%v' expected '%v'",
			conf.Path, "/path/to/database")
		t.FailNow()
	}
	if conf.TLS_key_path != "a_key_file.pem" {
		t.Log("conf.Path was '%v' expected '%v'",
			conf.TLS_key_path, "a_key_file.pem")
		t.FailNow()
	}
	if conf.TLS_cert_path != "a_cert_file.pem" {
		t.Log("conf.Path was '%v' expected '%v'",
			conf.TLS_cert_path, "a_cert_file.pem")
		t.FailNow()
	}
	if conf.TLS_ca_path != "ca_root_file.pem" {
		t.Log("conf.Path was '%v' expected '%v'",
			conf.TLS_ca_path, "ca_root_file.pem")
		t.FailNow()
	}
	if conf.Display_Address != "127.0.0.1" {
		t.Log("conf.Path was '%v' expected '%v'",
			conf.Display_Address, "127.0.0.1")
		t.FailNow()
	}
	if conf.Accept_Address != "0.0.0.0" {
		t.Log("conf.Path was '%v' expected '%v'",
			conf.Accept_Address, "0.0.0.0")
		t.FailNow()
	}
	if conf.Display_Hostname != "nosuchhost" {
		t.Log("conf.Path was '%v' expected '%v'",
			conf.Display_Hostname, "nosuchhost")
		t.FailNow()
	}
}
