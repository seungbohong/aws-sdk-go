#!/bin/bash

# Go 소스 코드 빌드
go build -o s3test main.go

# 빌드된 바이너리 실행
./s3test

# 빌드된 바이너리 파일 삭제
rm s3test