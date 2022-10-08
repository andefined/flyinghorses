package cell

import (
	"context"
	"strconv"
	"strings"

	cellv1 "github.com/andefined/flyinghorses/internal/flyinghorses/cell/v1"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// CellMeasurementService
type CellMeasurementService struct {
	log          *zap.SugaredLogger
	ctx          context.Context
	TopicChannel *redis.PubSub
	ErrorChannel chan error
}

// NewCellMeasurementSerice
func NewCellMeasurementSerice(log *zap.SugaredLogger, ctx context.Context, out *redis.PubSub, err chan error) *CellMeasurementService {
	return &CellMeasurementService{
		log, ctx, out, err,
	}
}

// Consume
func (s *CellMeasurementService) Consume() {
	for {
		msg, err := s.TopicChannel.ReceiveMessage(s.ctx)
		if err != nil {
			s.ErrorChannel <- err
			continue
		}

		s.Produce(msg.String())
	}
}

// Produce
func (s *CellMeasurementService) Produce(msg string) {
	replacer := strings.NewReplacer("<", "", ">", "", "Messagecell: ", "")
	cellString := replacer.Replace(msg)
	cellData := strings.Split(cellString, ",")
	if len(cellData) == 17 {
		mcc, _ := strconv.ParseInt(cellData[0], 10, 32)        // tower.mcc
		mnc, _ := strconv.ParseInt(cellData[1], 10, 32)        // tower.mnc
		tac, _ := strconv.ParseInt(cellData[2], 10, 32)        // tower.tac
		cid, _ := strconv.ParseInt(cellData[3], 10, 32)        // tower.cid
		phyid, _ := strconv.ParseInt(cellData[4], 10, 32)      // tower.phyid
		earfcn, _ := strconv.ParseInt(cellData[5], 10, 32)     // tower.earfcn
		rssi, _ := strconv.ParseFloat(cellData[6], 32)         // tower.rssi
		frequency, _ := strconv.ParseFloat(cellData[7], 32)    // tower.frequency
		enodeb_id, _ := strconv.ParseInt(cellData[8], 10, 32)  // tower.enodeb_id
		sector_id, _ := strconv.ParseInt(cellData[9], 10, 32)  // tower.sector_id
		cfo, _ := strconv.ParseFloat(cellData[10], 32)         // tower.cfo
		rsrq, _ := strconv.ParseFloat(cellData[11], 32)        // tower.rsrq
		snr, _ := strconv.ParseFloat(cellData[12], 32)         // tower.snr
		rsrp, _ := strconv.ParseFloat(cellData[13], 32)        // tower.rsrp
		tx_pwr, _ := strconv.ParseFloat(cellData[14], 32)      // tower.tx_pwr
		raw_sib1 := cellData[15]                               // tower.raw_sib1
		timestamp, _ := strconv.ParseInt(cellData[16], 10, 32) // seconds

		cell := &cellv1.Cell{
			Id:        uuid.New().String(),
			Mcc:       mcc,
			Mnc:       mnc,
			Tac:       tac,
			Cid:       cid,
			Phyid:     phyid,
			Earfcn:    earfcn,
			Rssi:      rssi,
			Frequency: frequency,
			EnodebId:  enodeb_id,
			SectorId:  sector_id,
			Cfo:       cfo,
			Rsrq:      rsrq,
			Snr:       snr,
			Rsrp:      rsrp,
			TxPwr:     tx_pwr,
			RawSib1:   raw_sib1,
			Timestamp: timestamp,
		}

		s.log.Infof("New Cell: %v", cell)
	} else {
		s.log.Debugf("Malformed Cell Meaasurement: %v", msg)
	}
}
