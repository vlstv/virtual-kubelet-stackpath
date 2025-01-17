package provider

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/go-openapi/strfmt"
	gomock "github.com/golang/mock/gomock"
	_ "github.com/golang/mock/mockgen/model"
	"github.com/google/uuid"
	"github.com/stackpath/vk-stackpath-provider/internal/api/workload/workload_client"
	"github.com/stackpath/vk-stackpath-provider/internal/api/workload/workload_client/instance"
	workloads "github.com/stackpath/vk-stackpath-provider/internal/api/workload/workload_client/workloads"
	"github.com/stackpath/vk-stackpath-provider/internal/api/workload/workload_models"
	"github.com/stackpath/vk-stackpath-provider/internal/config"
	mocks "github.com/stackpath/vk-stackpath-provider/internal/mocks"
	"github.com/virtual-kubelet/virtual-kubelet/node/nodeutil"
	"gotest.tools/assert"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

const (
	mockedNodeName   = "vk-mock"
	testStackID      = "test"
	testClientID     = "test"
	testClientSecret = "test"
	testCityCode     = "JFK"
)

func TestCreatePod(t *testing.T) {
	podName := fmt.Sprintf("test-pod-%s", uuid.New().String())
	podNamespace := fmt.Sprintf("test-ns-%s", uuid.New().String())
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	ctx := context.Background()

	wsc := mocks.NewWorkloadsClientService(mockController)
	stackPathClientMock := workload_client.EdgeCompute{Workloads: wsc}

	provider, err := createTestProvider(ctx, mocks.NewMockConfigMapLister(mockController), mocks.NewMockSecretLister(mockController), mocks.NewMockPodLister(mockController), &stackPathClientMock)
	if err != nil {
		t.Fatal("failed to create the test provider", err)
	}

	testPod := createTestPod(podName, podNamespace)
	badPod := testPod
	badPod.Spec.Containers[0].LivenessProbe = &v1.Probe{
		ProbeHandler: v1.ProbeHandler{
			TCPSocket: &v1.TCPSocketAction{
				Port: intstr.IntOrString{Type: intstr.String, StrVal: "http"},
			},
		},
	}

	testCases := []struct {
		description     string
		initMockedCalls func()
		pod             *v1.Pod
		expectedError   error
	}{
		{
			description: "successfully creates a workload",
			pod:         testPod,
			initMockedCalls: func() {
				w, _ := provider.getWorkloadFrom(testPod)
				params := workloads.CreateWorkloadParams{
					Body:    &workload_models.V1CreateWorkloadRequest{Workload: w},
					StackID: provider.apiConfig.StackID,
					Context: ctx,
				}
				wsc.EXPECT().CreateWorkload(&params, nil).Times(1)
			},
			expectedError: nil,
		},
		{
			description: "fails to create a workload due to bad probe port",
			pod:         badPod,
			initMockedCalls: func() {
				w, _ := provider.getWorkloadFrom(testPod)
				params := workloads.CreateWorkloadParams{
					Body:    &workload_models.V1CreateWorkloadRequest{Workload: w},
					StackID: provider.apiConfig.StackID,
					Context: ctx,
				}
				wsc.EXPECT().CreateWorkload(&params, nil).Return(nil, errors.New("unable to find named port")).Times(1)
			},
			expectedError: errors.New("unable to find named port"),
		},
	}

	for _, c := range testCases {
		t.Run(c.description, func(t *testing.T) {
			c.initMockedCalls()
			err := provider.CreatePod(context.Background(), c.pod)
			if err != nil {
				assert.Equal(t, c.expectedError.Error(), err.Error())
			} else {
				assert.Equal(t, c.expectedError, nil)
			}
		})
	}
}

