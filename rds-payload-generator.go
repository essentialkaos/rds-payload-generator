package main

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                     Copyright (c) 2009-2019 ESSENTIAL KAOS                         //
//        Essential Kaos Open Source License <https://essentialkaos.com/ekol>         //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"pkg.re/essentialkaos/ek.v11/fmtc"
	"pkg.re/essentialkaos/ek.v11/fmtutil"
	"pkg.re/essentialkaos/ek.v11/fsutil"
	"pkg.re/essentialkaos/ek.v11/options"
	"pkg.re/essentialkaos/ek.v11/rand"
	"pkg.re/essentialkaos/ek.v11/usage"

	"pkg.re/essentialkaos/redy.v4"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// Application info
const (
	APP  = "RDS Payload Generator"
	VER  = "1.1.0"
	DESC = "Payload generator for Redis-Split"
)

// Supported command line options
const (
	OPT_DIR      = "d:dir"
	OPT_KEYS     = "k:keys"
	OPT_RATIO    = "r:ratio"
	OPT_NO_COLOR = "nc:no-color"
	OPT_HELP     = "h:help"
	OPT_VER      = "v:version"
)

// START_PORT start port
const START_PORT = 63000

// ////////////////////////////////////////////////////////////////////////////////// //

// RedisStore is Redis connection pull
type RedisStore struct {
	clients map[string]*redy.Client
}

// ////////////////////////////////////////////////////////////////////////////////// //

var optMap = options.Map{
	OPT_DIR:      {Value: "/opt/redis-split"},
	OPT_KEYS:     {Type: options.INT, Value: 5000, Min: 10, Max: 1000000},
	OPT_RATIO:    {Type: options.INT, Value: 4, Min: 1, Max: 100},
	OPT_NO_COLOR: {Type: options.BOOL},
	OPT_HELP:     {Type: options.BOOL, Alias: "u:usage"},
	OPT_VER:      {Type: options.BOOL, Alias: "ver"},
}

// ////////////////////////////////////////////////////////////////////////////////// //

// main is main func
func main() {
	_, errs := options.Parse(optMap)

	if len(errs) != 0 {
		for _, err := range errs {
			printError(err.Error())
		}

		os.Exit(1)
	}

	if options.GetB(OPT_NO_COLOR) {
		fmtc.DisableColors = true
	}

	if options.GetB(OPT_VER) {
		showAbout()
		return
	}

	if options.GetB(OPT_HELP) {
		showUsage()
		return
	}

	checkRDSInstallation()

	fmtc.TPrintf("{s-}Starting…{!}")

	generatePayload()
}

// checkRDSInstallation checks Redis-Split installation
func checkRDSInstallation() {
	rdsDir := options.GetS(OPT_DIR)
	metaDir := rdsDir + "/meta"

	if !fsutil.IsExist(rdsDir) {
		printError("Directory %s doesn't exist", rdsDir)
		os.Exit(1)
	}

	if !fsutil.IsExist(metaDir) {
		printError("Directory %s doesn't exist", metaDir)
		os.Exit(1)
	}

	if !fsutil.IsDir(rdsDir) {
		printError("%s is not a directory", metaDir)
		os.Exit(1)
	}

	if !fsutil.IsDir(metaDir) {
		printError("%s is not a directory", metaDir)
		os.Exit(1)
	}

	if fsutil.IsEmptyDir(metaDir) {
		printError("No instances are created")
		os.Exit(1)
	}
}

// generatePayload generates payload to instances
func generatePayload() {
	var lastIDListUpdate time.Time
	var ids []string
	var num int

	var reads, writes int64

	store := &RedisStore{make(map[string]*redy.Client)}
	metaDir := options.GetS(OPT_DIR) + "/meta"
	maxKey := options.GetI(OPT_KEYS)
	ratio := options.GetI(OPT_RATIO)

	for {
		if time.Since(lastIDListUpdate) >= 5*time.Minute {
			ids = fsutil.List(metaDir, false)
			num = len(ids)
			lastIDListUpdate = time.Now()
		}

		time.Sleep(getPause())

		instanceID := ids[rand.Int(num)]

		if !isInstanceWorks(instanceID) {
			store.Remove(instanceID)
			continue
		}

		client := store.Get(instanceID)
		key := getKey(maxKey)

		switch rand.Int(ratio) {
		case 0:
			client.Cmd("SET", key)
			writes++
		default:
			client.Cmd("GET", key)
			reads++
		}

		fmtc.TPrintf(
			"{s}[{!} {c}↑ %s{!} {s}|{!} {m}↓ %s{!} {s}]{!}",
			fmtutil.PrettyNum(writes),
			fmtutil.PrettyNum(reads),
		)
	}
}

// getPause returns pause between requests
func getPause() time.Duration {
	r := 0.001 * float64(rand.Int(25))
	return time.Duration(r * float64(time.Second))
}

// getKey returns key name with random suffix
func getKey(max int) string {
	return "KEY" + strconv.Itoa(rand.Int(max))
}

// isInstanceWorks returns true if instance is works
func isInstanceWorks(id string) bool {
	pidDir := options.GetS(OPT_DIR) + "/pid"
	pidFile := fmt.Sprintf("%s/%s.pid", pidDir, id)

	return fsutil.IsExist(pidFile)
}

// printError prints error message to console
func printError(f string, a ...interface{}) {
	fmtc.Fprintf(os.Stderr, "{r}"+f+"{!}\n", a...)
}

// ////////////////////////////////////////////////////////////////////////////////// //

// Get returns client for given ID from store
func (rs *RedisStore) Get(id string) *redy.Client {
	client := rs.clients[id]

	if client != nil {
		return client
	}

	idInt, _ := strconv.Atoi(id)
	port := START_PORT + idInt

	client = &redy.Client{
		Addr:         "127.0.0.1:" + strconv.Itoa(port),
		DialTimeout:  time.Second,
		ReadTimeout:  time.Second,
		WriteTimeout: time.Second,
	}

	client.Connect()

	rs.clients[id] = client

	return client
}

// Remove removes client for given ID from store
func (rs *RedisStore) Remove(id string) {
	if rs.clients[id] == nil {
		return
	}

	rs.clients[id].Close()

	delete(rs.clients, id)
}

// ////////////////////////////////////////////////////////////////////////////////// //

// showUsage generates usage info
func showUsage() {
	info := usage.NewInfo()

	info.AddOption(OPT_DIR, "Redis-Split main dir", "dir")
	info.AddOption(OPT_KEYS, "Number of keys {s-}(10-1000000 default: 5000){!}")
	info.AddOption(OPT_RATIO, "Writes/reads ration {s-}(1-100 default: 4){!}")
	info.AddOption(OPT_NO_COLOR, "Disable colors in output")
	info.AddOption(OPT_HELP, "Show this help message")
	info.AddOption(OPT_VER, "Show version")

	info.AddExample(
		"-d /srv/redis-split -k 35000 -r 10",
		"Run tool with custom settings",
	)

	info.Render()
}

// showAbout print info about version
func showAbout() {
	about := &usage.About{
		App:     APP,
		Version: VER,
		Desc:    DESC,
		Year:    2006,
		Owner:   "ESSENTIAL KAOS",
		License: "Essential Kaos Open Source License <https://essentialkaos.com/ekol>",
	}

	about.Render()
}
