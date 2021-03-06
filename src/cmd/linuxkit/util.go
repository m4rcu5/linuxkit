package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func getStringValue(envKey string, flagVal string, defaultVal string) string {
	var res string

	// If defined, take the env variable
	if _, ok := os.LookupEnv(envKey); ok {
		res = os.Getenv(envKey)
	}

	// If a flag is specified, this value takes precedence
	// Ignore cases where the flag carries the default value
	if flagVal != "" && flagVal != defaultVal {
		res = flagVal
	}

	// if we still don't have a value, use the default
	if res == "" {
		res = defaultVal
	}
	return res
}

func getIntValue(envKey string, flagVal int, defaultVal int) int {
	var res int

	// If defined, take the env variable
	if _, ok := os.LookupEnv(envKey); ok {
		var err error
		res, err = strconv.Atoi(os.Getenv(envKey))
		if err != nil {
			res = 0
		}
	}

	// If a flag is specified, this value takes precedence
	// Ignore cases where the flag carries the default value
	if flagVal > 0 {
		res = flagVal
	}

	// if we still don't have a value, use the default
	if res == 0 {
		res = defaultVal
	}
	return res
}

func getBoolValue(envKey string, flagVal bool) bool {
	var res bool

	// If defined, take the env variable
	if _, ok := os.LookupEnv(envKey); ok {
		switch os.Getenv(envKey) {
		case "":
			res = false
		case "0":
			res = false
		case "false":
			res = false
		case "FALSE":
			res = false
		case "1":
			res = true
		default:
			// catches "true", "TRUE" or anything else
			res = true

		}
	}

	// If a flag is specified, this value takes precedence
	if res != flagVal {
		res = flagVal
	}

	return res
}

func stringToIntArray(l string, sep string) ([]int, error) {
	var err error
	if l == "" {
		return []int{}, err
	}
	s := strings.Split(l, sep)
	i := make([]int, len(s))
	for idx := range s {
		if i[idx], err = strconv.Atoi(s[idx]); err != nil {
			return nil, err
		}
	}
	return i, nil
}

// This function parses the "size" parameter of a disk specification
// and returns the size in MB. The "size" paramter defaults to GB, but
// the unit can be explicitly set with either a G (for GB) or M (for
// MB). It returns the disk size in MB.
func getDiskSizeMB(s string) (int, error) {
	if s == "" {
		return 0, nil
	}
	sz := len(s)
	if strings.HasSuffix(s, "M") {
		i, err := strconv.Atoi(s[:sz-1])
		if err != nil {
			return 0, err
		}
		return i, nil
	}
	if strings.HasSuffix(s, "G") {
		s = s[:sz-1]
	}
	return strconv.Atoi(s)
}

// DiskConfig is the config for a disk
type DiskConfig struct {
	Path   string
	Size   int
	Format string
}

// Disks is the type for a list of DiskConfig
type Disks []DiskConfig

func (l *Disks) String() string {
	return fmt.Sprint(*l)
}

// Set is used by flag to configure value from CLI
func (l *Disks) Set(value string) error {
	d := DiskConfig{}
	s := strings.Split(value, ",")
	for _, p := range s {
		c := strings.SplitN(p, "=", 2)
		switch len(c) {
		case 1:
			// assume it is a filename even if no file=x
			d.Path = c[0]
		case 2:
			switch c[0] {
			case "file":
				d.Path = c[1]
			case "size":
				size, err := getDiskSizeMB(c[1])
				if err != nil {
					return err
				}
				d.Size = size
			case "format":
				d.Format = c[1]
			default:
				return fmt.Errorf("Unknown disk config: %s", c[0])
			}
		}
	}
	*l = append(*l, d)
	return nil
}
