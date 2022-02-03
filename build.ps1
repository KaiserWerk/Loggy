$sourcecode = ".\main.go"
$target = "build\loggy"
$version = "1.0.0"
$commit = git rev-parse HEAD
$date = Get-Date -Format "yyyy-MM-dd HH:mm:ss K"

# Linux, 64-bit
$env:GOOS = 'linux';   $env:GOARCH = 'amd64';             go build -o "$($target)-v$($version)-linux-amd64"   -ldflags "-s -w -X 'main.Version=$($version)' -X 'main.Commit=$($commit)' -X 'main.Date=$($date)'" $sourcecode
# Raspberry Pi
$env:GOOS = 'linux';   $env:GOARCH = 'arm'; $env:GOARM=7; go build -o "$($target)-v$($version)-linux-arm7"    -ldflags "-s -w -X 'main.Version=$($version)' -X 'main.Commit=$($commit)' -X 'main.Date=$($date)'" $sourcecode
# macOS
$env:GOOS = 'darwin';  $env:GOARCH = 'amd64';             go build -o "$($target)-v$($version)-macos-amd64"   -ldflags "-s -w -X 'main.Version=$($version)' -X 'main.Commit=$($commit)' -X 'main.Date=$($date)'" $sourcecode
# Windows, 64-bit
$env:GOOS = 'windows'; $env:GOARCH = 'amd64';             go build -o "$($target)-v$($version)-win-amd64.exe" -ldflags "-s -w -X 'main.Version=$($version)' -X 'main.Commit=$($commit)' -X 'main.Date=$($date)'" $sourcecode