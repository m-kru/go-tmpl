**tmpl** is a simple stdout text template emitter.
Templates are regular text files stored in the file system.
`tmpl` supports annotations in templates and uses the go `"text/template"` package logic to execute templates.
To get more information on how annotations work, check the documentation of `"text/template",` for example:

```
go doc template
```

`tmpl` has no external dependencies, it utilizes only go standard library.

# How it works

Templates are regular text files stored in the directory pointed to by the `TMPL_DIR` environment variable.
Within the `TMPL_DIR`, the user can nest directories as desired.
For example, let's assume `TMPL_DIR=/home/user/templates`.
The content of the directory can look as follows:
```
/home/user/templates/
├── absent-mail.txt
├── c
│   ├── do-while.c
│   ├── main.c
│   ├── Makefile
│   └── switch.c
├── go
│   ├── main.go
│   └── test.go
├── kind-mail.txt
├── tex
│   ├── enumerate.tex
│   ├── itemize.tex
│   └── letter.tex
└── vhdl
    ├── entity.vhd
    ├── process
    │   ├── async.vhd
    │   └── sync.vhd
    └── tb.vhd
```

To print any template to the stdout the user simply calls `tmpl` program with arguments being template path subdirectories and template file name.
For example, to print `absent-mail.txt` user calls:
```
tmpl absent-mail.txt
```
For example, to print `sync.vhd` user calls:
```
tmpl vhdl process sync.vhd
```

To make things easier, the provided template file name is a template file pattern.
`tmpl` looks for a file containing the provided pattern in the specified subdirectory.
If only one file within the subdirectory matches the pattern, then it is assumed to be the template the user asks for.
For example, to print `absent-mail.txt` the user can call:
```
tmpl absent
```
For example, to print `sync.vhd` user can simply call:
```
tmpl vhdl process s
```
If the provided template file pattern is ambiguous, `tmpl` reports an error.

The pattern-matching rule doesn't apply to subdirectory names.
They must be provided by the user exactly as they are.

## Annotations

`tmpl` supports annotations in templates and uses the go `"text/template"` package logic to execute templates.
To pass the annotation value, simply append the `<key>=<value>` argument while calling the program.
Annotation assignments are separated from the subdirectory names and template file pattern arguments with the `--` argument.

Example:

Let's assume we have the following template named `mail`:
```
Dear {{.name}},

...

Good luck, {{.name}}!
```
To generate the template with the `name` annotation replaced with a specific value, one can execute:
```
tmpl mail -- name=Tom
```
The generated text looks as follows:
```
Dear Tom,

...

Good luck, Tom!
```

## Integration with other programs

If you want to integrate `tmpl` with other programs, for example, with text editors, `tmpl` is capable of printing custom text to stderr.
The text for the stderr must be placed in a file with the same base name as the template but must start with the '.'.
Moreover, the file with stderr content must end with a custom extension.

### Example

Let's assume there is a file named `for.c` with the following simple for-loop template:
```
for () {
}
```
The [enix](https://github.com/m-kru/enix) editor can set the cursor in an arbitrary place after inserting the template.
To do so, it must be provided with the following commands via stderr:
```
enix:sel-switch-cursor
enix:5 right
```
These commands are placed in the `.for.enix` file in the same directory as the `for.c` file.

To print extra content to the stderr, not only template to the stdout, `tmpl` must be called with an additional argument.
This argument is the extension of the file with stderr content.
It must be placed as the first argument for `tmpl` and must be prepended with the '-' character.

In the above example, to print template content to the stdout and extra editor control content to the stderr, the user must run:
```
tmpl -enix c for
```

# Installation

1. Clone the repo, run `go build`, copy `tmpl` binary wherever you want.
2. `go install github.com/m-kru/tmpl@latest`
