package cmd_test

import (
	"bytes"
	"os"
	"os/exec"
	"path"
	"testing"

	"github.com/gohxs/boiler/cli/bp/cmd"
)

func init() {
	os.Chdir("../_test")
}

// Helper function to ignore multi arg error
/*func aerr(args ...interface{}) interface{} {
	return args[0]
}*/

func TestInit(t *testing.T) {
	proj := "test-proj1"
	defer os.RemoveAll(proj)

	tu := testUtil{t}

	buf := bytes.NewBuffer([]byte("Testing\ntest")) // First "Testing" second "test"

	cmd.Stdin = buf
	os.Args = []string{"", "new", "boilerplate", proj}
	err := cmd.RootCmd.Execute()
	tu.eq(err, nil)

	cmd := exec.Command("go", "run", path.Join(proj, "hello.go"))
	outbuf := bytes.NewBuffer([]byte{})
	cmd.Stdout = outbuf

	err = cmd.Run()
	tu.eq(err, nil)
	tu.eq(outbuf.String(), "Hello Testing this is an test app")

}

func TestInit2(t *testing.T) {
	proj := "test-proj2"
	defer os.RemoveAll(proj)
	tu := testUtil{t}

	os.Args = []string{"", "new", "boilerplate2", proj}
	err := cmd.RootCmd.Execute()
	tu.eq(err, nil)

	_, err = os.Stat(path.Join(proj, "file.go"))
	tu.eq(err, nil)

}

type testUtil struct{ *testing.T }

func (t testUtil) eq(a1, a2 interface{}) {
	if a1 != a2 {
		t.Fatalf("Does not match '%v' != '%v'", a1, a2)
	}
}
