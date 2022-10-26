package sitemap

import (
	"bytes"
	"fmt"
	"time"
)

type Sitemap struct {
	Xmlns        string   `xml:"xmlns,attr"`
	Urls         []*Url   `xml:"url"`
	SitemapIndex []*Index `xml:"sitemap"`
}

func NewSitemap() *Sitemap {
	return &Sitemap{Xmlns: "http://www.sitemaps.org/schemas/sitemap/0.9"}
}

func (s *Sitemap) String() string {
	buf := bytes.NewBufferString("")
	buf.WriteString("<?xml version=\"1.0\" encoding=\"UTF-8\"?>")
	switch {
	case len(s.Urls) > 0:
		buf.WriteString(s.toUrlSet())
	case len(s.SitemapIndex) > 0:
		buf.WriteString(s.toSitemapSet())
	}
	return buf.String()
}

func (s *Sitemap) toUrlSet() string {
	buf := bytes.NewBufferString("")
	buf.WriteString(fmt.Sprintf("<urlset xmlns=\"%s\">", s.Xmlns))
	for _, url := range s.Urls {
		buf.WriteString("<url>")
		buf.WriteString("<loc>")
		buf.WriteString(url.Loc)
		buf.WriteString("</loc>")
		if !url.Lastmod.IsZero() {
			buf.WriteString("<lastmod>")
			buf.WriteString(url.Lastmod.Format("2006-01-02"))
			buf.WriteString("</lastmod>")
		}
		if url.Changefreq != "" {
			buf.WriteString("<changefreq>")
			buf.WriteString(string(url.Changefreq))
			buf.WriteString("</changefreq>")
		}
		if url.Changefreq != "" {
			buf.WriteString("<priority>")
			buf.WriteString(fmt.Sprintf("%.1f", url.Priority))
			buf.WriteString("</priority>")
		}
		buf.WriteString("</url>")
	}
	buf.WriteString("</urlset>")
	return buf.String()
}
func (s *Sitemap) toSitemapSet() string {
	buf := bytes.NewBufferString("")
	buf.WriteString(fmt.Sprintf("<sitemapindex xmlns=\"%s\">", s.Xmlns))
	for _, url := range s.Urls {
		buf.WriteString("<sitemap>")
		buf.WriteString("<loc>")
		buf.WriteString(url.Loc)
		buf.WriteString("</loc>")
		buf.WriteString("</sitemap>")
	}
	buf.WriteString("</sitemapindex>")
	return buf.String()
}

type Changefreq string

const (
	ChangefreqDaily  = "daily"
	ChangefreqWeekly = "weekly"
)

type Url struct {
	Loc        string     `xml:"loc"`
	Lastmod    time.Time  `xml:"lastmod"`
	Changefreq Changefreq `xml:"changefreq"`
	Priority   float32    `xml:"priority"`
}

func NewUrl(loc string, lastmod time.Time, changefreq Changefreq, priority float32) *Url {
	return &Url{Loc: loc, Lastmod: lastmod, Changefreq: changefreq, Priority: priority}
}

type Index struct {
	Loc string `xml:"loc"`
}

func NewIndex(loc string) *Index {
	return &Index{Loc: loc}
}
