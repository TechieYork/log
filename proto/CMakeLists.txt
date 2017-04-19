cmake_minimum_required(VERSION 2.8)

#Set project info
set(PROJECT_NAME log_proto)
project(${PROJECT_NAME})

#Set compiler info
set(CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} -std=c++11 -fPIC -shared -Wall -O2 -pipe -fno-ident -MMD -D_GNU_SOURCE -D_REENTRANT")
set(CMAKE_C_FLAGS "${CMAKE_C_FLAGS} -std=c++11 -fPIC -shared -Wall -O2 -pipe -fno-ident -MMD -D_GNU_SOURCE -D_REENTRANT")

#add include library path
include_directories(./)

#Get all cpp files
#file(GLOB CPP_FILES "*.cpp")
set(CPP_FILES 
    log.pb.cc )
	
add_library(${PROJECT_NAME} STATIC ${CPP_FILES})
