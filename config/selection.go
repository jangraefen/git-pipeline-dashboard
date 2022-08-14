package config

import (
	"fmt"

	"github.com/jangraefen/git-pipeline-dashboard/util"
)

const (
	EnvironmentVariablePrefix = "PIPELINE_DASHBOARD"
)

type GroupSelectionConfig struct {
	Title        string
	Source       string
	Namespaces   []string
	Repositories []string
	Users        []string
}

type GroupSelectionConfigList []GroupSelectionConfig

func (list GroupSelectionConfigList) ContainsSource(source string) bool {
	for _, selection := range list {
		if selection.Source == source {
			return true
		}
	}

	return false
}

func FromEnvironmentVariables() GroupSelectionConfigList {
	groups := util.GetEnvList(fmt.Sprintf("%s_GROUPS", EnvironmentVariablePrefix))
	selectionConfigs := make([]GroupSelectionConfig, len(groups))

	for i, group := range groups {
		envGroup := fmt.Sprintf("%s_%s", EnvironmentVariablePrefix, util.ToEnvironmentName(group))
		selectionConfigs[i] = GroupSelectionConfig{
			Title:        group,
			Source:       util.GetEnvOrDefault(fmt.Sprintf("%s_SOURCE", envGroup), "gitlab"),
			Namespaces:   util.GetEnvList(fmt.Sprintf("%s_NAMESPACES", envGroup)),
			Repositories: util.GetEnvList(fmt.Sprintf("%s_REPOSITORIES", envGroup)),
			Users:        util.GetEnvList(fmt.Sprintf("%s_USERS", envGroup)),
		}
	}

	return selectionConfigs
}
