// transactiontracker.go

package lib

import (
	"encoding/gob"
	"github.com/nats-io/go-nats"
	"sync"
)

// ensures transmissable types are registered for binary encoding
func init() {
	// make gob encoder aware of local types
	gob.Register(TxStatusUpdate{})
}

// track status in simple hash counter; key is transactionid, value is no. records processed
var txStatus = make(map[string]int)

// mutex to protect status hash for concurrent updates
var statusMutex = &sync.Mutex{}

// track transaction sizes to know when done
var txSize = make(map[string]int)

// mutes to protect size hash
var sizeMutex = &sync.Mutex{}

// TransactionTracker simple strucuture to capture details of
// transactions in system: tx size and tx progress
type TransactionTracker struct {
	C              *nats.EncodedConn
	ReportInterval int
}

// create a TransactionTracker
// progress interval is how many messages between status updates
// eg. report status every 500 messages
func NewTransactionTracker(report_interval int) *TransactionTracker {

	return &TransactionTracker{C: CreateNATSConnection(), ReportInterval: report_interval}

}

// message type for reporting progress
type TxStatusUpdate struct {
	TxID       string
	Message    string
	Progress   int
	Size       int
	TxComplete bool
	UIComplete bool
}

// Update progrress of transaction processing
func (tt *TransactionTracker) IncrementTracker(txID string) {

	statusMutex.Lock()
	txStatus[txID]++
	statusMutex.Unlock()

}

// set the overall size of the transaction when we know it
func (tt *TransactionTracker) SetTxSize(txID string, size int) {

	sizeMutex.Lock()
	txSize[txID] = size
	sizeMutex.Unlock()

}

func (tt *TransactionTracker) GetStatusReport(txID string) (significantChange bool, msg *NiasMessage) {

	// default assume nothing to report
	sigChange := false

	// report only if worthwhile progress, useful amount of messages
	// processed or txaction is complete.
	var progress int
	statusMutex.Lock()
	progress = txStatus[txID]
	statusMutex.Unlock()

	var size int
	sizeMutex.Lock()
	size = txSize[txID]
	sizeMutex.Unlock()

	// capture basic report details
	txu := TxStatusUpdate{TxID: txID,
		Progress: progress,
		Size:     size}

	// if progress < size but mod interval
	if (progress % tt.ReportInterval) == 0 {
		sigChange = true
	}

	// transaction is complete
	if (progress >= size) && (size > 0) {
		txu.TxComplete = true
		txu.Message = "Transaction complete."
		sigChange = true
		// notify any listeners that tx is complete
		tt.C.Publish(TRACK_TOPIC, txID)
		// remove from data stores
		removeTx(txID)
	}

	report := NiasMessage{}
	report.TxID = txID
	report.Body = txu

	// otherwise cahnge is small, no update
	return sigChange, &report

}

// release any memory associated with keeping track of the
// transaction to prevent slow resource leak
func removeTx(txID string) {

	sizeMutex.Lock()
	delete(txSize, txID)
	sizeMutex.Unlock()

	statusMutex.Lock()
	delete(txStatus, txID)
	statusMutex.Unlock()

}

// standard response to successful file upload
type IngestResponse struct {
	TxID    string
	Records int
}
