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
 * Copyright 2023 Red Hat, Inc.
 */

package sched

import (
	"encoding/json"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/klog/v2"

	"sigs.k8s.io/yaml"

	"github.com/k8stopologyawareschedwg/deployer/pkg/manifests"
)

func SchedulerConfig(cm *corev1.ConfigMap, schedulerName string, params *manifests.ConfigParams) error {
	if cm.Data == nil {
		return fmt.Errorf("no data found in ConfigMap: %s/%s", cm.Namespace, cm.Name)
	}

	data, ok := cm.Data[manifests.SchedulerConfigFileName]
	if !ok {
		return fmt.Errorf("no data key named: %s found in ConfigMap: %s/%s", manifests.SchedulerConfigFileName, cm.Namespace, cm.Name)
	}

	newData, _, err := RenderConfig([]byte(data), schedulerName, params)
	if err != nil {
		return err
	}

	cm.Data[manifests.SchedulerConfigFileName] = string(newData)
	return nil
}

func RenderConfig(data []byte, schedulerName string, params *manifests.ConfigParams) ([]byte, bool, error) {
	if schedulerName == "" || params == nil {
		klog.InfoS("missing parameters, passing through", "schedulerName", schedulerName, "params", toJSON(params))
		return data, false, nil
	}

	var r unstructured.Unstructured
	if err := yaml.Unmarshal(data, &r.Object); err != nil {
		klog.ErrorS(err, "cannot unmarshal scheduler config, passing through")
		return data, false, err
	}

	updated := false

	if params.LeaderElection != nil {
		lead, ok, err := unstructured.NestedMap(r.Object, "leaderElection")
		if !ok || err != nil {
			klog.ErrorS(err, "failed to process unstructured data", "leaderElection", ok)
			return data, false, err
		}

		leadUpdated, err := updateLeaderElection(lead, params)
		if err != nil {
			klog.ErrorS(err, "failed to update unstructured data", "leaderElection", lead, "params", params)
			return data, false, err
		}
		if leadUpdated {
			updated = true
		}

		if err := unstructured.SetNestedMap(r.Object, lead, "leaderElection"); err != nil {
			klog.ErrorS(err, "failed to override unstructured data", "data", "leaderElection")
			return data, false, err
		}

	}

	profiles, ok, err := unstructured.NestedSlice(r.Object, "profiles")
	if !ok || err != nil {
		klog.ErrorS(err, "failed to process unstructured data", "profiles", ok)
		return data, false, err
	}
	for _, prof := range profiles {
		profile, ok := prof.(map[string]interface{})
		if !ok {
			klog.InfoS("unexpected profile data")
			return data, false, nil
		}

		profileName, ok, err := unstructured.NestedString(profile, "schedulerName")
		if !ok || err != nil {
			klog.ErrorS(err, "failed to get profile name", "profileName", ok)
			return data, false, err
		}

		if profileName != schedulerName {
			continue
		}

		if params.ProfileName != "" {
			err = unstructured.SetNestedField(profile, params.ProfileName, "schedulerName")
			if err != nil {
				klog.ErrorS(err, "failed to update unstructured data", "schedulerName", params.ProfileName)
				return data, false, err
			}
			updated = true
		}

		pluginConfigs, ok, err := unstructured.NestedSlice(profile, "pluginConfig")
		if !ok || err != nil {
			klog.ErrorS(err, "failed to process unstructured data", "pluginConfig", ok)
			return data, false, err
		}
		for _, plConf := range pluginConfigs {
			pluginConf, ok := plConf.(map[string]interface{})
			if !ok {
				klog.V(1).InfoS("unexpected profile coonfig data")
				return data, false, nil
			}

			name, ok, err := unstructured.NestedString(pluginConf, "name")
			if !ok || err != nil {
				klog.ErrorS(err, "failed to process unstructured data", "name", ok)
				return data, false, err
			}
			if name != manifests.SchedulerPluginName {
				continue
			}
			args, ok, err := unstructured.NestedMap(pluginConf, "args")
			if !ok || err != nil {
				klog.ErrorS(err, "failed to process unstructured data", "args", ok)
				return data, false, err
			}

			argsUpdated, err := updateArgs(args, params)
			if err != nil {
				klog.ErrorS(err, "failed to update unstructured data", "args", args, "params", params)
				return data, false, err
			}
			if argsUpdated {
				updated = true
			}

			if err := unstructured.SetNestedMap(pluginConf, args, "args"); err != nil {
				klog.ErrorS(err, "failed to override unstructured data", "data", "args")
				return data, false, err
			}
		}

		if err := unstructured.SetNestedSlice(profile, pluginConfigs, "pluginConfig"); err != nil {
			klog.ErrorS(err, "failed to override unstructured data", "data", "pluginConfig")
			return data, false, err
		}
	}

	if err := unstructured.SetNestedSlice(r.Object, profiles, "profiles"); err != nil {
		klog.ErrorS(err, "failed to override unstructured data", "data", "profiles")
		return data, false, err
	}

	newData, err := yaml.Marshal(&r.Object)
	if err != nil {
		klog.ErrorS(err, "cannot re-encode scheduler config, passing through")
		return data, false, nil
	}
	return newData, updated, nil
}

