使用 cgo mingw 生成动态库 并使用 下面操作生成导入库

# 1. Build MinGW DLL.
go build -buildmode=c-shared -o foo.dll . # Generate MinGW DLL
# 2. Generate a def file out of the DLL.
gendef foo.dll
# 3. Generate the MSVC import library.
dlltool --input-def foo.def --output-lib foo.lib 
# 4. Compile the C file using MSVC with the import library previously generated.
# There is no need for foo.dll to be present when compiling C,
# but it must be available when running the resulting binary.
cl.exe main.c foo.lib 
