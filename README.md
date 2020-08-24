dp-elasticsearch
================

Elastic search client with healthcheck. It wraps a dp-net/http client to make requests against elasticsearch service.

### health package

Using elasticsearch checker function currently performs a GET request against elasticsearch 'cluster health' API (`/_cluster/health"`)

The healthcheck will only succeed if the request can be performend and the cluster is in `green` state.
If the cluster is in `yellow` state, a Checker in WARNING status will be returned. In any other case, a CRITICAL Checker will be returned.

Read the [Health Check Specification](https://github.com/ONSdigital/dp/blob/master/standards/HEALTH_CHECK_SPECIFICATION.md) for details.

More information about elasticsearch [cluster health API](https://www.elastic.co/guide/en/elasticsearch/reference/current/cluster-health.html)

Instantiate an elasticsearch client
```
import "github.com/ONSdigital/dp-elasticsearch/elasticsearch"

...
    cli := elasticsearch.NewClient(<url>, <signRequests>, <maxRetries>)
...
```

Call elasticsearch health checker with `cli.Checker(context.Background())` and this will return a check object like so:

```json
{
    "name": "string",
    "status": "string",
    "message": "string",
    "status_code": "int",
    "last_checked": "ISO8601 - UTC date time",
    "last_success": "ISO8601 - UTC date time",
    "last_failure": "ISO8601 - UTC date time"
}
```

### Contributing

See [CONTRIBUTING](CONTRIBUTING.md) for details.

### License

Copyright © 2020, Office for National Statistics (https://www.ons.gov.uk)

Released under MIT license, see [LICENSE](LICENSE.md) for details.
