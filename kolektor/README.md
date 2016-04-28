Build: `go install github.com/ezeev/golang/kolektor`

Run: `../bin/kolektor -config github.com/ezeev/golang/kolektor/config.json -collectors github.com/ezeev/golang/kolektor/collector-yaml`

View Profile: `go tool pprof -text ../bin/kolektor /var/folders/2y/_w_pspnj3f5ddlk7f01wsh9h0000gn/T/profile279781245/cpu.pprof`
