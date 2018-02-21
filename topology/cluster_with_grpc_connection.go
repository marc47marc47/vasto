package topology

import (
	"fmt"
	"github.com/chrislusf/glog"
	"github.com/chrislusf/vasto/pb"
	"google.golang.org/grpc"
)

func (cluster *Cluster) WithConnection(name string, serverId int, fn func(*pb.ClusterNode, *grpc.ClientConn) error) error {

	node, _, ok := cluster.GetNode(serverId)

	if !ok {
		glog.Errorf("cluster misses server %d: %+v", serverId, cluster.String())
		return fmt.Errorf("server %d not found", serverId)
	}

	return doWithConnect(name, node, serverId, fn)
}

type PrimaryShards []*pb.ClusterNode

func (nodes PrimaryShards) WithConnection(name string, serverId int, fn func(*pb.ClusterNode, *grpc.ClientConn) error) error {

	if serverId < 0 || serverId >= len(nodes) {
		return fmt.Errorf("server %d not found in %d servers: %+v", serverId, len(nodes), nodes)
	}

	node := nodes[serverId]

	return doWithConnect(name, node, serverId, fn)

}

func doWithConnect(name string, node *pb.ClusterNode, serverId int, fn func(*pb.ClusterNode, *grpc.ClientConn) error) error {

	if node == nil {
		return fmt.Errorf("%s: server %d is missing", name, serverId)
	}

	// glog.V(2).Infof("connecting to server %d at %s", serverId, node.GetAdminAddress())

	grpcConnection, err := grpc.Dial(node.StoreResource.AdminAddress, grpc.WithInsecure())
	if err != nil {
		return fmt.Errorf("%s: fail to dial %s: %v", name, node.StoreResource.AdminAddress, err)
	}
	defer grpcConnection.Close()

	// glog.V(2).Infof("%s: connect to shard %s on %s", name, node.ShardInfo.IdentifierOnThisServer(), node.StoreResource.AdminAddress)

	return fn(node, grpcConnection)
}
