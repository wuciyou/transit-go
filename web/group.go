package web

import "sync"

type group struct {
	sync.Mutex
	groupSession map[string]map[string]string
	sessionHaveGroup []string
}

func (g *group) add (groupName string, connectKey string){
	g.Lock()

	if gs,ok := g.groupSession[groupName]; !ok{
		g.groupSession[groupName] = make(map[string]string)
		g.groupSession[groupName][connectKey] = connectKey
	}else{
		gs[connectKey] = connectKey
	}

	g.sessionHaveGroup = append(g.sessionHaveGroup,groupName)

	g.Unlock()

}

func (g *group) del(groupName string, connectKey string){




}
