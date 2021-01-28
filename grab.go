package stacker

import (
	"fmt"
	"os"
	"path"

	"github.com/anuvu/stacker/types"
	"github.com/pkg/errors"
)

func Grab(sc types.StackerConfig, storage types.Storage, name string, source string, targetDir string, hash string) error {
	c, err := NewContainer(sc, storage, name)
	if err != nil {
		return err
	}
	defer c.Close()

	err = c.bindMount(targetDir, "/stacker", "")
	if err != nil {
		return err
	}
	defer os.Remove(path.Join(sc.RootFSDir, name, "rootfs", "stacker"))

	if len(hash) > 0 {
		if err = c.Execute(fmt.Sprintf("echo %s %s | sha256sum --check", hash, source), nil); err != nil {
			return errors.Errorf("The requested hash of %s import is different than the actual hash: %s",
				source, hash)
		}
	}

	return c.Execute(fmt.Sprintf("cp -a %s /stacker", source), nil)
}
