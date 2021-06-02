// Copyright 2021 Cisco Systems, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package hostagent

import (
	netpolicy "github.com/k8snetworkplumbingwg/network-attachment-definition-client/pkg/apis/k8s.cni.cncf.io/v1"
//	v1 "github.com/k8snetworkplumbingwg/network-attachment-definition-client/pkg/apis/k8s.cni.cncf.io/v1"
	netClient "github.com/k8snetworkplumbingwg/network-attachment-definition-client/pkg/client/clientset/versioned"
	netattclient "github.com/k8snetworkplumbingwg/network-attachment-definition-client/pkg/client/clientset/versioned/typed/k8s.cni.cncf.io/v1"
//	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/tools/cache"
	"k8s.io/kubernetes/pkg/controller"
)

const (
	defaultAnnot        = "aci-cni/default-network"
	netAttachDefCRDName = "network-attachment-definitions.k8s.cni.cncf.io"
)

type NetworkAttachmentData struct {
	Name      string
	Namespace string
	Config    string
	Annot     string
}

type ClientInfo struct {
	NetClient netattclient.K8sCniCncfIoV1Interface
}

//func (agent *HostAgent) netAttDefInit(stopCh <-chan struct{}) {
//	agent.log.Debug("Initializing network-attachment-definition client")
//	restconfig := agent.env.RESTConfig()
//	netAttClient, err := netattclient.NewForConfig(restconfig)
//	if err != nil {
////		agent.log.Errorf("Failed to intialize network-attachment-definition client")
//		return
//	}
//
//	agent.initNetworkAttachmentDefinitionFromClient(netAttClient)
//cont.netAttDefInformer.Run(stopCh)
//}

func (agent *HostAgent) initNetworkAttDefInformerFromClient(
	netClientSet *netClient.Clientset) {

	agent.log.Debug("running initNetworkAttachmentDefinitionFromClient")
	//list, err := netClientSet.K8sCniCncfIoV1().NetworkAttachmentDefinitions("default").List(context.TODO(), metav1.ListOptions{})
	//if err != nil {
	//	agent.log.Error("Error listing all network attachment definitions: %v", err)
	//}
	agent.initNetworkAttachmentDefinitionInformerBase(
		cache.NewListWatchFromClient(
			netClientSet.RESTClient(), "network-attachment-definitions",
			"default", fields.Everything()))
}

//*v1.NetworkAttachmentDefinitionList

func (agent *HostAgent) initNetworkAttachmentDefinitionInformerBase(listWatch *cache.ListWatch) {
	agent.netAttDefInformer = cache.NewSharedIndexInformer(
		listWatch, &netpolicy.NetworkAttachmentDefinition{}, controller.NoResyncPeriodFunc(),
		cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc},
	)
	agent.netAttDefInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			agent.networkAttDefAdded(obj)
		},
		UpdateFunc: func(oldobj interface{}, newobj interface{}) {
			agent.networkAttDefChanged(oldobj, newobj)
		},
		DeleteFunc: func(obj interface{}) {
			agent.networkAttDefDeleted(obj)
		},
	})
}

func (agent *HostAgent) networkAttDefAdded(obj interface{}) {
	ntd := obj.(*netpolicy.NetworkAttachmentDefinition)
	agent.log.Infof("network atttachment Added: %s", ntd.ObjectMeta.Name)
	netattdata := NetworkAttachmentData{
		Name:      ntd.ObjectMeta.Name,
		Namespace: ntd.ObjectMeta.Namespace,
		Config:    ntd.Spec.Config,
		//Annot:     ntd.Spec.Annotations[defaultAnnot],
	}
	agent.log.Debug("Name", netattdata.Name)
	agent.log.Debug("Namespace", netattdata.Namespace)
	agent.log.Debug("Config", netattdata.Config)
	//agent.log.Infof("Annotion", netattdata.Annot)

	ntdKey := agent.getnetattKey(&netattdata)
	agent.netattdefmap[ntdKey] = &netattdata
}

func (agent *HostAgent) getnetattKey(netdata *NetworkAttachmentData) string {
	return netdata.Name + "/" + netdata.Namespace
}

func (agent *HostAgent) networkAttDefChanged(oldobj interface{}, newobj interface{}) {
	//	ntd := newobj.(*netpolicy.NetworkAttachmentDefinition)
}

func (agent *HostAgent) networkAttDefDeleted(obj interface{}) {
	//	ntd := obj.(*netpolicy.NetworkAttachmentDefinition)
}
