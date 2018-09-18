# Transporter
A reverse proxy implementation using Go

This uses the [NewSingleHostReverseProxy](https://golang.org/pkg/net/http/httputil/#NewSingleHostReverseProxy) available in httpuitl package to implement the reverse proxy. This helps to filter and redirect to different URLs based on string patterns in URL. 

##### How to build

Clone the repository,

`$ git clone https://github.com/isurusiri/rvpxy.git`  

change directory,

`$ cd rvpxy`

##### How to provide configurations

It is required to have a json file like below:  

`
{
	"port": "",
	"defaultRoute": "",
	"domain": "",
	"routes": [
		{ "detect": "", "uri": "" }
	]
}
`  
`port ` - default port to listen  
`defaultRoute ` - default navigation path  
`domain ` - default application domain  
`detect ` - string pattern to detect when redirecting
`uri ` - target uri
