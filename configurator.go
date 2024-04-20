/*
**
** IPv6rs Configurator
** Provides a GUI menu interface dynamically created from a JSON.
**
** Distributed under the COOL License.
**
** Copyright (c) 2024 IPv6.rs <https://ipv6.rs>
** All Rights Reserved
**
*/

package configurator

import (
  "encoding/json"
  "fmt"
  "fyne.io/fyne/v2"
  "fyne.io/fyne/v2/canvas"
  "fyne.io/fyne/v2/container"
  "fyne.io/fyne/v2/dialog"
  "fyne.io/fyne/v2/layout"
  "fyne.io/fyne/v2/storage"
  "fyne.io/fyne/v2/widget"
  "fyne.io/fyne/v2/theme"
  "image/color"
  "io/ioutil"
  "net/url"
  "os"
  "os/exec"
  "runtime"
  "strconv"
  "strings"
  "regexp"
)


type Configuration struct {
  Header            string   `json:"Header"`
  Entries           []Entry  `json:"Entries"`
  SubmitButtonText  string   `json:"SubmitButtonText"`
  Exec              string   `json:"exec"`
}

type Entry struct {
  VariableName      string   `json:"VariableName"`
  Question          string   `json:"Question"`
  CanBeNull         bool     `json:"CanBeNull"`
  AcceptableAnswers []string `json:"AcceptableAnswers"`
  Placeholder       string   `json:"Placeholder"`
  Type              string   `json:"Type"`
  AllowedChars      string   `json:"AllowedChars,omitempty"`
  MinLength         *int      `json:"MinLength,omitempty"`
  MaxLength         *int      `json:"MaxLength,omitempty"`
}

