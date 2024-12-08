# Approach and Overview:

Based on my understanding of the assignment, the primary goal is to create a service in go which primarily

-> deduplicates incoming requests based on id param

-> logs count of unique request per minute to a log file

-> sends a POST request to a give "endpoint" param

For this, I decided to use go's echo framework. And initially for the deduplication aspect of the project decided to implement it using a sync.Map,
The map's key would be the id and in case of any repeated key, the unique counter would not be incremented. After a minute, this map would be cleared.
But since this was a local, in-memory data structure that cannot share state across multiple instances, it was unsuitable for our use case (considering Extension 2).

Redis, on the other hand, provides a centralized space accessible by all instances. The key of the cache would be in the format "minute+id" (yyyy-mm-dd:hh:mm + id).
This way, for a minute, we can capture unique ids even with multiple instances.

Additionally, we're making a POST(Extension 1) endpoint to the "endpoint" parameter passed to our GET api, along with the unique ids' count at that point in time.
[Made an assumption here that the query param is point in time, i.e. if we hit endpoint with "endpoint" param at 30s of the minute, we want to send
the unique id count till that point in the current minute].

The call being made to "endpoint" should be a non-blocking call, since the expectation here is a quick response, it made sense to place the code in a goroutine. So the processing method can return quickly after checking for uniqueness.

Regrettably, I could not implement extension 3 since I had very little experience with kafka and also, the time constraint. I did not want to add code to the proj which i didn't fully understand.



# Structure:

I structured the project so that we have a main file in the root used to register our routes via an echo client object and also start our server. The register routes
calls a method under a route folder, used to define our routes.

This leads to our handler where we define our handler function for each route. This handler extracts the query param, validates them and calls the service.
This is more or less our web layer.

The service holds the implementation logic. This our service layer. [If we had any db operations, the service would call an additional db layer method.]



# Logic of implementation:

I used Redis SETNX with a time-to-live (TTL) of one hour to determine uniqueness.

I used a Ticker to log the number of unique requests every minute into a log file. Log entries include timestamps and counts for better traceability.

Additionally, hit an HTTP POST request to a configurable endpoint with the unique request count for the current minute. Payload is structured as JSON, enabling consumption by any downstream services.

This call to POST is done in a goroutine since its a non-blocking call.

Lastly we count the length of keys at any given time to get the unique id count.



# Areas of improvement:

The validation of the input at handler can be done via json tags like "required". But wanted to be a little verbose here.

The values for the redis client object like ports, db, password should ideally come from a config or from a secret store. Since this was local proj decided against it.

Note:

Performance Testing done via the wrk framework:

wrk -t12 -c400 -d1s 'http://localhost:8080/api/verve/accept?id=400'
