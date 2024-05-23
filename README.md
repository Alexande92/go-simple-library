# Simple Go library 

Simple project for managing your books


## Actions

To run, build or test this project provided Makefile

Command **run** creates binary file in `dist` folder and runs it
 ```bash
    make run
 ```

Command **build** - just makes a binary
 ```bash
    make build
 ```

Command **test** - should run tests from files with `_test` in name 
All results will be shown with coverage per package. Required coverage per package ~80%
 ```bash
    make test
 ```