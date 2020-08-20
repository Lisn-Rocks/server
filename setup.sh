#!/usr/bin/bash

# Create the root folder at `~/Public/lisn`.
mkdir root root/logs
cp -r .pkg/pub root/
cp -r .pkg/storage root/
mkdir root/storage/songs root/storage/archives

# The `root/pub/lisn` folder and everything within, will be generated by
# `npm run build` from lisn-web-app folder ... given you set it up properly :)
# See https://github.com/sharpvik/lisn-web-app

if [ ! -d ~/Public ]; then
    mkdir ~/Public
fi
mv root ~/Public/lisn

# Unpack config files. You'll have to edit these!
cp -r .pkg/config config
mv config/.env.example config/.env

