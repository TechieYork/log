#include <limits.h>

#include <string>
using namespace std;

#include "log.pb.h"

#include "net_log.h"

//#include "Python.h"
#include "boost/python.hpp"

//Init net log
//project must less than 128 bytes
//service must less than 128 bytes
int net_log_init(const char* project, const char* service)
{
    return init_net_log(project, service);
}

//Uninit net log
int net_log_uninit()
{
    return uninit_net_log();
}

//Send net log
//level include INFO, WARNING, ERROR and FATAL(INFO=0, WARNING=1, ERROR=2, FATAL=3)
//log must less than 4M - 10K
int net_log_send(const int level, const char* log)
{
    return send_net_log(level, log);
}

BOOST_PYTHON_MODULE(net_log_ext)
{
    using namespace boost::python;

    def("init", net_log_init);
    def("uninit", net_log_uninit);
    def("send", net_log_send);
}
