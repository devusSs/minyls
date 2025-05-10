package cli

import (
	"fmt"
	"net/url"
	"os"
	"path"
	"text/tabwriter"
	"time"

	"github.com/devusSs/minyls/internal/log"
	"github.com/devusSs/minyls/internal/storage"
)

func List() error {
	err := initialize()
	if err != nil {
		return fmt.Errorf("failed to initialize: %w", err)
	}

	log.Log().Debug().Str("func", "cli.List").Msg("initialized")

	data, err := storage.Read()
	if err != nil {
		return fmt.Errorf("failed to read storage: %w", err)
	}

	log.Log().Debug().Str("func", "cli.List").Any("data", data).Msg("read data from storage")

	if len(data.Entries) == 0 {
		fmt.Println("NO DATA TO BE DISPLAYED")
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
	fmt.Fprintln(w, "ID\tTimestamp\tMinio ID\tYOURLS ID\tExpiry")

	for _, entry := range data.Entries {
		minioURL, err := url.Parse(entry.MinioLink)
		if err != nil {
			return fmt.Errorf("failed to parse minio url: %w", err)
		}

		yourlsURL, err := url.Parse(entry.YOURLSLink)
		if err != nil {
			return fmt.Errorf("failed to parse yourls url: %w", err)
		}

		fmt.Fprintf(w,
			"%d\t%s\t%s\t%s\t%s\n",
			entry.ID,
			entry.Timestamp.Format(time.DateTime),
			path.Base(minioURL.Path),
			path.Base(yourlsURL.Path),
			entry.Expiry,
		)
	}

	return w.Flush()
}
