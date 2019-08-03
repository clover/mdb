// Copyright 2019 Clover Network, Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Really basic logging stuff. We should eventually replace this with an
// actual logging package
//
// tbh, this is just a copypasta of a stupid logging package Jay tossed
// together for personal projects
package main

import (
	"fmt"
	"os"
	"regexp"
	"runtime"
)

const (
	LogLvlNone = iota
	LogLvlError
	LogLvlWarning
	LogLvlInfo
	LogLvlDebug
	LogLvlTrace
)

// Default loglevel
var Loglevel = LogLvlWarning
var Logcallers = false

func stdout(format string, vals ...interface{}) {
	outmsg := fmt.Sprintf(format, vals...)
	fmt.Fprint(os.Stdout, outmsg)
}

func stdoutl(format string, vals ...interface{}) {
	outmsg := fmt.Sprintf(format+"\n", vals...)
	fmt.Fprint(os.Stdout, outmsg)
}

func stderr(format string, vals ...interface{}) {
	outmsg := fmt.Sprintf(format, vals...)
	fmt.Fprint(os.Stderr, outmsg)
}

func stderrl(format string, vals ...interface{}) {
	outmsg := fmt.Sprintf(format+"\n", vals...)
	fmt.Fprint(os.Stderr, outmsg)
}

var reEndpath = regexp.MustCompile(`([^/]*)$`)

func addcaller(format string) string {
	if Logcallers {
		// pc, file, line, ok := runtime.Caller(2)
		pc, _, line, ok := runtime.Caller(2)
		f := runtime.FuncForPC(pc)
		fn := "NOFUNC"
		if ok {
			m := reEndpath.FindStringSubmatch(f.Name())
			if m != nil {
				fn = m[1]
			} else {
				fn = f.Name()
			}
		}

		// r := reEndpath.FindStringSubmatch(file)
		// if r != nil {
		// 	file = r[1]
		// }

		format = fmt.Sprintf("(%s#%d) %s", fn, line, format)
	}

	return format
}

func log(format string, vals ...interface{}) {
	stderr(format, vals...)
}

func logl(format string, vals ...interface{}) {
	log(format+"\n", vals...)
}

func logerr(format string, vals ...interface{}) {
	if Loglevel >= LogLvlError {
		log("ERROR: "+addcaller(format), vals...)
	}
}

func logerrl(format string, vals ...interface{}) {
	if Loglevel >= LogLvlError {
		logl("ERROR: "+addcaller(format), vals...)
	}
}

func logwarn(format string, vals ...interface{}) {
	if Loglevel >= LogLvlWarning {
		log("WARNING: "+addcaller(format), vals...)
	}
}

func logwarnl(format string, vals ...interface{}) {
	if Loglevel >= LogLvlWarning {
		logl("WARNING: "+addcaller(format), vals...)
	}
}

func loginfo(format string, vals ...interface{}) {
	if Loglevel >= LogLvlInfo {
		log("INFO: "+addcaller(format), vals...)
	}
}

func loginfol(format string, vals ...interface{}) {
	if Loglevel >= LogLvlInfo {
		logl("INFO: "+addcaller(format), vals...)
	}
}

func logdebug(format string, vals ...interface{}) {
	if Loglevel >= LogLvlDebug {
		log("DEBUG: "+addcaller(format), vals...)
	}
}

func logdebugl(format string, vals ...interface{}) {
	if Loglevel >= LogLvlDebug {
		logl("DEBUG: "+addcaller(format), vals...)
	}
}

func logtrace(format string, vals ...interface{}) {
	if Loglevel >= LogLvlTrace {
		logl("TRACE: "+addcaller(format), vals...)
	}
}

func logtracel(format string, vals ...interface{}) {
	if Loglevel >= LogLvlTrace {
		logl("TRACE: "+addcaller(format), vals...)
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