func NewWindow(app fyne.App, arg string, linkText string, linkUrl string) {
  if(arg == "") {
    os.Exit(1)
  }

  url, err := url.Parse(linkUrl)
  if err != nil {
    panic(err)
  }

  config, err := loadConfiguration(arg)
  if err != nil {
    fmt.Println("Failed to load configuration:", err)
    os.Exit(1)
  }

  w := app.NewWindow(config.Header)
  w.Resize(fyne.NewSize(500, 750))

  inputs := make(map[string]fyne.CanvasObject)
  errors := make(map[string]*canvas.Text)
  errorContainers := make(map[string]*fyne.Container)

  verticalSpacer := canvas.NewRectangle(color.Transparent)
  verticalSpacer.SetMinSize(fyne.NewSize(0, 40))

  smallVerticalSpacer := canvas.NewRectangle(color.Transparent)
  smallVerticalSpacer.SetMinSize(fyne.NewSize(0, 10))

  formItems := container.NewVBox(widget.NewLabelWithStyle(config.Header, fyne.TextAlignCenter, fyne.TextStyle{Bold: true}))
  formItems.Add(verticalSpacer)

  for _, entry := range config.Entries {
    label := widget.NewLabel(entry.Question)
    var errorContainer *fyne.Container = container.NewHBox()
    errorContainers[entry.VariableName] = errorContainer
    errorText := canvas.NewText("", theme.ErrorColor())
    errorText.Hide()
    errors[entry.VariableName] = errorText

    if entry.Placeholder == "FILEPICKER" {
      filePickerEntry := widget.NewEntry()
      filePickerEntry.Disable()
      filePickerButton := makeFilePickerButton(w, filePickerEntry, entry.AcceptableAnswers)
      inputs[entry.VariableName] = filePickerEntry
      formItems.Add(label)
      formItems.Add(container.NewHBox(filePickerEntry, filePickerButton))
      formItems.Add(errorContainer)
    } else if entry.Placeholder == "FOLDERPICKER" {
      folderPickerEntry := widget.NewEntry()
      folderPickerEntry.Disable()
      folderPickerButton := makeFolderPickerButton(w, folderPickerEntry)
      inputs[entry.VariableName] = folderPickerEntry
      formItems.Add(label)
      formItems.Add(container.NewHBox(folderPickerEntry, folderPickerButton))
      formItems.Add(errorContainer)
    } else {
      var input fyne.CanvasObject
      if len(entry.AcceptableAnswers) > 0 {
        selectInput := widget.NewSelect(entry.AcceptableAnswers, nil)
        selectInput.PlaceHolder = entry.Placeholder
        input = selectInput
      } else if entry.Type=="password" {
        input = widget.NewPasswordEntry();
        input.(*widget.Entry).SetPlaceHolder(entry.Placeholder)
      } else {
        input = widget.NewEntry()
        input.(*widget.Entry).SetPlaceHolder(entry.Placeholder)
      }
      inputs[entry.VariableName] = input
      formItems.Add(label)
      formItems.Add(input)
      formItems.Add(errorContainer)
      formItems.Add(verticalSpacer)
    }
  }
  submitButton := NewResizableButton(config.SubmitButtonText, 140, func() {
    allValid, envVars := validateFormEntries(config.Entries, inputs, errors, errorContainers)
    if allValid {
      executeCommand(config.Exec, envVars)
      w.Close()
    }
  })
  submitButtonContainer := container.NewHBox(layout.NewSpacer(), submitButton)
  formItems.Add(submitButtonContainer)

  ipv6rsLink := widget.NewHyperlinkWithStyle(linkText, url, fyne.TextAlignTrailing, fyne.TextStyle{})
  formItems.Add(smallVerticalSpacer)
  formItems.Add(ipv6rsLink)
  formItems.Add(verticalSpacer)


  paddingSize := fyne.NewSize(50, 50)
  paddedContainer := container.New(NewResizablePaddedLayout(paddingSize), formItems)

  scrollContainer := container.NewVScroll(paddedContainer)

  w.SetContent(scrollContainer)

  w.Show()
}
func makeFolderPickerButton(w fyne.Window, entry *widget.Entry) *widget.Button {
  return widget.NewButton("Browse", func() {
    dialog.ShowFolderOpen(func(uri fyne.ListableURI, err error) {
      if uri != nil && err == nil {
        entry.SetText(uri.Path())
        entry.Refresh()
      }
    }, w)
  })
}
func makeFilePickerButton(window fyne.Window, entry *widget.Entry, filters []string) *widget.Button {
  return widget.NewButton("Select File", func() {
    fileDialog := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
      if err != nil || reader == nil {
        return
      }
      entry.SetText(reader.URI().Path())
    }, window)

    if len(filters) > 0 {
      extFilters := make([]string, len(filters))
      for i, ext := range filters {
        extFilters[i] = "." + ext
      }
      fileDialog.SetFilter(storage.NewExtensionFileFilter(extFilters))
    }
    fileDialog.Show()
  })
}
func validateFormEntries(entries []Entry, inputs map[string]fyne.CanvasObject, errors map[string]*canvas.Text, errorContainers map[string]*fyne.Container) (bool, map[string]string) {
  allValid := true
  envVars := make(map[string]string)
  for _, entry := range entries {
    inputText := ""
    if entry.Placeholder == "FILEPICKER" || entry.Placeholder == "FOLDERPICKER" {
      inputText = inputs[entry.VariableName].(*widget.Entry).Text
    } else {
      switch input := inputs[entry.VariableName].(type) {
      case *widget.Entry:
        inputText = input.Text
      case *widget.Select:
        inputText = input.Selected
      default:
        fmt.Println("Unsupported input type for variable:", entry.VariableName)
        allValid = false
        continue
      }
    }

    errorLabel := errors[entry.VariableName]
    errorContainer := errorContainers[entry.VariableName]

    errorContainer.Objects = []fyne.CanvasObject{}
    if !entry.CanBeNull && inputText == "" {
      errorLabel.Text = "This field cannot be empty"
      allValid = false
    } else if entry.MinLength != nil && len(inputText) < *entry.MinLength {
      errorLabel.Text = fmt.Sprintf("Minimum length of %d characters is required", *entry.MinLength)
      allValid = false
    } else if entry.MaxLength != nil && len(inputText) > *entry.MaxLength {
      errorLabel.Text = fmt.Sprintf("Maximum length of %d characters is required", *entry.MaxLength)
      allValid = false
    } else if !isValidType(inputText, entry.Type) {
      errorLabel.Text = fmt.Sprintf("Expecting a %s", entry.Type)
      allValid = false
    } else if entry.AllowedChars != "" && !regexp.MustCompile(entry.AllowedChars).MatchString(inputText) {
      errorLabel.Text = "Input contains invalid characters"
      allValid = false
    } else {
      errorLabel.Text = ""
    }

    if errorLabel.Text != "" {
      errorContainer.Add(errorLabel)
      errorLabel.Show()
    } else {
      errorLabel.Hide()
    }
    errorContainer.Refresh()
    if allValid {
      envVars[entry.VariableName] = inputText
    }
  }
  return allValid, envVars
}

func isValidType(input, expectedType string) bool {
  switch expectedType {
    case "int":
      _, err := strconv.Atoi(input)
      return err == nil
    case "float":
      _, err := strconv.ParseFloat(input, 64)
      return err == nil
    case "string":
      return true
    case "password":
      return true
    default:
      return false
  }
}

