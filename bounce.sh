kc delete -f tmp/
sleep 1
echo "1..."
sleep 1
echo "2..."
sleep 1
echo "3..."
docker build -t reg.captainjustin.space/defrost-proxy:dev3 .
sleep 1
kc apply -f tmp/