func TestDeletePod(t *testing.T) {
	podName := fmt.Sprintf("test-pod-%s", uuid.New().String())
	podNamespace := fmt.Sprintf("test-ns-%s", uuid.New().String())
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	ctx := context.Background()

	wsc := mocks.NewWorkloadsClientService(mockController)
	stackPathClientMock := workload_client.EdgeCompute{Workloads: wsc}
	testPod := createTestPod(podName, podNamespace)

	provider, err := createTestProvider(ctx, mocks.NewMockConfigMapLister(mockController), mocks.NewMockSecretLister(mockController), mocks.NewMockPodLister(mockController), &stackPathClientMock)
	if err != nil {
		t.Fatal("failed to create the test provider", err)
	}

	params := workloads.DeleteWorkloadParams{
		StackID:    provider.apiConfig.StackID,
		WorkloadID: provider.getWorkloadSlug(podNamespace, podName),
		Context:    ctx,
	}

	testCases := []struct {
		description     string
		initMockedCalls func()
		expectedError   error
	}{
		{
			description: "successfully deletes a pod",
			initMockedCalls: func() {
				wsc.EXPECT().DeleteWorkload(&params, nil).Return(nil, nil).Times(1)
			},
			expectedError: nil,
		},
		{
			description: "fails to delete a pod",
			initMockedCalls: func() {
				wsc.EXPECT().DeleteWorkload(&params, nil).Return(nil, errors.New("API call failed")).Times(1)
			},
			expectedError: errors.New("API call failed"),
		},
	}
	for _, c := range testCases {
		t.Run(c.description, func(t *testing.T) {
			c.initMockedCalls()
			err := provider.DeletePod(context.Background(), testPod)
			if err != nil {
				assert.Equal(t, c.expectedError.Error(), err.Error())
			} else {
				assert.Equal(t, c.expectedError, nil)
			}
		})
	}
}

func TestGetPodStatus(t *testing.T) {
	podName := "test-pod"
	podNamespace := "test-ns"
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	ctx := context.Background()

	isc := mocks.NewInstanceClientService(mockController)
	stackPathClientMock := workload_client.EdgeCompute{Instance: isc}

	provider, err := createTestProvider(ctx, mocks.NewMockConfigMapLister(mockController), mocks.NewMockSecretLister(mockController), mocks.NewMockPodLister(mockController), &stackPathClientMock)
	if err != nil {
		t.Fatal("failed to create the test provider", err)
	}

	params := instance.GetWorkloadInstanceParams{
		Context:      ctx,
		StackID:      provider.apiConfig.StackID,
		WorkloadID:   provider.getWorkloadSlug(podNamespace, podName),
		InstanceName: provider.getInstanceName(podNamespace, podName),
	}

	phases := map[string]v1.PodPhase{
		string(workload_models.Workloadv1InstanceInstancePhaseSCHEDULING): v1.PodPending,
		string(workload_models.Workloadv1InstanceInstancePhaseSTARTING):   v1.PodPending,
		string(workload_models.Workloadv1InstanceInstancePhaseRUNNING):    v1.PodRunning,
		string(workload_models.Workloadv1InstanceInstancePhaseCOMPLETED):  v1.PodSucceeded,
		string(workload_models.Workloadv1InstanceInstancePhaseSTOPPED):    v1.PodSucceeded,
		string(workload_models.Workloadv1InstanceInstancePhaseFAILED):     v1.PodFailed,
	}

	for workloadPhase, podPhase := range phases {
		i := createTestInstance(
			provider.getInstanceName(podNamespace, podName),
			workload_models.NewWorkloadv1InstanceInstancePhase(workload_models.Workloadv1InstanceInstancePhase(workloadPhase)),
			&workload_models.V1ContainerStatus{
				Waiting: &workload_models.ContainerStatusWaiting{
					Reason:  "waiting",
					Message: "waiting message"},
			},
		)

		isc.EXPECT().GetWorkloadInstance(&params, nil).Return(
			&instance.GetWorkloadInstanceOK{
				Payload: &workload_models.V1GetWorkloadInstanceResponse{
					Instance: i,
				},
			},
			nil,
		)

		podStatus, err := provider.GetPodStatus(ctx, podNamespace, podName)
		if err != nil {
			t.Fatal("Failed to get pod status", err)
		}

		assert.Equal(t, podStatus.Phase, podPhase)
	}
}