func loadConfiguration(file string) (*Configuration, error) {
  var config Configuration

  file = expandPath(file)

  data, err := ioutil.ReadFile(file)
  if err != nil {
    return nil, err
  }
  err = json.Unmarshal(data, &config)
  if err != nil {
    return nil, err
  }
  return &config, nil
}

func executeCommand(scriptFilename string, envVars map[string]string) error {
  scriptFilename = expandPath(scriptFilename)


  var scriptFile string;
  var extension string;

  switch runtime.GOOS {
    case "windows":
      scriptFile = strings.TrimSuffix(scriptFilename, ".sh") + ".ps1"
      extension = "ps1"
    default:
      scriptFile = scriptFilename;
      extension = "sh"
  }

  scriptContent, err := ioutil.ReadFile(scriptFile)
  if err != nil {
    fmt.Println("Error reading script file:", err)
    return err
  }
  modifiedScript := string(scriptContent)

  for key, value := range envVars {
    placeholder := fmt.Sprintf("_%s", key)
    modifiedScript = strings.ReplaceAll(modifiedScript, placeholder, value)
  }

  tmpfile, err := ioutil.TempFile("", "script_*."+extension)
  if err != nil {
    fmt.Println("Error creating temporary script file:", err)
    return err
  }

  if _, err := tmpfile.Write([]byte(modifiedScript)); err != nil {
    fmt.Println("Error writing to temporary script file:", err)
    return err
  }
  if err := tmpfile.Close(); err != nil {
    fmt.Println("Error closing temporary script file:", err)
    return err
  }

  if err := os.Chmod(tmpfile.Name(), 0755); err != nil {
    fmt.Println("Error setting script file executable:", err)
    return err
  }

  if err := executeInTerminal(tmpfile.Name()); err != nil {
    return err
  }
  return nil
}

func executeInTerminal(scriptFilename string) error {
  var cmd *exec.Cmd

  switch runtime.GOOS {
    case "windows":
      escapedScriptFilename := strings.ReplaceAll(scriptFilename, `\`, `\\`)
      psCommand := fmt.Sprintf("powershell -ExecutionPolicy Bypass -File \"%s\"; Remove-Item -Path \"%s\"", escapedScriptFilename, escapedScriptFilename)
      cmd = exec.Command("cmd", "/C", "start", "PowerShell Script Window", "powershell", "-NoExit", "-Command", psCommand)
    case "darwin":
      cmd = exec.Command("osascript", "-e", fmt.Sprintf(`tell application "Terminal" to do script "sh %s && rm %s"`, scriptFilename, scriptFilename))
    case "linux":
      cmd = exec.Command("gnome-terminal", "--", "bash", "-c", fmt.Sprintf(`%s; rm %s; exec bash`, scriptFilename, scriptFilename))

    default:
      fmt.Errorf("unsupported platform")
  }

  cmd.Stdout = os.Stdout
  cmd.Stderr = os.Stderr
  err := cmd.Start()
  if err != nil {
    fmt.Errorf("Failed to execute command: %v", err)
    return err;
  }

  return nil;
}

func expandPath(path string) string {
  if path == "~" || path == "$HOME" {
    homeDir, err := os.UserHomeDir()
    if err != nil {
      fmt.Printf("Unable to find user's home directory: %v\n", err)
      return path
    }
    return homeDir
  }
  if strings.HasPrefix(path, "~/") || strings.HasPrefix(path, "$HOME/") {
    homeDir, err := os.UserHomeDir()
    if err != nil {
      fmt.Println("Unable to determine the home directory:", err)
      return path
    }
    var newPath string
    if runtime.GOOS == "windows" && (strings.HasPrefix(path, "~/.") || strings.HasPrefix(path, "$HOME/.")) {
      localAppData, exists := os.LookupEnv("LOCALAPPDATA")
      if !exists {
        fmt.Println("LOCALAPPDATA environment variable not set.")
        return path
      }
      newPath = strings.Replace(path, "~/.", localAppData+string(os.PathSeparator), 1)
      newPath = strings.Replace(newPath, "$HOME/.", localAppData+string(os.PathSeparator), 1)
    } else {
      newPath = strings.Replace(path, "~/", homeDir+string(os.PathSeparator), 1)
      newPath = strings.Replace(newPath, "$HOME/", homeDir+string(os.PathSeparator), 1)
    }
    return newPath
  }
  return path
}
