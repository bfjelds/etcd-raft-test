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
echo "Add 1 to cluster: $NEW_STATIC_LIST"
( $TEST_BINARY_PATH --id 1 --cluster "$STATIC_LIST" --port 12380 1>&1.log ) &
echo "Add 2 to cluster: $NEW_STATIC_LIST"
( $TEST_BINARY_PATH --id 2 --cluster "$STATIC_LIST" --port 22380 1>&2.log ) &
echo "Add 3 to cluster: $NEW_STATIC_LIST"
( $TEST_BINARY_PATH --id 3 --cluster "$STATIC_LIST" --port 32380 1>&3.log ) &
sleep 2s
cat *.log | grep "became leader" | sort

```

## Dynamically add to static cluster
```bash

STATIC_LIST_AS_ARRAY=($(echo $STATIC_LIST | tr "," "\n"))
LEADER_INDEX=$(cat *.log | grep "became leader" | sort | tail -1 | sed 's/.* \([0-9]*\) became leader at term .*/\1/')
LEADER_PORT=$(ps -aux | grep $TEST_BINARY_NAME | grep "\-\-id $LEADER_INDEX" | sed 's/.* --port \([0-9]*\).*/\1/')
LEADER_URL="http://127.0.0.1:$LEADER_PORT"

NEW_URL="http://127.0.0.1:42379"
NEW_STATIC_LIST="$STATIC_LIST,$NEW_URL"
echo "Add 4 to cluster: $NEW_STATIC_LIST"
curl -L "$LEADER_URL/4" -XPOST -d "$NEW_URL"
( $TEST_BINARY_PATH --id 4 --cluster "$NEW_STATIC_LIST" --port 42380 --join 1>&4.log ) &

```

## Force election by removing leader
```bash

LEADER_INDEX=$(cat *.log | grep "became leader" | sort | tail -1 | sed 's/.* \([0-9]*\) became leader at term .*/\1/')
NON_LEADER_PORT=$(ps -aux | grep $TEST_BINARY_NAME | grep "\-\-id" | grep -v "\-\-id $LEADER_INDEX" | tail -1 | sed 's/.* --port \([0-9]*\).*/\1/')
NON_LEADER_URL="http://127.0.0.1:$NON_LEADER_PORT"

echo "Non-leader URL: $NON_LEADER_URL"
echo "Deleting leader: $LEADER_INDEX"
curl -L "$NON_LEADER_URL/$LEADER_INDEX" -XDELETE

sleep 2s
cat *.log | grep "4 elected leader" | sort
cat *.log | grep "became leader" | sort

```

## Cleanup
```bash

kill $(ps | grep $TEST_BINARY_NAME | awk '{print $1}')
rm -rf *.log
rm -rf raftexample-*

```
