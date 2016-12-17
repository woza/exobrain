**Communications protocol**


The exobrain software is based around a client-server model, where the
server stores a mapping of tags to passwords and the client displays
elements of this mapping to the user.  This document describes how
these two components communicate.

Security considerations
=======================

The link between the client and the server carries sensitive
information (i.e: passwords).  We also want to protect it from
tampering, even when it's carrying non-sensitive information (i.e:
tags).  Both of these aims can be achieved using TLS to protect the
link.  The TLS handshake will be set up as follows:
   - It will only support TLS 1.2
   - It will require mutual authentication


Message formats
===============

Design considerations
---------------------
When selecting message formats, attention has been paid to the following factors:

   - Ease of implementation.  Size and count information is provided
     up front to facilitate allocating buffers to hold incoming
     messages.

   - Uniformity of structure.  Messages which convey the same
     informational elements have the same structure, so they can be
     parsed by a common chunk of code.

   - Platform / language independence.  All strings are accompanied by
     encoding information because client and server might be written
     in different languages with different encoding conventions (e.g:
     UTF-8 vs ASCII).  All integer quantities are transmitted in
     network byte order.

No attention has been paid to the following aspects:

    - Multi-part messages.  All messages consist of a few tens of
      bytes of context and then a tag and a password.  Tags and
      passwords are unlikely to exceed the hundred character mark,
      which means messages are probably limited <= 200 bytes.
      Multi-part-messages or streaming is needless added complexity
      here.

    - Checksums or other integrity checks.  For data in flight, this
      will be handled by standard TCP / TLS mechanisms.

Common fixed values
-------------------
  - STATUS
      o STATUS_OK = 0
      o STATUS_FAIL = 1
  - CMD
      o CMD_QUERY = 0
      o CMD_FETCH = 1
      o CMD_CUSTOM = 500 + x, where x >= 0 is endpoint-specific custom command.
  - ENCODING
      o ENCODE_ASCII = 0
      o ENCODE_UTF8  = 1


1. Query request / response

The "Query" message is used by the client to retrieve a list of tags
known to the server.  The query request consists of a size and a
command:
   | Size    | Command   |
   | 4 bytes | 4 bytes   |
   | 4       | CMD_QUERY |

The query response is more complicated, and comes in two flavours:
success and failure.  Success messages have the following format:

   |Msg Size | Status    | Encoding | Count   | Size    | Tag      | Size    | Tag      | ...
   | 4 bytes | 4 bytes   |  4 bytes | 4 bytes | 4 bytes | variable | 4 bytes | variable |
   |         | STATUS_OK | ENCODING |

The "count" element contains the number of tags in the response.  Each tag is transmitted as two elements:
   1. The size of the tag, when stored in the specified Encoding
   2. The tag itself

The field "Msg Size" contains the size of the entire message, including itself.

Fail messages have the following format:
   |Msg Size | Status      |
   |4 bytes  | 4 bytes     |
   | 8       | STATUS_FAIL |

2. Trigger request / response

The "Trigger" message is used by the client to cause the password for a
given tag to be sent to the display.  The request message has the following format:
   | Size    | Command   | Encoding | Tag      |
   | 4 bytes | 4 bytes   | 4 bytes  | variable |
             | CMD_TRIGGER | ENCODING |

There are again two possible responses, failure or success.  The fail message has the following format:
   | Status      |
   | 4 bytes     |
   | STATUS_FAIL |

The success message is:
   | Status      |
   | 4 bytes     |
   | STATUS_OK |

3. Display request / response

The "Display" message is sent from the server to the display to
request the display of a password.  It has the following format:

   | Size    | Command     | Encoding | Password |
   | 4 bytes | 4 bytes     | 4 bytes  | variable |
   | N       | CMD_DISPLAY | ENCODING |

The value of N (i.e: the size sent in the header) is 12 bytes + the
length of the password in the specified encoding.

The display will send back one of two messages, denoting succes or failure:

   | Status      |
   | 4 bytes     |
   | STATUS_OK |

Or
   | Status      |
   | 4 bytes     |
   | STATUS_FAIL |
