# In Memory TFTP Server

![alt tag](https://raw.github.com/gabrielhartmann/tftp/master/tftp_demo.gif)

To run the TFTP server:

```go
$ go run server.go
```

Output will indicate the port to which the client should connect:

```
INFO[0000] UDP local address: [::]:64372
```

Using the 'tftp' client included with OSX, the typical commands (keeping in mind warnings below) would be:

```sh
$ tftp
tftp> connect localhost 64372
tftp> put foo.txt
putting foo.txt to localhost:foo.txt [octet]
Sent 1942 bytes in 0.0 seconds [inf bits/sec]
tftp> get foo.txt
getting from localhost:foo.txt to foo.txt [octet]
Received 1942 bytes in 0.0 seconds [inf bits/sec]
```

To start reading the code, it is helpful to note that there are two major components: the tftp server and the file server.  They are located in the appropriately named packages / directories.

To start reading the tftp server code, a good place to start would be with the three session files: req_session.go, read_session.go, and write_session.go.  The request session (req_session.go) spawns read or write sessions for each request it gets from a client.  The main code driving the UDP connectivity is in reader_writer.go.  The main method in server.go consists entirely of spawning a request session.

The file server code is very straight forward.  There is a file defining a file server interface, and an in memory implementation of that interface.

A word of warning, this is only an in memory TFTP server, so files are not written to disk on the server side.  A different implementation of the file server interface could provide persistent storage.

Please note that the bundled OSX client has slightly odd behavior.  When it requests a file which does not exist, and it correctly receives an error packet indicating this, it still overwrites the local file with an empty file.

Choose any client, but this server is restricted to an early unextended spec of a TFTP server.  In the 'tftp' client that is packaged with OSX be sure to consult the '?' help menu.  Please set the mode to binary (octet) and turn off timeout, tsize, and non-standard (not 512B) block sizes.

This can be verified by toggling 'trace' and 'verbose' on and verifying that the packets being sent match the spec in Appendix I here: http://tools.ietf.org/html/rfc1350

```sh
tftp> binary
mode set to octet
tftp> tout
Timeout option off.
tftp> tsize
Tsize mode off.
tftp> trace
Packet tracing on.
tftp> verbose
Verbose mode on.
tftp> blksize
(blksize) 512
```
