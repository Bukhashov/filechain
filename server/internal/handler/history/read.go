package history

import (
	"context"
)

func(h *history) ReadAll(ctx context.Context)(err error){

	// s := NewStorage(h.client, h.logger)
	// s.Read(ctx)

	return nil
}

func(h *history) ReadOne(ctx context.Context, model *HistoryModel)(err error){
	s := NewStorage(h.client, h.logger, h.table)
	if err = s.ReadOne(ctx, model); err != nil {
		return err
	}
	return nil
}