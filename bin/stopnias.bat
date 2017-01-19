rem This is the NIAS batch file terminator.  

@echo off

rem Stop the NIAS services. 
taskkill /IM nats-streaming-server.exe
taskkill /IM napval.exe


