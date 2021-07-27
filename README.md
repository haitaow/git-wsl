# wsl-git

wsl-git is a utility for integration with git UI such as SourceTree, to workaround the problem that Windows' version of
git is slow against repo in \\wsl$\. It forwards the git commands issued by SourceTree etc. to the git running in WSL
which is fast.

## Note

*It's still a work-in-progress and under testing. After it stabilizes, usage steps will be updated.*

* Warning:
  * Some Windows paths in arguments are not translated to Linux paths yet. Some git commands may fail. Frequently used
    features in SourceTree seemed OK so far. If something is not working as expected, please inspect the log file.
    Please report issues (or you can try to fix it yourself :) ).

## Installation

* From binary
  * download wsl-git.zip
  * First time install: extract to a folder and fun install.bat as administrator
  * Update option 1: same as first time install but run as administrator is not required
  * Update option 2: just copy the files in the zip folder to the existing installation folder
* From source
  * git clone
  * go build wsl-git.go
  * (optional) copy wsl-git.exe to another location when playing around the source code.
  * go to a desired directory, e.g. \Users\username\bin,
  * mklink git.exe wsl-git.exe
  * a default json config file will be generated the first time wsl-git.exe is executed.
  * log will be written to the same directory as wsl-git.exe, so ensure it has write permission.

## SourceTree integration

* In SourceTree, go to Tools | Options, then Git tab, Git Version section
* Click "Embedded" if it's not grayed out
* Click "Clear path cache" if it's not grayed out
* Click "System", and browse to the git.exe symbolic link created above
* (optional, recommended) clear "Disable LibGit2 integration"
* Click OK
* Restart SourceTree
