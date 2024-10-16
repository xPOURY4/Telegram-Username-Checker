# Telegram Username Checker

This tool allows you to check the availability of Telegram usernames using the Telegram API. It handles rate limiting and various response types from the Telegram servers.

## Features

- Check multiple usernames in bulk
- Handle rate limiting with adaptive delay
- Categorize usernames as available, taken, invalid, or purchasable
- Save results to separate files based on status

## Prerequisites

- Go 1.16 or higher
- Telegram API credentials (API ID and API Hash)

## Installation

1. Clone this repository:
```
git clone https://github.com/xPOURY4/telegram-username-checker.git
cd telegram-username-checker
```

2. Install dependencies:
```
go mod tidy
```

## Configuration

1. Open `main.go` and replace the following constants with your Telegram API credentials:
```go
const (
    apiID   = "YOUR_API_ID"
    apiHash = "YOUR_API_HASH"
    phone   = "YOUR_PHONE_NUMBER"
   )
```

   You can obtain your API ID and API Hash from the [Telegram API Development Tools](https://my.telegram.org/apps).

2. Create a file named `usernames.txt` in the same directory as `main.go`. Add the usernames you want to check, one per line.

## Usage

1. Run the program:
```
go run main.go
```

2. When prompted, enter the authentication code sent to your Telegram account.

3. The program will start checking the usernames listed in `usernames.txt`.

4. Results will be displayed in the console and saved to separate files:
   - `available_usernames.txt`: Usernames that are available for registration
   - `taken_usernames.txt`: Usernames that are already in use
   - `invalid_usernames.txt`: Usernames that are invalid according to Telegram's rules
   - `purchasable_usernames.txt`: Usernames that are available for purchase
   - `error_usernames.txt`: Usernames that couldn't be checked due to errors

## Notes

- The tool implements an adaptive delay system to handle rate limiting. If you encounter `FLOOD_WAIT` errors, the program will automatically wait and retry.
- Be cautious when checking a large number of usernames, as it may trigger Telegram's anti-spam measures.
- This tool is for educational purposes only. Make sure to comply with Telegram's Terms of Service when using it.

## Troubleshooting

- If you encounter persistent `FLOOD_WAIT` errors, try increasing the `baseDelay` in the `main()` function.
- Make sure your `apiID`, `apiHash`, and `phone` are correctly set in the code.
- Ensure that your Telegram account is not limited or banned.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Disclaimer

This tool is not officially associated with Telegram. Use it at your own risk.
