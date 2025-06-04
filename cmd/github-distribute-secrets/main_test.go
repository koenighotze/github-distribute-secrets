package main

// func TestMain(t *testing.T) {
// 	t.Run("should use the dry run client if the flag is provided", func(t *testing.T) {
// 		origArgs := os.Args
// 		os.Args = []string{"cmd", "--dry-run"}
// 		defer func() { os.Args = origArgs }()

// 		exitCalled := false
// 		origOsExit := osExit
// 		osExit = func(code int) { exitCalled = true }
// 		defer func() { osExit = origOsExit }()

// 		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

// 		main()
// 	})
// 	t.Run("should use the default client if the flag is omitted", func(t *testing.T) {})
// 	t.Run("should distribute the secrets", func(t *testing.T) {})

// 	t.Run("should exit with -1 if setting the secrets fails", func(t *testing.T) {})
// }
