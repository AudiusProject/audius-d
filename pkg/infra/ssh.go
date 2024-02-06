package infra

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"golang.org/x/crypto/ssh"
)

func setupSSHClient(privateKeyPath, publicIP string) (*ssh.Client, error) {
	key, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("unable to read private key: %v", err)
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, fmt.Errorf("unable to parse private key: %v", err)
	}

	config := &ssh.ClientConfig{
		User:            "ubuntu",
		Auth:            []ssh.AuthMethod{ssh.PublicKeys(signer)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         10 * time.Second,
	}

	client, err := ssh.Dial("tcp", publicIP+":22", config)
	if err != nil {
		return nil, fmt.Errorf("failed to dial: %v", err)
	}

	return client, nil
}

func ExecuteSSHCommand(privateKeyPath, publicIP, command string) (string, error) {
	client, err := setupSSHClient(privateKeyPath, publicIP)
	if err != nil {
		return "", err
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		return "", fmt.Errorf("failed to create session: %v", err)
	}
	defer session.Close()

	out, err := session.CombinedOutput(command)
	if err != nil {
		return "", fmt.Errorf("failed to execute command: %v", err)
	}

	return string(out), nil
}

func ExecuteSCPCommand(privateKeyPath, publicIP, localFilePath, remoteFilePath string) error {
	client, err := setupSSHClient(privateKeyPath, publicIP)
	if err != nil {
		return err
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create session: %v", err)
	}
	defer session.Close()

	srcFile, err := os.Open(localFilePath)
	if err != nil {
		return fmt.Errorf("unable to open local file: %v", err)
	}
	defer srcFile.Close()

	fileInfo, err := srcFile.Stat()
	if err != nil {
		return fmt.Errorf("unable to stat local file: %v", err)
	}

	go func() {
		w, _ := session.StdinPipe()
		defer w.Close()

		fmt.Fprintf(w, "C0644 %d %s\n", fileInfo.Size(), filepath.Base(remoteFilePath))
		io.Copy(w, srcFile)
		fmt.Fprint(w, "\x00")
	}()

	err = session.Run(fmt.Sprintf("scp -t %s", remoteFilePath))
	if err != nil {
		return fmt.Errorf("failed to execute SCP command: %v, %s", err, remoteFilePath)
	}

	return nil
}

func WaitForUserDataCompletion(privateKeyPath, publicIP string) error {
	const timeout = 3 * time.Minute
	const checkInterval = 10 * time.Second
	const completionSignalCommand = "test -f /home/ubuntu/user-data-done && echo 'done' || echo 'not done'"

	startTime := time.Now()
	for {
		if time.Since(startTime) > timeout {
			return fmt.Errorf("timeout waiting for user data script to complete")
		}
		output, err := ExecuteSSHCommand(privateKeyPath, publicIP, completionSignalCommand)
		if err != nil {
			fmt.Println("Error checking for user data completion:", err)
		} else if output == "done\n" {
			fmt.Println("User data script completed successfully.")
			return nil
		}
		fmt.Println("Waiting for instance provisioning to complete...")
		time.Sleep(checkInterval)
	}
}
