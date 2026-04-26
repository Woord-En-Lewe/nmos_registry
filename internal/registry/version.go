package registry

import (
	"fmt"
	"strconv"
	"strings"
)

type Version struct {
	Seconds int64
	Nanos   int64
}

func ParseVersion(v string) (*Version, error) {
	parts := strings.Split(v, ":")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid version format: %s (expected seconds:nanoseconds)", v)
	}

	secs, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid seconds value: %s", parts[0])
	}

	nanos, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid nanoseconds value: %s", parts[1])
	}

	return &Version{Seconds: secs, Nanos: nanos}, nil
}

func CompareVersions(v1, v2 string) (int, error) {
	parsed1, err := ParseVersion(v1)
	if err != nil {
		return 0, fmt.Errorf("cannot compare versions: %w", err)
	}

	parsed2, err := ParseVersion(v2)
	if err != nil {
		return 0, fmt.Errorf("cannot compare versions: %w", err)
	}

	if parsed1.Seconds > parsed2.Seconds {
		return 1, nil
	}
	if parsed1.Seconds < parsed2.Seconds {
		return -1, nil
	}

	if parsed1.Nanos > parsed2.Nanos {
		return 1, nil
	}
	if parsed1.Nanos < parsed2.Nanos {
		return -1, nil
	}

	return 0, nil
}

func IsVersionEarlier(v1, v2 string) (bool, error) {
	cmp, err := CompareVersions(v1, v2)
	if err != nil {
		return false, err
	}
	return cmp < 0, nil
}

func IsVersionLater(v1, v2 string) (bool, error) {
	cmp, err := CompareVersions(v1, v2)
	if err != nil {
		return false, err
	}
	return cmp > 0, nil
}

func IsVersionEqual(v1, v2 string) (bool, error) {
	cmp, err := CompareVersions(v1, v2)
	if err != nil {
		return false, err
	}
	return cmp == 0, nil
}