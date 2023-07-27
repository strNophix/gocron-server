# gocron-server
Small job scheduling server/library based on [gocron](https://github.com/go-co-op/gocron). Besides providing cron-functionality it also offers a gRPC API to manually trigger/schedule jobs and listen for their completion.

Sample config:
```toml
[server]
host=":9092"

# This unit will be executed every minute or can be manually triggered.
[[unit]]
name="echo"
command="echo hello"
cron="* * * * *"

# This unit can only be manually triggered.
[[unit]]
name="notify"
command="notify-send hello"
```

## Usage
```sh
go run cmd/main.go <job-definition-toml>
# Example: go run cmd/main.go examples/config.toml
```

Jobs can also be defined using code by extending the server. For an example of extending the server see [this example](./examples/counter.go).
