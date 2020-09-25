# Duster

This is Duster. A small debugger for unikernels written in C that run on the xen hypervisor (PV only at the moment). 

## Installation 
Installing Duster should be easy. 

### Installation Ubuntu 
Download the `deb` file from release section of the repo. Then
    sudo dpkg -i duster.deb

### Building and Installing from source 
It is recomeneded you use the Dockerfile provided as it handles all the dependencies required. However, if you don't want to use docker you may just install the dependencies manually. Anyway to build with docker just carry out the following commands:
	sudo docker build -t duster-build .
	sudo docker run -v $PWD:/build:Z -it duster-build

Then just move the duster executable into your bin and you're done!

## Unikernel build requirements
When using Duster please make sure to turn off optimisations. Unfortunately, Duster can't handle these at the moment. Generally, setting the optimisation to `-O0`, disabling loop unrolling (`-fno-unroll-loops`) and making sure `-fomit-frame-pointer` isn't enable should be enough. If you're using Stardust then this has already been done for you. Just run the Makefile with `debug=y`.

## Start up 
Starting Duster is pretty straightforward. Before you start Duster you need to start your domain up. The domain needs to be paused on startup. To do this just run the command `xl create -p [domain name]`. Then to start Duster just do:
    duster -path=[Path to domain.gz] -id=[the domain id (e.g. 5)]

## Using Software 
The commands supported by Duster are:
1. break [filename.c]:[line number] - sets a breakpoint at specific line in the c program
2. remove [filenae.c]:[line number] - deletes a breakpoint
3. continue - runs until it hits a breakpoint or runs forever if there is no breakpoint.
4. read [variable name]- reads a variable (this should be compatible with C type. However, there slight issue with arrays of the form c[variable] which causes it crash).
5. quit - quits the debugger.
6. step - steps to the next source line
7. der [variable] - deferences a pointer (only works with variable not attributes, unfortunately)

## Demo 
The the following demo should help to clarify the above section. Assume the following code is being debugged after the initial startup.
### Code 
```
void demo() {
        struct test my_test;
        my_test.val = 10;
        my_test.my_pointer = NULL;
        my_test.no = 0.5;
        for (int i = 0; i < 3; i++) {
                printf("We are on the %d iteration of the loop\n", i);
                my_test.val += 100;
                my_test.no += 0.5;      
        }
        struct test *pointer_tester;
        pointer_tester = malloc(sizeof(struct test));
        pointer_tester->val = 20;
        pointer_tester->no = 0.110;
        pointer_tester->my_pointer = NULL;
        printf("%d\n", pointer_tester->val);
        printk("We're done\n");
}     
```
### Duster running
![Demo](images/Screencast-from-25-09-20-19_15_40.gif)

## Current Limitations
Unfortunately Duster is not prefect! The following details the issues you may have with Duster. These will be fixed eventually. If you find or think of anything else please put it in an issue and it will be looked at. 
* Currently, you cannot print off arrays whose length is defined by a variable
* No C++ (or any other language) support 
* When you quit you need to reset the domain. Duster does not clean up after itself!
* You cannot view contents of pointer type attributes 
* The step command will just step into a function. If you want to step over a function you must set a breakpoint after it
* You need to run the quit command, depending on your OS ctrl-C won't always work
