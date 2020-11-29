# grpc-rest-compare

gRPC is a high performant framework to exchange message for remote procedure call. It use protocol buffers for serializing structured data.
RESTful API is a very popular Web Service framework, which allow web developer to exchange json message through http call.
On the perspective of ease of use, REST is definitely easier than gRPC because gRPC requires the developer to define the message struct in the protocol buffer, and compile it to GO module with protoc compiler.

However, on the perspective of performance, gRPC is much faster than REST. Google Cloud run service can now serve both HTTP and gRPC request.Here I have done a simple performance benchmark to compare gRPC and REST using Google Cloud run.

The test program is written using GO lang. The client has 10 threads, which will send a large-size message with a name string and an array with10K integers, and each thread will repeat 5000 times for sending the message. The Cloud Run service is defined as 1 instance 1 CPU with concurrency = 10.

The gRPC throughput is 48 requests per sec which is > 10X than REST API. It uses much less CPU time to process each message ( cpu time ms / req per sec = 832 / 48 = 17.33 ms) than REST (404 ms / 4 = 101 ms). While the request latency is similar, which are 6ms and 8ms respectively.

The Winner is gRPC.

However, when we add Cloud pubsub in the match, pubsub easily win the game becaues it is async. We use a GO client to publish the similar message string and then use another GO client to pull the message. The pubsub topic publication throughput is 139 message per sec, and the pull subscription throughput is about 104 message per sec. So the final winner is pubsub.

https://iwasnothing.medium.com/grpc-vs-rest-performance-comparison-1fe5fb14a01c
