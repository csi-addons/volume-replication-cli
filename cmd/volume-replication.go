package cmd

import (
	"context"
	"fmt"
	"strconv"

	replication "github.com/kube-storage/spec/lib/go"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(volumereplication)
	cmds := []*cobra.Command{enable, disable, promote, demote, resync}
	for _, cmd := range cmds {
		volumereplication.AddCommand(cmd)
		cmd.Flags().StringToString("parameters", nil, "parameters to send to backend")
		cmd.Flags().StringToString("secrets", nil, "secrets to connect to backend storage")
		cmd.Flags().Bool("force", false, "Perform force operation")
		cmd.SilenceUsage = true
	}
	volumereplication.SilenceUsage = true
}

var volumereplication = &cobra.Command{
	Use:   `replication`,
	Short: "replication commands",
	Long:  "Commands for replication",
}

var enable = &cobra.Command{
	Use:     `enable`,
	Short:   "enable replication on a volume",
	Long:    "enable replication on a volume",
	Example: `enable --parameters=mirroringMode=snapshot --volumeID=xxx-xxxx-xxxx-xxxx-xxx`,
	RunE:    enableReplication,
}

var disable = &cobra.Command{
	Use:     `disable`,
	Short:   "disable replication on a volume",
	Long:    "disable replication on a volume",
	Example: `disable --force=true --volumeID=xxx-xxxx-xxxx-xxxx-xxx`,
	RunE:    disableReplication,
}

var promote = &cobra.Command{
	Use:     `promote`,
	Short:   "promote volume as primary",
	Long:    "promote volume as primary",
	Example: `promote --force=true --volumeID=xxx-xxxx-xxxx-xxxx-xxx`,
	RunE:    promoteVolume,
}

var demote = &cobra.Command{
	Use:     `demote`,
	Short:   "demote volume as secondary",
	Long:    "demote volume as secondary",
	Example: `demote --force=true --volumeID=xxx-xxxx-xxxx-xxxx-xxx`,
	RunE:    demoteVolume,
}

var resync = &cobra.Command{
	Use:     `resync`,
	Short:   "resync volume to correct split-brain",
	Long:    "resync volume to correct split-brain",
	Example: `resync --force=true --volumeID=xxx-xxxx-xxxx-xxxx-xxx`,
	RunE:    resyncVolume,
}

func enableReplication(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("volume id missing")

	}
	volumeID := args[0]

	parameters, err := cmd.Flags().GetStringToString("parameters")
	if err != nil {
		return err
	}
	secrets, err := cmd.Flags().GetStringToString("secrets")
	if err != nil {
		return err
	}
	grpcClient, err := newGRPCClient()
	if err != nil {
		return err
	}
	repClient := replication.NewControllerClient(grpcClient)

	req := &replication.EnableVolumeReplicationRequest{
		VolumeId:   volumeID,
		Parameters: parameters,
		Secrets:    secrets,
	}

	createCtx, cancel := context.WithTimeout(context.Background(), config.timeout)
	defer cancel()
	_, err = repClient.EnableVolumeReplication(createCtx, req)
	if err != nil {
		return err
	}
	fmt.Println("successfully enabled replication on the volume")
	return nil
}

func disableReplication(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("volume id missing")

	}
	volumeID := args[0]

	parameters, err := cmd.Flags().GetStringToString("parameters")
	if err != nil {
		return err
	}

	secrets, err := cmd.Flags().GetStringToString("secrets")
	if err != nil {
		return err
	}

	force, err := cmd.Flags().GetBool("force")
	if err != nil {
		return err
	}
	grpcClient, err := newGRPCClient()
	if err != nil {
		return err
	}
	parameters["force"] = strconv.FormatBool(force)
	repClient := replication.NewControllerClient(grpcClient)
	req := &replication.DisableVolumeReplicationRequest{
		VolumeId:   volumeID,
		Parameters: parameters,
		Secrets:    secrets,
	}

	createCtx, cancel := context.WithTimeout(context.Background(), config.timeout)
	defer cancel()
	_, err = repClient.DisableVolumeReplication(createCtx, req)
	if err != nil {
		return err
	}
	fmt.Println("successfully disabled replication on the volume")
	return nil
}

func promoteVolume(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("volume id missing")

	}
	volumeID := args[0]

	parameters, err := cmd.Flags().GetStringToString("parameters")
	if err != nil {
		return err
	}

	secrets, err := cmd.Flags().GetStringToString("secrets")
	if err != nil {
		return err
	}

	force, err := cmd.Flags().GetBool("force")
	if err != nil {
		return err
	}
	grpcClient, err := newGRPCClient()
	if err != nil {
		return err
	}
	parameters["force"] = strconv.FormatBool(force)
	repClient := replication.NewControllerClient(grpcClient)
	req := &replication.PromoteVolumeRequest{
		VolumeId:   volumeID,
		Parameters: parameters,
		Secrets:    secrets,
	}

	createCtx, cancel := context.WithTimeout(context.Background(), config.timeout)
	defer cancel()
	_, err = repClient.PromoteVolume(createCtx, req)
	if err != nil {
		return err
	}
	fmt.Println("successfully promoted the volume")
	return nil
}

func demoteVolume(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("volume id missing")

	}
	volumeID := args[0]

	parameters, err := cmd.Flags().GetStringToString("parameters")
	if err != nil {
		return err
	}

	force, err := cmd.Flags().GetBool("force")
	if err != nil {
		return err
	}
	secrets, err := cmd.Flags().GetStringToString("secrets")
	if err != nil {
		return err
	}
	grpcClient, err := newGRPCClient()
	if err != nil {
		return err
	}
	parameters["force"] = strconv.FormatBool(force)
	repClient := replication.NewControllerClient(grpcClient)
	req := &replication.DemoteVolumeRequest{
		VolumeId:   volumeID,
		Parameters: parameters,
		Secrets:    secrets,
	}
	createCtx, cancel := context.WithTimeout(context.Background(), config.timeout)
	defer cancel()
	_, err = repClient.DemoteVolume(createCtx, req)
	if err != nil {
		return err
	}
	fmt.Println("successfully demoted the volume")
	return nil
}

func resyncVolume(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("volume id missing")

	}
	volumeID := args[0]

	parameters, err := cmd.Flags().GetStringToString("parameters")
	if err != nil {
		return err
	}

	secrets, err := cmd.Flags().GetStringToString("secrets")
	if err != nil {
		return err
	}

	force, err := cmd.Flags().GetBool("force")
	if err != nil {
		return err
	}
	grpcClient, err := newGRPCClient()
	if err != nil {
		return err
	}
	parameters["force"] = strconv.FormatBool(force)
	repClient := replication.NewControllerClient(grpcClient)
	req := &replication.ResyncVolumeRequest{
		VolumeId:   volumeID,
		Parameters: parameters,
		Secrets:    secrets,
	}

	createCtx, cancel := context.WithTimeout(context.Background(), config.timeout)
	defer cancel()
	_, err = repClient.ResyncVolume(createCtx, req)
	if err != nil {
		return err
	}
	fmt.Println("successfully resynced the volume")
	return nil
}
