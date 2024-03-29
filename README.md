# Configurator for golang

Configurator provides a GUI interface for your application, dynamically generated from a JSON file.

## Screenshot

![Configurator](https://raw.githubusercontent.com/ipv6rslimited/configurator/main/screenshot.png)

## Features

- Custom GUI Dialog Generation from JSON

- Input Type Checking

- Supports drop down selection, text, password and files

- Runs a script with dynamic var replacement prior

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
