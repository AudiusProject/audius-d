package infra

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/AudiusProject/audius-d/pkg/conf"
	"github.com/AudiusProject/audius-d/pkg/logger"
	"github.com/pulumi/pulumi-cloudflare/sdk/v3/go/cloudflare"
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
	confCtxConfig, err = conf.ReadOrCreateContextConfig()
	if err != nil {
		logger.Error("Failed to retrieve context. ", err)
		return
	}

	fqStackName = fmt.Sprintf("%s/%s/%s", confCtxConfig.Network.PulumiUserName, confCtxConfig.Network.PulumiProjectName, confCtxConfig.Network.PulumiStackName)
	log.Println("pkg/infra init :: fqStackName: ", fqStackName)
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

		instanceName := "audius-d-devnet-test.sandbox"
		instance, privateKeyFilePath, err := CreateEC2Instance(pCtx, instanceName)
		if err != nil {
			return err
		}
		pCtx.Export("instancePublicIp", instance.PublicIp)
		pCtx.Export("instancePrivateKeyFilePath", pulumi.String(privateKeyFilePath))

		if confCtxConfig.Network.CloudflareAPIKey != "" && confCtxConfig.Network.CloudflareZoneId != "" {
			provider, err := cloudflare.NewProvider(pCtx, "cloudflareProvider", &cloudflare.ProviderArgs{
				ApiToken: pulumi.StringPtr(confCtxConfig.Network.CloudflareAPIKey),
			})
			if err != nil {
				return fmt.Errorf("failed to create cloudflare provider: %w", err)
			}

			record, err := cloudflare.NewRecord(pCtx, fmt.Sprintf("cf-record-%s", instanceName), &cloudflare.RecordArgs{
				Name:    pulumi.String(instanceName),
				Proxied: pulumi.Bool(true),
				Ttl:     pulumi.Int(1), // Set TTL to automatic (required for proxied)
				Type:    pulumi.String("A"),
				Value:   instance.PublicIp,
				ZoneId:  pulumi.String(confCtxConfig.Network.CloudflareZoneId),
			}, pulumi.Provider(provider))
			if err != nil {
				return fmt.Errorf("failed to create cloudflare record: %w", err)
			}
			pCtx.Export("cloudflareRecordHostname", record.Hostname)
			pCtx.Export("cloudflareRecordValue", record.Value)
		}
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
