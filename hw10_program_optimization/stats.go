package hw10programoptimization

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/valyala/fastjson"
)

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	domainSuffix := fmt.Sprintf(".%s", domain)

	result := make(DomainStat)

	s := bufio.NewScanner(r)
	s.Split(bufio.ScanLines)

	for s.Scan() {
		if len(s.Bytes()) == 0 {
			continue
		}

		if err := fastjson.ValidateBytes(s.Bytes()); err != nil {
			return nil, err
		}

		email := fastjson.GetString(s.Bytes(), "Email")

		if strings.HasSuffix(email, domainSuffix) {
			if emailParts := strings.SplitN(email, "@", 2); len(emailParts) > 1 {
				result[strings.ToLower(emailParts[1])]++
			}
		}
	}

	return result, nil
}
