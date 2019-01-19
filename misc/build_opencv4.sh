#!/usr/bin/env bash

# Config
MAINTAINER="Daniel Chote"
VERSION=4.0.1
                        
# Make a new directory
rm -rf ~/opencv-build
mkdir ~/opencv-build
cd ~/opencv-build

# Download OpenCV
curl -L https://github.com/opencv/opencv/archive/${VERSION}.tar.gz | tar xz
curl -L https://github.com/opencv/opencv_contrib/archive/${VERSION}.tar.gz | tar xz

cd opencv-${VERSION}
# Make build directory.
mkdir build
# Change to 
cd build 

# define settings for a pi zero compat build
cmake -D CMAKE_BUILD_TYPE=RELEASE -D CMAKE_INSTALL_PREFIX=/usr/local -D BUILD_SHARED_LIBS=ON\
        -D INSTALL_PYTHON_EXAMPLES=ON -D INSTALL_C_EXAMPLES=ON -D ENABLE_PRECOMPILED_HEADERS=ON \
        -D OPENCV_EXTRA_MODULES_PATH=~/opencv-build/opencv_contrib-${VERSION}/modules \
        -D PYTHON_DEFAULT_EXECUTABLE=$(which python3) -D PYTHON_EXECUTABLE=$(which python3) \
        -D ENABLE_NEON=OFF -D ENABLE_VFPV3=OFF -D BUILD_EXAMPLES=ON -D OPENCV_GENERATE_PKGCONFIG=ON ..

make -j $(shell nproc --all)
make

checkinstall --default \
--type debian --install=no \
--pkgname opencv4 \
--pkgversion "${VERSION}" \
--pkglicense BSD \
--deldoc --deldesc --delspec \
--requires "libtesseract3,ffmpeg,libjasper1" \
--pakdir ~ --maintainer "${MAINTAINER}" --provides opevcv4 \
--addso --autodoinst \
make install
