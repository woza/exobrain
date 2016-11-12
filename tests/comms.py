import struct

def read_all(sock, todo):
    ret = b''
    while len(ret) < todo:
        got = sock.recv(todo - len(ret))
        ret += got

    return ret

def get_u32(sock):
    raw = read_all(sock, 4)
    return struct.unpack('>I', raw)[0]
