package main

import (
    "fmt"
    "log"
    "os"
    "github.com/ethanrowe/slingpie/venv"
)

func showHelp() {
    fmt.Printf(`Usage: slingpie path-to-virtualenv > your.tar.gz

Slingpie version v%s

Streams to stdout a gzipped tar archive of the specified python
virtualenv, the contents of which having been wrapped in such a manner
that they should work properly when placed at any arbitrary filesystem
location.

The contents of the resulting tar archive do not include a top-level
directory; it is the responsibility of the user to unpack the archive
in an appropriate location.

Once unpacked, the transportable virtualenv can be used like most any
python virtualenv, except that the python packages should not be
changed (no "pip install", for instance).  The virtualenv can be run on
systems that have equivalent architectures, with the same version of
python installed to the system in the same location.

The transportable virtualenv does not need to be activated to work properly.

Example usage:

  # Create a virtualenv
  virtualenv some-python
  # Pack it up
  slingpie some-python > some-python.tar.gz
  # Put the transportable one at "transportable"
  mkdir transportable
  (cd transportable; tar xzf ../some-python.tar.gz)
  # Invoke the python in your transportable one.
  transportable/bin/python
`, releaseVersion)
}

func main() {
    args := os.Args[1:]
    if len(args) == 1 && args[0] != "-?" && args[0] != "--help"  {
        src, err := venv.WrapVenv(args[0])
        if err == nil {
            err = src.Construct()
            if err == nil {
                err = src.Stream(os.Stdout)
                // src.Destroy()
            }
        }
        if err != nil {
            log.Fatalln("An error occurred:", err)
        }
    } else {
        showHelp()
    }
}

