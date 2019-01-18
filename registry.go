package mregistry

import (
	"errors"
	"sync"

	"golang.org/x/sys/windows/registry"
)

// QKEY - Quattro KEY
// lpKey - key of registry.Key
// lpLocation - key location in registry
// lpAccess - access to registry key and values
type qkey struct {
	lpKey      registry.Key
	lpLocation string
	lpAccess   uint32
}

const (
	DWORD = iota + 1<<1 - 1
	QWORD
	REG_SZ
	EXPAND_SZ
)

// QID - Quattro Formaggi
// QID have 2 different values:
// ... 1 - DWordValue
// ... 2 - QWordValue
// ... 3 - StringValue
// ... 4 - StringExpandValue
type qid struct {
	vtype int
}

var (
	ErrUnknownType = errors.New("unknown provided type of registry value")
	ErrLocationArg = errors.New("your not set any data in `location` argument")
	ErrNamesArg    = errors.New("your not set name data in `names` argument")
)

var mu sync.RWMutex

// makeQkey makes automatic key structure for different usage
func makeQkey(key registry.Key, location string, access uint32) qkey {
	mu.RLock()
	defer mu.RUnlock()

	return qkey{
		lpKey:      key,
		lpLocation: location,
		lpAccess:   access,
	}
}

// checkArgumentsErrors checking errors of location and names which transmission as function argument
func checkArgumentsErrors(location string, names []string) error {
	mu.RLock()
	defer mu.RUnlock()

	if location == "" && len(location) == 0 {
		return ErrLocationArg
	}
	if len(names) == 0 && cap(names) == 0 {
		return ErrNamesArg
	}
	return nil
}

// registryMakeUniversal return universal qkey and parallel checks
func registryMakeUniversal(key registry.Key, location string, access uint32, names []string) (qkey, error) {
	mu.RLock()
	defer mu.RUnlock()

	if err := checkArgumentsErrors(location, names); err != nil {
		return makeQkey(key, location, access), err
	}
	return makeQkey(key, location, access), nil
}

func (q qkey) setCustomDQValues(id qid, names []string, dwordv []uint32, qwordv []uint64) error {
	key, err := registry.OpenKey(q.lpKey, q.lpLocation, q.lpAccess)
	defer key.Close()
	if err != nil {
		return err
	}

	/*if id.vtype > (1<<2)-1 {
		return ErrUnknownType
	}*/
	if id.vtype == DWORD {
		wait.Add(len(dwordv))
		if len(qwordv) == 0 && cap(qwordv) == 0 {
			for i, _ := range names {
				if err = key.SetDWordValue(names[i], dwordv[i]); err != nil {
					return err
				}
			}
		}
	} else if id.vtype == QWORD {
		if len(dwordv) == 0 && cap(dwordv) == 0 {
			for i, _ := range names {
				if err = key.SetQWordValue(names[i], qwordv[i]); err != nil {
					return err
				}
			}
		}
	} else {
		return ErrUnknownType
	}
	return nil
}

// SetMultipleDWordValues sets multiple new DWORD values in one registry key.
// Based on registry.SetDWordValue
func SetMultipleDWordValues(key registry.Key, location string, access uint32, names []string, values []uint32) error {
	q, err := registryMakeUniversal(key, location, access, names)
	if err != nil {
		return err
	}
	qidSt := qid{vtype: DWORD}
	if err := q.setCustomDQValues(qidSt, names, values, []uint64{}); err != nil {
		return err
	}
	return nil
}

// SetMultipleQWordValues sets multiple new QWORD(different of DWORD) values in one registry key.
// Based on registry.SetQWordValue
func SetMultipleQWordValues(key registry.Key, location string, access uint32, names []string, values []uint64) error {
	q, err := registryMakeUniversal(key, location, access, names)
	if err != nil {
		return err
	}
	qidSt := qid{vtype: QWORD}
	if err := q.setCustomDQValues(qidSt, names, []uint32{}, values); err != nil {
		return err
	}
	return nil
}

func (q qkey) setCustomBinaryValues(names []string, value [][]byte) error {
	key, err := registry.OpenKey(q.lpKey, q.lpLocation, q.lpAccess)
	defer key.Close()
	if err != nil {
		return err
	}
	if value != nil {
		for i, _ := range names {
			if err = key.SetBinaryValue(names[i], value[i]); err != nil {
				return err
			}
		}
	}
	return nil
}

// SetMultipleBinaryValues sets multiple new BINARY values in one registry key.
// Based on registry.SetBinaryValue
func SetMultipleBinaryValues(key registry.Key, location string, access uint32, names []string, values [][]byte) error {
	q, err := registryMakeUniversal(key, location, access, names)
	if err != nil {
		return err
	}
	if len(values) == 0 {
		values = append(values, []byte("\x00"))
	}
	if err := q.setCustomBinaryValues(names, values); err != nil {
		return err
	}
	return nil
}

func (q qkey) setCustomStringValues(id qid, names []string, value []string) error {
	key, err := registry.OpenKey(q.lpKey, q.lpLocation, q.lpAccess)
	defer key.Close()
	if err != nil {
		return err
	}
	if value != nil {
		if id.vtype == REG_SZ {
			for i, _ := range names {
				if err = key.SetStringValue(names[i], value[i]); err != nil {
					return err
				}
			}
		} else if id.vtype == EXPAND_SZ {
			for i, _ := range names {
				if err = key.SetExpandStringValue(names[i], value[i]); err != nil {
					return err
				}
			}
		} else {
			return ErrUnknownType
		}
	}
	return nil
}

// SetMultipleStringValues sets multiple new REG_SZ values in one registry key.
// Based on registry.SetStringValue
func SetMultipleStringValues(key registry.Key, location string, access uint32, names []string, values []string) error {
	q, err := registryMakeUniversal(key, location, access, names)
	if err != nil {
		return err
	}
	qidSt := qid{vtype: REG_SZ}
	if err := q.setCustomStringValues(qidSt, names, values); err != nil {
		return err
	}
	return nil
}

// SetMultipleExpandStringValues sets multiple new EXPAND_SZ values in one registry key.
// Based on registry.SetExpandStringValue
func SetMultipleExpandStringValues(key registry.Key, location string, access uint32, names []string, values []string) error {
	q, err := registryMakeUniversal(key, location, access, names)
	if err != nil {
		return err
	}
	qidSt := qid{vtype: EXPAND_SZ}
	if err := q.setCustomStringValues(qidSt, names, values); err != nil {
		return err
	}
	return nil
}
