package api

import (
	"fmt"
	"time"

	"github.com/ipfs/go-cid"
	"github.com/textileio/powergate/ffs"
)

type DealID uint64

func (i *API) ImportStorage(payloadCid cid.Cid, dealIDs []DealID) error {
	i.lock.Lock()
	defer i.lock.Unlock()

	// Create a storage config based on the default one but with
	// hot storage disabled.
	scfg := ffs.StorageConfig{
		Repairable: false,
		Hot: ffs.HotConfig{
			Enabled: false,
		},
		Cold: i.cfg.DefaultStorageConfig.Cold,
	}

	if len(dealIDs) == 0 {
		return fmt.Errorf("deal ids list is empty")
	}

	// - Get Deal 1 information.
	// - Iterate other deals and check are for the same PieceCid

	filStorage := make([]ffs.FilStorage, len(deals))
	for i, d := range deals {
		filStorage[i] = ffs.FilStorage{
			PieceCid: pieceCid,
			Miner:    d.MinerAddress,
		}
		if d.ProposalCid != nil {
			filStorage[i].ProposalCid = *d.ProposalCid
		}
	}

	cinfo := ffs.StorageInfo{
		JobID:   ffs.EmptyJobID,
		Cid:     payloadCid,
		Created: time.Now(),
		Hot:     ffs.HotInfo{Enabled: false},
		Cold: ffs.ColdInfo{
			Enabled: true,
			Filecoin: ffs.FilInfo{
				DataCid:   payloadCid,
				Proposals: filStorage,
			},
		},
	}

	if err := i.sched.ImportStorageInfo(cinfo); err != nil {
		return fmt.Errorf("importing cid info in scheduler: %s", err)
	}

	if err := i.is.putStorageConfig(payloadCid, scfg); err != nil {
		return fmt.Errorf("saving new imported config: %s", err)
	}

	return nil
}
