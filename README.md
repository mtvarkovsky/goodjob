# goodjob
A a small-scale job execution framework for micro-service-oriented architectures.

It provides building blocks to build small-scale orchestrators that execute jobs in reversible, revertible (or both at the same time) manner.

A job is a unit of work that consists of many tasks that have to be executed in a sequential order.

A task is an atomic unit of work within a job.
Tasks contain actual logic of executed job steps and are defined the end-users.

A job can have multiple flavors:
- it can be a simple job that stops when one of the tasks has failed
- it can be repeatable - it can try to continue its execution multiple times when one of the tasks has failed
- it can be reversible - it can try to reverse previous tasks' results when one of the tasks has failed
- it can be repeatable and reversible at the same time

Writing job boilerplate code can be tedious, so a [simple code-generation tool](/tools/generatejob) that creates jobs from .yaml specification is provided.

Examples of how to use this framework can be found [here](examples).
