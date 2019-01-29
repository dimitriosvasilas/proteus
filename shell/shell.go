package main

import (
	"errors"
	"flag"
	"strings"
	"time"

	"github.com/abiosoft/ishell"
	utils "github.com/dimitriosvasilas/proteus"
	pb "github.com/dimitriosvasilas/proteus/protos/qpu"
	pbQPU "github.com/dimitriosvasilas/proteus/protos/utils"
	cli "github.com/dimitriosvasilas/proteus/qpu/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type query struct {
	datatype  string
	attribute string
	lbound    string
	ubound    string
}

func find(q []query, c cli.Client) error {
	log.Debug("shell:find ", q)

	var query map[string][2]*pbQPU.Value
	lbound, ubound, err := utils.AttrBoundStrToVal(q[0].datatype, q[0].lbound, q[0].ubound)
	if err != nil {
		return errors.New("bound error")
	}
	query = map[string][2]*pbQPU.Value{q[0].attribute: {lbound, ubound}}

	log.Debug("shell:find ", query)

	return sendQuery(query, c)
}

func sendQuery(query map[string][2]*pbQPU.Value, c cli.Client) error {
	log.Debug("shell:sendQuery ", query)

	msg := make(chan *pb.QueryResultStream)
	done := make(chan bool)
	errs := make(chan error)
	errs1 := make(chan error)

	go queryConsumer(query, msg, done, errs, errs1)
	c.Find(time.Now().UnixNano(), query, msg, done, errs)
	err := <-errs1
	return err
}

func queryConsumer(query map[string][2]*pbQPU.Value, msg chan *pb.QueryResultStream, done chan bool, errs chan error, errs1 chan error) {
	for {
		if doneMsg := <-done; doneMsg {
			err := <-errs
			errs1 <- err
		}
		res := <-msg

		log.Debug("shell:queryConsumer received: ", res)

		displayResults(query, res.GetObject(), res.GetDataset())
	}
}

func displayResults(query map[string][2]*pbQPU.Value, obj *pbQPU.Object, ds *pbQPU.DataSet) {
	logMsg := log.Fields{
		"key":     obj.GetKey(),
		"dataset": ds,
	}
	for qAttr := range query {
		switch query[qAttr][0].Val.(type) {
		case *pbQPU.Value_Int:
			logMsg[qAttr] = obj.GetAttributes()[qAttr].GetInt()
		case *pbQPU.Value_Flt:
			attrK := "x-amz-meta-f-" + qAttr
			logMsg[qAttr] = obj.GetAttributes()[attrK].GetFlt()
		default:
			if qAttr == "key" {
				logMsg[qAttr] = obj.GetKey()
			} else {
				attrK := "x-amz-meta-" + qAttr
				logMsg[qAttr] = obj.GetAttributes()[attrK].GetStr()
			}
		}
	}
	log.WithFields(logMsg).Info("result")
}
func initShell(c cli.Client) {
	shell := ishell.New()
	shell.Println("QPU Shell")

	shell.AddCmd(&ishell.Cmd{
		Name: "find",
		Help: "Perform a query on object attribute",
		Func: func(ctx *ishell.Context) {
			query, err := processQueryString(ctx.Args[0])
			if err != nil {
				ctx.Err(err)
				return
			}
			err = find(query, c)
			if err != nil {
				ctx.Err(err)
				return
			}
		},
	})
	shell.Run()
}

func processQueryString(q string) ([]query, error) {
	log.Debug("shell:processQueryString: ", q)

	queryProcessed := make([]query, 0)
	predicate := strings.Split(q, "&")

	for _, p := range predicate {
		datatype := strings.Split(p, "_")
		if len(datatype) < 2 {
			return nil, errors.New("Query should have the form predicate[&predicate], where predicate=type_attrKey=lbound/ubound")
		}
		attrK := strings.Split(datatype[1], "=")
		if len(attrK) < 2 {
			return nil, errors.New("Query should have the form predicate&[predicate], where predicate=type_attrKey=lbound/ubound")
		}
		bound := strings.Split(attrK[1], "/")
		if len(bound) < 2 {
			return nil, errors.New("Query should have the form predicate&[predicate], where predicate=type_attrKey=lbound/ubound")
		}
		queryProcessed = append(queryProcessed, query{
			datatype:  datatype[0],
			attribute: attrK[0],
			lbound:    bound[0],
			ubound:    bound[1],
		})
	}

	log.Debug("shell:processQueryString: ", queryProcessed)

	return queryProcessed, nil
}

func main() {
	var endpoint string
	flag.StringVar(&endpoint, "endpoint", "noEndpoint", "QPU endpoint to send query")
	var queryIn string
	flag.StringVar(&queryIn, "query", "emptyQuery", "Query string")
	var mode string
	flag.StringVar(&mode, "mode", "noMode", "Script execution mode: cmd(command) / sh(shell) / http(http server)")
	flag.Parse()
	if endpoint == "noEndpoint" || mode == "noMode" || (mode != "cmd" && mode != "sh" && mode != "http") {
		flag.Usage()
		return
	}
	if mode == "cmd" && queryIn == "emptyQuery" {
		flag.Usage()
		return
	}
	c, conn, err := cli.NewClient(endpoint)
	defer conn.Close()

	err = viper.BindEnv("DEBUG")
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("BindEnv DEBUG failed")
	}
	debug := viper.GetBool("DEBUG")
	if debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	if err != nil {
		log.Fatal("failed to create Client %v", err)
	}
	if mode == "cmd" {
		query, err := processQueryString(queryIn)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Fatal("Find failed")
		}
		err = find(query, c)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Fatal("Find failed")
		}
	} else if mode == "sh" {
		initShell(c)
	} else if mode == "http" {
		log.WithFields(log.Fields{
			"error": errors.New("Not implemented"),
		}).Fatal("Find failed")
	}
}