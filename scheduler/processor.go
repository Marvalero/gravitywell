package scheduler

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/AlexsJones/gravitywell/configuration"
	"github.com/AlexsJones/gravitywell/platform"
	"github.com/AlexsJones/gravitywell/state"
	"github.com/AlexsJones/gravitywell/vcs"
	"github.com/fatih/color"
)

func process(opt Options, cluster configuration.Cluster) map[state.State]string {
	//---------------------------------
	color.Yellow(fmt.Sprintf("Switching to cluster: %s\n", cluster.Name))
	restclient, k8siface, err := platform.GetKubeClient(cluster.Name)
	if err != nil {
		color.Red(err.Error())
		os.Exit(1)
	}
	//---------------------------------
	stateMap := make(map[state.State]string)
	//---------------------------------
	for _, deployment := range cluster.Deployments {
		//---------------------------------
		color.Yellow(fmt.Sprintf("Fetching deployment %s into %s\n", deployment.Deployment.Name, path.Join(opt.TempVCSPath, deployment.Deployment.Name)))
		gvcs := new(vcs.GitVCS)
		_, err = vcs.Fetch(gvcs, path.Join(opt.TempVCSPath, deployment.Deployment.Name), deployment.Deployment.Git)
		if err != nil {
			color.Cyan("Project already exists in local directory")
		}
		//---------------------------------
		for _, a := range deployment.Deployment.Action {
			if a.Execute.Shell != "" {
				color.Yellow(fmt.Sprintf("Running shell command %s\n", a.Execute.Shell))
				if err := ShellCommand(a.Execute.Shell, path.Join(opt.TempVCSPath, deployment.Deployment.Name), true); err != nil {
					color.Red(err.Error())
				}
			}
			//---------------------------------
			if a.Execute.Kubectl.Command == "" {
				color.Red("No Kubernetes create action to run")
			}
			//---------------------------------
			fileList := []string{}
			err := filepath.Walk(path.Join(opt.TempVCSPath, deployment.Deployment.Name, a.Execute.Kubectl.Path), func(path string, f os.FileInfo, err error) error {
				fileList = append(fileList, path)
				return nil
			})
			if err != nil {
				color.Red(err.Error())

			}
			for _, file := range fileList {
				color.Yellow(fmt.Sprintf("Attempting to deploy %s\n", file))
				if _, err = os.Stat(file); os.IsNotExist(err) {
					continue

				}
				if sa, _ := os.Stat(file); sa.IsDir() {
					continue
				}
				var stateResponse state.State
				if stateResponse, err = platform.DeployFromFile(restclient, k8siface, file, deployment.Deployment.Namespace, opt.DryRun, opt.TryUpdate); err != nil {
					color.Red(err.Error())
				}
				stateMap[stateResponse] = deployment.Deployment.Name
			}
			//---------------------------------
		}
	}
	return stateMap
}
