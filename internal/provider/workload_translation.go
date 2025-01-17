// Package provider implements the stackpath virtual kubelet provider
package provider

import (
	"fmt"
	"strings"

	"github.com/stackpath/vk-stackpath-provider/internal/api/workload/workload_models"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/util/intstr"
)

const oneGi = 1024 * 1024 * 1024

func (p *StackpathProvider) getWorkloadFrom(pod *v1.Pod) (*workload_models.V1Workload, error) {
	spec, err := p.getWorkloadSpecFrom(pod)
	if err != nil {
		return nil, err
	}

	targets := p.getWorkloadTargetsFrom(pod)

	metadata := p.getWorkloadMetadataFrom(pod)

	w := workload_models.V1Workload{
		Name:     p.getWorkloadSlug(pod.Namespace, pod.Name),
		Slug:     p.getWorkloadSlug(pod.Namespace, pod.Name),
		Metadata: metadata,
		Spec:     spec,
		Targets:  targets,
	}
	return &w, nil
}

func (p *StackpathProvider) getWorkloadMetadataFrom(pod *v1.Pod) *workload_models.V1Metadata {
	metadata := workload_models.V1Metadata{
		Labels: workload_models.V1StringMapEntry{
			podNameLabelKey:      pod.Name,
			podNamespaceLabelKey: pod.Namespace,
			nodeNameLabelKey:     p.nodeName,
		},
	}
	return &metadata
}

func (p *StackpathProvider) getWorkloadSpecFrom(pod *v1.Pod) (*workload_models.V1WorkloadSpec, error) {
	containers, err := p.getWorkloadContainersFrom(pod.Spec.Containers)
	if err != nil {
		return nil, err
	}

	networkInterfaces := p.getWorkloadNetworkInterfacesFrom(pod)

	volumes, err := p.getWorkloadVolumesFrom(pod)
	if err != nil {
		return nil, err
	}

	imagePullCredentials, err := p.getImagePullCredentialsFrom(pod.Namespace, pod.Spec.ImagePullSecrets)
	if err != nil {
		return nil, err
	}

	spec := workload_models.V1WorkloadSpec{
		Containers:           containers,
		NetworkInterfaces:    networkInterfaces,
		ImagePullCredentials: imagePullCredentials,
		VolumeClaimTemplates: volumes,
	}
	return &spec, nil
}

func (p *StackpathProvider) getVolumeClaimSpecFrom(spec *v1.CSIVolumeSource) (*workload_models.V1VolumeClaimSpec, error) {
	var err error
	volumeSize := *resource.NewQuantity(1*oneGi, resource.BinarySI)
	if volumeAttributeSize, ok := spec.VolumeAttributes["size"]; ok {
		volumeSize, err = resource.ParseQuantity(volumeAttributeSize)
		if err != nil {
			return nil, err
		}
	}

	// force size to match Stackpath limit
	if providedSize, ok := volumeSize.AsInt64(); ok {
		if providedSize > int64(1000*oneGi) {
			p.logger.Info("adjusting volume size to match Stackpath 1000Gi limit")
			volumeSize = *resource.NewQuantity(1000*oneGi, resource.BinarySI)

		}
	}

	storage := volumeSize.String()

	return &workload_models.V1VolumeClaimSpec{
		Resources: &workload_models.V1ResourceRequirements{
			Limits:   workload_models.V1StringMapEntry{"storage": storage},
			Requests: workload_models.V1StringMapEntry{"storage": storage},
		},
	}, nil
}

func (p *StackpathProvider) getWorkloadVolumesFrom(pod *v1.Pod) ([]*workload_models.V1VolumeClaim, error) {
	workloadVolumeClaim := []*workload_models.V1VolumeClaim{}
	for _, volume := range pod.Spec.Volumes {
		if volume.CSI == nil || volume.CSI.Driver != stackpathVirtualKubeletCSIDriver {
			p.logger.Infof("skipping volume %s, only CSI driver of type %s volumes are supported", volume.Name, stackpathVirtualKubeletCSIDriver)
			continue
		}

		spec, err := p.getVolumeClaimSpecFrom(volume.CSI)
		if err != nil {
			return nil, err
		}
		workloadVolumeClaim = append(workloadVolumeClaim, &workload_models.V1VolumeClaim{
			Name:     volume.Name,
			Slug:     volume.Name,
			Metadata: &workload_models.V1Metadata{},
			Spec:     spec,
		})
	}
	return workloadVolumeClaim, nil
}

