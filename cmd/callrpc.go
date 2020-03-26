package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/ava-labs/avash/cfg"
	pmgr "github.com/ava-labs/avash/processmgr"
	"github.com/spf13/cobra"
	"github.com/ybbus/jsonrpc"
)

// CallRPCCmd issues an RPC to a node endpoint using JSONRPC protocol
var CallRPCCmd = &cobra.Command{
	Use:     "callrpc [node name] [endpoint] [method] [JSON params] [var scope] [var name]",
	Short:   "Issues an RPC to a node.",
	Long:    `Issues an RPC to a node endpoint for the specified method and params.
	Response is saved to the local varstore.`,
	Example: `callrpc n1 ext/avm avm.getBalance {"address":"X-KqpU28P2ipUxfTfwaT847wWxyXB4XuWad","assetID":"AVA"} s v`,
	Args: cobra.MinimumNArgs(6),
	Run: func(cmd *cobra.Command, args []string) {
		log := cfg.Config.Log
		meta, err := pmgr.ProcManager.Metadata(args[0])
		if err != nil {
			log.Error("process not found: %s", args[0])
			return
		}
		var md Metadata
		if err = json.Unmarshal([]byte(meta), &md); err != nil {
			log.Error("unable to unmarshal metadata for process %s: %s", args[0], err.Error())
			return
		}
		jrpcloc := fmt.Sprintf("http://%s:%s/%s", md.Serverhost, md.HTTPport, args[1])
		rpcClient := jsonrpc.NewClient(jrpcloc)
		argMap := make(map[string]interface{})
		if err = json.Unmarshal([]byte(args[3]), &argMap); err != nil {
			log.Error("invalid JSON object: %s", args[3])
			return
		}
		response, err := rpcClient.Call(args[2], argMap)
		if err != nil {
			log.Error("rpcClient returned error: %s", err.Error())
			return
		}
		if response.Error != nil {
			log.Error("rpcClient returned error: %d, %s", response.Error.Code, response.Error.Message)
			return
		}
		resBytes, err := json.Marshal(response.Result)
		if err != nil {
			log.Error("rpcClient returned invalid JSON object: %v", response.Result)
			return
		}
		resVal := string(resBytes)
		log.Info("Response: %s", resVal)
		VarStoreSetCmd.Run(VarStoreSetCmd, []string{args[4], args[5], resVal})
	},
}