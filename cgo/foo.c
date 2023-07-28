
#include <stdio.h>
#include "foo.h"

float foo_func(const char *str) {
    printf("foo_func() call with: %s\n", str);
    return 20.04 * 10.5;
}