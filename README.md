Exobrain is a password-management system.  There is no shortage of
these already, I chose to write exobrain regardless because I wanted
to.


Building
========

GO password database 
-------

   The following instructions assume:
   - You've checked out the code into the directory $EXOBRAIN.
   - You're using a Bash-like shell.  You will need to change step (2) if using a different shell (Csh, Powershell, etc)
   
1. cd $EXOBRAIN/src/modules/go_pwdb
2. export GOPATH=`pwd` 
3. go get -v golang.org/x/crypto/pbkdf2
5. go build
6. cd utils
7. go build put.go
8. go build list.go
9. go build initdb.go

C# UI / Display
---------------

Configuring
===========

TLS Credentials
---------------

   You will need TLS credentials (certificates + private keys) to
   fulfill the following roles:

   1. UI talking to server
   2. Server listening to UI
   3. Server talking to display
   4. Display listening to server

   You will also need CA certificates to validate each of the above
   certificates.

Note that you may be able to use a single certificate for more than
one role - for example, the same CA may have signed all certificates
involved, so you only need one CA certificate.  If you don't want to
use certificate extensions to indicate roles for the certificate (or
if you indicate both server auth and client auth roles) then you can
use the same certificates for components (2,3), and (given the current
implementation with UI/Display as one application) you can also use
one set of credentials for (1,4).

The script tests/gen_test_credentials.sh might be useful as a starting
point when generating TLS credentials.

TLS and Windows

Note that for the UI and Display, you will need the key and
certificate in a single PKCS#12-format file (usual extension .pfx).
This is due to a combination of the C# API and my own laziness.

Also note that the C# APIs I used do not allow me to explictly set
trusted root certificates - you will need to install your CA
certificate into the Windows trust store manually (right-click on your
certificate and select install).  Make sure you install into the XXXX
store, but be aware of the risks: once you do this, Windows will trust
any certificate signed by someone who has the CA's private key.

Also note that all components use TLS version 1.2 exclusively -
ancient versions of Windows have this disabled by default.  I don't
have a sufficiently old version of Windows to test with (I run patched
Windows 7) so I'm not going to give instructions on how to enable this
- I suggest consulting your friendly neighbourhood search engine.

Setting up
==========

Creating TLS certificates
-------------------------

Depending on your security requirements and which implementations you
select, you will need up to 4 TLS endpoint certificates and up to 4 CA
certificates.  The deployment of these certificates is shown in the
following diagram, in which numbers show endpoint certs and letters
show CA certs.

+----+ (A)         (B) +--------+ (C)        (D) +---------+
| UI | (1) <-----> (2) | Server | (3) <----> (4) | Display |
+----+                 +--------+                +---------+

In the figure:

   - (1) is a certificate with the "client auth role" extension,
     signed by CA (B) and used to identify the UI component.
   - (2) is a certificate with the "server auth role" extension,
     signed by CA (A) and used to identify the server component.
   - (3) is a certificate with the "client auth role" extension,
     signed by CA (D) and used to identify the server component.
   - (4) is a certificate with the "server auth role" extension,
     signed by CA (C) and used to identify the display component.

In some circumstances it is possible to reduce the number of
certificates required:

    - If all certificates are issued by the same CA, then only one CA
      cert is required rather than 4.

    - If certificate role extensions are not used, or are used to
      indicate that both client and server auth is supported, then
      only one certificate is needed for the server component.
    
    - If using the C# GUI which acts as UI and Display, and
      certificate role extensions (if present) allow both client and
      server authentication roles, then only one certificate is
      required for both the UI and Display components.

If all these circumstances apply, then the minimal setup is as follows:

+----+ (A)         (A) +--------+ (A)        (A) +---------+
| UI | (1) <-----> (2) | Server | (2) <----> (1) | Display |
+----+                 +--------+                +---------+

    - (A) is a CA that signs both (1) and (2)
    - (1) identifies the C# GUI, (2) identifies the server

You will need to decide how many of the special circumstances apply to
your situation and figure out which certificates you need.  Generating
the certificates is straightforward (on Linux at least) using OpenSSL
- if you need inspiration the script tests/gen_test_credentials.sh is
a good place to start.

Creating a fresh database
-------------------------

You can create a fresh (empty) database by using the utility
$EXOBRAIN/src/modules/go_pwdb/utils/initdb.  Run it and follow the
prompts - you'll need to know which TLS credentials you want to use
and the various network addresses you're going to employ.  When it's
finished, the utility will have written a database file and a
configuration file to the paths you specified.

Here's an example run:
<<< START EXAMPLE >>>
$ ./initdb 
Enter location for new configuration file: example.conf
Enter path for database: example.db
Enter key for talking to UI component: ui_key.pem
Enter certificate for talking to assert identity to UI component: ui_cert.pem
Enter certificate to validate identity of UI component: CA.pem  
Enter key for talking to display component: display_key.pem
Enter certificate for talking to assert identity to display component: display_cert.pem
Enter certificate to validate identity of display component: CA.pem
Enter address:port of display server: 127.0.0.1:6112
Enter adddress:port on which to accept connections from UI component: 0.0.0.0:5443
Enter server name of the display server: display
Enter database password: 
Confirm database password: 
Configuration file written to example.conf
<<< END EXAMPLE >>>


Adding passwords to the database
--------------------------------

Use the put utility - provide it with a configuration file and a username.
Eg: ./put example.conf fred

You will then be prompted for the user's password, a confirmation
thereof and the database password.  If the database password is correct then the user will be added to the database.

You can then inspect the usernames in the database using the list
utility (which requires a configuration file):

./list example.conf
Enter the database password when prompted



Architecture
============

Exobrain is designed as a group of three components: a UI, a server
and a display.  The components communicate over TLS-protected,
mutually-authenticated sockets.  This design allows for flexibility -
the components can reside on different physical systems.  In one
scenario, a central server can store the password database, a user can
request a password from their desktop computer and have the password
sent to their mobile phone.

A network interface is also language neutral.  This was a key
consideration because exobrain is also a project where I can tinker
with new languages - as a first pass I implemented the server in Go
and a combined UI/display in C#, both of which I am new to.  If I want
to try out a new language I can simply write a new server
implementation in it and have said implementation work with existing
UI / display components, or write a new UI component in the new
language and it already has a server to talk to.


