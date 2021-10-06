# Blatta

Blatta is a command-line utility designed to help with monitoring, troubleshooting and optimizing CockroachDB.

# Commands

## Setup a Configuration File

Blatta configuration parameters can be passed via command-line flags, but it is easier to define a configuration file.

By default, Blatta will look in the user home directory for the file `.blatta.yml`. It can use a different configuration
file that is passed in via the `--config` parameter with the configuration file path.

Here is an example configuration file:

```
url: http://localhost:8080
username: demo
password: demo10166
cacert: /var/folders/rg/_ltnftys4wj6ltzm0y45mhsh0000gq/T/demo657659925/ca.cert
pgurl: postgres://demo:demo10166@127.0.0.1:26257?sslmode=require
insecure: true
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

