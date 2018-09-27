package dockercreds

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	corev1 "k8s.io/api/core/v1"
	"kube.ci/kubeci/pkg/credentials"
)

const annotationPrefix = "credential.kube.ci/docker-"

var config dockerConfig

func flags(fs *flag.FlagSet) {
	config = dockerConfig{make(map[string]entry)}
	fs.Var(&config, "basic-docker", "List of secret=url pairs.")
}

func init() {
	flags(flag.CommandLine)
}

// As the flag is read, this status is populated.
// dockerConfig implements flag.Value
type dockerConfig struct {
	Entries map[string]entry `json:"auths"`
}

func (dc *dockerConfig) String() string {
	if dc == nil {
		// According to flag.Value this can happen.
		return ""
	}
	var urls []string
	for k, v := range dc.Entries {
		urls = append(urls, fmt.Sprintf("%s=%s", v.Secret, k))
	}
	return strings.Join(urls, ",")
}

func (dc *dockerConfig) Set(value string) error {
	parts := strings.Split(value, "=")
	if len(parts) != 2 {
		return fmt.Errorf("expect entries of the form secret=url, got: %v", value)
	}
	secret := parts[0]
	url := parts[1]

	if _, ok := dc.Entries[url]; ok {
		return fmt.Errorf("multiple entries for url: %v", url)
	}

	e, err := newEntry(secret)
	if err != nil {
		return err
	}
	dc.Entries[url] = *e
	return nil
}

type entry struct {
	Secret   string `json:"-"`
	Username string `json:"username"`
	Password string `json:"password"`
	Auth     string `json:"auth"`
	Email    string `json:"email"`
}

func newEntry(secret string) (*entry, error) {
	secretPath := credentials.VolumeName(secret)

	ub, err := ioutil.ReadFile(filepath.Join(secretPath, corev1.BasicAuthUsernameKey))
	if err != nil {
		return nil, err
	}
	username := string(ub)

	pb, err := ioutil.ReadFile(filepath.Join(secretPath, corev1.BasicAuthPasswordKey))
	if err != nil {
		return nil, err
	}
	password := string(pb)

	return &entry{
		Secret:   secret,
		Username: username,
		Password: password,
		Auth:     base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", username, password))),
		Email:    "not@val.id",
	}, nil
}

type dockerConfigBuilder struct{}

// NewBuilder returns a new builder for Docker credentials.
func NewBuilder() credentials.Builder { return &dockerConfigBuilder{} }

// MatchingAnnotations extracts flags for the credential helper
// from the supplied secret and returns a slice (of length 0 or
// greater) of applicable domains.
func (*dockerConfigBuilder) MatchingAnnotations(secret *corev1.Secret) []string {
	var flags []string
	switch secret.Type {
	case corev1.SecretTypeBasicAuth:
		// OK.
	default:
		return flags
	}

	for k, v := range secret.Annotations {
		if strings.HasPrefix(k, annotationPrefix) {
			flags = append(flags, fmt.Sprintf("--basic-docker=%s=%s", secret.Name, v))
		}
	}
	return flags
}

func (*dockerConfigBuilder) Write() error {
	dockerDir := filepath.Join(os.Getenv("HOME"), ".docker")
	dockerConfig := filepath.Join(dockerDir, "config.json")
	if err := os.MkdirAll(dockerDir, os.ModePerm); err != nil {
		return err
	}

	content, err := json.Marshal(config)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(dockerConfig, content, 0600)
}