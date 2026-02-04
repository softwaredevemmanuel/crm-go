// utils/time_parser.go
package utils

import (
    "errors"
    "strings"
    "time"
	"fmt"
)

// ParseTime parses a time string using multiple common formats
// Supports: RFC3339, ISO8601, common date formats, and Unix timestamps
func ParseTime(timeStr string) (time.Time, error) {
    if timeStr == "" {
        return time.Time{}, errors.New("time string is empty")
    }
    
    // Trim whitespace
    timeStr = strings.TrimSpace(timeStr)
    
    // List of formats to try (in order of preference)
    formats := []string{
        // RFC3339 formats (most common for APIs)
        time.RFC3339,                    // "2006-01-02T15:04:05Z07:00"
        "2006-01-02T15:04:05Z",         // UTC with Z
        "2006-01-02T15:04:05",          // No timezone
        
        // ISO8601 variations
        "2006-01-02 15:04:05Z07:00",    // Space instead of T
        "2006-01-02 15:04:05",          // Space, no timezone
        
        // Date only formats
        "2006-01-02",                   // Date only
        "2006/01/02",                   // Slash separator
        "02-01-2006",                   // DD-MM-YYYY
        "02/01/2006",                   // DD/MM/YYYY
        "01-02-2006",                   // MM-DD-YYYY (US format)
        "01/02/2006",                   // MM/DD/YYYY (US format)
        
        // Time with date
        "2006-01-02 15:04",             // Date with hours and minutes
        "2006-01-02 3:04 PM",           // 12-hour format
        "2006-01-02 3:04:05 PM",        // 12-hour with seconds
        
        // Common HTTP/Log formats
        time.RFC1123,                   // "Mon, 02 Jan 2006 15:04:05 MST"
        time.RFC1123Z,                  // "Mon, 02 Jan 2006 15:04:05 -0700"
        time.RFC822,                    // "02 Jan 06 15:04 MST"
        time.RFC822Z,                   // "02 Jan 06 15:04 -0700"
        
        // Unix timestamp (seconds)
        "unix",
        "unix_seconds",
        
        // Unix millisecond timestamp
        "unix_ms",
        "unix_milliseconds",
    }
    
    // Special case: Unix timestamps
    if isUnixTimestamp(timeStr) {
        return parseUnixTimestamp(timeStr)
    }
    
    // Try each format
    for _, format := range formats {
        t, err := time.Parse(format, timeStr)
        if err == nil {
            return t, nil
        }
    }
    
    // If still not parsed, try with location
    return parseWithLocation(timeStr)
}

// ParseTimeWithLocation parses time string with a specific timezone
func ParseTimeWithLocation(timeStr, locationStr string) (time.Time, error) {
    t, err := ParseTime(timeStr)
    if err != nil {
        return t, err
    }
    
    // If no location provided, return as-is
    if locationStr == "" {
        return t, nil
    }
    
    // Try to load location
    loc, err := time.LoadLocation(locationStr)
    if err != nil {
        // Try common timezone abbreviations
        loc = getLocationFromAbbreviation(locationStr)
        if loc == nil {
            // Return original time if location invalid
            return t, nil
        }
    }
    
    // Convert to specified location
    return t.In(loc), nil
}

// ParseTimeStrict parses only RFC3339/ISO8601 formats
func ParseTimeStrict(timeStr string) (time.Time, error) {
    if timeStr == "" {
        return time.Time{}, errors.New("time string is empty")
    }
    
    timeStr = strings.TrimSpace(timeStr)
    
    strictFormats := []string{
        time.RFC3339,
        time.RFC3339Nano,
        "2006-01-02T15:04:05Z",
        "2006-01-02T15:04:05",
        "2006-01-02 15:04:05Z",
        "2006-01-02 15:04:05",
    }
    
    for _, format := range strictFormats {
        t, err := time.Parse(format, timeStr)
        if err == nil {
            return t, nil
        }
    }
    
    return time.Time{}, errors.New("invalid time format, expected RFC3339 like '2024-03-15T14:00:00Z' or '2024-03-15T14:00:00+05:30'")
}

// Helper function to check if string is a Unix timestamp
func isUnixTimestamp(timeStr string) bool {
    // Check if it's all digits (Unix timestamp in seconds or milliseconds)
    for _, c := range timeStr {
        if c < '0' || c > '9' {
            return false
        }
    }
    return len(timeStr) <= 13 // Unix timestamp max 13 digits (milliseconds until year 2286)
}

