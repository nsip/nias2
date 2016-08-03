package nias2

import (
	"bufio"
	"github.com/beevik/etree"
	//"log"
	"os"
	"regexp"
	"strings"
)

type PrivacyFilter struct {
	Path etree.Path // path to filter
	Attr string     // attribute to filter if any (etree.Path does not support attribute queries)
	Repl string     // replacement string
}

// implementation of the psi service
type PrivacyService struct {
	filters [][]PrivacyFilter
}

var sensitivities = [5]string{"none", "low", "medium", "high", "extreme"}

// create a new psi service instance
func NewPrivacyService() (*PrivacyService, error) {
	pri := PrivacyService{}
	filterfiles := make([]string, 4)
	filterfiles[0] = "./privacyfilters/low.xpath"
	filterfiles[1] = "./privacyfilters/medium.xpath"
	filterfiles[2] = "./privacyfilters/high.xpath"
	filterfiles[3] = "./privacyfilters/extreme.xpath"
	filter_txt := make([][]string, 4)
	for i, filename := range filterfiles {
		var file *os.File
		var err error
		if file, err = os.Open(filename); err != nil {
			return nil, err
		}
		defer file.Close()
		var lines []string
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}
		err = scanner.Err()
		if err != nil {
			return nil, err
		}
		filter_txt[i] = make([]string, len(lines))
		copy(filter_txt[i], lines)
	}
	pri.filters = parse_filters(filter_txt)
	return &pri, nil
}

// Filter out elements or attributes. Return array of filtered XML strings,
// with cumulative filtering
func filter(xml string, filters [][]PrivacyFilter) ([]string, error) {
	doc := etree.NewDocument()
	if err := doc.ReadFromString(xml); err != nil {
		return nil, err
	}
	ret := make([]string, len(filters)+1)
	ret[0] = xml
	for i, filterslice := range filters {
		for _, filter := range filterslice {
			for _, elem := range doc.FindElementsPath(filter.Path) {
				if filter.Attr == "" {
					filter1(elem, filter.Repl)
				} else {
					if attr := elem.SelectAttr(filter.Attr); attr != nil {
						attr.Value = filter.Repl
					}
				}
			}
		}
		out, err := doc.WriteToString() // output cumulative effect of filtering
		if err != nil {
			return nil, err
		} else {
			ret[i+1] = out
		}
	}
	return ret, nil
}

// recursively replace all chardata in the element and its children with repl
func filter1(elem *etree.Element, repl string) {
	children := elem.ChildElements()
	if len(children) > 0 {
		for _, child := range children {
			filter1(child, repl)
		}
	} else {
		elem.SetText(repl)
	}
}

// Parse a slice of slices of filters specified as XXXX@attr:Replacement into PrivacyFilter.
// Each slice of filters is a privacy filter setting; these are applied cumulatively.
// Any trailing attributes are split out because etree does not understand attribute XPaths.
// If no Replacement text is supplied, ZZREDACTED is supplied by default
func parse_filters(filter_txt [][]string) [][]PrivacyFilter {
	re := regexp.MustCompile("\\@[^]/@]+$")
	ret := make([][]PrivacyFilter, len(filter_txt))
	for i := range filter_txt {
		ret[i] = make([]PrivacyFilter, 0, len(filter_txt[i]))
		for j := range filter_txt[i] {
			if len(filter_txt[i][j]) == 0 {
				continue // blank
			}
			if filter_txt[i][j][0] == '#' {
				continue // comment
			}
			var new PrivacyFilter
			parts := strings.Split(filter_txt[i][j], ":")
			if len(parts) >= 2 {
				new.Repl = parts[1]
			} else {
				new.Repl = "ZZREDACTED"
			}
			attr := re.FindString(parts[0])
			if attr == "" {
				new.Path = etree.MustCompilePath(parts[0])
				new.Attr = ""
			} else {
				new.Path = etree.MustCompilePath(string(re.ReplaceAll([]byte(parts[0]), []byte(""))))
				new.Attr = attr[1:len(attr)]
			}
			ret[i] = append(ret[i], new)
		}
	}
	return ret
}

// implement the nias Service interface
func (pri *PrivacyService) HandleMessage(req *NiasMessage) ([]NiasMessage, error) {

	responses := make([]NiasMessage, 0)
	out, err := filter(req.Body.(string), pri.filters)
	if err != nil {
		return nil, err
	}
	for i, out1 := range out {
		r := NiasMessage{}
		r.TxID = req.TxID
		r.SeqNo = req.SeqNo
		r.Target = STORE_AND_FORWARD_PREFIX + sensitivities[i] + "::"
		r.Body = out1
		responses = append(responses, r)
	}
	return responses, nil
}
