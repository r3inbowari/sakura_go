@echo off

set packageName=cmd
set appName=sakura
set buildVersion=v1.3.6
set major=1
set minor=3
set patch=6
set mode=REL

for /f "delims=" %%i in ('go version') do (set goVersion=%%i)
for /f "delims=" %%i in ('git show -s --format^=%%H') do (set gitHash=%%i)
for /f "delims=" %%i in ('git show -s --format^=%%cd') do (set buildTime=%%i)

echo ===================================================
echo                  Go build running
echo ===================================================
echo %goVersion%
echo build hash %gitHash%
echo build time %buildTime%
echo build tag %buildVersion%
echo ===================================================

cd %packageName%

set GOOS=windows
set GOARCH=amd64
go env -w GOOS=windows
go env -w GOARCH=amd64
go build -ldflags "-X 'main.major=%major%' -X 'main.minor=%minor%'-X 'main.patch=%patch%' -X 'main.releaseVersion=%buildVersion%' -X 'main.mode=%mode%' -X 'main.goVersion=%goVersion%' -X 'main.gitHash=%gitHash%' -X 'main.buildTime=%buildTime%'" -o ../build/%appName%_%GOOS%_%GOARCH%_%buildVersion%
echo Done %appName%_%GOOS%_%GOARCH%_%buildVersion%
set upxArgs=%upxArgs% %appName%_%GOOS%_%GOARCH%_%buildVersion%.exe

set GOOS=darwin
set GOARCH=amd64
go env -w GOOS=darwin
go env -w GOARCH=amd64
go build -ldflags "-X 'main.major=%major%' -X 'main.minor=%minor%'-X 'main.patch=%patch%' -X 'main.releaseVersion=%buildVersion%' -X 'main.mode=%mode%' -X 'main.goVersion=%goVersion%' -X 'main.gitHash=%gitHash%' -X 'main.buildTime=%buildTime%'" -o ../build/%appName%_%GOOS%_%GOARCH%_%buildVersion%
echo Done %appName%_%GOOS%_%GOARCH%_%buildVersion%
set upxArgs=%upxArgs% %appName%_%GOOS%_%GOARCH%_%buildVersion%

set GOOS=linux
set GOARCH=amd64
go env -w GOOS=linux
go env -w GOARCH=amd64
go build -ldflags "-X 'main.major=%major%' -X 'main.minor=%minor%'-X 'main.patch=%patch%' -X 'main.releaseVersion=%buildVersion%' -X 'main.mode=%mode%' -X 'main.goVersion=%goVersion%' -X 'main.gitHash=%gitHash%' -X 'main.buildTime=%buildTime%'" -o ../build/%appName%_%GOOS%_%GOARCH%_%buildVersion%
echo Done %appName%_%GOOS%_%GOARCH%_%buildVersion%
set upxArgs=%upxArgs% %appName%_%GOOS%_%GOARCH%_%buildVersion%

set GOOS=windows
set GOARCH=amd64
go env -w GOOS=windows
go env -w GOARCH=amd64
go build -ldflags "-X 'main.major=%major%' -X 'main.minor=%minor%'-X 'main.patch=%patch%' -X 'main.releaseVersion=%buildVersion%' -X 'main.mode=%mode%' -X 'main.goVersion=%goVersion%' -X 'main.gitHash=%gitHash%' -X 'main.buildTime=%buildTime%'" -o ../build/%appName%_%GOOS%_%GOARCH%_%buildVersion%.exe
echo Done %appName%_%GOOS%_%GOARCH%_%buildVersion%.exe
set upxArgs=%upxArgs% %appName%_%GOOS%_%GOARCH%_%buildVersion%.exe

set mode=AliyunFC
set GOOS=linux
set GOARCH=amd64
go env -w GOOS=linux
go env -w GOARCH=amd64
go build -ldflags "-X 'main.major=%major%' -X 'main.minor=%minor%'-X 'main.patch=%patch%' -X 'main.releaseVersion=%buildVersion%' -X 'main.mode=%mode%' -X 'main.goVersion=%goVersion%' -X 'main.gitHash=%gitHash%' -X 'main.buildTime=%buildTime%'" -o ../build/%appName%_%GOOS%_%GOARCH%_%mode%_%buildVersion%
echo Done %appName%_%GOOS%_%GOARCH%_%mode%_%buildVersion%
set upxArgs=%upxArgs% %appName%_%GOOS%_%GOARCH%_%buildVersion%

cd ../build
certutil -hashfile sakura_windows_amd64_%buildVersion% MD5
certutil -hashfile sakura_windows_amd64_%buildVersion%.exe MD5
certutil -hashfile sakura_darwin_amd64_%buildVersion% MD5
