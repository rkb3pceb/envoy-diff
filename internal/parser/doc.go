// Package parser implements readers for various environment variable
// source formats used by envoy-diff.
//
// Supported formats:
//
//   - .env files (KEY=VALUE, with optional quoting and comment support)
//
// Parsed results are returned as EnvMap (map[string]string) values that
// can be passed downstream to the differ and auditor components.
//
// Example usage:
//
//	f, err := os.Open(".env")
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer f.Close()
//
//	env, err := parser.ParseEnvFile(f)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println(env["APP_NAME"])
package parser
