@echo off

echo Building...
go build -o temp.exe

echo Running...
temp.exe who

echo Done.
del temp.exe