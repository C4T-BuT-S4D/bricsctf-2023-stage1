package main

import (
	"fmt"
	"os"
	"os/exec"
)

type Emulator struct {
	snapshotName string
} 

func NewEmulator(snapshot string) *Emulator {
	return &Emulator{snapshotName: snapshot}
}

func (e *Emulator) InstallApp(path string) error {
	cmd := exec.Command("adb", "install", path)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		stderr := string(err.(*exec.ExitError).Stderr)
		return fmt.Errorf("error installing app: %v\n", stderr)
	}
	return nil
}

func (e *Emulator) StartApp(pkgName string) error {
	cmd := exec.Command("adb", "shell", 
		fmt.Sprintf("monkey -p %s -c android.intent.category.LAUNCHER 1", pkgName))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		stderr := string(err.(*exec.ExitError).Stderr)
		return fmt.Errorf("error starting app: %v", stderr)
	}
	return nil
}

func (e *Emulator) Reset() error {
	cmd := exec.Command("adb", "emu", 
		fmt.Sprintf("avd snapshot load %s", e.snapshotName))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		stderr := string(err.(*exec.ExitError).Stderr)
		return fmt.Errorf("error restoring snapshot: %v", stderr)
	}
	return nil
}

