package converter

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/panjf2000/ants/v2"

	"github.com/sirupsen/logrus"
)

var BaseDir string

func Convert() {
	info, err := os.Stat(BaseDir)
	if err != nil {
		logrus.Fatal(err)
	}
	if !info.IsDir() {
		logrus.Fatalf("%s is not a dir", BaseDir)
	}

	var count int
	_ = filepath.Walk(BaseDir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			count++
		}
		return nil
	})

	var wg sync.WaitGroup
	wg.Add(count)

	pool, err := ants.NewPool(100, ants.WithPreAlloc(true))
	if err != nil {
		logrus.Fatal(err)
	}
	defer pool.Release()

	_ = filepath.Walk(BaseDir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			err = pool.Submit(func() {
				defer func() {
					wg.Done()
					logrus.Infof("[Convert] %s", path)
				}()

				fname := filepath.Join(filepath.Dir(path), strings.ReplaceAll(filepath.Base(path), ".scel", ".rime"))
				args := []string{
					"-os:linux",
					"-ft:\"len:1-50|rm:space\"",
					"-ct:pinyin",
					"-i:scel",
					path,
					"-o:rime",
					fname,
				}
				// ImeWlConverterCmd -i:scel 全国县市乡镇名字.scel -os:linux ft:"len:1-50|rm:space" -ct:pinyin -o:rime ~/test
				cmd := exec.Command("ImeWlConverterCmd", args...)
				out, err := cmd.CombinedOutput()
				if err != nil {
					logrus.Errorf("failed to converter %s, error: %s: %s", path, err, string(out))
				}
			})
		}
		if err != nil {
			logrus.Error(err)
		}
		return nil
	})

	wg.Wait()

}
