package profari

import "fmt"

type Test interface {
	Run()
	Teardown()
	GetName() string
	GetErrChan() chan error
}

func RunTests(tests ...Test) string {
	rc := make(chan result)

	for _, s := range tests {
		go func(s Test) {
			r := result{
				name: s.GetName(),
			}
			go s.Run()
			r.success = true
			for err := range s.GetErrChan() {
				if err != nil {
					r.success = false
					s.Teardown()
					rc <- r
					return
				}
			}
			s.Teardown()
			rc <- r
		}(s)
	}

	template := "%-20v%-20v%-20v\n"
	resultText := fmt.Sprintf(template, "ID", "Suite", "Success?")
	for j := 0; j < len(tests); j++ {
		r := <-rc
		r.id = j + 1
		row := fmt.Sprintf(template, r.id, r.name, r.success)
		resultText += row
	}
	return resultText
}

type result struct {
	name    string
	success bool
	id      int
}
