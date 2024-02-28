# Blatta

Blatta is a command-line utility designed to help with monitoring, troubleshooting and optimizing [CockroachDB](https://www.cockroachlabs.com/).

Why is it called Blatta? Blatta is the [Genus of cockroaches](https://en.wikipedia.org/wiki/Blatta). It means light-shunning insect.
This sounded like a great name for a tool to use with CockroachDB.

# Commands

## Command-line Flags

Blatta uses a number of command-line flags that vary depending on the command that is executed.

In most cases, a configuration file for each CockroachDB environment is the easiest
approach to managing the longer command-line flags.

To view all of the flags available to a specific command, run `blatta [COMMAND] [SUBCOMMAND] --help`.
For example, to view all of the command-line flags available to the `monitor hotRanges` command,
run `blatta monitor hotRanges --help`.

## Setup a Configuration File

Blatta configuration parameters can be passed via command-line flags, but it is easier to define a configuration file.

By default, Blatta will look in the user home directory for the file `.blatta.yml`. 
Here is an example configuration file:

```
url: http://localhost:8080
username: demo
password: demo10166
cacert: /var/folders/rg/_ltnftys4wj6ltzm0y45mhsh0000gq/T/demo657659925/ca.cert
pgurl: postgres://demo:demo10166@127.0.0.1:26257?sslmode=require
insecure: true
```

To load a different configuration file, use the `--config` parameter with 
the configuration file path.

You can easily use Blatta with different CockroachDB clusters by create a 
separate configuration file per cluster. For example, create the configuration file
`~/.blatta.foo.yml` and run the blatta command with `--config ~/.blatta.foo.yml` to
use the specific parameters needed to connect to the `foo` cluster.

## SQL Client using Configuration File

To simplify managing multiple cluster configurations, Blatta can connect you using 
`cockroach sql` using the following command:

```bash
blatta sql --configuration /path/to/configuration/file
```

## Monitor Hot Ranges

Blatta can monitor a CockroachDB cluster for hot ranges at specific intervals by leveraging the cluster API.

```bash
blatta monitor hotRanges --url [Cluster API URL] \
  --username [username] --password [password] \
  --pgurl [connection string] \
  (--insecure)
  (--wait [seconds=30]) \
  (--count [count=infinite]) (--insecure)
```

## Settings

### Populate

```
export CLIENT=jon-blatta
rp create $CLIENT -n 1 --clouds=aws --aws-zones=us-east-1b:1 --aws-machine-type-ssd=m5d.xlarge

GOOS=linux GOARCH=amd64 go build -o bin/blatta
rp put $CLIENT bin/blatta
export BLATTA_URL=[serverless instance URL]
rp ssh $CLIENT "./blatta settings save --version recent-50 --url $BLATTA_URL > blatta.log 2>&1 &"
rp ssh $CLIENT "tail -f blatta.log"

```
