Submission by Geoff Lamperd for the 'Article API' test.

My first thought was to provide this solution using Java, since I'm new to Go. I did a quick search at the outset,
with the requirements for this test in mind, and it appeared that Go and its ecosystem provide some advantages. Not least, 
a lightweight web server and HTTP router. I decided to throw caution to the wind and develop the solution in Go. That 
meant learning the basics of the language, concentrating on just the elements I'd need for this solution: code environment,
some data structures, HTTP classes. I relied heavily on  https://astaxie.gitbooks.io/build-web-application-with-golang/en/
which covered pretty much everything I needed.

The solution utilises https://github.com/julienschmidt/httprouter This is a HTTP request router that works with
net/http and implements routing with a minimum of fuss.  

The data structures used within the service are:
- articleList. A simple array of Article objects, and it provides the cache of article data. Designed to support loading from a text file at startup. This 
	capability was not requested, but I thought it would help with testing. N.B. the provided file has
	a single request, ID=1. A real-world implementation would
	use a map, and would enforce uniqueness of ID, which I have not done as it didn't seem to be explicitly required. Duplicate IDs will not be rejected by the POST /article handler, but requests will only 
	retrieve the first such article. No doubt a real-world version of this service would include
	end points for update and delete functions. I have not catered for these.
- tagDateMap. A nested map, keyed by tag then date, to support the GET /tag/{tag}/{date} request. The 
	structure is updated whenever a new article is added, and is designed for efficiency in handling a
	client's /tag GET request. For a real world
	solution we'd want to understand more about the distribution of tags and dates to arrive at the most efficient structure. Then again, in the real world, there would no doubt be a database backing the 
	request, not an internal temporary structure. 
	
Error handling: Errors encountered while handling a request will result in a response to the client with 
	an appropriate HTTP response code. Some details will also be logged to stdout. 
	
Tests: I tested basic operation of the 3 end points, along with a few edge cases and error cases. I added
	enough articles to test that the stated limit of 10 article IDs was honoured in the /tag response. I used
	the Postman tool for REST requests, which allowed me to quickly repeat a series of requests.

Timing: I took about 8 hours to set up Go, learn some basics, and work through some exercises. Time for implementation of the service was probably 3-4 hours.

Installing and running:

- Download from github to your local environment. The URL is https://github.com/DalaiLlaama/golang-sample 
- Open a terminal/command line session and cd to the 'golang-sample/rest-server' folder.
- Enter the command `go build`. This assumes you have the Go environment set up.
- Enter `rest-server` to start the service. The service will bind to localhost port 8080. 
- Using curl or your favourite REST request client, send requests to 'http://localhost:8080'

Code for the service is in `golang-sample/rest-server/rest-server.go`
The code for `httprouter` is included, in a separate folder, for convenience in installation. It is a clone of the
library mentioned above.  
