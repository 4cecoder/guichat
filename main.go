package main

import (
	"io/ioutil"
	"log"
	"os"
	"strings"
	"strconv"

	"github.com/gotk3/gotk3/gtk"
	//"github.com/gotk3/gotk3/pango"
)

const historyFilename = "chat_history.txt"

var username string = "Anonymous"
var serverAddress string = ""
var fontSize int = 10 // This is in points

func main() {
	gtk.Init(nil)

	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		log.Fatal("Unable to create window:", err)
	}
	win.SetTitle("Cute Cat Chat 😺💬")
	win.SetDefaultSize(450, 600)

	vbox, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 5)
	if err != nil {
		log.Fatal("Unable to create box:", err)
	}

	textview, err := gtk.TextViewNew()
	if err != nil {
		log.Fatal("Unable to create text view:", err)
	}
	textview.SetEditable(false)
	textview.SetCursorVisible(false)
	textview.SetWrapMode(gtk.WRAP_WORD_CHAR)

	// Apply font size using CSS
	cssProvider, err := gtk.CssProviderNew()
	if err != nil {
		log.Fatal("Unable to create css provider:", err)
	}
	err = cssProvider.LoadFromData("* { font-size: " + strconv.Itoa(fontSize) + "pt; }")
	if err != nil {
		log.Fatal("Unable to load css data:", err)
	}
	context, err := textview.GetStyleContext()
	if err != nil {
		log.Fatal("Unable to get style context:", err)
	}
	context.AddProvider(cssProvider, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)

	scrolledWindow, err := gtk.ScrolledWindowNew(nil, nil)
	if err != nil {
		log.Fatal("Unable to create scrolled window:", err)
	}
	scrolledWindow.SetPolicy(gtk.POLICY_NEVER, gtk.POLICY_AUTOMATIC)
	scrolledWindow.Add(textview)
	vbox.PackStart(scrolledWindow, true, true, 0)

	buffer, err := textview.GetBuffer()
	if err != nil {
		log.Fatal("Unable to get buffer:", err)
	}
	if file, err := ioutil.ReadFile(historyFilename); err == nil {
		buffer.SetText(string(file))
	}

	entry, err := gtk.EntryNew()
	if err != nil {
		log.Fatal("Unable to create entry:", err)
	}
	entry.Connect("activate", func() {
		message, err := entry.GetText()
		if err != nil {
			log.Fatal("Unable to get text:", err)
		}
		message = strings.TrimSpace(message)
		if message != "" {
			endIter := buffer.GetEndIter()
			buffer.Insert(endIter, username+": "+message+"\n")
			entry.SetText("")
		}
	})
	vbox.PackStart(entry, false, true, 0)

	settingsButton, err := gtk.ButtonNewWithLabel("Settings")
	if err != nil {
		log.Fatal("Unable to create settings button:", err)
	}
	settingsButton.Connect("clicked", func() {
		dialog, err := gtk.DialogNew()
		if err != nil {
			log.Fatal("Unable to create settings dialog:", err)
		}
		dialog.SetTitle("Settings")
		contentArea, err := dialog.GetContentArea()
		if err != nil {
			log.Fatal("Unable to get content area:", err)
		}

		usernameEntry, err := gtk.EntryNew()
		if err != nil {
			log.Fatal("Unable to create username entry:", err)
		}
		usernameEntry.SetText(username)
		contentArea.PackStart(usernameEntry, true, true, 0)

		serverAddressEntry, err := gtk.EntryNew()
		if err != nil {
			log.Fatal("Unable to create server address entry:", err)
		}
		serverAddressEntry.SetText(serverAddress)
		contentArea.PackStart(serverAddressEntry, true, true, 0)

		fontSizeAdjustment, err := gtk.AdjustmentNew(float64(fontSize), 8, 72, 1, 5, 0)
		if err != nil {
			log.Fatal("Unable to create font size adjustment:", err)
		}
		fontSizeSpinButton, err := gtk.SpinButtonNew(fontSizeAdjustment, 1, 0)
		if err != nil {
			log.Fatal("Unable to create font size spin button:", err)
		}
		contentArea.PackStart(fontSizeSpinButton, true, true, 0)

		clearHistoryButton, err := gtk.ButtonNewWithLabel("Clear History")
		if err != nil {
			log.Fatal("Unable to create clear history button:", err)
		}
		clearHistoryButton.Connect("clicked", func() {
			if err := os.Remove(historyFilename); err != nil {
				log.Fatal("Unable to clear history:", err)
			}
			buffer.Delete(buffer.GetStartIter(), buffer.GetEndIter())
		})
		contentArea.PackStart(clearHistoryButton, false, false, 0)

		okButton, err := gtk.ButtonNewWithLabel("OK")
		if err != nil {
			log.Fatal("Unable to create OK button:", err)
		}
		okButton.Connect("clicked", func() {
			newUsername, err := usernameEntry.GetText()
			if err != nil {
				log.Fatal("Unable to get text:", err)
			}
			if newUsername != "" {
				username = newUsername
			}
			newServerAddress, err := serverAddressEntry.GetText()
			if err != nil {
				log.Fatal("Unable to get text:", err)
			}
			if newServerAddress != "" {
				serverAddress = newServerAddress
			}
			newFontSize := int(fontSizeSpinButton.GetValue())
			if newFontSize != fontSize {
				fontSize = newFontSize
				err = cssProvider.LoadFromData("* { font-size: " + strconv.Itoa(fontSize) + "pt; }")
				if err != nil {
					log.Fatal("Unable to load css data:", err)
				}
			}
			dialog.Destroy()
		})
		contentArea.PackStart(okButton, false, false, 0)

		dialog.ShowAll()
	})
	vbox.PackStart(settingsButton, false, true, 0)

	win.Connect("destroy", func() {
		startIter, endIter := buffer.GetBounds()
		text, err := buffer.GetText(startIter, endIter, true)
		if err != nil {
			log.Fatal("Unable to get text:", err)
		}
		err = ioutil.WriteFile(historyFilename, []byte(text), 0644)
		if err != nil {
			log.Fatal("Unable to write file:", err)
		}
		gtk.MainQuit()
	})

	win.Add(vbox)
	win.ShowAll()

	gtk.Main()
}

