Exobrain is a password-management system.  There is no shortage of
these already, I chose to write exobrain regardless because I wanted
to.


Building
========

GO password database 
-------

   The following instructions assume that the reader has checked out the code into the directory $EXOBRAIN, and is using a Bash-like shell (for step 2).
   
- cd $EXOBRAIN/src/modules/go_pwdb
- export GOPATH=`pwd` 
- go get -v golang.org/x/crypto/pbkdf2
- go build
- cd utils
- go build put.go
- go build list.go
- go build initdb.go

C# UI / Display
---------------
TODO

Configuring
===========

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

- If all certificates are issued by the same CA, then only one CA cert
  is required rather than 4.

- If certificate role extensions are not used, or are used to indicate
  that both client and server auth is supported, then only one
  certificate is needed for the server component.
    
- If using the C# GUI which acts as UI and Display, and certificate
  role extensions (if present) allow both client and server
  authentication roles, then only one certificate is required for both
  the UI and Display components.

If all these circumstances apply, then the minimal setup is as follows:

    +----+ (A)         (A) +--------+ (A)        (A) +---------+
    | UI | (1) <-----> (2) | Server | (2) <----> (1) | Display |
    +----+                 +--------+                +---------+

- (A) is a CA that signs both (1) and (2)
- (1) identifies the C# GUI, (2) identifies the server

You will need to decide how many of the special circumstances apply to
your situation and figure out which certificates you need.  Generating
the certificates is straightforward (on Linux at least) using
OpenSSL. If you need inspiration the script
$EXOBRAIN/src/modules/go_pwd/simple_test.sh is a good place to start.

Creating a fresh database
-------------------------

You can create a fresh (empty) database by using the utility
$EXOBRAIN/src/modules/go_pwdb/utils/initdb.  Run it and follow the
prompts - you'll need to know which TLS credentials you want to use
and the various network addresses you're going to employ.  When it's
finished, the utility will have written a database file and a
configuration file to the paths you specified.

Here's an example run:
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

##Manipulating an existing database

### Adding passwords

Use the put utility - provide it with a configuration file and a username.
Eg: `./put example.conf fred`

You will then be prompted for the user's password, a confirmation
thereof and the database password.  If the database password is
correct then the user will be added to the database.

### Listing usernames for which passwords are known

You can then inspect the usernames in the database using the list
utility (which requires a configuration file):

`./list example.conf`
Enter the database password when prompted

### Listing known passwords
You cannot do this via the CLI, you must talk to the server using its network interface.

### Deleting a username and its associated password
`./delete example.conf <doomed username>`

### Changing the database encryption key
`./rekey example.conf` and follow the prompts

