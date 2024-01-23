package infra

import (
	"context"
	"fmt"
	"os"

	"github.com/pulumi/pulumi/sdk/v3/go/auto"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/optdestroy"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/optup"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

var (
	pulumiUserName = "endline"
	projectName    = "audius-d"
	stackName      = "devnet"
	fqStackName    = fmt.Sprintf("%s/%s/%s", pulumiUserName, projectName, stackName)
)

func Update() error {
	ctx := context.Background()

	s, err := auto.UpsertStackInlineSource(ctx, fqStackName, projectName, func(pCtx *pulumi.Context) error {
		instance, err := CreateEC2Instance(pCtx, "audius-d-devnet-test")
		if err != nil {
			return err
		}
		pCtx.Export(fmt.Sprintf("instance-publicIp"), instance.PublicIp)
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to create or select stack: %w", err)
	}

	_, err = s.Up(ctx, optup.ProgressStreams(os.Stdout))
	if err != nil {
		return fmt.Errorf("failed to run up: %w", err)
	}

	return nil
}

func Destroy() error {
	ctx := context.Background()

	s, err := auto.UpsertStackInlineSource(ctx, fqStackName, projectName, func(ctx *pulumi.Context) error {
		// We need to match the creation pattern to obtain the stack reference
		// as there is no Pulumi.yaml defined
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to create or select stack: %w", err)
	}

	_, err = s.Destroy(ctx, optdestroy.ProgressStreams(os.Stdout))
	if err != nil {
		return fmt.Errorf("failed to destroy stack: %w", err)
	}

	return nil
}
