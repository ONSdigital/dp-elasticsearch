# dp-elasticsearch

Elasticsearch library to create an elasticsearch client to be able to make requests to elasticsearch. Currently the library support elasticsearch 7.10. The version 7.10 uses go-elasticsearch library behind the scenes and this library can be viewed as a wrapper around go-elasticsearch library.  Please follow readme on how to create different versions of client and how to consume these.  

## elasticsearch package

Includes implementation of a health checker, that reuses the elasticsearch client to check requests can be made against elasticsearch cluster and known indexes.

### setup elasticsearch client

#### setup ES 7.10 client

Setting up ES7.10 client is similar as setting up es2.2, just that we need to specify client library as ```GoElasticV710```, as follows:

```golang
import (
    dpEs "github.com/ONSdigital/dp-elasticsearch/v3"
)

...
    esClient, esClientErr := dpEs.NewClient(dpEsClient.Config{
        ClientLib: dpEsClient.GoElasticV710,
        Address:   cfg.esURL,
    })
    if esClientErr != nil {
        log.Fatal(ctx, "Failed to create dp-elasticsearch client", esClientErr)
    }
...
```

##### Embedding custom http roundtripper with es7.10

You could create custom roundtripper (say if you have to sign requests if you are using es7.10), as follows:

```golang
import (
    dpEs "github.com/ONSdigital/dp-elasticsearch/v3"
    "github.com/ONSdigital/dp-net/v3/awsauth"
)

...
if cfg.signRequests {
        fmt.Println("Use a signing roundtripper client")
        awsSignerRT, err := awsauth.NewAWSSignerRoundTripper(awsauth.NewAWSSignerRoundTripper(<aws_filename_placeholder>, <aws_profile_placeholder>, <aws_region_placeholder>, "es",
            awsauth.Options{TlsInsecureSkipVerify: cfg.aws.tlsInsecureSkipVerify})
        if err != nil {
            log.Fatal(ctx, "Failed to create http signer", err)
        }

        esHTTPClient = dphttp.NewClientWithTransport(awsSignerRT)
    }
    
    esClient, esClientErr := dpEs.NewClient(dpEsClient.Config{
        ClientLib: dpEsClient.GoElasticV710,
        Address:   cfg.esURL,
        Transport: esHTTPClient,
    })
    if esClientErr != nil {
        log.Fatal(ctx, "Failed to create dp-elasticsearch client", esClientErr)
    }
...
```

More details on [aws signer roundtripper](https://github.com/ONSdigital/dp-net/tree/main/awsauth#round-tripper)

##### setting up bulk indexer with es7.10

```golang
import (
    dpEs "github.com/ONSdigital/dp-elasticsearch/v3"
)

...
    esClient, esClientErr := dpEs.NewClient(dpEsClient.Config{
        ClientLib: dpEsClient.GoElasticV710,
        Address:   cfg.esURL,
    })
    if esClientErr != nil {
        log.Fatal(ctx, "Failed to create dp-elasticsearch client", esClientErr)
    }
    
    if err := esClient.NewBulkIndexer(ctx); err != nil {
        log.Fatal(ctx, "Failed to create new bulk indexer")
    }
    
    // Adding docs to bulk indexer
    err := esClient.BulkIndexAdd(ctx, v710.Create, indexName, documentID, documentBody)
        if err != nil {
            log.Fatal(ctx, "Failed to add documents to bulk indexer")
        }
        
    // Close bulk indexer
    err := esClient.BulkIndexClose(ctx)
    if err != nil {
        log.Fatal(ctx, "Failed to close bulk indexer")
    }
...
```

#### health checker

Using elasticsearch checker function currently performs a GET request against elasticsearch 'cluster health' API (`/_cluster/health"`)

The healthcheck will only succeed if the request can be performend and the cluster is in `green` state.
If the cluster is in `yellow` state, a Checker in WARNING status will be returned. In any other case, a CRITICAL Checker will be returned.

Read the [Health Check Specification](https://github.com/ONSdigital/dp/blob/master/standards/HEALTH_CHECK_SPECIFICATION.md) for details.

More information about elasticsearch [cluster health API](https://www.elastic.co/guide/en/elasticsearch/reference/current/cluster-health.html)

Instantiate an elasticsearch client

```golang
import "github.com/ONSdigital/dp-elasticsearch/elasticsearch/v3"

...
    cli := elasticsearch.NewClient(<url>, <maxRetries>, <optional array of indexes>)
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

## Contributing

See [CONTRIBUTING](CONTRIBUTING.md) for details.

### License

Copyright © 2019-2025, Office for National Statistics <https://www.ons.gov.uk>

Released under MIT license, see [LICENSE](LICENSE.md) for details.
