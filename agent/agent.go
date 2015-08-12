package agent

import (
	"fmt"

	"github.com/hashicorp/vault/api"
	"github.com/mitchellh/mapstructure"
)

// SSHAgent is used to perform authentication related operations
type Agent struct {
	c    *api.Client
	Path string
}

// SSHAgent is used to return the client for authentication related API calls.
func SSHAgent(c *api.Client, path string) *Agent {
	return &Agent{
		c:    c,
		Path: path,
	}
}

// SSHVerifyResp is a structure representing the fields in vault server's
// response.
type SSHVerifyResponse struct {
	Message  string `mapstructure:"message"`
	Username string `mapstructure:"username"`
	IP       string `mapstructure:"ip"`
}

// Verifies if the key provided by user is present in vault server. If yes,
// the response will contain the IP address and username associated with the
// key.
func (c *Agent) Verify(otp string) (*SSHVerifyResponse, error) {
	data := map[string]interface{}{
		"otp": otp,
	}
	verifyPath := fmt.Sprintf("/v1/%s/verify", c.Path)
	r := c.c.NewRequest("PUT", verifyPath)
	if err := r.SetJSONBody(data); err != nil {
		return nil, err
	}

	resp, err := c.c.RawRequest(r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	secret, err := api.ParseSecret(resp.Body)
	if err != nil {
		return nil, err
	}

	if secret.Data == nil {
		return nil, nil
	}

	var verifyResp SSHVerifyResponse
	err = mapstructure.Decode(secret.Data, &verifyResp)
	if err != nil {
		return nil, err
	}
	return &verifyResp, nil
}