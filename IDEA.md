Planning
===============
Actions:
* init
* new
* create


#### init
```
boiler init
```
Initialize boiler plate in current dir, creates a .boiler directory
analyses vars in boiler files and automatically creates a config.yaml


#### new
```
boiler new [repo/folder] [target]
```
Generates files from boilerplate and filters trough template, makes questions to user
and stores in user.yaml to read later in create


#### create
```
boiler create gofile [name]
```
Creates a file from template/structure passing through user.yaml defaults

#### TODO:

* create a CoreContext from a folder either found root, or new template struct
* Support multi files in generators

```yaml
generators:
  go:
    desc: Creates a go file based on template
    target: "{{.curdir}}/{{.name}}"
    source: gofile.go.boiler
    ext: go # for files only
    flags: ["package, p","proj"]
  controller:
    desc: Creates a controller for app
    target: "{{.projroot}}/{{.name}}"
    source: controller/
```

To something:
```yaml
generators:
  go:
		files:
			- [gofile.go.boiler, "{{.curdir}}/{{.name}}"] // source, dest // Ext taken from src file
    desc: Creates a go file based on template
    flags: ["package, p","proj"]
  controller:
    desc: Creates a controller for app
    target: "{{.projroot}}/{{.name}}"
    source: controller/

```
* Improve Process function to process a file and then a dir(invert) - DONE


