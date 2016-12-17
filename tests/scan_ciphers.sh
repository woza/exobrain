#!/bin/bash


TODO=$(openssl ciphers | sed -e's/:/ /g')
HOST=192.168.0.104
PORT=8765
for c in $TODO; do
    used_cipher=$(openssl s_client -connect $HOST:$PORT -cert server_to_display.crt -key server_to_display.key -CAfile root.crt -tls1_2 -cipher $c 2> /dev/null | grep Cipher.is | awk '{print $NF}')
    if [ "$used_cipher" == "(NONE)" ]; then
	echo "$c - NO"
    else
	echo "$c - YES ($used_cipher)"
    fi
done
    
