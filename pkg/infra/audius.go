package infra

import (
	"fmt"
	"path/filepath"

	"github.com/AudiusProject/audius-d/pkg/conf"
)

func RunAudiusD(privateKeyFilePath, publicIp string) error {
	// bootstrap audius-ctl
	_, err := ExecuteSSHCommand(
		privateKeyFilePath,
		publicIp,
		"audius-ctl",
	)
	if err != nil {
		return err
	}

	// copy local configs (current context)
	contextBaseDir, err := conf.GetContextBaseDir()
	if err != nil {
		return err
	}
	currentContext, err := conf.GetCurrentContextName()
	if err != nil {
		return err
	}
	err = ExecuteSCPCommand(
		privateKeyFilePath,
		publicIp,
		filepath.Join(contextBaseDir, currentContext),
		filepath.Join("/home/ubuntu/.audius/contexts", currentContext), // TODO remove hardcode base dir this will break
	)
	if err != nil {
		return err
	}
	_, err = ExecuteSSHCommand(
		privateKeyFilePath,
		publicIp,
		fmt.Sprintf("audius-ctl config use-context %s", currentContext),
	)
	if err != nil {
		return err
	}

	// run stack
	_, err = ExecuteSSHCommand(
		privateKeyFilePath,
		publicIp,
		"audius-ctl up",
	)
	if err != nil {
		return err
	}
	return nil
}