func TestGetPod(t *testing.T) {
	podName := "test-pod"
	podNamespace := "test-ns"
	nodeName := "test-node"
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	ctx := context.Background()

	wsc := mocks.NewWorkloadsClientService(mockController)
	isc := mocks.NewInstanceClientService(mockController)
	activePodsLister := mocks.NewMockPodLister(mockController)
	mockPodsNamespaceLister := mocks.NewMockPodNamespaceLister(mockController)
	stackPathClientMock := workload_client.EdgeCompute{Workloads: wsc, Instance: isc}

	pod := createTestPod(podName, podNamespace)
	pod.Status.Phase = v1.PodPending

	provider, err := createTestProvider(ctx, mocks.NewMockConfigMapLister(mockController), mocks.NewMockSecretLister(mockController), activePodsLister, &stackPathClientMock)
	if err != nil {
		t.Fatal("failed to create the test provider", err)
	}

	testCases := []struct {
		description      string
		initMockedCalls  func()
		updatedPodStatus v1.PodPhase
		expectedError    error
	}{
		{
			description: "successfully gets a pod and updates its status accordingly",
			initMockedCalls: func() {
				wsc.EXPECT().GetWorkload(gomock.Any(), gomock.Any()).Return(&workloads.GetWorkloadOK{
					Payload: &workload_models.V1GetWorkloadResponse{
						Workload: &workload_models.V1Workload{
							Name: podName,
							Metadata: &workload_models.V1Metadata{
								Labels: workload_models.V1StringMapEntry{
									nodeNameLabelKey:     nodeName,
									podNamespaceLabelKey: podNamespace,
									podNameLabelKey:      podName,
								},
							},
						},
					},
				}, nil).Times(1)

				isc.EXPECT().GetWorkloadInstance(gomock.Any(), gomock.Any()).Return(&instance.GetWorkloadInstanceOK{
					Payload: &workload_models.V1GetWorkloadInstanceResponse{
						Instance: &workload_models.Workloadv1Instance{
							Name:  podName,
							Phase: workload_models.Workloadv1InstanceInstancePhaseRUNNING.Pointer(),
						},
					},
				}, nil).Times(1)

				activePodsLister.EXPECT().Pods(podNamespace).Return(mockPodsNamespaceLister).Times(1)
				mockPodsNamespaceLister.EXPECT().Get(podName).Return(pod, nil).Times(1)
			},
			expectedError:    nil,
			updatedPodStatus: v1.PodRunning,
		},
		{
			description: "fails to get a pod due to workload API failure",
			initMockedCalls: func() {
				wsc.EXPECT().GetWorkload(gomock.Any(), gomock.Any()).Return(nil, errors.New("API error")).Times(1)
			},
			expectedError:    errors.New("API error"),
			updatedPodStatus: v1.PodUnknown,
		},
		{
			description: "fails to get a pod due to instance API failure",
			initMockedCalls: func() {
				wsc.EXPECT().GetWorkload(gomock.Any(), gomock.Any()).Return(&workloads.GetWorkloadOK{
					Payload: &workload_models.V1GetWorkloadResponse{
						Workload: &workload_models.V1Workload{
							Name: podName,
							Metadata: &workload_models.V1Metadata{
								Labels: workload_models.V1StringMapEntry{
									nodeNameLabelKey:     nodeName,
									podNamespaceLabelKey: podNamespace,
									podNameLabelKey:      podName,
								},
							},
						},
					},
				}, nil).Times(1)

				isc.EXPECT().GetWorkloadInstance(gomock.Any(), gomock.Any()).Return(nil, errors.New("API error")).Times(1)
			},
			expectedError:    errors.New("API error"),
			updatedPodStatus: v1.PodUnknown,
		},
		{
			description: "fails to get a pod due to an error occurred while retrieving the pod form the indexer",
			initMockedCalls: func() {
				wsc.EXPECT().GetWorkload(gomock.Any(), gomock.Any()).Return(&workloads.GetWorkloadOK{
					Payload: &workload_models.V1GetWorkloadResponse{
						Workload: &workload_models.V1Workload{
							Name: podName,
							Metadata: &workload_models.V1Metadata{
								Labels: workload_models.V1StringMapEntry{
									nodeNameLabelKey:     nodeName,
									podNamespaceLabelKey: podNamespace,
									podNameLabelKey:      podName,
								},
							},
						},
					},
				}, nil).Times(1)

				isc.EXPECT().GetWorkloadInstance(gomock.Any(), gomock.Any()).Return(&instance.GetWorkloadInstanceOK{
					Payload: &workload_models.V1GetWorkloadInstanceResponse{
						Instance: &workload_models.Workloadv1Instance{
							Name:  podName,
							Phase: workload_models.Workloadv1InstanceInstancePhaseRUNNING.Pointer(),
						},
					},
				}, nil).Times(1)

				activePodsLister.EXPECT().Pods(podNamespace).Return(mockPodsNamespaceLister).Times(1)
				mockPodsNamespaceLister.EXPECT().Get(podName).Return(nil, errors.New("indexer error")).Times(1)
			},
			expectedError:    errors.New("indexer error"),
			updatedPodStatus: v1.PodUnknown,
		},
	}

	for _, c := range testCases {
		t.Run(c.description, func(t *testing.T) {
			c.initMockedCalls()

			updatedPod, err := provider.GetPod(ctx, podNamespace, podName)
			if err != nil {
				assert.Equal(t, c.expectedError.Error(), err.Error())
			} else {
				assert.Equal(t, c.expectedError, nil)
				assert.Equal(t, pod.Status.Phase, v1.PodPending)
				assert.Equal(t, updatedPod.Name, podName)
				assert.Equal(t, updatedPod.Namespace, podNamespace)
				assert.Equal(t, updatedPod.Status.Phase, c.updatedPodStatus)
			}
		})
	}
}

