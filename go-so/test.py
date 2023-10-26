import ctypes
import time

so = ctypes.cdll.LoadLibrary('./libfib.so')
i_ = 10

start = time.time()
r = so.Fib(i_)
end = time.time()
print(f'go fib {i_} -> {r}, cost: {end - start}s')


def fib(value: int) -> int:
    if value == 1 or value == 2:
        return value
    return fib(value - 1) + fib(value - 2)


start = time.time()
r = fib(i_)
end = time.time()
print(f'python fib {i_} -> {r}, cost: {end - start}s')
