#!/usr/bin/env python3

import sys
import os
import json
import pdb
import argparse
from collections import defaultdict

default_net_bandwidth = {'LAN':10,'WAN':100,'T3': 1000, 'T2': 10000, 'T1': 1000000}
default_net_latency = {'LAN':10,'WAN':50,'T3': 100, 'T2': 1000, 'T1': 1000}
default_intrfc_bandwidth = {'LAN':10,'WAN':100,'T3': 1000, 'T2': 10000, 'T1': 1000000}

Networks = {}
Routers = {}
Hosts   = {}

object_id = 0

def getArgs(cmdline):
    parser = argparse.ArgumentParser()

    parser.add_argument('-nets', metavar='file with networks and connections', dest='net_file', required=True)
    parser.add_argument('-autodev', metavar='output file with automatically created routers', dest='autodev_file', required=True)

    args = parser.parse_args(cmdline)
    return args.net_file, args.autodev_file

def gatherNetworks(net_dict_list):
    global object_id
    nets_by_name = {}
    for net_dict in net_dict_list:
        if 'ext_conn' not in net_dict:
            net_dict['ext_conn'] = []
        if 'bandwidth' not in net_dict:
            net_dict['bandwidth'] = default_net_bandwidth[ net_dict['level'] ]
        if 'latency' not in net_dict:
            net_dict['latency'] = default_net_latency[ net_dict['level'] ]

        net_dict['number'] = object_id
        object_id += 1
        net_name = net_dict['name']
        nets_by_name[net_name] = net_dict

    return nets_by_name

def gatherRouters(router_dict_list):
    global object_id
    rtrs_by_name = {}
    for rtr_dict in router_dict_list:
        rtr_dict['number'] = object_id
        object_id += 1
        rtr_name = rtr_dict['name']
        rtrs_by_name[rtr_name] = rtr_dict 

    return rtrs_by_name

def gatherHosts(host_dict_list):
    global object_id
    hosts_by_name = {}
    for host_dict in host_dict_list:
        host_dict['number'] = object_id
        object_id += 1
        host_name = host_dict['name']
        hosts_by_name[host_name] = host_dict 

    return hosts_by_name

def createHostHomes(hosts_by_name, nets_by_name):
    # separate out hosts that have IP addresses
    withIP = []
    withoutIP = []
    
    for hn, hd in hosts_by_name.items():
        if 'IP_Addr' in hd:
            withIP.append(hd)
        else:
            withoutIP.append(hd)
 
    # cluster the hosts that have IP addresses in the same /24
    cluster = defaultdict(list)

    for host_dict in withIP:
        IP_Addr = host_dict['IP_Addr']

        # strip off the /32 if present 
        IP_Addr = IP_Addr.replace('/32','')

        octets = IP_Addr.split('.')
        cluster_name = '.'.join(octets[:3]) 
        cluster[cluster_name].append(host_dict)

    # if any one of the hosts in a cluster declares a home then they all get the same
    # home and that home gets an IP_Addr.  
    # N.B. this assumes that at most one of these has declared a home
    #
    for cn, clstr_list in cluster.items():
        home = None
        for host_dict in clstr_list: 
            if 'home' in host_dict: 
                home = host_dict['home']
                break

        if home is None:
            # create a home LAN for the host based on the host's name
            home = host_dict['name']+'-LAN'
            home_LAN = {'name': home, 'level':'LAN'}
            nets_by_name[home] = home_LAN        
              
        for host_dict in clstr_list:
            host_dict['home'] = home

