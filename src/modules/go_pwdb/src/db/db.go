package db

import (
	"os"
	"config"
	"errors"
	"crypto/aes"
	"crypto/rand"
	"crypto/cipher"
	"encoding/json"
	"crypto/sha256"
	"encoding/binary"
	"golang.org/x/crypto/pbkdf2"
)

var pw_map map[string]string = make( map[string]string )

func Load( config config.Config ) (config.Config, error) {
	existing, fetched_salt, err := fetch_decrypted( config );
	config.Salt = fetched_salt
	if err != nil {
		return config,err;
	}

	if existing != nil {
		err = json.Unmarshal(existing, &pw_map);
		if err != nil {
			return config,err;
		}
	}
	return config,nil
}

func Save( config config.Config ) error {
	bytes,err := json.Marshal(pw_map)
	if err != nil{
		return err
	}
	return put_encrypted( config, bytes )
}

func Get( tag string ) (string, error){
	pw,exists := pw_map[tag]
	if exists {
		return pw,nil
	}
	return "", errors.New("Specified tag does not exist in database")
}

func Put( tag string, pw string ) {
	pw_map[tag] = pw
}

func Remove( tag string ){
	delete(pw_map, tag)
}

func GetAll() []string {
	// Sets up a constant order - in Go, iteration order over
	//  maps is randomised
	var ret []string
	for k := range pw_map{
		ret = append(ret, k)
	}
	return ret
}

func put_encrypted( config config.Config, data []byte ) error {
	key := pbkdf2.Key(
		[]byte(config.Password),
		config.Salt,
		1000000, 256/8, sha256.New)

	block,err := aes.NewCipher( key )
	if err != nil {
		return err
	}
	
	keystream,err := cipher.NewGCM( block )
	if err != nil {
		return err
	}

	new_nonce := make( []byte, keystream.NonceSize() )
	_,err = rand.Read( new_nonce )
	if err != nil{
		return err
	}

	out_handle, err := os.OpenFile(config.Path,
		os.O_WRONLY | os.O_CREATE | os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer out_handle.Close()
	version := []byte{1}
	_,err = out_handle.Write( version )
	if err != nil{
		return err
	}
	salt_len := make([]byte, 4)
	binary.BigEndian.PutUint32(salt_len, uint32(len(config.Salt)))
	_,err = out_handle.Write( salt_len )
	if err != nil{
		return err
	}
 	
	_,err = out_handle.Write( config.Salt )
	if err != nil{
		return err
	}
 	_, err = out_handle.Write( new_nonce )
	if err != nil{
		return err
	}
	
	if len(data) > 0 {
		cipher_bytes := []byte{}
		extra_bytes := []byte{}
		cipher_bytes = keystream.Seal( cipher_bytes, new_nonce, data,
			extra_bytes )

		_, err = out_handle.Write( cipher_bytes )
	}
	return err
}
	
func fetch_decrypted( config config.Config ) ([]byte, []byte, error) {
	raw_bytes, err := get_file_contents(config.Path);
	if err != nil{
		return nil,nil,err;
	}

	if len(raw_bytes) == 0 {
		return []byte{},[]byte{},nil
	}
	off := 0
	version := raw_bytes[off]
	if version != 1{
		return nil, nil, errors.New("Unhandled database version")
	}
	off += 1
	salt_len := int(binary.BigEndian.Uint32(raw_bytes[off:off + 4]))
	off += 4
	config.Salt = raw_bytes[off:off+salt_len]
	off += salt_len
	
	key := pbkdf2.Key(
		[]byte(config.Password),
		config.Salt,
		1000000,
		256/8,
		sha256.New)

	block,err := aes.NewCipher( key )
	if err != nil {
		return nil,config.Salt,err
	}
	
	keystream,err := cipher.NewGCM( block )
	if err != nil {
		return nil,config.Salt,err
	}
	plain_bytes := []byte{}
	nonce_bytes := raw_bytes[off:off +keystream.NonceSize()]
	off += keystream.NonceSize()
	if off == len(raw_bytes){
		return nil, config.Salt,nil
	}
	cipher_bytes := raw_bytes[off:]
	extra_bytes := []byte{}
	plain_bytes,err = keystream.Open( plain_bytes, nonce_bytes,
		cipher_bytes, extra_bytes )
	if err != nil {
		return nil,config.Salt,err
	}
	return plain_bytes,config.Salt,nil
}	

func get_file_contents( path string ) ([]byte,error) {
	info,err := os.Stat(path);
	if err != nil {		
		return nil,err;
	}

	var ret = make([]byte, info.Size());
	handle,err := os.Open(path);
	if err != nil {
		return nil,err
	}
	defer handle.Close()

	n,err := handle.Read(ret);
	if err != nil {
		return nil,err
	}

	if int64(n) != info.Size() {
		return nil, errors.New("Bad number of bytes returned")
	}

	return ret,nil;
}


