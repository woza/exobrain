#!/bin/bash

# This script will set up, run and tear down a simple client-server
# interaction test of the Go server.

function prep_CA()
{
    # Don't assume that the default OpenSSL config is still configured as default out-of-the-box.
    ca_dir=simple_test_ca
    rm -rf $ca_dir
    mkdir -p $ca_dir/certs
    mkdir -p $ca_dir/newcerts
    mkdir -p $ca_dir/crl
    touch $ca_dir/index.txt
    openssl req -x509 -config simple_openssl.cnf -subj /C=AU/ST=ACT/L=Canberra/O=Exobrain/OU=Testing/CN=CA/ -newkey rsa:4096 -keyout simple_test_root.key -nodes -sha256 -days 3650 -set_serial $RANDOM -out simple_test_root.crt -batch > /dev/null 2>&1
}

function generate_credential()
{
    name=$1
    openssl req  -config simple_openssl.cnf -new -subj /C=AU/ST=ACT/L=Canberra/O=Exobrain/OU=Testing/CN=$name/ -newkey rsa:2048 -keyout $name.key -nodes -sha256 -days $[365*2] -out $name.req -batch > /dev/null 2>&1
    printf "%08x" $RANDOM > $ca_dir/serial
    openssl ca  -config simple_openssl.cnf -in $name.req -out $name.crt -cert simple_test_root.crt -keyfile simple_test_root.key -batch > /dev/null 2>&1
}

function create_tls_credentials()
{
    echo -n "Creating TLS credentials..."
    prep_CA
    generate_credential simple_test_python
    generate_credential simple_test_go
    echo "done"
}

function spawn_go_server()
{
    echo -n "Configuring new database..."
    cat > simple_test.in <<EOF
simple_test.conf
simple_test.db
simple_test_go.key
simple_test_go.crt
simple_test_root.crt
simple_test_go.key
simple_test_go.crt
simple_test_root.crt
127.0.0.1:6677
127.0.0.1:7766
simple_test_python
simple_testpass
simple_testpass
EOF
    ./utils/initdb < simple_test.in > /dev/null
    echo -e "aaaa\naaaa\nsimple_testpass\n" | ./utils/put simple_test.conf alpha  > /dev/null
    echo -e "bbbb\nbbbb\nsimple_testpass\n" | ./utils/put simple_test.conf beta  > /dev/null
    echo "done"
    echo -n "Spawning new server..."
    echo -e "simple_testpass\n" | ./go_pwdb -conf simple_test.conf  > /dev/null &
    SERVER_PID=$!
    # Server takes a few seconds to spin up, wait for it
    sleep 5
    echo "done"
}

echo "Running simple test"
create_tls_credentials
spawn_go_server
echo "Running test client."
./test_cli.py --key simple_test_python.key \
	      --cert simple_test_python.crt \
	      --ca simple_test_root.crt \
	      --tag alpha --tag beta \
	      --pw alpha:aaaa --pw beta:bbbb \
	      --server-port 7766 --display-port 6677
kill -1 $SERVER_PID
