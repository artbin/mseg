module visualize

go 1.25

replace (
	github.com/artbin/mlist/mlist => ../../mlist
	github.com/artbin/mlist/segment => ../../segment
)

require (
	github.com/artbin/mlist/mlist v0.0.0-00010101000000-000000000000
	github.com/artbin/mlist/segment v0.0.0-00010101000000-000000000000
	github.com/bradleyjkemp/memviz v0.2.3
)
