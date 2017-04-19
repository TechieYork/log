# -*- coding:utf-8 -*-

import socket
import os
import time
import sys
import string

import net_log_ext

#Init net log
#project must less than 128 bytes
#service must less than 128 bytes
def init_net_log(project, service):
    return net_log_ext.init(project, service) 

#Uninit net log
def uninit_net_log():
    return net_log_ext.uninit();

#Send net log
#level include INFO, WARNING, ERROR and FATAL(INFO=0, WARNING=1, ERROR=2, FATAL=3)
#log must less than 4M - 10K
def send_net_log(level, log):
    return net_log_ext.send(level, log)

#Escape
def escape_log(log):
    return log.replace('\r', "[@RETURN]").replace('\n', "[@NEWLINE]")

#Unescape
def unescape_log(log):
    return log.replace("[@RETURN]", "\r").replace("[@NEWLINE]", "\n")
