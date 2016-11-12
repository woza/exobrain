#!/usr/bin/python

import sys
from Tkinter import *
import socket
import struct
import comms

class comms_link:
    CMD_QUERY_ALL = 0
    CMD_TRIGGER = 1
    CMD_EXIT = 2
    
    STATUS_OK = 0
    def __init__(self, server_loc):
        self.server_loc = server_loc
        self.link = socket.create_connection( server_loc )

    def exit(self):
        msg = struct.pack('>II', 4, comms_link.CMD_EXIT)
        self.link.send( msg )
        
    def refresh(self):
        msg = struct.pack('>II', 4, comms_link.CMD_QUERY_ALL)
        self.link.send( msg )
        status = comms.get_u32(self.link)
        if status != comms_link.STATUS_OK:
            print "Failed to receive OK response from server, got %d" % status
            sys.exit(1)
        sz = comms.get_u32(self.link)
        msg = comms.read_all(self.link, sz)
        self.link.close()
        return msg

    def trigger(self, key):
        sz = 4 + len(key)
        msg = struct.pack('>II', sz, comms_link.CMD_TRIGGER)
        self.link.send( msg )
        self.link.send( key )        
        status = comms.get_u32(self.link)
        if status != comms_link.STATUS_OK:
            print "Failed to receive OK response from server, got %d" % status
            sys.exit(1)
        self.link.close()
        return msg
        

class UI( Frame ):
    def __init__(self, root, server_loc):
        Frame.__init__(self, root)
        self.server_loc = server_loc
        self.pack()
        self.refresh_button = Button(self, text="Refresh", command=self.do_refresh)
        self.refresh_button.grid(row=0, column=1)
        self.quit_button = Button(self, text="Quit", command=self.do_quit)
        self.quit_button.grid(row=0, column=0)
        self.options = Listbox(self)
        self.options.grid(row=1, column=0)
        self.refresh_button = Button(self, text="Send Password", command=self.send_password)
        self.refresh_button.grid(row=2, column=0)

        self.do_refresh()

    def do_quit( self ):
        link = comms_link( self.server_loc )
        link.exit()
        self.quit()
        
    def do_refresh( self ):
        link = comms_link( self.server_loc )
        info = link.refresh()
        self.options.delete(0, END)
        for element in self.parse_data_stream(info):
            self.options.insert(END, element)

    def send_password(self):
        link = comms_link( self.server_loc )
        idx = self.options.curselection()[0]
        key = self.options.get(idx)
        link.trigger(key)
        
    def parse_data_stream(self, info):
        off = 0
        while off < len(info):
            sz = struct.unpack('>I', info[:4])[0]
            off += 4
            value = info[off:off + sz]
            off += sz
            yield value


root = Tk()
ui = UI(root,('127.0.0.1', 6512))
root.mainloop()
