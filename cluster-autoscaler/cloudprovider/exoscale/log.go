/*
Copyright 2021 The Kubernetes Authors.

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

package exoscale

import (
	"k8s.io/klog/v2"
)

func fatalf(format string, args ...interface{}) {
	klog.Fatalf("exoscale-provider: "+format, args...)
}

func errorf(format string, args ...interface{}) {
	klog.Errorf("exoscale-provider: "+format, args...)
}

func infof(format string, args ...interface{}) {
	klog.Infof("exoscale-provider: "+format, args...)
}

func debugf(format string, args ...interface{}) {
	klog.V(3).Infof("exoscale-provider: "+format, args...)
}
