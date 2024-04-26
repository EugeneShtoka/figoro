/*
Copyright © 2024 Eugene Shtoka <eshtoka@gmail.com>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package typedkeyring

import (
	"encoding/json"
	"fmt"

	"github.com/zalando/go-keyring"
)

type Keyring[T any] struct {
	ServiceName		string
}

func New[T any](serviceName string) *Keyring[T] {
    return &Keyring[T]{ ServiceName: serviceName }
}

func (k *Keyring[T]) Delete(name string) error {
    err := keyring.Delete(k.ServiceName, name)
    if err != nil {
        return fmt.Errorf("failed to delete token '%s': %w", name, err)
    }

    return nil
}

func (k *Keyring[T]) Load(name string) (*T, error) {
	data, err := keyring.Get(k.ServiceName, name)
	if err != nil {
        return nil, fmt.Errorf("failed to load token: %w", err)
    }

	var value T
    err = json.Unmarshal([]byte(data), &value)
    if err != nil {
        return nil, fmt.Errorf("error deserializing object: %w" + err.Error())
    }

    return &value, nil
}

func (k *Keyring[T])Save(name string, value *T) error {
    jsonData, err := json.Marshal(value)
    if err != nil {
        return fmt.Errorf("failed to serialize token to JSON: %w", err)
    }

    jsonStr := string(jsonData) 
	err = keyring.Set(k.ServiceName, name, jsonStr)
	if err != nil {
        return fmt.Errorf("failed to save token: %w", err)
    }

	return nil
}