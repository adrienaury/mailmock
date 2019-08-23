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
//
// Linking this library statically or dynamically with other modules is
// making a combined work based on this library.  Thus, the terms and
// conditions of the GNU General Public License cover the whole
// combination.
//
// As a special exception, the copyright holders of this library give you
// permission to link this library with independent modules to produce an
// executable, regardless of the license terms of these independent
// modules, and to copy and distribute the resulting executable under
// terms of your choice, provided that you also meet, for each linked
// independent module, the terms and conditions of the license of that
// module.  An independent module is a module which is not derived from
// or based on this library.  If you modify this library, you may extend
// this exception to your version of the library, but you are not
// obligated to do so.  If you do not wish to do so, delete this
// exception statement from your version.

package log

// Logger is the fundamental interface for all log operations.
type Logger interface {
	Trace(msg string, fields ...map[string]interface{})
	Debug(msg string, fields ...map[string]interface{})
	Info(msg string, fields ...map[string]interface{})
	Warn(msg string, fields ...map[string]interface{})
	Error(msg string, fields ...map[string]interface{})
}

// NoopLogger is a Logger implementation that does nothing.
type NoopLogger struct{}

// Trace for NoopLogger does nothing.
func (NoopLogger) Trace(msg string, fields ...map[string]interface{}) {}

// Debug for NoopLogger does nothing.
func (NoopLogger) Debug(msg string, fields ...map[string]interface{}) {}

// Info for NoopLogger does nothing.
func (NoopLogger) Info(msg string, fields ...map[string]interface{}) {}

// Warn for NoopLogger does nothing.
func (NoopLogger) Warn(msg string, fields ...map[string]interface{}) {}

// Error for NoopLogger does nothing.
func (NoopLogger) Error(msg string, fields ...map[string]interface{}) {}

// DefaultLogger is the logger user by this package.
var DefaultLogger Logger = NoopLogger{}
