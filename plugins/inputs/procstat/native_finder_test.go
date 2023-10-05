package procstat

import (
	"fmt"
	"os/exec"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
)

func BenchmarkPattern(b *testing.B) {
	f, err := NewNativeFinder()
	require.NoError(b, err)
	for n := 0; n < b.N; n++ {
		_, err := f.Pattern(".*")
		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkFullPattern(b *testing.B) {
	f, err := NewNativeFinder()
	require.NoError(b, err)
	for n := 0; n < b.N; n++ {
		_, err := f.FullPattern(".*")
		if err != nil {
			panic(err)
		}
	}
}

func TestChildPattern(t *testing.T) {
	if runtime.GOOS == "linux" || runtime.GOOS == "darwin" {
		cmd := exec.Command("/bin/bash", "-c", "sleep 30")
		if err := cmd.Start(); err != nil {
			fmt.Printf("Error starting command: %s\n", err)
			return
		}

		f, err := NewNativeFinder()
		require.NoError(t, err)

		childpids, err := f.ChildPattern("sleep")
		fmt.Println("pid in childpids which pattern sleep...")
		for _, p := range childpids {
			t.Log(string(p))
			fmt.Println(string(p))
		}
		fmt.Println("cmd infos ...")
		fmt.Println(cmd.Path)
		for _, arg := range cmd.Args {
			fmt.Println(arg)
		}

		fmt.Println("cmd infos get by ps -ef ...")
		if len(childpids) > 0 {
			cmd2 := exec.Command("/bin/bash", "-c", "ps -ef |grep sleep")
			fmt.Println(cmd2.Stdout)
		}

		fmt.Println("require.Contains ...")
		require.Contains(t, childpids, PID(cmd.Process.Pid))
		//require.Equal(t, []PID{PID(cmd.Process.Pid)}, childpids)
		cmd.Process.Kill()
		if err != nil {
			panic(err)
		}

		var nilpids []PID
		childpids, err = f.ChildPattern("sleep 30")
		for _, p := range childpids {
			t.Log(string(p))
		}

		require.Equal(t, nilpids, childpids)
		if err != nil {
			panic(err)
		}
	}
}
