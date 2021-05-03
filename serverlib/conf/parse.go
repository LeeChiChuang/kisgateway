package conf

import (
	"bufio"
	"io"
	"os"
	"strconv"
	"strings"
)

type Conf map[string]map[string]string

const configLen = 2

var ConfigInfo Conf

func ReadConf(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	rd := bufio.NewReader(f)
	ConfigInfo = make(Conf)
	for {
		line, err := rd.ReadString('\n')
		if err != nil || err == io.EOF {
			break
		}
		line = strings.TrimSpace(line)
		if len(line) == 0 || line[0] == ';' || line[0] == '[' {
			continue
		}
		t := strings.Split(line, "=")
		if len(t) != configLen || strings.TrimSpace(t[1]) == "" {
			continue
		}

		keys := strings.Split(strings.TrimSpace(t[0]), ".")
		prefix := keys[0]
		key := keys[1]
		if _, ok := ConfigInfo[prefix]; !ok {
			ConfigInfo[prefix] = make(map[string]string)
		}
		ConfigInfo[prefix][key] = strings.TrimSpace(t[1])
	}

	return nil
}

func GetConf(key string) string {
	return getConf(key)
}

func GetConfInt(key string) int {
	i, _ := strconv.Atoi(getConf(key))
	return i
}

func getConf(key string) string {
	keys := strings.Split(strings.TrimSpace(key), ".")
	prefix := keys[0]
	k := keys[1]
	if v, ok := ConfigInfo[prefix]; !ok {
		return ""
	} else {
		return v[k]
	}
}
