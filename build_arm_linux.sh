rm tracker
rm tracker-arm-linux.zip

# Needs ARM compiler (gcc-arm-linux-gnueabihf)
CC=arm-linux-gnueabihf-gcc GOOS=linux GOARCH=arm GOARM=6 CGO_ENABLED=1 go build

zip tracker-arm-linux.zip tracker
zip -u tracker-arm-linux.zip config/config.ini
zip -u tracker-arm-linux.zip data/
zip -u tracker-arm-linux.zip scripts/db/*.*