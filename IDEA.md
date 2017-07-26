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



