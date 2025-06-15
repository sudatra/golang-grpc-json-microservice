package main

import (
	"flag"
)

func main() {
	listenAddr := flag.String("listenaddr", ":3000", "listen address the service is running on");
	flag.Parse();

	svc := NewLoggingService(
		NewMetricService(
			&priceFetcher{},
		),
	);
	
	server := newJSONAPIServer(*listenAddr, svc);
	server.Run();
}