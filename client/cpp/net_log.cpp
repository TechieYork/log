#include "net_log.h"

//Net log path
string net_log_path = "/var/tmp/net_log.sock";

//Net log socket send buffer length
int net_log_send_buffer_len = 4 * 1024 * 1024;

//Net log max length
int net_log_max_len = net_log_send_buffer_len - 10240;

//Net log unix domain socket descriptor
int net_log_socket = 0;

//Net log unix domain socket address
struct sockaddr_un net_log_socket_addr;

//Net log unix domain socket address length
int net_log_socket_addr_len = 0;

//Net log is initialized
bool net_log_is_init = false;

//Net log project
string net_log_project;

//Net log service
string net_log_service;

//Net log escape tags
string return_tag = "[@RETURN]";
string newline_tag = "[@NEWLINE]";

string generate_log_position(const string& file, const int line)
{
	return "[" + file.substr(file.find_last_of("//") + 1) + ":" + to_string(line) + "]";
}

const string generate_log(const string& log)
{
    string log_temp = log;
    escape_log(log_temp);
    return log_temp;
}

//Init net log
//project must less than 128 bytes
//service must less than 128 bytes
int init_net_log(const string& project, const string& service)
{
	net_log_socket = socket(PF_LOCAL, SOCK_DGRAM, 0);

	if (0 == net_log_socket)
	{
		return -1;
	}

	net_log_socket_addr_len = sizeof(net_log_socket_addr);
	net_log_socket_addr.sun_family = AF_UNIX;

	if (0 > snprintf(net_log_socket_addr.sun_path, PATH_MAX, net_log_path.c_str()))
	{
		return -2;
	}

	if (project.empty() || service.empty()
		|| 128 < project.size() || 128 < service.size())
	{
		return -3;
	}

	if (0 != setsockopt(net_log_socket, SOL_SOCKET, SO_SNDBUF, &net_log_send_buffer_len, sizeof(net_log_send_buffer_len)))
	{
		return -4;
	}

	net_log_project = project;
	net_log_service = service;

	net_log_is_init = true;

	return 0;
}

//Uninit net log
int uninit_net_log()
{
	net_log_is_init = false;

	close(net_log_socket);

	net_log_socket = 0;
	net_log_socket_addr_len = 0;

	return 0;
}

//Send net log
//level include INFO, WARNING, ERROR and FATAL(INFO=0, WARNING=1, ERROR=2, FATAL=3)
//log must less than 4M - 10K
int send_net_log(const int level, const string& log)
{
	//Check net log is initialized
	if (!net_log_is_init)
	{
		return -1;
	}

	//Check params
	if ((size_t)net_log_max_len < log.size())
	{
		return -2;
	}

	//Generate log buffer
	log_proto::LogPackage log_package;

	log_package.set_level(level);
	log_package.set_project(net_log_project);
	log_package.set_service(net_log_service);
	log_package.set_log(log);

	string log_buffer;

	if (!log_package.SerializeToString(&log_buffer))
	{
		return -3;
	}

	//Send log
	int ret = sendto(net_log_socket, log_buffer.c_str(), log_buffer.size(), 0, (struct sockaddr *)&net_log_socket_addr, net_log_socket_addr_len);

	if (0 > ret)
	{
		return -4;
	}

	return 0;
}

//Escape
int escape_log(string& log)
{
	std::string::size_type pos = 0;

	while (true)
	{
		pos = log.find('\r', pos);

		if (pos == string::npos)
		{
			break;
		}

		log.replace(pos, 1, return_tag);
		pos += return_tag.length();
	}

	pos = 0;

	while (true)
	{
		pos = log.find('\n', pos);

		if (pos == string::npos)
		{
			break;
		}

		log.replace(pos, 1, newline_tag);
		pos += newline_tag.length();
	}

	return 0;
}

//Unescape
int unescape_log(string& log)
{
	std::string::size_type pos = 0;

	while (true)
	{
		pos = log.find(return_tag, pos);

		if (pos == string::npos || pos >= log.length())
		{
			break;
		}

		log.replace(pos, return_tag.length(), "\r");
		pos += 1;
	}

	pos = 0;

	while (true)
	{
		pos = log.find(newline_tag, pos);

		if (pos == string::npos || pos >= log.length())
		{
			break;
		}

		log.replace(pos, newline_tag.length(), "\n");
		pos += 1;
	}

	return 0;
}
