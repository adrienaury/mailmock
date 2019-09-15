// Copyright (C) 2019  Adrien Aury
//
// This file is part of Mailmock.
//
// Mailmock is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// Mailmock is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with Mailmock.  If not, see <https://www.gnu.org/licenses/>.

// Package repository handles storage access for Mailmock REST API.
package repository

var storedObjects = []interface{}{}

// Store stores th object and gives it an ID.
func Store(o interface{}) int {
	id := len(storedObjects)
	storedObjects = append(storedObjects, o)
	return id
}

// Use returns the object with ID or nil.
func Use(id int) interface{} {
	if id < len(storedObjects) {
		return storedObjects[id]
	}
	return nil
}

// All returns all objects currently stored.
func All(from, limit int) (map[int]interface{}, bool) {
	if from < len(storedObjects) {
		if from+limit < len(storedObjects) {
			return tomap(from, from+limit), false
		}
		return tomap(from, len(storedObjects)), from == 0
	}
	if from > len(storedObjects) {
		return nil, false
	}
	return map[int]interface{}{}, from == 0
}

// Reset removes all objects in storage.
func Reset() {
	storedObjects = []interface{}{}
}

// Len gives the total number of objects stored.
func Len() int {
	return len(storedObjects)
}

func tomap(start, end int) map[int]interface{} {
	m := make(map[int]interface{})
	for i := start; i < end; i++ {
		m[i] = storedObjects[i]
	}
	return m
}
