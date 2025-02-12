from __future__ import print_function
from ctypes import *

lib = cdll.LoadLibrary("./monkeyc.so")


# define class GoString to map:
# C type struct { const char *p; GoInt n; }
class GoString(Structure):
    _fields_ = [("p", c_char_p), ("n", c_longlong)]


# describe and invoke Add()
lib.CompileAndRun.argtypes = [GoString]
lib.CompileAndRun.restype = GoString

code = b"let add = fn(a, b) { a + b; }; add(1, 2);"
nativecode = GoString(code, len(code))

res = lib.CompileAndRun(nativecode)
# p_field = res._fields_[0][1].value
# print(type(p_field), p_field.__get__(res._fields_[0]))

print("compiled: ", getattr(res, "p"))
