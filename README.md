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
7. def [variable] - deferences a pointer (only works with variable not attribute, unfortunately)

 
