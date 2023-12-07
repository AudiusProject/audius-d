package orchestration

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/AudiusProject/audius-d/pkg/conf"
	"github.com/joho/godotenv"
)

type OverrideEnv = map[string]string

// deploys a server node generically
func RunNode(nconf conf.NetworkConfig, serverConfig conf.BaseServerConfig, override OverrideEnv, containerName string, nodeType string, internalVolumes []string) error {
	if isContainerRunning(containerName) {
		fmt.Printf("container %s already running\n", containerName)
		return nil
	}

	if isContainerNameInUse(containerName) {
		fmt.Printf("container %s already exists, removing and starting with current config\n", containerName)
		if err := removeContainerByName(containerName); err != nil {
			return err
		}
	}

	imageTag := fmt.Sprintf("audius/audius-docker-compose:%s", nconf.Tag)
	externalVolume := fmt.Sprintf("audius-d-%s", containerName)
	httpPorts := fmt.Sprintf("-p %d:%d", serverConfig.ExternalHttpPort, serverConfig.InternalHttpPort)
	httpsPorts := fmt.Sprintf("-p %d:%d", serverConfig.ExternalHttpsPort, serverConfig.InternalHttpsPort)
	formattedInternalVolumes := " -v " + strings.Join(internalVolumes, " -v ")

	// assemble wrapper command and run
	// todo: handle https port
	upCmd := fmt.Sprintf("docker run --privileged -d -v %s:/var/lib/docker %s %s --name %s %s %s", externalVolume, httpPorts, httpsPorts, containerName, formattedInternalVolumes, imageTag)
	if err := Sh(upCmd); err != nil {
		return err
	}

	// initialize override.env file
	localOverridePath := fmt.Sprintf("./%s-override.env", containerName)
	if err := godotenv.Write(override, localOverridePath); err != nil {
		return err
	}

	envCmd := fmt.Sprintf("docker cp %s %s:/root/audius-docker-compose/%s/override.env", localOverridePath, containerName, nodeType)
	if err := Sh(envCmd); err != nil {
		return err
	}

	cmd := fmt.Sprintf(`docker exec %s sh -c "while ! docker ps &> /dev/null; do echo 'starting up' && sleep 1; done"`, containerName)
	if err := runCommand("/bin/sh", "-c", cmd); err != nil {
		return err
	}

	if err := os.Remove(localOverridePath); err != nil {
		return err
	}

	// assemble inner command and run
	startCmd := fmt.Sprintf(`docker exec %s sh -c "cd %s && docker compose up -d"`, containerName, nodeType)
	if err := Sh(startCmd); err != nil {
		return err
	}

	return nil
}

func Sh(cmd string) error {
	fmt.Println(cmd)
	return runCommand("/bin/sh", "-c", cmd)
}

func isContainerRunning(containerName string) bool {
	cmd := exec.Command("docker", "ps", "-q", "-f", "name="+containerName)
	output, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
		return false
	}
	return len(strings.TrimSpace(string(output))) > 0
}

func isContainerNameInUse(containerName string) bool {
	cmd := exec.Command("docker", "ps", "-a", "--format", "{{.Names}}")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
		return false
	}

	// Split the output into individual container names
	containers := strings.Split(strings.TrimSpace(string(output)), "\n")

	// Check if the given container name exists in the list
	for _, name := range containers {
		if name == containerName {
			return true
		}
	}
	return false
}

