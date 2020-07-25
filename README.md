# Objective

Implement a gRPC server using the language of your choice, delivery should be made in the form of a single Pull Request containing all your proposed changes.

The service consists of a single RPC endpoint `Exists`, which should return whether or not the kaiko platform supports a given pair on a given exchange.

You must use the following endpoint `https://reference-data-api.kaiko.io/v1/instruments`, in order to build the response rules.

Response rules are as follows:

* if we find a match for (`exchange_code`, `exchange_pair_code`) the endpoint should return `YES`
* if we don't find a match for (`exchange_code`, `exchange_pair_code`) the endpoint should return `NO`
* if we don't have any match for `exchange_code` the endpoint should return `UNKNOWN`

# Hard Constraints

* response time must be under 100ms
* protobuf definition must not be modified (except for language-related changes such as options)
* client code written in go must not be modified and must be compatible with your implementation
* your implementation must be documented
* your implementation must be ready to be packaged in a docker image (Dockerfile must be implemented)
* your implementation must be integration-tested

# Soft Constraints

* you should implement a build tool (Gradle, Rake, Make, etc.) to ease the build and testing of your implementation
* your commit history should be clean and reflect atomic changes, you are allowed to rewrite history