func (p *StackpathProvider) getWorkloadContainersFrom(k8sContainers []v1.Container) (workload_models.V1ContainerSpecMapEntry, error) {
	containers := make(workload_models.V1ContainerSpecMapEntry)
	for _, k8sContainer := range k8sContainers {
		container, err := p.getWorkloadContainerSpecFrom(&k8sContainer)
		if err != nil {
			return nil, err
		}
		containers[k8sContainer.Name] = *container
	}
	return containers, nil
}

func (p *StackpathProvider) getWorkloadTargetsFrom(pod *v1.Pod) workload_models.V1TargetMapEntry {
	// creating one target just for provided city code
	// No autoscaling allowed, min and max replicas = 1
	target := workload_models.V1Target{
		Spec: &workload_models.V1TargetSpec{
			DeploymentScope: "cityCode",
			Deployments: &workload_models.V1DeploymentSpec{
				MinReplicas: 1,
				MaxReplicas: 1,
				Selectors:   []*workload_models.V1MatchExpression{{Key: "cityCode", Operator: "in", Values: []string{p.apiConfig.CityCode}}},
			},
		},
	}
	targets := workload_models.V1TargetMapEntry{targetName: target}
	return targets
}

func (p *StackpathProvider) getWorkloadNetworkInterfacesFrom(pod *v1.Pod) []*workload_models.V1NetworkInterface {
	networkInterface := workload_models.V1NetworkInterface{
		Network: "default",
		IPFamilies: []*workload_models.V1IPFamily{
			workload_models.NewV1IPFamily(workload_models.V1IPFamilyIPV4),
		},
		Subnet:            "",
		IPV6Subnet:        "",
		EnableOneToOneNat: true,
	}
	networkInterfaces := []*workload_models.V1NetworkInterface{&networkInterface}
	return networkInterfaces

}

func (p *StackpathProvider) getWorkloadContainerSpecFrom(k8sContainer *v1.Container) (*workload_models.V1ContainerSpec, error) {
	ports := p.getWorkloadContainerPortsFrom(k8sContainer.Ports)

	env := p.getWorkloadContainerEnvFrom(k8sContainer.Env)

	resources := p.getWorkloadContainerResourcesFrom(k8sContainer.Resources)

	volumeMounts := p.getWorkloadContainerVolumeMountsFrom(k8sContainer.VolumeMounts)

	livenessProbe, err := p.getWorkloadContainerProbeFrom(k8sContainer.LivenessProbe, k8sContainer.Ports)
	if err != nil {
		return nil, err
	}

	readinessProbe, err := p.getWorkloadContainerProbeFrom(k8sContainer.ReadinessProbe, k8sContainer.Ports)
	if err != nil {
		return nil, err
	}

	k8sCommand := k8sContainer.Command
	k8sArgs := k8sContainer.Args
	if len(k8sArgs) != 0 {
		if len(k8sCommand) == 0 {
			return nil, fmt.Errorf("failed to create workload from pod. args not supported without command")
		}
		k8sCommand = append(k8sCommand, k8sArgs...)
	}

	workloadContainerSpec := workload_models.V1ContainerSpec{
		Image:          k8sContainer.Image,
		Command:        k8sCommand,
		Ports:          ports,
		Env:            env,
		Resources:      resources,
		VolumeMounts:   volumeMounts,
		LivenessProbe:  livenessProbe,
		ReadinessProbe: readinessProbe,
	}

	return &workloadContainerSpec, nil
}

