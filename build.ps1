$sourcecode = ".\main.go"
$target = "build\loggy"
$version = "1.0.0"
$date = Get-Date -Format "yyyy-MM-dd HH:mm:ss K"
# Windows, 64-bit
$env:GOOS = 'windows'; $env:GOARCH = 'amd64';             go build -o "$($target)-win64.exe" -ldflags "-s -w -X 'main.Version=$($version)' -X 'main.Date=$($date)'" $sourcecode
# Linux, 64-bit
$env:GOOS = 'linux';   $env:GOARCH = 'amd64';             go build -o "$($target)-linux64"   -ldflags "-s -w -X 'main.Version=$($version)' -X 'main.Date=$($date)'" $sourcecode
# Raspberry Pi
$env:GOOS = 'linux';   $env:GOARCH = 'arm'; $env:GOARM=5; go build -o "$($target)-raspi32"   -ldflags "-s -w -X 'main.Version=$($version)' -X 'main.Date=$($date)'" $sourcecode
# macOS
$env:GOOS = 'darwin';  $env:GOARCH = 'amd64';             go build -o "$($target)-macos64"   -ldflags "-s -w -X 'main.Version=$($version)' -X 'main.Date=$($date)'" $sourcecode