package main

import (
	"flag"

	diskapi "github.com/kubernetes-csi/csi-proxy/internal/os/disk"
	filesystemapi "github.com/kubernetes-csi/csi-proxy/internal/os/filesystem"
	smbapi "github.com/kubernetes-csi/csi-proxy/internal/os/smb"
	volumeapi "github.com/kubernetes-csi/csi-proxy/internal/os/volume"
	"github.com/kubernetes-csi/csi-proxy/internal/server"
	disksrv "github.com/kubernetes-csi/csi-proxy/internal/server/disk"
	filesystemsrv "github.com/kubernetes-csi/csi-proxy/internal/server/filesystem"
	smbsrv "github.com/kubernetes-csi/csi-proxy/internal/server/smb"
	srvtypes "github.com/kubernetes-csi/csi-proxy/internal/server/types"
	volumesrv "github.com/kubernetes-csi/csi-proxy/internal/server/volume"
	"k8s.io/klog"
)

var (
	kubeletCSIPluginsPath = flag.String("kubelet-csi-plugins-path", `C:\var\lib\kubelet`, "Prefix path of the Kubelet plugin directory in the host file system")
	kubeletPodPath        = flag.String("kubelet-pod-path", `C:\var\lib\kubelet`, "Prefix path of the kubelet pod directory in the host file system")
)

func main() {
	defer klog.Flush()
	klog.InitFlags(nil)

	flag.Parse()
	apiGroups, err := apiGroups()
	if err != nil {
		panic(err)
	}
	s := server.NewServer(apiGroups...)

	klog.Info("Starting CSI-Proxy Server ...")
	klog.Infof("Version: %s", version)
	if err := s.Start(nil); err != nil {
		panic(err)
	}
}

// apiGroups returns the list of enabled API groups.
func apiGroups() ([]srvtypes.APIGroup, error) {
	fssrv, err := filesystemsrv.NewServer(*kubeletCSIPluginsPath, *kubeletPodPath, filesystemapi.New())
	if err != nil {
		return []srvtypes.APIGroup{}, err
	}

	volumesrv, err := volumesrv.NewServer(volumeapi.New())
	if err != nil {
		return []srvtypes.APIGroup{}, err
	}

	disksrv, err := disksrv.NewServer(diskapi.New())
	if err != nil {
		return []srvtypes.APIGroup{}, err
	}

	smbsrv, err := smbsrv.NewServer(smbapi.New(), fssrv)
	if err != nil {
		return []srvtypes.APIGroup{}, err
	}

	return []srvtypes.APIGroup{
		fssrv,
		disksrv,
		volumesrv,
		smbsrv,
	}, nil
}
