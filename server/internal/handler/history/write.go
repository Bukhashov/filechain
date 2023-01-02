package history

import (
	"context"
	"fmt"
)

func (h *history) Write(ctx context.Context)(err error){
	
	s := NewStorage(h.client, h.logger, h.table)

	readModel := HistoryModel{
		Addres: *h.address,
	}

	if err = s.ReadOne(ctx, &readModel); err != nil {
		h.logger.Info("wirte histore err not fund prove hash")
		return err
	}

	Hash := fmt.Sprintf("%s%s", *h.user, *h.target)

	writeModel := HistoryModel{
		Hash: 		[]byte(Hash),
		TimeStamp:	h.t,
		Addres: 	*h.address,
		ProveHash: 	readModel.Hash,
		User: 		[]byte(*h.user),
		Target: 	[]byte(*h.target),
	}

	writeModel.setHash()

	if err = s.Write(ctx, &writeModel); err != nil {
		return err
	}
	
	return nil
}