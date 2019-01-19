# gopicamera

RaspberryPi camera streamer with VueJS frontend and automatic discovery of all the other cameras running on the local network.  

![screenshot](images/screenshot.png "Screenshot")

gopicamera is a work in progress, but it will eventually include authentication and wifi configuration from the web interface. The end goal is to roll an SD card image for the Raspberry Pi Zero W that can get you booted in to a functional gopicamera instance as quickly as possible.

I started this project because I want to be able to aggregate many camera sources into a single backend service that can do various object/state detection on all the incoming streams at the same time. 
One of the first models I will be training is the detection of failed 3d prints (ie. spaghetti mess), when the system detects a failed print it can notify you... hopefully saving you from an even bigger mess.

## Supporting packages

A precompiled package for OpenCV4 can be found in [/packages](/packages/), please also install all of the other dependancies listed.


## Running gopicamera

Ensure the bcm2835-v4l2 kernel module is loaded `sudo modprobe bcm2835-v4l2`, edit config.json to suite your needs, if using a MIPI camera ensure `camera.deviceID = -1`. Run as root `sudo ./gopicamera`