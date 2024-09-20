#!/bin/bash

cd /root/upfast-tf

go build -o main .

systemctl reload upfast
