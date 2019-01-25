package command

import (
	"github.com/alecthomas/kingpin"
	"github.com/levertonai/kubor/common"
	"github.com/levertonai/kubor/kubernetes"
	"gopkg.in/yaml.v2"
	"k8s.io/client-go/tools/clientcmd"
	"os"
)

func init() {
	cmd := &Contexts{}
	common.RegisterCliFactory(cmd)
}

type Contexts struct{}

func (instance *Contexts) ConfigureCliCommands(context string, hc common.HasCommands) error {
	if context != "" {
		return nil
	}
	hc.Command("contexts", "List available contexts").
		Action(func(context *kingpin.ParseContext) error {
			return instance.Run()
		})
	return nil
}

func (instance *Contexts) Run() error {
	config, currentContext, err := kubernetes.NewKubeClientConfig()
	if err != nil {
		return err
	}
	_, err = config.ClientConfig()
	if err != nil {
		return err
	}
	information, err := instance.toContextInformation(config, currentContext)
	if err != nil {
		return err
	}
	encoder := yaml.NewEncoder(os.Stdout)
	return encoder.Encode(information)
}

func (instance *Contexts) toContextInformation(config clientcmd.ClientConfig, currentContext string) (map[string]contextInformation, error) {
	rc, err := config.RawConfig()
	if err != nil {
		return nil, err
	}
	result := map[string]contextInformation{}
	for name, context := range rc.Contexts {
		info := contextInformation{
			User:     context.AuthInfo,
			Selected: currentContext == name,
		}
		cluster := rc.Clusters[context.Cluster]
		if cluster != nil {
			info.Cluster = context.Cluster
			info.Server = cluster.Server
		}
		result[name] = info
	}
	return result, nil
}

type contextInformation struct {
	Selected bool   `yaml:"selected,omitempty"`
	Cluster  string `yaml:"cluster,omitempty"`
	Server   string `yaml:"server,omitempty"`
	User     string `yaml:"user,omitempty"`
}
