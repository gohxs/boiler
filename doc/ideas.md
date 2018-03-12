#### New method:
Build project by going uptree and reading .boiler/boiler.yml


#### Automatic templating
parse files to find any keys `{{.var}}` it would ask use var name

#### Automatic generators
Use fs to find generators: 
```
boiler/
└── generators
    └── [genname]
        ├── boiler.yml
        ├── file2.go
        └── file.go
```

```bash
bp add genname target
```
It would find the generator by [genname] and process all files in that folder

##### Cons
* we would miss the ability to template target file names
##### Pros
* Easy starting point configuration, just add the folder and go


#### Global/Per user boiler templating
have a $HOME/.boiler, if project not found we could fetch generators from that folder
