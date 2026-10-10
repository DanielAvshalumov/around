#!/bin/bash

tor \
  --SocksPort 127.0.0.1:9050 \
  --ControlPort 127.0.0.1:9051 \
  --CookieAuthentication 1 &

echo "tor finished"
exec /bin/around