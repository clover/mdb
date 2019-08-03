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

package main

import (
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/gobwas/glob"
	"github.com/spf13/cobra"
)

// findCmd represents the find command
var tagsFindCmd = &cobra.Command{
	Use: "search <conditional>[,<jointype><conditional>, ...]",
	Long: `Search for a given tag, tag value, or tag combination

    <conditional> can be one of:
        <tagname> - match tag that exists with any value
        <tagname>=<value> - match tag that exists with an exact value

    Multiple conditionals can be combined using a comma and a join type.
    Joins can be chained to an arbitrary depth, and are evaluated left-to-
    right.

    Join types are:
        & or A: result list contains endpoints that exist on both sides of join
        - or S: result of right side of join is removed from left side of join
        + or M: merge results of left and right sides of join`,
	Example: `Show all endpoints that are puppet servers:

    mdb tags search 'role.infra.puppet'

Show all endpoints that are puppet servers located on p101:

    mdb tags search 'role.infra.puppet,&hn=p101.example.com'

Show all endpoints located on p101 that aren't puppet servers:

    mdb tags search 'hn=p101.example.com,-role.infra.puppet'`,
	Aliases: []string{"find"},
	Short:   "find endpoints with specific tag/tag value/tag set",
	Args:    cobra.ExactArgs(1),
	Run:     tagsFindCmdRun,
}

var reChainConditional = regexp.MustCompile(`(?i)([MAS+&-])(.*)`)

func tagsFindCmdRun(cmd *cobra.Command, args []string) {
	accum := map[string]bool{} // accumulated results

	conditionals := strings.Split(args[0], ",")
	endpoints := getendpointlist(kk)

	for i, conditional := range conditionals {
		// Set defaults for the first conditional (which conveniently defines
		// the variables and sets scope)
		jointype := "+"
		c := conditional

		// If not the first conditional, we have to know how to join them
		if i > 0 {
			m := reChainConditional.FindStringSubmatch(conditional)

			if m == nil {
				logerrl("Invalid join type or conditional string: %s", conditional)
				os.Exit(MDB_EXIT_USERFAIL)
			}

			jointype = m[1]
			c = m[2]
		}
		logtracel("processing conditional '%s' with join type '%s'", c, jointype)

		//		res := getmatches(endpoints, c)
		//		logtracel("found %d conditional matches", len(res))

		switch jointype {
		case "M": // "merge"
			fallthrough
		case "+":
			join_merge(accum, getmatches(endpoints, c))
		case "S": // "subtract"
			fallthrough
		case "-":
			join_sub(accum, getmatches(rkeys(accum), c))
		case "A": // "and"
			fallthrough
		case "&":
			join_and(accum, getmatches(rkeys(accum), c))
		default:
			logerrl("unknown join type '%s' (this shouldn't happen)", jointype)
			os.Exit(MDB_EXIT_USERFAIL)
		}
	}

	// We made it! Create our endpoint list
	el := []string{}
	for endpoint, _ := range accum {
		// Kind of a hack ... okay, really a hack. Stop hostnames starting
		// with _ from appearing in search results (they'll be templates)
		if !strings.HasPrefix(endpoint, "@") {
			el = append(el, endpoint)
		}
	}

	rettag, _ := cmd.Flags().GetString("rettag")

	// If we don't want a specific tag, output the match list
	if rettag == "" {
		logtracel("No rettag requested, outputting sorted list of %d matches",
			len(el))
		sort.Strings(el)
		for _, e := range el {
			stdoutl("%s", e)
		}
		return
	}

	// We want a specific tag, so lets iterate the list and fetch them
	// FIXME: Move to new function
	tel := []string{} // "tagged endpoint list"
	for _, e := range el {
		tagval, err := gettag(kk, e, rettag)
		if isKeyNotFound(err) {
			logtracel("endpoint '%s' has no tag '%s'", e, rettag)
			continue
		}

		if err != nil {
			logerr("etcd communication failure: %s", err)
			os.Exit(MDB_EXIT_ETCD_COMM)
		}

		tel = append(tel, tagval)
	}

	logtracel("rettag '%s' requested, returning %d entries (of %d)",
		rettag, len(tel), len(el))

	sort.Strings(tel)
	for _, e := range tel {
		stdoutl("%s", e)
	}
	return
}

