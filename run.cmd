@echo off

go build -o temp.exe
temp.exe su --keep
rm temp.exe