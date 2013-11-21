package venv

import (
    "fmt"
    "io"
    "io/ioutil"
    "log"
    "os"
    "path"
    "github.com/ethanrowe/slingpie/util"
    "github.com/ethanrowe/slingpie/exec"
)

type Venv struct {
    Source string
    Target string
}

type VenvTwoSix struct {
    Venv
}

type VenvTwoSeven struct {
    VenvTwoSix
}

func (v *Venv) String() string {
    return fmt.Sprintf("Venv{%s, %s}", v.Source, v.Target)
}

func (v *Venv) Destroy() error {
    log.Println("Cleaning up", v.Target)
    return os.RemoveAll(v.Target)
}

type TransportableVenv interface {
    Construct() error
    Destroy() error
    Stream(io.Writer) error
    fmt.Stringer
}

func (v *Venv) Stream(outstream io.Writer) (err error) {
    log.Println("Streaming", v.TargetPath())
    err = util.TarGz(v.TargetPath(), outstream)
    return
}

func buildPath(root string, paths []string) string {
    if len(paths) > 0 {
        all := make([]string, len(paths) + 1)
        all[0] = root
        copy(all[1:], paths)
        root = path.Join(all...)
    }
    return root
}

func (v *Venv) SourcePath(paths ...string) string {
    return buildPath(v.Source, paths)
}

func (v *Venv) TargetPath(paths ...string) string {
    return buildPath(v.Target, paths)
}

func (v *VenvTwoSix) exec(funcs ...func() error) (err error) {
    for i := 0; err == nil && i < len(funcs); i++ {
        err = funcs[i]()
    }
    return err
}

func (v *VenvTwoSix) Construct() error {
    log.Println("Assembling 2.6 transportable at", v.TargetPath())
    return v.exec(
        v.HandleInclude,
        v.HandleLib,
        v.HandleMan,
        v.HandleBin)
}

func (v *VenvTwoSeven) Construct() error {
    log.Println("Assembling 2.7 transportable at", v.TargetPath())
    err := v.VenvTwoSix.Construct()
    v.HandleLocal()
    return err
}

func (v *VenvTwoSix) HandleInclude() error {
    log.Println("Handling include")
    return util.CpRecursive(
        v.SourcePath("include"),
        v.TargetPath("include"))
}

func (v *VenvTwoSix) HandleLib() error {
    log.Println("Handling lib")
    return util.CpRecursive(
        v.SourcePath("lib"),
        v.TargetPath("lib"))
}

func (v *VenvTwoSix) HandleMan() error {
    srcPath := v.SourcePath("man")
    _, err := os.Stat(srcPath)
    if err == nil {
        log.Println("Handling man")
        return util.CpRecursive(
            srcPath,
            v.TargetPath("man"))
    } else if os.IsNotExist(err) {
        log.Println("Skipping man")
        return nil
    }
    return err
}

func (v *VenvTwoSix) HandleBin() (err error) {
    log.Println("Handling bin")
    altDir := v.TargetPath("bin.original")
    err = util.CpRecursive(v.SourcePath("bin"), altDir)
    if err == nil {
        wrapper := v.TargetPath("exec-wrapper")
        err = exec.ConstructWrapper(wrapper, "bin.original")
        if err == nil {
            exec.ConstructWrapperLinks(v.TargetPath("bin"), altDir, wrapper)
        }
    }
    return
}

func (v *VenvTwoSeven) HandleLocal() (err error) {
    log.Println("Handling local")
    src := v.SourcePath("local")
    dest := v.TargetPath("local")
    var info os.FileInfo
    info, err = os.Stat(src)
    if err == nil {
        err = os.Mkdir(dest, info.Mode().Perm())
        if err == nil {
            entries, err := ioutil.ReadDir(src)
            if err == nil {
                for _, entry := range entries {
                    err = os.Symlink(
                        path.Join("..", entry.Name()),
                        path.Join(dest, entry.Name()))
                    if err != nil {
                        return err
                    }
                }
            }
        }
    }
    return err
}

func build(src string, dest string, verMaj int, verMin int) TransportableVenv {
    if verMaj == 2 && verMin == 6 {
        log.Println("Assembling 2.6 virtualenv wrapper.")
        x := new(VenvTwoSix)
        x.Source = src
        x.Target = dest
        return x
    } else if verMaj == 2 && verMin == 7 {
        log.Println("Assembling 2.7 virtualenv wrapper.")
        x := new(VenvTwoSeven)
        x.Source = src
        x.Target = dest
        return x
    }
    return nil
}

func WrapVenv(pth string) (TransportableVenv, error) {
    maj, min, err := DetermineVersion(pth)
    if err != nil {
        return nil, err
    }
    target, err := ioutil.TempDir("", "slingpie-")
    if err != nil {
        return nil, err
    }
    return build(pth, target, maj, min), nil
}

