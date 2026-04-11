// Package merge provides multi-source environment variable merging
// with configurable conflict resolution strategies.
//
// # Strategies
//
// Three strategies are available:
//
//   - StrategyLast  – the last source to define a key wins (default)
//   - StrategyFirst – the first definition is kept; later sources are ignored
//   - StrategyError – any duplicate key causes Merge to return an error
//
// # Usage
//
//	base := map[string]string{"PORT": "8080", "DEBUG": "false"}
//	override := map[string]string{"DEBUG": "true", "LOG_LEVEL": "info"}
//
//	res, err := merge.Merge(merge.StrategyLast, base, override)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println(res.Env["DEBUG"]) // "true"
package merge
