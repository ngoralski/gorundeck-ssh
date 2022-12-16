package main

import (
	"bytes"
	"fmt"
	"github.com/helloyi/go-sshclient"
	"golang.org/x/crypto/ssh"
	"log"
	"os"
	"regexp"
)

func main() {

	var (
		stdin          bytes.Buffer
		stdout         bytes.Buffer
		stderr         bytes.Buffer
		BastionKeyFile ssh.Signer
		TargetKeyFile  ssh.Signer
		err            error
		targetHost     string
		targetUser     string
		targetPort     string
		targetCmd      string
		bastionUser    string
		bastionHost    string
		bastionPort    string
	)

	targetUser = os.Getenv("RD_NODE_USERNAME")
	targetHost = os.Getenv("RD_NODE_HOSTNAME")
	targetPort = os.Getenv("RD_NODE_SSH_PORT")
	targetCmd = os.Getenv("RD_EXEC_COMMAND")
	bastionUser = os.Getenv("RD_CONFIG_BASTION_USERNAME")
	bastionHost = os.Getenv("RD_CONFIG_BASTION_HOST")
	bastionPort = os.Getenv("RD_CONFIG_BASTION_PORT")

	if targetPort == "" {
		targetPort = "22"
	}

	//cfg, _ := ssh_config.Decode(strings.NewReader(os.Getenv("RD_CONFIG_SSH_CONFIG")))

	BastionKeyFile, err = ssh.ParsePrivateKey([]byte(os.Getenv("RD_CONFIG_BASTION_SSH_KEY_STORAGE_PATH")))
	if err != nil {
		log.Fatal(err)
	}

	if os.Getenv("RD_CONFIG_SSH_KEY_STORAGE_PATH") != "" {
		TargetKeyFile, err = ssh.ParsePrivateKey([]byte(os.Getenv("RD_CONFIG_SSH_KEY_STORAGE_PATH")))
		if err != nil {
			log.Fatal(err)
		}
	} else {
		TargetKeyFile = BastionKeyFile
	}

	var configBastion = &ssh.ClientConfig{
		User: bastionUser,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(BastionKeyFile),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	var configTarget = &ssh.ClientConfig{
		User: targetUser,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(TargetKeyFile),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// connect ot ssh server
	connBastion, err := sshclient.Dial("tcp", fmt.Sprintf("%s:%s", bastionHost, bastionPort), configBastion)
	if err != nil {
		log.Fatal(err)
	}

	connTarget, err := connBastion.Dial("tcp", fmt.Sprintf("%s:%s", targetHost, targetPort), configTarget)
	if err != nil {
		log.Fatal(err)
	}

	if err := connTarget.Shell().Start(); err != nil {
		log.Fatal(err)
	}

	if err := connTarget.Shell().SetStdio(&stdin, &stdout, &stderr).Start(); err != nil {
		log.Fatal(err)
	}

	stdin.WriteString(targetCmd)

	err = connTarget.Shell().SetStdio(&stdin, &stdout, &stderr).Start()

	if err != nil {

		sErr := fmt.Sprintf("%v", err)
		switch sErr {
		case "Process exited with status 1":
			matched, _ := regexp.Match(`sudo: a terminal is required to read the password`, []byte(stderr.String()))
			if matched {
				log.Println("Sorry this plugin require that the sudo command is configured without a password")
				os.Exit(3)
			} else {
				fmt.Printf("Command %s, return an error : %v\n", targetCmd, err)
				fmt.Println(&stdout)
				fmt.Println(&stderr)
				log.Fatal(err)
			}

		default:
			fmt.Printf("Command %s, return an error : %v\n", targetCmd, err)
			fmt.Println(&stdout)
			fmt.Println(&stderr)
			os.Exit(2)
		}

	} else {
		fmt.Printf("Command executed sucessfully\n")
		fmt.Println(&stdout)
	}

	// Reset buffer to avoid to re(read|write) previous data
	//stdout.Reset()
	//stdin.Reset()
	//stderr.Reset()

	//targetCmd = fmt.Sprintf("sudo %s", targetCmd)
	//fmt.Printf("Test SUDO %s\n", targetCmd)
	//matched, _ := regexp.Match(`^sudo\s.*`, []byte(targetCmd))
	//if matched {
	//	fmt.Printf("We have a sudo cmd : %s\n", targetCmd)
	//	stdin.WriteString(targetCmd)
	//
	//	err := connTarget.Shell().SetStdio(&stdin, &stdout, &stderr).Start()
	//	//stdin.WriteString("logout")
	//	//err = connTarget.Terminal(nil).SetStdio(&stdin, &stdout, &stderr).Start()
	//	connTarget.Cmd("logout")
	//
	//	fmt.Printf("STDOUT : %s\n", stdout.String())
	//	fmt.Printf("STDERR : %s\n", stderr.String())
	//
	//	if err != nil {
	//
	//		sErr := fmt.Sprintf("%v", err)
	//		switch sErr {
	//		case "Process exited with status 1":
	//			matched, _ := regexp.Match(`sudo: a terminal is required to read the password`, []byte(stderr.String()))
	//			if matched {
	//				log.Fatal("Sorry this plugin require that the sudo command is configured without a password")
	//			} else {
	//				log.Fatal(err)
	//			}
	//
	//		default:
	//			fmt.Printf("Command %s, return an error : %v\n", targetCmd, err)
	//			log.Fatal(err)
	//		}
	//
	//	}
	//	fmt.Printf("Command executed\n")
	//}

	defer connTarget.Close()
	defer connBastion.Close()

}
