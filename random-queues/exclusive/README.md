# random queue name: amq.gen-PhvTLHYLYYV0LgSWX86mcQ

In the amqp client, when we supply queue name as an empty string, we create a
non-durable queue with a generated name:

    random queue name: amq.gen-PhvTLHYLYYV0LgSWX86mcQ


When the connection that declared it closes, the queue will be deleted because
it is declared as exclusive.

to test:

    $ go run mail.go
