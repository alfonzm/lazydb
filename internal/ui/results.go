package ui

import (
	"github.com/alfonzm/lazydb/internal/db"
	"github.com/rivo/tview"
)

type Results struct {
	view *tview.Box
	db   *db.DBClient
}

func NewResults(db *db.DBClient) (*Results, error) {
	view := tview.NewBox()
	view.SetBorder(true).SetTitle("Results")

	return &Results{
		view: view,
		db:   db,
	}, nil
}
