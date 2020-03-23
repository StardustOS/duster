
all: test variable_data simple globalvars different-scopes

test: test.c
	gcc -g -O0 test.c -o test

variable_data: variable_data.c
	gcc -g -O0 variable_data.c -o variable_data

simple: simple.c
	gcc -g -O0 simple.c -o simple

globalvars: globalvars.c
	gcc -g -O0 globalvars.c -o globalvars

different-scopes: different-scopes.c
	gcc -g -O0 different-scopes.c -o different-scopes


clean:
	rm test
	rm variable_data
	rm simple
	rm different-scopes