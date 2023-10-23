package main

import (
	"bufio"
	_ "embed"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"regexp"
	"strings"
)

//go:embed sample.audius.conf
var confExample string
var confFilePath string
var port int
var tlsPort int

// with the intent of reducing configuration,
// the latest audius-docker-compose sha (from stage branch) is set at build time via ci.
// this bakes the (tested) image dependency in, so we know that the built binary will always work.
var imageTag = "stage"

func main() {
	fmt.Println(fmt.Sprintf("imageTag: audius/audius-docker-compose:%s", imageTag))

	flag.StringVar(&confFilePath, "c", "", "Path to the .conf file")
	flag.IntVar(&port, "port", 80, "specify a custom http port")
	flag.IntVar(&tlsPort, "tls", 443, "specify a custom https port")

	if !regexp.MustCompile(`^[a-zA-Z0-9_\-]+$`).MatchString(imageTag) {
		exitWithError("Invalid image tag:", imageTag)
	}
	cmdName := "up"
	if len(os.Args) > 1 {
		cmdName = os.Args[1]
	}
	flag.Parse()

	switch cmdName {
	case "down":
		runDown()
	default:
		runUp(checkConfigFile())
	}
}

func checkConfigFile() string {
	nodeType := "discovery-provider"

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
	defer file.Close()
	if err != nil {
		exitWithError("Error opening config file:", err)
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "creatorNodeEndpoint") {
			nodeType = "creator-node"
			break
		}
	}
	if err := scanner.Err(); err != nil {
		exitWithError("Error reading config file:", err)
	}

	return nodeType
}

func runUp(nodeType string) {
	ensureDirectory("/tmp/dind")

	volumeFlag := ""
	if confFilePath != "" {
		volumeFlag = fmt.Sprintf("-v %s:/root/audius-docker-compose/%s/override.env", confFilePath, nodeType)
	}

	var cmd string
	baseCmd := fmt.Sprintf(`docker run --privileged -d -v /tmp/dind:/var/lib/docker %s -p %d:80 -p %d:443`, volumeFlag, port, tlsPort)

	if nodeType == "creator-node" {
		cmd = fmt.Sprintf(baseCmd + ` \
        --name creator-node \
        -v /var/k8s/mediorum:/var/k8s/mediorum \
        -v /var/k8s/creator-node-backend:/var/k8s/creator-node-backend \
        -v /var/k8s/creator-node-db:/var/k8s/creator-node-db \
        audius/audius-docker-compose:` + imageTag)
	} else {
		cmd = fmt.Sprintf(baseCmd + ` \
        --name discovery-provider \
        -v /var/k8s/discovery-provider-db:/var/k8s/discovery-provider-db \
        -v /var/k8s/discovery-provider-chain:/var/k8s/discovery-provider-chain \
        audius/audius-docker-compose:` + imageTag)
	}

	execCmd := fmt.Sprintf(`docker exec %s sh -c "while ! docker ps &> /dev/null; do echo 'starting up' && sleep 1; done && cd %s && docker compose up -d"`, nodeType, nodeType)

	if err := runCommand("/bin/sh", "-c", cmd+" && "+execCmd); err != nil {
		exitWithError("Error executing command:", err)
	}
}

func runDown() {
	runCommand("docker", "rm", "-f", "creator-node", "discovery-provider")
}

func ensureDirectory(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.MkdirAll(path, 0755); err != nil {
			exitWithError("Failed to create directory:", err)
		}
	}
}

func audiusCli(args ...string) {
	audCli := []string{"exec", "discovery-provider", "audius-cli"}
	cmds := append(audCli, args...)
	err := runCommand("docker", cmds...)
	if err != nil {
		exitWithError("Error with audius-cli:", err)
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
