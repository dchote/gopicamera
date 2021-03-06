# Building OpenCV 4 on Raspbian

Install required dependencies
```
sudo apt-get install --yes libjpeg-dev libtiff5-dev libjasper-dev libpng12-dev libavcodec-dev libavformat-dev \
	libswscale-dev libxvidcore-dev libx264-dev libv4l-dev liblapacke-dev libgtk-3-dev \
	libopenblas-dev libhdf5-dev libtesseract-dev libleptonica-dev \
	python3-numpy python3-dev checkinstall cmake gfortran curl \
	libavcodec57 libavformat57 libavutil-dev libavutil55 libraw1394-dev \
	libraw1394-tools libswresample-dev libswresample2 libswscale4
```
---

Build and install an updated version of `checkinstall` to fix the segfault that currently happens on raspbian (https://github.com/opencv/opencv/issues/8897).
```
git clone https://github.com/giuliomoro/checkinstall
cd checkinstall
make install
```
---

Then start the build process and go have a cup of tea or something...
```
./build_opencv4.sh
```
