package infra

import (
	"context"
	"fmt"
	"os"

	"github.com/AudiusProject/audius-d/pkg/conf"
	"github.com/AudiusProject/audius-d/pkg/logger"
	"github.com/pulumi/pulumi/sdk/v3/go/auto"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/optdestroy"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/optpreview"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/optup"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

var (
	fqStackName   string
	confCtxConfig *conf.ContextConfig
)

func init() {
	var err error

	// TODO: local pulumi stacks require this passphrase env var
	// we are not using pulumi encrypted configs (as yet) so ideally we could remove this
	err = os.Setenv("PULUMI_CONFIG_PASSPHRASE", "")
	if err != nil {
		fmt.Println("Error setting environment variable:", err)
	}

	confCtxConfig, err = conf.ReadOrCreateContextConfig()
	if err != nil {
		logger.Error("Failed to retrieve context. ", err)
		return
	}

	fqStackName = fmt.Sprintf("%s-%s-%s", confCtxConfig.Network.PulumiUserName, confCtxConfig.Network.PulumiProjectName, confCtxConfig.Network.PulumiStackName)
	logger.Debug("pkg/infra init :: fqStackName: ", fqStackName)
}

func getStack(ctx context.Context, pulumiFunc pulumi.RunFunc) (*auto.Stack, error) {
	s, err := auto.UpsertStackInlineSource(ctx, fqStackName, confCtxConfig.Network.PulumiProjectName, pulumiFunc)
	if err != nil {
		return nil, fmt.Errorf("failed to create or select stack: %w", err)
	}
	return &s, nil
}

func Update(ctx context.Context, preview bool) error {

	s, err := getStack(ctx, func(pCtx *pulumi.Context) error {
		// TODO: dns hostname is tied to audius sandbox
		instanceName := fmt.Sprintf("%s-%s.sandbox", confCtxConfig.Network.PulumiProjectName, confCtxConfig.Network.PulumiStackName)
		instance, privateKeyFilePath, err := CreateEC2Instance(pCtx, instanceName)
		if err != nil {
			return err
		}

		// TODO: refactor config checks
		if confCtxConfig.Network.CloudflareAPIKey != "" && confCtxConfig.Network.CloudflareZoneId != "" {
			err := ConfigureCloudflare(pCtx, instanceName, instance.PublicIp, confCtxConfig.Network.CloudflareAPIKey, confCtxConfig.Network.CloudflareZoneId)
			if err != nil {
				return err
			}
		}

		// TODO: handle async errors
		pulumi.All(instance.PublicIp).ApplyT(func(all []interface{}) error {
			publicIp := all[0].(string)
			err = WaitForUserDataCompletion(privateKeyFilePath, publicIp)
			if err != nil {
				return fmt.Errorf("%w", err)
			}
			err = RunAudiusD(privateKeyFilePath, publicIp)
			if err != nil {
				return fmt.Errorf("%w", err)
			}
			return nil
		})

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
	s, err := getStack(ctx, func(pCtx *pulumi.Context) error {
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

func Cancel(ctx context.Context) error {
	s, err := getStack(ctx, func(pCtx *pulumi.Context) error {
		return nil
	})
	if err != nil {
		return err
	}

	err = s.Cancel(ctx)
	if err != nil {
		return fmt.Errorf("failed to cancel stack update: %w", err)
	}

	return nil
}

func GetStackOutput(ctx context.Context, outputName string) (string, error) {
	s, err := getStack(ctx, func(pCtx *pulumi.Context) error {
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
