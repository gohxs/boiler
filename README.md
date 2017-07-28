Boiler (name WIP)
=======================
Simple tool written in go to assist project creation from boilerplates with go text/templates

This tool can help in the create of new projects with a template project structure 
and also 

Follow this [walkthrough](/WALKTHROUGH.md) to create a simple boilerplate

Content
--------------
* [Install](#install)
* [bash completion](#bash-completion)
* [boilerplate structure](#boilerplate-structure)
	* [.boiler/config.yml](#config-file)
* [Usage](#usage)
  * [create command](cmd-create)
  * [add command](cmd-add)


<a name="install"></a>
Install:
----------------
```
go get github.com/gohxs/boiler/cli/bp
```

Basic usage:
---------------
```
bp create boilerplate/path proj1
```

```
bp create http://github.com/gohxs/boiler-example-boilerplate proj1
```

<a name="bash-completion"></a>
Bash autocompletion
--------------------
generating bash completion file
```
bp --bashcompletion=bp.sh
```
place result file `bp.sh` in ~/.bash_completion or any bash_completion.d folder in OS
> bp app sould be in $PATH

<a name="boilerplate-structure"></a>
boilerplate structure example:
-------------------------------
```
boiler-example
├── .boiler
│   ├── config.yml
│   └── templates
│       ├── controller
│       │   ├── pkg.go.boiler
│       │   └── pkg_test.go.boiler
│       ├── gofile.go.boiler
│       └── pure.txt
├── hello.go.boiler
├── README.md.boiler
└── staticfile.txt
```
<a name="config-file"></a>
####File: .boiler/config.yml
```yaml
description: |  ## Descriptive text to show when creating a project from this boileplate
  ______       ___________            
  ___  /__________(_)__  /____________
  __  __ \  __ \_  /__  /_  _ \_  ___/
  _  /_/ / /_/ /  / _  / /  __/  /    
  /_.___/\____//_/  /_/  \___//_/.
  Test boilerplate for boiler cli app

vars: #list of variables for initialization
  - {name: author , default: No author, question: Ask something about author}
  - {name: description, default: Awesome app, question: Ask something about description}

generators:
  gofile:
    aliases: [.go]
    description: Creates a go file based on template                # Description for help page
    files:                                                          # Files to be processed, this files are in .boiler/template/... folder
      - {source: gofile.go.boiler, target: "{{.curdir}}/{{.name}}"} # target is processed with template
    vars:                                                           # List of variables for unit creation
      - name: package
        default: "{{.projName}}"                                    # default is processed with template
        flag: "package, p"                                          # flags that can be used in bp add gofile --package main
        question: package name of the new file                      # question to be shown on interactive input
        description: Package to place on package thing              # description of the var
  puretxt:
    description: Clones pure.txt file
    aliases: [.txt]
    files:
      - {source: pure.txt, target: "{{.curdir}}/{{.name}}" }
  package:
    description: Creates a package for app
    files:
      - {source: controller/pkg.go.boiler, target: "{{.projRoot}}/{{.name}}/{{.name}}"}
      - {source: controller/pkg_test.go.boiler, target: "{{.projRoot}}/{{.name}}/{{.name}}_test"}
``` 

config.yml root object

* `description` description of the overall boilerplate usually show on creation
* `vars` variables to be used in further templates
* `generators` list of *generator* objects 
  
**var** object - can be used in init or inside a generator: 

* `name` Name of variable, this name will be subsituted in {{.*varname*}}
* `default` Default value for variable
* `flag` Flag in command line (currently for generator only)
* `question` Text to display when running interactively
  
**generator** object describes a unit creation

* `description` - description of generator, generally used in cli help 
* `aliases`     - Array containing aliases while calling `bp add [name/aliases]`
* `files`       - Array of File objects {source, target} files used as source to be copied to destination
* `vars`        - vars describing either flag/user input these will be used in templates
	


---------------------
<a name="usage"></a>
<a name="cmd-create"></a>
####create command
```
Create new project from a boilerplate

Usage:
  bp create [repository/source] [projname] [flags]

Aliases:
  create, c

Flags:
  -h, --help   help for create
```

Considering source boilerplate example:
_.boiler_ folder contains the config.yml and templates for unit creations, 
its not necessary to have this folder to create a project  

all \*.boiler files outside .boiler folder will be processed through templates and extension will be removed
in this example a `hello.go.boiler` will become `hello.go` but `gofile.go.boiler` inside `.boiler` folder will remain the same

```
$ bp create ./boiler-example proj
Loading boilerplate from boiler-example
______       ___________
___  /__________(_)__  /____________
__  __ \  __ \_  /__  /_  _ \_  ___/
_  /_/ / /_/ /  / _  / /  __/  /
/_.___/\____//_/  /_/  \___//_/.
Test boilerplate for boiler cli app

-----
Ask something about author [author] (No author)? Luis Figueiredo
Ask something about description [description] (Awesome app)?
Generating project...
Created project: proj proj

```
Result tree 
```
proj
├── .boiler
│   ├── config.yml
│   ├── templates
│   │   ├── controller
│   │   │   ├── pkg.go.boiler
│   │   │   └── pkg_test.go.boiler
│   │   ├── gofile.go.boiler
│   │   └── pure.txt
│   └── user.yml
├── hello.go
├── README.md
└── staticfile.txt
```
`user.yml` was added with initialization vars prompted to user and will be used in unit generators


Analysing README.md.boiler and checking with result:
```
Project {{.projName}}
===================
{{.projDate.Format "02-01-2006"}}

Author: {{.author}}

Project variables
=================
{{range $k, $v := .}}{{$k}} = {{$v}} 
{{end}}
```
Result:
```
Project proj
===================
28-07-2017

Author: Luis Figueiredo

Project variables
=================
author = Luis Figueiredo 
description = Awesome app 
projDate = 2017-07-28 18:57:07.95796051 +0000 UTC 
projName = proj 
```

--------------
<a name="cmd-add"></a>
####add command
```
Add a file based on boilerplate generator

Usage:
  bp add [file] [flags]
  bp add [command]

Aliases:
  add, a

Available Commands:
  gofile      
  package     
  puretxt     

Flags:
  -h, --help   help for add

Use "bp add [command] --help" for more information about a command.
```
**Available commands** are specific for each project, it will list the generators available on the project

While inside our new created project 
```bash
bp add gofile other.go
```
Considering that we have vars in the generator `gofile`, these will be requested or fetched by flags from user:

Output:
```
package name of the new file [package] (proj)?
Generating file: .../proj/other.go
```

Using flags:
```bash
bp add gofile other.go --package proj
```

**Special case .ext**: With a special alias in the generator we can create files based on extension
```yaml
...
gofile:
    aliases: [.go]
    description: Creates a go file based on template
...
```
basically using the command as `bp add myfile.go` is the same as `bp add gofile myfile.go` bp checks the extension of desired file `.go`
and fetches a generator with that name/alias



###TODO:
- [X] Improve package naming of core
- [-] Maybe merge config package into (core) package? 
- [ ] Create `init` command which initializes a boiler folder in current work dir
  - [ ] Create commands to add variables to initialization
- [ ] Facilitate way to Add new generators from command instead of editing config.yml `bp generator create name`
  - [ ] Add files to generator `bp generator file sourcefile.ext
  - [ ] Add vars to generator `bp generator var name --question "" --default ""`
- [ ] Add more test cases within boiler (currently there is tests on cli)
- [ ] Remove global core from boiler and use it on cmd/main
  

