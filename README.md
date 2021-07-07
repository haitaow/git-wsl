# wsl-git

wsl-git is a utility for integration with git UI such as SourceTree, to workaround the problem that Windows' version of
git is slow against repo in \\wsl$\. It forwards the git commands issued by SourceTree etc. to the git running in WSL
which is fast.

## Usage

*It's still a work-in-progress and under testing. After it stabilizes, usage steps will be updated.*

Temporary steps for SourceTree integration

* Warning:
  * The current version will log all git commands passed through. Clear log file once in a while may be a good idea.
  * Windows paths in arguments are not translated to Linux paths yet. Some git commands may fail. Frequently used
    features in SourceTree seemed OK so far. If something is not working as expected, please inspect the log file.
    Please report issues (or you can try to fix it yourself :) ).
* download wsl-git.exe, or "go build wsl-git.go".
* (optional) copy wsl-git.exe to another location when playing around the source code.
* go to a desired directory, e.g. \Users\username\bin
* mklink git.exe wsl-git.exe
* log will be written to the same directory, so ensure it has write permission.
* In SourceTree, go to Tools | Options, then Git tab, Git Version section
  * Click "Embedded" if it's not grayed out
  * Click "Clear path cache" if it's not grayed out
  * Click "System", and browse to the git.exe symbolic link created above
  * (optional, recommended) clear "Disable LibGit2 integration"
  * Click OK
  * Restart SourceTree
    