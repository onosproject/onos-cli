// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package topo

import (
	"strings"

	topoapi "github.com/onosproject/onos-api/go/onos/topo"
	"github.com/spf13/cobra"
)

func compileFilters(cmd *cobra.Command, objectType topoapi.Object_Type) *topoapi.Filters {
	aspects, _ := cmd.Flags().GetStringSlice("with-aspect")
	filters := &topoapi.Filters{WithAspects: aspects}
	lq, _ := cmd.Flags().GetString("label")
	filters.LabelFilters = compileLabelFilters(lq)
	filters.ObjectTypes = []topoapi.Object_Type{objectType}
	if objectType == topoapi.Object_ENTITY || objectType == topoapi.Object_RELATION {
		kq, _ := cmd.Flags().GetString("kind")
		filters.KindFilter = compileKindFilter(kq)
	}
	return filters
}

func compileLabelFilters(query string) []*topoapi.Filter {
	filters := make([]*topoapi.Filter, 0)
	fields := strings.Split(query, " && ")
	for _, field := range fields {
		filter, _ := compileLabelFilter(strings.TrimSpace(field))
		if filter != nil {
			filters = append(filters, filter)
		}
	}
	return filters
}

func compileLabelFilter(field string) (*topoapi.Filter, error) {
	if strings.Contains(field, " !in (") {
		key := extractKey(field, " !in (")
		values := extractValues(field)
		return &topoapi.Filter{
			Filter: &topoapi.Filter_Not{Not: &topoapi.NotFilter{
				Inner: &topoapi.Filter{Filter: &topoapi.Filter_In{In: &topoapi.InFilter{Values: values}}}},
			},
			Key: key,
		}, nil

	} else if strings.Contains(field, " in (") {
		key := extractKey(field, " in (")
		values := extractValues(field)
		return &topoapi.Filter{
			Filter: &topoapi.Filter_In{In: &topoapi.InFilter{Values: values}},
			Key:    key,
		}, nil

	} else if strings.Contains(field, "!=") {
		key := extractKey(field, "!=")
		value := extractValue(field)
		return &topoapi.Filter{
			Filter: &topoapi.Filter_Not{Not: &topoapi.NotFilter{
				Inner: &topoapi.Filter{Filter: &topoapi.Filter_Equal_{Equal_: &topoapi.EqualFilter{Value: value}}}},
			},
			Key: key,
		}, nil

	} else if strings.Contains(field, "=") {
		key := extractKey(field, "=")
		value := extractValue(field)
		return &topoapi.Filter{
			Filter: &topoapi.Filter_Equal_{Equal_: &topoapi.EqualFilter{Value: value}},
			Key:    key,
		}, nil

	}
	return nil, nil
}

func extractKey(field string, sep string) string {
	return strings.TrimSpace(strings.Split(field, sep)[0])
}

func extractValue(field string) string {
	return strings.TrimSpace(strings.Split(field, "=")[1])
}

func extractValues(field string) []string {
	gs := strings.Split(strings.Split(strings.Split(field, "(")[1], ")")[0], ",")
	values := make([]string, 0, len(gs))
	for _, v := range gs {
		values = append(values, strings.TrimSpace(v))
	}
	return values
}

func compileKindFilter(query string) *topoapi.Filter {
	if len(query) == 0 {
		return nil
	}

	// parse queries of form "!in (a, b, c)"
	// do this before the positive form of inclusion b/c containing !in ( is more restrictive and precise
	if strings.Contains(query, "!in (") {
		values := extractValues(query)
		return &topoapi.Filter{
			Filter: &topoapi.Filter_Not{Not: &topoapi.NotFilter{
				Inner: &topoapi.Filter{Filter: &topoapi.Filter_In{In: &topoapi.InFilter{Values: values}}}},
			},
		}
		// parse queries of form "in ( a,b,c )"
	} else if strings.Contains(query, "in (") {
		values := extractValues(query)
		return &topoapi.Filter{
			Filter: &topoapi.Filter_In{In: &topoapi.InFilter{Values: values}},
		}
		// parse queries of form "!= a", again do more restrictive Contains check first
	} else if strings.Contains(query, "!=") {
		value := extractValue(query)
		return &topoapi.Filter{
			Filter: &topoapi.Filter_Not{Not: &topoapi.NotFilter{
				Inner: &topoapi.Filter{Filter: &topoapi.Filter_Equal_{Equal_: &topoapi.EqualFilter{Value: value}}}},
			},
		}
		// parse queries of form "= a"
	} else if strings.Contains(query, "=") {
		value := extractValue(query)
		return &topoapi.Filter{
			Filter: &topoapi.Filter_Equal_{Equal_: &topoapi.EqualFilter{Value: value}},
		}
		// parse queries of the form "a" or "a, b, c"
		// shortcut for "in(a)" (which is equivalent to "=a") or "in(a, b, c)"
	} else if !strings.Contains(query, "(") && !strings.Contains(query, ")") && !strings.Contains(query, "!") {
		values := extractValues(" (" + query + ") ")
		return &topoapi.Filter{
			Filter: &topoapi.Filter_In{In: &topoapi.InFilter{Values: values}},
		}
	}
	return nil
}
