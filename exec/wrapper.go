package exec

import "fmt"

const wrapperText = `#!/bin/bash
#!/bin/bash

INVOCATION_PATH=$0

if [ "${INVOCATION_PATH:0:1}" != '/' ]; then
    INVOCATION_PATH=$(pwd)/$INVOCATION_PATH
fi

TARGET_NAME=$(basename "$INVOCATION_PATH")
VIRTUALENV_BIN=$(dirname "$INVOCATION_PATH")
VIRTUALENV_ROOT=$(dirname "$VIRTUALENV_BIN")
VIRTUALENV_BIN_REAL="$VIRTUALENV_ROOT/%s"
VIRTUALENV_PYTHON="$VIRTUALENV_BIN_REAL/python"
TARGET="$VIRTUALENV_BIN_REAL/$TARGET_NAME"

PY_SHEBANG=$(head -n 1 $TARGET | grep -E '^#!.*\bpython([0-9]\.[0-9])?\b')

if [[ -z $PY_SHEBANG ]]; then
    PREFIX=""
else
    PREFIX="$VIRTUALENV_PYTHON "
fi


if [[ ! $VIRTUALENV_PYTHON -ef $VIRTUALENV_ACTIVE_PYTHON ]]; then
    if [ -n "$_OLD_VIRTUAL_PATH" ]; then
        PATH="$_OLD_VIRTUAL_PATH"
        unset _OLD_VIRTUAL_PATH
    fi

    VIRTUALENV_ACTIVE_PYTHON="$VIRTUALENV_PYTHON"
    _OLD_VIRTUAL_PATH="$PATH"
    VIRTUAL_ENV="$VIRTUALENV_ROOT"

    PATH="$VIRTUALENV_BIN:$PATH"
    export _OLD_VIRTUAL_PATH
    export VIRTUALENV_ACTIVE_PYTHON
    export VIRTUAL_ENV
    export PATH
fi

exec $PREFIX$TARGET "$@"
`

func WrapperBody(realbin string) string {
    return fmt.Sprintf(wrapperText, realbin)
}
