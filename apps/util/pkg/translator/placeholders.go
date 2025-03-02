package translator

import "fmt"

type placeholders map[string]string

func (p *placeholders) add(counter int, value string) string {
	for _, template := range []string{"ph", "pH", "Ph", "PH"} {
		for _, lSpace := range []string{"", " "} {
			for _, rSpace := range []string{"", " "} {
				placeholder := fmt.Sprintf("(%s%s_%d%s)", lSpace, template, counter, rSpace)
				(*p)[placeholder] = value
			}
		}
	}

	return fmt.Sprintf("(PH_%d)", counter)
}
