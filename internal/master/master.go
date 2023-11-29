package master

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/bugfixes/go-bugfixes/logs"
	"github.com/google/uuid"
	ConfigBuilder "github.com/keloran/go-config"
	"net/http"
)

type Master struct {
	ConfigBuilder.Config
}

type AccountDetails struct {
	Username string
	Password string
	VHost    string
}

func NewMaster(cfg ConfigBuilder.Config) *Master {
	return &Master{
		Config: cfg,
	}
}

func (m *Master) Build() (*AccountDetails, error) {
	v, err := m.CreateVHost()
	if err != nil {
		return nil, logs.Errorf("master: unable to create vhost: %w", err)
	}

	u, p, err := m.CreateAccount(v)
	if err != nil {
		return nil, logs.Errorf("master: unable to create account: %w", err)
	}

	if err := m.AssignVHost(u, v); err != nil {
		return nil, logs.Errorf("master: unable to assign vhost: %w", err)
	}

	return &AccountDetails{
		Username: u,
		Password: p,
		VHost:    v,
	}, nil
}

func (m *Master) CreateAccount(vhost string) (string, string, error) {
	username := uuid.NewString()
	pw := uuid.NewString()
	client := &http.Client{}
	apiURL := fmt.Sprintf("https://%s:%d/api/users/%s", m.Rabbit.ManagementHost, m.Rabbit.Port, username)
	if m.Config.Queue.Port == 0 {
		apiURL = fmt.Sprintf("https://%s/api/users/%s", m.Rabbit.ManagementHost, username)
	}

	type userDetails struct {
		Password string `json:"password"`
		Tags     string `json:"tags"`
		Vhost    string `json:"vhost"`
	}
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(userDetails{
		Password: pw,
		Tags:     "administrator",
	}); err != nil {
		return "", "", logs.Errorf("queue: unable to marshal request: %w", err)
	}

	req, err := http.NewRequest("PUT", apiURL, buf)
	if err != nil {
		return "", "", logs.Errorf("queue: unable to create request: %w", err)
	}
	req.SetBasicAuth(m.Rabbit.Username, m.Rabbit.Password)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return "", "", logs.Errorf("queue: unable to create request: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			logs.Infof("queue: unable to close request: %v", err)
		}
	}()

	return username, pw, nil
}

func (m *Master) CreateVHost() (string, error) {
	vhost := uuid.NewString()
	client := &http.Client{}
	apiURL := fmt.Sprintf("https://%s:%d/api/vhosts/%s", m.Queue.ManagementHost, m.Queue.Port, vhost)
	if m.Config.Queue.Port == 0 {
		apiURL = fmt.Sprintf("https://%s/api/vhosts/%s", m.Queue.ManagementHost, vhost)
	}

	type vHostDetails struct {
		Description string `json:"description"`
		Tags        string `json:"tags"`
	}
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(vHostDetails{
		Description: "vhost created by queue-manager-service",
		Tags:        "production",
	}); err != nil {
		return "", logs.Errorf("CreateVHost: unable to marshal request: %w", err)
	}

	req, err := http.NewRequest("PUT", apiURL, buf)
	if err != nil {
		return "", logs.Errorf("qCreateVHostueue: unable to create request: %w", err)
	}
	req.SetBasicAuth(m.Queue.Username, m.Queue.Password)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return "", logs.Errorf("CreateVHost: unable to create request: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			logs.Infof("CreateVHost: unable to close request: %v", err)
		}
	}()

	return vhost, nil
}

func (m *Master) AssignVHost(username, vhost string) error {
	client := &http.Client{}
	apiURL := fmt.Sprintf("https://%s:%d/api/permissions/%s/%s", m.Queue.ManagementHost, m.Queue.Port, vhost, username)
	if m.Config.Queue.Port == 0 {
		apiURL = fmt.Sprintf("https://%s/api/permissions/%s/%s", m.Queue.ManagementHost, vhost, username)
	}

	type perms struct {
		Configure string `json:"configure"`
		Write     string `json:"write"`
		Read      string `json:"read"`
	}
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(perms{
		Configure: ".*",
		Write:     ".*",
		Read:      ".*",
	}); err != nil {
		return logs.Errorf("CreateVHost: unable to marshal request: %w", err)
	}

	req, err := http.NewRequest("PUT", apiURL, buf)
	if err != nil {
		return logs.Errorf("qCreateVHostueue: unable to create request: %w", err)
	}
	req.SetBasicAuth(m.Queue.Username, m.Queue.Password)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return logs.Errorf("CreateVHost: unable to create request: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			logs.Infof("CreateVHost: unable to close request: %v", err)
		}
	}()

	return nil
}
