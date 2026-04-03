package cmd

import (
	"fmt"
	"log"

	"github.com/jroimartin/gocui"
	"github.com/spf13/cobra"
)

func CreateCUI() *cobra.Command {
	cmd := &cobra.Command{
		Use:   `cui`,
		Short: `example for github.com/jroimartin/gocui`,
		Run: func(cmd *cobra.Command, args []string) {
			g, err := gocui.NewGui(gocui.OutputNormal)
			if err != nil {
				log.Panicln(err)
			}
			defer g.Close()

			layout := func(g *gocui.Gui) error {
				maxX, maxY := g.Size()
				if v, err := g.SetView("hello", maxX/2-7, maxY/2, maxX/2+7, maxY/2+2); err != nil {
					if err != gocui.ErrUnknownView {
						return err
					}
					fmt.Fprintln(v, "Hello world!")
				}
				return nil
			}

			quit := func(g *gocui.Gui, v *gocui.View) error {
				return gocui.ErrQuit
			}

			g.SetManagerFunc(layout)

			if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
				log.Panicln(err)
			}

			if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
				log.Panicln(err)
			}
		},
	}
	return cmd
}