func removeContainerByName(containerName string) error {
	cmd := exec.Command("docker", "rm", "-f", containerName)
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func startDevnetDocker() {
	fmt.Println("Starting local eth, sol, and acdc chains")
	runCommand("docker", "compose", "-f", "./deployments/docker-compose.devnet.yml", "up", "-d")
}

func downDevnetDocker() {
	runCommand("docker", "compose", "-f", "./deployments/docker-compose.devnet.yml", "down")
}

func audiusCli(container string, args ...string) {
	audCli := []string{"exec", container, "audius-cli"}
	cmds := append(audCli, args...)
	err := runCommand("docker", cmds...)
	if err != nil {
		exitWithError("Error with audius-cli:", err)
	}
}

func dockerExec(arg ...string) string {
	nodeType := ""
	baseCmd := []string{"exec", nodeType}
	cmds := append(baseCmd, arg...)
	out, err := exec.Command("docker", cmds...).Output()
	if err != nil {
		exitWithError("Error with cmd:", err, cmds)
	}
	return string(out)
}

func awaitDockerStart() {
	nodeType := ""
	cmd := fmt.Sprintf(`docker exec %s sh -c "while ! docker ps &> /dev/null; do echo 'starting up' && sleep 1; done"`, nodeType)
	if err := runCommand("/bin/sh", "-c", cmd); err != nil {
		exitWithError("Error awaiting docker start:", err)
	}

}

func exitWithError(msg ...interface{}) {
	fmt.Println(msg...)
	os.Exit(1)
}

// generates relevant nethermind chain configuration files
// logic ported over from audius-docker-compose https://github.com/AudiusProject/audius-docker-compose/blob/stage/audius-cli#L848
func configureChainSpec(nodeType string, network string) {
	extraVanity := "0x22466c6578692069732061207468696e6722202d204166726900000000000000"
	extraSeal := "0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"

	// gather env config
	// discovery-provider/stage.env for example
	networkEnvPath := fmt.Sprintf("discovery-provider/%s.env", network)
	networkEnv := dockerExec("cat", networkEnvPath)
	networkEnvMap, err := godotenv.Unmarshal(networkEnv)
	if err != nil {
		exitWithError("Error unmarshalling network env:", err)
	}

	signers := networkEnvMap["audius_genesis_signers"]
	extraData := fmt.Sprintf("%s%s%s", extraVanity, signers, extraSeal)

	specTemplatePath := fmt.Sprintf("discovery-provider/chain/%s_spec_template.json", network)
	specInput := dockerExec("cat", specTemplatePath)
	var specData map[string]interface{}
	err = json.Unmarshal([]byte(specInput), &specData)
	if err != nil {
		exitWithError("Unmarshall error:", err)
	}

	networkId := specData["params"].(map[string]interface{})["networkID"].(string)
	fmt.Printf("Network id: %s\n", networkId)

	specData["genesis"].(map[string]interface{})["extraData"] = extraData

	specOutput, err := json.MarshalIndent(specData, "", "    ")
	if err != nil {
		exitWithError("Error marshalling specData:", err)
	}

	peersStr := networkEnvMap["audius_static_nodes"]
	peers := strings.Split(peersStr, ",")
	peersOutput, err := json.MarshalIndent(peers, "", "    ")
	if err != nil {
		exitWithError("Error marshalling peers output:", err)
	}

	err = os.WriteFile("spec.json", specOutput, 0644)
	if err != nil {
		exitWithError("Error writing spec", err)
	}

	err = os.WriteFile("static-nodes.json", peersOutput, 0644)
	if err != nil {
		exitWithError("Error writing static nodes", err)
	}

	// docker cp ./spec.json creator-node:/root/audius-docker-compose/discovery-provider/chain
	err = exec.Command("docker", "cp", "./spec.json", fmt.Sprintf("%s:/root/audius-docker-compose/discovery-provider/chain", nodeType)).Run()
	if err != nil {
		exitWithError("Error with spec docker cp:", err)
	}

	err = exec.Command("docker", "cp", "./static-nodes.json", fmt.Sprintf("%s:/root/audius-docker-compose/discovery-provider/chain", nodeType)).Run()
	if err != nil {
		exitWithError("Error with static nodes docker cp:", err)
	}

	// cleanup, remove temp files from filesystem
	err = os.Remove("spec.json")
	if err != nil {
		exitWithError(err)
	}
	err = os.Remove("static-nodes.json")
	if err != nil {
		exitWithError(err)
	}
}
