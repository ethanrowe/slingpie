package util

import (
    "compress/gzip"
    "io"
    "os"
    "path/filepath"
    "archive/tar"
)

const TARBUFFSIZE = 4096 // 4kb

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

func tarWalker(src string, w *tar.Writer) func(string, os.FileInfo, error) error {
    return func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        var link string
        if info.Mode() & os.ModeSymlink > 0 {
            link, err = os.Readlink(path)
        }
        if err == nil {
            var head *tar.Header
            head, err = tar.FileInfoHeader(info, link)
            if err == nil {
                head.Name, err = filepath.Rel(src, path)
                if err == nil {
                    err = w.WriteHeader(head)
                    if err == nil && info.Mode().IsRegular() {
                        var in *os.File
                        in, err = os.Open(path)
                        if err == nil {
                            var bytes [TARBUFFSIZE]byte
                            var count int
                            count, err = in.Read(bytes[:])
                            for err == nil {
                                _, err = w.Write(bytes[:count])
                                if err == nil {
                                    count, err = in.Read(bytes[:])
                                }
                            }
                            if err == io.EOF {
                                err = nil
                                in.Close()
                            }
                        }
                    }
                }
            }
        }
        return err
    }
}

func Tar(src string, w io.Writer) (err error) {
    src, err = filepath.Abs(src)
    if err == nil {
        out := tar.NewWriter(w)
        err = filepath.Walk(src, tarWalker(src, out))
        if err == nil {
            err = out.Close()
        }
    }
    return
}

func TarGz(src string, w io.Writer) (err error) {
    zipper := gzip.NewWriter(w)
    err = Tar(src, zipper)
    if err == nil {
        err = zipper.Close()
    }
    return
}

