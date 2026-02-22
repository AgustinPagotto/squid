package main

type Choice int

const (
	ChoiceSaveNote Choice = iota
	ChoiceSeeNotes
	ChoiceDashboard
)

func (c Choice) String() string {
	return [...]string{"Save Note", "See Notes", "Dashboard"}[c]
}

func (c Choice) PageName() string {
	return [...]string{"notes-create", "notes", "dashboard"}[c]
}
