@echo off

REM chcp 65001

echo Hello World!
echo %0
echo %1
echo %2

REM choice /t 10 /c YN /d Y
REM Wscript ./data/cmd/sleep.vbs
REM start /wait ./data/cmd/sleep.vbs
REM timeout /T 5 /NOBREAK
cmd /c ping -n 3 127.0.0.1
echo ===============

echo %TEST1%
echo %TEST2%
