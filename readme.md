# todo

```text
[x] implement cloudwatch metrics for application and url
[x] implement parallelism to finish requests faster
[] terraform should have
    - eks/ecs
        should we use container service?
        we are using ECS with Fargate
    - cloudwatch
        - seems to be hooked up for container logs, will have to push custom metrics for URL responses
    - s3
    - amazon container registry
        - created. will need pipelines to build the code and push the updated image in
    - route53
[] pipelines
    - github actions
    - 2 repos (1 infra and 1 application and push image to acr)
```

## design considerations

the idea would be to write all the responses to a single file inside S3 and serve the file whenever the endpoint is
called

- for the endpoint, redirect to the s3 json file to reduce load on the container since we do not serve the file from the container itself
- should we do 1 file inside S3 bucket for 1 URL?
    - this allows multiple pods to work on the same thing
    - use md5hash to check if the file has been modified
      - this will only be helpful if we load entire url csv into memory
    - might not make sense since you will need a page that contains all the responses
- cloudwatch
    - trigger alarms if websites constantly sends a non-200 response
    - use time taken for requests as a metric to push to cloudwatch metrics
        - for improvement, can use a time series database (influxdb etc.)
    - application health
- security
  - ideally, we will want to place the cluster inside a private subnet and place the ALB and use an IGW for internet access
  - probably need some form of IAM to ensure only the application can use the bucket
  - we should probably use secrets to pass in the s3 key then as future improvement, write a function that generates
        the key based on its own instance profile so you don't have to do key rotation
- error handling
  - if input file is invalid, trigger alarm?
- secrets
  - s3 key
  - cloudwatch logs
  - cloudwatch metrics
- permissions for iam 
  - s3
  - cloudwatch
  - acr?

## stuff

```text
curl 169.254.170.2$AWS_CONTAINER_CREDENTIALS_RELATIVE_URI
```

check which directory to write the input.csv to then if check md5hash every 10 minutes, if md5hash doesnt change, we dont need to download and just use back the old file
if hash changed then we call the download function, potentially not md5hash but some other form of checksum.

- use environment variables for bucket name and input file name then declare inside environment variables inside task definition
- enable load balancer

[x] vpc requires internet gateway and route 0.0.0.0/0 to igw
- create healthcheck for container by using healthcheck api 

210123

we still lack the ecs service and task definition for terraform
after testing, can see how to optimize terraform to be more modular
missing subnet association to route table