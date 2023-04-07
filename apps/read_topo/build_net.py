#!/usr/bin/env python3

import sys
import os
import json
import pdb
import argparse
import random
from collections import defaultdict

default_net_Bandwidth = {'LAN':10,'WAN':100,'T3': 1000, 'T2': 10000, 'T1': 1000000}
default_net_latency = {'LAN':10,'WAN':50,'T3': 100, 'T2': 1000, 'T1': 1000}
default_intrfc_Bandwidth = {'LAN':10,'WAN':100,'T3': 1000, 'T2': 10000, 'T1': 1000000}

Networks = {}
Routers = {}
Hosts   = {}

object_id = 0
T1 = 1
T2 = 4
xT2 = 4.0
T3d = 4
WANd = 10
LANd = 10
hosts_file = None 
rtrs_file = None
store_nets_file = None 
store_topo_file = None
init_seed = '123455'


def getArgs(cmdline):
    global init_seed
    parser = argparse.ArgumentParser()

    parser.add_argument('-T1', metavar='does a T1 backbone exist', dest='T1', action='store_const', const=True, required=False)
    parser.add_argument('-T2', metavar='number of T2 networks', dest='T2', required=True)
    parser.add_argument('-xT2', metavar='average number of cross T2 connections per T2', dest='xT2', required=True)
    parser.add_argument('-T3d', metavar='density of T3 networks per T2', dest='T3d', required=True)
    parser.add_argument('-WANd', metavar='density of WAN networks per T3', dest='WANd', required=True)
    parser.add_argument('-LANd', metavar='density of LAN networks per WAN', dest='LANd', required=True)
    parser.add_argument('-hosts', metavar='file with description of hosts', dest='hosts_file', required=True)
    parser.add_argument('-rtrs', metavar='file with description of routers', dest='rtrs_file', required=False)
    parser.add_argument('-nets', metavar='output file with description of networks', dest='store_nets_file', required=False)
    parser.add_argument('-topo', metavar='output file with description of topology', dest='store_topo_file', required=False)
    parser.add_argument('-seed', metavar='RNG seed', dest='seed', required=False)

    args = parser.parse_args(cmdline)

    seed = args.seed if args.seed else init_seed
    hosts_file = None if not args.hosts_file else args.hosts_file
    rtrs_file = None if not args.rtrs_file else args.rtrs_file
    store_nets_file = None if not args.store_nets_file else args.store_nets_file
    store_topo_file = None if not args.store_topo_file else args.store_topo_file

    return args.T1, int(args.T2), float(args.xT2), float(args.T3d), float(args.WANd), float(args.LANd),\
            hosts_file, rtrs_file, store_nets_file, store_topo_file, seed

