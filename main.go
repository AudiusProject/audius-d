package main

import (
	"bufio"
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/AudiusProject/audius-d/conf"
	"github.com/davecgh/go-spew/spew"
	"github.com/joho/godotenv"
)

//go:embed sample.audius.conf
var confExample string
var confFilePath string
var port int
var tlsPort int
var network string
var nodeType string
var seed bool
var autoUpgrade bool

// with the intent of reducing configuration,
// the latest audius-docker-compose sha (from stage branch) is set at build time via ci.
// this bakes the (tested) image dependency in, so we know that the built binary will always work.
var imageTag string

func main() {
	fmt.Println("again")
	toml := conf.ReadTomlUnsafe("devnet.toml")
	spew.Dump(toml)
	// flag.StringVar(&confFilePath, "c", "", "Path to the .conf file")
	// flag.IntVar(&port, "port", 80, "specify a custom http port")
	// flag.IntVar(&tlsPort, "tls", 443, "specify a custom https port")
	// flag.StringVar(&network, "network", "prod", "specify the network to run on")
	// flag.StringVar(&nodeType, "node", "creator-node", "specify the node type to run")
	// flag.BoolVar(&seed, "seed", false, "seed data (only applicable to discovery-provider)")
	// flag.BoolVar(&autoUpgrade, "autoUpgrade", true, "runs cron job to keep node on the latest version")

	// fmt.Printf("imageTag: audius/audius-docker-compose:%s\n", imageTag)

	// if !regexp.MustCompile(`^[a-zA-Z0-9_\-]+$`).MatchString(imageTag) {
	// 	exitWithError("Invalid image tag:", imageTag)
	// }
	// cmdName := "up"
	// if len(os.Args) > 1 {
	// 	cmdName = os.Args[1]
	// }
	// flag.Parse()

	// switch cmdName {
	// case "down":
	// 	runDown()
	// 	downDevnet()
	// case "devnet":
	// 	startDevnet()
	// default:
	// 	fmt.Printf("standing up %s on network %s\n", nodeType, network)
	// 	readConfigFile()
	// 	runUp()
	// }
}

func startDevnet() {
	fmt.Println("starting devnet")
	runCommand("docker", "compose", "-f", "./devnet/docker-compose.yml", "up", "-d")
}

func downDevnet() {
	runCommand("docker", "compose", "-f", "./devnet/docker-compose.yml", "down")
}

func readConfigFile() {
	if confFilePath == "" {
		if usr, err := user.Current(); err != nil {
			exitWithError("Error retrieving current user:", err)
		} else {
			confFilePath = filepath.Join(usr.HomeDir, ".audius", "audius.conf")
		}
	} else {
		if absPath, err := filepath.Abs(confFilePath); err != nil {
			exitWithError("Error creating absolute path to config file:", err)
		} else {
			confFilePath = absPath
		}
	}

	if _, err := os.Stat(confFilePath); os.IsNotExist(err) {
		exitWithError("Config not found at provided location:", confFilePath, confExample)
	}

	file, err := os.Open(confFilePath)
	if err != nil {
		exitWithError("Error opening config file:", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if err := scanner.Err(); err != nil {
		exitWithError("Error reading config file:", err)
	}
}

func runUp() {
	volumeFlag := ""
	if confFilePath != "" {
		volumeFlag = fmt.Sprintf("-v %s:/root/audius-docker-compose/%s/override.env", confFilePath, nodeType)
	}

	// volume create is idempotent
	if err := runCommand("/bin/sh", "-c", "docker volume create audius-d"); err != nil {
		exitWithError("Error executing command:", err)
	}

	var cmd string
	baseCmd := fmt.Sprintf(`docker run --privileged -d -v audius-d:/var/lib/docker %s -p %d:80 -p %d:443`, volumeFlag, port, tlsPort)

	switch nodeType {
	case "creator-node":
		cmd = fmt.Sprintf(baseCmd + ` \
        --name creator-node \
        -v /var/k8s/mediorum:/var/k8s/mediorum \
        -v /var/k8s/creator-node-backend:/var/k8s/creator-node-backend \
        -v /var/k8s/creator-node-db:/var/k8s/creator-node-db \
        audius/audius-docker-compose:` + imageTag)
	case "discovery-provider":
		baseCmd = baseCmd + " -p 5000:5000"
		cmd = fmt.Sprintf(baseCmd + ` \
        --name discovery-provider \
        -v /var/k8s/discovery-provider-db:/var/k8s/discovery-provider-db \
        -v /var/k8s/discovery-provider-chain:/var/k8s/discovery-provider-chain \
        audius/audius-docker-compose:` + imageTag)
	case "identity-service":
		baseCmd = baseCmd + " -p 7000:7000"
		cmd = fmt.Sprintf(baseCmd + ` \
        --name identity-service \
        audius/audius-docker-compose:` + imageTag)
	default:
		exitWithError(fmt.Sprintf("provided node type is not supported: %s", nodeType))
	}

	if err := runCommand("/bin/sh", "-c", cmd); err != nil {
		exitWithError("Error executing command:", err)
	}

	awaitDockerStart()
	audiusCli("set-network", network)

	if nodeType == "discovery-provider" {
		configureChainSpec()
	}

	if autoUpgrade {
		fmt.Println("setting auto-upgrade")
		audiusCli("auto-upgrade")
		fmt.Println("auto-upgrade enabled")
	}

	switch nodeType {
	case "creator-node":
		execCmd := fmt.Sprintf(`docker exec %s sh -c "cd %s && docker compose up -d"`, nodeType, nodeType)
		if err := runCommand("/bin/sh", "-c", execCmd); err != nil {
			exitWithError("Error executing command:", err)
		}
	case "discovery-provider":
		audiusCli("launch-chain")
		launchCmd := []string{"launch", "discovery-provider", "-y"}
		if seed {
			launchCmd = append(launchCmd, "--seed")
		}
		audiusCli(launchCmd...)
	case "identity-service":
		audiusCli("launch", "identity-service", "-y")
	default:
		exitWithError(fmt.Sprintf("provided node type is not supported: %s", nodeType))
	}
}

func runDown() {
	runCommand("docker", "rm", "-f", "creator-node", "discovery-provider", "identity-service")
}

func audiusCli(args ...string) {
	audCli := []string{"exec", nodeType, "audius-cli"}
	cmds := append(audCli, args...)
	err := runCommand("docker", cmds...)
	if err != nil {
		exitWithError("Error with audius-cli:", err)
	}
}

func dockerExec(arg ...string) string {
	baseCmd := []string{"exec", nodeType}
	cmds := append(baseCmd, arg...)
	out, err := exec.Command("docker", cmds...).Output()
	if err != nil {
		exitWithError("Error with cmd:", err, cmds)
	}
	return string(out)
}

func awaitDockerStart() {
	cmd := fmt.Sprintf(`docker exec %s sh -c "while ! docker ps &> /dev/null; do echo 'starting up' && sleep 1; done"`, nodeType)
	if err := runCommand("/bin/sh", "-c", cmd); err != nil {
		exitWithError("Error awaiting docker start:", err)
	}

}

func runCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func exitWithError(msg ...interface{}) {
	fmt.Println(msg...)
	os.Exit(1)
}

// generates relevant nethermind chain configuration files
// logic ported over from audius-docker-compose https://github.com/AudiusProject/audius-docker-compose/blob/stage/audius-cli#L848
func configureChainSpec() {
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
