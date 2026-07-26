package connector

import "strings"

// resolveAlias maps a configured alias to its session name so that
// `sesh connect wp` behaves exactly like `sesh connect wallpaper`, running the
// full strategy chain against the real name. Anything that isn't an alias is
// returned untouched.
func (c *RealConnector) resolveAlias(name string) string {
	if name == "" {
		return name
	}
	for _, session := range c.config.SessionConfigs {
		if session.Alias != "" && strings.EqualFold(session.Alias, name) {
			return session.Name
		}
	}
	return name
}
