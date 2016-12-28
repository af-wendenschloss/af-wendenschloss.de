@echo off
if "%~1"=="-FIXED_CTRL_C" (
   REM Remove the -FIXED_CTRL_C parameter
   SHIFT
) ELSE (
   REM Run the batch with <NUL and -FIXED_CTRL_C
   CALL <NUL %0 -FIXED_CTRL_C %*
   GOTO :EOF
)
set OLD_CD=%CD%
set SCRIPT=%~dp0
set DIRECTORY=%SCRIPT:~0,-1%
"%DIRECTORY%\tools\bin\update_hugo.exe"
cd "%DIRECTORY%\src"
"%DIRECTORY%\tools\bin\hugo.exe"  server --buildDrafts
cd "%OLD_CD%"

:EOF