func createTestProvider(ctx context.Context, configMapMocker *mocks.MockConfigMapLister, secretMocker *mocks.MockSecretLister, podMocker *mocks.MockPodLister, stackpathClient *workload_client.EdgeCompute) (*StackpathProvider, error) {

	cfg := nodeutil.ProviderConfig{
		ConfigMaps: configMapMocker,
		Secrets:    secretMocker,
		Pods:       podMocker,
	}

	cfg.Node = &v1.Node{}

	cfg.Node.Name = mockedNodeName

	err := os.Setenv("SP_STACK_ID", testStackID)
	if err != nil {
		return nil, err
	}
	err = os.Setenv("SP_CLIENT_ID", testClientID)
	if err != nil {
		return nil, err
	}
	err = os.Setenv("SP_CLIENT_SECRET", testClientSecret)
	if err != nil {
		return nil, err
	}
	err = os.Setenv("SP_CITY_CODE", testCityCode)
	if err != nil {
		return nil, err
	}

	apiConfig, err := config.NewConfig(ctx)
	if err != nil {
		return nil, err
	}

	provider, err := NewStackpathProvider(ctx, stackpathClient, apiConfig, cfg, "127.0.0.1")

	if err != nil {
		return nil, err
	}

	return provider, nil
}

func createContainerState(name string, containerState v1.ContainerState) *workload_models.V1ContainerStatus {
	cs := workload_models.V1ContainerStatus{Name: name}
	if containerState.Running != nil {
		cs.Running = &workload_models.ContainerStatusRunning{StartedAt: strfmt.DateTime(time.Now())}
	} else if containerState.Waiting != nil {
		cs.Waiting = &workload_models.ContainerStatusWaiting{}
	} else {
		// terminated
		cs.Terminated = &workload_models.ContainerStatusTerminated{
			ExitCode:   1,
			Message:    "A message",
			Reason:     "A reason",
			FinishedAt: strfmt.DateTime(time.Now()),
			StartedAt:  strfmt.DateTime(time.Now().AddDate(0, 0, -1)),
		}
	}

	return &cs
}

func createTestContainerSpec() workload_models.V1ContainerSpec {
	return workload_models.V1ContainerSpec{
		Image: "nginx:latest",
	}
}

func createTestInstance(name string, phase *workload_models.Workloadv1InstanceInstancePhase, containerStatus *workload_models.V1ContainerStatus) *workload_models.Workloadv1Instance {

	spec := createTestContainerSpec()

	i := workload_models.Workloadv1Instance{
		ContainerStatuses: []*workload_models.V1ContainerStatus{containerStatus},
		Containers:        workload_models.V1ContainerSpecMapEntry{name: spec},
		ID:                uuid.New().String(),
		Metadata:          nil,
		IPAddress:         "127.0.0.1",
		IPV6Address:       "::1",
		Name:              name,
		Phase:             phase,
	}
	return &i
}

func createTestPod(podName, podNamespace string) *v1.Pod {
	return &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      podName,
			Namespace: podNamespace,
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name: "nginx",
					Ports: []v1.ContainerPort{
						{
							Name:          "http",
							ContainerPort: 8080,
						},
					},
					Resources: v1.ResourceRequirements{
						Requests: v1.ResourceList{
							"cpu":    resource.MustParse("0.99"),
							"memory": resource.MustParse("1.5G"),
						},
						Limits: v1.ResourceList{
							"cpu":    resource.MustParse("3999m"),
							"memory": resource.MustParse("8010M"),
						},
					},

					LivenessProbe: &v1.Probe{
						ProbeHandler: v1.ProbeHandler{
							HTTPGet: &v1.HTTPGetAction{
								Port: intstr.FromString("http"),
								Path: "/",
							},
						},
						InitialDelaySeconds: 10,
						PeriodSeconds:       5,
						TimeoutSeconds:      60,
						SuccessThreshold:    3,
						FailureThreshold:    5,
					},
					ReadinessProbe: &v1.Probe{
						ProbeHandler: v1.ProbeHandler{
							HTTPGet: &v1.HTTPGetAction{
								Port: intstr.FromInt(8080),
								Path: "/",
							},
						},
						InitialDelaySeconds: 10,
						PeriodSeconds:       5,
						TimeoutSeconds:      60,
						SuccessThreshold:    3,
						FailureThreshold:    5,
					},
				},
			},
		},
	}
}