def createHostHomes(hosts_list, nets_by_name):
    # separate out hosts that have IP addresses
    withIP = []
    withoutIP = []
    
    for hd in hosts_list:
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
            if 'Home' in host_dict: 
                home = host_dict['Home']
                break

        if home is None:
            # create a home LAN for the host based on the host's name
            home = host_dict['Name']+'-LAN'
            home_LAN = {'Name': home, 'Level':'LAN'}
            nets_by_name[home] = home_LAN        
              
        for host_dict in clstr_list:
            host_dict['Home'] = home

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
def augmentRoutersList(nets_list, nets_by_name, rtrs_list):
  
    all_rtr_names = set() 
    # separate nets by levels
    Level = {'LAN':[],'WAN':[],'T3':[],'T2':[],'T1':[]}
    for net_dict in nets_list:
        Level[ net_dict['Level']].append(net_dict)

    # for every connection with every LAN, see if there is already a router.
    # If not, create one
    #
    for lan in Level['LAN']:
        for wan_name in lan['Ext_conn']:
            rtr_name = 'rtr-'+lan['Name']+'-'+wan_name
            if rtr_name not in all_rtr_names:
                rtr = {'Name':rtr_name,'Intrfc':[]}
                rtr['Intrfc'].append({'Faces':lan['Name'], 
                    'Bandwidth':default_intrfc_Bandwidth['LAN']})
                    
                rtr['Intrfc'].append({'Faces':wan_name,
                    'Bandwidth':default_intrfc_Bandwidth['LAN']})
           
                all_rtr_names.add(rtr_name) 
                rtrs_list.append(rtr)

    # for every connection with every WAN, see if there is already a router.
    # If not, create one
    #
    for wan in Level['WAN']:
        for t3_name in wan['Ext_conn']:
            rtr_name = 'rtr-'+wan['Name']+'-'+t3_name
            if rtr_name not in all_rtr_names:
                rtr = {'Name':rtr_name,'Intrfc':[]}
                rtr['Intrfc'].append({'Faces':wan['Name'], 
                    'Bandwidth':default_intrfc_Bandwidth['WAN']})
                rtr['Intrfc'].append({'Faces':t3_name,
                    'Bandwidth':default_intrfc_Bandwidth['WAN']})

                all_rtr_names.add(rtr_name)
                rtrs_list.append( rtr )

    # for every connection with every T3, see if there is already a router.
    # If not, create one
    #
    for t3 in Level['T3']:
        for t2_name in t3['Ext_conn']:
            rtr_name = 'rtr-'+t3['Name']+'-'+t2_name
            if rtr_name not in all_rtr_names:
                rtr = {'Name':rtr_name,'Intrfc':[]}
                rtr['Intrfc'].append({'Faces':t3['Name'], 
                    'Bandwidth':default_intrfc_Bandwidth['T3']})
                rtr['Intrfc'].append({'Faces':t2_name,
                    'Bandwidth':default_intrfc_Bandwidth['T3']})

                all_rtr_names.add(rtr_name)
                rtrs_list.append(rtr)

    # for every connection with every T3, see if there is already a router.
    # If not, create one
    #
    for t2 in Level['T2']:
        for conn_name in t3['Ext_conn']:
            if nets_by_name[conn_name]['Level'] == 'T2':
                rtr_name_1 = 'rtr-'+t2['Name']+'-'+conn_name 
                rtr_name_2 = 'rtr-'+conn_name+'-'+t2['Name']
                rtr_name = rtr_name_1 if rtr_name_1 < rtr_name_2 else rtr_name_2
                if rtr_name not in all_rtr_names:
                    rtr = {'Name':rtr_name,'Intrfc':[]}
                    rtr['Intrfc'].append({'Faces':t2['Name'], 
                        'Bandwidth':default_intrfc_Bandwidth['T2']})
                    rtr['Intrfc'].append({'Faces':conn_name,
                        'Bandwidth':default_intrfc_Bandwidth['T2']})

                    all_rtr_names.add(rtr_name)
                    rtrs_list.append(rtr)

            elif nets_by_name[conn_name]['Level'] == 'T1':
                rtr_name = 'rtr-'+t2['Name']+'-'+conn_name
                if rtr_name not in all_rtr_names:
                    rtr = {'Name':rtr_name,'Intrfc':[]}
                    rtr['Intrfc'].append({'Faces':t2['Name'], 
                        'Bandwidth':default_intrfc_Bandwidth['T2']})
                    rtr['Intrfc'].append({'Faces':conn_name,
                        'Bandwidth':default_intrfc_Bandwidth['T2']})

                    all_rtr_names.add(rtr_name)
                    rtrs_list.append( rtr )


# net_dict_list is a list of dictionaries, each of which
# describes a network, with attributes
#       name (string)
#       Ext_conn (list of strings)
#       level (string)
#       Bandwidth (float)
#       latency (float)
#       number (int)
#
# from this create and return a dictionary with highest level attributes
# 'Networks', 'Routers', 'Hosts'
#  Each returns a list of object dictionaries of the given type
#
#  The hosts are carried in from expression in the hosts_file
#
#  The list of routers is initialized by whatever is in rtrs_file (if not empty),
#  and then computed from interactions between networks and hosts
#
#
def buildTopology(net_dict_list, rtrs_file, hosts_file):
    rtn_dict = {'Networks':net_dict_list, 'Routers':[], 'Hosts':[]}

    if rtrs_file:
        with open(rtrs_file,'r') as rf:
            rtn_dict['Routers'] = json.load(rf)

    if hosts_file:
        with open(hosts_file,'r') as rf:
            rtn_dict['Hosts'] = json.load(rf)

    nets_by_name = {}
    for net_dict in net_dict_list:
        nets_by_name[net_dict['Name']] = net_dict

    createHostHomes(rtn_dict['Networks'], nets_by_name)

    # build up routers from network topology
    #
    augmentRoutersList(net_dict_list, nets_by_name, rtn_dict['Routers'])

    # number all the objects
    object_id = 0
    for net_dict in rtn_dict['Networks']:
        net_dict['Number'] = object_id
        object_id += 1

    for rtr_dict in rtn_dict['Routers']:
        rtr_dict['Number'] = object_id
        object_id += 1

    for host_dict in rtn_dict['Hosts']:
        host_dict['Number'] = object_id
        object_id += 1

    return rtn_dict
    

