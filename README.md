## cli

    This Go Cli makes https requests for the given urls in a parallel manner. The default number of parallel requests is 10.

## Prerequisite

    You need to have go installed. To install go, please see https://go.dev/doc/install.

## Compile

    To compile the code , run `make build` in the root directory.This will compile the code and create a binary ./myhttp.

## Run

    Once you compiled the code as I mentioned above, makefile is configured to create a binary `myhttp` inside the root directory , CLI is now ready to be executed.
    Running `./myhttp --help` will provide you the usage and help you understand how to use this cli.

    for example, you can run the cli as follows:
    `./myhttp  --parallel 3  --urls="google.com facebook.com yahoo.com yandex.com twitter.com"`

    The cli will request all the urls provided in --urls with 3 parallel requests at a time.

## Test

    To run the tests, run `make test` in the root directory.This w

## Test Coverage

    To run the tests coverage, run `make test_coverage` in the root directory. This will display how much of the code is tested using the unit tests.

## Implementation Details

    -   I have used workerpool pattern with buffered channel to do a job , in this case, making requests to provided urls,  This will ensure that all the requests will be executed parallely with a limitation.
    -   I have also applied SOLID principles using interfaces as [davecheney mentioned] (https://dave.cheney.net/2016/08/20/solid-go-design) to be able to write testcases effectively. This provided loose coupling and higher cohesion as I created the client and workerpools as a separate packages
    -   I tried not to use any third party packags as much as possible to have this as lighter as possible.
    -   I could have used dockerfile for golang, however since this is a simple cli, I thought that it might look overengineering.
    -   I have added code coverage to give an understanding on how much of the test coverage I have provided.

## Improvements

    -   we could make the DefaultMaxParallelWorkers as an env config so that it is easy to change this count without having to replace the binary.
    -   we could have better logging system for this code if we want to have production ready.
