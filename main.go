package main

import (
	"fmt"
	"log"

	"github.com/projectdiscovery/subfinder/v2/pkg/resolve"
	"github.com/projectdiscovery/subfinder/v2/pkg/runner"
	"hamravesh.ir/mehrdad-khojastefar/subfinder"
)

func main() {
	r, err := subfinder.NewRunner("test", &runner.Options{
		Threads:            10,
		Timeout:            30,
		MaxEnumerationTime: 10,
		Resolvers:          resolve.DefaultResolvers,
	})
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(r.GetSubdomainArray("iran.ir"))
}
