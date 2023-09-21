This is a more realistic example of usage for [goodjob](https://github.com/mtvarkovsky/goodjob).

Let's imagine that we have three services that handle user data:
- Auth - handles authentication and authorization data of users
- Orders - handles some type of orders created by users
- Users - handles user profile data

The dummy code for those services' RPC clients is defiled in the [dummyservice package](dummyservices).

Now, let's imagine that we want to safe-delete all user data.

We can do it in multiple fashions:
- call RPC methods to delete user data in each service and hope for the best, i.e. no errors would occur
- call RPC methods and retry operations that have failed for a reasonable amount of attempts
- call RPC methods and revert successful operations if one of the operations has failed.
- or we can do both, retry failed operations and revert unsuccessful ones if some retry threshold is reached.

Examples of all aforementioned approaches can be found here:
- [simple safe-delete](simple)
- [retry on failure](retryable)
- [revert on failure](revertible)
- [both retry and revert](retryable_revertible)

Each of those packages contains:
- job definitions in the job.yaml files
- generated job boilerplate code in the job.gen.go files
- tests in the job_test.go files

Tasks for those jobs are in the [tasks.go](tasks.go) file.

