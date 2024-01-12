package deploy

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// Subgraph yaml
type Subgraph struct {
	DataSources []DataSources `yaml:"dataSources"`
	Schema      Schema        `yaml:"schema"`
	SpecVersion string        `yaml:"specVersion"`
}

type DataSources struct {
	Kind    string  `yaml:"kind"`
	Mapping Mapping `yaml:"mapping"`
	Name    string  `yaml:"name"`
	Network string  `yaml:"network"`
	Source  Source  `yaml:"source"`
}

type Schema struct {
	File string `yaml:"file"`
}

type Source struct {
	Abi        string `yaml:"abi"`
	Address    string `yaml:"address"`
	StartBlock int    `yaml:"startBlock"`
}

type Mapping struct {
	Abis          []Abis         `yaml:"abis"`
	ApiVersion    string         `yaml:"apiVersion"`
	Entities      []string       `yaml:"entities"`
	EventHandlers []EventHandler `yaml:"eventHandlers"`
	BlockHandlers []BlockHandler `yaml:"blockHandlers"`
	File          string         `yaml:"file"`
	Kind          string         `yaml:"kind"`
	Language      string         `yaml:"language"`
}

type Abis struct {
	File string `yaml:"file"`
	Name string `yaml:"name"`
}

type EventHandler struct {
	Event   string `yaml:"event"`
	Handler string `yaml:"handler"`
}

type BlockHandler struct {
	Handler string `yaml:"handler"`
}

type Networks map[string]map[string]map[string]any

type subgraphUpdater interface {
	UpdateSubgraph(s *Subgraph, startBlock int)
	UpdateNetworks(n Networks, startBlock int)
}

type zerogDAOperatorStateSubgraphUpdater struct {
	c *Config
}

func (u zerogDAOperatorStateSubgraphUpdater) UpdateSubgraph(s *Subgraph, startBlock int) {
	s.DataSources[0].Source.Address = strings.TrimPrefix(u.c.ZGDA.RegistryCoordinatorWithIndices, "0x")
	s.DataSources[0].Source.StartBlock = startBlock
	s.DataSources[1].Source.Address = strings.TrimPrefix(u.c.ZGDA.PubkeyRegistry, "0x")
	s.DataSources[1].Source.StartBlock = startBlock
	s.DataSources[2].Source.Address = strings.TrimPrefix(u.c.ZGDA.PubkeyCompendium, "0x")
	s.DataSources[2].Source.StartBlock = startBlock
	s.DataSources[3].Source.Address = strings.TrimPrefix(u.c.ZGDA.PubkeyCompendium, "0x")
	s.DataSources[3].Source.StartBlock = startBlock
	s.DataSources[4].Source.Address = strings.TrimPrefix(u.c.ZGDA.RegistryCoordinatorWithIndices, "0x")
	s.DataSources[4].Source.StartBlock = startBlock
	s.DataSources[5].Source.Address = strings.TrimPrefix(u.c.ZGDA.PubkeyRegistry, "0x")
	s.DataSources[5].Source.StartBlock = startBlock
}

func (u zerogDAOperatorStateSubgraphUpdater) UpdateNetworks(n Networks, startBlock int) {
	n["devnet"]["BLSRegistryCoordinatorWithIndices"]["address"] = u.c.ZGDA.RegistryCoordinatorWithIndices
	n["devnet"]["BLSRegistryCoordinatorWithIndices"]["startBlock"] = startBlock
	n["devnet"]["BLSRegistryCoordinatorWithIndices_Operator"]["address"] = u.c.ZGDA.RegistryCoordinatorWithIndices
	n["devnet"]["BLSRegistryCoordinatorWithIndices_Operator"]["startBlock"] = startBlock

	n["devnet"]["BLSPubkeyRegistry"]["address"] = u.c.ZGDA.PubkeyRegistry
	n["devnet"]["BLSPubkeyRegistry"]["startBlock"] = startBlock
	n["devnet"]["BLSPubkeyRegistry_QuorumApkUpdates"]["address"] = u.c.ZGDA.PubkeyRegistry
	n["devnet"]["BLSPubkeyRegistry_QuorumApkUpdates"]["startBlock"] = startBlock

	n["devnet"]["BLSPubkeyCompendium"]["address"] = u.c.ZGDA.PubkeyCompendium
	n["devnet"]["BLSPubkeyCompendium"]["startBlock"] = startBlock
	n["devnet"]["BLSPubkeyCompendium_Operator"]["address"] = u.c.ZGDA.PubkeyCompendium
	n["devnet"]["BLSPubkeyCompendium_Operator"]["startBlock"] = startBlock
}

