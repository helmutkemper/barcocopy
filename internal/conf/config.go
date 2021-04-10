package conf

import (
	"os"
	"path/filepath"
)

const allocationPoolSize = 100 * 1024 * 1024 // 100Mib
const filePermissions = 0755

// Config represents the application configuration
type Config interface {
	LocalDbConfig
	GossipConfig
	ProducerConfig
	ConsumerPort() int
	AdminPort() int
	HomePath() string
	CreateAllDirs() error
}

type LocalDbConfig interface {
	LocalDbPath() string
}

type ProducerConfig interface {
	ProducerPort() int
	FlowController() FlowController
	MaxMessageSize() int
	// MaxGroupSize is the maximum size of an uncompressed group of messages
	MaxGroupSize() int
}

type GossipConfig interface {
	GossipPort() int
}

func NewConfig() Config {
	return &config{
		flowControl: newFlowControl(allocationPoolSize),
	}
}

type config struct {
	flowControl *flowControl
}

func (c *config) ProducerPort() int {
	return 8081
}

func (c *config) ConsumerPort() int {
	return 8082
}

func (c *config) AdminPort() int {
	return 8083
}

func (c *config) GossipPort() int {
	return 8084
}

func (c *config) MaxMessageSize() int {
	return 1024 * 1024
}

func (c *config) MaxGroupSize() int {
	return 2 * 1024 * 1024
}

func (c *config) FlowController() FlowController {
	return c.flowControl
}

func (c *config) HomePath() string {
	homePath := os.Getenv("SODA_HOME")
	if homePath == "" {
		return filepath.Join("/var", "lib", "soda")
	}
	return homePath
}

func (c *config) dataPath() string {
	return filepath.Join(c.HomePath(), "data")
}

func (c *config) LocalDbPath() string {
	return filepath.Join(c.dataPath(), "local.db")
}

func (c *config) CreateAllDirs() error {
	return os.MkdirAll(c.dataPath(), filePermissions)
}
