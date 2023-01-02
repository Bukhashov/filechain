package history

import (
	"context"
	"bytes"
	"crypto/sha256"
	"strconv"
	"fmt"
)

func(h *HistoryModel) setHash() {
	timeStamp := []byte(strconv.FormatInt(h.TimeStamp, 10))
	headers := bytes.Join([][]byte{h.ProveHash, h.Target, timeStamp}, []byte{})
	hash := sha256.Sum256(headers)

	h.Hash = hash[:]
}

func (h *history) Genesis(ctx context.Context, )(err error) {
	forHash := fmt.Sprintf("%s%s", *h.user, *h.target)
	model := HistoryModel{
		Hash: 		[]byte(forHash),
		TimeStamp: 	h.t,
		Addres: 	*h.address,
		ProveHash: 	[]byte{},
		User: 		[]byte(*h.user),
		Target: 	[]byte(*h.target),
	}
	
	model.setHash()

	s := NewStorage(h.client, h.logger, h.table)
	if err = s.Write(ctx, &model); err != nil {
		return err
	}

	return nil
}