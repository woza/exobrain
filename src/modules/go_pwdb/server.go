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
	err = db.Load(conf)
	if  err != nil {
		fmt.Println("Failed to load database: ",err)
		return
	}
	fmt.Println("Listening on ",conf.Accept_Address)
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
		        fmt.Println("Accepted client")
			handle_client( client, conf );
		}
	}
}

func handle_client( tcp_client net.Conn, conf config.Config ){
        client := upgrade_to_tls(tcp_client,conf)
	ok := protocol.Input_Response{ protocol.STATUS_OK }
	fail := protocol.Input_Response{ protocol.STATUS_FAIL }
	
	for{
		fmt.Println("Awaiting next request")
		req,err := protocol.Fetch_Input_Request(client)
		if err != nil{
			fmt.Println("Failed to receive input: ",err)
			return
		}
		fmt.Println("Received request")
		switch in_msg := req.(type){
		case protocol.Input_Request_Query:
			fmt.Println("QUERY_ALL command");
			response := protocol.Input_Response_Query{
				ENCODE_UTF8,
				db.GetAll(),
			}
			fmt.Println("Tags in response: ",response.Tags)
			_  = response.Put(client)
			fmt.Println("Finished putting response")
		case protocol.Input_Request_Exit:
			fmt.Println("EXIT command");
		case protocol.Input_Request_Trigger:			
			fmt.Println("TRIGGER command")
			tag := string(in_msg.GetPayload())
			fmt.Println(tag)
			pw,err := db.Get(tag)
			if err != nil{
				_ = fail.Put(client)
				break
			}
			display,err := connect_to_display(conf)
			if err != nil{
				_ = fail.Put(client)
				break
			}
			out_msg := protocol.Display_Request_Trigger{
				ENCODE_UTF8, pw,
			}
			err = out_msg.Put(display)
			if err != nil {
				fmt.Println("Display returned failure")
				_ = fail.Put(client)
			}else{
				fmt.Println("Display returned success")
				_ = ok.Put(client)
			}
			display.Close()
			fmt.Println("Trigger processing finished")
		}
	}
}

func upgrade_to_tls( src net.Conn, conf config.Config ) *tls.Conn {
        tls_conf,err := prep_tls_config( conf )
	if err != nil{
		fmt.Println("Failed to prep TLS config: ", err)
		return nil
	}

	return tls.Server( src, tls_conf)
}
	


func connect_to_display( conf config.Config ) (*tls.Conn, error){
        tls_conf,err := prep_tls_config( conf )
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

func prep_tls_config( conf config.Config ) (*tls.Config, error ){
	cert, err := tls.LoadX509KeyPair(
		conf.TLS_cert_path,
		conf.TLS_key_path)
	if err != nil{
		fmt.Println("Failed to read TLS credentials: ", err)
		return nil, err
	}

	pool := x509.NewCertPool()
	pem_certs, err := ioutil.ReadFile(conf.TLS_ca_path)
	fmt.Println("Loading CA certs from ",conf.TLS_ca_path)
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
	fmt.Println("Found subject count: ",len(pool.Subjects()))
	
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
