package store

import (
	"github.com/chrislusf/vasto/pb"
	"github.com/chrislusf/vasto/storage/change_log"
	"time"
)

func (ss *storeServer) processDelete(deleteRequest *pb.DeleteRequest) *pb.DeleteResponse {

	resp := &pb.DeleteResponse{
		Ok: true,
	}
	err := ss.nodes[0].db.Delete(deleteRequest.Key)
	if err != nil {
		resp.Ok = false
		resp.Status = err.Error()
	} else {
		ss.logDelete(deleteRequest.Key, deleteRequest.PartitionHash, uint64(time.Now().UnixNano()))
	}
	return resp

}

func (ss *storeServer) logDelete(key []byte, partitionHash uint64, updatedAtNs uint64) {

	if ss.nodes[0].lm == nil {
		return
	}

	ss.nodes[0].lm.AppendEntry(change_log.NewLogEntry(
		partitionHash,
		updatedAtNs,
		0,
		true,
		key,
		nil,
	))

}