func (p *StackpathProvider) getWorkloadContainerPortsFrom(k8sPorts []v1.ContainerPort) workload_models.V1InstancePortMapEntry {
	portsToReturn := workload_models.V1InstancePortMapEntry{}
	for _, k8sPort := range k8sPorts {
		name := k8sPort.Name
		if name == "" {
			// provide a default name if not provided
			name = "default"
		}

		// Default protocol is TCP
		protocol := k8sPort.Protocol
		if protocol != v1.ProtocolTCP && protocol != v1.ProtocolUDP {
			protocol = v1.ProtocolTCP
		}

		portsToReturn[name] = workload_models.V1InstancePort{
			Port:     k8sPort.ContainerPort,
			Protocol: strings.ToUpper(string(protocol)),
		}
	}

	return portsToReturn
}

func (p *StackpathProvider) getWorkloadContainerEnvFrom(k8sEnv []v1.EnvVar) workload_models.V1EnvironmentVariableMapEntry {
	envToReturn := workload_models.V1EnvironmentVariableMapEntry{}
	for _, k8sEnvVar := range k8sEnv {
		if k8sEnvVar.ValueFrom == nil {
			ignore := false
			// k8s adds default environment variables that aren't useful in the stackpath environment, so don't use them
			for _, ignoreVar := range k8sEnvVariablesToIgnore {
				if k8sEnvVar.Name == ignoreVar {
					ignore = true
				}
			}
			if ignore {
				continue
			}
			envToReturn[k8sEnvVar.Name] = workload_models.V1EnvironmentVariable{
				Value: k8sEnvVar.Value,
			}
		} else {
			// TODO Handle environment variable from source
			p.logger.Warnf("Value From is not supported for env var %s", k8sEnvVar.Name)
		}
	}

	return envToReturn
}

func (p *StackpathProvider) getWorkloadContainerResourcesFrom(k8sResource v1.ResourceRequirements) *workload_models.V1ResourceRequirements {

	// set default
	resources := workload_models.V1ResourceRequirements{
		Requests: containerResourcesSP1,
		Limits:   containerResourcesSP1,
	}

	resources.Limits = toSPInstanceSize(k8sResource)
	resources.Requests = toSPInstanceSize(k8sResource)
	return &resources

}

// toSPInstanceSize replaces the resource allocation requests to one of the 5 supported instance sizes provided by SP
func toSPInstanceSize(k8sResource v1.ResourceRequirements) workload_models.V1StringMapEntry {
	var requestCPU, requestMEM, limitCPU, limitMEM *resource.Quantity

	mem2Gi := resource.NewQuantity(2*oneGi, resource.BinarySI)
	mem4Gi := resource.NewQuantity(4*oneGi, resource.BinarySI)
	mem8Gi := resource.NewQuantity(8*oneGi, resource.BinarySI)
	mem16Gi := resource.NewQuantity(16*oneGi, resource.BinarySI)

	cpu1 := resource.NewQuantity(1, resource.DecimalSI)
	cpu2 := resource.NewQuantity(2, resource.DecimalSI)
	cpu4 := resource.NewQuantity(4, resource.DecimalSI)

	if k8sResource.Requests != nil {
		requestCPU = k8sResource.Requests.Cpu()
		requestMEM = k8sResource.Requests.Memory()
	} else {
		// default values matching SP1
		requestCPU = resource.NewQuantity(1, resource.DecimalSI)
		requestMEM = mem2Gi
	}

	if k8sResource.Limits != nil {
		limitCPU = k8sResource.Limits.Cpu()
		limitMEM = k8sResource.Limits.Memory()
	} else {
		// default values matching SP1
		limitCPU = resource.NewQuantity(1, resource.DecimalSI)
		limitMEM = mem2Gi
	}

	cpu := maxResource(requestCPU, limitCPU)
	mem := maxResource(requestMEM, limitMEM)

	if cpu.Cmp(*cpu4) == 1 || mem.Cmp(*mem16Gi) == 1 {
		return containerResourcesSP5
	}

	if (cpu.Cmp(*cpu2) == 1 && cpu.Cmp(*cpu4) != 1) || (mem.Cmp(*mem8Gi) == 1 && mem.Cmp(*mem16Gi) != 1) {
		return containerResourcesSP4
	}

	if cpu.Cmp(*cpu1) == 1 && cpu.Cmp(*cpu2) != 1 {
		if mem.Cmp(*mem4Gi) == 1 {
			return containerResourcesSP3
		}
		return containerResourcesSP2
	}

	if cpu.Cmp(*cpu1) != 1 {
		if mem.Cmp(*mem4Gi) == 1 && mem.Cmp(*mem8Gi) != 1 {
			return containerResourcesSP3
		}
		if mem.Cmp(*mem2Gi) == 1 && mem.Cmp(*mem4Gi) != 1 {
			return containerResourcesSP2
		}
	}

	return containerResourcesSP1
}

