module github.com/bfjelds/etcd-raft-test

go 1.14

require (
	github.com/coreos/etcd v3.3.13+incompatible
	github.com/mitchellh/go-homedir v1.1.0
	github.com/spf13/cobra v0.0.3
	github.com/spf13/viper v1.7.1
	go.etcd.io/etcd v0.0.0-20191023171146-3cf2f69b5738
	go.uber.org/multierr v1.5.0 // indirect
	go.uber.org/zap v1.10.0
)

replace go.etcd.io/etcd => go.etcd.io/etcd v0.0.0-20191023171146-3cf2f69b5738
