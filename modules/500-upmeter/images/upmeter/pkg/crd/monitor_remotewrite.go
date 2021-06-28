/*
Copyright 2021 Flant CJSC

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package crd

import (
	"context"
	"fmt"

	"github.com/flant/shell-operator/pkg/kube"
	"github.com/flant/shell-operator/pkg/kube_events_manager"
	"github.com/flant/shell-operator/pkg/kube_events_manager/types"
	log "github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"

	v1 "d8.io/upmeter/pkg/crd/v1"
)

type RemoteWriteMonitor struct {
	monitor kube_events_manager.Monitor
	logger  *log.Entry
}

func NewRemoteWriteMonitor(kubeClient kube.KubernetesClient, logger *log.Entry) *RemoteWriteMonitor {
	monitor := kube_events_manager.NewMonitor()
	monitor.WithKubeClient(kubeClient)

	return &RemoteWriteMonitor{
		monitor: monitor,
		logger:  logger,
	}
}

func (m *RemoteWriteMonitor) Start(ctx context.Context) error {
	config := &kube_events_manager.MonitorConfig{
		Metadata: struct {
			MonitorId    string
			DebugName    string
			LogLabels    map[string]string
			MetricLabels map[string]string
		}{
			"upmeterremotewrite-crd",
			"upmeterremotewrite-crd",
			map[string]string{},
			map[string]string{},
		},
		EventTypes: []types.WatchEventType{
			types.WatchEventAdded,
			types.WatchEventModified,
			types.WatchEventDeleted,
		},
		ApiVersion:              "deckhouse.io/v1",
		Kind:                    "UpmeterRemoteWrite",
		LogEntry:                m.logger.WithField("component", "upmeterremotewrite-monitor"),
		KeepFullObjectsInMemory: true,
	}

	m.monitor.WithContext(ctx)
	m.monitor.WithConfig(config)

	// Load initial CRD list
	err := m.monitor.CreateInformers()
	if err != nil {
		return fmt.Errorf("create informers: %v", err)
	}

	m.monitor.Start(ctx)
	return ctx.Err()
}

func (m *RemoteWriteMonitor) Stop() {
	m.monitor.Stop()
}

func (m *RemoteWriteMonitor) getLogger() *log.Entry {
	return m.monitor.GetConfig().LogEntry
}

func (m *RemoteWriteMonitor) Subscribe(handler RemoteWriteChangeHandler) {
	m.monitor.WithKubeEventCb(func(ev types.KubeEvent) {
		// One event and one object per change, we always have single items in these lists.
		evType := ev.WatchEvents[0]
		obj := ev.Objects[0].Object

		rw, err := convert(obj)
		if err != nil {
			m.getLogger().Errorf("cannot convert UpmeterRemoteWrite object: %v", err)
			return
		}

		switch evType {
		case types.WatchEventAdded:
			handler.OnAdd(rw)
		case types.WatchEventModified:
			handler.OnModify(rw)
		case types.WatchEventDeleted:
			handler.OnDelete(rw)
		}
	})
}

func (m *RemoteWriteMonitor) List() ([]*v1.RemoteWrite, error) {
	res := make([]*v1.RemoteWrite, 0)
	for _, obj := range m.monitor.GetExistedObjects() {
		rw, err := convert(obj.Object)
		if err != nil {
			return nil, err
		}
		res = append(res, rw)
	}
	return res, nil
}

func convert(o *unstructured.Unstructured) (*v1.RemoteWrite, error) {
	var rw v1.RemoteWrite
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(o.UnstructuredContent(), &rw)
	if err != nil {
		return nil, fmt.Errorf("cannot convert unstructured to v1.RemoteWrite: %v", err)
	}
	return &rw, nil
}

type RemoteWriteChangeHandler interface {
	OnAdd(*v1.RemoteWrite)
	OnModify(*v1.RemoteWrite)
	OnDelete(*v1.RemoteWrite)
}
