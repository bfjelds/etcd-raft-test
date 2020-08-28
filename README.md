# etcd-raft-test
try out etcd raft implementation

## Setup
From the directory containing this README.md:
```bash

TEST_BINARY_NAME="etcd-raft-test"
go build -o $TEST_BINARY_NAME
mkdir test
cd test

```

## Usage for static cluster
```bash

TEST_BINARY_PATH="../$TEST_BINARY_NAME"
STATIC_LIST="http://127.0.0.1:12379,http://127.0.0.1:22379,http://127.0.0.1:32379"
( $TEST_BINARY_PATH --id 1 --cluster "$STATIC_LIST" --port 12380 1>&1.log ) &
( $TEST_BINARY_PATH --id 2 --cluster "$STATIC_LIST" --port 22380 1>&2.log ) &
( $TEST_BINARY_PATH --id 3 --cluster "$STATIC_LIST" --port 32380 1>&3.log ) &
sleep 2s
cat *.log | grep leader

```

## Dynamically add to static cluster
```bash

STATIC_LIST_AS_ARRAY=($(echo $STATIC_LIST | tr "," "\n"))
LEADER_INDEX=$(cat 1.log | grep leader | tail -1 | sed 's/.*elected leader \([0-9]*\).*/\1/')
LEADER_PORT=$(ps -aux | grep $TEST_BINARY_NAME | grep "\-\-id $LEADER_INDEX" | sed 's/.* --port \([0-9]*\).*/\1/')
LEADER_URL="http://127.0.0.1:$LEADER_PORT"

NEW_URL="http://127.0.0.1:42379"
NEW_STATIC_LIST="$STATIC_LIST,$NEW_URL"
curl -L "$LEADER_URL/4" -XPOST -d "$NEW_URL"
( $TEST_BINARY_PATH --id 4 --cluster "$NEW_STATIC_LIST" --port 42380 --join 1>&4.log ) &

```

## Force election
```bash

LEADER_INDEX=$(cat 1.log | grep leader | tail -1 | sed 's/.*elected leader \([0-9]*\).*/\1/')
LEADER_PID=$(ps -aux | grep $TEST_BINARY_NAME | grep "\-\-id $LEADER_INDEX" | awk '{print $2}')
kill $LEADER_PID

sleep 2s
cat *.log | grep leader

```

## Cleanup
```bash

kill $(ps | grep $TEST_BINARY_NAME | awk '{print $1}')
rm -rf *.log
rm -rf raftexample-*

```
