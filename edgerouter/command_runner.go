package edgerouter

import (
	"log"

	"golang.org/x/crypto/ssh"
)

type Client struct {
	session *ssh.Session
}

func (client *Client) Session() *ssh.Session {
	if client.session != nil {
		return client.session
	}

	client.session = nil // TODO Create the session here
	return client.session
}

func (client *Client) read(configPath string) (*string, error) {
	output, err := client.Session().CombinedOutput("/opt/terraform-provider-edgerouter/read.sh " + configPath)

	if err != nil {
		log.Println("[WARN] Could not read config path " + configPath + ". Error was: " + err.Error())
		return nil, err
	}

	node := string(output)

	if node == "Specified configuration path is not valid" {
		log.Println("[DEBUG] No node found at " + configPath)
		return nil, nil
	}

	return &node, nil
}
