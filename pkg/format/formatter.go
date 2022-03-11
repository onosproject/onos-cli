// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package format

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"text/tabwriter"
	"text/template"
	"text/template/parse"
)

var nameFinder = regexp.MustCompile(`\.([\._A-Za-z0-9]*)}}`)

// Format defines a type for a string that can be used as template to format data.
type Format string

// TrimAndPad modifies `s` so that it is exactly `l` characters long, removing
// characters from the end, or adding spaces as necessary.
func TrimAndPad(s string, l int) string {
	// TODO: support right justification if a negative number is passed
	if len(s) > l {
		s = s[:l]
	}
	return s + strings.Repeat(" ", l-len(s))
}

// GetHeaderString extract the set of column names from a template.
func GetHeaderString(tmpl *template.Template, nameLimit int) string {
	var header string
	for _, n := range tmpl.Tree.Root.Nodes {
		switch n.Type() {
		case parse.NodeText:
			header += n.String()
		case parse.NodeString:
			header += n.String()
		case parse.NodeAction:
			found := nameFinder.FindStringSubmatch(n.String())
			if len(found) == 2 {
				if nameLimit > 0 {
					parts := strings.Split(found[1], ".")
					start := len(parts) - nameLimit
					if start < 0 {
						start = 0
					}
					header += strings.ToUpper(strings.Join(parts[start:], "."))
				} else {
					header += strings.ToUpper(found[1])
				}
			}
		}
	}
	return header
}

// IsTable returns a bool if the template is a table
func (f Format) IsTable() bool {
	return strings.HasPrefix(string(f), "table")
}

// Execute compiles the template and prints the output
func (f Format) Execute(writer io.Writer, withHeaders bool, nameLimit int, data interface{}) error {
	var tabWriter *tabwriter.Writer
	format := f

	if f.IsTable() {
		tabWriter = tabwriter.NewWriter(writer, 0, 4, 4, ' ', 0)
		format = Format(strings.TrimPrefix(string(f), "table"))
	}

	funcmap := template.FuncMap{
		"timestamp": formatTimestamp,
		"since":     formatSince,
		"gosince":   formatGoSince}

	tmpl, err := template.New("output").Funcs(funcmap).Parse(string(format))
	if err != nil {
		return err
	}

	if f.IsTable() && withHeaders {
		header := GetHeaderString(tmpl, nameLimit)

		if _, err = tabWriter.Write([]byte(header)); err != nil {
			return err
		}
		if _, err = tabWriter.Write([]byte("\n")); err != nil {
			return err
		}

		slice := reflect.ValueOf(data)
		if slice.Kind() == reflect.Slice {
			for i := 0; i < slice.Len(); i++ {
				if err = tmpl.Execute(tabWriter, slice.Index(i).Interface()); err != nil {
					return err
				}
				if _, err = tabWriter.Write([]byte("\n")); err != nil {
					return err
				}
			}
		} else {
			if err = tmpl.Execute(tabWriter, data); err != nil {
				return err
			}
			if _, err = tabWriter.Write([]byte("\n")); err != nil {
				return err
			}
		}
		tabWriter.Flush()
		return nil
	}

	slice := reflect.ValueOf(data)
	if slice.Kind() == reflect.Slice {
		for i := 0; i < slice.Len(); i++ {
			if err = tmpl.Execute(writer, slice.Index(i).Interface()); err != nil {
				return err
			}
			if _, err = writer.Write([]byte("\n")); err != nil {
				return err
			}
		}
	} else {
		if err = tmpl.Execute(writer, data); err != nil {
			return err
		}
		if _, err = writer.Write([]byte("\n")); err != nil {
			return err
		}
	}
	return nil

}

// ExecuteFixedWidth Formats a table row using a set of fixed column widths.
// Used for streaming
// output where column widths cannot be automatically determined because only
// one line of the output is available at a time.
//
// Assumes the format uses tab as a field delimiter.
//
// columnWidths: struct that contains column widths
// header: If true return the header. If false then evaluate data and return data.
// data: Data to evaluate
func (f Format) ExecuteFixedWidth(columnWidths interface{}, header bool, data interface{}) (string, error) {
	if !f.IsTable() {
		return "", errors.New("Fixed width is only available on table format")
	}

	outputAs := strings.TrimPrefix(string(f), "table")
	tmpl, err := template.New("output").Parse(outputAs)
	if err != nil {
		return "", fmt.Errorf("Failed to parse template: %v", err)
	}

	var buf bytes.Buffer
	var tabSepOutput string

	if header {
		// Caller wants the table header.
		tabSepOutput = GetHeaderString(tmpl, 1)
	} else {
		// Caller wants the data.
		err = tmpl.Execute(&buf, data)
		if err != nil {
			return "", fmt.Errorf("Failed to execute template: %v", err)
		}
		tabSepOutput = buf.String()
	}

	// Extract the column width constants by running the template on the
	// columnWidth structure. This will cause text.template to split the
	// column widths exactly like it did the output (i.e. separated by
	// tab characters)
	buf.Reset()
	err = tmpl.Execute(&buf, columnWidths)
	if err != nil {
		return "", fmt.Errorf("Failed to execute template on widths: %v", err)
	}
	tabSepWidth := buf.String()

	// Loop through the fields and widths, printing each field to the
	// preset width.
	output := ""
	outParts := strings.Split(tabSepOutput, "\t")
	widthParts := strings.Split(tabSepWidth, "\t")
	for i, outPart := range outParts {
		width, err := strconv.Atoi(widthParts[i])
		if err != nil {
			return "", fmt.Errorf("Failed to parse width %s: %v", widthParts[i], err)
		}
		output = output + TrimAndPad(outPart, width) + " "
	}

	// remove any trailing spaces
	output = strings.TrimRight(output, " ")

	return output, nil
}
