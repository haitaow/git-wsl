@echo off
set WSLDIR=%HOMEDRIVE%%HOMEPATH%\bin

if exist %WSLDIR% (
  echo Install to existing directory %WSLDIR%
) else (
  echo Install to new directory %WSLDIR%
  md %WSLDIR%
)

copy wsl-git.exe %WSLDIR%
copy wsl-git.json %WSLDIR%
if exist %WSLDIR%\git.exe (
  echo %WSLDIR%\git.exe already exists.  Delete it and re-run install.bat will create it as a link to %WSLDIR%\wsl-git.exe.
) else (
  echo Creating %WSLDIR%\git.exe as link to %WSLDIR%\wsl-git.exe.
  echo If it fails, please run install.bat as administrator.
  mklink %WSLDIR%\git.exe %WSLDIR%\wsl-git.exe
)
