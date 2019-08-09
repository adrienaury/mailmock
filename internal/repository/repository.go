package repository

var storedObjects = []interface{}{}

// Store store th object and give an ID for
func Store(o interface{}) int {
	id := len(storedObjects)
	storedObjects = append(storedObjects, o)
	return id
}

// Use returns the object with ID or nil
func Use(ID int) interface{} {
	if ID < len(storedObjects) {
		return storedObjects[ID]
	}
	return nil
}

// All returns all objects currently stored
func All() []interface{} {
	return storedObjects
}
