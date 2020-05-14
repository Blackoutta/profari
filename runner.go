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
				name:    s.GetName(),
				success: true,
			}

			go func() {
				for {
					select {
					case err := <-s.GetErrChan():
						if err != nil {
							if err.Error() == "done" {
								fmt.Println("end test!")
								rc <- r
								return
							}
							r.success = false
						}
					}
				}
			}()

			s.Run()
			s.Teardown()
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
