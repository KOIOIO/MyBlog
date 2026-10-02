package website

import (
	"context"
	"testing"

	"server/config"
	"server/internal/domain/website"
)

type stubImages struct {
	urls      []string
	changed   []string
	changeCat string
	inited    []string
}

func (s *stubImages) CarouselURLs(context.Context) ([]string, error) { return s.urls, nil }
func (s *stubImages) ChangeCategory(_ context.Context, urls []string, category string) error {
	s.changed = append(s.changed, urls...)
	s.changeCat = category
	return nil
}
func (s *stubImages) InitCategory(_ context.Context, urls []string) error {
	s.inited = append(s.inited, urls...)
	return nil
}

var _ website.WebsiteImagePort = (*stubImages)(nil)

type stubFooters struct {
	list    []*website.FooterLink
	saved   []*website.FooterLink
	deleted []*website.FooterLink
}

func (s *stubFooters) List(context.Context) ([]*website.FooterLink, error) { return s.list, nil }
func (s *stubFooters) Save(_ context.Context, l *website.FooterLink) error {
	s.saved = append(s.saved, l)
	return nil
}
func (s *stubFooters) Delete(_ context.Context, l *website.FooterLink) error {
	s.deleted = append(s.deleted, l)
	return nil
}

var _ website.FooterLinkRepository = (*stubFooters)(nil)

type stubNews struct {
	data website.HotSearchData
}

func (s *stubNews) GetHotSearchData(_ context.Context, source string) (website.HotSearchData, error) {
	s.data.Source = source
	return s.data, nil
}
func (s *stubNews) WarmAll(context.Context) error { return nil }

var _ website.HotSearchProvider = (*stubNews)(nil)

type stubCal struct {
	data website.Calendar
}

func (s *stubCal) GetCalendarByDate(_ context.Context, dateStr string) (website.Calendar, error) {
	s.data.Date = dateStr
	return s.data, nil
}

var _ website.CalendarProvider = (*stubCal)(nil)

func newTestService(imgs *stubImages, footers *stubFooters, news *stubNews, cal *stubCal) *Service {
	return NewService(imgs, footers, news, cal, &config.Config{})
}

func TestCarouselAndFooterLink(t *testing.T) {
	imgs := &stubImages{urls: []string{"/a.png", "/b.png"}}
	footers := &stubFooters{list: []*website.FooterLink{{ID: 1, Name: "n"}}}
	svc := newTestService(imgs, footers, &stubNews{}, &stubCal{})

	urls, err := svc.Carousel(context.Background())
	if err != nil || len(urls) != 2 {
		t.Fatalf("carousel: %v %v", urls, err)
	}
	links, err := svc.FooterLink(context.Background())
	if err != nil || len(links) != 1 || links[0].Name != "n" {
		t.Fatalf("footer: %v %v", links, err)
	}

	if err := svc.AddCarousel(context.Background(), "/c.png"); err != nil {
		t.Fatal(err)
	}
	if len(imgs.changed) != 1 || imgs.changed[0] != "/c.png" || imgs.changeCat != "背景" {
		t.Fatalf("add carousel: %v %q", imgs.changed, imgs.changeCat)
	}
	if err := svc.CancelCarousel(context.Background(), "/c.png"); err != nil {
		t.Fatal(err)
	}
	if len(imgs.inited) != 1 || imgs.inited[0] != "/c.png" {
		t.Fatalf("cancel carousel: %v", imgs.inited)
	}

	fl := &website.FooterLink{Logo: "/l.png", Name: "x"}
	if err := svc.CreateFooterLink(context.Background(), fl); err != nil {
		t.Fatal(err)
	}
	if len(footers.saved) != 1 || footers.saved[0].Name != "x" {
		t.Fatalf("create footer: %v", footers.saved)
	}
	if err := svc.DeleteFooterLink(context.Background(), fl); err != nil {
		t.Fatal(err)
	}
	if len(footers.deleted) != 1 {
		t.Fatalf("delete footer: %v", footers.deleted)
	}
}

func TestNewsAndCalendar(t *testing.T) {
	news := &stubNews{data: website.HotSearchData{UpdateTime: "t"}}
	cal := &stubCal{data: website.Calendar{LunarDate: "八月廿二"}}
	svc := newTestService(&stubImages{}, &stubFooters{}, news, cal)

	d, err := svc.News(context.Background(), "baidu")
	if err != nil || d.Source != "baidu" {
		t.Fatalf("news: %+v %v", d, err)
	}
	c, err := svc.Calendar(context.Background(), "2026/1002")
	if err != nil || c.LunarDate != "八月廿二" {
		t.Fatalf("calendar: %+v %v", c, err)
	}
}
