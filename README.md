# Crawler

Service to recursively count hosts in the body of a given url

* [Design](DESIGN.md)

## Requirements

* Go 1.13
* Docker
* docker-compose

## Build

The build process places binaries in the top level `./bin` directory.

```bash
make dep
make build
```

## Run Tests

Unit tests can be run with the following commands.

```bash
make dep
make test
```

To run integration testing we need to spin up a redis instance and a webserver. Normally I'd build out mocks for everything I could to avoid having to run any services outside of my code. I do this to ensure my tests run quickly and to avoid any dependency or flaky test scenarios (like race conditions when running tests concurrently). I haven't built this out as I normally would due to time constraints.

Follow these steps to spin up and run the application.

```bash
make build
```

This will run your redis container and expose it on port 6379

```bash
docker-compose up
```

Run a test server serving the example pages, you'll need to add a host entry for `example.com` pointing to `127.0.0.1`. Normally I wouldn't do it this way but again, time constraints. This will allow `example.com` on port `80` to your local system.

```bash
./bin/testserver
```

Next we run the apiserver application, again normally I'd automate this kind of testing but due to time constraints we're just going to validate the example scenarios. In my design I had intended for this to be a webserver with a RESTlike API that we could `POST` these requests to, this should simulate that behavior and create a new job with a unique ID.

```bash
./bin/apiserver
```

Next we run the worker(s), we can run any number of them and they should process jobs and crawl pages. They currently will only batch 10 items at a time but if I had more time I would use environment variables to make that configurable.

```bash
./bin/worker
```

We should see the new job pick up, new crawl jobs triggered and the results of the crawl jobs output to stdout. These results are being stored in redis and if I had more time I would have created an API to display both the state of the running jobs and the results via API calls from the apiserver process.

Additionaly for testing I would expand the docker-compose to build each component indivudually and allow for testing with multiple workers (using the docker-compose scaling features).

## Deployment

To deploy this in a production ready manner we would make several changes.

* Build multi-node redis cluster across multiple availability zones, using managed services if available.
* Secure access to redis
* Build docker images for each component
* Use k8s or some other container orchestration mechanism to deploy at least 3 replicas of the API component, and at least 3 of the worker components. Each of these components should be auto-scaled horizontally as necessary.
* We would also want to add some kind of authentication to the API

## Improvements

There are lots of things that can be improved.

* I'd probably replace redis entirely. Maybe use a kvstore like dynamodb to keep state for my application and a different message queue service, maybe something like RabbitMQ that can deal with timeouts, retries and priority etc. or a managed service like SQS.
* There are lots of opportunities to refactor in the current code base, there is a fair amount of code duplication and since it's a v0 design I'd probably iterate on some of the patterns to make it more robust and easier to test.
* More unit testing, mocking for integrations and perf testing.
* I would want to also put together a pipeline that would build, test and validate changes and deploy progressively.
* The application is also lacking instrumentation and logging, we would need to add metrics, tracing and some logging framework to ensure we have visibility in case things go wrong.
* We need health endpoints and we need to capture signals so that we can respond to graceful shutdown events without data loss. This is important to ensure the lifecycle of our containers is stable and usable in a production environment.
* Finally, I would probably change my design a little bit as well, in distributed systems tracking state changes can be expensive and difficult. If it made sense for the use case I would push everything to a stateful message bus and consume/produce all my changes to it to ensure my state gets updated asynchronously. This model would allow us to scale globally and it would remove single points of failure in the system. Obviously such a system could be more complex but it would provide a way to track state changes in a repeatable way.

## Scaling

I think I've covered most of the design changes I would make that would impact scaling, but there are certainly ways to improve the current design.

* Currently redis is going to be our bottleneck, it is the single source of truth and while it's very fast it can only vertically scale so far. We could leverage the different partitioning and clustering capabilities of redis to squeeze some more performance out of it.
* Right now the batch size hasn't been optimized in any way, we would have to do some perf testing to validate appropriate values.
* Currently we expect our requests to complete within a reasonable amount of time, we have no mechanism to timeout, retry or backoff. This could lead to blocking requests making our workers sit idle waiting for responses. Ideally we would timeout quickly and retry, maybe with some kind of backoff or priority queing.