func updateLeaderElection(lead map[string]interface{}, params *manifests.ConfigParams) (bool, error) {
	var updated int
	var err error

	err = unstructured.SetNestedField(lead, params.LeaderElection.LeaderElect, "leaderElect")
	if err != nil {
		return updated > 0, err
	}
	updated++

	err = unstructured.SetNestedField(lead, params.LeaderElection.ResourceName, "resourceName")
	if err != nil {
		return updated > 0, err
	}
	updated++

	err = unstructured.SetNestedField(lead, params.LeaderElection.ResourceNamespace, "resourceNamespace")
	if err != nil {
		return updated > 0, err
	}
	updated++

	return updated > 0, nil

}

func updateArgs(args map[string]interface{}, params *manifests.ConfigParams) (bool, error) {
	var updated int
	var err error

	if params.Cache != nil {
		if params.Cache.ResyncPeriodSeconds != nil {
			resyncPeriod := *params.Cache.ResyncPeriodSeconds // shortcut
			err = unstructured.SetNestedField(args, resyncPeriod, "cacheResyncPeriodSeconds")
			if err != nil {
				return updated > 0, err
			}
			updated++
		}
	}

	cacheArgs, ok, err := unstructured.NestedMap(args, "cache")
	if !ok {
		cacheArgs = make(map[string]interface{})
	}
	if err != nil {
		return updated > 0, err
	}

	var cacheArgsUpdated int
	if params.Cache != nil {
		cacheArgsUpdated, err = updateCacheArgs(cacheArgs, params)
		if err != nil {
			return updated > 0, err
		}
	}
	updated += cacheArgsUpdated

	if cacheArgsUpdated > 0 {
		if err := unstructured.SetNestedMap(args, cacheArgs, "cache"); err != nil {
			return updated > 0, err
		}
	}

	scoringStratArgs, ok, err := unstructured.NestedMap(args, "scoringStrategy")
	if !ok {
		scoringStratArgs = make(map[string]interface{})
	}
	if err != nil {
		return updated > 0, err
	}

	var scoringStratArgsUpdated int
	if params.ScoringStrategy != nil {
		scoringStratArgsUpdated, err = updateScoringStrategyArgs(scoringStratArgs, params)
		if err != nil {
			return updated > 0, err
		}
	}
	updated += scoringStratArgsUpdated

	if scoringStratArgsUpdated > 0 {
		if err := unstructured.SetNestedMap(args, scoringStratArgs, "scoringStrategy"); err != nil {
			return updated > 0, err
		}
	}

	return updated > 0, ensureBackwardCompatibility(args)
}

func updateCacheArgs(args map[string]interface{}, params *manifests.ConfigParams) (int, error) {
	var updated int
	var err error

	if params.Cache.ResyncMethod != nil {
		resyncMethod := *params.Cache.ResyncMethod // shortcut
		err = manifests.ValidateCacheResyncMethod(resyncMethod)
		if err != nil {
			return updated, err
		}
		err = unstructured.SetNestedField(args, resyncMethod, "resyncMethod")
		if err != nil {
			return updated, err
		}
		updated++
	}
	if params.Cache.ForeignPodsDetectMode != nil {
		foreignPodsMode := *params.Cache.ForeignPodsDetectMode // shortcut
		err = manifests.ValidateForeignPodsDetectMode(foreignPodsMode)
		if err != nil {
			return updated, err
		}
		err = unstructured.SetNestedField(args, foreignPodsMode, "foreignPodsDetect")
		if err != nil {
			return updated, err
		}
		updated++
	}
	if params.Cache.InformerMode != nil {
		informerMode := *params.Cache.InformerMode // shortcut
		err = manifests.ValidateCacheInformerMode(informerMode)
		if err != nil {
			return updated, err
		}
		err = unstructured.SetNestedField(args, informerMode, "informerMode")
		if err != nil {
			return updated, err
		}
		updated++
	}

	return updated, nil
}

func updateScoringStrategyArgs(args map[string]interface{}, params *manifests.ConfigParams) (int, error) {
	var updated int
	var err error

	if params.ScoringStrategy.Type != "" {
		scoringStratType := params.ScoringStrategy.Type // shortcut
		err = manifests.ValidateScoringStrategyType(scoringStratType)
		if err != nil {
			return updated, err
		}
		err = unstructured.SetNestedField(args, scoringStratType, "type")
		if err != nil {
			return updated, err
		}
		updated++
	}

	if len(params.ScoringStrategy.Resources) > 0 {
		var resources []interface{}
		for _, scRes := range params.ScoringStrategy.Resources {
			resources = append(resources, map[string]interface{}{
				"name":   scRes.Name,
				"weight": scRes.Weight,
			})
		}
		if err := unstructured.SetNestedSlice(args, resources, "resources"); err != nil {
			return updated, err
		}
		updated++
	}

	return updated, nil
}

func ensureBackwardCompatibility(args map[string]interface{}) error {
	resyncPeriod, ok, err := unstructured.NestedInt64(args, "cacheResyncPeriodSeconds")
	if !ok {
		// nothing to do
		return nil
	}
	if err != nil {
		return fmt.Errorf("cannot process field cacheResyncPeriodSeconds: %w", err)
	}
	if resyncPeriod > 0 {
		// nothing to do
	} else {
		// remove for backward compatibility
		delete(args, "cacheResyncPeriodSeconds")
	}
	return nil
}

func toJSON(v any) string {
	if v == nil {
		return "<nil>"
	}
	data, err := json.Marshal(v)
	if err != nil {
		return fmt.Sprintf("<err=%v>", err)
	}
	return string(data)
}
