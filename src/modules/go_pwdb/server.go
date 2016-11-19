package main

import ("fmt"
        "net"
        "encoding/binary")

const CMD_QUERY_ALL = 0
const CMD_TRIGGER = 1
const CMD_EXIT = 2

const CMD_DISPLAY_SHOW = 0
const CMD_DISPLAY_EXIT = 1

const STATUS_OK = 0
const STATUS_FAIL = 1

func read_database() map[string]string {
    var ret map[string]string = make(map[string]string)
    ret["major"] = "bloodnok"
    ret["neddy"] = "seagoon"
    return ret
}

func main() {
    var db = read_database();
    ln, err := net.Listen("tcp", ":6512");
    if err != nil {
        fmt.Println("Failed to create listening socket");
        return;
    }
    for {
        client, err := ln.Accept();
        if err != nil {
            fmt.Println("Failed to accept client");
        } else{
            handle_client( client, db );
        }
    }
}

func handle_client( client net.Conn, db map[string]string ){
    var size_buff[4]byte;
    var cmd_buff[4]byte;

    for{
        client.Read( size_buff[:] )
        client.Read( cmd_buff[:] )
        var size = binary.BigEndian.Uint32( size_buff[:] )
	var cmd = binary.BigEndian.Uint32( cmd_buff[:] )
        if cmd == CMD_QUERY_ALL {
            fmt.Println("QUERY_ALL command");
            var status_buff[4]byte
            binary.BigEndian.PutUint32(status_buff[:], STATUS_OK)
	    client.Write( status_buff[:] )
            var tot_len uint32 = 0
            for k := range db {
                tot_len += 4 + uint32(len(k))
            }
            var size_buf[4]byte
            var encoding_buf[4]byte

	    binary.BigEndian.PutUint32( size_buf[:], tot_len )
            client.Write( size_buff[:] )
            for k := range db {
  	        binary.BigEndian.PutUint32( size_buf[:], uint32(len(k)) )
                client.Write( size_buff[:] )
                client.Write( []byte(k) )
            }
            continue
        }           
        if cmd == CMD_TRIGGER {            
            var payload = make([]byte, size, size)
            client.Read(payload)
            var tag = string(payload)
            fmt.Println("TRIGGER command")
            fmt.Println(tag)
            continue
       }
       if cmd == CMD_EXIT {            
           fmt.Println("EXIT command");
       }          
    }
}

