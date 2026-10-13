module github.com/artbin/mseg/pfqueue

go 1.25

replace github.com/artbin/mseg/segment => ../segment

require github.com/artbin/mseg/segment v0.0.2

retract v0.0.1 // published tag does not match the checksum database
