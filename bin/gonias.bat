rem This is the NIAS batch file launcher. Add extra validators to the bottom of this list. 
rem Change the directory as appropriate (go-nias)
rem gnatsd MUST be the first program launched

@echo off

rem Change to the NIAS install directory
cd\go-nias

rem Run the NIAS services. Add to the BOTTOM of this list
start gnatsd
start harness

rem Run the web client (launch browser here)
start http://localhost:1325/nias
