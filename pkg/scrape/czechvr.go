package scrape

import (
	"log"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
	"github.com/mozillazg/go-slugify"
	"github.com/nleeper/goment"
	"github.com/thoas/go-funk"
)

func ScrapeCzechVR(knownScenes []string, out *[]ScrapedScene) error {
	siteCollector := colly.NewCollector(
		colly.AllowedDomains("www.czechvrnetwork.com"),
		colly.CacheDir(siteCacheDir),
		colly.UserAgent(userAgent),
		colly.MaxDepth(5),
	)

	sceneCollector := colly.NewCollector(
		colly.AllowedDomains("www.czechvrnetwork.com"),
		colly.CacheDir(sceneCacheDir),
		colly.UserAgent(userAgent),
	)

	siteCollector.OnRequest(func(r *colly.Request) {
		log.Println("visiting", r.URL.String())
	})

	sceneCollector.OnRequest(func(r *colly.Request) {
		log.Println("visiting", r.URL.String())
	})

	sceneCollector.OnHTML(`html`, func(e *colly.HTMLElement) {
		sc := ScrapedScene{}
		sc.SceneType = "VR"
		sc.Studio = "CzechVR"
		sc.Site = "Czech VR"
		sc.HomepageURL = strings.Split(e.Request.URL.String(), "?")[0]

		// Title
		e.ForEach(`div.nazev h2 span.desktop`, func(id int, e *colly.HTMLElement) {
			sc.Site = strings.TrimSpace(e.Text)
		})

		e.ForEach(`div.nazev h2`, func(id int, e *colly.HTMLElement) {
			fullTitle := strings.TrimSpace(e.Text)
			sc.Title = strings.Split(fullTitle, " - ")[1]
			tmp := strings.Split(strings.Split(fullTitle, " - ")[0], " ")
			sc.SiteID = tmp[len(tmp)-1]
		})
		sc.SceneID = slugify.Slugify(sc.Site) + "-" + sc.SiteID

		e.ForEach(`dl8-video`, func(id int, e *colly.HTMLElement) {
			sc.Covers = append(sc.Covers, e.Request.AbsoluteURL(e.Attr("poster")))
		})

		// Gallery
		e.ForEach(`div.galerka a`, func(id int, e *colly.HTMLElement) {
			base := e.Request.AbsoluteURL(e.Attr("href"))
			base = strings.Split(base, "?")[0]
			sc.Gallery = append(sc.Gallery, base)
		})

		// Synopsis
		e.ForEach(`div.textDetail`, func(id int, e *colly.HTMLElement) {
			sc.Synopsis = strings.TrimSpace(e.Text)
		})

		// Tags
		e.ForEach(`div.tagyall div.tag`, func(id int, e *colly.HTMLElement) {
			sc.Tags = append(sc.Tags, strings.TrimSpace(e.Text))
		})

		// Cast
		e.ForEach(`div.featuring a`, func(id int, e *colly.HTMLElement) {
			sc.Cast = append(sc.Cast, strings.TrimSpace(e.Text))
		})

		// Date
		e.ForEach(`div.nazev div.datum`, func(id int, e *colly.HTMLElement) {
			tmpDate, _ := goment.New(e.Text, "MMM DD, YYYY")
			sc.Released = tmpDate.Format("YYYY-MM-DD")
		})

		// Duration
		e.ForEach(`div.nazev div.casDetail span.desktop`, func(id int, e *colly.HTMLElement) {
			tmpDuration, err := strconv.Atoi(strings.Split(e.Text, ":")[0])

			if err == nil {
				sc.Duration = tmpDuration
			}
		})

		// Filenames
		e.ForEach(`div.download a.trailer`, func(id int, e *colly.HTMLElement) {
			if (id == 1) {
				tmp := strings.Split(e.Attr("href"), "/")
				parts := strings.Split(tmp[len(tmp)-1], "-")
				base := parts[0] + "-" + parts[1] + "-" + parts[2]

				filenames := []string{
					"1920x960-30fps-smartphone_lq",
					"1920x960-60fps-psvr_med_h264_180_sbs",
					"1920x960-60fps-smartphone_hq",
					"2160x1080-60fps-smartphone_hq",
					"2880x1440-60fps-gearvr_lq_h264-180x180_3dh",
					"2880x1440-60fps-gearvr_med_h264-180x180_3dh",
					"2880x1440-60fps-psvr_hq_h264_180_sbs",
					"3840x1920-60fps-gearvr_hq_h264-180x180_3dh",
					"3840x1920-60fps-gearvr_med_h264-180x180_3dh",
					"3840x1920-60fps-oculusrift_hq_h264",
					"3840x1920-60fps-oculusrift_med_h264",
					"5400x2700-60fps-gearvr_hq_h265-180x180_3dh",
					"5400x2700-60fps-oculusrift_hq_h265",
				}

				for i := range filenames {
					filenames[i] = base + "-" + filenames[i] + ".mp4"
				}

				sc.Filenames = filenames
			}
		})

		*out = append(*out, sc)
	})

	siteCollector.OnHTML(`div#StrankovaniDesktop span.stred a,div#StrankovaniDesktopHome span.stred a`, func(e *colly.HTMLElement) {
		pageURL := e.Request.AbsoluteURL(e.Attr("href"))
		siteCollector.Visit(pageURL)
	})

	siteCollector.OnHTML(`div.postTag div.foto a`, func(e *colly.HTMLElement) {
		sceneURL := e.Request.AbsoluteURL(e.Attr("href"))

		// If scene exist in database, there's no need to scrape
		if !funk.ContainsString(knownScenes, sceneURL) {
			sceneCollector.Visit(sceneURL)
		}
	})

	siteCollector.Visit("https://www.czechvrnetwork.com/vr-porn-videos?next=1")

	return nil
}
