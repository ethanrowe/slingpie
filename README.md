# slingpie

A little helper for making python virtualenvs more easily relocatable

# Summary

Python virtualenvs (at least in the 2.x space; cannot speak to the 3.x
space just yet) are handy and all, but the manner of relocating them
is rather messy.  This is due to the writing of absolute paths to the
underlying python executable within all entry points.

We sling python virtualenvs around as one of the artifacts of our build
process, as this makes dealing with certain projects (like numpy, scipy,
etc.) easier.  We can also push a whole tarred virtualenv through a hadoop
job and have it work nicely in a streaming task.

The entry point thing makes this a little less reliable.  Enter *slingpie*,
the greatest project in the go language ever written by this particular
human being as of this moment.

Slingpie takes a virtualenv on your filesystem and produces a gzipped tarball
of it that is naturally relocatable.  It wraps all of the entry points to
ensure the virtualenv is activated, using relative symlinks to accomplish
this.

The relocatable virtualenv should be regarded as immutable; you shouldn't
`pip install` into it, as things won't work properly.  But you can move it
around, tar it up and plop it elsewhere, etc.

# Usage

Prepare your virtualenv:

        % virtualenv your/python

Make a tarball:

        % slingpie your/python > your-python.tar.gz
        2013/11/21 12:03:58 Finding python version for /home/you/your/python
        2013/11/21 12:03:58 Assembling 2.7 virtualenv wrapper.
        2013/11/21 12:03:58 Assembling 2.7 transportable at /tmp/slingpie-249937419
        2013/11/21 12:03:58 Assembling 2.6 transportable at /tmp/slingpie-249937419
        2013/11/21 12:03:58 Handling include
        2013/11/21 12:03:58 Handling lib
        2013/11/21 12:03:58 Handling man
        2013/11/21 12:03:58 Handling bin
        2013/11/21 12:03:58 Handling local
        2013/11/21 12:03:58 Streaming /tmp/slingpie-249937419
        %

The tarball does *not* have a top-level directory.  It's up to you to create
the directory you want to use as the transportable/relocated virtualenv.

Thus:

        % mkdir /mnt/somewhere/my-python
        % cd /mnt/somewhere/my-python
        % tar xzf ~/your-python.tar.gz

Now you can use it:

        % /mnt/somewhere/my-python/bin/python

If you import things and inspect their help information within an interactive
session, you'll see that all the files are coming from the relocated paths.

# Why go?

A reasonable question is: why on earth did you implement this in go?
It's a python-oriented utility, so why not implement it in python?  If not
python, why not something like bash?

I seriously considered these things.  The problem is that I didn't want to
distribute a python package that would potentially involve installing into
virtualenvs in order to have the software for bundling up those virtualenvs.

I also didn't want to use bash because I felt dirty about it, which is
completely irrational.

I also wanted to try go on something small, and come out the other side with
something portable.  So here we are.

