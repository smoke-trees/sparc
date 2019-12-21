echo "Generating grpc protofiles for Authentication Service"
echo "Language: Golang"
protoc -I ./proto ./proto/auth.proto --go_out=plugins=grpc:.
if [ $? -eq 0 ]; then
  echo "Files Generated Successfully"
  else
    echo "Files Genration Failed"
   fi