// Helper to parse Unix timestamp (seconds or milliseconds)
func parseUnixTimestamp(timeStr string) (time.Time, error) {
    var seconds int64
    var milliseconds int64
    
    // Parse as int64
    if len(timeStr) <= 10 {
        // Seconds timestamp (up to 10 digits)
        _, err := fmt.Sscanf(timeStr, "%d", &seconds)
        if err != nil {
            return time.Time{}, err
        }
        return time.Unix(seconds, 0), nil
    } else {
        // Milliseconds timestamp (11-13 digits)
        _, err := fmt.Sscanf(timeStr, "%d", &milliseconds)
        if err != nil {
            return time.Time{}, err
        }
        seconds = milliseconds / 1000
        nanos := (milliseconds % 1000) * 1e6
        return time.Unix(seconds, nanos), nil
    }
}

// Helper to parse time with location detection
func parseWithLocation(timeStr string) (time.Time, error) {
    // Try to detect common timezone abbreviations in the string
    tzAbbreviations := map[string]string{
        "UTC":   "UTC",
        "GMT":   "GMT",
        "EST":   "America/New_York",
        "EDT":   "America/New_York",
        "CST":   "America/Chicago",
        "CDT":   "America/Chicago",
        "MST":   "America/Denver",
        "MDT":   "America/Denver",
        "PST":   "America/Los_Angeles",
        "PDT":   "America/Los_Angeles",
        "IST":   "Asia/Kolkata",
        "JST":   "Asia/Tokyo",
        "CET":   "Europe/Paris",
        "CEST":  "Europe/Paris",
        "AEST":  "Australia/Sydney",
        "AEDT":  "Australia/Sydney",
    }
    
    // Check if string ends with timezone abbreviation
    for tzAbbr, tzName := range tzAbbreviations {
        if strings.HasSuffix(strings.ToUpper(timeStr), tzAbbr) {
            // Remove timezone from string
            cleanStr := strings.TrimSuffix(timeStr, tzAbbr)
            cleanStr = strings.TrimSuffix(cleanStr, " ")
            
            // Try parsing without timezone
            t, err := ParseTime(cleanStr)
            if err != nil {
                return t, err
            }
            
            // Convert to detected timezone
            loc, _ := time.LoadLocation(tzName)
            if loc != nil {
                return t.In(loc), nil
            }
            return t, nil
        }
    }
    
    return time.Time{}, errors.New("unable to parse time string: " + timeStr)
}

// Helper to get location from common abbreviations
func getLocationFromAbbreviation(abbr string) *time.Location {
    tzMap := map[string]string{
        // US Timezones
        "EST": "America/New_York",
        "EDT": "America/New_York",
        "CST": "America/Chicago", 
        "CDT": "America/Chicago",
        "MST": "America/Denver",
        "MDT": "America/Denver",
        "PST": "America/Los_Angeles",
        "PDT": "America/Los_Angeles",
        
        // World Timezones
        "UTC": "UTC",
        "GMT": "GMT",
        "IST": "Asia/Kolkata",
        "JST": "Asia/Tokyo",
        "CET": "Europe/Paris",
        "CEST": "Europe/Paris",
        "AEST": "Australia/Sydney",
        "AEDT": "Australia/Sydney",
        "BST": "Europe/London",
        "BRT": "America/Sao_Paulo",
        "HKT": "Asia/Hong_Kong",
        "SGT": "Asia/Singapore",
    }
    
    if tzName, exists := tzMap[strings.ToUpper(abbr)]; exists {
        loc, _ := time.LoadLocation(tzName)
        return loc
    }
    
    return nil
}

// FormatTimeForAPI formats time for API responses (RFC3339)
func FormatTimeForAPI(t time.Time) string {
    return t.Format(time.RFC3339)
}

// FormatTimeForDisplay formats time for human-readable display
func FormatTimeForDisplay(t time.Time, timezone string) string {
    layout := "Monday, January 2, 2006 at 3:04 PM"
    
    if timezone != "" {
        loc, err := time.LoadLocation(timezone)
        if err == nil {
            t = t.In(loc)
        }
    }
    
    return t.Format(layout)
}

// IsValidTimeFormat checks if a string is a valid time format
func IsValidTimeFormat(timeStr string) bool {
    _, err := ParseTime(timeStr)
    return err == nil
}

// GetSupportedTimeFormats returns a list of supported time formats
func GetSupportedTimeFormats() []string {
    return []string{
        "RFC3339: 2006-01-02T15:04:05Z07:00",
        "UTC with Z: 2006-01-02T15:04:05Z", 
        "No timezone: 2006-01-02T15:04:05",
        "Space separator: 2006-01-02 15:04:05",
        "Date only: 2006-01-02",
        "US format: 01/02/2006",
        "Unix timestamp (seconds): 1678886400",
        "Unix milliseconds: 1678886400000",
    }
}