package cow

import (
	"errors"
	"reflect"
	"sync"
)

type CopyOnWriteSlice struct {
	values   []any
	typeKind reflect.Kind

	sync.Mutex
}

func NewCopyOnWriteSlice(initSize int) *CopyOnWriteSlice {
	return &CopyOnWriteSlice{
		values:   make([]any, 0, initSize),
		typeKind: 0,
	}
}

func (c *CopyOnWriteSlice) Append(values ...any) {
	if len(values) == 0 {
		return
	}
	c.Lock()
	defer func() {
		c.Unlock()
	}()

	newSlice := make([]any, len(c.values), len(c.values)+len(values))
	copy(newSlice, c.values)
	newSlice = append(newSlice, values)
	c.setValues(newSlice)
}

func (c *CopyOnWriteSlice) Set(index int, value any) error {
	err := c.checkIndex(index)
	if err != nil {
		return err
	}
	err = c.checkType(value)
	if err != nil {
		return err
	}

	c.Lock()
	defer func() {
		c.Unlock()
	}()
	newSlice := make([]any, 0, len(c.values))
	newSlice[index] = value
	c.setValues(newSlice)
	return nil
}

func (c *CopyOnWriteSlice) Get(index int) (any, error) {
	err := c.checkIndex(index)
	if err != nil {
		return nil, err
	}

	return c.values[index], nil
}

func (c *CopyOnWriteSlice) All() ([]any, error) {
	res := make([]any, len(c.values))
	copy(res, c.values)
	return res, nil
}

func (c *CopyOnWriteSlice) ForEach(f func(index int, value any)) {
	for index, value := range c.values {
		f(index, value)
	}
}

func (c *CopyOnWriteSlice) checkType(value any) error {
	typeKind := reflect.TypeOf(value).Kind()
	if typeKind == 0 {
		c.typeKind = typeKind
	} else {
		if c.typeKind != typeKind {
			return errors.New("input type is inconsistent with the existing type")
		}
	}
	return nil
}

func (c *CopyOnWriteSlice) checkIndex(index int) error {
	if index >= len(c.values) {
		return errors.New("index out of len")
	}
	return nil
}

func (c *CopyOnWriteSlice) setValues(values []any) {
	c.values = values
}
