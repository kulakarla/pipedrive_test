# PipeDrive API Proxy

This is an API proxy application which forwards specified allowed requests to the PipeDrive developer API. This program has been created as a test assignment for the PipeDrive DevOps Engineer Intern position.

The allowed end-points are:
- `GET`, `POST` /deals
- `GET`, `POST` /deals/
- `PUT` /deals/`<id>`
- `PUT` /deals/`<id>`/

All other requests return `405 Method Not Allowed`

## Running & Installation

To run the program, clone the repository:

```sh
git clone https://github.com/kulakarla/pipedrive_test.git
```
To get the application working, your own PipeDrive API token needs to be configured. Inside `config/config.go`, replace the constant APIToken with your PipeDrive API token. The token can be found while logged in to your PipeDrive account, navigate to the top-right. Click on the little avatar, open up `Personal preferences`, open the `API` tab from where you can find your personal API token.

For easier running, the application can be ran using a [Docker](https://www.docker.com/) container.

Inside the folder cloned repository folder, run:

Create the docker image:  
`docker-compose build`

Run the container:  
`docker-compose up`

Now you should have a Docker container up and running.
The API runs on `localhost:8080`

*Alternatively*, if you have Go installed and do not want to use Docker, you can run:  
`go run main.go` for a quick compile & run, or `go build` to build and then run the program.


### Explanations for Work Done

### `main.go`

The `main` method declares handler wrappers for the HTTP end-points and set-ups the server at `localhost:8080`.

### `config/config.go`

This package is a helper/utility package to define the actual PipeDrive API constants. This includes the configurable API token, which must be configured by user. Otherwise, all the requests will fail.

### `handlers/handlers.go`

This is the main controller for handling incoming requests. 

func `Handler` checks if for the request path an allowed request method is used. If `true`, the request is forwarded to the proxy. In case of `false`, a `405 Method Not Allowed` will be returned to the user. This handler takes in requests for path's `/deals`, `/deals/`, `/deals/<id>`, `/deals/<id>/`.

func `MetricsHandler` is the handler for `GET` `/metrics` endpoint. It will get the metrics for the current API session and return them to the user.

func `InvalidPathHandler` is the handler for all various non-allowed requests. Examples like `GET /dsa`, `DELETE /heyjude/23` will get handled by this method. The handler will return a `404 Status Not Found` to the user.

func `RequestMetricsMiddleWare` is a wrapper around an end-point request handler used for gathering metric data. It is used to track and calculate request duration metrics.

### `handlers/handlers_test.go`

This unit test suite includes test for testing individual handlers.

**TestGetHandler** tests that `GET /deals`  request returns HTTP Status 200 OK  
**TestPostHandler** tests that `POST /deals`  request returns HTTP Status 201 CREATED and that the response body includes the correct title  
**TestPutHandler** tests that `PUT /deals/<id>` request changes the currency of a changed deal
**TestMetricsHandler** tests that `GET /metrics` returns the correct count of requests
**TestGetDealsByIDNotAllowed** tests that `GET /deals/<id>` is not allowed  
**TestDeleteDealsNotAllowed** tests that `DELETE /deals` is not allowed  
**TestPatchMetricsNotAlloed** tests that `PATCH /metrics` is not allowed  
**TestInvalidPathHandler** tests that `GET /whatever/1251` is an invalid request  

### `proxy/proxy.go`

func `Request` does the main work of the program - forwards requests to the actual PipeDrive API. From the user sent request to the proxy API, headers and request body are copied and forwarded to the actual PipeDrive API. Response headers and body are then copied and returned back to the user. Additionally, the method calculates the latency of a request.

### `proxy/proxy_test.go`

This unit test sutie includes a test suite for testing the request forwarding.

**TestProxyRequestResponseEqualToDirect** tests that `GET /deals` request response body is the same when calling the PipeDrive API directly and through the proxy, while also checking that the returned header keys and count are the same

### `metrics/metrics.go`

This package is designed to track and manage metrics for a running API session. It tracks total number of requests, total and average latency and request duration metrics for the allowed API endpoints.

struct `MethodMetrics` defines the structure, including JSON, for an allowed API request type.

struct `Metrics` defines the structure for the entire metrics request itself, showing information about the allowed `GET`, `POST` and `PUT` requests

variables `Metrics` is used for creating the metrics instance to track and keep the information, `Mutex` is used for thread-safety.

func `GetMetrics` returns the current metrics of all the endpoints

func `UpdateMetrics` is used for updating the latency metric for given endpoint

func `UpdateDuration` is used for updating the duration and total requests metrics for given endpoint

func `ResetMetrics` is an utility tool for reseting the metrics, used for testing purposes

### `utils/test_utils.go`

This package is created as an utility tool for tests, specifically for testing `POST` and `PUT` methods. As testing these methods create a resource in the actual PipeDrive account for the user and the proxy itself does not allow a `DELETE` endpoint, this package was created.

func `DeleteCreatedResourceInTests` sends a direct `DELETE /deals/<id>` request to the PipeDrive API with the ID given as the function parameter. It is used for deleting the resource created when testing `POST` and `PUT` handlers to avoid unnecessary redundant "trash" deals from accumulating as tests are run frequently. 

Note that this package is not included when building the program. It is only used when compiling for tests and using the tag `testing`.

### `Dockerfile` and `docker-compose.yml`

To ensure problem-free running of the program, Docker capabilities are added.

`docker-compose.yml` creates the image specified in the Dockerfile together with port-mapping to ensure no extra work is needed from the user to get the program working locally

`Dockerfile` specifies the base Go image, the workdir, copies and downoads the module dependencies, copies the code to create the executable. Exposes the `8080` port for the container and will run the main program executable.

### `.github/workflows`

This GitHub Actions specific folder contains 2 workflows:

`print-deploy.yml` is a simple workflow that is ran when a pull-request is merged onto `main` (when something is pushed to `main`). It just echo-s `Deployed!`, nothing else.

`test-and-lint.yml` is another simple workflow that runs Go linting and tests when something is pushed to branch that is part of a pull-request.