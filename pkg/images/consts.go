/*
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * Copyright 2021 Red Hat, Inc.
 */

package images

const (
	SchedulerPluginSchedulerDefaultImageTag  = "registry.k8s.io/scheduler-plugins/kube-scheduler:v0.26.7"
	SchedulerPluginControllerDefaultImageTag = "registry.k8s.io/scheduler-plugins/controller:v0.26.7"
	NodeFeatureDiscoveryDefaultImageTag      = "registry.k8s.io/nfd/node-feature-discovery:v0.13.3"
	ResourceTopologyExporterDefaultImageTag  = "quay.io/k8stopologyawareschedwg/resource-topology-exporter:v0.13.0"
)

const (
	SchedulerPluginSchedulerDefaultImageSHA  = "registry.k8s.io/scheduler-plugins/kube-scheduler@sha256:1eb894c0743d5d01e18c74150645354554201aa6e0dcae62d49a8d3657326f1a"
	SchedulerPluginControllerDefaultImageSHA = "registry.k8s.io/scheduler-plugins/controller@sha256:6acd9f6ec8a1cd0b23f458577cf287fd70ce71febdfb843c08790ebfa3a8ebef"
	NodeFeatureDiscoveryDefaultImageSHA      = "registry.k8s.io/nfd/node-feature-discovery@sha256:7b56c7c5aa062607b8656c060befb8d9b32f7f32b99ae422b1da882a1b50adce"
	ResourceTopologyExporterDefaultImageSHA  = "quay.io/k8stopologyawareschedwg/resource-topology-exporter@sha256:046defed112c09f353c62d45c94265283004db6d42cc1c0f7904c824d75361e5"
)
