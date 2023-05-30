# http benchmarking

A simple tool for measuring http request latency

Metrics exposed:

- Number of total requests
- Total request latency grouped by status code

Available at /metrics for prometheus integration

Configuration is loaded from an .env file at the root of the project. (example file available)

- Application config:

| Option    | Type    | Description                                | Default               |
|-----------|---------|--------------------------------------------|-----------------------|
| ip        | string  | server ip                                  | localhost             |
| port      | string  | server port                                | 8080                  |
| threads   | int     | # of threads                               | 1                     |
| requests  | int     | # of requests                              | 3                     |
| endpoint  | string  | request endpoint url                       | [https://www.google.gr](https://www.google.gr) |
| client_timeout | int | http client request timeout in seconds    | 10                    |
| method    | string  | http method                                | GET                   |
| frequency | string  | the frequency to repeat the check          | 10s                   |
| buckets   | string  | comma separated latency metric buckets (ms)| 50,100,500,1000       |
| verbose   | bool    | verbosity                                  | false                 |
| uuidParam | bool    | include uuid parameter in request          | false                 |

- Prometheus config:

| Option    | Type    | Description                           | Default  |
|-----------|---------|---------------------------------------|----------|
| app       | string  | app name                              | -        |
| port_name | string  | port name                             | -        |
| namespace | string  | prometheus namespace                  | -        |

:warning: If the method is POST/PUT/PATCH, the app searches for a file named req.json containing the request body in json format.

After creating the env file, you can:

```bash
# build docker image
make build
# build and run the docker image
make run
# create the yml files needed for prometheus monitoring, under ./k8s/deploy
make prometheus
# apply the k8s files
make apply
