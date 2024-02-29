package infra

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/AudiusProject/audius-d/pkg/conf"
	"github.com/AudiusProject/audius-d/pkg/logger"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/auto"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/optdestroy"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/optpreview"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/optup"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

var (
	confCtxConfig *conf.ContextConfig
)

func init() {
	var err error

	baseDir, err := conf.GetConfigBaseDir()
	if err != nil {
		logger.Error("Failed to retrieve config base dir. ", err)
		return
	}

	envVars := map[string]string{
		// local pulumi stacks require this passphrase env var
		// we are not using pulumi secrets - so this is not a security risk
		"PULUMI_CONFIG_PASSPHRASE": "",
		// use a single ~/.audius/.pulumi for all pulumi state management files
		// the .pulumi dir is created upon pulumi login below
		"PULUMI_HOME": fmt.Sprintf("%s/.pulumi", baseDir),
	}
	err = setMultipleEnvVars(envVars)
	if err != nil {
		logger.Error("Error setting environment variables: %v", err)
		return
	}
	cmd := exec.Command("pulumi", "login", fmt.Sprintf("file://%s", baseDir))
	_, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error executing pulumi login:", err)
		return
	}

	confCtxConfig, err = conf.ReadOrCreateContextConfig()
	if err != nil {
		logger.Error("Failed to retrieve context. ", err)
		return
	}
	if confCtxConfig.Network.Infra != nil {
		if confCtxConfig.Network.Infra.PulumiUserName == "" ||
			confCtxConfig.Network.Infra.PulumiProjectName == "" ||
			confCtxConfig.Network.Infra.PulumiStackName == "" {
			logger.Error("Incomplete Pulumi config. ", err)
			return
		}
	}
}

func setMultipleEnvVars(vars map[string]string) error {
	for key, value := range vars {
		err := os.Setenv(key, value)
		if err != nil {
			return err
		}
	}
	return nil
}

func getStack(ctx context.Context, pulumiFunc pulumi.RunFunc) (*auto.Stack, error) {
	stack, err := auto.UpsertStackInlineSource(ctx, confCtxConfig.Network.Infra.PulumiStackName, confCtxConfig.Network.Infra.PulumiProjectName, pulumiFunc)
	if err != nil {
		return nil, fmt.Errorf("failed to create or select stack: %w", err)
	}
	return &stack, nil
}

func Update(ctx context.Context, preview bool) error {
	s, err := getStack(ctx, func(pCtx *pulumi.Context) error {
		// TODO: domain name config
		instanceName := fmt.Sprintf("%s-%s.sandbox.audius.co", confCtxConfig.Network.Infra.PulumiProjectName, confCtxConfig.Network.Infra.PulumiStackName)
		bucketName := fmt.Sprintf("mediorum--%s", instanceName)
		//
		var instance *ec2.Instance
		if awsCredentialsValid(&confCtxConfig.Network) {
			provider, err := awsAuthProvider(pCtx)
			if err != nil {
				return err
			}
			ec2Instance, privateKeyFilePath, err := CreateEC2Instance(pCtx, provider, instanceName)
			if err != nil {
				return err
			}
			instance = ec2Instance
			ec2Instance.PublicIp.ApplyT(func(ip string) error {
				// once we have an IP. block until provisioning completes.
				if err := WaitForUserDataCompletion(privateKeyFilePath, ip); err != nil {
					fmt.Printf("Error waiting for user data completion: %v\n", err)
				}
				return nil
			})
			_, err = CreateS3Bucket(pCtx, provider, bucketName)
			if err != nil {
				return err
			}
		}
		if cloudflareCredentialsValid(&confCtxConfig.Network) {
			provider, err := cloudflareAuthProvider(pCtx)
			if err != nil {
				return err
			}
			zoneId := confCtxConfig.Network.Infra.CloudflareZoneId
			recordName := strings.Replace(instanceName, ".audius.co", "", 1)
			recordIp := instance.PublicIp
			err = CloudflareAddDNSRecord(pCtx, provider, zoneId, recordName, recordIp)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	if preview {
		_, err = s.Preview(ctx, optpreview.ProgressStreams(os.Stdout))
	} else {
		_, err = s.Up(ctx, optup.ProgressStreams(os.Stdout))
	}
	if err != nil {
		return fmt.Errorf("failed to run Update: %w", err)
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
		return fmt.Errorf("failed to run Destroy: %w", err)
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
		return fmt.Errorf("failed to run Cancel: %w", err)
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