// just join two maps (without duplicates)
func join_merge(dest map[string]bool, src map[string]bool) {
	for e, _ := range src {
		dest[e] = true
	}
}

// join two maps and keep only the common elements
func join_and(dest map[string]bool, src map[string]bool) {
	// Look through existing endpoints
	for e, _ := range dest {
		// Does the endpoint also exist in second map?
		if _, ok := src[e]; !ok {
			// Nope, remove it from the existing endpoints
			// (this is safe to do in a range in golang)
			delete(dest, e)
		}
	}
}

// join two maps, removing from the first things in the second
func join_sub(dest map[string]bool, src map[string]bool) {
	// Look through things we might remove
	for e, _ := range src {
		// Does the endpoint exist in the first map?
		if _, ok := dest[e]; ok {
			// Yep, so remove it
			delete(dest, e)
		}
	}
}

// FIXME: Should we hand an error upwards rather than exiting here?
func getmatches(endpoints []string, conditional string) map[string]bool {
	s := strings.SplitN(conditional, "=", 2)
	tag, _ := standardize_tagname(s[0])

	var g glob.Glob
	var err error
	var searchval string

	if len(s) > 1 {
		searchval = s[1]
		g, err = glob.Compile(searchval)
		if err != nil {
			logerrl("Invalid glob: %s", err)
			os.Exit(MDB_EXIT_USERFAIL)
		}
	} else {
		g = glob.MustCompile("*")
	}

	res := map[string]bool{}
	for _, endpoint := range endpoints {
		t, err := gettag(kk, endpoint, tag)
		if err != nil && isKeyNotFound(err) {
			logtracel("Tag %s not present on %s during search", tag, endpoint)
			continue
		}

		// Any other kind of error, abort to avoid improper results
		// FIXME: Should this exit logic be higher up?
		if err != nil {
			logerr("Error while searching endpoint tags on %s: %s", endpoint, err)
			os.Exit(MDB_EXIT_ETCD_COMM)
		}

		if len(searchval) > 0 && !g.Match(t) {
			logtracel("Tag %s present on %s during search, but doesn't match", tag, endpoint)
			continue
		}

		// We have a null searchval or we have a matching searchval
		res[endpoint] = true
	}

	return res
}

// func tagsFindCmdRun(cmd *cobra.Command, args []string) {
// 	s := strings.SplitN(args[0], "=", 2)
// 	endpoints := getendpointlist(kk)

// 	tag, _ := standardize_tagname(s[0])
// 	var searchval string
// 	if len(s) > 1 {
// 		searchval = s[1]
// 	}

// 	for _, endpoint := range endpoints {
// 		t, err := gettag(kk, endpoint, tag)
// 		if isKeyNotFound(err) {
// 			//			fmt.Printf("Tag not present on %s\n", endpoint)
// 			continue
// 		}

// 		if len(searchval) > 0 && t != searchval {
// 			//			fmt.Printf("Tag present but wrong value (%s != %s) on %s\n",
// 			//				t, searchval, endpoint)
// 			continue
// 		}
// 		fmt.Printf("%s\n", endpoint)
// 	}

// 	return
// }

func rkeys(resmap map[string]bool) []string {
	keys := []string{}

	for k := range resmap {
		keys = append(keys, k)
	}

	return keys
}

func init() {
	tagsCmd.AddCommand(tagsFindCmd)
	tagsFindCmd.Flags().StringP("rettag", "r", "",
		"return the named tag for each field")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// findCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// findCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
