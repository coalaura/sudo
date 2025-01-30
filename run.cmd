@echo off

go build -o temp.exe
temp.exe su
rm temp.exe