type zerogDAUIMonitoringUpdater struct {
	c *Config
}

func (u zerogDAUIMonitoringUpdater) UpdateSubgraph(s *Subgraph, startBlock int) {
	s.DataSources[0].Source.Address = strings.TrimPrefix(u.c.ZGDA.ServiceManager, "0x")
	s.DataSources[0].Source.StartBlock = startBlock
}

func (u zerogDAUIMonitoringUpdater) UpdateNetworks(n Networks, startBlock int) {
	n["devnet"]["ZGDAServiceManager"]["address"] = u.c.ZGDA.ServiceManager
	n["devnet"]["ZGDAServiceManager"]["startBlock"] = startBlock
}

func (env *Config) deploySubgraphs(startBlock int) {
	if !env.Environment.IsLocal() {
		return
	}

	currDir, _ := os.Getwd()
	changeDirectory("../subgraphs")
	defer changeDirectory(currDir)

	fmt.Println("Deploying Subgraph")
	env.deploySubgraph(zerogDAOperatorStateSubgraphUpdater{c: env}, "zgda-operator-state", startBlock)
	env.deploySubgraph(zerogDAUIMonitoringUpdater{c: env}, "zgda-batch-metadata", startBlock)
}

func (env *Config) deploySubgraph(updater subgraphUpdater, path string, startBlock int) {
	execBashCmd(fmt.Sprintf(`cp -r "%s" "%s"`, path, env.Path))

	currDir, _ := os.Getwd()
	changeDirectory(filepath.Join(env.Path, path))
	defer changeDirectory(currDir)

	subgraphPath := filepath.Join(env.Path, path)
	env.updateSubgraph(updater, subgraphPath, startBlock)

	execYarnCmd("install")
	execYarnCmd("codegen")
	execYarnCmd("remove-local")
	execYarnCmd("create-local")
	execBashCmd("yarn deploy-local --version-label=v0.0.1")
}

func (env *Config) updateSubgraph(updater subgraphUpdater, path string, startBlock int) {
	const (
		networkFile  = "networks.json"
		subgraphFile = "subgraph.yaml"
	)

	currDir, _ := os.Getwd()
	changeDirectory(path)
	defer changeDirectory(currDir)

	networkData := readFile(networkFile)

	var networkTemplate Networks
	if err := json.Unmarshal([]byte(networkData), &networkTemplate); err != nil {
		log.Panicf("Failed to unmarshal networks.json. Error: %s", err)
	}
	updater.UpdateNetworks(networkTemplate, startBlock)
	networkJson, err := json.MarshalIndent(networkTemplate, "", "  ")
	if err != nil {
		log.Panicf("Error: %s", err.Error())
	}
	writeFile(networkFile, networkJson)
	log.Print("networks.json written")

	subgraphTemplateData := readFile(subgraphFile)

	var sub Subgraph
	if err := yaml.Unmarshal(subgraphTemplateData, &sub); err != nil {
		log.Panicf("Error %s:", err.Error())
	}
	updater.UpdateSubgraph(&sub, startBlock)
	subgraphYaml, err := yaml.Marshal(&sub)
	if err != nil {
		log.Panic(err)
	}
	writeFile(subgraphFile, subgraphYaml)
	log.Print("subgraph.yaml written")
}

func (env *Config) StartGraphNode() {
	changeDirectory(filepath.Join(env.rootPath, "inabox"))
	err := execCmd("./bin.sh", []string{"start-graph"}, []string{})
	if err != nil {
		log.Panicf("Failed to start graph node. Err: %s", err)
	}
}

func (env *Config) StopGraphNode() {
	changeDirectory(filepath.Join(env.rootPath, "inabox"))
	err := execCmd("./bin.sh", []string{"stop-graph"}, []string{})
	if err != nil {
		log.Panicf("Failed to stop graph node. Err: %s", err)
	}
}
