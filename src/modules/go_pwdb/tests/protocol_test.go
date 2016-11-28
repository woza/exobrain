package tests

import (
	"fmt"
	"testing"
	"protocol"
	"bytes"
	"encoding/binary"
)

func TestQueryInputFetch( t *testing.T ){
	var query_input[8]byte;
	binary.BigEndian.PutUint32(query_input[:4], 4)
	binary.BigEndian.PutUint32(query_input[4:8], protocol.CMD_QUERY_ALL)
	buff := bytes.NewBuffer( query_input[:] )
	msg,err := protocol.Fetch_Input_Request(buff)
	if err != nil{
		fmt.Println("Failed to fetch query input:",err)
		t.FailNow()
	}
	_,ok := msg.(protocol.Input_Request_Query)
	if !ok {
		fmt.Println("Incorrect type returned when fetching query input")
		t.FailNow()
	}

}

func TestExitInputFetch( t *testing.T ){
	var exit_input[8]byte;
	binary.BigEndian.PutUint32(exit_input[:4], 4)
	binary.BigEndian.PutUint32(exit_input[4:8], protocol.CMD_EXIT)
	buff := bytes.NewBuffer( exit_input[:] )
	msg,err := protocol.Fetch_Input_Request(buff)
	if err != nil{
		fmt.Println("Failed to fetch exit input:",err)
		t.FailNow()
	}
	_,ok := msg.(protocol.Input_Request_Exit)
	if !ok {
		fmt.Println("Incorrect type returned when fetching exit input")
		t.FailNow()
	}

}

func TestTriggerInputFetch( t *testing.T ){

	var trigger_input[13]byte;
	binary.BigEndian.PutUint32(trigger_input[:4], 9)
	binary.BigEndian.PutUint32(trigger_input[4:8], protocol.CMD_TRIGGER)
	copy(trigger_input[8:], []byte("hello"))
	buff := bytes.NewBuffer( trigger_input[:] )
	msg,err := protocol.Fetch_Input_Request(buff)
	if err != nil{
		fmt.Println("Failed to fetch trigger input:",err)
		t.FailNow()
	}
	_,ok := msg.(protocol.Input_Request_Trigger)
	if !ok {
		fmt.Println("Incorrect type returned when fetching trigger input")
		t.FailNow()
	}

	p := msg.GetPayload()
	ps := string(p)
	if ps != "hello"{
		fmt.Printf("Trigger message yielded invalid payload '%v'",ps)
		t.FailNow()
	}	
}


	

