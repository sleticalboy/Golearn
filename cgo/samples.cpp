
#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <thread>
#include <iostream>

extern "C" {

#include "samples.h"

float hello(const char *str) {
    std::cout << "hello() thread id: " << std::this_thread::get_id() << std::endl;

    printf("hello() call with: %s\n", str);
    return 20.04 * 10.5;
}

void crash() {
    printf("%s() will crash after 1 seconds\n", __func__);
    int *p = NULL;
    *p = 0; // 引发一个空指针异常
}

// 回调函数，在 go 中实现
extern void cgoCallback(int value);

void start_loop(int counter) {
    std::cout << "start_loop() thread id: " << std::this_thread::get_id() << std::endl;
    if (counter == 0) return;

    auto task = [&]() {
        std::cout << "start_loop() sub thread id: " << std::this_thread::get_id() << std::endl;
        for(int i = 0; i < counter; i++) {
            // std::cout << "loop() counter: " << i << std::endl;
            cgoCallback(i);
            usleep(450 * 1000);
        }
    };
    auto t = std::thread(task);
    t.join();
}

}
