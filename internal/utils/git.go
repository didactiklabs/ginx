package utils

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"go.uber.org/zap"
)

func IsRepoCloned(dirName string) bool {
	if _, err := os.Stat(dirName); os.IsNotExist(err) {
		return false
	}
	return true
}

func PullRepo(r *git.Repository) error {
	w, err := r.Worktree()
	if err != nil {
		return err
	}
	err = w.Pull(&git.PullOptions{RemoteName: "origin"})
	if err != nil {
		return err
	}
	ref, err := r.Head()
	if err != nil {
		return err
	}
	commit, err := r.CommitObject(ref.Hash())
	if err != nil {
		return err
	}
	Logger.Debug("Pulled latest commit.", zap.String("commit", commit.Hash.String()))
	return nil
}

func CloneRepo(source, branch, dir string) (*git.Repository, error) {
	branchRefName := plumbing.NewBranchReferenceName(branch)
	r, err := git.PlainClone(dir, false, &git.CloneOptions{
		URL:               source,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
		ReferenceName:     plumbing.ReferenceName(branchRefName),
	})
	if err != nil {
		return r, err
	}
	ref, err := r.Head()
	if err != nil {
		return r, err
	}
	commit, err := r.CommitObject(ref.Hash())
	if err != nil {
		return r, err
	}
	Logger.Debug("Cloned repository.", zap.String("commit", commit.Hash.String()))
	return r, nil
}

func RunCommand(dirName, command string, args ...string) error {
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = dirName
	return cmd.Run()
}

func GetLatestRemoteCommit(r *git.Repository, branch string) (string, error) {
	rem, err := r.Remote("origin")
	if err != nil {
		return "", err
	}
	refs, err := rem.List(&git.ListOptions{
		PeelingOption: git.IgnorePeeled,
	})
	if err != nil {
		return "", err
	}
	target := fmt.Sprintf("refs/heads/%s", branch)
	for _, ref := range refs {
		if ref.Name().String() == target {
			return ref.Hash().String(), nil
		}
	}
	return "", nil
}

func GetLatestLocalCommit(dir string) (string, error) {
	r, err := git.PlainOpen(dir)
	if err != nil {
		return "", err
	}
	ref, err := r.Head()
	if err != nil {
		return "", err
	}
	return ref.Hash().String(), nil
}
