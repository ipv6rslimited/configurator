# Configurator for golang

Configurator provides a GUI interface for your application, dynamically generated from a JSON file.

## Screenshot

![Configurator](https://raw.githubusercontent.com/ipv6rslimited/configurator/main/screenshot.png)

This was generated from:

```
{
  "Header": "Title",
  "Entries": [
    {
      "VariableName": "NAME",
      "Question": "What's your name?",
      "CanBeNull": false,
      "AcceptableAnswers": [],
      "Type":"string",
      "Placeholder": "Enter name"
    },
    {
      "VariableName": "PASS",
      "Question": "What's your password?",
      "CanBeNull": false,
      "AcceptableAnswers": [],
      "Type":"password",
      "Placeholder": "Enter pass"
    },
    {
      "VariableName": "COOL",
      "Question": "Are you cool?",
      "CanBeNull": false,
      "Type":"string",
      "AcceptableAnswers": ["Y", "N"],
      "Placeholder": "Choose Yes or No"
    },
    {
      "VariableName": "FILE",
      "Question": "Filename:",
      "CanBeNull": false,
      "Type":"string",
      "AcceptableAnswers": ["json"],
      "Placeholder": "FILEPICKER"
    }
  ],
  "SubmitButtonText": "Create",
  "exec": "example.sh"
}
```

It will call example.sh, wherein, the variables will be automatically filled in and usable.

```
#!/bin/bash
echo _NAME _PASS _COOL _FILE
```


## Features

- Custom GUI Dialog Generation from JSON

- Input Type Checking

- Supports drop down selection, text, password and files

- Runs a script with dynamic var replacement prior

- Checks for null, regex and min/max length

## Use Case

- We needed a dialog for our 1-click self hosting installer for [IPv6rs](https://ipv6.rs)

## Example

- See the configurator-test folder to see how to use Configurator

- Build the configurator-test by typing:
```
cd test
go run test.go example.json
```

## License

Distributed under the COOL License.

Copyright (c) 2024 IPv6.rs <https://ipv6.rs>
All Rights Reserved
