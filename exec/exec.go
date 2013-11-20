package exec

import (
    "io"
    "io/ioutil"
    "os"
    "path"
)

func ConstructWrapper(dest, wrappedDir string) (err error) {
    var file *os.File
    file, err = os.Create(dest)
    if err != nil {
        return
    }

    _, err = io.WriteString(file, WrapperBody(wrappedDir))
    if err != nil {
        return
    }

    err = file.Close()
    if err != nil {
        return
    }

    return os.Chmod(dest, 0755)
}

func ConstructWrapperLinks(target, real, wrapper string) (err error) {
    var info os.FileInfo
    info, err = os.Stat(real)
    if err != nil {
        return
    }
    err = os.Mkdir(target, info.Mode().Perm())
    if err != nil {
        return
    }

    var entries []os.FileInfo
    entries, err = ioutil.ReadDir(real)
    if err != nil {
        return
    }

    for _, entry := range entries {
        err = os.Symlink(
            path.Join("..", path.Base(real), entry.Name()),
            path.Join(target, entry.Name()))
        if err != nil {
            return
        }
    }
    return
}

