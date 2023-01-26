package main

import (
	"context"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/logrusorgru/aurora"

	"github.com/projectdiscovery/nuclei/v2/pkg/catalog/config"
	"github.com/projectdiscovery/nuclei/v2/pkg/catalog/disk"
	"github.com/projectdiscovery/nuclei/v2/pkg/catalog/loader"
	"github.com/projectdiscovery/nuclei/v2/pkg/core"
	"github.com/projectdiscovery/nuclei/v2/pkg/core/inputs"
	"github.com/projectdiscovery/nuclei/v2/pkg/model/types/severity"
	"github.com/projectdiscovery/nuclei/v2/pkg/output"
	"github.com/projectdiscovery/nuclei/v2/pkg/parsers"
	"github.com/projectdiscovery/nuclei/v2/pkg/protocols"
	"github.com/projectdiscovery/nuclei/v2/pkg/protocols/common/contextargs"
	"github.com/projectdiscovery/nuclei/v2/pkg/protocols/common/hosterrorscache"
	"github.com/projectdiscovery/nuclei/v2/pkg/protocols/common/interactsh"
	"github.com/projectdiscovery/nuclei/v2/pkg/protocols/common/protocolinit"
	"github.com/projectdiscovery/nuclei/v2/pkg/protocols/common/protocolstate"
	"github.com/projectdiscovery/nuclei/v2/pkg/reporting"
	"github.com/projectdiscovery/nuclei/v2/pkg/testutils"
	"github.com/projectdiscovery/nuclei/v2/pkg/types"
	"github.com/projectdiscovery/ratelimit"
)

func StartNuclei(domains []string) ([]*output.ResultEvent, error) {
	// // testing
	// domains = []string{"localhost"}
	fmt.Println("Starting Nuclei...")
	results := []*output.ResultEvent{}
	cache := hosterrorscache.New(30, hosterrorscache.DefaultMaxHostsCount)
	defer cache.Close()

	mockProgress := &testutils.MockProgressClient{}
	reportingClient, _ := reporting.New(&reporting.Options{}, "")
	defer reportingClient.Close()

	outputWriter := testutils.NewMockOutputWriter()
	outputWriter.WriteCallback = func(event *output.ResultEvent) {
		results = append(results, event)
		fmt.Println(results)
	}

	defaultOpts := types.DefaultOptions()
	protocolstate.Init(defaultOpts)
	protocolinit.Init(defaultOpts)

	// defaultOpts.IncludeIds = goflags.StringSlice{"cname-service", "cve"}
	// defaultOpts.IncludeTags = goflags.StringSlice{"cve", "exposure", "tech", "xss", "lfi"}
	defaultOpts.Severities = severity.Severities{severity.High, severity.Critical}
	// defaultOpts.ExcludeTags = config.ReadIgnoreFile().Tags

	interactOpts := interactsh.NewDefaultOptions(outputWriter, reportingClient, mockProgress)
	interactClient, err := interactsh.New(interactOpts)
	if err != nil {
		return results, err
	}
	defer interactClient.Close()

	home, _ := os.UserHomeDir()
	catalog := disk.NewCatalog(path.Join(home, "nuclei-templates"))
	executerOpts := protocols.ExecuterOptions{
		Output:          outputWriter,
		Options:         defaultOpts,
		Progress:        mockProgress,
		Catalog:         catalog,
		IssuesClient:    reportingClient,
		RateLimiter:     ratelimit.New(context.Background(), 150, time.Second),
		Interactsh:      interactClient,
		HostErrorsCache: cache,
		Colorizer:       aurora.NewAurora(true),
		ResumeCfg:       types.NewResumeCfg(),
	}
	engine := core.New(defaultOpts)
	engine.SetExecuterOptions(executerOpts)

	workflowLoader, err := parsers.NewLoader(&executerOpts)
	if err != nil {
		return results, err
	}
	executerOpts.WorkflowLoader = workflowLoader
	defaultOpts.TemplateList = true
	defaultOpts.AutomaticScan = true
	configObject, err := config.ReadConfiguration()
	if err != nil {
		return results, err
	}

	store, err := loader.New(loader.NewConfig(defaultOpts, configObject, catalog, executerOpts))
	if err != nil {
		return results, err
	}
	store.Load()

	inputArgs := []*contextargs.MetaInput{}
	for _, d := range domains {
		inputArgs = append(inputArgs, &contextargs.MetaInput{Input: d})

	}

	input := &inputs.SimpleInputProvider{Inputs: inputArgs}
	_ = engine.Execute(store.Templates(), input)
	fmt.Println("Started Nuclei...")
	engine.WorkPool().Wait() // Wait for the scan to finish
	fmt.Println("finished Nuclei...")
	fmt.Println(results)
	return results, nil
}
