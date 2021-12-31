# goutils

A `go` module aggregating various utilities that are used across my go applications

### listflag

Some of my `go` executables accept an argument that is specified multiple times to form a list (eg. `--platform ps3 --platform wii` will result in a list of strings `["ps3", "wii"]`). Since I'm not a fan of inventing cutom syntaxes (comma separated strings or `=` assignment or whatever else), as I always forget such custom solutions to the list issue, I did a very simple util type to be used with `flag` go package to use it across my applications that can handle multiple specifications of the same named argument

### cil - Controlled Infinite Loop

An infinite loop that can be controlled by with callback or context. Useful for calling something infinitely until a certain condition is met or an error raises that should cancel the loop.
