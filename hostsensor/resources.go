package hostsensor

// HostSensorResource is the enumerated type listing all resources from the host-scanner.
type HostSensorResource string

const (
	// host-scanner resources
	KubeletConfiguration         HostSensorResource = "KubeletConfiguration"
	OsReleaseFile                HostSensorResource = "OsReleaseFile"
	KernelVersion                HostSensorResource = "KernelVersion"
	LinuxSecurityHardeningStatus HostSensorResource = "LinuxSecurityHardeningStatus"
	OpenPortsList                HostSensorResource = "OpenPortsList"
	LinuxKernelVariables         HostSensorResource = "LinuxKernelVariables"
	KubeletCommandLine           HostSensorResource = "KubeletCommandLine"
	KubeletInfo                  HostSensorResource = "KubeletInfo"
	KubeProxyInfo                HostSensorResource = "KubeProxyInfo"
	ControlPlaneInfo             HostSensorResource = "ControlPlaneInfo"
	CloudProviderInfo            HostSensorResource = "CloudProviderInfo"
	CNIInfo                      HostSensorResource = "CNIInfo"
)

func (r HostSensorResource) String() string {
	return string(r)
}

// MapResourceToPlural maps scanner resource types to their CRD plural names
func MapResourceToPlural(r HostSensorResource) string {
	switch r {
	case OsReleaseFile:
		return "osreleasefiles"
	case KernelVersion:
		return "kernelversions"
	case LinuxSecurityHardeningStatus:
		return "linuxsecurityhardeningstatuses"
	case OpenPortsList:
		return "openportslists"
	case LinuxKernelVariables:
		return "linuxkernelvariables"
	case KubeletInfo:
		return "kubeletinfos"
	case KubeProxyInfo:
		return "kubeproxyinfos"
	case ControlPlaneInfo:
		return "controlplaneinfos"
	case CloudProviderInfo:
		return "cloudproviderinfos"
	case CNIInfo:
		return "cniinfos"
	default:
		return ""
	}
}

// MapHostSensorResourceToApiGroup returns the API group for the resource.
// Note: This returns the default group/version. Specific implementations might use different versions.
func MapHostSensorResourceToApiGroup(r HostSensorResource) string {
	switch r {
	case
		KubeletConfiguration,
		OsReleaseFile,
		KubeletCommandLine,
		KernelVersion,
		LinuxSecurityHardeningStatus,
		OpenPortsList,
		LinuxKernelVariables,
		KubeletInfo,
		KubeProxyInfo,
		ControlPlaneInfo,
		CloudProviderInfo,
		CNIInfo:
		return "hostdata.kubescape.cloud/v1beta0"
	default:
		return ""
	}
}
