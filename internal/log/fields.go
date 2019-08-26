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

package log

// Name of fields, exhaustive list.
const (
	FieldApp       = "app"        // Name of the application.
	FieldVersion   = "version"    // Version of the application.
	FieldCommit    = "commit"     // Commit hash.
	FieldBuildDate = "build_date" // Build date.
	FieldBuiltBy   = "built_by"   // Identity of build owner.

	FieldService = "service" // Name of the service.
	FieldListen  = "addr"    // Listened address and port.
	FieldServer  = "server"  // Name of the server.
	FieldError   = "error"   // Error causing the event.

	FieldSession  = "session"  // Current session.
	FieldCommand  = "command"  // Current command being processed.
	FieldResponse = "response" // Current response (to be) emited.
)

// Fields is used to define the content of an event with structured fields.
// It can be used as a shorthand for map[string]interface{}.
type Fields map[string]interface{}
