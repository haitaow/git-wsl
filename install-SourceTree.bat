@echo off

if "" == "%GIT_DIR%" (
    set GIT_DIR=%LOCALAPPDATA%\Atlassian\SourceTree\git_local\mingw32\bin
)
if not exist "%GIT_DIR%\git.exe" (
    echo "%GIT_DIR%\git.exe not found!  Set GIT_DIR to where git.exe is located before running install."
    exit 1
)

if not exist "%GIT_DIR%\win-git.exe" (
    echo copy "%GIT_DIR%\git.exe" "%GIT_DIR%\win-git.exe"
    copy "%GIT_DIR%\git.exe" "%GIT_DIR%\win-git.exe"
)

@echo on
copy wsl-git.exe %GIT_DIR%\
copy wsl-git.json %GIT_DIR%\
copy wsl-git.exe %GIT_DIR%\git.exe
