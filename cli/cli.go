package cli

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2023 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/essentialkaos/ek/v12/fmtc"
	"github.com/essentialkaos/ek/v12/fmtutil"
	"github.com/essentialkaos/ek/v12/fsutil"
	"github.com/essentialkaos/ek/v12/options"
	"github.com/essentialkaos/ek/v12/rand"
	"github.com/essentialkaos/ek/v12/usage"
	"github.com/essentialkaos/ek/v12/usage/completion/bash"
	"github.com/essentialkaos/ek/v12/usage/completion/fish"
	"github.com/essentialkaos/ek/v12/usage/completion/zsh"
	"github.com/essentialkaos/ek/v12/usage/man"

	"github.com/essentialkaos/redy/v4"

	"github.com/essentialkaos/rds-payload-generator/cli/support"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// Application info
const (
	APP  = "RDS Payload Generator"
	VER  = "1.2.0"
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

	OPT_VERB_VER     = "vv:verbose-version"
	OPT_COMPLETION   = "completion"
	OPT_GENERATE_MAN = "generate-man"
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
	OPT_DIR:      {},
	OPT_KEYS:     {Type: options.INT, Value: 5000, Min: 10, Max: 1000000},
	OPT_RATIO:    {Type: options.INT, Value: 4, Min: 1, Max: 100},
	OPT_NO_COLOR: {Type: options.BOOL},
	OPT_HELP:     {Type: options.BOOL},
	OPT_VER:      {Type: options.BOOL},

	OPT_VERB_VER:     {Type: options.BOOL},
	OPT_COMPLETION:   {},
	OPT_GENERATE_MAN: {Type: options.BOOL},
}

// ////////////////////////////////////////////////////////////////////////////////// //

// Run is main utility function
func Run(gitRev string, gomod []byte) {
	preConfigureUI()

	_, errs := options.Parse(optMap)

	if len(errs) != 0 {
		for _, err := range errs {
			printError(err.Error())
		}

		os.Exit(1)
	}

	configureUI()

	switch {
	case options.Has(OPT_COMPLETION):
		os.Exit(printCompletion())
	case options.Has(OPT_GENERATE_MAN):
		printMan()
		os.Exit(0)
	case options.GetB(OPT_VER):
		genAbout(gitRev).Print()
		os.Exit(0)
	case options.GetB(OPT_VERB_VER):
		support.Print(APP, VER, gitRev, gomod)
		os.Exit(0)
	case options.GetB(OPT_HELP):
		genUsage().Print()
		os.Exit(0)
	}

	checkRDSInstallation()
	generatePayload()
}

// preConfigureUI preconfigures UI based on information about user terminal
func preConfigureUI() {
	term := os.Getenv("TERM")

	fmtc.DisableColors = true

	if term != "" {
		switch {
		case strings.Contains(term, "xterm"),
			strings.Contains(term, "color"),
			term == "screen":
			fmtc.DisableColors = false
		}
	}

	if !fsutil.IsCharacterDevice("/dev/stdout") && os.Getenv("FAKETTY") == "" {
		fmtc.DisableColors = true
	}

	if os.Getenv("NO_COLOR") != "" {
		fmtc.DisableColors = true
	}
}

// configureUI configures user interface
func configureUI() {
	if options.GetB(OPT_NO_COLOR) {
		fmtc.DisableColors = true
	}
}

// checkRDSInstallation checks Redis-Split installation
func checkRDSInstallation() {
	rdsDir := getRDSMainDir()
	err := fsutil.ValidatePerms("DRX", rdsDir)

	if err != nil {
		printError(err.Error())
		os.Exit(1)
	}

	metaDir := rdsDir + "/meta"
	err = fsutil.ValidatePerms("DRX", metaDir)

	if err != nil {
		printError(err.Error())
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
	metaDir := getRDSMainDir() + "/meta"
	maxKey := options.GetI(OPT_KEYS)
	ratio := options.GetI(OPT_RATIO)

	fmtc.TPrintf("{s-}Starting…{!}")

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

// getRDSMainDir returns path to main Redis-Split directory
func getRDSMainDir() string {
	return fsutil.ProperPath("DRX",
		[]string{
			options.GetS(OPT_DIR),
			"/opt/redis-split",
			"/srv/redis-split",
			"/srv2/redis-split",
			"/srv3/redis-split",
			"/srv4/redis-split",
		},
	)
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
	pidDir := getRDSMainDir() + "/pid"
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

// printCompletion prints completion for given shell
func printCompletion() int {
	info := genUsage()

	switch options.GetS(OPT_COMPLETION) {
	case "bash":
		fmt.Printf(bash.Generate(info, "rds-payload-generator"))
	case "fish":
		fmt.Printf(fish.Generate(info, "rds-payload-generator"))
	case "zsh":
		fmt.Printf(zsh.Generate(info, optMap, "rds-payload-generator"))
	default:
		return 1
	}

	return 0
}

// printMan prints man page
func printMan() {
	fmt.Println(
		man.Generate(
			genUsage(),
			genAbout(""),
		),
	)
}

// genUsage generates usage info
func genUsage() *usage.Info {
	info := usage.NewInfo()

	info.AddOption(OPT_DIR, "Path to Redis-Split main dir", "dir")
	info.AddOption(OPT_KEYS, "Number of keys {s-}(10-1000000 default: 5000){!}")
	info.AddOption(OPT_RATIO, "Writes/reads ratio {s-}(1-100 default: 4){!}")
	info.AddOption(OPT_NO_COLOR, "Disable colors in output")
	info.AddOption(OPT_HELP, "Show this help message")
	info.AddOption(OPT_VER, "Show version")

	info.AddExample(
		"-d /srv/redis-split -k 35000 -r 10",
		"Run tool with custom settings",
	)

	return info
}

// genAbout generates info about version
func genAbout(gitRev string) *usage.About {
	about := &usage.About{
		App:     APP,
		Version: VER,
		Desc:    DESC,
		Year:    2006,
		Owner:   "ESSENTIAL KAOS",
		License: "Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>",
	}

	if gitRev != "" {
		about.Build = "git:" + gitRev
	}

	return about
}
