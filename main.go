package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"slices"

	v1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/util/sets"
)

type resource struct {
	verbs sets.Set[string]
	name  []string
}

func main() {
	clusterRole, err := loadResource()
	if err != nil {
		slog.Error("error loading resource", "error", err)
		return
	}

	ruleStore := make(map[struct{ apiGroup, resource string }]resource)

	for _, rule := range clusterRole.Rules {
		for _, apiGroup := range rule.APIGroups {
			if apiGroup == "" {
				apiGroup = "core"
			}

			for _, r := range rule.Resources {
				key := struct{ apiGroup, resource string }{apiGroup, r}

				res, ok := ruleStore[key]
				if !ok {
					ruleStore[key] = resource{
						verbs: sets.New(rule.Verbs...),
						name:  rule.ResourceNames,
					}
				} else {
					res.verbs.Insert(rule.Verbs...)

					ruleStore[key] = res
				}
			}
		}
	}

	for key, value := range ruleStore {
		verbs := value.verbs.UnsortedList()

		slices.Sort(verbs)

		fmt.Printf("Composite Key: %v/%v Verbs: %v Resource Names: %v\n", key.apiGroup, key.resource, verbs, value.name)
	}
}

func loadResource() (*v1.ClusterRole, error) {
	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		return nil, err
	}

	var clusterRole v1.ClusterRole

	err = json.Unmarshal(data, &clusterRole)
	if err != nil {
		return nil, err
	}

	return &clusterRole, nil
}
