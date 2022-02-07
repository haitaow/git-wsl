# wsl-git

wsl-git is a utility for integration with git UI such as SourceTree, to workaround the problem that Windows' version of
git is slow against repo in \\wsl$\. It forwards the git commands issued by SourceTree etc. to the git running in WSL
which is fast.

## Note

*It's still a work-in-progress and under testing. After it stabilizes, usage steps will be updated.*

It forward the git command to the default distro of wsl at this time.  Run "wsl -l" to verify your default distro, and
run "wsl git version" in Windows PowerShell or cmd window to verify you have git installed in the default distro.

* Warning:
    * Some Windows paths in arguments are not translated to Linux paths yet. Some git commands may fail. Frequently used
      features in SourceTree seemed OK so far. If something is not working as expected, please inspect the log file.
      Please report issues (or you can try to fix it yourself :) ).

## Get files

* From binary
    * download wsl-git.zip
    * extract to some-dir
    * cd some-dir
* From source
    * git clone
    * cd wsl-git
    * go build wsl-git.go

## SourceTree integration

* Note: by default install-SourceTree.bat only works with the embedded version of git in SourceTree, i.e. in the above
  options section, "Embedded" should be used. If "Update Embedded" is clicked, ./install-SourceTree.bat needs to be
  re-run.  If using another version is desired, set GIT_DIR before the next step. 
* Note: the next step will backup the original git.exe to win-git.exe and replace it.  If win-git.exe already exists,
  it will not be overwritten, so re-run is safe. 
* ./install-SourceTree.bat
* (optional, recommended) In SourceTree, go to Tools | Options, then Git tab, Git Version section, clear "Disable
  LibGit2 integration".

## Standalone install (not required for SourceTree integration)

* First time install: run ./install-standalone.bat as administrator
* Update option 1: same as first time install but run as administrator is not required
* Update option 2: just copy the files in the zip folder to the existing installation folder

## Notes

* If a json config file is not found in the same directory as the binary, a default one will be generated the first time
  wsl-git.exe is executed,
* Log(s) will be written to the same directory as wsl-git.exe, so ensure it has write permission.
* By default, any git command exited with error will be logged, unless it's a known error to be skipped. Add to
  IgnoreErrors section in the json config file if needed. Some common errors are already included in the default config
  file, and can be removed if desired.
* Setting LogAll to true in the config file will log all git commands regardless of success or error.  
