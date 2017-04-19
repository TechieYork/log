#ifndef NET_LOG_H
#define NET_LOG_H

#include <sys/types.h>
#include <sys/socket.h>
#include <sys/un.h>

#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>

#include <limits.h>

#include <string>
using namespace std;

#include "log.pb.h"

extern "C"
{

#define SEND_LOG(level, log) send_net_log(level, generate_log_position(__FILE__, __LINE__) + generate_log(log))

#define SEND_LOG_INFO(log) SEND_LOG(INFO, log)
#define SEND_LOG_WARNING(log) SEND_LOG(WARNING, log)
#define SEND_LOG_ERROR(log) SEND_LOG(ERROR, log)
#define SEND_LOG_FATAL(log) SEND_LOG(FATAL, log)

	enum LogLevel
	{
		//Log severity definition
		INFO = 0,
		WARNING = 1,
		ERROR = 2,
		FATAL = 3,
	};

	string generate_log_position(const string& file, const int line);
    const string generate_log(const string& log);

	//Init net log
	//project must less than 128 bytes
	//service must less than 128 bytes
	int init_net_log(const string& project, const string& service);

	//Uninit net log
	int uninit_net_log();

	//Send net log
	//level include INFO, WARNING, ERROR and FATAL(INFO=0, WARNING=1, ERROR=2, FATAL=3)
	//log must less than 4M - 10K
	int send_net_log(const int level, const string& log);

	//Escape
	int escape_log(string& log);

	//Unescape
	int unescape_log(string& log);

}

#endif
