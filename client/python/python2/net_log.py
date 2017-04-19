# -*- coding:utf-8 -*-

import socket
import os
import log_pb2
import time
import sys
import string

print("------- Begin test collector -------")

#Net log path
net_log_path = "/var/tmp/net_log.sock";

#Net log socket
net_log_socket = socket.socket(socket.AF_UNIX, socket.SOCK_DGRAM)

#Net log socket send buffer max length
net_log_send_buffer_len = 4 * 1024 * 1024

#Net log content max length
net_log_max_len = net_log_send_buffer_len - 10240

net_log_socket.setsockopt(socket.SOL_SOCKET, socket.SO_SNDBUF, net_log_send_buffer_len)

#Net log is initialized
net_log_is_init = False

#Net log project
net_log_project = ""

#Net log service
net_log_service = ""

#Net log escape tags
return_tag = "[@RETURN]"
newline_tag = "[@NEWLINE]"

#Init net log
#project must less than 128 bytes
#service must less than 128 bytes
def init_net_log(project, service):
    global net_log_socket
    global net_log_project
    global net_log_service
    global net_log_is_init

    if None == net_log_socket:
        return -1

    if 128 < len(project) or 128 < len(service):
        return -3

    net_log_project = project
    net_log_service = service

    net_log_is_init = True
    return 0

#Uninit net log
def uninit_net_log():
    global net_log_is_init
    global net_log_socket
    
    net_log_is_init = False

    if None != net_log_socket:
        net_log_socket.close()
        net_log_socket = None

    return 0

#Send net log
#level include INFO, WARNING, ERROR and FATAL(INFO=0, WARNING=1, ERROR=2, FATAL=3)
#log must less than 4M - 10K
def send_net_log(level, log):
    global net_log_is_init
    global net_log_project
    global net_log_service
    global net_log_socket

    #Check net log is initialized
    if False == net_log_is_init:
        return -1

    #Check params
    if net_log_max_len < len(log):
        return -2

    #Generate log buffer
    log_package = log_pb2.LogPackage()

    log_package.project = net_log_project
    log_package.service = net_log_service
    log_package.level = level
    log_package.log = log

    log_buffer = log_package.SerializeToString()

    ret = net_log_socket.sendto(log_buffer, net_log_path)

    if 0 >= ret:
        return -4

    return 0

#Escape
def escape_log(log):
    return log.replace('\r', return_tag).replace('\n', newline_tag)

#Unescape
def unescape_log(log):
    return log.replace(return_tag, "\r").replace(newline_tag, "\n")