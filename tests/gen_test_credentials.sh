#!/bin/bash

function prep_CA()
{
    mkdir -p demoCA/certs
    mkdir -p demoCA/newcerts
    mkdir -p demoCA/crl
    touch demoCA/index.txt
    #openssl req -x509 -subj /C=AU/ST=ACT/L=Canberra/O=Exobrain/OU=Testing/CN=CA/ -newkey rsa:4096 -keyout root.key -nodes -sha256 -config `pwd`/openssl.cnf -days 3650 -set_serial $RANDOM  -out root.crt -batch
}

function generate_cert() # config file, name, extension_section
{
    conf=$1
    name=$2
    ext=$3
    openssl req -new -subj /C=AU/ST=ACT/L=Canberra/O=Exobrain/OU=Testing/CN=$name/ -newkey rsa:2048 -keyout $name.key -nodes -sha256 -config $conf -days $[365*2] -out $name.req -batch
    echo "REQ DONE"
    printf "%08x" $RANDOM > demoCA/serial
    openssl ca -config $conf -in $name.req -out $name.crt -cert root.crt -keyfile root.key -batch -extensions $ext

    openssl pkcs12 -export -password pass:password -in $name.crt -inkey $name.key -out $name.pfx
}

function generate_certs()
{
    sed -e's/<<CLIENT_IP>>/192.168.0.104/g' openssl.cnf > client.cnf
    generate_cert client.cnf ui_to_server client_cert
    generate_cert client.cnf server_to_display client_cert

    sed -e's/<<SERVER_IP>>/192.168.0.104/g' openssl.cnf > server.cnf
    generate_cert server.cnf server_from_ui server_cert
    generate_cert server.cnf display_from_server server_cert


    sed -e's/<<BOTH_IP>>/192.168.0.104/g' openssl.cnf > both.cnf
    generate_cert both.cnf ui both_cert
    generate_cert both.cnf server both_cert
    generate_cert both.cnf display both_cert
}

prep_CA
generate_certs
