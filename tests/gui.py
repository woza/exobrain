#!/usr/bin/python

import sys
from Tkinter import *
import socket
import struct

class comms_link:
    CMD_QUERY_ALL = 0

    STATUS_OK = 0
    def __init(self, server_loc):
        self.server_loc = server_loc
        self.link = socket.create_connection( server_loc )

    def refresh(self):
        msg = struct.pack('>II', 4, CMD_QUERY_ALL)
        self.link.send( msg )
        msg = self.link.recv(4)
        sz = struct.unpack('>I', msg)[0]
        msg = self.link.recv(sz)
        status = struct.unpack('>I', msg[:3])[0]
        if status != STATUS_OK:
            print "Failed to receive OK response from server, got %d" % status
            sys.exit(1)
        ret = msg[3:]
        self.link.close()
        return ret

class UI( Frame ):
    def __init__(self, root, server_loc):
        Frame.__init__(self, root)
        self.server_loc = server_loc
#        self.pack()
        self.refresh_button = Button(self, text="Refresh", command=self.do_refresh)
        self.refresh_button.grid(row=0, column=1)
        self.quit_button = Button(self, text="Quit", command=Frame.quit)
        self.refresh_button.grid(row=0, column=0)

        self.do_refresh()

    def do_refresh( self ):
#        link = comms_link( self.server_loc )
 #       info = link.refresh()
        self.pw_buttons = []
        for element in self.parse_data_stream(info):
            b = Button(self, text=element, command=self.send_password)
            self.pw_buttons += [b]
            b.grid(row=len(self.pw_buttons), column=0)

    def send_password(self):
        for i,b in enumerate(self.pw_buttons):
            print "Checking button %d get %s" % (i, str(b.get()))

    # def parse_data_stream(self, info):
    #     off = 0
    #     while off < len(info):
    #         sz = struct.unpack('>I', info[:3])[0]
    #         off += 4
    #         value = info[off:off + sz]
    #         off += sz
    #         yield value

    def parse_data_stream(self, info):
        yield "Major Bloodnok"
        yield "Henry Crun"
        yield "Count Moriarty"
