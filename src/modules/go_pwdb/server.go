package main

import ("fmt"
        "net"
	"db"
	"crypto/tls"
	"config"
	"protocol"
	"io/ioutil"
	"crypto/x509"
        )



const (
	ENCODE_ASCII = iota
	ENCODE_UTF8 = iota
)

func main() {
	conf,err := config.New()
	if  err != nil {
		fmt.Println("Failed to load config: ",err)
		return
	}
	conf,err = db.Load(conf)
	if  err != nil {
		fmt.Println("Failed to load database: ",err)
		return
	}
	ln, err := net.Listen("tcp", conf.Accept_Address)
	if err != nil {
		fmt.Println("Failed to create listening socket");
		return;
	}
	
	for {
		client, err := ln.Accept()
		if err != nil {
			fmt.Println("Failed to accept client");
		} else{
			handle_client( client, conf );
		}
	}
}

func handle_client( tcp_client net.Conn, conf config.Config ){
        client := upgrade_to_tls(tcp_client,conf)
	ok := protocol.Input_Response{ protocol.STATUS_OK }
	fail := protocol.Input_Response{ protocol.STATUS_FAIL }
	
	for{
		req,err := protocol.Fetch_Input_Request(client)
		if err != nil{
			return
		}
		switch in_msg := req.(type){
		case protocol.Input_Request_Query:
			response := protocol.Input_Response_Query{
				ENCODE_UTF8,
				db.GetAll(),
			}
			_  = response.Put(client)
		case protocol.Input_Request_Exit:
		case protocol.Input_Request_Trigger:			
			tag := string(in_msg.GetPayload())
			pw,err := db.Get(tag)
			if err != nil{
				fmt.Println("Failed to retrieve password from database: ",
					err)
				_ = fail.Put(client)
				break
			}
			display,err := connect_to_display(conf)
			if err != nil{
				fmt.Println("Failed to connect to display")
				_ = fail.Put(client)
				break
			}
			out_msg := protocol.Display_Request_Trigger{
				ENCODE_UTF8, pw,
			}
			err = out_msg.Put(display)
			if err != nil {
				fmt.Println("Display returned failure",err)
				_ = fail.Put(client)
			}else{
				_ = ok.Put(client)
			}
			display.Close()
		}
	}
}

func upgrade_to_tls( src net.Conn, conf config.Config ) *tls.Conn {
        tls_conf,err := prep_tls_config( conf.From_ui )
	if err != nil{
		fmt.Println("Failed to prep TLS config: ", err)
		return nil
	}

	return tls.Server( src, tls_conf)
}

func connect_to_display( conf config.Config ) (*tls.Conn, error){
        tls_conf,err := prep_tls_config( conf.To_display )
	if err != nil{
		fmt.Println("Failed to prep TLS config: ", err)
		return nil, err
	}

	tls_conf.ServerName = conf.Display_Hostname
	display,err := tls.Dial( "tcp", conf.Display_Address, tls_conf )
	if err != nil{
		fmt.Println("Failed to connect to display")
		fmt.Println(err)
	}
	return display,err
}

func prep_tls_config( cred config.Credentials ) (*tls.Config, error ){
	cert, err := tls.LoadX509KeyPair(
		cred.TLS_cert_path,
		cred.TLS_key_path)
	if err != nil{
		fmt.Println("Failed to read TLS credentials: ", err)
		return nil, err
	}

	pool := x509.NewCertPool()
	pem_certs, err := ioutil.ReadFile(cred.TLS_ca_path)
	ioutil.WriteFile("/tmp/CA.back", []byte(pem_certs), 0644)
	
	if err != nil{
		fmt.Println("Failed to read TLS CA file")
		return nil, err
	}
	ok := pool.AppendCertsFromPEM( pem_certs )
	if !ok {
		fmt.Println("Failed to append CA certs to pool")
		return nil,err
	}
	ioutil.WriteFile("/tmp/CA.subjects", pool.Subjects()[0], 0644)
	
	display_conf := &tls.Config{
		Certificates : []tls.Certificate{cert},
		RootCAs : pool,
		ClientCAs : pool,
		ClientAuth: tls.RequireAndVerifyClientCert,
		InsecureSkipVerify : false,
	};
	display_conf.BuildNameToCertificate()
	return display_conf,nil;
}	
