package util

import (
    "io"
    "os"
    "path/filepath"
)

func cpLink(srcLink, destLink string) (err error) {
    var target string
    target, err = os.Readlink(srcLink)
    if err == nil {
        err = os.Symlink(target, destLink)
    }
    return
}

func cpFile(srcFile, destFile string, info os.FileInfo) error {
    var in, out *os.File
    var err error
    in, err = os.Open(srcFile)
    if err != nil {
        return err
    }

    defer in.Close()

    out, err = os.Create(destFile)
    if err != nil {
        return err
    }
    
    _, err = io.Copy(out, in)
    if err != nil {
        return err
    }

    err = out.Close()
    if err != nil {
        return err
    }

    return os.Chmod(destFile, info.Mode().Perm())
}

func cpWalker(src, dest string) func(string, os.FileInfo, error) error {
    return func(path string, info os.FileInfo, err error) error {
        target := dest + path[len(src):]
        if info.IsDir() {
            return os.Mkdir(target, info.Mode().Perm())
        } else if info.Mode() & os.ModeSymlink > 0 {
            return cpLink(path, target)
        } else if filepath.Ext(path) != ".pyc" {
            return cpFile(path, target, info)
        }
        return nil
    }
}

func CpRecursive(src, dest string) (err error) {
    src, err = filepath.Abs(src)
    if err == nil {
        dest, err = filepath.Abs(dest)
    }
    if err == nil {
        return filepath.Walk(src, cpWalker(src, dest))
    }
    return
}

