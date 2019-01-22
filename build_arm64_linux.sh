rm tracker
rm tracker-arm64-linux.zip

# Needs ARM64 compiler (gcc-aarch64-linux-gnu)
CC=aarch64-linux-gnu-gcc GOOS=linux GOARCH=arm64 CGO_ENABLED=1 go build

zip tracker-arm64-linux.zip tracker
zip -u tracker-arm64-linux.zip config/config.ini
zip -u tracker-arm64-linux.zip data/
zip -u tracker-arm64-linux.zip scripts/db/*.*