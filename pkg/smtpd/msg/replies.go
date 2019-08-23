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

package msg

var (
	// GreetingBanner is the first reply after a client opens a new connection.
	GreetingBanner = "Service ready"

	// Success is the default reply after successful processing.
	Success = "OK"

	// AskForData is replied to a DATA command before receiving mail content.
	AskForData = "Start mail input; end with <CRLF>.<CRLF>"

	// Goodbye is the last reply sent before closing the transmission channel.
	Goodbye = "Service closing transmission channel"

	// ServiceNotAvailable is sent when the server cannot accept new connection.
	ServiceNotAvailable = "Service not available, closing transmission channel"

	// RequestedActionAborted is sent when an error occured processing the request.
	RequestedActionAborted = "Requested action aborted: error in processing"

	// BadSequence is sent when the last command was correct, but not expected.
	BadSequence = "Bad sequence of commands"

	// NoValidRecipients is sent when no valid recipients were issued.
	NoValidRecipients = "No valid recipients"

	// NotImplemented is sent when the command is recognized, but not implemented by the server.
	NotImplemented = "Command not implemented"

	// NotRecognized is sent when the command is not recognized.
	NotRecognized = "Syntax error, command unrecognized"

	// ParameterError is sent when the command is recognized, but arguments or parameters are invalid.
	ParameterError = "Syntax error in parameters or arguments"
)
