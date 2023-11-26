## Example

```sh
docker run -d --name routemonitor-e2e \
  --volume $(pwd):/routemonitor:ro \
  --workdir /routemonitor \
  --cap-add=NET_ADMIN \
  golang:1.21.2-bullseye /bin/bash -c 'tail -f /dev/null'

docker exec -it routemonitor-e2e bash -c '''
apt update && apt install -y iproute2
'''

docker exec -it routemonitor-e2e bash -c '''
ip link add xeth1 type dummy
ip link set dev xeth1 up
ip addr add 10.10.10.25/24 dev xeth1
ip route add 10.10.10.24/32 via 172.17.0.2 dev eth0

# ip route get 1.1.1.1
# 1.1.1.1 via 172.17.0.1 dev eth0 src 172.17.0.2 uid 0
#     cache
# ip route get 10.10.10.1
# 10.10.10.1 dev xeth1 src 10.10.10.25 uid 0
#     cache
# ip route get 10.10.10.24
# 10.10.10.24 dev eth0 src 172.17.0.2 uid 0
#     cache
# ip route get 10.10.10.25
# local 10.10.10.25 dev lo src 10.10.10.25 uid 0
#     cache <local>
'''

docker exec -it routemonitor-e2e go run . -config=example/config.yml
```
