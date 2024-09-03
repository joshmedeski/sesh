package tmuxinator

import (
	"github.com/joshmedeski/sesh/model"
)

func (t *RealTmuxinator) ListSessions() ([]*model.TmuxinatorSession, error) {
  res, err :=  t.shell.ListCmd("bash", "-c", "tmuxinator list | tail -n +2 | tr -s '[:space:]' '\\n'")
  if err != nil {
    return []*model.TmuxinatorSession{}, err   
  }
  return parseTmuxinatorSessionsOutput(res) 
}

func parseTmuxinatorSessionsOutput(rawList []string) ([]*model.TmuxinatorSession, error) {
	sessions := make([]*model.TmuxinatorSession, 0, len(rawList))
	for _, line := range rawList {
    if len(line) > 0 {
      session := &model.TmuxinatorSession{
        Name: line,
      }
      sessions = append(sessions, session)
    }
	}

	return sessions, nil
}
