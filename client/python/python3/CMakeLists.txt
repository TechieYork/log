cmake_minimum_required(VERSION 2.8)
SET(CMAKE_BUILE_TYPE Debug)
SET(CMAKE_CXX_FLAGS_DEBUG "$ENV{CXXFLAGS} -O0 -Wall -g3 -ggdb")  
SET(CMAKE_CXX_FLAGS_RELEASE "$ENV{CXXFLAGS} -O3 -Wall")

set(PROJECT_NAME net_log_ext)

project(${PROJECT_NAME})

message(STATUS "Project name:${PROJECT_NAME}\r\n")

set(CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} -DBOOST_PYTHON_STATIC_LIB -std=c++11 -fPIC -Wall -O2 -pipe -fno-ident -MMD -D_GNU_SOURCE -D_REENTRANT")
set(CMAKE_C_FLAGS "${CMAKE_C_FLAGS} -DBOOST_PYTHON_STATIC_LIB -std=c++11 -fPIC -Wall -O2 -pipe -fno-ident -MMD -D_GNU_SOURCE -D_REENTRANT")

file(GLOB CPP_FILES "*.cpp")

foreach(SOURCE ${CPP_FILES})
    string(REPLACE ${PROJECT_SOURCE_DIR} "." CPP_FILE ${SOURCE})
    message(${CPP_FILE})
endforeach()

include_directories(./ 
    /usr/local/include 
    /usr/local/src/Python-3.5.1/
    /usr/local/src/Python-3.5.1/Include)

link_directories(/usr/local/lib  )
		
add_library(${PROJECT_NAME} SHARED ${CPP_FILES})

target_link_libraries(${PROJECT_NAME} libnet_log.a liblog_proto.a libboost_python.a libprotobuf.a) 

