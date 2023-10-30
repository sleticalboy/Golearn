import ctypes
import time

so = ctypes.cdll.LoadLibrary('./libfib.so')

so.HelloGo(b"Python 3")

i_ = 93

start = time.time()
r = so.Fib(i_)
end = time.time()
print(f'go fib {i_} -> {r}, type: {type(r)}, cost: {end - start}s')


def fib(value: int) -> int:
    if value < 2:
        return 1
    f = 1
    s = 1
    for i in range(value):
        if i < 2:
            continue
        temp = s
        s = f + s
        if i > 90:
            print(f"{i} -> {s}")
        f = temp
    return s


start = time.time()
r = fib(i_)
end = time.time()
print(f'python fib {i_} -> {r}, type: {type(r)}, cost: {end - start}s')


class PyPerson(ctypes.Structure):
    _fields_ = [
        ('name', ctypes.c_char_p),
        ('age', ctypes.c_int),
        ('home', ctypes.c_char_p),
    ]


go_home = so.GoHome
p = PyPerson(b'li ming', 28, b'chain henan')
# 指针类型
so.GoHome(ctypes.byref(p))
# 非指针类型
so.GoHome2(p)
