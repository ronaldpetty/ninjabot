package main

import (
	"log"
	"os"

	"github.com/rodrigo-brito/ninjabot/pkg/data"
	"github.com/rodrigo-brito/ninjabot/pkg/exchange"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "download",
		Usage: "download historical data",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "pair",
				Aliases:  []string{"p"},
				Usage:    "eg. BTCUSDT",
				Required: true,
			},
			&cli.IntFlag{
				Name:     "days",
				Aliases:  []string{"d"},
				Usage:    "eg. 100 (default 30 days)",
				Required: false,
			},
			&cli.TimestampFlag{
				Name:     "start",
				Aliases:  []string{"s"},
				Usage:    "eg. 2021-12-01",
				Layout:   "2006-01-02",
				Required: false,
			},
			&cli.TimestampFlag{
				Name:     "end",
				Aliases:  []string{"e"},
				Usage:    "eg. 2020-12-31",
				Layout:   "2006-01-02",
				Required: false,
			},
			&cli.StringFlag{
				Name:     "timeframe",
				Aliases:  []string{"t"},
				Usage:    "eg. 1h",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "output",
				Aliases:  []string{"o"},
				Usage:    "eg. ./btc.csv",
				Required: true,
			},
		},
		Action: func(c *cli.Context) error {
			exc, err := exchange.NewBinance(c.Context)
			if err != nil {
				return err
			}

			var options []data.Option
			if days := c.Int("days"); days > 0 {
				options = append(options, data.WithDays(days))
			}

			start := c.Timestamp("start")
			end := c.Timestamp("end")
			if start != nil && end != nil && !start.IsZero() && !end.IsZero() {
				options = append(options, data.WithInterval(*start, *end))
			} else if start != nil || end != nil {
				log.Fatal("START and END must be informed together")
			}

			return data.NewDownloader(exc).Download(c.Context, c.String("pair"),
				c.String("timeframe"), c.String("output"), options...)
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
