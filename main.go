package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"syscall"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var (
	app              *tview.Application
	fooView, barView *tview.TextView
	menuList         *tview.List
	currentPage      string = "foo"
)

var globalCmdFoo, globalCmdBar *exec.Cmd

func NewLayout(list tview.Primitive, panel tview.Primitive) *tview.Flex {
	bodyLayout := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(list, 20, 1, true).
		AddItem(panel, 0, 1, false)

	layout := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(bodyLayout, 0, 1, true)

	return layout
}

func NewList(logPages *tview.Pages) *tview.List {
	list := tview.NewList()
	menuList = list
	list.SetBorder(true)
	list.SetTitle("Application")

	list.AddItem("foo", "", '1', func() {
		currentPage = "foo"
		logPages.SwitchToPage("foo")
		list.SetCurrentItem(0)
	})
	list.AddItem("bar", "", '2', func() {
		currentPage = "bar"
		logPages.SwitchToPage("bar")
		list.SetCurrentItem(1)
	})
	list.AddItem("Quit", "", 'q', func() {
		if globalCmdFoo != nil && globalCmdFoo.Process != nil {
			globalCmdFoo.Process.Signal(syscall.SIGINT)
		}
		if globalCmdBar != nil && globalCmdBar.Process != nil {
			globalCmdBar.Process.Signal(syscall.SIGINT)
		}
		app.Stop()
	})
	return list
}

func NewLogPanel(app *tview.Application) *tview.Pages {
	fooView = tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true).
		SetChangedFunc(func() { app.Draw() })
	fooView.SetTitle("foo").SetBorder(true)

	barView = tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true).
		SetChangedFunc(func() { app.Draw() })
	barView.SetTitle("bar").SetBorder(true)

	pages := tview.NewPages().
		AddPage("foo", fooView, true, true).
		AddPage("bar", barView, true, false)
	return pages
}

func setupInputCapture(app *tview.Application, pages *tview.Pages) {
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTAB {
			if currentPage == "foo" {
				currentPage = "bar"
				menuList.SetCurrentItem(1)
			} else {
				currentPage = "foo"
				menuList.SetCurrentItem(0)
			}
			pages.SwitchToPage(currentPage)
			return nil
		}
		return event
	})
}

func runCommand(cmdPath string, view *tview.TextView, app *tview.Application, wg *sync.WaitGroup) *exec.Cmd {
	cmd := exec.Command(cmdPath)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatalf("Error creating stdout pipe for %s: %v", cmdPath, err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Fatalf("Error creating stderr pipe for %s: %v", cmdPath, err)
	}

	if err := cmd.Start(); err != nil {
		log.Fatalf("Error starting command %s: %v", cmdPath, err)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			line := scanner.Text()
			app.QueueUpdateDraw(func() {
				fmt.Fprintln(view, line)
			})
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			line := scanner.Text()
			app.QueueUpdateDraw(func() {
				fmt.Fprintln(view, line)
			})
		}
	}()

	return cmd
}

func NewApp() *tview.Application {
	app = tview.NewApplication()

	logPanel := NewLogPanel(app)
	list := NewList(logPanel)
	setupInputCapture(app, logPanel)

	layout := NewLayout(list, logPanel)
	pages := tview.NewPages().AddPage("main", layout, true, true)
	app.SetRoot(pages, true)
	return app
}

func main() {
	app := NewApp()

	var wg sync.WaitGroup

	globalCmdFoo = runCommand("./bin/foo", fooView, app, &wg)
	globalCmdBar = runCommand("./bin/bar", barView, app, &wg)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-signals
		if globalCmdFoo != nil && globalCmdFoo.Process != nil {
			globalCmdFoo.Process.Signal(syscall.SIGINT)
		}
		if globalCmdBar != nil && globalCmdBar.Process != nil {
			globalCmdBar.Process.Signal(syscall.SIGINT)
		}
		app.Stop()
	}()

	if err := app.Run(); err != nil {
		log.Fatalf("Error running application: %v", err)
	}

	wg.Wait()
}
