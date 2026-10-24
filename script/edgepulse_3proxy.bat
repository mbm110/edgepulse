:: --------------------------------------------------------------
::	: EdgePulse  3Proxy
::	: 1.0.6
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


:: nowip_3proxy.txt 
if not exist "nowip_3proxy.txt" (
    echo EdgePulse: operation completed.
    echo EdgePulse: operation completed.
    echo EdgePulse: operation completed.
    echo.
    set /p nowip="Message"
    echo !nowip!>nowip_3proxy.txt
    echo.
)  

:: nowip_3proxy.txt  Cloudflare CDN IP
set /p nowip=<nowip_3proxy.txt
echo EdgePulse: operation completed.


::  RESET  " IP " 
::  3  goto :STOP  goto :RESET 
:RESET


echo EdgePulse: operation completed.
echo.|edgepulse.exe -o "result_3proxy.txt"


::  0
if not exist result_3proxy.txt (
    echo.
    echo EdgePulse: operation completed.
    goto :STOP
)

::  IP
for /f "skip=1 tokens=1 delims=," %%i in ('more result_3proxy.txt') do (
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
::for /f %%a in ('type result_3proxy.txt') do set /a v+=1
::if %v% GTR 11 (
::    echo.
echo EdgePulse: operation completed.
::    goto :RESET
::)


echo %bestip%>nowip_3proxy.txt
echo.
echo EdgePulse: operation completed.
echo EdgePulse: operation completed.



::  D:\Program Files\3Proxy  3Proxy 
CD /d "D:\Program Files\3Proxy"
::  3Proxy 



echo.
echo EdgePulse: operation completed.
copy 3proxy.cfg 3proxy.cfg_backup
echo.
echo EdgePulse: operation completed.
(
    for /f "tokens=*" %%i in (3proxy.cfg_backup) do (
        set s=%%i
        set s=!s:%nowip%=%bestip%!
        echo !s!
        )
)>3proxy.cfg

net stop 3proxy
net start 3proxy

echo ...
echo.
:STOP
pause 