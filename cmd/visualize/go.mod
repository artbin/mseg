module visualize

go 1.25

replace (
	github.com/artbin/mseg/mlist => ../../mlist
	github.com/artbin/mseg/segment => ../../segment
)

require (
	github.com/artbin/mseg/mlist v0.0.0-00010101000000-000000000000
	github.com/artbin/mseg/segment v0.0.0-00010101000000-000000000000
	github.com/bradleyjkemp/memviz v0.2.3
)