def main():
    global T1, T2, xT2, T3, WAN, LAN, seed, hosts_file, store_nets_file
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

    T1_exists, T2n, xT2, T3d, WANd, LANd, hosts_file, rtrs_file, store_nets_file, store_topo_file, seed = getArgs(cmdline)

    random.seed(seed)

    if T1:
        T1 = {'Name':'Backbone', 'Level':'T1'}

    T2_list = []
    for idx in range(0,int(T2n)):
        T2_list.append({'Name':'T2-'+str(idx),'Ext_conn':['Backbone'],'Level':'T2',
                        'Bandwidth':default_net_Bandwidth['T2'],\
                        'Latency':default_net_latency['T2']})

    # from the average number of xT2 connections per T2 compute
    # what the probability of a given edge being present is.  
    #
    # number of T2 networks
    L2 = len(T2_list)

    # average connections per T2 is probability a different T2
    # is a connection, times the number of different T2's
    #
    # xT2 = p*(L2-1)
    # so
    # p = xT2/(L2-1)
    #
    prxT2 = float(xT2)/(L2-1)

    for idx in range(0, L2):
        for jdx in range(idx+1, L2):    
            p = random.random()
            if p < prxT2:
                T2_list[idx]['Ext_conn'].append( T2_list[jdx]['Name'] )
                T2_list[jdx]['Ext_conn'].append( T2_list[idx]['Name'] )
 

    # T3d is the ratio of number of T3 networks to T2 networks
    T3_list = []
    for idx in range(0, int(L2*T3d + 1)):
        T3_list.append({'Name':'T3-'+str(idx), 'Ext_conn':[], 'Level':'T3',\
            'Bandwidth':default_net_Bandwidth['T3'],
            'Latency':default_net_latency['T3']})
    
    L3 = len(T3_list)
    for idx in range(0, L3):
        # choose a random T2 and its reflection as the T2 points of connection
        if L2 > 1:
            T2_base = random.randint(0,L2-1)
            offset = L2
            while (T2_base+offset)%L2 == T2_base:
                offset += 1

            T2_reflect = (T2_base+offset)%L2
            T3_list[idx]['Ext_conn'].append(T2_list[T2_base]['Name'])
            T3_list[idx]['Ext_conn'].append(T2_list[T2_reflect]['Name'])
        else:
            T3_list[idx]['Ext_conn'].append(T2_list[0]['Name'])

    # WANd is the ratio of number of WAN networks to T3 networks
    WAN_list = []
    for idx in range(0, int(L3*WANd + 1)):
        WAN_list.append({'Name':'WAN-'+str(idx), 'Ext_conn':[], 'Level':'WAN',\
            'Bandwidth':default_net_Bandwidth['WAN'],
            'Latency':default_net_latency['WAN']})
   
    # A WAN typically connects to one T3, choose one modulo index identity
    LWAN = len(WAN_list)
    for idx in range(0, LWAN):
        T3id = idx%L3 
        WAN_list[idx]['Ext_conn'].append(T3_list[T3id]['Name'])

    # LANd is the ratio of number of LAN networks to WAN networks
    LAN_list = []
    for idx in range(0, int(LWAN*LANd + 1)):
        LAN_list.append({'Name':'LAN-'+str(idx), 'Ext_conn':[], 'Level':'LAN',\
            'Bandwidth':default_net_Bandwidth['LAN'],
            'Latency':default_net_latency['LAN']})
   
    # A WAN typically connects to one T3, choose one modulo index identity
    LLAN = len(LAN_list)
    for idx in range(0, LLAN):
        WANid = idx%LWAN 
        LAN_list[idx]['Ext_conn'].append(WAN_list[WANid]['Name'])

    nets_json = {}
    nets_json['Networks'] = []
    if T1:
        nets_json['Networks'].append(T1)
    nets_json['Networks'].extend(T2_list)
    nets_json['Networks'].extend(T3_list)
    nets_json['Networks'].extend(WAN_list)
    nets_json['Networks'].extend(LAN_list)

    if store_nets_file is not None:
        with open(store_nets_file,'w') as wf:
            json.dump(nets_json, wf, sort_keys=True, indent=4)

    topo_dict = buildTopology(nets_json['Networks'], rtrs_file, hosts_file)

    if store_topo_file is not None:
        with open(store_topo_file,'w') as wf:
            json.dump(topo_dict, wf, sort_keys=True, indent=4)

if __name__ == '__main__':
    main()
        


