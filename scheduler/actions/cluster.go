package actions

import (
	"cloud.google.com/go/container/apiv1"
	"context"
	"errors"
	"github.com/AlexsJones/gravitywell/configuration"
	"github.com/AlexsJones/gravitywell/kinds"
	"github.com/AlexsJones/gravitywell/platform/provider"
	awsprovider "github.com/AlexsJones/gravitywell/platform/provider/aws"
	"github.com/AlexsJones/gravitywell/platform/provider/gcp"
	"github.com/AlexsJones/gravitywell/platform/provider/minikube"
	"github.com/AlexsJones/gravitywell/scheduler/actions/shell"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/fatih/color"
	logger "github.com/sirupsen/logrus"
	"os"
)

func NewMinikubeConfig() (*minikube.MiniKubeProvider, error) {


	return nil, nil
}
func MinikubeClusterProcessor(minikubeprovider *minikube.MiniKubeProvider,
	commandFlag configuration.CommandFlag, cluster kinds.ProviderCluster) error {

	create := func() error {

		err := provider.Create(minikubeprovider, cluster)

		if err != nil {
			return err
		}
		// Run post install -----------------------------------------------------
		for _, executeCommand := range cluster.PostInstallHook {
			if executeCommand.Execute.Shell != "" {
				err := shell.ShellCommand(executeCommand.Execute.Shell,
					executeCommand.Execute.Path, false)
				if err != nil {
					return err
				}

			}
		}
		return nil
	}
	delete := func() error {
		err := provider.Delete(minikubeprovider, cluster)
		if err != nil {
			return err
		}
		// Run post delete -----------------------------------------------------
		for _, executeCommand := range cluster.PostDeleteHook {
			if executeCommand.Execute.Shell != "" {
				err := shell.ShellCommand(executeCommand.Execute.Shell,
					executeCommand.Execute.Path, false)
				if err != nil {
					return err
				}
			}
		}
		return nil
	}
	// Run Command ------------------------------------------------------------------
	switch commandFlag {
	case configuration.Create:
		return create()
	case configuration.Apply:
		logger.Info("Ignoring apply on cluster - no such option")
		return nil
	case configuration.Replace:
		if err := delete(); err != nil {
			return err
		}
		return create()
	case configuration.Delete:
		return delete()
	}
	return nil
}
func NewAmazonWebServicesConfig() (*awsprovider.AWSProvider, error) {
	awsp := awsprovider.AWSProvider{}

	awsP := os.Getenv("AWS_DEFAULT_PROFILE")
	if awsP == "" {
		return nil, errors.New("no AWS_DEFAULT_PROFILE")
	}
	awsR := os.Getenv("AWS_DEFAULT_REGION")
	if awsR == "" {
		return nil, errors.New("no AWS_DEFAULT_REGION")
	}
	awsp.Region = awsR
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(awsR),
		Credentials: credentials.NewSharedCredentials("", awsP),
	})
	if err != nil {
		return nil, err
	}

	awsp.AWSClient = sess

	return &awsp, err
}
func AmazonWebServicesClusterProcessor(awsprovider *awsprovider.AWSProvider,
	commandFlag configuration.CommandFlag,
	cluster kinds.ProviderCluster) error {

	create := func() error {

		err := provider.Create(awsprovider, cluster)

		if err != nil {
			return err
		}
		// Run post install -----------------------------------------------------
		for _, executeCommand := range cluster.PostInstallHook {
			if executeCommand.Execute.Shell != "" {
				err := shell.ShellCommand(executeCommand.Execute.Shell,
					executeCommand.Execute.Path, false)
				if err != nil {
					return err
				}

			}
		}
		return nil
	}
	delete := func() error {
		err := provider.Delete(awsprovider, cluster)
		if err != nil {
			color.Red(err.Error())
		}
		// Run post delete -----------------------------------------------------
		for _, executeCommand := range cluster.PostDeleteHook {
			if executeCommand.Execute.Shell != "" {
				err := shell.ShellCommand(executeCommand.Execute.Shell,
					executeCommand.Execute.Path, false)
				if err != nil {
					return err
				}
			}
		}
		return nil
	}
	// Run Command ------------------------------------------------------------------
	switch commandFlag {
	case configuration.Create:
		return create()
	case configuration.Apply:
		logger.Info("Ignoring apply on cluster - no such option")
		return nil
	case configuration.Replace:
		if err := delete(); err != nil {
			return err
		}
		return create()
	case configuration.Delete:
		return delete()
	}
	return nil
}

func NewGoogleCloudPlatformConfig() (*gcp.GCPProvider, error) {

	gcpProviderClient := &gcp.GCPProvider{}

	gcpProviderClient.Context = context.Background()

	cmc, err := container.NewClusterManagerClient(gcpProviderClient.Context)
	if err != nil {
		logger.Error(err)
		os.Exit(1)
	}
	gcpProviderClient.ClusterManagerClient = cmc

	return gcpProviderClient, nil
}

func GoogleCloudPlatformClusterProcessor(gcpProvider *gcp.GCPProvider,
	commandFlag configuration.CommandFlag,
	cluster kinds.ProviderCluster) error {

	create := func() error {

		err := provider.Create(gcpProvider, cluster)

		if err != nil {
			return err

		}
		// Run post install -----------------------------------------------------
		for _, executeCommand := range cluster.PostInstallHook {
			if executeCommand.Execute.Shell != "" {
				err := shell.ShellCommand(executeCommand.Execute.Shell,
					executeCommand.Execute.Path, false)
				if err != nil {
					return err
				}

			}
		}
		return nil
	}
	delete := func() error {
		err := provider.Delete(gcpProvider, cluster)
		if err != nil {
			return err
		}
		// Run post delete -----------------------------------------------------
		for _, executeCommand := range cluster.PostDeleteHook {
			if executeCommand.Execute.Shell != "" {
				err := shell.ShellCommand(executeCommand.Execute.Shell,
					executeCommand.Execute.Path, false)
				if err != nil {
					return err
				}
			}
		}
		return nil
	}
	// Run Command ------------------------------------------------------------------
	switch commandFlag {
	case configuration.Create:
		return create()
	case configuration.Apply:
		logger.Info("Ignoring apply on cluster - no such option")
		return nil
	case configuration.Replace:
		if err := delete(); err != nil {
			return err
		}
		return create()
	case configuration.Delete:
		return delete()
	}
	return nil
}
