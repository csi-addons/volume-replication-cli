package cmd

import (
	"os"
	"time"

	client "github.com/kube-storage/volume-replication-cli/pkg"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

var (
	config configuration
)

var (
	defaultSockerPath = "/run/csi/socket"
)

// configurations for CLI
type configuration struct {
	timeout    time.Duration
	csiAddress string
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "volume-replication-cli",
	Short: "CLI for volume replication",
	Long: `CLI tool to communicate with the volume replication GRPC server
	to enable/disable the unidirectional or bidirection replication`,
	Run: func(cmd *cobra.Command, args []string) {
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&config.csiAddress, "csi-address", "",
		"The gRPC endpoint to connect to CSI Driver."+
			"\nEnvironment variable CSI_ADDRESS (the default one is /run/csi/socket)")
	rootCmd.PersistentFlags().DurationVar(&config.timeout, "timeout", time.Second*30, "client connection timeout")
	rootCmd.SilenceUsage = true
}

func newGRPCClient() (*grpc.ClientConn, error) {
	// create new client
	grpcClient, err := client.Connect(config.csiAddress)
	if err != nil {
		return nil, err
	}
	return grpcClient, nil
}

// initConfig to set default values
func initConfig() {
	// set driver address from ENV
	if config.csiAddress == "" {
		config.csiAddress = os.Getenv("CSI_ADDRESS")
		if config.csiAddress == "" {
			config.csiAddress = defaultSockerPath
		}
	}

}
