# A brief user guide

This is a pretty basic version of the debugger. The code still needs to be cleaned up. However, most of the functionality has been implemented.

## Changes to Stardust
To use the debugger, all the optimisations must be turned off. This is done by going into the stardust.mk file and commenting out all the optimisation. Then change -O3 to -O0 and also add -g (so the debug symbols are added to the executable). Additionally, in the mini-os.conf file change the number of vcpu to one (multi-core also messes up the debugger).

## Ruuning stardust
Once the changes above have been made, then run stardust using the following command: 

    sudo xl create -p mini-os.conf

This will startup stardust and pause it. This enables use to debug from the first line of c. 

## Running duster
    ./duster -path=[Path to mini-os.gz] -id=[the domain id (e.g. 5)]

## Commands 
The commands supported by the program are as follow:
    
1. break [filename.c]:[line number] - sets a breakpoint at specific line in the c program
2. remove [filenae.c]:[line number] - deletes a breakpoint
3. continue - runs until it hits a breakpoint or runs forever if there is no breakpoint.
4. read [variable name]- reads a variable (this should be compatible with C type. However, there slight issue with arrays of the form c[variable] which causes it crash).
5. quit - quits the debugger.
6. step - steps to the next source line
7. def [variable] - deferences a pointer (only works with variable not attribute, unfortunately)

(The executable can be found in the bin folder. I still need to add a build system).
 