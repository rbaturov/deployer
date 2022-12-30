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
	SchedulerPluginSchedulerDefaultImageTag  = "registry.k8s.io/scheduler-plugins/kube-scheduler:v0.24.9"
	SchedulerPluginControllerDefaultImageTag = "registry.k8s.io/scheduler-plugins/controller:v0.24.9"
	ResourceTopologyExporterDefaultImageTag  = "quay.io/k8stopologyawareschedwg/resource-topology-exporter:v0.6.0"
	NodeFeatureDiscoveryDefaultImageTag      = "gcr.io/k8s-staging-nfd/node-feature-discovery:v0.10.1"
)

const (
	SchedulerPluginSchedulerDefaultImageSHA  = "registry.k8s.io/scheduler-plugins/kube-scheduler@sha256:7e5681d6ee55da2a371111401fafd1ba371df83c4f1da088a0a7b20a2951eb72"
	SchedulerPluginControllerDefaultImageSHA = "registry.k8s.io/scheduler-plugins/controller@sha256:28ade406054565a06a7585acd839c420e1a0b8969f5bbe81ca7f0d97a4a483f3"
	ResourceTopologyExporterDefaultImageSHA  = "quay.io/k8stopologyawareschedwg/resource-topology-exporter@sha256:c4721e940250ef3a31600e70d4c933551492b22ed54c9593ad2a246da976d6f1"
	NodeFeatureDiscoveryDefaultImageSHA      = "gcr.io/k8s-staging-nfd/node-feature-discovery@sha256:4aebf17c8b72ee91cb468a6f21dd9f0312c1fcfdf8c86341f7aee0ec2d5991d7"
)
