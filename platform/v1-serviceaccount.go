package platform

import (
	"errors"
	"fmt"

	"github.com/AlexsJones/gravitywell/configuration"
	"github.com/AlexsJones/gravitywell/state"
	log "github.com/Sirupsen/logrus"
	"k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

func execV1ServiceAccountResouce(k kubernetes.Interface, cm *v1.ServiceAccount, namespace string, opts configuration.Options, commandFlag configuration.CommandFlag) (state.State, error) {
	log.Info("Found ServiceAccount resource")
	cmclient := k.CoreV1().ServiceAccounts(namespace)

	if opts.DryRun {
		_, err := cmclient.Get(cm.Name, v12.GetOptions{})
		if err != nil {
			log.Error(fmt.Sprintf("DRY-RUN: ServiceAccount resource %s does not exist\n", cm.Name))
			return state.EDeploymentStateNotExists, err
		} else {
			log.Info(fmt.Sprintf("DRY-RUN: ServiceAccount resource %s exists\n", cm.Name))
			return state.EDeploymentStateExists, nil
		}
	}
	//Replace -------------------------------------------------------------------
	if commandFlag == configuration.Replace {
		log.Debug("Removing resource in preparation for redeploy")
		graceperiod := int64(0)
		cmclient.Delete(cm.Name, &meta_v1.DeleteOptions{GracePeriodSeconds: &graceperiod})
		_, err := cmclient.Create(cm)
		if err != nil {
			log.Error(fmt.Sprintf("Could not deploy ServiceAccount resource %s due to %s", cm.Name, err.Error()))
			return state.EDeploymentStateError, err
		}
		log.Debug("ServiceAccount deployed")
		return state.EDeploymentStateOkay, nil
	}
	//Create ---------------------------------------------------------------------
	if commandFlag == configuration.Create {
		_, err := cmclient.Create(cm)
		if err != nil {
			log.Error(fmt.Sprintf("Could not deploy ServiceAccount resource %s due to %s", cm.Name, err.Error()))
			return state.EDeploymentStateError, err
		}
		log.Debug("ServiceAccount deployed")
		return state.EDeploymentStateOkay, nil
	}
	//Apply --------------------------------------------------------------------
	if commandFlag == configuration.Apply {
		_, err := cmclient.Update(cm)
		if err != nil {
			log.Error("Could not update ServiceAccount")
			return state.EDeploymentStateCantUpdate, err
		}
		log.Debug("ServiceAccount updated")
		return state.EDeploymentStateUpdated, nil
	}
	//Delete -------------------------------------------------------------------
	if commandFlag == configuration.Delete {
		err := cmclient.Delete(cm.Name, &meta_v1.DeleteOptions{})
		if err != nil {
			log.Error("Could not delete Service Account")
			return state.EDeploymentStateCantUpdate, err
		}
		log.Debug("Service Account deleted")
		return state.EDeploymentStateOkay, nil
	}
	return state.EDeploymentStateNil, errors.New("No kubectl command")

}