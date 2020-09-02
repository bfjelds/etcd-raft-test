/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package root

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
	"go.etcd.io/etcd/raft"
	"go.etcd.io/etcd/raft/raftpb"

	raftSupport "github.com/bfjelds/etcd-raft-test/raft"
)

var (
	RaftNode raft.Node

	proposeC    chan string
	confChangeC chan raftpb.ConfChange
	cfgFile     string
	cluster     string
	id          int
	port        int
	join        bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "etcd-raft-test",
	Short: "Create a test of the etcd raft implementation",
	Long: `Create a test of the etcd 
	raft implementation. Usage:

etc-raft-test foo &
etc-raft-test bar &
etc-raft-test zoo &
`,
	Run: func(cmd *cobra.Command, args []string) {

		// Start raft business

		proposeC := make(chan string)
		defer close(proposeC)
		confChangeC := make(chan raftpb.ConfChange)
		defer close(confChangeC)

		// raft provides a commit stream for the proposals from the http api
		var kvs *raftSupport.KvStore
		getSnapshot := func() ([]byte, error) { return kvs.GetSnapshot() }
		commitC, errorC, snapshotterReady, rc := raftSupport.NewRaftNode(id, strings.Split(cluster, ","), join, getSnapshot, proposeC, confChangeC)

		kvs = raftSupport.NewKVStore(<-snapshotterReady, proposeC, commitC, errorC)

		// How to poll for current leader
		go func() {
			for {
				select {
				case <-time.After(1 * time.Second):
					if rc.Node != nil {
						status := rc.Node.Status()
						log.Printf("LEADERSHIP-TRACKER-POLLING - Me=%d ... Leader=%d", status.ID, status.SoftState.Lead)
					}

					//do something
				}
			}

		}()

		// the key-value http handler will propose updates to raft
		raftSupport.ServeHttpKVAPI(kvs, port, confChangeC, errorC)

		// 	storage := raft.NewMemoryStorage()

		// 	c := &raft.Config{
		// 		ID:              uint64(id),
		// 		ElectionTick:    10,
		// 		HeartbeatTick:   1,
		// 		Storage:         storage,
		// 		MaxSizePerMsg:   4096,
		// 		MaxInflightMsgs: 256,
		// 	}

		// 	// Create storage and config as shown above.
		// 	// Set peer list to itself, so this node can become the leader of this single-node cluster.
		// 	peers := []raft.Peer{{ID: 0x01}}
		// 	RaftNode = raft.StartNode(c, peers)

		// 	ticker := time.NewTicker(1 * time.Second)
		// 	for {
		// 		select {
		// 		case <- ticker.C:
		// 			RaftNode.Tick()
		// 		case readyc := <- RaftNode.Ready():
		// 			fmt.Println("node is ready")
		// 			storage.Append(readyc.Entries)

		// 			for _, ent := range readyc.CommittedEntries {
		// 				fmt.Println("commit entry:", string(ent.Data))
		// 				switch ent.Type {
		// 				case raftpb.EntryNormal:
		// 					fmt.Println("Normal entry:", ent)

		// 				case raftpb.EntryConfChange:
		// 					fmt.Println("Conf Change entry:", ent)
		// 					var cc raftpb.ConfChange
		// 					cc.Unmarshal(ent.Data)
		// 					RaftNode.ApplyConfChange(cc)
		// 				}
		// 			}
		// 			RaftNode.Advance()
		// 		}
		// 	}
	},
}

// func recvRaftRPC(ctx context.Context, m raftpb.Message) {
// 	n.Step(ctx, m)
// }

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.etcd-raft-test.yaml)")

	rootCmd.PersistentFlags().StringVar(&cluster, "cluster", "http://127.0.0.1:9021", "comma separated cluster peers")
	rootCmd.PersistentFlags().IntVar(&id, "id", 1, "node ID")
	rootCmd.PersistentFlags().IntVar(&port, "port", 9121, "key-value server port")
	rootCmd.PersistentFlags().BoolVar(&join, "join", false, "join an existing cluster")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".etcd-raft-test" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".etcd-raft-test")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
