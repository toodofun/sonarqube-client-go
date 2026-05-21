// Copyright 2025 The Toodofun Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http:www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

/*
This is a tool to generate client library for SonarQube (https://www.sonarqube.org/) web-api
It allows to generate client library based on your server version.
*/
package main

import (
	"flag"
	"log"
	"os"
)

var (
	host          = flag.String("host", "http://localhost:9000", "SonarQube server")
	deprecated    = flag.Bool("deprecated", false, "generate code for deprecated api methods")
	internal      = flag.Bool("internal", false, "generate code for internal methods")
	targetVersion = flag.String("target", "", "set target api version (default: server's version)")
	help          = flag.Bool("help", false, "show usage")
	out           = flag.String("out", ".", "output directory")
	auth          = flag.String("auth", "", "Authorization header value, e.g. Basic YWRtaW46YWRtaW4=")
	token         = flag.String("token", "", "SonarQube user token (squ_...); builds Basic auth when -auth is unset")
	packageName   = flag.String("package", "", "package name, if not set will be sonarqube_client")
	templateDir   = flag.String("template", "", "template directory under tool root (default: embedded tpl/)")
)

func main() {
	flag.Parse()
	if *help {
		flag.Usage()
		os.Exit(0)
	}

	authHeader := resolveAuthorization(*auth, *token)

	def, err := loadAPI(nil, *host, *deprecated, *internal, *targetVersion, authHeader)
	if err != nil {
		log.Fatal(err)
	}

	if err = generateCode(def, *out); err != nil {
		log.Fatal(err)
	}
}
