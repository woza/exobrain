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
    STATUS_FAIL = 1

    ENCODE_ASCII = 0
    ENCODE_UTF8 = 1
    
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
        print "Success status"
        encoding = comms.get_u32(self.link)
        print "Encoding read %d" % encoding
        nbyte = comms.get_u32(self.link)
        print "Nbyte read %d" % nbyte
        msg = comms.read_all(self.link, nbyte)
        print "Msg read"
        self.link.close()
        return (encoding, msg)

    def trigger(self, key):
        sz = 4 + len(key)
        msg = struct.pack('>II', sz, comms_link.CMD_TRIGGER)
        print "Sending trigger"
        self.link.send( msg )
        print "Sending key"
        self.link.send( key )
        print "Receiving status"
        status = comms.get_u32(self.link)
        print "Received status %d" % status
        if status != comms_link.STATUS_OK:
            print "Failed to receive OK response from server, got %d" % status
            sys.exit(1)
        self.link.close()
        print "Finished trigger"
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
        encoding,info = link.refresh()
        self.options.delete(0, END)
        for element in self.parse_data_stream(info, encoding):
            self.options.insert(END, element)

    def send_password(self):
        link = comms_link( self.server_loc )
        idx = self.options.curselection()[0]
        key = self.options.get(idx)
        link.trigger(key)
        print "Send password"
        
    def parse_data_stream(self, info, encoding):
        '''
        Generator function returning UTF-8 encoding labels
        '''
        
        off = 0
        while off < len(info):
            sz = struct.unpack('>I', info[:4])[0]
            off += 4
            if encoding == comms_link.ENCODE_ASCII:
                print "ASCII encoding"
                value = info[off:off + sz].decode("utf-8")
            else:
                if encoding == comms_link.ENCODE_UTF8:
                    print "UTF-8 encoding"
                    value = unicode(info[off:off + sz], 'utf-8')
                else:
                    print "Unknown encoding"
                    sys.exit(1)
            off += sz
            yield value


root = Tk()
ui = UI(root,('127.0.0.1', 6512))
root.mainloop()
