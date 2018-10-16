package action

import (
	"github.com/jaitl/goEnglishBot/app/telegram/command"
)

type Executor struct {
	session         *SessionModel
	typeToAction    map[Type]Action
	commandToAction map[command.Type]Action
}

func NewExecutor(session *SessionModel, actions []Action) *Executor {
	type2acts := make(map[Type]Action)
	cmd2acts := make(map[command.Type]Action)

	for _, act := range actions {
		type2acts[act.GetType()] = act

		for _, cmd := range act.GetStartCommands() {
			cmd2acts[cmd] = act
		}
	}

	return &Executor{
		session:         session,
		typeToAction:    type2acts,
		commandToAction: cmd2acts,
	}
}

func (e *Executor) Execute(cmd command.Command) error {
	ses, err := e.session.FindSession(cmd.GetUserId())

	if err != nil {
		return err
	}

	if ses != nil {
		action, ok := e.typeToAction[ses.ActionType]
		if ok {
			waitCommands := action.GetWaitCommands(ses.Stage)
			ok = waitCommands[cmd.GetType()]

			if ok {
				err = action.Execute(ses.Stage, cmd, ses)
				return err
			} else {
				err = e.session.ClearSession(cmd.GetUserId())

				if err != nil {
					return err
				}

				action = e.commandToAction[cmd.GetType()]
				err = action.Execute(action.GetStartStage(), cmd, nil)

				return err
			}
		}
	} else {
		action := e.commandToAction[cmd.GetType()]
		err = action.Execute(action.GetStartStage(), cmd, nil)
		return err
	}

	return nil
}
