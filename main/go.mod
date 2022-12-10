module example.com/main

go 1.19

replace example.com/timeAndDate => ../timeAndDate

replace example.com/weatherApi => ../weatherApi

replace example.com/openWeatherApi => ../openWeatherApi

require example.com/openWeatherApi v0.0.0-00010101000000-000000000000

require example.com/weatherApi v0.0.0-00010101000000-000000000000

require example.com/timeAndDate v0.0.0-00010101000000-000000000000

require (
	github.com/PuerkitoBio/goquery v1.8.0 // indirect
	github.com/VividCortex/gohistogram v1.0.0 // indirect
	github.com/andybalholm/cascadia v1.3.1 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/chromedp/cdproto v0.0.0-20221126224343-3a0787b8dd28 // indirect
	github.com/chromedp/chromedp v0.8.6 // indirect
	github.com/chromedp/sysutil v1.0.0 // indirect
	github.com/dariubs/percent v1.0.0 // indirect
	github.com/geziyor/geziyor v0.0.0-20220429000531-738852f9321d // indirect
	github.com/go-kit/kit v0.12.0 // indirect
	github.com/gobwas/httphead v0.1.0 // indirect
	github.com/gobwas/pool v0.2.1 // indirect
	github.com/gobwas/ws v1.1.0 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.4 // indirect
	github.com/prometheus/client_golang v1.14.0 // indirect
	github.com/prometheus/client_model v0.3.0 // indirect
	github.com/prometheus/common v0.37.0 // indirect
	github.com/prometheus/procfs v0.8.0 // indirect
	github.com/temoto/robotstxt v1.1.2 // indirect
	golang.org/x/net v0.4.0 // indirect
	golang.org/x/sys v0.3.0 // indirect
	golang.org/x/text v0.5.0 // indirect
	golang.org/x/time v0.3.0 // indirect
	google.golang.org/protobuf v1.28.1 // indirect
)
