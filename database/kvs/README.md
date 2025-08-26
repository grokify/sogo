# Document

Sogo `github.com/sogo/database/kvs` provides a storage abstraction layer for various storage engines with a common configuration struct and getter/setter interfaces.

Initially, it's goal is to provide services for simple key value string storage. This use case is very simple but often times you may want to use a different solution depending on your provider, e.g. DynamoDB for AWS, Redis for Heroku, and the local file system for testing.

Supported engines include:

* DynamoDB
* Files
* Redis
