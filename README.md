# GoPing
My Practice
deterioration version of ping
# DEMO 
I use raw socket, so require root privs.
`$ sudo go run main.go google.com`
Using tshark, tcpdump and so on, you can confirm GoPing's packet
1. running the GoPing
`$ sudo go run main.go 8.8.8.8`
2. capture the packet 
```
$ sudo tshark -ni en0 -f icmp
 Capturing on 'Wi-Fi: en0'
     1   0.000000 10.138.165.122 → 8.8.8.8      ICMP 57 Echo (ping) request  id=0x1a7d, seq=0/0, ttl=64
     2   0.053929      8.8.8.8 → 10.138.165.122 ICMP 60 Echo (ping) reply    id=0x1a7d, seq=0/0, ttl=113 (request in 1)
     3   1.000725 10.138.165.122 → 8.8.8.8      ICMP 57 Echo (ping) request  id=0x1a7d, seq=1/256, ttl=64
     4   1.013489      8.8.8.8 → 10.138.165.122 ICMP 60 Echo (ping) reply    id=0x1a7d, seq=1/256, ttl=113 (request in 3)
     5   2.005206 10.138.165.122 → 8.8.8.8      ICMP 57 Echo (ping) request  id=0x1a7d, seq=2/512, ttl=64
     6   2.049537      8.8.8.8 → 10.138.165.122 ICMP 60 Echo (ping) reply    id=0x1a7d, seq=2/512, ttl=113 (request in 5)
     7   3.009508 10.138.165.122 → 8.8.8.8      ICMP 57 Echo (ping) request  id=0x1a7d, seq=3/768, ttl=64
     8   3.044697      8.8.8.8 → 10.138.165.122 ICMP 60 Echo (ping) reply    id=0x1a7d, seq=3/768, ttl=113 (request in 7)
     9   4.014061 10.138.165.122 → 8.8.8.8      ICMP 57 Echo (ping) request  id=0x1a7d, seq=4/1024, ttl=64
     ....
```
# Author
[井上大誠](https://github.com/taise-hub)
