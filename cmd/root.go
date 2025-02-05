package cmd

import (
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/didactiklabs/ginx/internal/utils"
)

var (
	versionFlag      bool
	nowFlag          bool
	exitFailFlag     bool
	version          string
	logLevelFlag     string
	sourceFlag       string
	branchFlag       string
	pollIntervalFlag int
)

var RootCmd = &cobra.Command{
	Use:   "ginx [flags] -- <command>",
	Short: "ginx",
	Long: `
Ginx is a cli tool that watch a remote repository and run an arbitrary command on changes/updates.
`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Initialize configuration here
		initConfig()
	},
	Run: func(cmd *cobra.Command, args []string) {
		if versionFlag {
			fmt.Printf("%s", version)
			os.Exit(0)
		}

		var r *git.Repository
		var err error
		source := sourceFlag
		branch := branchFlag
		interval := time.Duration(pollIntervalFlag) * time.Second
		projectName := path.Base(strings.TrimSuffix(source, "/"))
		dir, err := os.MkdirTemp("", fmt.Sprintf("ginx-%s-*", projectName))
		if err != nil {
			utils.Logger.Fatal("Failed to create temporary directory.", zap.Error(err))
		}

		if !utils.IsRepoCloned(source) {
			utils.Logger.Info("Cloning repository.", zap.String("url", source), zap.String("branch", branch))
			r, err = utils.CloneRepo(source, branch, dir)
			if err != nil {
				utils.Logger.Fatal("Failed to clone repository.", zap.Error(err))
				err := os.RemoveAll(dir)
				if err != nil {
					utils.Logger.Fatal("error removing directory.", zap.Error(err))
				}
			}
		} else {
			r, err = git.PlainOpen(dir)
			utils.Logger.Info("Repository already exist, open directory repository.", zap.String("directory", dir))
			if err != nil {
				utils.Logger.Fatal("Failed to open existing directory repository.", zap.Error(err))
				err := os.RemoveAll(dir)
				if err != nil {
					utils.Logger.Fatal("error removing directory.", zap.Error(err))
				}
			}
		}
		if nowFlag {
			if len(args) > 0 {
				utils.Logger.Info("Running command.", zap.String("command", args[0]), zap.Any("args", args[1:]))
				if err := utils.RunCommand(dir, args[0], args[1:]...); err != nil {
					utils.Logger.Error("Failed to run command.", zap.Error(err))
					err := os.RemoveAll(dir)
					if err != nil {
						utils.Logger.Fatal("error removing directory.", zap.Error(err))
					}
				}
			}
			err := os.RemoveAll(dir)
			if err != nil {
				utils.Logger.Fatal("error removing directory.", zap.Error(err))
			}
			os.Exit(0)
		}

		for {
			// Get the latest commit hash from the remote repository
			remoteCommit, err := utils.GetLatestRemoteCommit(r, branch)
			utils.Logger.Debug("Fetched remote commit.", zap.String("remoteCommit", remoteCommit))
			if err != nil {
				utils.Logger.Fatal("error fetching local commit.", zap.Error(err))
				err := os.RemoveAll(dir)
				if err != nil {
					utils.Logger.Fatal("error removing directory.", zap.Error(err))
				}
			}

			// Get the latest commit hash from the local repository
			localCommit, err := utils.GetLatestLocalCommit(dir)
			utils.Logger.Debug("Fetched local commit.", zap.String("localCommit", localCommit))
			if err != nil {
				utils.Logger.Fatal("error fetching local commit.", zap.Error(err))
				err := os.RemoveAll(dir)
				if err != nil {
					utils.Logger.Fatal("error removing directory.", zap.Error(err))
				}
			}

			if remoteCommit != localCommit {
				utils.Logger.Info("Detected remote changes.", zap.String("url", source), zap.String("branch", branch))
				if err := utils.PullRepo(r); err != nil {
					utils.Logger.Info("Failed to pull. Recloning repository.", zap.String("url", source))
					err := os.RemoveAll(dir)
					if err != nil {
						utils.Logger.Fatal("error removing directory.", zap.Error(err))
					}
					_, err = utils.CloneRepo(source, branch, dir)
					if err != nil {
						utils.Logger.Fatal("Failed to clone repository.", zap.Error(err))
						err := os.RemoveAll(dir)
						if err != nil {
							utils.Logger.Fatal("error removing directory.", zap.Error(err))
						}
					}
				}
				if len(args) > 0 {
					utils.Logger.Info("Running command.", zap.String("command", args[0]), zap.Any("args", args[1:]))
					if err := utils.RunCommand(dir, args[0], args[1:]...); err != nil {
						if exitFailFlag {
							utils.Logger.Fatal("Failed to run command.", zap.Error(err))
							err := os.RemoveAll(dir)
							if err != nil {
								utils.Logger.Fatal("error removing directory.", zap.Error(err))
							}
						}
						utils.Logger.Error("Failed to run command.", zap.Error(err))
						err := os.RemoveAll(dir)
						if err != nil {
							utils.Logger.Fatal("error removing directory.", zap.Error(err))
						}
					}
				}
			} else {
				utils.Logger.Info("No changes detected in remote repository.", zap.String("url", source), zap.String("branch", branch))
			}
			time.Sleep(interval)
		}
	},
	Args: cobra.ArbitraryArgs,
}

func initConfig() {
	// Your configuration initialization logic
	logLevel := zapcore.InfoLevel //nolint:all
	switch logLevelFlag {
	case "debug":
		logLevel = zapcore.DebugLevel
	case "error":
		logLevel = zapcore.ErrorLevel
	default:
		logLevel = zapcore.InfoLevel
	}
	utils.InitializeLogger(logLevel)
}

func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	RootCmd.Flags().BoolVarP(&versionFlag, "version", "v", false, "display version information")
	RootCmd.Flags().BoolVarP(&nowFlag, "now", "", false, "run the command on the targeted branch now")
	RootCmd.Flags().BoolVarP(&exitFailFlag, "exit-on-fail", "", false, "exit on command fail")
	RootCmd.PersistentFlags().StringVarP(&logLevelFlag, "log-level", "l", "info", "override log level (debug, info, error)")
	RootCmd.PersistentFlags().StringVarP(&sourceFlag, "source", "s", "", "git repository to watch")
	RootCmd.PersistentFlags().StringVarP(&branchFlag, "branch", "b", "main", "branch to watch")
	RootCmd.PersistentFlags().IntVarP(&pollIntervalFlag, "interval", "n", 30, "interval in seconds to poll the remote repo")
}
