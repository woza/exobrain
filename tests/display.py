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
        self.client, _ = self.server.accept()
        self.client.setblocking(False)
        self.need_cmd()
        self.root.after(100, self.network_callback)


    def network_callback(self):
        try:
            if self.net_state == Display.NET_STATE_PARSE_CMD:
                self.read_loop(4)
                cmd = struct.unpack('>I', self.buff)[0]
                if cmd == Display.CMD_EXIT:
                    sys.exit(0)
                if cmd == Display.CMD_SHOW:
                    self.need_size()                    
            if self.net_state == Display.NET_STATE_PARSE_SIZE:
                self.read_loop(4)
                self.msg_size = struct.unpack('>I', self.buff)[0]
                print "Parsed size %d" % self.msg_size
                self.need_payload()
            if self.net_state == Display.NET_STATE_PARSE_MSG:
                self.read_loop(self.msg_size)
                print "Received password '%s'" % self.buff
                self.msg.set(self.buff)
                self.need_cmd()
        except would_block:
            pass

        self.root.after(100, self.network_callback)

    def read_loop(self, req_len):
        try:
            while len(self.buff) < req_len :
                got = self.client.recv(req_len-len(self.buff))
                print "Read " + str(len(got))
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

display = Display(Tk())
display.root.mainloop()
