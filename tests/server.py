
import sys
import SocketServer
import ssl
import struct
import socket
import comms

def parse_table():
    ret = {}
    ret["major"] = "bloodnok"
    ret["harry"] = "seagoon"
    return ret


class listener(SocketServer.BaseRequestHandler):
    CMD_QUERY_ALL = 0
    CMD_TRIGGER = 1
    CMD_EXIT = 2

    CMD_DISPLAY_SHOW = 0
    CMD_DISPLAY_EXIT = 1
    
    STATUS_OK = 0
    STATUS_FAIL = 1

    def handle(self):
        global table
        global display
        size = comms.get_u32(self.request)
        msg = comms.read_all(self.request, size)
        code = struct.unpack('>I', msg[:4])[0]
        if code == listener.CMD_QUERY_ALL:
            self.request.send(struct.pack('>I', listener.STATUS_OK))
            msg = b''
            for k in table.keys():
                msg += struct.pack('>I', len(k))
                msg += k            
            self.request.send(struct.pack('>I', len(msg)))
            self.request.send(msg)
            return

        if code == listener.CMD_TRIGGER:
            key = msg[4:]
            if key not in table:
                self.request.send(struct.pack('>I', listener.STATUS_FAIL))
                return
            
            pw = table[key]
            print "Displaying password '%s'" %pw
            display.send(struct.pack('>II', listener.CMD_DISPLAY_SHOW, len(pw)))
            display.send(pw)
            status = comms.read_all(self.request, 4)
            self.request.send(status)
            return
        
        if code == listener.EXIT:
            display.send(struct.pack('>I', listener.CMD_DISPLAY_EXIT))
            sys.exit(0)
                         
        print "Received unknown command %d" % code
                

cert_file='server.crt'
key_file='server.key'
ca_file='root.crt'

client_sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
display = ssl.wrap_socket(client_sock, key_file, cert_file, False,
                         ssl.CERT_REQUIRED, ssl.PROTOCOL_TLSv1_2,
                         ca_file)
display.connect(('127.0.0.1', 6655))
print "TLS connection established"
table = parse_table()
SocketServer.allow_reuse_address = True
server = SocketServer.TCPServer(('127.0.0.1', 6512), listener)
server.serve_forever()
