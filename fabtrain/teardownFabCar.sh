echo "==================="
echo "Teardown everything"
echo "==================="
cd ./fabtrain-network
docker-compose -f docker-compose.yaml down
sleep 10
docker rm $(docker ps -aq)
sleep 10
docker rmi $(docker images dev-* -q)
docker ps -a
sleep 10
echo
echo "================="
echo "Teardown complete"
echo "================="

