package workers

import (
	"crypto/tls"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRedisPoolConfig(t *testing.T) {
	// Tests redis pool size which defaults to 1
	opts, err := processOptions(Options{
		ServerAddr: "localhost:6379",
		ProcessID:  "2",
	})

	assert.NoError(t, err)
	assert.Equal(t, 1, opts.client.Options().PoolSize)

	opts, err = processOptions(Options{
		ServerAddr: "localhost:6379",
		ProcessID:  "1",
		PoolSize:   20,
	})

	assert.NoError(t, err)
	assert.Equal(t, 20, opts.client.Options().PoolSize)
}

func TestRedisPoolConfigTLS(t *testing.T) {
	opts, err := processOptions(Options{
		ServerAddr: "localhost:6379",
		ProcessID:  "1",
		PoolSize:   20,
	})

	assert.NoError(t, err)
	assert.Nil(t, opts.client.Options().TLSConfig)

	opts, err = processOptions(Options{
		ServerAddr:     "localhost:6379",
		ProcessID:      "1",
		PoolSize:       20,
		RedisTLSConfig: &tls.Config{ServerName: "test_tls"},
	})

	assert.NoError(t, err)
	assert.NotNil(t, opts.client.Options().TLSConfig)
	assert.Equal(t, "test_tls", opts.client.Options().TLSConfig.ServerName)
}

func TestCustomProcessConfig(t *testing.T) {
	opts, err := processOptions(Options{
		ServerAddr: "localhost:6379",
		ProcessID:  "1",
	})

	assert.NoError(t, err)
	assert.Equal(t, "1", opts.ProcessID)

	opts, err = processOptions(Options{
		ServerAddr: "localhost:6379",
		ProcessID:  "2",
	})

	assert.NoError(t, err)
	assert.Equal(t, "2", opts.ProcessID)
}

func TestRequiresRedisConfig(t *testing.T) {
	_, err := processOptions(Options{ProcessID: "2"})

	assert.Error(t, err, "Configure requires either the Server or Sentinels option")
}

func TestRequiresProcessConfig(t *testing.T) {
	_, err := processOptions(Options{ServerAddr: "localhost:6379"})

	assert.Error(t, err, "Configure requires a ProcessID, which uniquely identifies this instance")
}

func TestAddsColonToNamespace(t *testing.T) {
	opts, err := processOptions(Options{
		ServerAddr: "localhost:6379",
		ProcessID:  "1",
	})

	assert.NoError(t, err)
	assert.Equal(t, "", opts.Namespace)

	opts, err = processOptions(Options{
		ServerAddr: "localhost:6379",
		ProcessID:  "1",
		Namespace:  "prod",
	})

	assert.NoError(t, err)
	assert.Equal(t, "prod:", opts.Namespace)
}

func TestDefaultPollIntervalConfig(t *testing.T) {
	opts, err := processOptions(Options{
		ServerAddr: "localhost:6379",
		ProcessID:  "1",
	})

	assert.NoError(t, err)
	assert.Equal(t, 15*time.Second, opts.PollInterval)

	opts, err = processOptions(Options{
		ServerAddr:   "localhost:6379",
		ProcessID:    "1",
		PollInterval: time.Second,
	})

	assert.NoError(t, err)
	assert.Equal(t, time.Second, opts.PollInterval)
}

func TestSentinelConfigGood(t *testing.T) {
	opts, err := processOptions(Options{
		SentinelAddrs:   "localhost:26379,localhost:46379",
		RedisMasterName: "123",
		ProcessID:       "1",
		PollInterval:    time.Second,
	})

	assert.NoError(t, err)
	assert.Equal(t, "FailoverClient", opts.client.Options().Addr)
	assert.Nil(t, opts.client.Options().TLSConfig)
}

func TestSentinelConfigGoodTLS(t *testing.T) {
	opts, err := processOptions(Options{
		SentinelAddrs:   "localhost:26379,localhost:46379",
		RedisMasterName: "123",
		ProcessID:       "1",
		PollInterval:    time.Second,
		RedisTLSConfig:  &tls.Config{ServerName: "test_tls"},
	})

	assert.NoError(t, err)
	assert.Equal(t, "FailoverClient", opts.client.Options().Addr)
	assert.NotNil(t, opts.client.Options().TLSConfig)
	assert.Equal(t, "test_tls", opts.client.Options().TLSConfig.ServerName)
}

func TestSentinelConfigNoMaster(t *testing.T) {
	_, err := processOptions(Options{
		SentinelAddrs: "localhost:26379,localhost:46379",
		ProcessID:     "1",
		PollInterval:  time.Second,
	})

	assert.Error(t, err)
}
