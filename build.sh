#!/bin/bash

cd ~/upfast-tf

go build -o main .

systemctl restart upfast
