# Blatta

Blatta is a command-line utility designed to help with monitoring, troubleshooting and optimizing CockroachDB.

# Commands

## Monitor Hot Ranges

Blatta can monitor a CockroachDB cluster for hot ranges at specific intervals by leveraging the cluster API.

```bash
blatta monitor hotRanges --url [Cluster API URL] --username [username] --password [password] \
  --interval [seconds=30] --qps-threshold [number=750] \
  (--count [count=infinite]) (--insecure)
```