# every LAN has a router to at most one WAN, named as a connection.
# If that router was not declared in the input file we create one.
# 
# If a WAN connects to a T3 it has a router.  If that router was
# not declared in the input file we create one
#
# If a T3 connects to a T2 and there is no declared router, we create one
#
# If two T2's connect and no router is declared for that connection, we make one
#
# If a T2 connects to a T1 we either find a router in the input file, or create one.
#
def createRouters(nets_by_name, rtrs_by_name):
    rtrs = {}
    # separate nets by levels
    Level = {'LAN':[],'WAN':[],'T3':[],'T2':[],'T1':[]}
    for _, net_dict in nets_by_name.items():
        Level[ net_dict['level']].append(net_dict)

    # for every connection with every LAN, see if there is already a router.
    # If not, create one
    #
    for lan in Level['LAN']:
        for wan_name in lan['ext_conn']:
            rtr_name = 'rtr-'+lan['name']+'-'+wan_name
            if rtr_name not in rtrs_by_name:
                rtr = {'name':rtr_name,'intrfc':[]}
                rtr['intrfc'].append({'faces':lan['name'], 
                    'bandwidth':default_intrfc_bandwidth['LAN']})
                    
                rtr['intrfc'].append({'faces':wan_name,
                    'bandwidth':default_intrfc_bandwidth['LAN']})

                rtrs_by_name[rtr_name] = rtr

    # for every connection with every WAN, see if there is already a router.
    # If not, create one
    #
    for wan in Level['WAN']:
        for t3_name in wan['ext_conn']:
            rtr_name = 'rtr-'+wan['name']+'-'+t3_name
            if rtr_name not in rtrs_by_name:
                rtr = {'name':rtr_name,'intrfc':[]}
                rtr['intrfc'].append({'faces':wan['name'], 
                    'bandwidth':default_intrfc_bandwidth['WAN']})
                rtr['intrfc'].append({'faces':t3_name,
                    'bandwidth':default_intrfc_bandwidth['WAN']})

                rtrs_by_name[rtr_name] = rtr

    # for every connection with every T3, see if there is already a router.
    # If not, create one
    #
    for t3 in Level['T3']:
        for t2_name in t3['ext_conn']:
            rtr_name = 'rtr-'+t3['name']+'-'+t2_name
            if rtr_name not in rtrs_by_name:
                rtr = {'name':rtr_name,'intrfc':[]}
                rtr['intrfc'].append({'faces':t3['name'], 
                    'bandwidth':default_intrfc_bandwidth['T3']})
                rtr['intrfc'].append({'faces':t2_name,
                    'bandwidth':default_intrfc_bandwidth['T3']})

                rtrs_by_name[rtr_name] = rtr

    # for every connection with every T3, see if there is already a router.
    # If not, create one
    #
    for t2 in Level['T2']:
        for conn_name in t3['ext_conn']:
            if nets_by_name[conn_name]['level'] == 'T2':
                rtr_name_1 = 'rtr-'+t2['name']+'-'+conn_name 
                rtr_name_2 = 'rtr-'+conn_name+'-'+t2['name']
                rtr_name = rtr_name_1 if rtr_name_1 < rtr_name_2 else rtr_name_2
                if rtr_name not in rtrs_by_name:
                    rtr = {'name':rtr_name,'intrfc':[]}
                    rtr['intrfc'].append({'faces':t2['name'], 
                        'bandwidth':default_intrfc_bandwidth['T2']})
                    rtr['intrfc'].append({'faces':conn_name,
                        'bandwidth':default_intrfc_bandwidth['T2']})

                    rtrs_by_name[rtr_name] = rtr

            elif nets_by_name[conn_name]['level'] == 'T1':
                rtr_name = 'rtr-'+t2['name']+'-'+conn_name
                if rtr_name not in rtrs_by_name:
                    rtr = {'name':rtr_name,'intrfc':[]}
                    rtr['intrfc'].append({'faces':t2['name'], 
                        'bandwidth':default_intrfc_bandwidth['T2']})
                    rtr['intrfc'].append({'faces':conn_name,
                        'bandwidth':default_intrfc_bandwidth['T2']})

                    rtrs_by_name[rtr_name] = rtr


def main():
    if len(sys.argv) < 2:
        print('arguments required')
        exit()

    cmdline = []
    if sys.argv[1] == '-is':
        with open(os.path.abspath(sys.argv[2]),'r', encoding='latin-1') as cf:
            for line in cf.readlines():
                line = line.strip()
                vect = line.split()
                for idx in range(0, len(vect)):
                    cmdline.append( vect[idx].strip() )
    else:
        cmdline = sys.argv[1:]

    net_file, auto_dev_file = getArgs(cmdline)
    
    with open(net_file, 'r', encoding='latin-1') as rf:
        nets_json = json.load(rf)

    # get the components of the network graph
    #
    net_dict  = gatherNetworks(nets_json['Networks'])
    if not 'Routers' in nets_json:
        nets_json['Routers'] = []
    if not 'Hosts' in nets_json:
        nets_json['Hosts'] = []

    # create the home devices for the hosts.
    # if two hosts have IP addresses and are within /24 of each other they have the same
    # home router
    host_dict = gatherHosts(nets_json['Hosts'])
    rtr_dict  = gatherRouters(nets_json['Routers'])
   
    createHostHomes(host_dict, net_dict)

    import pdb
    pdb.set_trace()
    x=1

    createRouters(net_dict, rtr_dict)

if __name__ == '__main__':
    main()
        


