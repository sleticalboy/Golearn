import ctypes
import time

so = ctypes.cdll.LoadLibrary('./a.so')

start = time.time()
for i in range(1000000):
    so.Fib(199)
end = time.time()
print(f'python cost: {end - start}s')

so.GoFib()

