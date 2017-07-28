Creating a boilerplate with boiler
---------------------------------------

if not, set $GOPATH/bin in $PATH to be easier to use any cli tool installed from go  
  
Getting the tool:
```
go get github.com/gohxs/boiler/cli/bp
```

Creating a boiler plate structure
```
mkdir my-boilerplate
mkdir my-boilerplate/.boiler
```

Add the following to my-boilerplate/.boiler/config.yml
```yaml
description: My plate project
vars:  # Optionally request information to user while initializing
  - {name: author, question: What is your name}                           # Random vars for example purpose
  - {name: thing, question: What is the thing name, default: beepboop}            
```

Add a simple file to be processed trough templates:
file: my-boilerplate/hello.sh.boiler
```bash
echo "hello {{.author}}, I'm {{.thing}}"
```

Tree should be something like this:
```
my-boilerplate
├── .boiler
│   └── config.yml
└── hello.sh.boiler
```


Consuming the boilerplate:

```bash
bp create location/my-boilerplate proj1
```
where location should be the folder where your boilerplate is  

Result should look like:
```
Loading boilerplate from my-boilerplate/
My plate project
-----
What is your name [author] ()? me
What is the thing name [thing] (beepboop)?
Generating project...
Created project: proj1 proj1
```

tree:
```
proj1
├── .boiler
│   ├── config.yml
│   └── user.yml
└── hello.sh
```

Executing hello.sh and checking results
```
$ . proj1/hello.sh
Hello me I'm beepboop
```

###Adding a generator:
in config.yml we add the generator field with the list of generators,
Generators can be added/edited in on going projects by changing proper .boiler/config.yml in each project
We can even create .boiler/config.yml in existing projects with required generators/templates

Create the templates folder
```
mkdir my-template/.boiler/templates
```  

We will add a template file inside our my-boilerplate
**my-template/.boiler/templates/controller.sh.boiler**
```
# copyright (C) {{.author}} {{.time.Format "2006"}}

echo "I'm a controller I do {{.action}}"
```

lets add a generator named `controller`
```yaml
description: My plate project
vars:  # Optionally request information to user while initializing
  - {name: author, question: What is your name}  
  - {name: thing, question: What is the thing name, default: beepboop}            # Random vars for example purpose
# New generators field here:
generators:
  controller:                                                           # A generator named controller
    aliases: [c]                                                        # It will work as 'bp add controller' or 'bp add c'
    files:                                                              # File list one or several files, can also be a directory
      - {source: controller.sh.boiler, target: "{{.name}}"}             # It will look for controller.sh.boiler inside .boiler/templates folder
    vars:                                                               # vars list
      - {name: action, question: What is the action, flag: "action,a" }  # Similar to init vars but here we optionally add flag field
```

After generating other project with the new config.yml we can see descriptions from newly created controller generator
```
$ bp create location/my-boilerplate proj2
$ cd proj 2
$ bp add
Add a file based on boilerplate generator

Usage:
  bp add [file] [flags]
  bp add [command]

Aliases:
  add, a

Available Commands:
  controller

Flags:
  -h, --help   help for add

Use "bp add [command] --help" for more information about a command.

$ bp add controller -h
Usage:
  bp add controller [name] [flags]

Aliases:
  controller, c

Flags:
  -a, --action string   What is the action
  -h, --help            help for controller
```  


Lets generate a new unit within proj2 folder:
```
bp add controller mycontroller1
What's the action [action] ()? things
Generating file: mycontroller1.sh
```



Can also be generated with flags, which bp will not ask if flags are present:
```
$ bp add c mycontroller2 --action "other things"
Generating file: mycontroller2.sh
```


Results:
```
$ . mycontroller1.sh
I'm a controller I do things
$ . mycontroller2.sh
I'm a controller I do other things
```



Generated controller content:
```
# copyright (C) me 2017

echo "I'm a controller I do things"
```


Generator can have many files as it requires, target will be processed through template
so if we add generator with a file that containts a target like `"controllers/{{.author}}{{.name}}"` and using
`bp add c control` it will generate that file within `controllers/authorcontrol.sh`  
  
As a special case we can add a dot alias to this generator i.e. '.sh' this way we can pass a single parameter to bp
`bp add file.sh` it will be the same as `bp add controller file.sh`






