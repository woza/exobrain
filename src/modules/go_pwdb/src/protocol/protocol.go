package protocol

import (
	"encoding/binary"
	"errors"
	"io"
        "fmt"
)

const (
	CMD_QUERY_ALL = iota
	CMD_TRIGGER = iota
	CMD_DISPLAY_SHOW = iota
	CMD_EXIT = iota
)

const (
	STATUS_OK = iota
	STATUS_FAIL = iota
)

const (
	CMD_DISPLAY_EXIT = iota
)
type Display_Request_Trigger struct {
	Encoding uint32
	Password string
}

type Input_Request_Header struct{
	size uint32
	cmd uint32
}

type Input_Request_Nil struct{}

type Input_Request_Query struct{
	header Input_Request_Header
}

type Input_Response_Query struct{
	Encoding uint32
	Tags []string
}

type Input_Response struct {
	Status uint32
}

type Input_Request_Exit struct{
	header Input_Request_Header
}

type Input_Request_Trigger struct{
	header Input_Request_Header
	Encoding uint32
	payload []byte
}

type Input_Request interface{
	GetPayload() []byte
}

func (self *Input_Request_Header) UnmarshalBinary(data[] byte) error{
	self.size = binary.BigEndian.Uint32( data[:4] )
	self.cmd  = binary.BigEndian.Uint32( data[4:8] )
	return nil
}

func (self *Input_Request_Trigger) UnmarshalBinary(data[] byte) error{
	// Assumes header already unmarshalled
	self.Encoding = binary.BigEndian.Uint32( data[:4] )
	self.payload = data[4:]
	return nil
}

func (self Input_Response_Query) MarshalBinary() ([]byte, error){
	body_sz := 0
	for _,k := range self.Tags {
		body_sz += 4 + len(k)
	}
	status_len := 4
	encoding_len := 4
	msg_len_len := 4
	count_len := 4
	msg_len := msg_len_len + status_len + encoding_len + count_len + body_sz
	ret := make([]byte, msg_len)
	off := 0
	binary.BigEndian.PutUint32(ret[off:msg_len_len], uint32(msg_len))
	off = off + msg_len_len	
	binary.BigEndian.PutUint32(ret[off:off+status_len], STATUS_OK)
	off = off + status_len
	binary.BigEndian.PutUint32(ret[off:off + encoding_len], self.Encoding)
	off = off + encoding_len
	binary.BigEndian.PutUint32(ret[off:off + count_len], uint32(len(self.Tags)))
	off = off + count_len
	for _,k := range self.Tags {
		binary.BigEndian.PutUint32(ret[off:off+4], uint32(len(k)))
		off = off + 4
		copy(ret[off:off+len(k)], []byte(k))
		off = off + len(k)
	}
	return ret,nil
}

func (self Display_Request_Trigger) MarshalBinary() ([]byte, error){
	raw_pw := []byte(self.Password)
	pw_len := len(raw_pw)
	len := 12 + pw_len
	ret := make([]byte, len)
	binary.BigEndian.PutUint32(ret[0:4], uint32(len))
	binary.BigEndian.PutUint32(ret[4:8], CMD_DISPLAY_SHOW)
	binary.BigEndian.PutUint32(ret[8:12], self.Encoding)
	copy(ret[12:len], raw_pw)
	return ret,nil
}
		
func (self Input_Request_Query) GetPayload() []byte {
	return []byte{}
}

func (self Input_Request_Exit) GetPayload() []byte {
	return []byte{}
}

func (self Input_Request_Nil) GetPayload() []byte {
	return []byte{}
}

func (self Input_Request_Trigger) GetPayload() []byte {
	return self.payload
}

func Fetch_Input_Request( src io.Reader ) (Input_Request, error){
	head_buff := make([]byte, 8)
	n,err := io.ReadFull(src, head_buff)
	if err != nil{
		return Input_Request_Nil{},err
	}
	fmt.Println("Fetch_Input_Request read count: ",n)
	head := Input_Request_Header{}
	head.UnmarshalBinary(head_buff)
	if head.cmd == CMD_QUERY_ALL{
		return Input_Request_Query{head},nil
	}
	if head.cmd == CMD_EXIT{
		return Input_Request_Exit{head},nil
	}
	if head.cmd == CMD_TRIGGER{		
		payload := make([]byte, head.size-4)
		fmt.Println("Payload size ",len(payload))
		_,err := io.ReadFull(src, payload)
		if err != nil{
			return Input_Request_Nil{},err
		}
		ret := Input_Request_Trigger{}
		ret.header = head
		ret.UnmarshalBinary(payload)
		return ret,nil
	}
	return Input_Request_Nil{},errors.New("Unrecognised command")	
}

func (self Input_Response)Put( dest io.Writer ) error{
	var status_buff[4]byte
	binary.BigEndian.PutUint32(status_buff[:], self.Status)
	_,err := dest.Write( status_buff[:] )
	return err
}

func (self Input_Response_Query) Put( dest io.Writer ) error{
	msg,err := self.MarshalBinary()
	if err != nil{
		return err
	}
	_,err = dest.Write(msg)
	return err
}

func (self Display_Request_Trigger) Put( link io.ReadWriter ) error{
	msg,err := self.MarshalBinary()
	if err != nil{
		return err
	}

	_,err = link.Write(msg)
	if err != nil{
		return err
	}
	var res_buff[4]byte
	_,err = io.ReadFull(link, res_buff[:])
	if err != nil{
		return err
	}
	code := binary.BigEndian.Uint32(res_buff[:])
	if code == STATUS_FAIL{
		return errors.New("Display returned Failure")
	}

	return nil
}
