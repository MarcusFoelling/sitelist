package main

import (
	"encoding/xml"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	app     = "sitelist"
	version = "1.0.0"
	author  = "Marcus Andreas Fölling <marcus.foelling@businessbytes.de>"
	github  = "https://github.com/marcusfoelling/sitelist"
)

type CreatedBy struct {
	Tool        string `xml:"tool"`
	Version     string `xml:"version"`
	DateCreated string `xml:"date-created"`
}

type Site struct {
	URL        string `xml:"url,attr"`
	CompatMode string `xml:"compat-mode,omitempty"`
	OpenIn     string `xml:"open-in,omitempty"`
}

type SiteList struct {
	XMLName   xml.Name  `xml:"site-list"`
	Version   string    `xml:"version,attr"`
	CreatedBy CreatedBy `xml:"created-by"`
	Sites     []Site    `xml:"site"`
}

func main() {
	fmt.Printf("%s v%s | %s | %s\n\n", app, version, author, github)

	file := flag.String("file", "./sitelist.xml", "Pfad zur XML-Datei")
	add := flag.String("add", "", "Domain zur Sitelist hinzufügen")
	delete := flag.String("delete", "", "Domain aus der Sitelist entfernen")
	flag.Parse()

	if *add != "" && *delete != "" {
		exitError(errors.New("add und delete dürfen nicht gleichzeitig verwendet werden"))
	}
	if *add == "" && *delete == "" {
		flag.Usage()
		exitError(errors.New("entweder -add oder -delete muss angegeben werden"))
	}

	sitelist, err := loadSiteList(*file)
	if err != nil {
		exitError(err)
	}

	switch {
	case *add != "":
		domain := normalizeDomain(*add)
		if err := validateDomain(domain); err != nil {
			exitError(err)
		}
		for _, site := range sitelist.Sites {
			if site.URL == domain {
				fmt.Printf("Domain bereits vorhanden: %s\n", domain)
				return
			}
		}
		site := Site{
			URL:        domain,
			CompatMode: "IE11",
			OpenIn:     "IE11",
		}
		sitelist.Sites = append(sitelist.Sites, site)
		fmt.Printf("Domain hinzugefügt: %s\n", domain)

	case *delete != "":
		domain := normalizeDomain(*delete)
		if err := validateDomain(domain); err != nil {
			exitError(err)
		}
		updatedSite := make([]Site, 0, len(sitelist.Sites))
		found := false
		for _, site := range sitelist.Sites {
			if site.URL == domain {
				found = true
				continue
			}
			updatedSite = append(updatedSite, site)
		}
		if !found {
			fmt.Printf("Domain nicht gefunden: %s\n", domain)
			return
		}
		sitelist.Sites = updatedSite
		fmt.Printf("Domain entfernt: %s\n", domain)
	}

	if err := saveSiteList(*file, sitelist); err != nil {
		exitError(err)
	}
}

func loadSiteList(filename string) (SiteList, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return SiteList{
				Version: "2",
				CreatedBy: CreatedBy{
					Tool:        filepath.Base(os.Args[0]),
					Version:     "0",
					DateCreated: time.Now().Format(time.RFC3339)},
			}, nil
		}
		return SiteList{}, fmt.Errorf("XML-Datei konnte nicht gelesen werden: %w", err)
	}

	var sitelist SiteList
	if err := xml.Unmarshal(data, &sitelist); err != nil {
		return SiteList{}, fmt.Errorf("XML-Datei ist ungültig: %w", err)
	}
	return sitelist, nil
}

func saveSiteList(filename string, sitelist SiteList) error {
	newVersion, err := incrementVersion(sitelist.CreatedBy.Version)
	if err != nil {
		newVersion = "1"
	}
	sitelist.CreatedBy.Version = newVersion
	sitelist.CreatedBy.Tool = filepath.Base(os.Args[0])
	sitelist.CreatedBy.DateCreated = time.Now().Format(time.RFC3339)
	sort.Slice(sitelist.Sites, func(i, j int) bool {
		return sitelist.Sites[i].URL < sitelist.Sites[j].URL
	})

	data, err := xml.MarshalIndent(sitelist, "", "  ")
	if err != nil {
		return fmt.Errorf("XML konnte nicht erzeugt werden: %w", err)
	}
	data = append([]byte(xml.Header), data...)

	dir := filepath.Dir(filename)
	tempFile, err := os.CreateTemp(dir, ".sitelist-*.tmp")
	if err != nil {
		return fmt.Errorf("temporäre Datei konnte nicht erstellt werden: %w", err)
	}
	tempName := tempFile.Name()

	defer func() {
		_ = tempFile.Close()
		_ = os.Remove(tempName)
	}()

	if _, err := tempFile.Write(data); err != nil {
		return fmt.Errorf("temporäre Datei konnte nicht geschrieben werden: %w", err)
	}
	if err := tempFile.Close(); err != nil {
		return fmt.Errorf("temporäre Datei konnte nicht geschlossen werden: %w", err)
	}
	if err := os.Rename(tempName, filename); err != nil {
		return fmt.Errorf("XML-Datei konnte nicht ersetzt werden: %w", err)
	}
	return nil
}

func incrementVersion(version string) (string, error) {
	v, err := strconv.Atoi(version)
	if err != nil {
		return "", fmt.Errorf("ungültige Version %q: %w", version, err)
	}
	v++
	return strconv.Itoa(v), nil
}

func normalizeDomain(domain string) string {
	domain = strings.TrimSpace(domain)
	domain = strings.ToLower(domain)
	domain = strings.TrimSuffix(domain, ".")
	return domain
}

func validateDomain(domain string) error {
	if domain == "" {
		return errors.New("Domain darf nicht leer sein")
	}
	if strings.ContainsAny(domain, " /\\@:") {
		return fmt.Errorf("ungültiger Domainname: %q", domain)
	}
	if !strings.Contains(domain, ".") {
		return fmt.Errorf("ungültiger Domainname: %q", domain)
	}
	return nil
}

func exitError(err error) {
	fmt.Fprintln(os.Stderr, "Fehler:", err)
	os.Exit(1)
}
