@echo off
if "%~1"=="-FIXED_CTRL_C" (
   rem Remove the -FIXED_CTRL_C parameter
   SHIFT
) else (
   rem Run the batch with <NUL and -FIXED_CTRL_C
   call <NUL %0 -FIXED_CTRL_C %*
   goto :EOF
)
set OLD_CD=%CD%
set SCRIPT=%~dp0
set DIRECTORY=%SCRIPT:~0,-1%
"%DIRECTORY%\tools\bin\update_hugo.exe"
if not errorlevel 0 goto :EOF
cd "%DIRECTORY%\src"
"%DIRECTORY%\tools\bin\hugo.exe"  server --buildDrafts
cd "%OLD_CD%"

:EOF
