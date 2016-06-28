rem This is the NIAS batch file terminator.  

@echo off

rem Stop the NIAS services. 
taskkill /IM gnatsd.exe
taskkill /IM  harness.exe


