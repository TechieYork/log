echo "Begin make golang proto files"
/usr/local/bin/protoc -I=./ --go_out=. log.proto

echo "Begin make cpp proto files"
/usr/local/bin/protoc -I=./ --cpp_out=./ log.proto

echo "Begin make python proto files"
/usr/local/bin/protoc -I=./ --python_out=./ log.proto

echo "Begin make java proto files"
/usr/local/bin/protoc -I=./ --java_out=./ log.proto

echo "End make proto files"

echo "Begin make liblog_proto.a"
cmake .
make -j

echo "End make liblog_proto.a"

#echo "Begin install to /usr/local/include and /usr/local/lib"
#sudo cp ./log.pb.h /usr/local/include
#sudo cp ./log_pb2.py /usr/local/include
#sudo cp ./liblog_proto.a /usr/local/lib

#echo "End install to /usr/local/include and /usr/local/lib"
