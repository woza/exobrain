#!/usr/bin/python

import sys
from Tkinter import *
import struct
import socket
import ssl

class Display( Frame ):
    NET_STATE_ACCEPTING = 0
    NET_STATE_TALKING = 1
    
    def __init__(self, root):
        Frame.__init__(self, root, addr)
        self.msg = StringVar()
        self.status = StringVar()
        self.msg_out = Label(root, textvariable=self.msg)
        self.msg_out.pack()
        self.status_out = Label(root, textvariable=self.status)
        self.status_out.pack()
        self.net_state = NET_STATE_ACCEPTING
        self.raw_sock = socket.socket( socket.SF_INET, socket.SOCK_STREAM )
        self.raw_sock.bind(addr)
        self.raw_sock.listen(3)
        self.raw_sock.setblocking(False)

    def network_callback(self):
        if self.net_state == NET_STATE_ACCEPTING:
            info = self.raw_sock.accept()
            print "Attempted accept, info " + str(info)
            if info is not None:
                self.raw_client,_ = info
                flags = ssl.OP_NO_SSLv2 | ssl.OP_NO_SSLv3 | ssl.OP_NO_TLSv1 | ssl.OP_NO_TLSv1_1
                self.client = ssl.wrap_socket(self.raw_client, 
                                              self.ssl_keyfile,
                                              True,
                                              self.ssl_certfile,
                                              ssl.CERT_REQUIRED,
                                              ssl.PROTOCOL_SSLv23
                                              self.ca_list,



