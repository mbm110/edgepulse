:: --------------------------------------------------------------
::	: EdgePulse  Hosts
::	: 1.0.5
::	: mbm110
::	: https://github.com/mbm110/edgepulse
:: --------------------------------------------------------------
@echo off
Setlocal Enabledelayedexpansion

::

>nul 2>&1 "%SYSTEMROOT%\system32\cacls.exe" "%SYSTEMROOT%\system32\config\system" 

if '%errorlevel%' NEQ '0' (  
    goto UACPrompt  
) else ( goto gotAdmin )  

:: vbs bat

:UACPrompt  
    echo Set UAC = CreateObject^("Shell.Application"^) > "%temp%\getadmin.vbs" 
    echo UAC.ShellExecute "%~s0", "", "", "runas", 1 >> "%temp%\getadmin.vbs" 
    "%temp%\getadmin.vbs" 
    exit /B  

:: vbs 
  
:gotAdmin  
    if exist "%temp%\getadmin.vbs" ( del "%temp%\getadmin.vbs" )  
    pushd "%CD%" 
    CD /D "%~dp0" 


::


:: nowip_hosts.txt 
if not exist "nowip_hosts.txt" (
    echo EdgePulse: operation completed.
    echo EdgePulse: operation completed.
    echo.
    echo EdgePulse: operation completed.
    set /p nowip="Message"
    echo !nowip!>nowip_hosts.txt
    echo.
)  

:: nowip_hosts.txt  Hosts  Cloudflare CDN IP
set /p nowip=<nowip_hosts.txt
echo EdgePulse: operation completed.


::  RESET  " IP " 
::  3  goto :STOP  goto :RESET 
:RESET


echo EdgePulse: operation completed.
echo.|edgepulse.exe -o "result_hosts.txt"


::  0
if not exist result_hosts.txt (
    echo.
    echo EdgePulse: operation completed.
    goto :STOP
)

::  IP
for /f "skip=1 tokens=1 delims=," %%i in ('more result_hosts.txt') do (
    SET bestip=%%i
    goto :END
)

:END

::  IP  IP 
if "%bestip%"=="" (
    echo.
    echo EdgePulse: operation completed.
    goto :STOP
)
if "%bestip%"=="%nowip%" (
    echo.
    echo EdgePulse: operation completed.
    goto :STOP
)


::  " IP " 
::  IP EdgePulse  IP 
::  -sl  :: 10  11
::set /a v=0
::for /f %%a in ('type result_hosts.txt') do set /a v+=1
::if %v% GTR 11 (
::    echo.
echo EdgePulse: operation completed.
::    goto :RESET
::)


echo %bestip%>nowip_hosts.txt
echo.
echo EdgePulse: operation completed.
echo EdgePulse: operation completed.

CD /d "C:\Windows\System32\drivers\etc"
echo.
echo EdgePulse: operation completed.
copy hosts hosts_backup
echo.
echo EdgePulse: operation completed.
(
    for /f "tokens=*" %%i in (hosts_backup) do (
        set s=%%i
        set s=!s:%nowip%=%bestip%!
        echo !s!
        )
)>hosts

echo ...
echo.
:STOP
pause 