func maxResource(x, y *resource.Quantity) *resource.Quantity {
	if x.Cmp(*y) == -1 {
		return y
	}
	return x
}

func (p *StackpathProvider) getWorkloadContainerVolumeMountsFrom(k8sVolumeMounts []v1.VolumeMount) []*workload_models.V1InstanceVolumeMount {
	workloadVolumeMounts := []*workload_models.V1InstanceVolumeMount{}
	for _, k8sVolumeMount := range k8sVolumeMounts {
		if k8sVolumeMount.MountPath == defaultK8sServiceAccountMountPath {
			// skip mounting the default service account created by k8s
			continue
		}
		workloadVolumeMounts = append(workloadVolumeMounts, &workload_models.V1InstanceVolumeMount{
			Slug:      k8sVolumeMount.Name,
			MountPath: k8sVolumeMount.MountPath,
		})
	}
	return workloadVolumeMounts
}

func (p *StackpathProvider) getProbeHTTPHeadersFrom(httpHeaders []v1.HTTPHeader) workload_models.V1StringMapEntry {
	httpHeadersToReturn := workload_models.V1StringMapEntry{}
	for _, header := range httpHeaders {
		httpHeadersToReturn[header.Name] = header.Value
	}
	return httpHeadersToReturn
}

func (p *StackpathProvider) getWorkloadContainerProbeFrom(k8sProbe *v1.Probe, containerPorts []v1.ContainerPort) (*workload_models.V1Probe, error) {

	if k8sProbe == nil || *k8sProbe == (v1.Probe{}) {
		return nil, nil
	}

	if k8sProbe.GRPC != nil {
		p.logger.Infof("probe of type GRPC is not supported, skipping")
		return nil, nil
	}

	if k8sProbe.Exec != nil {
		p.logger.Info("probe of type Exec is not supported, skipping")
		return nil, nil
	}

	probe := &workload_models.V1Probe{
		FailureThreshold:    k8sProbe.FailureThreshold,
		InitialDelaySeconds: k8sProbe.InitialDelaySeconds,
		PeriodSeconds:       k8sProbe.PeriodSeconds,
		SuccessThreshold:    k8sProbe.SuccessThreshold,
		TimeoutSeconds:      k8sProbe.TimeoutSeconds,
	}

	if k8sProbe.HTTPGet != nil {
		port, err := getPortValue(k8sProbe.HTTPGet.Port, containerPorts)
		if err != nil {
			return nil, err
		}
		probe.HTTPGet = &workload_models.V1HTTPGetAction{
			Path:        k8sProbe.HTTPGet.Path,
			Port:        port,
			Scheme:      string(k8sProbe.HTTPGet.Scheme),
			HTTPHeaders: p.getProbeHTTPHeadersFrom(k8sProbe.HTTPGet.HTTPHeaders),
		}
	}

	if k8sProbe.TCPSocket != nil {
		port, err := getPortValue(k8sProbe.TCPSocket.Port, containerPorts)
		if err != nil {
			return nil, err
		}
		probe.TCPSocket = &workload_models.V1TCPSocketAction{
			Port: port,
		}
	}

	return probe, nil
}

func getPortValue(port intstr.IntOrString, containerPorts []v1.ContainerPort) (int32, error) {
	var portValue int32
	switch port.Type {
	case intstr.Int:
		portValue = int32(port.IntValue())
	case intstr.String:
		portName := port.String()
		for _, p := range containerPorts {
			if portName == p.Name {
				portValue = p.ContainerPort
				break
			}
		}
		if portValue == 0 {
			return 0, fmt.Errorf("unable to find named port: %s", portName)
		}
	}
	return portValue, nil
}
