package registrar

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	//"sync"
	"time"

	helper "github.com/elastic/beats/libbeat/common/file"
	"github.com/elastic/beats/libbeat/logp"
	//"github.com/elastic/beats/libbeat/monitoring"
	"github.com/elastic/beats/libbeat/paths"

	"github.com/jcsuscriptor/elmahbeat/models"

)

type Registrar struct {

	Channel      chan []models.State
	out          successLogger
	done         chan struct{}
	registryFile string // Path to the Registry File

	states               *models.States 
	flushTimeout         time.Duration
}

type successLogger interface {
	Published(n int) bool
}

 

func New(registryFile string, flushTimeout time.Duration) (*Registrar, error) {
	r := &Registrar{
		registryFile: registryFile,
		done:         make(chan struct{}),
		states:       models.NewStates(),
		Channel:      make(chan []models.State, 1),
		flushTimeout: flushTimeout,
	}
	err := r.Init()

	return r, err
}

// Init sets up the Registrar and make sure the registry file is setup correctly
func (r *Registrar) Init() error {
	// The registry file is opened in the data path
	r.registryFile = paths.Resolve(paths.Data, r.registryFile)

	// Create directory if it does not already exist.
	registryPath := filepath.Dir(r.registryFile)
	err := os.MkdirAll(registryPath, 0750)
	if err != nil {
		return fmt.Errorf("Failed to created registry file dir %s: %v", registryPath, err)
	}

	// Check if files exists
	fileInfo, err := os.Lstat(r.registryFile)
	if os.IsNotExist(err) {
		logp.Info("No registry file found under: %s. Creating a new registry file.", r.registryFile)
		// No registry exists yet, write empty state to check if registry can be written
		return r.writeRegistry()
	}
	if err != nil {
		return err
	}

	// Check if regular file, no dir, no symlink
	if !fileInfo.Mode().IsRegular() {
		// Special error message for directory
		if fileInfo.IsDir() {
			return fmt.Errorf("Registry file path must be a file. %s is a directory.", r.registryFile)
		}
		return fmt.Errorf("Registry file path is not a regular file: %s", r.registryFile)
	}

	logp.Debug("registrar", "Registry file set to: %s", r.registryFile)

	return nil
}

 
// loadStates fetches the previous reading state from the configure RegistryFile file
// The default file is `registry` in the data path.
func (r *Registrar) loadStates() error {
	f, err := os.Open(r.registryFile)
	if err != nil {
		return err
	}

	defer f.Close()

	logp.Info("Loading registrar data from %s", r.registryFile)

	decoder := json.NewDecoder(f)
	states := []models.State{}
	err = decoder.Decode(&states)
	if err != nil {
		return fmt.Errorf("Error decoding states: %s", err)
	}


	logp.Info("States Loaded from registrar: %+v", len(states))

	return nil
}
 
func (r *Registrar) Start() error {
	// Load the previous log file locations now, for use in input
	err := r.loadStates()
	if err != nil {
		return fmt.Errorf("Error loading state: %v", err)
	}

	go r.Run()

	return nil
}

func (r *Registrar) Run() {
	logp.Debug("registrar", "Starting Registrar")
	// Writes registry on shutdown
	defer func() {
		r.writeRegistry()
		//r.wg.Done()
	}()

	var (
		timer  *time.Timer
		flushC <-chan time.Time
	)

	for {
		select {
		case <-r.done:
			logp.Info("Ending Registrar")
			return
		case <-flushC:
			flushC = nil
			timer.Stop()
			r.flushRegistry()
		case states := <-r.Channel:
			r.onEvents(states)
			if r.flushTimeout <= 0 {
				r.flushRegistry()
			} else if flushC == nil {
				timer = time.NewTimer(r.flushTimeout)
				flushC = timer.C
			}
		}
	}
}

// onEvents processes events received from the publisher pipeline
func (r *Registrar) onEvents(states []models.State) {
	 

	logp.Debug("registrar", "Registrar state updates processed. Count: %v", len(states))

	// new set of events received -> mark state registry ready for next
	// cleanup phase in case gc'able events are stored in the registry.
	//r.gcRequired = r.gcEnabled
}

// gcStates runs a registry Cleanup. The bool returned indicates wether more
// events in the registry can be gc'ed in the future.
func (r *Registrar) gcStates() {
 
}

  

// writeRegistry writes the new json registry file to disk.
func (r *Registrar) writeRegistry() error {
	r.gcStates()

	logp.Debug("registrar", "Write registry file: %s", r.registryFile)

	tempfile := r.registryFile + ".new"
	f, err := os.OpenFile(tempfile, os.O_RDWR|os.O_CREATE|os.O_TRUNC|os.O_SYNC, 0600)
	if err != nil {
		logp.Err("Failed to create tempfile (%s) for writing: %s", tempfile, err)
		return err
	}

	// First clean up states
	states := r.states.GetStates()

	encoder := json.NewEncoder(f)
	err = encoder.Encode(states)
	if err != nil {
		f.Close()
		logp.Err("Error when encoding the states: %s", err)
		return err
	}

	// Directly close file because of windows
	f.Close()

	err = helper.SafeFileRotate(r.registryFile, tempfile)

	logp.Debug("registrar", "Registry file updated. %d states written.", len(states))
	 

	return err
}
