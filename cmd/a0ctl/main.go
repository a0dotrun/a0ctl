package main

import "github.com/a0dotrun/a0ctl/internal/command/root"

func main() {
	cmd := root.New()
	cmd.Execute()
}

//
// func main() {
// 	os.Exit(run())
// }
//
// func run() (exitCode int) {
// 	ctx, cancel := newContext()
// 	defer cancel()
//
// 	exitCode = cli.Run(ctx, os.Args[1:]...)
// 	return
// }
//
// func newContext() (context.Context, context.CancelFunc) {
// 	// NOTE: when signal.Notify is called for os.Interrupt it traps both
// 	// ^C (Control-C) and ^BREAK (Control-Break) on Windows.
// 	signals := []os.Signal{os.Interrupt}
// 	if runtime.GOOS != "windows" {
// 		signals = append(signals, syscall.SIGTERM)
// 	}
//
// 	return signal.NotifyContext(context.Background(), signals...)
// }
