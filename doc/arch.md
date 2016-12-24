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
