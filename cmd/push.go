package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"strings"

	"fmt"
	"github.com/cmoulliard/k8s-supervisor/pkg/common/oc"
)

var (
	mode      string
	artefacts = []string{"src", "pom.xml"}
	modes     = []string{"source", "binary"}
)

var pushCmd = &cobra.Command{
	Use:     "push",
	Short:   "Push local code to the development pod",
	Long:    `Push local code to the development pod.`,
	Example: ` sb push`,
	Args:    cobra.RangeArgs(0, 1),
	Run: func(cmd *cobra.Command, args []string) {
		var valid bool
		for _, value := range modes {
			if mode == value {
				valid = true
			}
		}

		if !valid {
			log.WithField("mode", mode).Fatal("The provided mode is not supported: ")
		}

		log.Infof("Push command called with mode '%s'", mode)

		setup, pod := SetupAndWaitForPod()
		podName := pod.Name
		containerName := setup.Application.Name

		log.Info("Copy files from the local developer project to the pod")

		switch mode {
		case "source":
			for i := range artefacts {
				log.Debug("Artefact : ", artefacts[i])
				args := []string{"cp", oc.Client.Pwd + "/" + artefacts[i], podName + ":/tmp/src/", "-c", containerName}
				log.Infof("Copy cmd : %s", args)
				oc.ExecCommand(oc.Command{Args: args})
			}
		case "binary":
			args := []string{"cp", oc.Client.Pwd + "/target/*.jar", podName + ":/deployments", "-c", containerName}
			log.Infof("Copy cmd : %s", args)
			oc.ExecCommand(oc.Command{Args: args})
		}
	},
}

func init() {
	pushCmd.Flags().StringVarP(&mode, "mode", "", "source",
		fmt.Sprintf("Mode used to push the code to the development pod. Supported modes are '%s'", strings.Join(modes, ",")))
	pushCmd.MarkFlagRequired("mode")
	pushCmd.Annotations = map[string]string{"command": "push"}

	rootCmd.AddCommand(pushCmd)
}
