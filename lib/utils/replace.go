package utils

import (
	"Kotone-DiVE/lib/config"
	"Kotone-DiVE/lib/db"
	"regexp"
	"strconv"
	"time"
)

var (
	configRegexBefore map[*regexp.Regexp]string
	configRegexAfter  map[*regexp.Regexp]string
)

func init() {
	configRegexBefore = map[*regexp.Regexp]string{}
	for k, v := range config.CurrentConfig.Replace.Before {
		configRegexBefore[regexp.MustCompile(k)] = v
	}
	configRegexAfter = map[*regexp.Regexp]string{}
	for k, v := range config.CurrentConfig.Replace.After {
		configRegexAfter[regexp.MustCompile(k)] = v
	}
}

func Replace(id string, list map[string]string, content string, trace bool) (string, string) {
	var (
		logStr string
		start  time.Time
		oldStr string
	) // debug
	if trace {
		start = time.Now()
		logStr = "Regex Replace() trace started at " + start.String() + " with string \"" + content + "\".\nGuildId is: " + id + ".\n"
	}

	for k, v := range configRegexBefore {
		if trace {
			oldStr = content
		}
		content = k.ReplaceAllString(content, v)
		if trace && content != oldStr {
			logStr += "Regex hit!\n|-Regex: \"" + k.String() + "\"\n|-Replace: \"" + v + "\"\n|-oldStr: \"" + oldStr + "\"\n|-New: " + content + "\n|-Time: " + strconv.FormatInt(time.Since(start).Nanoseconds(), 10) + "ns\n\n"
		}
	}
	if trace {
		logStr += "Processed config before regex(s) in " + strconv.FormatInt(time.Since(start).Nanoseconds(), 10) + "ns.\n"
	}

	val, exists := db.RegexCache[id]
	compiled := map[*regexp.Regexp]*string{}
	if exists {
		compiled = *val
		if trace {
			logStr += "Guild regex cache found.\nGot cache in " + strconv.FormatInt(time.Since(start).Nanoseconds(), 10) + "ns.\n"
		}
	} else {
		if trace {
			logStr += "Guild regex cache not found.\nCompiling regex...\n\n"
		}
		for k, v := range list {
			if trace {
				logStr += "Compiling \"" + k + "\" ...\n"
			}
			regex, err := regexp.Compile(k)
			if err == nil {
				text := v //Let's encrypt knows everything
				compiled[regex] = &text
			} else {
				if trace {
					logStr += "|-Error occurred while compiling.\n" + err.Error() + ".\n|-Skipping...\n"
				}
			}
		}
		db.RegexCache[id] = &compiled
		if trace {
			logStr += "Compiled regex in " + strconv.FormatInt(time.Since(start).Nanoseconds(), 10) + "ns.\n\n"
		}
	}
	if trace {
		logStr += "Starting process of " + strconv.Itoa(len(compiled)) + " user regex(s).\n\nRegex(s):\n"
		for k, v := range compiled {
			logStr += "|- \"" + k.String() + "\" => \"" + *v + "\"\n"
		}
		logStr += "\n"
	}
	for k, v := range compiled {
		if trace {
			oldStr = content
		}
		content = k.ReplaceAllString(content, *v)
		if trace && content != oldStr {
			logStr += "Regex hit!\n|-Regex: \"" + k.String() + "\"\n|-Replace: \"" + *v + "\"\n|-oldStr: \"" + oldStr + "\"\n|-New: \"" + content + "\"\n|-Time: " + strconv.FormatInt(time.Since(start).Nanoseconds(), 10) + "ns\n\n"
		}
	}
	if trace {
		logStr += "Processed user regex(s) in " + strconv.FormatInt(time.Since(start).Nanoseconds(), 10) + "ns.\n"
	}

	for k, v := range configRegexAfter {
		if trace {
			oldStr = content
		}
		content = k.ReplaceAllString(content, v)
		if trace && content != oldStr {
			logStr += "Regex hit!\n|-Regex: \"" + k.String() + "\"\n|-Replace: \"" + v + "\"\n|-oldStr: \"" + oldStr + "\"\n|-New: " + content + "\n|-Time: " + strconv.FormatInt(time.Since(start).Nanoseconds(), 10) + "ns\n\n"
		}
	}
	if trace {
		logStr += "Processed config after regex(s) in " + strconv.FormatInt(time.Since(start).Nanoseconds(), 10) + "ns.\nReplace() ended at " + time.Now().String() + " with string \"" + content + "\".\n"
	}
	return content, logStr
}
