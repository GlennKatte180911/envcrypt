// Package env — Group
//
// The Group family of functions partitions a []Entry slice into named buckets,
// making it straightforward to process subsets of an environment file as a
// unit.
//
// # Grouping by key prefix
//
//	groups := env.GroupByPrefix(entries, "_")
//	dbEntries := groups["DB"] // DB_HOST, DB_PORT, …
//
// # Grouping by key suffix
//
//	groups := env.GroupBySuffix(entries, "_")
//	hosts := groups["HOST"] // DB_HOST, APP_HOST, …
//
// # Custom grouping
//
//	groups := env.Group(entries, func(e env.Entry) string {
//		if strings.HasPrefix(e.Key, "INTERNAL_") {
//			return "internal"
//		}
//		return "public"
//	})
//
// # Flattening back to a slice
//
//	all := env.Flatten(groups)
//
// # Filtering a group by name
//
//	// GroupKeys returns the sorted list of bucket names present in a grouped map,
//	// which is useful for iterating over groups in a deterministic order.
//	names := env.GroupKeys(groups) // ["DB", "APP", "INTERNAL", …]
//	for _, name := range names {
//		process(name, groups[name])
//	}
package env
