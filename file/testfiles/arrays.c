#include <stdio.h>

int main(void) {
    char string[13];
    for (int i = 0; i < 13; i++) {
        string[i] = 'c';
    }
    printf("%s\n", string);
}