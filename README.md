dp-elasticsearch
================

Elasticsearch library to create an elasticsearch client to be able to make requests to elasticsearch.

### elasticsearch package

Includes implementation of a health checker, that reuses the elasticsearch client to check requests can be made against elasticsearch cluster and known indexes.

#### health checker

Using elasticsearch checker function currently performs a GET request against elasticsearch 'cluster health' API (`/_cluster/health"`)

The healthcheck will only succeed if the request can be performend and the cluster is in `green` state.
If the cluster is in `yellow` state, a Checker in WARNING status will be returned. In any other case, a CRITICAL Checker will be returned.

Read the [Health Check Specification](https://github.com/ONSdigital/dp/blob/master/standards/HEALTH_CHECK_SPECIFICATION.md) for details.

More information about elasticsearch [cluster health API](https://www.elastic.co/guide/en/elasticsearch/reference/current/cluster-health.html)

Instantiate an elasticsearch client
```
import "github.com/ONSdigital/dp-elasticsearch/elasticsearch/v2"

...
    cli := elasticsearch.NewClient(<url>, <signRequests>, <maxRetries>, <optional array of indexes>)
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


### awsauth package

Using AWS SDK library to create a signer function to successfully sign elasticsearch requests hosted in AWS.
The function adds multiple providers to the credentials chain that is used by the [AWS SDK V4 signer method `Sign`](https://docs.aws.amazon.com/sdk-for-go/api/aws/signer/v4/#Signer.Sign).

1) **Environment Provider** will attempt to retrieve credentials from `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY` set on the environment.

    Requires `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY` variables to be set/exported onto the environment.

    ```go
    import esauth "github.com/ONSdigital/dp-elasticsearch/v2/awsauth"
    ...
        signer, err := esauth.NewAwsSigner("", "", "eu-west-1", "es")
        if err != nil {
            ... // Handle error
        }
    ...
    ```

2) **Shared Credentials Provider** will attempt to retrieve credentials from the absolute path to credentials file, the default value if the filename is set to empty string `""` will be `~/.aws/credentials` and the default profile will be `default` if set to empty string.

    Requires credentials file to exist in the location specified in NewAwsSigner func.
    File must contain the keys necessary under the matching Profile heading, see example below:
    
    ```
        [development]
        aws_access_key_id=<access key id>
        aws_secret_access_key=<secret access key>
        region=<region>
    ```

    ```go
    import esauth "github.com/ONSdigital/dp-elasticsearch/v2/awsauth"
    ...
        signer, err := esauth.NewAwsSigner("~/.aws/credentials", "development", "eu-west-1", "es")
        if err != nil {
            ...
        }
    ...
    ```

3) **EC2 Role Provider** will attempt to retrieve credentials using an EC2 metadata client (this is created using an AWS SDK session).

    Requires Code is run on EC2 instance.

    ```go
    import esauth "github.com/ONSdigital/dp-elasticsearch/v2/awsauth"
    ...
        signer, err := esauth.NewAwsSigner("", "", "eu-west-1", "es")
        if err != nil {
            ...
        }
    ...
    ```

For more information on Providers for obtaining credentials, [see AWS documentation](https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/configuring-sdk.html#specifying-credentials).

The signer object should be created once on startup of an application and reused for each request, otherwise you will experience performance issues due to creating a session for every request.

To sign elasticsearch requests, one can use the signer like so:

```go
    ...

    var req *http.Request
    
    // TODO set request

    var bodyReader io.ReadSeeker

    if payload != <zero value of type> { // Check for a payload
        bodyReader = bytes.NewReader(<payload in []byte>)
        req, err = http.NewRequest(<method>, <path>, bodyReader)
    } else { // No payload (request body is empty)
        req, err = http.NewRequest(<method>, <path>, nil)
    }
    
    if err = signer.Sign(req, bodyReader, time.Now()); err != nil {
        ... // handle error
    }
    
    ...
```

### Contributing

See [CONTRIBUTING](CONTRIBUTING.md) for details.

### License

Copyright © 2019-2021, Office for National Statistics (https://www.ons.gov.uk)

Released under MIT license, see [LICENSE](LICENSE.md) for details.
