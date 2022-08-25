package monitor

import (
	"context"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	cellv1 "github.com/andefined/flyinghorses/internal/flyinghorses/cell/v1"
	"github.com/andefined/flyinghorses/pkg/config"
	"github.com/andefined/flyinghorses/pkg/logger"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// CellMeasurementService
type CellMeasurementService struct {
	log           *zap.SugaredLogger
	Command       string
	Args          string
	OutputChannel chan string
	ErrorChannel  chan error
}

// NewCellMeasurementSerice
func NewCellMeasurementSerice(log *zap.SugaredLogger, cmd string, args string, out chan string, err chan error) *CellMeasurementService {
	return &CellMeasurementService{
		log, cmd, args, out, err,
	}
}

// Exec
func (s *CellMeasurementService) Exec() {
	// Create the command
	cmd := exec.Command(s.Command, strings.Split(s.Args, " ")...)
	// attach stderr to stdout for combined results
	cmd.Stderr = cmd.Stdout
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		s.ErrorChannel <- err
		return
	}

	if err := cmd.Start(); err != nil {
		s.ErrorChannel <- err
		return
	}

	for {
		b := make([]byte, 1024)
		_, err := stdout.Read(b)

		if strings.Contains(string(b), "**** sending packet: ") {
			replacer := strings.NewReplacer("<", "", ">", "")
			msg := replacer.Replace(strings.Split(string(b), "**** sending packet: ")[1])
			s.OutputChannel <- strings.TrimSpace(strings.Split(msg, "set")[0])
		}

		if err != nil {
			s.ErrorChannel <- err
			break
		}
	}
}

// Consume
func (s *CellMeasurementService) Consume() {
	for {
		msg := <-s.OutputChannel
		if msg != "" {
			s.log.Debugf("Cell measurement: %s", msg)
			s.Produce(msg)
		}
	}
}

// Produce
// 202,10,40700,2617101,164,3350,
// -52.9655,2680,10223,13,-0.383168,
// -14.5692,-5.42715,-18.314,8.50311,
// 48b8183c972c9c7ed16210001ea81c543f22aa91d48719adf145,1661453161
// os << tower.mcc << "," << tower.mnc << "," << tower.tac << "," << tower.cid << "," << tower.phyid << ","
// << tower.earfcn << "," << tower.rssi << "," << tower.frequency << "," << tower.enodeb_id << "," << tower.sector_id
// << "," << tower.cfo << "," << tower.rsrq << "," << tower.snr << "," << tower.rsrp << "," << tower.tx_pwr << ","
// << tower.raw_sib1 << "," << seconds;
func (s *CellMeasurementService) Produce(msg string) {
	cellData := strings.Split(msg, ",")
	if len(cellData) == 17 {
		mcc, _ := strconv.ParseInt(cellData[0], 10, 32)
		mnc, _ := strconv.ParseInt(cellData[1], 10, 32)
		tac, _ := strconv.ParseInt(cellData[2], 10, 32)
		cid, _ := strconv.ParseInt(cellData[3], 10, 32)
		phyid, _ := strconv.ParseInt(cellData[4], 10, 32)
		earfcn, _ := strconv.ParseInt(cellData[5], 10, 32)
		rssi, _ := strconv.ParseFloat(cellData[6], 32)
		frequency, _ := strconv.ParseFloat(cellData[7], 32)
		enodeb_id, _ := strconv.ParseInt(cellData[8], 10, 32)
		sector_id, _ := strconv.ParseInt(cellData[9], 10, 32)
		cfo, _ := strconv.ParseFloat(cellData[10], 32)
		rsrq, _ := strconv.ParseFloat(cellData[11], 32)
		snr, _ := strconv.ParseFloat(cellData[12], 32)
		rsrp, _ := strconv.ParseFloat(cellData[13], 32)
		tx_pwr, _ := strconv.ParseFloat(cellData[14], 32)
		est_dist, _ := strconv.ParseFloat(cellData[15], 32)
		timestamp, _ := strconv.ParseInt(cellData[16], 10, 32)

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
			EstDist:   est_dist,
			RawSib1:   cellData[16],
			Timestamp: timestamp,
		}

		s.log.Infof("New Cell: %v", cell)
	} else {
		s.log.Debugf("Malformed Cell Meaasurement: %v", msg)
	}

}

// NewMonitorService
func NewMonitorService(ctx context.Context, cfg *config.Config) error {
	// ** LOGGER
	// Create a reusable zap logger
	log := logger.NewLogger(cfg.Env, cfg.Log.Level, cfg.Log.Path)
	log.Info("Starting srsLTE->cell_measurement monitoring service")

	// Create the output channel
	outputChannel := make(chan string)
	defer close(outputChannel)

	// Create the error channel
	errorChannel := make(chan error, 1)
	defer close(errorChannel)

	// Listen for os signals
	osSignals := make(chan os.Signal, 1)
	defer close(osSignals)

	// ** TERMINATION
	// Listen for manual termination
	signal.Notify(osSignals, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP, syscall.SIGQUIT)

	cellMeasurementService := NewCellMeasurementSerice(log, cfg.SRS.Command, cfg.GetSRSCommandArgs(), outputChannel, errorChannel)

	go cellMeasurementService.Consume()
	go cellMeasurementService.Exec()

	select {
	case err := <-errorChannel:
		log.Errorf("srsLTE->cell_measurement error: %s", err.Error())
		return err
	case signal := <-osSignals:
		log.Fatalf("srsLTE->cell_measurement shutdown signal: %s", signal)
	}

	return nil
}
