# goloadbalancers

A couple of go [http.Handler](https://golang.org/pkg/net/http/#Handler)
middleware implementing various load balancing algorithms.

To play with the toy test server:
 - Get glide (https://github.com/Masterminds/glide) on your local
 - `glide install`
 - ``go test `go list ./... | grep -v vendor` ``
 - `go build`
 - Set your hosts file to include testa, testb, testc, pointing at your localhost
 - `./goloadbalancers`
 - [separate terminal] `node testServer.js`
 
##random
Choose randomly between a set of balancees.

##jsq (JoinShortestQueue)
Choose from balancees the balancee with the currently lowest number of outstanding
requests.

##bestof
After I attended a talk given by Tyler McMullen (it appears a video of another
rendition of it is [here](https://www.youtube.com/watch?v=kpvbOzHUakA))
I was inspired to write this bit of software up. The talk mentioned
[this paper](https://www.eecs.harvard.edu/~michaelm/postscripts/mythesis.pdf)
which apparently advocates a modified Least Connections/Join Shortest queue
algorithm, where, instead of blindly choosing the server with the least
outstanding connections, instead there is a configured number of choices of
random servers, of which are compared the number of outstanding connections.

This helps smooth out the potential for the 'finger of death', whereby in a
standard JSQ implementation, a server newly joining (due to newly happy health
check or growth of the load balancees) the load balancer will suddenly receive
all requests to 'catch it up' to the rest of the servers; this has the potential
to fully tank the newly happy server. This is because it is more likely than 0%
that the new server will not be chosen, and that an existing server will be
taking the load.

The implementation is done such that a golang consumer which can deal with an
[http.Handler](https://golang.org/pkg/net/http/#Handler) will be able to balance
between several specified balancees.
