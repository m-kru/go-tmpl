**tmpl** is a simple stdout text template emitter.
Templates are regular text files stored in the file system.
`tmpl` supports annotations in templates and uses the go `"text/template"` package logic to execute templates.
To get more information on how annotations work, check the documentation of `"text/template",` for example

`tmpl` has no external dependencies, it utilizes only go standard library.

```
go doc template
```

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

## Setting cursor position

`tmpl` supports settings cursor position.
This feature might be useful if the user would like to call `tmpl` from within a text editor.
To specify cursor position, add the following line at the beginning of a template file.
```
cursor@<line>:<column>
```
The cursor position line is filtered out and not printed to the stdout.
Information about line number and column is redirected to the stderr.
However, the `cursor@` prefix is removed.
It is easy to parse cursor position information in the stderr stream, as except from the cursor position, only error messages are printed to the stderr.
However, error messages always start with the `error:` prefix.
For example, let's consider the following template:
```
cursor@2:9
\begin{itemize}
  \item
\end{itemize}
```
The first line informs that the cursor should be placed in line number 2, column number 9.
The output from the:
```
tmpl tex item
```
looks as follows:
```
\begin{itemize}
  \item
\end{itemize}
2:9 # <- This line is printed to the stderr, without this comment.
```
The cursor setting feature, of course, requires support from the editor.
I am open to modifying this feature if you know how to handle it better.

# Installation

1. Clone the repo, run `go build`, copy `tmpl` binary wherever you want.
2. `go install github.com/m-kru/tmpl@latest`
