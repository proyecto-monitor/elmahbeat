package beater

import (
	"fmt"
	"time"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"

	"github.com/jcsuscriptor/elmahbeat/config"

	//"github.com/jcsuscriptor/elmahbeat/models"
	"github.com/jcsuscriptor/elmahbeat/dao"
	"github.com/jcsuscriptor/elmahbeat/registrar"

	mgo "gopkg.in/mgo.v2"
	//"gopkg.in/mgo.v2/bson"
)

type Elmahbeat struct {
	done   chan struct{}
	config config.Config
	client beat.Client
	db *mgo.Database
}

var elmahDao = dao.ElmahDAO{}

 

// Creates beater
func New(b *beat.Beat, cfg *common.Config) (beat.Beater, error) {
	c := config.DefaultConfig
	if err := cfg.Unpack(&c); err != nil {
		return nil, fmt.Errorf("Error reading config file: %v", err)
	}

	//Connect
	elmahDao.Server = c.Url
	elmahDao.Database = c.Database
	elmahDao.Collection = c.Collection
	elmahDao.Connect()

	bt := &Elmahbeat{
		done:   make(chan struct{}),
		config: c,
	}
	return bt, nil
}

func (bt *Elmahbeat) Run(b *beat.Beat) error {
	logp.Info("elmahbeat is running! Hit CTRL-C to stop it.")

	var err error
	bt.client, err = b.Publisher.Connect()
	if err != nil {
		return err
	}

	//Registrar
	reg,err := registrar.New(bt.config.RegistryFile,time.Second)
	if err != nil {
		return err
	}
	
	ticker := time.NewTicker(bt.config.Period)
	counter := 1
	for {
		select {
		case <-bt.done:
			return nil
		case <-ticker.C:
		}

		elmahErrors, err := elmahDao.FindAll()
		if err != nil {
			return err
		}
		
		logp.Debug("ElmahBeat","Event Count %d", len(elmahErrors)) 

	    for _, item := range elmahErrors {
			
			event := beat.Event{
				Timestamp:   item.Time, // time.Now(),
				Fields: common.MapStr{
					"application":    item.Application,
					"host":    item.Host,
					"type":    item.Type,
					"message":    item.Message,	
					"source":    item.Source,	
					"detail":    item.Detail,	
					"user":    item.User,
					"statusCode":    item.StatusCode,
					"time":    item.Time,
					"webHostHtmlMessage":    item.WebHostHtmlMessage,	
					//"serverVariables":    item.ServerVariables,	
				},
			}

			bt.client.Publish(event)
		}
  
		logp.Info("Event sent %d", len(elmahErrors)) 
		counter++
	}
}

func (bt *Elmahbeat) Stop() {
	bt.client.Close()
	close(bt.done)
}
