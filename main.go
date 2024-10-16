package main

import (
    "context"
    "fmt"
    "os"
    "strings"
    "time"

    "github.com/gotd/td/telegram"
    "github.com/gotd/td/telegram/auth"
    "github.com/gotd/td/tg"
)

const (
    apiID   = "YOUR_API_ID"
    apiHash = "YOUR_API_HASH"
    phone   = "YOUR_PHONE_NUMBER"
)

func displayHeader() {
    header := `
╔════════════════════════════════════════════╗
║        Telegram Username Checker           ║
║         Created by: xPOURY4                ║
║     GitHub: github.com/xPOURY4             ║
╚════════════════════════════════════════════╝
`
    fmt.Println(header)
}

func main() {
    displayHeader()

    ctx := context.Background()
    client := telegram.NewClient(apiID, apiHash, telegram.Options{})

    if err := client.Run(ctx, func(ctx context.Context) error {
        flow := auth.NewFlow(
            auth.Constant(phone, "", auth.CodeAuthenticatorFunc(func(ctx context.Context, _ *tg.AuthSentCode) (string, error) {
                fmt.Print("Enter the code you received: ")
                var code string
                _, err := fmt.Scan(&code)
                return code, err
            })),
            auth.SendCodeOptions{},
        )

        if err := flow.Run(ctx, client.Auth()); err != nil {
            return fmt.Errorf("auth flow: %w", err)
        }

        fmt.Println("Successfully logged in!")

        usernames, err := readUsernamesFromFile("usernames.txt")
        if err != nil {
            return fmt.Errorf("reading usernames: %w", err)
        }

        results := make(map[string][]string)
        baseDelay := time.Second * 3

        for _, username := range usernames {
            status, err := checkUsername(ctx, client.API(), username)
            if err != nil {
                fmt.Printf("Error checking %s: %v\n", username, err)
                results["error"] = append(results["error"], username)
                if strings.Contains(err.Error(), "FLOOD_WAIT") {
                    waitTime, _ := time.ParseDuration(strings.Split(err.Error(), " ")[4] + "s")
                    fmt.Printf("Rate limit hit. Waiting for %v\n", waitTime)
                    time.Sleep(waitTime)
                    baseDelay *= 2 // Increase delay for subsequent requests
                }
            } else {
                fmt.Printf("%s is %s\n", username, status)
                results[status] = append(results[status], username)
            }

            time.Sleep(baseDelay)
        }

        for status, usernames := range results {
            if len(usernames) > 0 {
                if err := saveToFile(fmt.Sprintf("%s_usernames.txt", status), usernames); err != nil {
                    fmt.Printf("Error saving %s usernames: %v\n", status, err)
                } else {
                    fmt.Printf("%s usernames have been saved to %s_usernames.txt\n", strings.Title(status), status)
                }
            }
        }

        return nil
    }); err != nil {
        fmt.Printf("Error: %v\n", err)
        os.Exit(1)
    }
}

func checkUsername(ctx context.Context, api *tg.Client, username string) (string, error) {
    _, err := api.AccountUpdateUsername(ctx, username)
    if err != nil {
        if strings.Contains(err.Error(), "USERNAME_OCCUPIED") {
            return "taken", nil
        }
        if strings.Contains(err.Error(), "USERNAME_INVALID") {
            return "invalid", nil
        }
        if strings.Contains(err.Error(), "USERNAME_PURCHASE_AVAILABLE") {
            return "purchasable", nil
        }
        return "error", err
    }
    _, err = api.AccountUpdateUsername(ctx, "") // Reset username
    if err != nil {
        return "error", err
    }
    return "available", nil
}

func readUsernamesFromFile(filename string) ([]string, error) {
    content, err := os.ReadFile(filename)
    if err != nil {
        return nil, err
    }
    lines := strings.Split(string(content), "\n")
    var usernames []string
    for _, line := range lines {
        if username := strings.TrimSpace(line); username != "" {
            usernames = append(usernames, username)
        }
    }
    return usernames, nil
}

func saveToFile(filename string, lines []string) error {
    content := strings.Join(lines, "\n")
    return os.WriteFile(filename, []byte(content), 0644)
}
