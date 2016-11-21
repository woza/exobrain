package main

import (
     "crypto/aes"
     "crypto/cipher"
     "fmt"
     "os"
)

func main() {
     key := []byte{1,2,3,4,5,6,7,8,
                   9,0, 10, 11, 12, 13, 14, 15, 16,
                   17,18,19,20,21,22,23,24,
                   25,26,27,28,29,30,31}
     block,err := aes.NewCipher( key )
     if err != nil {
        fmt.Println("Failed to create block")
        return
     }

     keystream,err := cipher.NewGCM( block )
     if err != nil {
        fmt.Println("Failed to create GCM")
        return
     }
    
     handle, err := os.Open("ciphertext")
     if err != nil {
        fmt.Println("Failed to open input")
        return
     }
    
     defer handle.Close()

     cipher_bytes := []byte{};
     _,err = handle.Read(cipher_bytes)
     if err != nil {
        fmt.Println("Failed to read input")
        return
     }
    
     plain_bytes := []byte{}
     nonce_bytes := []byte{0,1,2,3,4,5,6,7,8,9,10,11}
     extra_bytes := []byte{44}
     plain_bytes,err = keystream.Open( plain_bytes, nonce_bytes, cipher_bytes, extra_bytes )
     if err != nil {
        fmt.Println("Failed to decrypt", err)
        return
     }
     
     out_handle, err := os.OpenFile("plaintext.txt", os.O_WRONLY | os.O_CREATE | os.O_TRUNC, 0600)
     if err != nil {
        fmt.Println("Failed to open output")
        return
     }
     out_handle.Write( plain_bytes )    
}




