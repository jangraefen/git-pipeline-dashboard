package util

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetEnvOrDefault(t *testing.T) {
	os.Setenv("DEFAULT_ENV_TEST", "some-value")
	defer os.Unsetenv("DEFAULT_ENV_TEST")

	assert.Equal(t, "some-value", GetEnvOrDefault("DEFAULT_ENV_TEST", "default-value"))
	assert.Equal(t, "default-value", GetEnvOrDefault("NONEXISTING_ENV_TEST", "default-value"))
}

func TestGetEnvOrDefaultUint(t *testing.T) {
	os.Setenv("DEFAULT_ENV_TEST", "1337")
	defer os.Unsetenv("DEFAULT_ENV_TEST")
	os.Setenv("DEFAULT_ENV_TEST_INVALID", "some-value")
	defer os.Unsetenv("DEFAULT_ENV_TEST_INVALID")

	assert.Equal(t, uint(1337), GetEnvOrDefaultUint("DEFAULT_ENV_TEST", 42))
	assert.Equal(t, uint(42), GetEnvOrDefaultUint("NONEXISTING_ENV_TEST", 42))
	assert.Equal(t, uint(42), GetEnvOrDefaultUint("DEFAULT_ENV_TEST_INVALID", 42))
}

func TestGetEnvList(t *testing.T) {
	os.Setenv("ENV_LIST_TEST_EMPTY", "")
	defer os.Unsetenv("ENV_LIST_TEST_EMPTY")
	os.Setenv("ENV_LIST_TEST_SINGLE", "some-value")
	defer os.Unsetenv("ENV_LIST_TEST_SINGLE")
	os.Setenv("ENV_LIST_TEST_MULTIPLE", "some-value,other-value")
	defer os.Unsetenv("ENV_LIST_TEST_MULTIPLE")
	os.Setenv("ENV_LIST_TEST_MULTIPLE_TRIM", " some-value , other-value ")
	defer os.Unsetenv("ENV_LIST_TEST_MULTIPLE_TRIM")

	assert.Equal(t, []string{}, GetEnvList("ENV_LIST_TEST_EMPTY"))
	assert.Equal(t, []string{"some-value"}, GetEnvList("ENV_LIST_TEST_SINGLE"))
	assert.Equal(t, []string{"some-value", "other-value"}, GetEnvList("ENV_LIST_TEST_MULTIPLE"))
	assert.Equal(t, []string{"some-value", "other-value"}, GetEnvList("ENV_LIST_TEST_MULTIPLE_TRIM"))
	assert.Equal(t, []string{}, GetEnvList("ENV_LIST_TEST_NONEXISTING"))
}

func TestToEnvironmentVariable(t *testing.T) {
	assert.Equal(t, "TESTINPUT", ToEnvironmentName("TestInput"))
	assert.Equal(t, "TEST_INPUT", ToEnvironmentName("test_input"))
	assert.Equal(t, "TEST_INPUT_NUMBER_2", ToEnvironmentName("Test  Input -- Number 2"))
}
