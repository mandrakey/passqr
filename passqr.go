package main

import (
    "fmt"
    "os"
    "errors"
    "os/exec"
    "strings"

    "github.com/urfave/cli"
    qrcode "github.com/skip2/go-qrcode"
)

func main() {
    app := cli.NewApp()
    app.Name = "passqr"
    app.Usage = "create a qr code from pass output"
    app.Flags = []cli.Flag {
        cli.StringFlag{
            Name: "filename,f",
            Value: "-",
            Usage: "Save the qr code to value. Pass - to display on console.",
        },
        cli.IntFlag{
            Name: "size,s",
            Value: 400,
            Usage: "Set the output image size if written to image file.",
        },
    }
    app.Action = run
    app.Run(os.Args)
}

func handleError(ctx *cli.Context, err error) {
    fmt.Printf("ERROR: %s\n", err)
}

func run(ctx *cli.Context) error {
    args := ctx.Args()
    if len(args) == 0 {
        return errors.New("missing path to password")
    }

    path := args[0]
    filename := ctx.String("filename")

    // Get the password
    cmd := exec.Command("pass", path)
    password, err := cmd.Output(); if err != nil {
        return errors.New(fmt.Sprintf("Retrieving the password failed: %s", err))
    }

    // Create qr code
    q, err := qrcode.New(string(password), qrcode.Medium); if err != nil {
        return errors.New(fmt.Sprintf("Failed to generate QR code: %s", err))
    }

    // Display or save
    tmp := strings.Split(filename, ".")
    ext := tmp[len(tmp) - 1]
    if contains([]string{"-","png"}, &ext) == false {
        fmt.Printf("Invalid file extension %s. Falling back to 'display on console'.", ext)
        ext = "-"
    }

    switch ext {
    case "-":
        displayQrCode(q)
    case "png":
        q.WriteFile(ctx.Int("size"), filename)
    }
    return nil
}

func displayQrCode(q *qrcode.QRCode) {
    bmp := q.Bitmap()
    for i := range(bmp) {
        for j := range(bmp[i]) {
            var c string
            if bmp[i][j] == true {
                c = "\033[40m  "
            } else {
                c = "\033[47m  "
            }
            fmt.Printf("%s", c)
        }
        fmt.Printf("\033[40m\n")
    }
    fmt.Printf("\033[40m")
}

func contains(s []string, v *string) bool {
    for _, i := range(s) {
        if i == *v {
            return true
        }
    }
    return false
}
