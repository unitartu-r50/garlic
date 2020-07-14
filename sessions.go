package main

import (
	"github.com/google/uuid"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

// sessions represent all sessions available and is the main type for a user which simplifies
// the control of a Pepper robot.
var sessions []*Session

// moves represent all ready-made moves located somewhere on the disk. This presented in the
// web UI as a library of moves which can be called any time by a user.
var moves Moves

// moveGroups is a helper variable for the "/sessions/" route to list moves by a group.
var moveGroups []string

func collectSessions(sayDir string, moves *Moves) ([]*Session, error) {
	var sessions = []*Session{
		{
			Name: "Session 1",
			Items: []*SessionItem{
				{
					Question: &SayAndMoveAction{
						SayItem: &SayAction{
							Phrase:   "Tere, mina olen robot Pepper. Mina olen 6-aastane ja tahan sinuga tuttavaks saada. Mis on sinu nimi?",
							FilePath: "1out_tutvustus.wav",
						},
						MoveItem: &MoveAction{
							Name:  "Hello_01",
							Delay: 0,
						},
					},
					PositiveAnswer: &SayAndMoveAction{
						SayItem: &SayAction{
							Phrase: "Nice",
						},
						MoveItem: &MoveAction{
							Name: "NiceReaction_01",
						},
					},
					NegativeAnswer: &SayAndMoveAction{
						SayItem: &SayAction{
							Phrase: "Sad",
						},
						MoveItem: &MoveAction{
							Name: "SadReaction_01",
						},
					},
				},
				{
					Question: &SayAndMoveAction{
						SayItem: &SayAction{
							Phrase:   "Kui vana sa oled?",
							FilePath: "2out_vanus.wav",
						},
						MoveItem: &MoveAction{
							Name:  "Show_Hand_Right_02",
							Delay: 0,
						},
					},
				},
				{
					Question: &SayAndMoveAction{
						SayItem: &SayAction{
							Phrase:   "Kas Sul on vendi või õdesid?",
							FilePath: "3out_vennad.wav",
						},
						MoveItem: &MoveAction{
							Name:  "Show_Hand_Both_02",
							Delay: 0,
						},
					},
				},
				{
					Question: &SayAndMoveAction{
						SayItem: &SayAction{
							Phrase:   "Ma tulin siia üksi, kuid mu pere on suur ja mööda maailma laiali.",
							FilePath: "3out_vennadVV.wav",
						},
						MoveItem: &MoveAction{
							Name:  "Show_Hand_Both_01",
							Delay: 0,
						},
					},
				},
				{
					Question: &SayAndMoveAction{
						SayItem: &SayAction{
							Phrase:   "Mina olen pärit Pariisist ja nüüd meeldib mulle väga Eestis elada. Mis sulle Sinu Eestimaa juures meeldib?",
							FilePath: "4out_päritolu.wav",
						},
						MoveItem: &MoveAction{
							Name:  "Show_Self_01",
							Delay: time.Second * 5,
						},
					},
				},
				{
					Question: &SayAndMoveAction{
						SayItem: &SayAction{
							Phrase:   "Jaa, see on väike ja sõbralik maa ja teil on 4 aastaaega",
							FilePath: "5out_eestimaavastus.wav",
						},
						MoveItem: &MoveAction{
							Name:  "NiceReaction_01",
							Delay: 0,
						},
					},
				},
			},
		},
		{
			ID:   uuid.Must(uuid.NewRandom()),
			Name: "Session 2",
			Items: []*SessionItem{
				{
					ID: uuid.Must(uuid.NewRandom()),
					Question: &SayAndMoveAction{
						SayItem: &SayAction{
							Phrase: "Q1",
						},
					},
				},
				{
					ID: uuid.Must(uuid.NewRandom()),
					Question: &SayAndMoveAction{
						SayItem: &SayAction{
							Phrase: "Q2",
						},
					},
				},
			},
		},
	}

	for _, s := range sessions {
		s.ID = uuid.Must(uuid.NewRandom())

		for _, item := range s.Items {
			// initiate unique IDs
			item.Question.SetID(uuid.Must(uuid.NewRandom()))
			if !item.Question.IsNil() {
				item.Question.SayItem.SetID(uuid.Must(uuid.NewRandom()))
				item.Question.MoveItem.SetID(uuid.Must(uuid.NewRandom()))
			}

			item.PositiveAnswer.SetID(uuid.Must(uuid.NewRandom()))
			if !item.PositiveAnswer.IsNil() {
				item.PositiveAnswer.SayItem.SetID(uuid.Must(uuid.NewRandom()))
				item.PositiveAnswer.MoveItem.SetID(uuid.Must(uuid.NewRandom()))
			}

			item.NegativeAnswer.SetID(uuid.Must(uuid.NewRandom()))
			if !item.NegativeAnswer.IsNil() {
				item.NegativeAnswer.SayItem.SetID(uuid.Must(uuid.NewRandom()))
				item.NegativeAnswer.MoveItem.SetID(uuid.Must(uuid.NewRandom()))
			}

			if item.Question != nil {
				if item.Question.SayItem.FilePath != "" {
					item.Question.SayItem.FilePath = path.Join(sayDir, s.Name, item.Question.SayItem.FilePath)
					if _, err := os.Stat(item.Question.SayItem.FilePath); os.IsNotExist(err) {
						return nil, err
					}
				}

				if item.Question.MoveItem != nil {
					if v := moves.GetByName(item.Question.MoveItem.Name); v != nil {
						m := *v                                // copy values from the library
						m.Delay = item.Question.MoveItem.Delay // copy delay from a user provided variable
						item.Question.MoveItem = &m
					}
				}
			}

			if item.PositiveAnswer != nil {
				if item.PositiveAnswer.SayItem.FilePath != "" {
					item.PositiveAnswer.SayItem.FilePath = path.Join(sayDir, s.Name,
						item.PositiveAnswer.SayItem.FilePath)
					if _, err := os.Stat(item.PositiveAnswer.SayItem.FilePath); os.IsNotExist(err) {
						return nil, err
					}
				}

				if item.PositiveAnswer.MoveItem != nil {
					if v := moves.GetByName(item.PositiveAnswer.MoveItem.Name); v != nil {
						m := *v                                      // copy values from the library
						m.Delay = item.PositiveAnswer.MoveItem.Delay // copy delay from a user provided variable
						item.PositiveAnswer.MoveItem = &m
					}
				}
			}

			if item.NegativeAnswer != nil {
				if item.NegativeAnswer.SayItem.FilePath != "" {
					item.NegativeAnswer.SayItem.FilePath = path.Join(sayDir, s.Name, item.NegativeAnswer.SayItem.FilePath)
					if _, err := os.Stat(item.NegativeAnswer.SayItem.FilePath); os.IsNotExist(err) {
						return nil, err
					}
				}

				if item.NegativeAnswer.MoveItem != nil {
					if v := moves.GetByName(item.NegativeAnswer.MoveItem.Name); v != nil {
						m := *v                                      // copy values from the library
						m.Delay = item.NegativeAnswer.MoveItem.Delay // copy delay from a user provided variable
						item.NegativeAnswer.MoveItem = &m
					}
				}
			}
		}
	}
	return sessions, nil
}

func collectMoves(dataDir string) ([]*MoveAction, error) {
	query := path.Join(dataDir, "**/*.qianim")
	matches, err := filepath.Glob(query)
	if err != nil {
		return nil, err
	}

	var items = make([]*MoveAction, len(matches))
	for i := range matches {
		// parsing the parent folder as a motion group name
		parts := strings.Split(matches[i], "/")
		parent := parts[len(parts)-2]

		// parsing the basename as a motion name
		basename := parts[len(parts)-1]
		name := strings.Replace(basename, filepath.Ext(basename), "", 1)

		// appending a motion
		items[i] = &MoveAction{
			ID:       uuid.Must(uuid.NewRandom()),
			FilePath: matches[i],
			Group:    parent,
			Name:     name,
		}
	}

	return items, err
}

// Session

// Session represents a session with a child, a set of questions and simple answers which
// are accompanied with moves by a robot to make the conversation a lively one.
type Session struct {
	ID          uuid.UUID
	Name        string
	Description string
	Items       []*SessionItem
}

// Sessions is a wrapper struct around an array of sessions with helpful methods.
type Sessions []*Session

// GetInstructionByID looks for a top level instruction, which unites Say and Move actions
// and presents them as a union of two actions, so both actions should be executed.
func (ss Sessions) GetInstructionByID(id uuid.UUID) *SayAndMoveAction {
	for _, session := range ss {
		for _, item := range session.Items {
			if !item.NegativeAnswer.IsNil() && item.NegativeAnswer.GetID() == id {
				return item.NegativeAnswer
			}
			if !item.PositiveAnswer.IsNil() && item.PositiveAnswer.GetID() == id {
				return item.PositiveAnswer
			}
			if !item.Question.IsNil() && item.Question.GetID() == id {
				return item.Question
			}
		}
	}
	return nil
}

// SessionItem represents a single unit of a session, it's a question and positive and negative
// answers accompanied with a robot's moves which are represented in the web UI as a set of buttons.
type SessionItem struct {
	ID             uuid.UUID
	Question       *SayAndMoveAction
	PositiveAnswer *SayAndMoveAction
	NegativeAnswer *SayAndMoveAction
}
