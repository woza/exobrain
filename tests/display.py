#!/usr/bin/python

import sys
from Tkinter import *
import struct
import socket
import ssl
import comms

class would_block (Exception):
    pass

class Display( Frame ):
    NET_STATE_ACCEPTING = 0
    NET_STATE_PARSE_SIZE = 1
    NET_STATE_PARSE_MSG = 2
    NET_STATE_SEND_RESPONSE = 3
    NET_STATE_PARSE_CMD = 4

    CMD_SHOW = 0
    CMD_EXIT = 1

    ENCODE_ASCII = 0
    ENCODE_UTF8 = 1
    def __init__(self, root):
        Frame.__init__(self, root)
        self.root = root
        self.msg = StringVar()
        self.status = StringVar()
        self.msg_out = Label(root, textvariable=self.msg)
        self.msg_out.pack()
        self.status_out = Label(root, textvariable=self.status)
        self.status_out.pack()
        self.net_state = Display.NET_STATE_ACCEPTING
        self.raw_sock = socket.socket( socket.AF_INET,
                                       socket.SOCK_STREAM)

        cert_file='display.crt'
        key_file='display.key'
        ca_file='root.crt'
        self.raw_sock.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)

        self.server = ssl.wrap_socket(self.raw_sock,
                                      key_file, cert_file, True,
                                      ssl.CERT_REQUIRED,
                                      ssl.PROTOCOL_TLSv1_2,
                                      ca_file)
        self.server.bind(('127.0.0.1', 6655))
        self.server.listen(3)
        self.server.setblocking(False)
        self.root.after(100, self.network_callback)


    def network_callback(self):
        try:
            if self.net_state == Display.NET_STATE_ACCEPTING:
                try:
                    self.client,_ = self.server.accept()
                    self.client.setblocking(False)
                    self.need_cmd()
                except socket.error as e:
                    if e.errno != 11:
                        raise e                    
            if self.net_state == Display.NET_STATE_PARSE_CMD:
                self.read_loop(4)
                cmd = struct.unpack('>I', self.buff)[0]
                print "CMD == " + str(cmd)
                if cmd == Display.CMD_EXIT:
                    print "CMD EXIT received"
                    sys.exit(0)
                if cmd == Display.CMD_SHOW:
                    print "CMD SHOW received"
                    self.need_size()            
            if self.net_state == Display.NET_STATE_PARSE_SIZE:
                self.read_loop(4)
                self.msg_size = struct.unpack('>I', self.buff)[0]
                print "Parsed size %d" % self.msg_size
                self.need_payload()
            if self.net_state == Display.NET_STATE_PARSE_MSG:
                self.read_loop(self.msg_size)
                encoding = struct.unpack('>I', self.buff[:4])[0]
                raw_msg = self.buff[4:]
                codec = None
                if encoding == Display.ENCODE_ASCII:
                    codec = 'ascii'
                if encoding == Display.ENCODE_UTF8:
                    codec = 'utf-8'
                if codec is None:
                    print "Unable to detect codec for encoding %d" % encoding
                    sys.exit(1)
                pw = unicode(raw_msg, codec)
                print "Received password '%s'" % pw
                self.write_success()
                print "Wrote success"
                self.need_client()
        except would_block:
            pass

        self.root.after(100, self.network_callback)
        
    def write_success(self):
        msg = struct.pack('>I', 0)
        todo = len(msg)
        off = 0
        while todo > 0:            
            sz = self.client.send(msg[off:])
            off += sz
            todo -= sz
        
    def read_loop(self, req_len):
        try:
            while len(self.buff) < req_len :
                got = self.client.recv(req_len-len(self.buff))
                if len(got) == 0:
                    raise would_block()
                self.buff += got
        except ssl.SSLWantReadError:
            raise would_block()
        
    def need_payload(self):
        self.net_state = Display.NET_STATE_PARSE_MSG
        self.buff = b''

    def need_size(self):
        self.net_state = Display.NET_STATE_PARSE_SIZE
        self.buff = b''

    def need_cmd(self):
        self.net_state = Display.NET_STATE_PARSE_CMD
        self.buff = b''

    def need_client(self):
        self.net_state = Display.NET_STATE_ACCEPTING
        self.buff = b''

display = Display(Tk())
display.root.mainloop()
