#! /bin/sh

sudo mkdir /usr/local/clippy
sudo chown $(whoami) /usr/local/clippy
cd .. && go install
# this starts the background process
clippy start
