# Application Design

## Base Functionality

The basic functionality should crawl a url for hostnames and and return a count as well as any urls found that can be crawled. The returned urls should have depth associated with them.

## Scaling Base Functionality

The results returned by the base functionality should be stored in a data store. State and results are harder to query from distributed hosts/sources so we should keep this centralized for now.

Counts should be kept in a shared data store.

New urls to crawl should be kept in some kind of queue.

For this design we'll use Redis as our data store and work queue. There are various options for queues and data stores that could provide additional functionality but to keep things simple we'll use redis.

### Why Redis

Redis provides a light weight store that should perform well under signficiant load. It also has features that allow for scaling as this application grows, like replication and clustering.

## Workers

The workers should pull new work items off the queue in configurable batches, each work item should represent a single goroutine. Each work item is associated with a Job.

## Queuing Jobs, Returning Results and Status

To scale horizontally across multiple hosts we need an API to add new crawl jobs, view the state of running jobs and view the results of those jobs.

This API is decoupled from the workers and will be accessible via an http endpoint.
