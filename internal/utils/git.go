package utils

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

func IsRepoCloned(dirName string) bool {
	if _, err := os.Stat(dirName); os.IsNotExist(err) {
		return false
	}
	return true
}

func PullRepo(r *git.Repository) error {
	// Get the working directory for the repository
	w, err := r.Worktree()
	if err != nil {
		return err
	}
	// Pull the latest changes from the origin remote and merge into the current branch
	err = w.Pull(&git.PullOptions{RemoteName: "origin"})
	if err != nil {
		return err
	}
	// Print the latest commit that was just pulled
	ref, err := r.Head()
	if err != nil {
		return err
	}
	commit, err := r.CommitObject(ref.Hash())
	if err != nil {
		return err
	}
	fmt.Println(commit)
	return nil
}

func CloneRepo(source, branch, dir string) (*git.Repository, error) {
	url := source
	directory := dir
	// Clone the given repository to the given directory
	branchRefName := plumbing.NewBranchReferenceName(branch)
	r, err := git.PlainClone(directory, false, &git.CloneOptions{
		URL:               url,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
		ReferenceName:     plumbing.ReferenceName(branchRefName),
	})
	if err != nil {
		return r, err
	}
	// Print the latest commit that was just pulled
	ref, err := r.Head()
	if err != nil {
		return r, err
	}
	commit, err := r.CommitObject(ref.Hash())
	if err != nil {
		return r, err
	}
	fmt.Println(commit)
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
	// Create the remote with repository URL
	rem, err := r.Remote("origin")
	if err != nil {
		return "", err
	}
	refs, err := rem.List(&git.ListOptions{
		// Returns all references, including peeled references.
		PeelingOption: git.IgnorePeeled,
	})
	if err != nil {
		return "", err
	}
	var refHash string
	for _, ref := range refs {
		if ref.Name().String() == fmt.Sprintf("refs/heads/%s", branch) {
			refHash = ref.Hash().String()
			break
		}
	}
	return refHash, nil
}

func GetLatestLocalCommit(dir string) (string, error) {
	r, err := git.PlainOpen(dir)
	if err != nil {
		return "", err
	}
	// Print the latest commit.
	ref, err := r.Head()
	if err != nil {
		return "", err
	}
	return ref.Hash().String(), nil
}
