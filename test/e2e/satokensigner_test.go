package e2e

import (
	"testing"

	"github.com/stretchr/testify/require"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"

	configv1 "github.com/openshift/api/config/v1"
	configclient "github.com/openshift/client-go/config/clientset/versioned/typed/config/v1"
	"github.com/openshift/cluster-kube-controller-manager-operator/pkg/operator/operatorclient"
	test "github.com/openshift/cluster-kube-controller-manager-operator/test/library"
)

func TestSATokenSignerControllerSyncCerts(t *testing.T) {
	// initialize clients
	kubeConfig, err := test.NewClientConfigForTest()
	require.NoError(t, err)
	kubeClient, err := clientcorev1.NewForConfig(kubeConfig)
	require.NoError(t, err)
	configClient, err := configclient.NewForConfig(kubeConfig)
	require.NoError(t, err)

	// wait for the operator readiness
	test.WaitForKubeControllerManagerClusterOperator(t, configClient, configv1.ConditionTrue, configv1.ConditionFalse, configv1.ConditionFalse)

	err = kubeClient.Secrets(operatorclient.TargetNamespace).Delete("service-account-private-key", &metav1.DeleteOptions{})
	require.NoError(t, err)

	// wait for the operator reporting progressing
	test.WaitForKubeControllerManagerClusterOperator(t, configClient, configv1.ConditionTrue, configv1.ConditionTrue, configv1.ConditionFalse)

	// and check for secret being synced from next-service-private-key
	_, err = kubeClient.Secrets(operatorclient.TargetNamespace).Get("service-account-private-key", metav1.GetOptions{})
	require.NoError(t, err)

	test.WaitForKubeControllerManagerClusterOperator(t, configClient, configv1.ConditionTrue, configv1.ConditionFalse, configv1.ConditionFalse)
}
