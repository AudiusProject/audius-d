package infra

import (
	"context"
	"fmt"
	"os"

	"github.com/pulumi/pulumi/sdk/v3/go/auto"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/optdestroy"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/optpreview"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/optup"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

var (
	pulumiUserName = "endline"
	projectName    = "audius-d"
	stackName      = "devnet-tiki"
	fqStackName    = fmt.Sprintf("%s/%s/%s", pulumiUserName, projectName, stackName)
)

func getOrInitStack(ctx context.Context, pulumiFunc pulumi.RunFunc) (*auto.Stack, error) {
	s, err := auto.UpsertStackInlineSource(ctx, fqStackName, projectName, pulumiFunc)
	if err != nil {
		return nil, fmt.Errorf("failed to create or select stack: %w", err)
	}
	return &s, nil
}

func Update(ctx context.Context, preview bool) error {
	s, err := getOrInitStack(ctx, func(pCtx *pulumi.Context) error {
		instance, privateKeyFilePath, err := CreateEC2Instance(pCtx, "audius-d-devnet-test")
		if err != nil {
			return err
		}
		pCtx.Export("instancePublicIp", instance.PublicIp)
		pCtx.Export("privateKeyFilePath", pulumi.String(privateKeyFilePath))
		return nil
	})
	if err != nil {
		return err
	}

	if preview {
		_, err = s.Preview(ctx, optpreview.ProgressStreams(os.Stdout))
		if err != nil {
			return fmt.Errorf("failed to preview changes: %w", err)
		}
	} else {
		_, err = s.Up(ctx, optup.ProgressStreams(os.Stdout))
		if err != nil {
			return fmt.Errorf("failed to run up: %w", err)
		}
	}

	return nil
}

func Destroy(ctx context.Context) error {
	s, err := getOrInitStack(ctx, func(pCtx *pulumi.Context) error {
		return nil
	})
	if err != nil {
		return err
	}

	_, err = s.Destroy(ctx, optdestroy.ProgressStreams(os.Stdout))
	if err != nil {
		return fmt.Errorf("failed to destroy stack: %w", err)
	}

	return nil
}

func GetStackOutput(ctx context.Context, outputName string) (string, error) {
	s, err := getOrInitStack(ctx, func(pCtx *pulumi.Context) error {
		return nil
	})
	if err != nil {
		return "", fmt.Errorf("failed to get or init stack: %w", err)
	}
	outputs, err := s.Outputs(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get outputs: %w", err)
	}
	output, ok := outputs[outputName]
	if !ok {
		return "", fmt.Errorf("output %s not found", outputName)
	}
	value, ok := output.Value.(string)
	if !ok {
		return "", fmt.Errorf("output %s is not a string", outputName)
	}

	return value, nil
}
