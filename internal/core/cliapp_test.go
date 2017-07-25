package core_test

import (
	"bytes"
	"os"
	"os/exec"
	"testing"

	"github.com/gohxs/boiler/internal/core"
)

// Helper function to ignore multi arg error
func aerr(args ...interface{}) interface{} {
	return args[0]
}

func TestInit(t *testing.T) {

	tu := testUtil{t}

	buf := bytes.NewBuffer([]byte("Testing\ntest"))

	app := core.NewApp(buf)

	os.Chdir("../../test")
	err := app.Run([]string{"", "init", "boilerplate", "test-proj1"})
	tu.eq(err, nil)

	defer os.RemoveAll("test-proj1")

	cmd := exec.Command("go", "run", "test-proj1/hello.go")
	outbuf := bytes.NewBuffer([]byte{})
	cmd.Stdout = outbuf

	err = cmd.Run()
	tu.eq(err, nil)

	tu.eq(outbuf.String(), "Hello Testing")

}

type testUtil struct{ *testing.T }

func (t testUtil) eq(a1, a2 interface{}) {
	if a1 != a2 {
		t.Fatalf("Does not match '%v' != '%v'", a1, a2)
	}
}
