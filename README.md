## TypeEase

TypeEase is a straightforward tool designed to help you manage custom keyboard shortcuts and automate text input. It allows you to bind specific key combinations to predefined text, making it easier to handle repetitive typing tasks and improve your productivity.

## Features

- **Custom Key Bindings**: Assign custom text to specific keyboard shortcuts.
- **Real-time Key Listening**: Monitors your keyboard inputs and triggers text replacement when a bound shortcut is typed.
- **Simple GUI**: A user-friendly interface built with Fyne, making it easy to create, edit, and delete key bindings.
- **File Storage**: Save your key bindings to a file for persistence across sessions.

## Installation

### Prerequisites

- Go 1.20 or later installed on your machine.
- A Linux-based system (required for `/dev/input` access).
- Fyne (a Go-based GUI toolkit).
- RobotGo (for simulating keyboard inputs).

### Clone the Repository

```bash
git clone https://github.com/medalha01/typeease.git
cd typeease
```

### Install Dependencies

```bash
go mod tidy
```

### Build the Project

```bash
go build -o typeease ./cmd/keymanager
```

This will create an executable named `typeease`.

### Run the Application

```bash
sudo ./typeease
```

## Usage

### Creating a New Key Binding

1. Open the TypeEase application.
2. In the "Key Details" section, enter the key combination you want to bind in the "Activation Key" field.
3. Enter the text you want to be inserted when this key combination is pressed in the "Substitute text" field.
4. Toggle the "Active" checkbox if you want this binding to be active immediately.
5. Click the "Save" button to store the binding.

### Editing an Existing Binding

1. Select the binding you want to edit from the list.
2. Modify the "Activation Key" or "Substitute text" fields as needed.
3. Click "Save" to update the binding.

### Deleting a Binding

1. Select the binding you want to delete from the list.
2. Click the "Delete" button.
3. Confirm the deletion when prompted.

### Persistent Storage

All key bindings are saved to a file in the application's directory, allowing them to persist across sessions.

## Limitations

- **Linux-only**: The application currently relies on `/dev/input` devices, which are specific to Linux systems.
- **Modifier Keys**: Some modifier keys (e.g., SHIFT, CTRL, ALT) are ignored in binding detection to avoid unintended triggers.
- **Basic Text Replacement**: The application is designed for simple text replacement and does not support complex scripting or macros.

## Contributing

Contributions are welcome! If you have suggestions, feel free to open an issue or submit a pull request.

### How to Contribute

1. Fork the repository.
2. Create a new branch for your feature or bugfix.
3. Commit your changes with clear descriptions.
4. Push to your branch and submit a pull request.

## Acknowledgments

- [Fyne](https://fyne.io/) for the GUI framework.
- [RobotGo](https://github.com/go-vgo/robotgo) for the keyboard input automation.

---

TypeEase is designed to be a simple, no-frills solution for managing keyboard shortcuts and text automation. It's ideal for users who want to streamline their typing tasks with minimal setup.
