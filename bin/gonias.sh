# This is the *nix NIAS batch file launcher. Add extra validators to the bottom of this list. 
# Change the directory as appropriate (go-nias)
# gnatsd MUST be the first program launched

if [ -f "nias.pid" ]
then
echo "There is a nias.pid file in place; run stopnias.sh"
exit
fi

#rem Run the NIAS services. Add to the BOTTOM of this list
# store each PID in pid list
./nats-streaming-server -p 4223 -sc nias_nss.cfg & echo $! > nias.pid

# give the nats server time to come up
sleep 5

./go-nias8 & echo $! >> nias.pid

echo "Run the web client (launch browser here):"
echo "http://localhost:1325/nias (NAPLAN Registration)"

