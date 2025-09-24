package devpreview_test

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/spf13/cobra"

	"github.com/red-hat-storage/odf-cli/cmd/odf/devpreview"
)

const output = "OK"

func TestConfigure(t *testing.T) {
	t.Run("root", func(t *testing.T) {
		root := testCommand()
		devpreview.Configure(root)
		checkCommand(t, root)
	})

	t.Run("child", func(t *testing.T) {
		root := testCommand()
		devpreview.Configure(root)
		checkCommand(t, root, "child")
	})

	t.Run("grandchild1", func(t *testing.T) {
		root := testCommand()
		devpreview.Configure(root)
		checkCommand(t, root, "child", "grandchild1")
	})

	t.Run("grandchild2", func(t *testing.T) {
		root := testCommand()
		devpreview.Configure(root)
		checkCommand(t, root, "child", "grandchild2")
	})
}

// Cobra remembers state from Execute() (e.g. --help flag), so we must use a
// new command for every test.
func testCommand() *cobra.Command {
	root := &cobra.Command{
		Use:   "root",
		Short: "root short",
		Long:  "root long description",
	}

	child := &cobra.Command{
		Use:   "child",
		Short: "child short",
		Long:  "child long description",
	}

	grandChild1 := &cobra.Command{
		Use:   "grandchild1",
		Short: "grand child 1 short",
		Long:  "grand child 1 long description",
		Run:   printOutput,
	}

	grandChild2 := &cobra.Command{
		Use:   "grandchild2",
		Short: "grand child 2 short",
		Long:  "grand child 2 long description",
		Run:   printOutput,
	}

	child.AddCommand(grandChild1, grandChild2)
	root.AddCommand(child)

	return root
}

func checkCommand(t *testing.T, root *cobra.Command, args ...string) {
	args = append(args, "--help")
	cmd, help, err := executeCommand(root, args...)
	if err != nil {
		t.Fatalf("running %q with args %q failed: %v", cmd.CommandPath(), args, err)
	}

	if cmd == root {
		if !cmd.Hidden {
			t.Fatal("root command should not be hidden")
		}
	} else if cmd.Hidden {
		t.Fatal("child commands should not be hidden")
	}

	if !strings.HasSuffix(cmd.Short, devpreview.Suffix) {
		t.Fatalf("command %q short does not include the developer preview suffix: %q",
			cmd.CommandPath(), cmd.Short)
	}

	if !strings.HasSuffix(cmd.Long, "\n"+devpreview.Note) {
		t.Fatalf("command %q long does not include the developer preview note: %q",
			cmd.CommandPath(), cmd.Long)
	}

	if !strings.Contains(help, devpreview.Note) {
		t.Fatalf("command %q help text does not include the developer preview note: %v",
			cmd.CommandPath(), help)
	}

	// For manual inspection
	fmt.Println(help)
}

func printOutput(cmd *cobra.Command, args []string) {
	cmd.Print(output)
}

func executeCommand(root *cobra.Command, args ...string) (*cobra.Command, string, error) {
	var buf bytes.Buffer

	root.SetOut(&buf)
	root.SetErr(&buf)
	root.SetArgs(args)

	cmd, err := root.ExecuteC()

	return cmd, buf.String(), err
}
