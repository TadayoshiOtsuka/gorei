package root

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
)

var (
	ErrNotEnoughArgs    = errors.New("not enough arguments e.g.) 'gorei <new package name>'")
	ErrModuleCanNotRead = errors.New("can not read module definition in 'go.mod'")
	ErrInterrupt        = errors.New("interrupted")
)

func Exec(args []string) error {
	if len(args) == 0 {
		return ErrNotEnoughArgs
	}
	newPkg := args[0]
	oldPkg, err := readOldPackageName()
	if err != nil {
		return err
	}
	res, err := confirm(oldPkg, newPkg)
	if err != nil {
		return err
	}
	if !res {
		return ErrInterrupt
	}

	cd, err := os.Getwd()
	if err != nil {
		return err
	}
	if err := scan(cd, oldPkg, newPkg); err != nil {
		return err
	}

	fmt.Printf("[%v] %v🎉\n", color.GreenString("gorei"), color.GreenString("done"))
	return nil
}

func readOldPackageName() (string, error) {
	fp, err := os.Open("./go.mod")
	if err != nil {
		return "", err
	}
	defer fp.Close()
	s := bufio.NewScanner(fp)
	s.Scan()
	c := s.Text()
	if !strings.Contains(c, "module") {
		return "", ErrModuleCanNotRead
	}

	return c[7:], nil
}

func confirm(old, new string) (bool, error) {
	label := fmt.Sprintf("[%v] Replace '%v' => '%v' ", color.GreenString("gorei"), old, new)
	p := promptui.Prompt{
		Label:     label,
		IsConfirm: true,
	}
	res, err := p.Run()
	if err != nil {
		return false, err
	}
	if res == "N" {
		return false, nil
	}

	return true, nil
}

func scan(src, oldPkg, newPkg string) error {
	fs, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, f := range fs {
		if strings.HasPrefix(f.Name(), ".") {
			continue
		}
		if f.IsDir() {
			scan(filepath.Join(src, f.Name()), oldPkg, newPkg)
		} else {
			genFile(src, src, f.Name(), oldPkg, newPkg)
		}
	}

	return nil
}

func genFile(src, dst, name, oldPkg, newPkg string) error {
	fs, fd := filepath.Join(src, name), filepath.Join(dst, name)
	file, err := os.ReadFile(fs)
	if err != nil {
		return err
	}

	f, err := os.Create(fd)
	if err != nil {
		return err
	}
	defer f.Close()

	file = replacePackageName(file, oldPkg, newPkg)
	if _, err = f.Write(file); err != nil {
		return err
	}
	fmt.Printf("[%v] replaced: %v\n", color.GreenString("gorei"), fs)

	return nil
}

func replacePackageName(file []byte, old, new string) []byte {
	c := string(file)
	c = strings.ReplaceAll(c, old, new)
	return []byte(c)
}
