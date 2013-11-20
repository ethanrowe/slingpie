package venv

import (
    "log"
    "os/exec"
    "path"
    "strconv"
    "strings"
)

const pythonPrefix = "Python "
const versionSeparator = "."

func DetermineVersion(p string) (maj, min int, err error) {
    log.Println("Finding python version for", p)
    cmd := exec.Command(path.Join(p, "bin", "python"), "-V")
    var out []byte
    out, err = cmd.CombinedOutput()
    if err == nil {
        outStr := string(out)
        if strings.HasPrefix(outStr, pythonPrefix) {
            vers := strings.Split(outStr[7:], versionSeparator)
            maj, err = strconv.Atoi(vers[0])
            if err == nil {
                min, err = strconv.Atoi(vers[1])
            }
        }
    }
    return
}
