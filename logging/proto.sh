echo "Generating protofiles for Logging Service"
echo "Language: Golang"
protoc -I ./proto ./proto/logging.proto ./proto/smartmeterdb.proto --go_out=plugins=grpc:.
if [ $? -eq 0 ]; then
  echo "Files Generated Successfully"
  else
    echo "Files Genration Failed"
   fi
