package config

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFromEnvironmentVariables(t *testing.T) {
	envGroups := fmt.Sprintf("%s_GROUPS", EnvironmentVariablePrefix)
	os.Setenv(envGroups, "my Selection")
	defer os.Unsetenv(envGroups)

	envGroupSource := fmt.Sprintf("%s_MY_SELECTION_SOURCE", EnvironmentVariablePrefix)
	os.Setenv(envGroupSource, "gitlab")
	defer os.Unsetenv(envGroupSource)

	envGroupNamespaces := fmt.Sprintf("%s_MY_SELECTION_NAMESPACES", EnvironmentVariablePrefix)
	os.Setenv(envGroupNamespaces, "foo,bar")
	defer os.Unsetenv(envGroupNamespaces)

	envGroupRepositories := fmt.Sprintf("%s_MY_SELECTION_REPOSITORIES", EnvironmentVariablePrefix)
	os.Setenv(envGroupRepositories, "abc,cdf")
	defer os.Unsetenv(envGroupRepositories)

	envGroupUsers := fmt.Sprintf("%s_MY_SELECTION_USERS", EnvironmentVariablePrefix)
	os.Setenv(envGroupUsers, "Alice,Bob")
	defer os.Unsetenv(envGroupUsers)

	selections := FromEnvironmentVariables()
	assert.Len(t, selections, 1)
	selection := selections[0]
	assert.Equal(t, selection.Title, "my Selection")
	assert.Equal(t, selection.Source, "gitlab")
	assert.Len(t, selection.Namespaces, 2)
	assert.Contains(t, selection.Namespaces, "foo")
	assert.Contains(t, selection.Namespaces, "bar")
	assert.Len(t, selection.Repositories, 2)
	assert.Contains(t, selection.Repositories, "abc")
	assert.Contains(t, selection.Repositories, "cdf")
	assert.Len(t, selection.Users, 2)
	assert.Contains(t, selection.Users, "Alice")
	assert.Contains(t, selection.Users, "Bob")
}
