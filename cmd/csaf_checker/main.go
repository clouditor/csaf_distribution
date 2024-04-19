// This file is Free Software under the MIT License
// without warranty, see README.md and LICENSES/MIT.txt for details.
//
// SPDX-License-Identifier: MIT
//
// SPDX-FileCopyrightText: 2022 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
// Software-Engineering: 2022 Intevation GmbH <https://intevation.de>

// Package main implements the csaf_checker tool.
package main

import (
	"log"

	"github.com/csaf-poc/csaf_distribution/v3/csaf/client/checker"
	"github.com/csaf-poc/csaf_distribution/v3/internal/options"
)

// run uses a processor to check all the given domains or direct urls
// and generates a report.
func run(cfg *checker.Config, domains []string) (*checker.Report, error) {
	p, err := checker.NewProcessor(cfg)
	if err != nil {
		return nil, err
	}
	defer p.Close()
	return p.Run(domains)
}

func main() {
	domains, cfg, err := checker.ParseArgsConfig()
	options.ErrorCheck(err)
	options.ErrorCheck(cfg.Prepare())

	if len(domains) == 0 {
		log.Println("No domain or direct url given.")
		return
	}

	report, err := run(cfg, domains)
	options.ErrorCheck(err)

	options.ErrorCheck(report.Write(cfg.Format, cfg.Output))
}
