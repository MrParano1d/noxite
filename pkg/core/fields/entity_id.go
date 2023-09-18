package fields

import (
	"fmt"
	"strconv"
	"strings"
)

type EntityID int

func (id EntityID) Int() int {
	return int(id)
}

func (id EntityID) String() string {
	return strconv.Itoa(id.Int())
}

func EntityIDFromInt(id int) (EntityID, error) {
	if id < 1 {
		return EntityID(0), InvalidEntityIDErr{ID: id, Reason: "must be greater than 0"}
	}
	return EntityID(id), nil
}

func EntityIDFromString(id string) (EntityID, error) {
	id = strings.TrimSpace(id)
	if len(id) == 0 {
		return EntityID(0), InvalidEntityIDErr{ID: 0, Reason: "empty string"}
	}

	idNum, err := strconv.Atoi(id)
	if err != nil {
		return EntityID(0), InvalidEntityIDErr{ID: 0, Reason: "not a number"}
	}

	return EntityIDFromInt(idNum)
}

type InvalidEntityIDErr struct {
	ID     int
	Reason string
}

func (e InvalidEntityIDErr) Error() string {
	return fmt.Sprintf("invalid entity id %d: %s", e.ID, e.Reason)
}
