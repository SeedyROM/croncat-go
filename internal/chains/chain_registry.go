package chains

import (
	"encoding/json"
	"errors"
	"os"
	"path"

	"github.com/mattn/go-zglob"

	"github.com/go-git/go-git/v5"
	"github.com/sirupsen/logrus"
)

type ChainRegistry struct {
	path string
}

func NewChainRegistry() (*ChainRegistry, error) {
	// Clone the chain registry locally
	clonePath, err := cloneChainRegistryLocally()
	if err != nil {
		return nil, err
	}

	// Return the chain registry
	return &ChainRegistry{
		path: clonePath,
	}, nil
}

func (c *ChainRegistry) GetChainInfo(id string) (*ChainInfo, error) {
	// Get the chain
	chainInfo, err := getChainInfo(c.path, id)
	if err != nil {
		return nil, err
	}

	// Return the chain
	return chainInfo, nil
}

func cloneChainRegistryLocally() (string, error) {
	// Get the current working directory of the process
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// Clone the given repository to a local hidden directory
	clonePath := path.Join(cwd, ".chain-registry")
	_, err = git.PlainClone(clonePath, false, &git.CloneOptions{
		URL:      "https://github.com/cosmos/chain-registry",
		Progress: os.Stdout,
	})

	// If the directory already exists, just pull the latest changes
	if err != nil && err == git.ErrRepositoryAlreadyExists {
		repo, err := git.PlainOpen(clonePath)
		if err != nil {
			return "", err
		}

		// Get the worktree
		w, err := repo.Worktree()
		if err != nil {
			return "", err
		}

		// Pull the latest changes
		err = w.Pull(&git.PullOptions{RemoteName: "origin"})

		if err != git.NoErrAlreadyUpToDate {
			return "", err
		}
	}

	return clonePath, nil
}

func getChainInfo(pathStr string, id string) (*ChainInfo, error) {
	chainFiles, err := zglob.Glob(path.Join(pathStr, "**/chain.json"))

	if err != nil {
		return nil, err
	}

	for _, chainFile := range chainFiles {
		chainInfo, err := loadChainInfo(chainFile)

		if err != nil {
			logrus.WithError(err).Error("Failed to load chain info for file: ", pathStr)
		}

		if chainInfo.ChainID == id {
			return chainInfo, nil
		}
	}

	return nil, errors.New("chain not found")
}

func loadChainInfo(path string) (*ChainInfo, error) {
	// Read the file
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// Unmarshal the chain info
	chainInfo := &ChainInfo{}
	err = json.Unmarshal(data, chainInfo)

	// Handle invalid JSON
	if err != nil {
		return nil, err
	}

	return chainInfo, nil
}
