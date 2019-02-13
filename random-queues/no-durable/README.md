Creating a queue using queue_declare is idempotent ‒ we can run the command as many times as we like, and only one will be created.

In this example 50 concurrent "clients" connect and 1000 random queues are created
