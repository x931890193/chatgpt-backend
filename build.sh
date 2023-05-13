rm -rf "target"
mkdir "target"
mkdir -p "target"
GOOS=linux GOARCH=arm GOARM=7 go build -o target/chatgpt-backend .
cp config/.config.yml target
