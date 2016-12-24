#!/usr/bin/python3

import struct
import time
import argparse
import ssl
import socket
import sys
import os
from concurrent.futures import ProcessPoolExecutor

class TestFail (Exception):
    pass

def get_all( src, todo ):
    ret = src.recv(todo)
    while len(ret) < todo:
        chunk = src.recv(todo-len(ret))
        if len(chunk) == 0:
            print ("Remote server closed mid-receive")
            raise TestFail()
        ret += chunk
    return ret

def decode(encoding, raw):
    encode_ascii = 0
    encode_utf8 = 1
    if encoding == encode_ascii:
        return str(raw, 'ascii')
    if encoding == encode_utf8:
        return str(raw, 'utf-8')
    print ("Received unknown encoding '%d'\n", encoding)
    raise TestFail()

def receive_password(config):
    sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)                
    display_server = ssl.wrap_socket(sock,
                                     config.key, config.cert, True,
                                     ssl.CERT_REQUIRED,
                                     ssl.PROTOCOL_TLSv1_2,
                                     config.ca)
    display_server.bind(('127.0.0.1', config.dport))
    display_server.listen(3)

    # Private, do not call directly
    remote,_ = display_server.accept()
    msg = get_all(remote, 4)
    sz = struct.unpack('>I', msg)[0]
    msg = get_all(remote,sz - 4)
    cmd,encoding = struct.unpack('>II', msg[:8])
    password = decode(encoding, msg[8:])
    status=0
    remote.send(struct.pack('>II', 4, status))
    remote.close()
    return password

class Server:
    '''
    This class represents a remote server implementing the exobrain
    protocol It performs the UI and Display roles of said protocol, to
    enable testing of the protocol as spoken by the server.  It is
    *NOT* intended for use with sensitive data or in a production
    environment.
    '''

    def __init__(self, config):
        self.async_server_pool = ProcessPoolExecutor(2)
        self.config = config

    def new_client(self):
        sock = socket.create_connection(('127.0.0.1', self.config.sport))
        return ssl.wrap_socket(sock,
                               self.config.key, self.config.cert, False,
                               ssl.CERT_REQUIRED,
                               ssl.PROTOCOL_TLSv1_2,
                               self.config.ca)


    def get_password(self, tag):
        token = self.async_server_pool.submit(receive_password, (self.config))
        # Send trigger command to server
        cmd=1
        
        payload=bytes(tag, 'utf-8')
        encode_utf8 = 1
        sz = 8 + len(payload)
        msg = struct.pack('>III', sz, cmd, encode_utf8)
        msg += payload
        dest = self.new_client()
        dest.send(msg)
        # Get trigger response
        msg = get_all(dest, 4)
        status = struct.unpack('>I', msg)[0]
        if status != 0:
            print ("Bad status %d received when trying to fetch password." % status)
            raise TestFail()
        timeout=5
        while not token.done():
            if timeout == 0:
                print ("Did not receive display message before timeout expired")
                raise TestFail()
            time.sleep(1)
            timeout - 1

        # Fetch password from our display-server thread
        pw = token.result()
        dest.close()
        return pw

    def list_tags(self):
        dest = self.new_client()
        cmd=0
        msg = struct.pack('>II', 4, cmd)
        dest.write(msg)
        sz = struct.unpack('>I', get_all(dest, 4))[0]
        msg = get_all(dest, sz-4)
        dest.close()

        status = struct.unpack('>I', msg[:4])[0]
        if status != 0:
            print ("Remote side indicated query operation failed")
            raise TestFail()

        encoding,count = struct.unpack('>II', msg[4:12])
        off=12
        tags=[]
        for i in range(count):
            nbyte = struct.unpack('>I', msg[off:off+4])[0]
            off += 4
            tags += [decode(encoding, msg[off:off+nbyte])]
            off += nbyte
        return tags


def check_tags(actual, expected):
    aset = set(actual)
    eset = set(expected)

    diff = aset.difference(eset)
    if len(diff) > 0:
        print ("Tags %s were unexpectedly received" % ','.join(diff))
        raise TestFail()
    diff = eset.difference(aset)
    if len(diff) > 0:
        print ("Tags %s were not returned and they should have been" % ','.join(diff))
        raise TestFail()
    
    
# Process args and build conf, expectations
parser = argparse.ArgumentParser(description="Exobrain Test CLI")
parser.add_argument("--key",  required=True)
parser.add_argument("--cert", required=True)
parser.add_argument("--ca", required=True)
parser.add_argument("--tag", action='append', required=True)
parser.add_argument("--pw", action='append', required=True)
parser.add_argument("--server-port", type=int, required=True, dest='sport')
parser.add_argument("--display-port", type=int, required=True, dest='dport')
conf = parser.parse_args(sys.argv[1:])

server = Server(conf)

actual_tags = server.list_tags()
check_tags(actual_tags, conf.tag)

# Each input password is actually a tag:password
expected_passwords = [p.split(':') for p in conf.pw]

for tag,expected in expected_passwords:
    actual = server.get_password(tag)
    if actual != expected:
        print ("Password for tag '%s' was '%s', but should have been '%s'" % (tag, actual, expected))
        sys.exit(1)

print ("Test passed")
sys.exit(0)
