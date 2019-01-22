rm tracker
rm tracker-linux.zip

go build

zip tracker-linux.zip tracker
zip -u tracker-linux.zip config/config.ini
zip -u tracker-linux.zip data/
zip -u tracker-linux.zip scripts/db/*.*