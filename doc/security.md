Password database
=================

The password database stores a series of key:value pairs in a text
file.  The contents of the text file are encrypted using a key based
on a password entered by the user.  The process of encryption is as
follows:

1. Get password, salt from user
2. Use PBKDF2 with SHA256 and 10^6 iterations to derive the key K
   (The number 10^6 was selected to require about 1 second of computation
    on my test platform.  This computation only needs to be done when loading
    the database - it's not frequent and so a 1-second pause is tolerable.)
3. Generate a new, randomised nonce value N
4. Encrypt the raw password info using AES-256 in GCM mode with K and N,
   producing ciphertext C
5. Write into a file for persistent storage - the format of the file
   is as follows:
       - A one-byte unsigned integer indicating database format version
       
       - A four-byte unsigned big-endian integer indicating the size
         of the password salt

       - The password salt, as a series of bytes

       - The nonce, as a series of bytes - the size of the nonce
         depends on the encryption algorithm used in a particular
         database version and is not coded explicitly

       - The ciphertext

The use of GCM provides for data integrity as well as privacy -
authentication information is generated on encryption and checked on
decryption.  Because we re-use keys, it is important that we use a new
nonce each time the database is encrypted.

Storing the salt in the database means that an old database can be
recovered from backups provided the user knows the password - it is
not necessary to keep a config file (with salt) synched to the
database when handling backups.

Data in transit
================

Data in transit is protected by TLS 1.2, and each connection requires
mutual authentication between its endpoints.

Note that exobrain is intended for deployment in a LAN environment
where a system administrator has access to reconfigure it as required.
It therefore does not implement CRL checking, OCSP checking, etc.  If
a certificate is no longer trusted, the administrator needs to issue a
new one and change the configuration file(s) to stop the old ones
being used.
