
#include <iostream>
using namespace std;

#include "jni.h"
#include "log_NetLog.h"

#include "net_log.h"

void getJniString(JNIEnv *env, jstring js, std::string& out)
{
	const char *str = env->GetStringUTFChars(js, 0);
	out = str;
	env->ReleaseStringUTFChars(js, str);
}

JNIEXPORT jint JNICALL Java_log_NetLog_init
(JNIEnv *env, jobject obj, jstring jsProject, jstring jsService)
{
	//return 0;
	std::string project;
	getJniString(env, jsProject, project);

	std::string service;
	getJniString(env, jsService, service);

	return init_net_log(project, service);
}

JNIEXPORT void JNICALL Java_log_NetLog_uninit
(JNIEnv *env, jobject obj)
{
	uninit_net_log();
}

JNIEXPORT void JNICALL Java_log_NetLog_logInfo
(JNIEnv *env, jobject obj, jint jLevel, jstring jsLog)
{
	std::string log;
	getJniString(env, jsLog, log);

	int logLevel = jLevel;
	send_net_log(logLevel, log);
}
